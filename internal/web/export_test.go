package web

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// init makes password hashes cheap, logins in tests would otherwise take a quarter second each.
func init() { bcryptCost = bcrypt.MinCost }

// TOTPCodeForTest exposes the current code of a secret to the web_test package.
func TOTPCodeForTest(secret string, now time.Time) (string, error) {
	return totpCode(secret, now.Unix()/totpPeriod)
}
