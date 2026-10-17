package web

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/i18n"
	tlmail "github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/store"
)

// inviteTTL is how long the link in an invitation mail sets the password.
const inviteTTL = 7 * 24 * time.Hour

// InviteResult says what happened to one invited email address.
type InviteResult struct {
	Email  string `json:"email"`
	Status string `json:"status" enum:"invited,added,member,invalid" doc:"invited: new account, mail with the link to set the password; added: existing account joined the club; member: was a member already; invalid: not an email address"`
}

// registerInvites registers the email invitation endpoint of a club.
func (s *Server) registerInvites(api huma.API) {
	type inviteIn struct {
		TenantPath
		Body struct {
			Emails []string   `json:"emails" minItems:"1" maxItems:"300" doc:"One address per entry; the UI splits a comma or newline separated list. Mails leave at most four per second"`
			Role   store.Role `json:"role" enum:"admin,staff"`
		}
	}
	type inviteOut struct{ Body []InviteResult }
	huma.Register(api, huma.Operation{
		OperationID: "invite-users", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/invitations", Tags: []string{"users"},
		Summary: "Invite people by email: existing accounts join the club, new ones get a mail with a link to set their password",
	}, func(ctx context.Context, in *inviteIn) (*inviteOut, error) {
		by, _, err := s.access(ctx, in.TenantID, true)
		if err != nil {
			return nil, err
		}
		tenant, err := s.Store.GetTenant(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		out := &inviteOut{Body: []InviteResult{}}
		seen := map[string]bool{}
		for _, raw := range in.Body.Emails {
			email := strings.ToLower(strings.TrimSpace(raw))
			if email == "" || seen[email] {
				continue
			}
			seen[email] = true
			status, err := s.invite(ctx, in.TenantID, tenant.Name, email, in.Body.Role)
			if err != nil {
				return nil, err
			}
			out.Body = append(out.Body, InviteResult{Email: email, Status: status})
		}
		s.Log.Info("users invited", "tenant", in.TenantID, "count", len(out.Body), "by", by.Name)
		return out, nil
	})
}

// invite adds one email address to a club: an existing account becomes a member and is told by mail,
// a new one is created without a usable password and gets the link to set it.
func (s *Server) invite(ctx context.Context, tenantID int64, club, email string, role store.Role) (string, error) {
	if a, err := mail.ParseAddress(email); err != nil || a.Address != email {
		return "invalid", nil
	}
	l := mailLang(ctx, "")
	u, err := s.Store.UserByEmail(ctx, email)
	switch {
	case err == nil:
		if _, err := s.Store.MemberRole(ctx, tenantID, u.ID); err == nil {
			return "member", nil
		} else if !errors.Is(err, store.ErrNotFound) {
			return "", err
		}
		if err := s.Store.SetMembership(ctx, tenantID, u.ID, role); err != nil {
			return "", err
		}
		s.mail(tlmail.Message{To: email, Subject: i18n.Tf(l, "%s added you on Standblatt", club), Body: strings.Join([]string{
			greeting(l, s.firstname(ctx, u.ID)), "",
			i18n.Tf(l, "%s gave your account access to its events on Standblatt. Log in as usual:", club), "",
			s.link(ctx, "/login"),
		}, "\n")})
		return "added", nil
	case !errors.Is(err, store.ErrNotFound):
		return "", err
	}
	// a random password nobody knows; the invitation link sets the real one
	hash, err := HashPassword(randomToken())
	if err != nil {
		return "", err
	}
	id, err := s.Store.CreateAccount(ctx, email, hash, reqInfo(ctx).lang)
	if errors.Is(err, store.ErrExists) {
		// the address is somebody's user name but not their email: treat it like an invalid entry
		return "invalid", nil
	}
	if err != nil {
		return "", err
	}
	if err := s.Store.SetMembership(ctx, tenantID, id, role); err != nil {
		return "", err
	}
	token := randomToken()
	if err := s.Store.CreateUserToken(ctx, sessionID(token), id, "reset", time.Now().Add(inviteTTL)); err != nil {
		return "", err
	}
	s.mail(tlmail.Message{To: email, Subject: i18n.Tf(l, "%s invites you to Standblatt", club), Body: strings.Join([]string{
		i18n.T(l, "Hello,"), "",
		i18n.Tf(l, "%s uses Standblatt for its shooting events and created an account for you. Set your password with this link:", club), "",
		s.link(ctx, "/reset-password?token="+token), "",
		i18n.T(l, "The link is valid for seven days. Afterwards you log in with this email address and your password."),
	}, "\n")})
	return "invited", nil
}
