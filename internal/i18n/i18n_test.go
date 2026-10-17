package i18n

import (
	"os"
	"regexp"
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	got := Parse(`msgid ""
msgstr "Language: de\n"

#: a.go
msgid "Lane"
msgstr "Scheibe"

msgid "Long"
msgstr ""
"zwei "
"Zeilen"

#, fuzzy
msgid "Old"
msgstr "Alt"

msgctxt "menu"
msgid "Open"
msgstr "Öffnen"

msgid "Empty"
msgstr ""

#~ msgid "Gone"
#~ msgstr "Weg"
`)
	want := map[string]string{"Lane": "Scheibe", "Long": "zwei Zeilen"}
	if len(got) != len(want) || got["Lane"] != want["Lane"] || got["Long"] != want["Long"] {
		t.Fatalf("got %v", got)
	}
}

func TestCatalogs(t *testing.T) {
	pot, err := os.ReadFile("locales/messages.pot")
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, m := range regexp.MustCompile(`(?m)^msgid "(.+)"$`).FindAllStringSubmatch(string(pot), -1) {
		ids = append(ids, unquote(`"`+m[1]+`"`))
	}
	verbs := regexp.MustCompile(`%[a-z]`)
	for _, locale := range []string{"de", "fr", "it"} {
		l := Get(locale)
		if len(l.msgs) == 0 {
			t.Fatalf("%s: no catalog", locale)
		}
		for _, id := range ids {
			tr, ok := l.msgs[id]
			if !ok {
				t.Errorf("%s: %q is not translated", locale, id)
				continue
			}
			if !slices.Equal(verbs.FindAllString(id, -1), verbs.FindAllString(tr, -1)) {
				t.Errorf("%s: %q and %q differ in format verbs", locale, id, tr)
			}
		}
	}
	if T(Get("de"), "Rank") != "Rang" || T(Get("en"), "Rank") != "Rank" || Tf(Lang{}, "S%d", 2) != "S2" {
		t.Fatal("lookup")
	}
}
