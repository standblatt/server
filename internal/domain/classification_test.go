package domain

import (
	"fmt"
	"strings"
	"testing"
)

// res is an attempt of a shooter at a programme with the given series points.
func res(shooter, programme int64, cat AgeCategory, series ...float64) AttemptResult {
	total := 0.0
	for _, v := range series {
		total += v
	}
	return AttemptResult{ShooterID: shooter, ProgrammeID: programme, Rule: AttemptBest,
		Entry: RankingEntry{Total: total, Series: series, Category: cat, Birthday: "1980-01-01"}}
}

// order returns the shooter IDs of the standings.
func order(list []Standing) []int64 {
	out := []int64{}
	for _, s := range list {
		out = append(out, s.ShooterID)
	}
	return out
}

func TestClassifySectionChampionship(t *testing.T) {
	c := ClassificationDefinition{Name: "Section", Complete: true, TieBreak: SSVTieBreak,
		Parts: []ClassificationPart{{ProgrammeID: 1, Count: 1}, {ProgrammeID: 2, Count: 1}, {ProgrammeID: 3, Count: 1}}}
	got := Classify(c, []AttemptResult{
		res(1, 1, Elite, 90), res(1, 2, Elite, 80), res(1, 3, Elite, 70),
		res(2, 1, Elite, 95), res(2, 2, Elite, 85), res(2, 3, Elite, 60),
		res(3, 1, Elite, 99), res(3, 2, Elite, 99),
	})
	if len(got) != 1 || got[0].Name != "" {
		t.Fatalf("one class without classes: %+v", got)
	}
	if o := order(got[0].Ranked); len(o) != 2 || o[0] != 1 || o[1] != 2 {
		t.Fatalf("ranked: %v", o)
	}
	if got[0].Ranked[0].Entry.Total != 240 || got[0].Ranked[0].Rank != 1 || got[0].Ranked[1].Rank != 2 {
		t.Fatalf("totals and ranks: %+v", got[0].Ranked)
	}
	if o := order(got[0].Unranked); len(o) != 1 || o[0] != 3 || got[0].Unranked[0].Parts[2] != nil {
		t.Fatalf("incomplete shooter: %+v", got[0].Unranked)
	}
}

func TestClassifyBestAttemptsAndParts(t *testing.T) {
	season := ClassificationDefinition{Name: "Club championship", TieBreak: SSVTieBreak, Parts: []ClassificationPart{{ProgrammeID: 1, Count: 3}}}
	got := Classify(season, []AttemptResult{
		res(1, 1, Elite, 80), res(1, 1, Elite, 95), res(1, 1, Elite, 90), res(1, 1, Elite, 70),
		res(2, 1, Elite, 99), res(2, 1, Elite, 99),
	})
	if r := got[0].Ranked; len(r) != 2 || r[0].ShooterID != 1 || r[0].Entry.Total != 265 || *r[0].Parts[0] != 265 || r[1].Complete {
		t.Fatalf("best three attempts: %+v", r)
	}

	best := ClassificationDefinition{Name: "Best two", Best: 2, Complete: true, TieBreak: SSVTieBreak,
		Parts: []ClassificationPart{{ProgrammeID: 1, Count: 1}, {ProgrammeID: 2, Count: 1}, {ProgrammeID: 3, Count: 1}}}
	got = Classify(best, []AttemptResult{res(1, 1, Elite, 50), res(1, 2, Elite, 90), res(1, 3, Elite, 80), res(2, 1, Elite, 100)})
	r := got[0].Ranked
	if len(r) != 1 || r[0].Entry.Total != 170 || r[0].Counted[0] || !r[0].Counted[1] || !r[0].Counted[2] {
		t.Fatalf("best two of three parts: %+v", r)
	}
	if len(got[0].Unranked) != 1 {
		t.Fatalf("one part of two is incomplete: %+v", got[0].Unranked)
	}
}

func TestClassifyTieBreaksAndClasses(t *testing.T) {
	shots := func(r AttemptResult, s ...float64) AttemptResult { r.Entry.Shots = s; return r }
	results := []AttemptResult{
		shots(res(1, 1, Elite, 50, 40), 10, 10, 10, 10, 10, 8, 8, 8, 8, 8),
		shots(res(2, 1, Elite, 40, 50), 10, 10, 10, 10, 9, 9, 9, 9, 9, 5),
		res(3, 1, U21, 45, 45), res(4, 1, Senior, 30, 30),
	}
	ssv := ClassificationDefinition{Name: "SSV", TieBreak: SSVTieBreak, Parts: []ClassificationPart{{ProgrammeID: 1, Count: 1}},
		Classes: []ClassificationClass{{Name: "Adults", Categories: []AgeCategory{Elite, Senior}}, {Name: "Juniors", Categories: []AgeCategory{U21}}}}
	got := Classify(ssv, results)
	if len(got) != 2 || got[0].Name != "Adults" || fmt.Sprint(order(got[0].Ranked)) != "[2 1 4]" {
		t.Fatalf("series backwards: last series 50 beats 40: %+v", got)
	}
	if o := order(got[1].Ranked); len(o) != 1 || o[0] != 3 {
		t.Fatalf("juniors: %v", o)
	}
	ssv.TieBreak = []TieBreak{TieShots}
	got = Classify(ssv, results)
	if o := order(got[0].Ranked); o[0] != 1 {
		t.Fatalf("five tens beat four: %v", o)
	}
}

func TestClassificationValidate(t *testing.T) {
	c := ClassificationDefinition{Parts: []ClassificationPart{{ProgrammeID: 1, Count: 1}, {ProgrammeID: 1, Count: 0}, {ProgrammeID: 9, Count: 1}}, Best: 4,
		TieBreak: []TieBreak{TieAge, TieAge, "luck"}, Classes: []ClassificationClass{{Name: "A", Categories: []AgeCategory{Elite}}, {Categories: []AgeCategory{Elite}}}}
	got := Join(c.Validate(map[int64]bool{1: true}), "; ")
	for _, want := range []string{"Name is required", "chosen twice", "Programme 9", "between 0 and 3", `"age"`, `"luck"`, "Class 2: name", "more than one class"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}
}
