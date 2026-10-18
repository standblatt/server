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
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
)

// inviteTTL is how long the link in an invitation mail works.
const inviteTTL = 7 * 24 * time.Hour

// invitePrefix marks the secret of an invitation link.
const invitePrefix = "inv_"

// InviteResult says what happened to one invited email address.
type InviteResult struct {
	Email  string `json:"email"`
	Status string `json:"status" enum:"invited,member,invalid" doc:"invited: mail with the link sent (or sent again); member: is a member already; invalid: not an email address"`
}

// InvitationView is what the link of an invitation shows before it is accepted.
type InvitationView struct {
	Club     string `json:"club"`
	Email    string `json:"email"`
	RoleKey  string `json:"roleKey,omitempty" required:"false" enum:"admin,board,staff,member"`
	RoleName string `json:"roleName" doc:"Empty for a default role that keeps its translated name"`
	Account  bool   `json:"account" doc:"An account with this address exists; its owner logs in to accept"`
}

// InvitationPath identifies an invitation of a tenant in the URL.
type InvitationPath struct {
	TenantID     int64 `path:"tenantId"`
	InvitationID int64 `path:"invitationId"`
}

// InvitationToken is the secret of an invitation link in the URL.
type InvitationToken struct {
	Token string `path:"token"`
}

// registerInvites registers the invitations of a club and the public endpoints of their links.
func (s *Server) registerInvites(api huma.API) {
	type inviteIn struct {
		TenantPath
		Body struct {
			Emails []string `json:"emails" minItems:"1" maxItems:"300" doc:"One address per entry; the UI splits a comma or newline separated list. Mails leave at most four per second"`
			RoleID int64    `json:"roleId"`
		}
	}
	type inviteOut struct{ Body []InviteResult }
	huma.Register(api, huma.Operation{
		OperationID: "invite-users", Metadata: needs("", permissions.ClubManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/invitations", Tags: []string{"users"},
		Summary: "Invite people by email; they become members when they accept through the link in the mail",
	}, func(ctx context.Context, in *inviteIn) (*inviteOut, error) {
		by, _, err := s.access(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		tenant, err := s.Store.GetTenant(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		if err := s.roleWithin(ctx, in.TenantID, in.Body.RoleID); err != nil {
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
			status, err := s.invite(ctx, in.TenantID, tenant.Name, email, in.Body.RoleID, by.ID)
			if err != nil {
				return nil, err
			}
			out.Body = append(out.Body, InviteResult{Email: email, Status: status})
		}
		s.Log.Info("users invited", "tenant", in.TenantID, "count", len(out.Body), "by", by.Name)
		return out, nil
	})

	type listOut struct{ Body []store.Invitation }
	huma.Register(api, huma.Operation{
		OperationID: "list-invitations", Metadata: needs("", permissions.ClubManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/invitations", Tags: []string{"users"},
	}, func(ctx context.Context, in *TenantPath) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListInvitations(ctx, in.TenantID)
		return &listOut{Body: list}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "resend-invitation", Metadata: needs("", permissions.ClubManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/invitations/{invitationId}/resend", Tags: []string{"users"},
		Summary: "Send an invitation again with a new link valid for seven days", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *InvitationPath) (*struct{}, error) {
		by, _, err := s.access(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		inv, err := s.Store.GetInvitation(ctx, in.TenantID, in.InvitationID)
		if err != nil {
			return nil, notFound(err, "Invitation not found")
		}
		if err := s.roleWithin(ctx, in.TenantID, inv.RoleID); err != nil {
			return nil, err
		}
		tenant, err := s.Store.GetTenant(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		_, err = s.invite(ctx, in.TenantID, tenant.Name, inv.Email, inv.RoleID, by.ID)
		return nil, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-invitation", Metadata: needs("", permissions.ClubManage), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/invitations/{invitationId}", Tags: []string{"users"},
		Summary: "Withdraw an invitation; its link stops working", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *InvitationPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		inv, err := s.Store.GetInvitation(ctx, in.TenantID, in.InvitationID)
		if err != nil {
			return nil, notFound(err, "Invitation not found")
		}
		if err := s.roleWithin(ctx, in.TenantID, inv.RoleID); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteInvitation(ctx, in.TenantID, in.InvitationID), "Invitation not found")
	})

	type viewOut struct{ Body InvitationView }
	huma.Register(api, huma.Operation{
		OperationID: "get-invitation", Method: http.MethodGet, Path: "/api/invitations/{token}", Tags: []string{"users"},
		Summary: "What an invitation link offers; 410 once it is accepted or expired",
	}, func(ctx context.Context, in *InvitationToken) (*viewOut, error) {
		inv, err := s.openInvitation(ctx, in.Token)
		if err != nil {
			return nil, err
		}
		_, err = s.Store.UserByEmail(ctx, inv.Email)
		if err != nil && !errors.Is(err, store.ErrNotFound) {
			return nil, err
		}
		return &viewOut{Body: InvitationView{Club: inv.TenantName, Email: inv.Email, RoleKey: inv.RoleKey, RoleName: inv.RoleName, Account: err == nil}}, nil
	})

	type acceptBody struct {
		Password string `json:"password,omitempty" required:"false" doc:"Sets the password of the new account when nobody is logged in and no account has the address"`
	}
	type acceptIn struct {
		InvitationToken
		Body *acceptBody `required:"false"`
	}
	huma.Register(api, huma.Operation{
		OperationID: "accept-invitation", Method: http.MethodPost, Path: "/api/invitations/{token}/accept", Tags: []string{"users"},
		Summary: "Accept an invitation: as the logged-in user, or by creating the account of the invited address with a password (logs it in)",
	}, func(ctx context.Context, in *acceptIn) (*AcceptOutput, error) {
		inv, err := s.openInvitation(ctx, in.Token)
		if err != nil {
			return nil, err
		}
		if u, err := currentUser(ctx); err == nil {
			if err := s.Store.AcceptInvitation(ctx, inv.ID, u.ID); err != nil {
				return nil, gone("This invitation is no longer valid")
			}
			s.Log.Info("invitation accepted", "tenant", inv.TenantID, "user", u.Name)
			me, err := s.me(ctx, u)
			return &AcceptOutput{Body: me}, err
		}
		if _, err := s.Store.UserByEmail(ctx, inv.Email); err == nil {
			return nil, conflict("An account with this address exists: log in to accept")
		} else if !errors.Is(err, store.ErrNotFound) {
			return nil, err
		}
		password := ""
		if in.Body != nil {
			password = in.Body.Password
		}
		if errs := ValidatePassword(password); len(errs) > 0 {
			return nil, invalidList(errs)
		}
		hash, err := HashPassword(password)
		if err != nil {
			return nil, err
		}
		id, err := s.Store.CreateAccount(ctx, inv.Email, hash, reqInfo(ctx).lang)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("The address is already somebody's user name")
		}
		if err != nil {
			return nil, err
		}
		if err := s.Store.VerifyEmail(ctx, id); err != nil {
			return nil, err
		}
		if err := s.Store.AcceptInvitation(ctx, inv.ID, id); err != nil {
			return nil, gone("This invitation is no longer valid")
		}
		token := randomToken()
		if err := s.Store.CreateSession(ctx, sessionID(token), id, sessionIdle, sessionTTL); err != nil {
			return nil, err
		}
		s.Log.Info("invitation accepted with a new account", "tenant", inv.TenantID, "user", inv.Email)
		u, err := s.Store.UserByName(ctx, inv.Email)
		if err != nil {
			return nil, err
		}
		me, err := s.me(ctx, u)
		cookie := sessionCookieFor(ctx, token, int(sessionTTL.Seconds()))
		return &AcceptOutput{SetCookie: cookie.String(), Body: me}, err
	})
}

// AcceptOutput is the account after accepting an invitation, with the session cookie when a new account was made.
type AcceptOutput struct {
	SetCookie string `header:"Set-Cookie"`
	Body      Me
}

// openInvitation returns the invitation of a link token; 404 when unknown, 410 once accepted or expired.
func (s *Server) openInvitation(ctx context.Context, token string) (store.InvitationInfo, error) {
	inv, err := s.Store.InvitationByToken(ctx, sessionID(token))
	switch {
	case err != nil:
		return inv, notFound(err, "Invitation not found")
	case inv.AcceptedAt != nil:
		return inv, gone("This invitation was accepted already")
	case time.Now().After(inv.ExpiresAt):
		return inv, gone("This invitation has expired; ask the club for a new one")
	}
	return inv, nil
}

// invite sends an email address the link of an invitation to a club, renewing an open invitation of the address;
// members get none.
func (s *Server) invite(ctx context.Context, tenantID int64, club, email string, roleID, by int64) (string, error) {
	if a, err := mail.ParseAddress(email); err != nil || a.Address != email {
		return "invalid", nil
	}
	if u, err := s.Store.UserByEmail(ctx, email); err == nil {
		if _, err := s.Store.MemberRole(ctx, tenantID, u.ID); err == nil {
			return "member", nil
		} else if !errors.Is(err, store.ErrNotFound) {
			return "", err
		}
	} else if !errors.Is(err, store.ErrNotFound) {
		return "", err
	}
	token := invitePrefix + randomToken()
	if err := s.Store.SaveInvitation(ctx, tenantID, email, roleID, by, sessionID(token), time.Now().Add(inviteTTL)); err != nil {
		return "", notFound(err, "Role not found")
	}
	l := mailLang(ctx, "")
	s.mail(tlmail.Message{To: email, Subject: i18n.Tf(l, "%s invites you to Standblatt", club), Body: strings.Join([]string{
		i18n.T(l, "Hello,"), "",
		i18n.Tf(l, "%s uses Standblatt for its shooting events and invites you to join. Accept the invitation with this link:", club), "",
		s.link(ctx, "/invite/"+token), "",
		i18n.T(l, "The link is valid for seven days."),
	}, "\n")})
	return "invited", nil
}
