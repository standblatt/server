package modules

import (
	"fmt"
	"testing"
)

// TestRequires checks that every required module exists.
func TestRequires(t *testing.T) {
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
		{[]string{Reservations}, nil, true, "[lanes reservations]"},
		{[]string{Reservations, Kiosk}, nil, false, "[kiosk]"},
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
