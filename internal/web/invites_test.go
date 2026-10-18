package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestInviteUsers(t *testing.T) {
	e, box := personalEnv(t, nil)
	ctx := context.Background()
	tid := e.tenant("PS Example")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	e.signup("lea@example.org")
	box.next(t) // Lea's own verification mail
	path := fmt.Sprintf("/api/tenants/%d/invitations", tid)

	code, b := e.do(staff, "POST", path, map[string]any{"emails": []string{"x@example.org"}, "roleId": e.role(tid, "staff")})
	expect(t, "staff invites", code, 403, b)

	code, b = e.do(admin, "POST", path, map[string]any{"emails": []string{"New.Shooter@Example.org", "lea@example.org", "boss", " new.shooter@example.org "}, "roleId": e.role(tid, "staff")})
	expect(t, "invite", code, 200, b)
	var results []web.InviteResult
	json.Unmarshal(b, &results)
	got := map[string]string{}
	for _, r := range results {
		got[r.Email] = r.Status
	}
	if len(results) != 3 || got["new.shooter@example.org"] != "invited" || got["lea@example.org"] != "added" || got["boss"] != "invalid" {
		t.Fatalf("results: %s", b)
	}
	code, b = e.do(admin, "POST", path, map[string]any{"emails": []string{"lea@example.org"}, "roleId": e.role(tid, "admin")})
	expect(t, "invite again", code, 200, b)
	if !strings.Contains(string(b), `"status":"member"`) {
		t.Fatalf("second invite: %s", b)
	}

	// Lea was told, the newcomer got the link to set the password.
	var inviteToken string
	for i := 0; i < 2; i++ {
		m := box.next(t)
		switch m.To {
		case "lea@example.org":
			if !strings.Contains(m.Subject, "PS Example added you") {
				t.Fatalf("lea mail: %q", m.Subject)
			}
		case "new.shooter@example.org":
			if match := tokenRe.FindStringSubmatch(m.Body); match != nil {
				inviteToken = match[1]
			}
		default:
			t.Fatalf("unexpected mail to %s", m.To)
		}
	}
	if inviteToken == "" {
		t.Fatal("no invitation mail")
	}
	code, b = e.do(e.client(), "POST", "/api/auth/reset", map[string]string{"token": inviteToken, "password": "a-chosen-password-1"})
	expect(t, "set password from invitation", code, 204, b)
	c := e.client()
	code, b = e.do(c, "POST", "/api/auth/login", map[string]string{"name": "new.shooter@example.org", "password": "a-chosen-password-1"})
	expect(t, "login", code, 200, b)
	if !strings.Contains(string(b), `"tenantName":"PS Example"`) || !strings.Contains(string(b), `"emailVerified":true`) {
		t.Fatalf("new member: %s", b)
	}
	members, _ := e.st.Members(ctx, tid)
	for _, m := range members {
		if m.Name == "new.shooter@example.org" && m.Email != "new.shooter@example.org" {
			t.Fatalf("member without email: %+v", m)
		}
	}
}
