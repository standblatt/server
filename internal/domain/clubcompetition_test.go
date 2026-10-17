package domain

import "testing"

func club(id int64, scores ...float64) []ClubScore {
	out := make([]ClubScore, len(scores))
	for i, s := range scores {
		out[i] = ClubScore{ClubID: id, Score: s}
	}
	return out
}

func TestClubCompetition(t *testing.T) {
	var scores []ClubScore
	// 12 participants: required = max(5, 6) = 6; (6×190 + 0.02 × 6×100) / 6 = 192
	scores = append(scores, club(1, 190, 190, 190, 190, 190, 190, 100, 100, 100, 100, 100, 100)...)
	// 5 participants: (180+180+180+180+181)/5 = 180.2
	scores = append(scores, club(2, 181, 180, 180, 180, 180)...)
	// 6 participants: required 5, (5×180 + 0.02×100) / 5 = 180.4
	scores = append(scores, club(3, 180, 180, 180, 180, 180, 100)...)
	// excluded, best result
	scores = append(scores, club(4, 200, 200, 200, 200, 200)...)
	// unranked and a zero score that does not count
	scores = append(scores, club(5, 150, 150, 150, 150, 0)...)

	ranked, unranked := ClubCompetition(scores, 4)
	type row struct {
		id     int64
		rank   int
		result float64
	}
	want := []row{{4, 0, 200}, {1, 1, 192}, {3, 2, 180.4}, {2, 3, 180.2}}
	if len(ranked) != len(want) {
		t.Fatalf("ranked: %+v", ranked)
	}
	for i, w := range want {
		g := ranked[i]
		if g.ClubID != w.id || g.Rank != w.rank || g.Result != w.result {
			t.Errorf("position %d: got club %d rank %d result %v, want %+v", i, g.ClubID, g.Rank, g.Result, w)
		}
	}
	if ranked[1].Required != 6 || ranked[1].Participants != 12 {
		t.Errorf("club 1: %+v", ranked[1])
	}
	if len(unranked) != 1 || unranked[0].ClubID != 5 || unranked[0].Participants != 4 {
		t.Fatalf("unranked: %+v", unranked)
	}
}

func TestClubCompetitionTruncatesAndBreaksTies(t *testing.T) {
	// (100+100+100+100+101)/5 = 100.2; 0.02×(1) adds 0.004 → 100.204; with 7 participants required is 5
	ranked, _ := ClubCompetition(club(1, 101, 100, 100, 100, 100, 0.1, 0.1), 0)
	if ranked[0].Result != 100.2 {
		t.Fatalf("truncation: %v", ranked[0].Result)
	}
	a := club(1, 100, 100, 100, 100, 100)
	b := club(2, 100, 100, 100, 100, 100)
	b[4].Score, b[0].Score = 99, 101
	ranked, _ = ClubCompetition(append(a, b...), 0)
	if ranked[0].ClubID != 2 {
		t.Fatalf("equal result and size: the higher best score wins, got %+v", ranked)
	}
}
