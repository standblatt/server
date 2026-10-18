package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

// TestRoles checks the role API and that a member can do exactly what the permissions of their role allow.
func TestRoles(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Roles")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	memberID := e.user("mia", tid, store.RoleMember)
	boss, helper, mia := e.login("boss"), e.login("helper"), e.login("mia")
	base := fmt.Sprintf("/api/tenants/%d", tid)
	ev := e.event(tid, "Annual")

	code, b := e.do(boss, "GET", base+"/roles", nil)
	expect(t, "list roles", code, 200, b)
	var roles []store.Role
	json.Unmarshal(b, &roles)
	if len(roles) != 4 || roles[0].Key != "admin" || len(roles[0].Permissions) != len(permissions.All) || roles[3].Members != 1 {
		t.Fatalf("default roles: %+v", roles)
	}
	code, b = e.do(helper, "GET", base+"/roles", nil)
	expect(t, "helpers manage no roles", code, 403, b)

	code, b = e.do(mia, "GET", base+"/events", nil)
	expect(t, "member sees events", code, 200, b)
	code, b = e.do(mia, "GET", fmt.Sprintf("%s/events/%d/registrations", base, ev), nil)
	expect(t, "member sees no registrations", code, 403, b)
	code, b = e.do(mia, "GET", base+"/shooters?q=", nil)
	expect(t, "member sees no shooters", code, 403, b)
	code, b = e.do(helper, "GET", base+"/shooters?q=", nil)
	expect(t, "helper role searches shooters", code, 200, b)

	code, b = e.do(boss, "POST", base+"/roles", store.RoleInput{Name: " ", Permissions: []string{}})
	expect(t, "role without name", code, 422, b)
	code, b = e.do(boss, "POST", base+"/roles", store.RoleInput{Name: "Secretary", Permissions: []string{permissions.EventsView, permissions.ShootersView}})
	expect(t, "create role", code, 201, b)
	var secretary store.Role
	json.Unmarshal(b, &secretary)
	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/users/%d", base, memberID), map[string]int64{"roleId": secretary.ID})
	expect(t, "give mia the role", code, 204, b)
	code, b = e.do(mia, "GET", base+"/shooters?q=", nil)
	expect(t, "secretary sees shooters", code, 200, b)
	code, b = e.do(mia, "POST", base+"/shooters", map[string]any{"firstname": "A", "lastname": "B"})
	expect(t, "secretary edits no shooters", code, 403, b)
	code, b = e.do(mia, "GET", "/api/me", nil)
	expect(t, "me", code, 200, b)
	var me web.Me
	json.Unmarshal(b, &me)
	if me.Tenants[0].RoleName != "Secretary" || fmt.Sprint(me.Tenants[0].Permissions) != "[events.view shooters.view]" {
		t.Fatalf("me: %+v", me.Tenants[0])
	}

	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/roles/%d", base, secretary.ID), nil)
	expect(t, "delete a role with members", code, 409, b)
	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/roles/%d", base, roles[0].ID), store.RoleInput{Name: "Chief", Permissions: []string{}})
	expect(t, "rename admin", code, 204, b)
	code, b = e.do(boss, "GET", base+"/roles", nil)
	expect(t, "list roles", code, 200, b)
	json.Unmarshal(b, &roles)
	if roles[0].Name != "Chief" || len(roles[0].Permissions) != len(permissions.All) {
		t.Fatalf("admin keeps every permission: %+v", roles[0])
	}
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/roles/%d", base, roles[0].ID), nil)
	expect(t, "delete the admin role", code, 409, b)
	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/roles/%d", base, secretary.ID), store.RoleInput{Name: "", Permissions: []string{}})
	expect(t, "own role without name", code, 422, b)

	other := e.tenant("Other")
	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/users/%d", base, memberID), map[string]int64{"roleId": e.role(other, "staff")})
	expect(t, "role of another club", code, 404, b)

	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/users/%d", base, memberID), map[string]int64{"roleId": e.role(tid, "member")})
	expect(t, "mia back to member", code, 204, b)
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/roles/%d", base, secretary.ID), nil)
	expect(t, "delete a role without members", code, 204, b)
	code, b = e.do(boss, "GET", base+"/roles", nil)
	expect(t, "list roles", code, 200, b)
	if json.Unmarshal(b, &roles); len(roles) != 4 {
		t.Fatalf("roles after delete: %+v", roles)
	}
}

// TestManagersStayWithinTheirPermissions checks that a role with club.manage but without the admin role cannot make,
// change or remove admins, cannot take a stronger role and cannot give permissions it does not have.
func TestManagersStayWithinTheirPermissions(t *testing.T) {
	e, box := personalEnv(t, nil)
	ctx := context.Background()
	tid := e.tenant("Managed")
	bossID := e.user("boss", tid, store.RoleAdmin)
	helperID := e.user("helper", tid, store.RoleStaff)
	memberID := e.user("mia", tid, store.RoleMember)
	boss := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d", tid)
	code, b := e.do(boss, "POST", base+"/roles", store.RoleInput{Name: "Manager", Permissions: []string{permissions.ClubManage, permissions.EventsView}})
	expect(t, "manager role", code, 201, b)
	var manager store.Role
	json.Unmarshal(b, &manager)
	managerID := e.user("max", 0, "")
	if err := e.st.SetMembership(ctx, tid, managerID, manager.ID); err != nil {
		t.Fatal(err)
	}
	max := e.login("max")
	admin, board, staff, member := e.role(tid, store.RoleAdmin), e.role(tid, store.RoleBoard), e.role(tid, store.RoleStaff), e.role(tid, store.RoleMember)
	user := func(id int64) string { return fmt.Sprintf("%s/users/%d", base, id) }
	role := func(id int64) string { return fmt.Sprintf("%s/roles/%d", base, id) }

	for _, c := range []struct {
		what, method, path string
		body               any
	}{
		{"demote an admin", "PUT", user(bossID), map[string]int64{"roleId": member}},
		{"remove an admin", "DELETE", user(bossID), nil},
		{"make someone admin", "PUT", user(memberID), map[string]int64{"roleId": admin}},
		{"take the admin role", "PUT", user(managerID), map[string]int64{"roleId": admin}},
		{"take the board role", "PUT", user(managerID), map[string]int64{"roleId": board}},
		{"change a stronger member", "PUT", user(helperID), map[string]int64{"roleId": member}},
		{"remove a stronger member", "DELETE", user(helperID), nil},
		{"add an admin", "POST", base + "/users", map[string]any{"name": "newadmin", "roleId": admin, "password": "long-enough-1"}},
		{"invite an admin", "POST", base + "/invitations", map[string]any{"emails": []string{"a@example.org"}, "roleId": admin}},
		{"make a stronger role", "POST", base + "/roles", store.RoleInput{Name: "Root", Permissions: []string{permissions.RangesManage}}},
		{"extend the own role", "PUT", role(manager.ID), store.RoleInput{Name: "Manager", Permissions: []string{permissions.ClubManage, permissions.EventsView, permissions.EventsConfigure}}},
		{"change a stronger role", "PUT", role(staff), store.RoleInput{Permissions: []string{permissions.EventsView}}},
		{"delete a stronger role", "DELETE", role(board), nil},
		{"rename the admin role", "PUT", role(admin), store.RoleInput{Name: "Chief", Permissions: []string{}}},
	} {
		for user, role := range map[int64]int64{bossID: admin, helperID: staff, memberID: member, managerID: manager.ID} {
			if err := e.st.SetMembership(ctx, tid, user, role); err != nil {
				t.Fatal(err)
			}
		}
		if code, b := e.do(max, c.method, c.path, c.body); code != 403 {
			t.Errorf("%s: got %d, want 403: %s", c.what, code, b)
		}
	}

	code, b = e.do(boss, "POST", base+"/invitations", map[string]any{"emails": []string{"b@example.org"}, "roleId": admin})
	expect(t, "admin invites an admin", code, 200, b)
	box.next(t)
	code, b = e.do(max, "GET", base+"/invitations", nil)
	expect(t, "list invitations", code, 200, b)
	var open []store.Invitation
	json.Unmarshal(b, &open)
	code, b = e.do(max, "DELETE", fmt.Sprintf("%s/invitations/%d", base, open[0].ID), nil)
	expect(t, "withdraw an admin invitation", code, 403, b)

	code, b = e.do(max, "PUT", user(memberID), map[string]int64{"roleId": member})
	expect(t, "manage a member within the own permissions", code, 204, b)
	code, b = e.do(max, "POST", base+"/roles", store.RoleInput{Name: "Viewer", Permissions: []string{permissions.EventsView}})
	expect(t, "make a weaker role", code, 201, b)
	code, b = e.do(max, "PUT", user(bossID), map[string]int64{"roleId": admin})
	expect(t, "the admin itself is out of reach", code, 403, b)
	code, b = e.do(boss, "PUT", user(helperID), map[string]int64{"roleId": admin})
	expect(t, "admins make admins", code, 204, b)
	code, b = e.do(boss, "DELETE", user(helperID), nil)
	expect(t, "admins remove admins", code, 204, b)
}
