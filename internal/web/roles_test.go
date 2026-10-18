package web_test

import (
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
}
