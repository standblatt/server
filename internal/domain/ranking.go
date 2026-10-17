package domain

import (
	"math"
	"sort"
)

// RankingEntry is one programme result to rank. Series holds the points per series in series order.
type RankingEntry struct {
	Total    float64
	Mouches  int
	Series   []float64
	Birthday string
	Category AgeCategory
}

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

// compareRanking returns a negative number when a ranks before b. Order: higher total; then the series
// compared from the last one backwards; then more mouches; then age: among juniors (U10–U21) the younger,
// among the others the older shooter ranks first, and non-juniors rank before juniors. An unknown birthday
// does not decide.
func compareRanking(a, b RankingEntry) int {
	if d := pointsDiff(a.Total, b.Total); d != 0 {
		return -d
	}
	for i, j := len(a.Series)-1, len(b.Series)-1; i >= 0 || j >= 0; i, j = i-1, j-1 {
		var x, y float64
		if i >= 0 {
			x = a.Series[i]
		}
		if j >= 0 {
			y = b.Series[j]
		}
		if d := pointsDiff(x, y); d != 0 {
			return -d
		}
	}
	if a.Mouches != b.Mouches {
		return b.Mouches - a.Mouches
	}
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
	order := make([]int, len(entries))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return compareRanking(entries[order[i]], entries[order[j]]) < 0
	})
	return order
}
