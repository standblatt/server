package web

import (
	"context"
	"time"

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
