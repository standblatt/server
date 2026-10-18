package web

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
)

// pairingAlphabet has no characters that look alike (no I, L, O, 0, 1).
const pairingAlphabet = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"

// Rate limits of pairing: new codes per IP, and code lookups and confirmations per IP, so codes cannot be guessed
// within their ten minutes.
const (
	pairWindow   = 15 * time.Minute
	pairMax      = 20
	pairCodesMax = 30
)

// pairingCode returns eight random characters of pairingAlphabet.
func pairingCode() string {
	var b strings.Builder
	buf := make([]byte, 1)
	for b.Len() < 8 {
		rand.Read(buf)
		// 248 = 8 × 31: rejecting the rest keeps every character equally likely.
		if buf[0] < 248 {
			b.WriteByte(pairingAlphabet[int(buf[0])%len(pairingAlphabet)])
		}
	}
	return b.String()
}

// normalizeCode upper-cases a typed code and keeps only its letters and digits.
func normalizeCode(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(s) {
		if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// formatCode shows a code as two groups of four, e.g. K7QD-4MZP.
func formatCode(code string) string {
	return code[:4] + "-" + code[4:]
}

// PairingCodePath identifies a pairing by the code the user typed.
type PairingCodePath struct {
	Code string `path:"code" maxLength:"20"`
}

// registerPairing registers pairing by code: the agent asks for a code and polls with its secret, a member looks
// the code up and confirms it for a range.
func (s *Server) registerPairing(api huma.API) {
	type pairIn struct {
		Body struct {
			Name    string `json:"name" minLength:"1" maxLength:"100" doc:"Name of the computer, proposed as the agent's name"`
			OS      string `json:"os" enum:"linux,windows,darwin"`
			Version string `json:"version" minLength:"1" maxLength:"50"`
		}
	}
	type pairOut struct {
		Body struct {
			Code      string    `json:"code" doc:"Shown to the user, e.g. K7QD-4MZP; confirmed on /connect?code=…"`
			Secret    string    `json:"secret" doc:"Poll secret for POST /api/agent/pair/status"`
			ExpiresAt time.Time `json:"expiresAt"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "agent-pair", Method: http.MethodPost, Path: "/api/agent/pair", Tags: []string{"agent"},
		Summary: "Start pairing: a code to confirm on /connect within ten minutes and a secret to poll with", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *pairIn) (*pairOut, error) {
		if err := s.rateLimited(ctx, "pair:"+reqInfo(ctx).ip, pairWindow, pairMax); err != nil {
			return nil, err
		}
		out := &pairOut{}
		out.Body.Secret = randomToken()
		out.Body.ExpiresAt = time.Now().Add(store.PairingTTL)
		p := store.Pairing{Name: strings.TrimSpace(in.Body.Name), OS: in.Body.OS, Version: in.Body.Version, ExpiresAt: out.Body.ExpiresAt}
		var err error
		for range 5 {
			code := pairingCode()
			if err = s.Store.CreatePairing(ctx, code, sessionID(out.Body.Secret), p); !errors.Is(err, store.ErrExists) {
				out.Body.Code = formatCode(code)
				break
			}
		}
		if err != nil {
			return nil, err
		}
		return out, nil
	})

	type statusIn struct {
		Body struct {
			Secret string `json:"secret" minLength:"1" maxLength:"100"`
		}
	}
	type statusOut struct {
		Body struct {
			Status string `json:"status" enum:"pending,connected"`
			Token  string `json:"token,omitempty" doc:"The agent token, given once when connected"`
			Club   string `json:"club,omitempty"`
			Range  string `json:"range,omitempty"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "agent-pair-status", Method: http.MethodPost, Path: "/api/agent/pair/status", Tags: []string{"agent"},
		Summary: "Poll a pairing: pending until confirmed, then connected with the agent token, given once; 410 when the code expired",
	}, func(ctx context.Context, in *statusIn) (*statusOut, error) {
		out := &statusOut{}
		token := randomToken()
		club, rng, err := s.Store.FetchPairing(ctx, sessionID(in.Body.Secret), sessionID(token))
		switch {
		case errors.Is(err, store.ErrPairingPending):
			out.Body.Status = "pending"
			return out, nil
		case errors.Is(err, store.ErrPairingExpired):
			return nil, gone("The code expired. Start connecting again.")
		case err != nil:
			return nil, notFound(err, "Unknown pairing")
		}
		out.Body.Status, out.Body.Token, out.Body.Club, out.Body.Range = "connected", token, club, rng
		return out, nil
	})

	type pairingOut struct{ Body store.Pairing }
	huma.Register(api, huma.Operation{
		OperationID: "get-agent-pairing", Method: http.MethodGet, Path: "/api/agent-pairings/{code}", Tags: []string{"printing"},
		Summary: "The computer behind an open pairing code, for the logged-in user who is about to confirm it",
	}, func(ctx context.Context, in *PairingCodePath) (*pairingOut, error) {
		if _, err := currentUser(ctx); err != nil {
			return nil, err
		}
		if err := s.rateLimited(ctx, "pair-code:"+reqInfo(ctx).ip, pairWindow, pairCodesMax); err != nil {
			return nil, err
		}
		p, err := s.Store.OpenPairing(ctx, normalizeCode(in.Code))
		if err != nil {
			return nil, notFound(err, "Unknown or expired code")
		}
		return &pairingOut{Body: p}, nil
	})

	type confirmIn struct {
		RangePrintPath
		Body struct {
			Code string `json:"code" minLength:"1" maxLength:"20"`
			Name string `json:"name" minLength:"1" maxLength:"100"`
		}
	}
	type confirmOut struct{ Body store.PrintAgent }
	huma.Register(api, huma.Operation{
		OperationID: "pair-print-agent", Metadata: needs("", permissions.RangesManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-agents/pair", Tags: []string{"printing"},
		Summary: "Confirm a pairing code: the computer becomes a print agent of the range", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *confirmIn) (*confirmOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := s.rateLimited(ctx, "pair-code:"+reqInfo(ctx).ip, pairWindow, pairCodesMax); err != nil {
			return nil, err
		}
		a, err := s.Store.ConfirmPairing(ctx, in.TenantID, in.RangeID, normalizeCode(in.Body.Code), strings.TrimSpace(in.Body.Name), sessionID(randomToken()))
		if err != nil {
			return nil, notFound(err, "Unknown or expired code")
		}
		return &confirmOut{Body: a}, nil
	})
}
