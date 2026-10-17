package domain

import (
	"strings"
	"testing"
)

func ptr(n int) *int { return &n }

func TestValidateResults(t *testing.T) {
	series := []Series{{Shots: 10, MaxPoints: 100}, {Shots: 5, MaxPoints: 50, TrackMouches: true}}
	ok := []SeriesResult{{Points: 100}, {Points: 47.5, Mouches: ptr(5)}}
	if errs := ValidateResults(series, ok); len(errs) != 0 {
		t.Fatalf("valid results rejected: %v", errs)
	}
	cases := map[string][]SeriesResult{
		"Enter all 2 series":              {{Points: 10}},
		"between 0 and 100":               {{Points: 100.5}, {Points: 1}},
		"Series 1 does not track":         {{Points: 1, Mouches: ptr(1)}, {Points: 1}},
		"mouches must be between 0 and 5": {{Points: 1}, {Points: 1, Mouches: ptr(6)}},
	}
	for want, results := range cases {
		if errs := ValidateResults(series, results); !strings.Contains(Join(errs, "|"), want) {
			t.Errorf("%s: got %v", want, errs)
		}
	}
}

func TestTotalsAndCombination(t *testing.T) {
	results := []SeriesResult{{Points: 95}, {Points: 88.1, Mouches: ptr(2)}, {Points: 91.2}, {Points: 97, Mouches: ptr(1)}}
	if p, m := Totals(results); p != 371.3 || m != 3 {
		t.Fatalf("totals: %v %d", p, m)
	}
	got := CombinedResults(results, []int{4, 2})
	if len(got) != 2 || got[0].Points != 88.1 || got[1].Points != 97 {
		t.Fatalf("combined: %+v", got)
	}
}

func TestEligibleAwards(t *testing.T) {
	tiers := []Tier{
		{ID: 1, Position: 1, Name: "Badge", Kind: AwardNatural, Value: 0, Limits: []AwardLimit{{AgeCategory: Elite, Points: 80}}},
		{ID: 2, Position: 2, Name: "Wreath", Kind: AwardCard, Value: 25, Limits: []AwardLimit{{AgeCategory: Elite, Points: 90}, {AgeCategory: U21, Points: 85}}},
		{ID: 3, Position: 3, Name: "Gold", Kind: AwardCard, Value: 25, Limits: []AwardLimit{{AgeCategory: Elite, Points: 95}}},
		{ID: 4, Position: 4, Name: "Card 6", Kind: AwardCard, Value: 6, Limits: []AwardLimit{{AgeCategory: Senior, Points: 80}}},
		{ID: 5, Position: 5, Name: "Card 8", Kind: AwardCard, Value: 8, Limits: []AwardLimit{{AgeCategory: Senior, Points: 90}}},
	}
	names := func(ts []Tier) string {
		var s []string
		for _, t := range ts {
			s = append(s, t.Name)
		}
		return strings.Join(s, ",")
	}
	cases := []struct {
		cat   AgeCategory
		total float64
		want  string
	}{
		{Elite, 79.9, ""},
		{Elite, 80, "Badge"},
		{Elite, 96, "Wreath,Gold,Badge"},
		{U21, 86, "Wreath"},
		{Senior, 85, "Card 6"},
		{Senior, 100, "Card 8"},
		{Veteran, 100, ""},
	}
	for _, c := range cases {
		if got := names(EligibleAwards(tiers, c.cat, c.total)); got != c.want {
			t.Errorf("%s %.1f: got %q, want %q", c.cat, c.total, got, c.want)
		}
	}
}
