package domain

import (
	"math"
	"sort"
)

// ClubMinParticipants is the smallest number of results a club needs to be ranked.
const ClubMinParticipants = 5

// ClubScore is one result counting for the club competition: the first two series of a programme.
type ClubScore struct {
	ClubID int64
	Score  float64
}

// ClubResult is a club in the club competition.
type ClubResult struct {
	ClubID       int64     `json:"clubId"`
	Rank         int       `json:"rank,omitempty" doc:"0 for clubs out of competition"`
	Excluded     bool      `json:"excluded"`
	Participants int       `json:"participants"`
	Required     int       `json:"required"`
	Result       float64   `json:"result"`
	Scores       []float64 `json:"scores"`
}

// ClubCompetition applies the SSV club competition rules. Only scores above zero count; each one is a
// participant. Clubs with fewer than five participants are unranked. Of n participants the best
// max(5, floor(n/2)) are mandatory and the rest optional; the club result is
// (sum of mandatory + 0.02 × sum of optional) / mandatory count, truncated to three decimals. Ties are broken by
// more participants, then by comparing the sorted scores one by one. The excluded club keeps its place but
// gets no rank.
func ClubCompetition(scores []ClubScore, excludedClub int64) (ranked, unranked []ClubResult) {
	byClub := map[int64][]float64{}
	var order []int64
	for _, s := range scores {
		if s.Score <= 0 {
			continue
		}
		if _, ok := byClub[s.ClubID]; !ok {
			order = append(order, s.ClubID)
		}
		byClub[s.ClubID] = append(byClub[s.ClubID], s.Score)
	}
	for _, id := range order {
		list := byClub[id]
		sort.Sort(sort.Reverse(sort.Float64Slice(list)))
		r := ClubResult{ClubID: id, Participants: len(list), Scores: list, Excluded: id == excludedClub}
		if len(list) < ClubMinParticipants {
			unranked = append(unranked, r)
			continue
		}
		r.Required = max(ClubMinParticipants, len(list)/2)
		var mandatory, optional float64
		for i, s := range list {
			if i < r.Required {
				mandatory += s
			} else {
				optional += s
			}
		}
		r.Result = math.Floor((mandatory+0.02*optional)/float64(r.Required)*1000+1e-9) / 1000
		ranked = append(ranked, r)
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.Result != b.Result {
			return a.Result > b.Result
		}
		if a.Participants != b.Participants {
			return a.Participants > b.Participants
		}
		for k := 0; k < len(a.Scores) && k < len(b.Scores); k++ {
			if a.Scores[k] != b.Scores[k] {
				return a.Scores[k] > b.Scores[k]
			}
		}
		return false
	})
	rank := 0
	for i := range ranked {
		if ranked[i].Excluded {
			continue
		}
		rank++
		ranked[i].Rank = rank
	}
	return ranked, unranked
}
