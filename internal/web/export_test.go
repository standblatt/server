package web

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/standblatt/server/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// init makes password hashes cheap, logins in tests would otherwise take a quarter second each.
func init() { bcryptCost = bcrypt.MinCost }

// TOTPCodeForTest exposes the current code of a secret to the web_test package.
func TOTPCodeForTest(secret string, now time.Time) (string, error) {
	return totpCode(secret, now.Unix()/totpPeriod)
}

// FormRangeForTest returns the range the print form of an entry is printed on.
func (s *Server) FormRangeForTest(ctx context.Context, tenantID, entryID int64) (int64, error) {
	r, err := s.renderForm(ctx, tenantID, entryID)
	return r.rangeID, err
}

// SessionIDForTest exposes the hash under which tokens are stored.
func SessionIDForTest(token string) string { return sessionID(token) }

// OperationForTest is an operation of a club with the permissions it needs.
type OperationForTest struct {
	ID, Method, Path string
	Permissions      []string
}

// ClubOperationsForTest lists every operation under /api/tenants/{tenantId}, sorted by path and method.
func ClubOperationsForTest() []OperationForTest {
	s := &Server{}
	api := humago.New(http.NewServeMux(), s.apiConfig())
	s.registerAPI(api)
	var out []OperationForTest
	for path, item := range api.OpenAPI().Paths {
		if !strings.HasPrefix(path, "/api/tenants/{tenantId}") {
			continue
		}
		for method, op := range map[string]*huma.Operation{"GET": item.Get, "POST": item.Post, "PUT": item.Put, "PATCH": item.Patch, "DELETE": item.Delete} {
			if op != nil {
				out = append(out, OperationForTest{ID: op.OperationID, Method: method, Path: path, Permissions: permissionsOf(op.Metadata)})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path || out[i].Path == out[j].Path && out[i].Method < out[j].Method
	})
	return out
}

// HelperMayForTest tells whether a helper role may call an operation.
func HelperMayForTest(role store.HelperRole, op string) bool { return helperOps[role][op] }
