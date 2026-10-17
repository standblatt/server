package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/i18n"
)

// apiError is an error answer whose texts are translated when it is written (see translateErrors).
type apiError struct {
	status int
	parts  []domain.Problem
}

// Error returns the untranslated parts of the error joined with semicolons.
func (e *apiError) Error() string { return domain.Join(e.parts, "; ") }

// GetStatus makes apiError a huma.StatusError.
func (e *apiError) GetStatus() int { return e.status }

// detail returns the parts of the error translated into l, joined with semicolons.
func (e *apiError) detail(l i18n.Lang) string {
	s := make([]string, len(e.parts))
	for i, p := range e.parts {
		s[i] = i18n.Tf(l, p.Format, p.Args...)
	}
	return strings.Join(s, "; ")
}

// one returns an apiError with status and a single problem.
func one(status int, format string, a []any) error {
	return &apiError{status, []domain.Problem{domain.Problemf(format, a...)}}
}

// The error constructors of the handlers; `make i18n` extracts their first argument.
func unauthorized(format string, a ...any) error { return one(http.StatusUnauthorized, format, a) }

// forbidden returns a 403 error with a translatable message.
func forbidden(format string, a ...any) error { return one(http.StatusForbidden, format, a) }

// missing returns a 404 error with a translatable message.
func missing(format string, a ...any) error { return one(http.StatusNotFound, format, a) }

// conflict returns a 409 error with a translatable message.
func conflict(format string, a ...any) error { return one(http.StatusConflict, format, a) }

// gone returns a 410 error with a translatable message.
func gone(format string, a ...any) error { return one(http.StatusGone, format, a) }

// invalid returns a 422 error with a translatable message.
func invalid(format string, a ...any) error { return one(http.StatusUnprocessableEntity, format, a) }

// tooMany returns a 429 error with a translatable message.
func tooMany(format string, a ...any) error { return one(http.StatusTooManyRequests, format, a) }

// invalidList answers 422 with validation problems, each translated on its own.
func invalidList(errs []domain.Problem) error {
	return &apiError{http.StatusUnprocessableEntity, errs}
}

// acceptLang is the two-letter language of an Accept-Language header, "" when there is none.
func acceptLang(header string) string {
	tag, _, _ := strings.Cut(header, ",")
	tag, _, _ = strings.Cut(strings.TrimSpace(tag), ";")
	if len(tag) >= 2 {
		return strings.ToLower(tag[:2])
	}
	return ""
}

// requestLang is the catalog for the Accept-Language header; the SPA sends its UI locale, scripts usually nothing.
func requestLang(header string) i18n.Lang {
	return i18n.Get(acceptLang(header))
}

// mailLang is the catalog for a mail: the user's chosen locale, else the language of the request that triggers the mail.
func mailLang(ctx context.Context, userLocale string) i18n.Lang {
	if userLocale != "" {
		return i18n.Get(userLocale)
	}
	return i18n.Get(reqInfo(ctx).lang)
}

// translateErrors is the huma transformer that turns error answers into the language of the request.
func translateErrors(ctx huma.Context, _ string, v any) (any, error) {
	l := requestLang(ctx.Header("Accept-Language"))
	switch e := v.(type) {
	case *apiError:
		return &huma.ErrorModel{Status: e.status, Title: http.StatusText(e.status), Detail: e.detail(l)}, nil
	case *huma.ErrorModel:
		e.Detail = i18n.T(l, e.Detail)
		for _, d := range e.Errors {
			d.Message = i18n.T(l, d.Message)
		}
	}
	return v, nil
}
