package domain

import (
	"math"
	"sort"
)

// SeriesResult is the score of one series. Mouches is nil for series that do not track them.
type SeriesResult struct {
	Points  float64 `json:"points" minimum:"0"`
	Mouches *int    `json:"mouches,omitempty" required:"false" minimum:"0"`
}

// ValidateResults checks one result per series against the series definitions.
func ValidateResults(series []Series, results []SeriesResult) []Problem {
	if len(results) != len(series) {
		return []Problem{Problemf("Enter all %d series", len(series))}
	}
	var errs []Problem
	for i, r := range results {
		s := series[i]
		if r.Points < 0 || r.Points > float64(s.MaxPoints) || math.IsNaN(r.Points) {
			errs = append(errs, Problemf("Series %d: points must be between 0 and %d", i+1, s.MaxPoints))
		}
		switch {
		case r.Mouches != nil && !s.TrackMouches:
			errs = append(errs, Problemf("Series %d does not track mouches", i+1))
		case r.Mouches != nil && (*r.Mouches < 0 || *r.Mouches > s.Shots):
			errs = append(errs, Problemf("Series %d: mouches must be between 0 and %d", i+1, s.Shots))
		}
	}
	return errs
}

// Totals sums points and mouches.
func Totals(results []SeriesResult) (points float64, mouches int) {
	for _, r := range results {
		points += r.Points
		if r.Mouches != nil {
			mouches += *r.Mouches
		}
	}
	return math.Round(points*100) / 100, mouches
}

// CombinedResults picks the given series (1-based) of a target programme's results for a combined programme,
// numbered from 1 in ascending order of the picked series.
func CombinedResults(target []SeriesResult, picked []int) []SeriesResult {
	sorted := append([]int(nil), picked...)
	sort.Ints(sorted)
	out := make([]SeriesResult, 0, len(sorted))
	for _, n := range sorted {
		if n >= 1 && n <= len(target) {
			out = append(out, target[n-1])
		}
	}
	return out
}

// Tier is an award tier of a stored programme.
type Tier struct {
	ID       int64        `json:"id"`
	Position int          `json:"position"`
	Name     string       `json:"name"`
	Kind     AwardKind    `json:"kind" enum:"natural,card,cash"`
	Value    float64      `json:"value"`
	Limits   []AwardLimit `json:"limits"`
}

// EligibleAwards returns the tiers a total reaches in an age category, most valuable first. A tier without
// a limit for the category is not available to it, and a tier of a kind whose more valuable tier is also
// reached is left out: nobody picks the CHF 6 wreath card when the CHF 8 one is reached.
func EligibleAwards(tiers []Tier, category AgeCategory, total float64) []Tier {
	var reached []Tier
	best := map[AwardKind]float64{}
	for _, t := range tiers {
		for _, l := range t.Limits {
			if l.AgeCategory == category && total+1e-9 >= l.Points {
				reached = append(reached, t)
				best[t.Kind] = max(best[t.Kind], t.Value)
				break
			}
		}
	}
	var out []Tier
	for _, t := range reached {
		if t.Value >= best[t.Kind] {
			out = append(out, t)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Value != out[j].Value {
			return out[i].Value > out[j].Value
		}
		return out[i].Position < out[j].Position
	})
	return out
}
