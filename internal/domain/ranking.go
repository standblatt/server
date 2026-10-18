package domain

import (
	"math"
	"slices"
	"sort"
)

// RankingEntry is one programme result to rank. Series holds the points per series in series order, Shots the
// single shots where they were recorded.
type RankingEntry struct {
	Total    float64
	Mouches  int
	Series   []float64
	Shots    []float64
	Birthday string
	Category AgeCategory
}

// TieBreak is a criterion that decides between equal totals.
type TieBreak string

// Tie-breaks: the series compared from the last one backwards, more of the highest single shots, more mouches, age.
const (
	TieSeries  TieBreak = "series"
	TieShots   TieBreak = "shots"
	TieMouches TieBreak = "mouches"
	TieAge     TieBreak = "age"
)

// SSVTieBreak is the order of the SSV ranking rules.
var SSVTieBreak = []TieBreak{TieSeries, TieMouches, TieAge}

const pointsTolerance = 0.005

// pointsDiff compares two point values within pointsTolerance and returns -1, 0 or 1.
func pointsDiff(a, b float64) int {
	switch d := a - b; {
	case math.Abs(d) < pointsTolerance:
		return 0
	case d > 0:
		return 1
	default:
		return -1
	}
}

// compareRanking returns a negative number when a ranks before b: the higher total, then the tie-breaks in order.
// Series compare from the last one backwards; shots compare the single shots from the highest down, so more tens
// win, then more nines; age ranks the younger first among juniors (U10–U21), the older among the others, and
// non-juniors before juniors, while an unknown birthday does not decide.
func compareRanking(a, b RankingEntry, criteria []TieBreak) int {
	if d := pointsDiff(a.Total, b.Total); d != 0 {
		return -d
	}
	for _, c := range criteria {
		if d := tieBreak(c, a, b); d != 0 {
			return d
		}
	}
	return 0
}

// tieBreak compares a and b by one criterion; negative when a ranks before b.
func tieBreak(c TieBreak, a, b RankingEntry) int {
	switch c {
	case TieSeries:
		return compareBackwards(a.Series, b.Series)
	case TieShots:
		return compareHighest(a.Shots, b.Shots)
	case TieMouches:
		return b.Mouches - a.Mouches
	case TieAge:
		if a.Birthday == "" || b.Birthday == "" {
			return 0
		}
		aj, bj := a.Category.Junior(), b.Category.Junior()
		switch {
		case aj && !bj:
			return 1
		case !aj && bj:
			return -1
		case aj:
			return compareStrings(b.Birthday, a.Birthday)
		default:
			return compareStrings(a.Birthday, b.Birthday)
		}
	}
	return 0
}

// compareBackwards compares two point lists from the last value backwards, a missing value counting 0.
func compareBackwards(a, b []float64) int {
	for i, j := len(a)-1, len(b)-1; i >= 0 || j >= 0; i, j = i-1, j-1 {
		var x, y float64
		if i >= 0 {
			x = a[i]
		}
		if j >= 0 {
			y = b[j]
		}
		if d := pointsDiff(x, y); d != 0 {
			return -d
		}
	}
	return 0
}

// compareHighest compares two lists of single shots sorted from the highest down, a missing shot counting 0.
func compareHighest(a, b []float64) int {
	x := slices.Sorted(slices.Values(a))
	y := slices.Sorted(slices.Values(b))
	slices.Reverse(x)
	slices.Reverse(y)
	for i := 0; i < len(x) || i < len(y); i++ {
		var p, q float64
		if i < len(x) {
			p = x[i]
		}
		if i < len(y) {
			q = y[i]
		}
		if d := pointsDiff(p, q); d != 0 {
			return -d
		}
	}
	return 0
}

// compareStrings returns -1, 0 or 1 as a sorts before, equal to or after b.
func compareStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// Rank orders the entries by the SSV ranking rules and returns their indexes in rank order;
// the rank of the entry at position i of the result is i+1.
func Rank(entries []RankingEntry) []int {
	return RankBy(entries, SSVTieBreak)
}

// RankBy orders the entries by total and then the given tie-breaks and returns their indexes in rank order.
func RankBy(entries []RankingEntry, criteria []TieBreak) []int {
	order := make([]int, len(entries))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return compareRanking(entries[order[i]], entries[order[j]], criteria) < 0
	})
	return order
}

// Counting returns the indexes of the n attempts that count under the rule, the best or earliest counting one
// first; attempts must be given in attempt order. Fewer are returned when there are fewer attempts.
func Counting(rule AttemptRule, attempts []RankingEntry, n int) []int {
	n = min(n, len(attempts))
	switch rule {
	case AttemptFirst:
		out := make([]int, n)
		for i := range out {
			out[i] = i
		}
		return out
	case AttemptLast:
		out := make([]int, n)
		for i := range out {
			out[i] = len(attempts) - 1 - i
		}
		return out
	default:
		return Rank(attempts)[:n]
	}
}
