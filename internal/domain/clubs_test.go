package domain

import (
	"slices"
	"testing"
)

func TestClubKey(t *testing.T) {
	if got := ClubKey("  SG   Hilltown "); got != "sg hilltown" {
		t.Fatalf("got %q", got)
	}
}

func TestSimilarClubs(t *testing.T) {
	names := []string{"SG Hilltown", "S.G. Hilltown", "Pistolenclub Lakeside", "FSG Riverbank", "PC Lakeside", "Valley"}
	for _, c := range []struct {
		name string
		want []string
	}{
		{"sg hilltown", []string{"S.G. Hilltown"}},
		{"SG Hiltown", []string{"SG Hilltown", "S.G. Hilltown"}},
		{"Lakeside", []string{"PC Lakeside", "Pistolenclub Lakeside"}},
		{"FSG Riverbnak", []string{"FSG Riverbank"}},
		{"Mountain", nil},
		{"", nil},
	} {
		var got []string
		for _, i := range SimilarClubs(c.name, names) {
			got = append(got, names[i])
		}
		if !slices.Equal(got, c.want) {
			t.Errorf("%q: got %v, want %v", c.name, got, c.want)
		}
	}
}
