package web

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
)

const totpPeriod = 30

var totpEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// newTOTPSecret returns a new random 160-bit TOTP secret in base32.
func newTOTPSecret() string {
	b := make([]byte, 20)
	rand.Read(b)
	return totpEncoding.EncodeToString(b)
}

// totpCode is RFC 6238 with SHA-1, six digits and 30 s steps, what every authenticator app does by default.
func totpCode(secret string, step int64) (string, error) {
	key, err := totpEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(step))
	mac := hmac.New(sha1.New, key)
	mac.Write(msg[:])
	sum := mac.Sum(nil)
	off := sum[len(sum)-1] & 0x0f
	n := binary.BigEndian.Uint32(sum[off:off+4]) & 0x7fffffff
	return fmt.Sprintf("%06d", n%1000000), nil
}

// totpStep returns the time step a code is valid for (now, or one step either side for clock drift), or 0.
func totpStep(secret, code string, now time.Time) int64 {
	code = strings.ReplaceAll(strings.TrimSpace(code), " ", "")
	step := now.Unix() / totpPeriod
	for _, d := range []int64{0, -1, 1} {
		if want, err := totpCode(secret, step+d); err == nil && hmac.Equal([]byte(want), []byte(code)) {
			return step + d
		}
	}
	return 0
}

// qrDataURL renders a URL as a QR code PNG data URL (second factor setup, booking link).
func qrDataURL(otpauth string) (string, error) {
	code, err := qr.Encode(otpauth, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	code, err = barcode.Scale(code, 240, 240)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, code); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// checkTOTP verifies a login code against the user's secret and burns its time step; "" when 2FA is off.
func (s *Server) checkTOTP(ctx context.Context, userID int64, code string) (required bool, err error) {
	if !s.TOTP {
		return false, nil
	}
	secret, _, err := s.Store.TOTP(ctx, userID)
	if err != nil || secret == "" {
		return false, err
	}
	if code == "" {
		return true, &apiError{http.StatusPreconditionRequired, []domain.Problem{domain.Problemf("Enter the code from your authenticator app.")}}
	}
	step := totpStep(secret, code, time.Now())
	if step == 0 {
		return true, unauthorized("Wrong code.")
	}
	if err := s.Store.UseTOTP(ctx, userID, step); errors.Is(err, store.ErrNotFound) {
		return true, unauthorized("This code was used already; wait for the next one.")
	} else if err != nil {
		return true, err
	}
	return true, nil
}

// registerTOTP registers the endpoints that set up, enable and disable the second factor.
func (s *Server) registerTOTP(api huma.API) {
	type setupOut struct {
		Body struct {
			Secret string `json:"secret" doc:"Base32 secret to type in when the QR code cannot be scanned"`
			URL    string `json:"url" doc:"otpauth:// URL encoded in the QR code"`
			QR     string `json:"qr" doc:"PNG data URL of the QR code"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "totp-setup", Method: http.MethodPost, Path: "/api/me/totp/setup", Tags: []string{"auth"},
		Summary: "Propose a new authenticator secret; nothing changes until it is confirmed with a code",
	}, func(ctx context.Context, _ *struct{}) (*setupOut, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if !s.TOTP {
			return nil, conflict("The second factor is switched off on this server")
		}
		out := &setupOut{}
		out.Body.Secret = newTOTPSecret()
		out.Body.URL = "otpauth://totp/" + url.PathEscape("Standblatt:"+u.Name) + "?" + url.Values{
			"secret": {out.Body.Secret}, "issuer": {"Standblatt"}, "algorithm": {"SHA1"}, "digits": {"6"}, "period": {"30"},
		}.Encode()
		out.Body.QR, err = qrDataURL(out.Body.URL)
		return out, err
	})

	type enableIn struct {
		Body struct {
			Secret string `json:"secret" minLength:"16" maxLength:"64"`
			Code   string `json:"code" minLength:"6" maxLength:"8" doc:"Current code from the app, proves the secret was saved"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "totp-enable", Method: http.MethodPut, Path: "/api/me/totp", Tags: []string{"auth"},
		Summary: "Switch the second factor on", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *enableIn) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if !s.TOTP {
			return nil, conflict("The second factor is switched off on this server")
		}
		if _, err := totpCode(in.Body.Secret, 0); err != nil {
			return nil, invalid("Secret: not a valid base32 string")
		}
		step := totpStep(in.Body.Secret, in.Body.Code, time.Now())
		if step == 0 {
			return nil, unauthorized("Wrong code.")
		}
		if err := s.Store.SetTOTP(ctx, u.ID, in.Body.Secret); err != nil {
			return nil, err
		}
		s.Log.Info("totp enabled", "user", u.Name)
		return nil, s.Store.UseTOTP(ctx, u.ID, step)
	})

	type disableIn struct {
		Body struct {
			Password string `json:"password" minLength:"1"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "totp-disable", Method: http.MethodPost, Path: "/api/me/totp/disable", Tags: []string{"auth"},
		Summary: "Switch the second factor off; needs the password", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *disableIn) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		_, hash, err := s.Store.Credentials(ctx, u.Name)
		if err != nil {
			return nil, err
		}
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Body.Password)) != nil {
			return nil, unauthorized("Wrong password.")
		}
		s.Log.Info("totp disabled", "user", u.Name)
		return nil, s.Store.SetTOTP(ctx, u.ID, "")
	})
}
