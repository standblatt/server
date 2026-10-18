package domain

import (
	"cmp"
	"slices"
	"strings"
	"unicode"
)

// ClubKey is what club names are matched by: lower case with single spaces, so "SG  hilltown " is "sg hilltown".
func ClubKey(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), " "))
}

// looseKey is the club key with only letters and digits, so "S.G. Hilltown" and "SG Hilltown" meet.
func looseKey(key string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, key)
}

// SimilarClubs returns the indexes of names that look like name without matching its key: the same letters and
// digits, a few typos apart, or one containing the other. Best first, at most five.
func SimilarClubs(name string, names []string) []int {
	key := ClubKey(name)
	loose := looseKey(key)
	limit := 2
	if len([]rune(loose)) >= 12 {
		limit = 3
	}
	type hit struct{ i, score int }
	var hits []hit
	for i, n := range names {
		k := ClubKey(n)
		if k == key || loose == "" {
			continue
		}
		l := looseKey(k)
		short, long := l, loose
		if len(short) > len(long) {
			short, long = long, short
		}
		switch d := levenshtein(loose, l); {
		case d <= limit:
			hits = append(hits, hit{i, d})
		case len([]rune(short)) >= 4 && strings.Contains(long, short):
			hits = append(hits, hit{i, limit + 1})
		}
	}
	slices.SortStableFunc(hits, func(a, b hit) int { return cmp.Compare(a.score, b.score) })
	out := make([]int, 0, 5)
	for _, h := range hits[:min(len(hits), 5)] {
		out = append(out, h.i)
	}
	return out
}

// levenshtein counts the single-rune edits between a and b.
func levenshtein(a, b string) int {
	x, y := []rune(a), []rune(b)
	prev := make([]int, len(y)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(x); i++ {
		cur := make([]int, len(y)+1)
		cur[0] = i
		for j := 1; j <= len(y); j++ {
			cost := 1
			if x[i-1] == y[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(y)]
}

// NormalizeLicense trims a licence number and pads a purely numeric one to six digits, as the SSV numbers are stored.
func NormalizeLicense(l string) string {
	l = strings.TrimSpace(l)
	if l == "" || len(l) >= 6 || strings.Trim(l, "0123456789") != "" {
		return l
	}
	return strings.Repeat("0", 6-len(l)) + l
}
