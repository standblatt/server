package web

import "time"

// TOTPCodeForTest exposes the current code of a secret to the web_test package.
func TOTPCodeForTest(secret string, now time.Time) (string, error) {
	return totpCode(secret, now.Unix()/totpPeriod)
}
