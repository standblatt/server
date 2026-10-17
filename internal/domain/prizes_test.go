package domain

import (
	"fmt"
	"testing"
)

func TestSpecialPrizes(t *testing.T) {
	// Elite of the 2026 Lupi championship (sport): the first three come from clubs with fewer than 4 participants.
	var c []PrizeCandidate
	add := func(cat AgeCategory, club string, n int) {
		for i := 0; i < n; i++ {
			c = append(c, PrizeCandidate{Category: cat, Club: club})
		}
	}
	add(Elite, "St. Gallen", 1)
	add(Elite, "Kloten", 1)
	add(Elite, "Igis", 1)
	add(Elite, "Rüthi", 1)
	add(Elite, "Diepoldsau", 1)
	add(Elite, "Neckertal", 1)
	add(Elite, "", 1)
	add(Elite, "Diessenhofen", 1)
	add(Elite, "Filler", 8)
	add(Senior, "Rüthi", 3)
	add(Senior, "Neckertal", 3)
	add(Senior, "Diessenhofen", 3)
	add(Senior, "Kloten", 2)
	add(Senior, "Filler", 1)
	add(U17, "Rüthi", 3)
	add(Veteran, "Igis", 2)
	add(Veteran, "Neckertal", 8)

	got := SpecialPrizes(c, SpecialPrizeRules{Ranks: 3, MinCategory: 10, MinClub: 4})
	want := []string{
		"elite 16 winners [3 5 7] skipped [0 1 2 4 6]",
		"u17 3 winners [] skipped []",
		"senior 12 winners [16 17 18] skipped []",
		"veteran 10 winners [33 34 35] skipped [31 32]",
	}
	order := map[AgeCategory]int{}
	for i, cat := range AgeCategories {
		order[cat] = i
	}
	if len(got) != len(want) {
		t.Fatalf("got %d categories: %+v", len(got), got)
	}
	for i := 1; i < len(got); i++ {
		if order[got[i-1].Category] > order[got[i].Category] {
			t.Fatalf("categories out of order: %+v", got)
		}
	}
	for _, g := range got {
		line := fmt.Sprintf("%s %d winners %v skipped %v", g.Category, g.Participants, nonNil(g.Winners), nonNil(g.Skipped))
		found := false
		for _, w := range want {
			found = found || w == line
		}
		if !found {
			t.Errorf("unexpected %s", line)
		}
	}

	if got := SpecialPrizes(c, SpecialPrizeRules{}); len(got[0].Winners) != 0 {
		t.Fatalf("ranks 0 turns prizes off: %+v", got[0])
	}
	if got := SpecialPrizes(c[:8], SpecialPrizeRules{Ranks: 1}); len(got[0].Winners) != 1 || got[0].Winners[0] != 0 {
		t.Fatalf("without minimums the first wins: %+v", got[0])
	}
}

func nonNil(s []int) []int {
	if s == nil {
		return []int{}
	}
	return s
}
