package domain

// SpecialPrizeRules award the best ranks of every age category of a programme with a special prize. A category
// needs MinCategory participants in the programme and a shooter's club MinClub participants in the same
// programme; ineligible shooters are skipped and the next one moves up.
type SpecialPrizeRules struct {
	Ranks       int `json:"ranks" minimum:"0" maximum:"10" doc:"Ranks per age category with a special prize; 0 turns special prizes off"`
	MinCategory int `json:"minCategory" minimum:"0" doc:"Participants an age category needs in the programme"`
	MinClub     int `json:"minClub" minimum:"0" doc:"Participants a shooter's club needs in the same programme"`
}

// PrizeCandidate is a ranked result of a programme; Club is empty for a shooter without club.
type PrizeCandidate struct {
	Category AgeCategory
	Club     string
}

// CategoryPrizes is the outcome of one age category: indexes into the candidates, in ranking order.
type CategoryPrizes struct {
	Category     AgeCategory
	Participants int
	Winners      []int
	Skipped      []int
}

// SpecialPrizes picks the prize winners from candidates given in ranking order (see Rank). Categories appear in
// AgeCategories order; a category below MinCategory has neither winners nor skipped shooters.
func SpecialPrizes(candidates []PrizeCandidate, r SpecialPrizeRules) []CategoryPrizes {
	clubs := map[string]int{}
	categories := map[AgeCategory][]int{}
	for i, c := range candidates {
		if c.Club != "" {
			clubs[c.Club]++
		}
		categories[c.Category] = append(categories[c.Category], i)
	}
	var out []CategoryPrizes
	for _, cat := range AgeCategories {
		members := categories[cat]
		if len(members) == 0 {
			continue
		}
		cp := CategoryPrizes{Category: cat, Participants: len(members)}
		if r.Ranks > 0 && len(members) >= r.MinCategory {
			for _, i := range members {
				if len(cp.Winners) == r.Ranks {
					break
				}
				if c := candidates[i]; c.Club != "" && clubs[c.Club] >= r.MinClub || c.Club == "" && r.MinClub == 0 {
					cp.Winners = append(cp.Winners, i)
				} else {
					cp.Skipped = append(cp.Skipped, i)
				}
			}
		}
		out = append(out, cp)
	}
	return out
}
