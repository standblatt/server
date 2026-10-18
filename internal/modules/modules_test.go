package modules

import (
	"fmt"
	"strings"
	"testing"
)

// TestEnum keeps the enum tag in step with All.
func TestEnum(t *testing.T) {
	if got := strings.Join(Keys(Event), ","); got != EventEnum {
		t.Fatalf("EventEnum = %q, want %q", EventEnum, got)
	}
	if got := strings.Join(Keys(Club), ","); got != ClubEnum {
		t.Fatalf("ClubEnum = %q, want %q", ClubEnum, got)
	}
	for _, m := range All {
		for _, r := range m.Requires {
			if _, ok := Get(r); !ok {
				t.Errorf("%s requires unknown %s", m.Key, r)
			}
		}
	}
}

// TestEffective checks requirements, used modules, the range rule and the order.
func TestEffective(t *testing.T) {
	cases := []struct {
		chosen, used []string
		hasRange     bool
		want         string
	}{
		{[]string{Booking}, nil, true, "[reservations booking]"},
		{[]string{Booking, Kiosk}, nil, false, "[kiosk]"},
		{[]string{SpecialPrizes}, []string{Helpers}, false, "[awards special_prizes helpers]"},
		{[]string{"nope", Kiosk, Kiosk}, nil, true, "[kiosk]"},
		{nil, nil, true, "[]"},
	}
	for _, c := range cases {
		if got := fmt.Sprint(Effective(Event, c.chosen, c.used, c.hasRange)); got != c.want {
			t.Errorf("Effective(%v, %v, %v) = %s, want %s", c.chosen, c.used, c.hasRange, got, c.want)
		}
	}
	if got := fmt.Sprint(Missing([]string{Kiosk, Awards}, []string{Awards})); got != "[kiosk]" {
		t.Errorf("Missing = %s", got)
	}
}
