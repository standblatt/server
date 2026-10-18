package web_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

var updatePermissions = flag.Bool("update-permissions", false, "rewrite testdata/permissions.txt")

// TestPermissionTable keeps every operation of a club on known permissions and the whole table in
// testdata/permissions.txt, so any change to who may call what shows up in review. After a deliberate change run
// go test ./internal/web -run TestPermissionTable -update-permissions.
func TestPermissionTable(t *testing.T) {
	var b strings.Builder
	for _, op := range web.ClubOperationsForTest() {
		if len(op.Permissions) == 0 {
			t.Errorf("%s %s (%s) requires no permission", op.Method, op.Path, op.ID)
		}
		for _, p := range op.Permissions {
			if !permissions.Known(p) {
				t.Errorf("%s: unknown permission %q", op.ID, p)
			}
		}
		fmt.Fprintf(&b, "%-6s %-78s %s\n", op.Method, op.Path, strings.Join(op.Permissions, " | "))
	}
	const file = "testdata/permissions.txt"
	if *updatePermissions {
		if err := os.WriteFile(file, []byte(b.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(want) != b.String() {
		t.Fatalf("the permissions of the operations changed; review and run with -update-permissions:\n%s", b.String())
	}
}

var pathParam = regexp.MustCompile(`\{(\w+)\}`)

// clubPath fills the path of an operation with the tenant and IDs nobody has, so no request changes real data.
func clubPath(path string, tenantID int64) string {
	return pathParam.ReplaceAllStringFunc(path, func(p string) string {
		switch {
		case p == "{tenantId}":
			return fmt.Sprint(tenantID)
		case strings.HasSuffix(p, "Id}"):
			return "999999"
		}
		return "x"
	})
}

// TestPermissionLeaks calls every operation of a club as a member whose role lacks its permissions (403 before any
// validation) and with one of them (no 403 for the role), and as each helper link role outside its operations.
func TestPermissionLeaks(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Leaks")
	e.user("boss", tid, store.RoleAdmin)
	e.user("probe", tid, store.RoleMember)
	probe := e.login("probe")
	role := e.role(tid, store.RoleMember)
	set := func(perms []string) {
		t.Helper()
		if err := e.st.UpdateRole(ctx, tid, role, store.RoleInput{Permissions: perms}); err != nil {
			t.Fatal(err)
		}
	}
	call := func(c func(method, path string) (int, []byte), op web.OperationForTest) (int, []byte) {
		return c(op.Method, clubPath(op.Path, tid))
	}
	asProbe := func(method, path string) (int, []byte) {
		var body any
		if method != "GET" && method != "DELETE" {
			body = map[string]any{}
		}
		return e.do(probe, method, path, body)
	}

	for _, op := range web.ClubOperationsForTest() {
		others := slices.DeleteFunc(permissions.Keys(), func(p string) bool { return slices.Contains(op.Permissions, p) })
		for _, perms := range [][]string{{}, others} {
			set(perms)
			if code, b := call(asProbe, op); code != 403 || !strings.Contains(string(b), "Not allowed for your role") {
				t.Errorf("%s %s with %v: %d %s", op.Method, op.Path, perms, code, b)
			}
		}
		for _, p := range op.Permissions {
			set([]string{p})
			if code, b := call(asProbe, op); strings.Contains(string(b), "Not allowed for your role") {
				t.Errorf("%s %s refused with %s: %d", op.Method, op.Path, p, code)
			}
		}
	}

	ev := e.event(tid, "Leak test")
	event, _ := e.st.GetEvent(ctx, tid, ev)
	set(permissions.Keys())
	for _, r := range []store.HelperRole{store.HelperScorer, store.HelperResults, store.HelperAmmunition} {
		status, b := e.do(e.login("boss"), "POST", fmt.Sprintf("/api/tenants/%d/events/%d/helper-codes", tid, ev),
			map[string]any{"name": string(r), "role": r, "rangeId": event.Ranges[0].ID, "lanes": "1-2"})
		expect(t, "helper link", status, 201, b)
		var link struct{ Secret string }
		json.Unmarshal(b, &link)
		code := link.Secret
		asHelper := func(method, path string) (int, []byte) { return e.helper(method, path, code, nil) }
		for _, op := range web.ClubOperationsForTest() {
			if web.HelperMayForTest(r, op.ID) {
				continue
			}
			if status, b := call(asHelper, op); status != 403 {
				t.Errorf("%s link may call %s %s: %d %s", r, op.Method, op.Path, status, b)
			}
		}
	}
}
