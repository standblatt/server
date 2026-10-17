// Package i18n translates the texts of printed documents and exports with the gettext catalogs in locales/.
// Source strings are English; `make i18n` extracts calls of T and Tf into locales/messages.pot and merges it
// into the catalogs.
package i18n

import (
	"bufio"
	"embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed locales/*.po
var files embed.FS

// Lang is the catalog of one language. The zero value returns the English source strings.
type Lang struct {
	msgs map[string]string
}

var langs = map[string]Lang{}

// init loads the embedded PO catalogs, keyed by file name without .po.
func init() {
	entries, err := files.ReadDir("locales")
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		b, err := files.ReadFile("locales/" + e.Name())
		if err != nil {
			panic(err)
		}
		langs[strings.TrimSuffix(e.Name(), ".po")] = Lang{msgs: Parse(string(b))}
	}
}

// Get returns the catalog of a locale such as "de"; unknown locales get English.
func Get(locale string) Lang {
	return langs[locale]
}

// T translates s.
func T(l Lang, s string) string {
	if v := l.msgs[s]; v != "" {
		return v
	}
	return s
}

// Tf translates format and formats it with a like fmt.Sprintf.
func Tf(l Lang, format string, a ...any) string {
	return fmt.Sprintf(T(l, format), a...)
}

// Parse reads the translated singular messages of a PO file, skipping fuzzy entries and entries with a context.
func Parse(po string) map[string]string {
	out := map[string]string{}
	var msgid, msgstr string
	var cur *string
	skip := false
	flush := func() {
		if !skip && msgid != "" && msgstr != "" {
			out[msgid] = msgstr
		}
		msgid, msgstr, cur, skip = "", "", nil, false
	}
	sc := bufio.NewScanner(strings.NewReader(po))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		key, value, _ := strings.Cut(line, " ")
		switch {
		case line == "":
			flush()
		case key == "#,":
			skip = skip || strings.Contains(line, "fuzzy")
		case strings.HasPrefix(line, "#"):
		case key == "msgctxt":
			skip, cur = true, nil
		case key == "msgid":
			msgid, cur = unquote(value), &msgid
		case key == "msgstr":
			msgstr, cur = unquote(value), &msgstr
		case strings.HasPrefix(line, `"`) && cur != nil:
			*cur += unquote(line)
		default:
			cur = nil
		}
	}
	flush()
	return out
}

// unquote decodes a quoted PO string, or returns an empty string if it is malformed.
func unquote(s string) string {
	v, err := strconv.Unquote(s)
	if err != nil {
		return ""
	}
	return v
}
