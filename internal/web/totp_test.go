package web

import (
	"testing"
	"time"
)

func TestTOTPVectors(t *testing.T) {
	// RFC 6238 appendix B, SHA-1, secret "12345678901234567890", last six of the eight digits.
	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	for _, c := range []struct {
		at   int64
		want string
	}{{59, "287082"}, {1111111109, "081804"}, {1234567890, "005924"}, {20000000000, "353130"}} {
		got, err := totpCode(secret, c.at/totpPeriod)
		if err != nil || got != c.want {
			t.Errorf("at %d: %s %v, want %s", c.at, got, err, c.want)
		}
		now := time.Unix(c.at, 0)
		if totpStep(secret, c.want, now) == 0 || totpStep(secret, c.want, now.Add(totpPeriod*time.Second)) == 0 || totpStep(secret, c.want, now.Add(-totpPeriod*time.Second)) == 0 {
			t.Errorf("at %d: code not accepted within one step", c.at)
		}
		if totpStep(secret, c.want, now.Add(2*totpPeriod*time.Second)) != 0 || totpStep(secret, "000000", now) != 0 {
			t.Errorf("at %d: stale or wrong code accepted", c.at)
		}
	}
	if _, err := totpCode("not base32!", 0); err == nil {
		t.Error("bad secret accepted")
	}
}
