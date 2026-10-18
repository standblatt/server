package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

var inviteRe = regexp.MustCompile(`/invite/(inv_[A-Za-z0-9_-]+)`)

// inviteToken returns the token of the next mail, which must be an invitation to an address.
func (m mailbox) inviteToken(t *testing.T, to string) string {
	t.Helper()
	msg := m.next(t)
	match := inviteRe.FindStringSubmatch(msg.Body)
	if msg.To != to || match == nil {
		t.Fatalf("invitation to %s: %+v", to, msg)
	}
	return match[1]
}

func TestInviteUsers(t *testing.T) {
	e, box := personalEnv(t, nil)
	ctx := context.Background()
	tid := e.tenant("PS Example")
	bossID := e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	lea := e.signup("lea@example.org")
	box.next(t) // Lea's own verification mail
	path := fmt.Sprintf("/api/tenants/%d/invitations", tid)
	staffRole := e.role(tid, "staff")

	code, b := e.do(staff, "POST", path, map[string]any{"emails": []string{"x@example.org"}, "roleId": staffRole})
	expect(t, "staff invites", code, 403, b)
	code, b = e.do(admin, "POST", path, map[string]any{"emails": []string{"x@example.org"}, "roleId": e.role(e.tenant("Other"), "staff")})
	expect(t, "role of another club", code, 404, b)

	code, b = e.do(admin, "POST", path, map[string]any{"emails": []string{"New.Shooter@Example.org", "lea@example.org", "boss", " new.shooter@example.org "}, "roleId": staffRole})
	expect(t, "invite", code, 200, b)
	var results []web.InviteResult
	json.Unmarshal(b, &results)
	got := map[string]string{}
	for _, r := range results {
		got[r.Email] = r.Status
	}
	if len(results) != 3 || got["new.shooter@example.org"] != "invited" || got["lea@example.org"] != "invited" || got["boss"] != "invalid" {
		t.Fatalf("results: %s", b)
	}
	newcomer := box.inviteToken(t, "new.shooter@example.org")
	leaToken := box.inviteToken(t, "lea@example.org")
	if members, _ := e.st.Members(ctx, tid); len(members) != 2 {
		t.Fatalf("nobody joins before accepting: %+v", members)
	}
	code, b = e.do(admin, "GET", path, nil)
	expect(t, "list invitations", code, 200, b)
	var open []store.Invitation
	json.Unmarshal(b, &open)
	if len(open) != 2 || open[0].Expired || open[0].RoleID != staffRole {
		t.Fatalf("open invitations: %s", b)
	}

	anon := e.client()
	code, b = e.do(anon, "GET", "/api/invitations/"+newcomer, nil)
	expect(t, "view invitation", code, 200, b)
	var view web.InvitationView
	json.Unmarshal(b, &view)
	if view.Club != "PS Example" || view.Email != "new.shooter@example.org" || view.RoleKey != "staff" || view.Account {
		t.Fatalf("view: %+v", view)
	}
	code, b = e.do(anon, "POST", "/api/invitations/"+newcomer+"/accept", map[string]string{"password": "short"})
	expect(t, "too short a password", code, 422, b)
	code, b = e.do(anon, "POST", "/api/invitations/"+newcomer+"/accept", map[string]string{"password": "a-chosen-password-1"})
	expect(t, "accept with a new account", code, 200, b)
	code, b = e.do(anon, "GET", "/api/me", nil)
	expect(t, "logged in", code, 200, b)
	if !strings.Contains(string(b), `"tenantName":"PS Example"`) || !strings.Contains(string(b), `"emailVerified":true`) || !strings.Contains(string(b), `"roleKey":"staff"`) {
		t.Fatalf("new member: %s", b)
	}
	code, b = e.do(e.client(), "GET", "/api/invitations/"+newcomer, nil)
	expect(t, "accepted link", code, 410, b)

	code, b = e.do(e.client(), "GET", "/api/invitations/"+leaToken, nil)
	expect(t, "view as existing account", code, 200, b)
	if !strings.Contains(string(b), `"account":true`) {
		t.Fatalf("lea has an account: %s", b)
	}
	code, b = e.do(e.client(), "POST", "/api/invitations/"+leaToken+"/accept", map[string]string{"password": "a-chosen-password-1"})
	expect(t, "existing account without login", code, 409, b)
	code, b = e.do(lea, "POST", "/api/invitations/"+leaToken+"/accept", nil)
	expect(t, "accept as lea", code, 200, b)
	code, b = e.do(admin, "POST", path, map[string]any{"emails": []string{"lea@example.org"}, "roleId": staffRole})
	expect(t, "invite a member", code, 200, b)
	if !strings.Contains(string(b), `"status":"member"`) {
		t.Fatalf("second invite: %s", b)
	}

	code, b = e.do(admin, "POST", path, map[string]any{"emails": []string{"x@example.org"}, "roleId": staffRole})
	expect(t, "invite x", code, 200, b)
	first := box.inviteToken(t, "x@example.org")
	code, b = e.do(admin, "GET", path, nil)
	expect(t, "list", code, 200, b)
	json.Unmarshal(b, &open)
	if len(open) != 1 {
		t.Fatalf("accepted invitations are no longer open: %s", b)
	}
	code, b = e.do(admin, "POST", fmt.Sprintf("%s/%d/resend", path, open[0].ID), nil)
	expect(t, "resend", code, 204, b)
	second := box.inviteToken(t, "x@example.org")
	code, b = e.do(e.client(), "GET", "/api/invitations/"+first, nil)
	expect(t, "old link after resending", code, 404, b)
	code, b = e.do(e.client(), "GET", "/api/invitations/"+second, nil)
	expect(t, "new link", code, 200, b)
	code, b = e.do(admin, "DELETE", fmt.Sprintf("%s/%d", path, open[0].ID), nil)
	expect(t, "withdraw", code, 204, b)
	code, b = e.do(e.client(), "GET", "/api/invitations/"+second, nil)
	expect(t, "withdrawn link", code, 404, b)

	if err := e.st.SaveInvitation(ctx, tid, "late@example.org", staffRole, bossID, web.SessionIDForTest("inv_late"), time.Now().Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(e.client(), "GET", "/api/invitations/inv_late", nil)
	expect(t, "expired link", code, 410, b)
}
