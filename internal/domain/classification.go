package domain

import (
	"slices"
	"sort"
	"strings"
)

// ClassificationPart is a programme of a classification and how many of a shooter's attempts at it count,
// chosen by the programme's attempt rule.
type ClassificationPart struct {
	ProgrammeID int64 `json:"programmeId"`
	Count       int   `json:"count" minimum:"1" maximum:"100" doc:"Attempts of the programme that count"`
}

// ClassificationClass ranks the shooters of some age categories on their own.
type ClassificationClass struct {
	Name       string        `json:"name" minLength:"1" maxLength:"100"`
	Categories []AgeCategory `json:"categories" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
}

// ClassificationDefinition ranks shooters over one or more programmes of an event, e.g. a club championship over the best
// three attempts of an annual programme or a section championship as the sum of several programmes.
type ClassificationDefinition struct {
	Name        string                `json:"name" minLength:"1" maxLength:"100"`
	Parts       []ClassificationPart  `json:"parts"`
	Best        int                   `json:"best,omitempty" required:"false" minimum:"0" doc:"Only the best parts count, this many; 0 = all"`
	Complete    bool                  `json:"complete,omitempty" required:"false" doc:"Only shooters with every counting part complete are ranked; the others are listed unranked"`
	TieBreak    []TieBreak            `json:"tieBreak" enum:"series,shots,mouches,age" doc:"Criteria after the total, in order"`
	Classes     []ClassificationClass `json:"classes" doc:"Separate rankings by age categories; empty = one ranking"`
	ShowOnKiosk bool                  `json:"showOnKiosk,omitempty" required:"false"`
}

// Normalize trims the names and turns nil lists into empty ones.
func (c *ClassificationDefinition) Normalize() {
	c.Name = strings.TrimSpace(c.Name)
	if c.Parts == nil {
		c.Parts = []ClassificationPart{}
	}
	if c.TieBreak == nil {
		c.TieBreak = []TieBreak{}
	}
	if c.Classes == nil {
		c.Classes = []ClassificationClass{}
	}
	for i := range c.Classes {
		c.Classes[i].Name = strings.TrimSpace(c.Classes[i].Name)
		if c.Classes[i].Categories == nil {
			c.Classes[i].Categories = []AgeCategory{}
		}
	}
}

// Validate returns every problem with the classification; programmes lists the event's programme IDs.
func (c ClassificationDefinition) Validate(programmes map[int64]bool) []Problem {
	var errs []Problem
	if c.Name == "" {
		errs = append(errs, Problemf("Name is required"))
	}
	if len(c.Parts) == 0 {
		errs = append(errs, Problemf("Choose at least one programme"))
	}
	seen := map[int64]bool{}
	for i, p := range c.Parts {
		switch {
		case !programmes[p.ProgrammeID]:
			errs = append(errs, Problemf("Programme %d does not belong to this event", p.ProgrammeID))
		case seen[p.ProgrammeID]:
			errs = append(errs, Problemf("Part %d: the programme is chosen twice", i+1))
		case p.Count < 1:
			errs = append(errs, Problemf("Part %d: at least one attempt must count", i+1))
		}
		seen[p.ProgrammeID] = true
	}
	if c.Best < 0 || c.Best > len(c.Parts) {
		errs = append(errs, Problemf("The best parts must be between 0 and %d", len(c.Parts)))
	}
	criteria := map[TieBreak]bool{}
	for _, t := range c.TieBreak {
		if !slices.Contains([]TieBreak{TieSeries, TieShots, TieMouches, TieAge}, t) || criteria[t] {
			errs = append(errs, Problemf("Unknown or repeated tie-break %q", t))
		}
		criteria[t] = true
	}
	inClass := map[AgeCategory]string{}
	for i, cl := range c.Classes {
		if cl.Name == "" {
			errs = append(errs, Problemf("Class %d: name is required", i+1))
		}
		for _, cat := range cl.Categories {
			switch {
			case !cat.Valid():
				errs = append(errs, Problemf("Class %d: unknown age category %q", i+1, cat))
			case inClass[cat] != "":
				errs = append(errs, Problemf("Age category %s is in more than one class", cat))
			}
			inClass[cat] = cl.Name
		}
	}
	return errs
}

// AttemptResult is a completed attempt of a shooter at a programme, as input of Classify.
type AttemptResult struct {
	ShooterID   int64
	ProgrammeID int64
	Rule        AttemptRule
	Entry       RankingEntry
}

// Standing is a shooter in a classification. Parts holds the points of every part (nil without a counting
// attempt) and Counted whether the part counts towards the total (see ClassificationDefinition.Best).
type Standing struct {
	ShooterID int64
	Rank      int
	Parts     []*float64
	Counted   []bool
	Complete  bool
	Entry     RankingEntry
}

// ClassStanding is the ranking of one class; Unranked are the shooters who are not complete.
type ClassStanding struct {
	Name     string
	Ranked   []Standing
	Unranked []Standing
}

// Classify ranks the shooters of a classification. Results must be in attempt order per shooter and programme.
// Every part takes the counting attempts by the programme's rule; with Best only that many parts with the highest
// points count; the total, series, shots and mouches of the counting attempts are compared by RankBy. Without
// classes there is one class named "".
func Classify(c ClassificationDefinition, results []AttemptResult) []ClassStanding {
	type key struct{ shooter, programme int64 }
	attempts := map[key][]RankingEntry{}
	rules := map[int64]AttemptRule{}
	seen := map[int64]bool{}
	var shooters []int64
	for _, r := range results {
		if !slices.ContainsFunc(c.Parts, func(p ClassificationPart) bool { return p.ProgrammeID == r.ProgrammeID }) {
			continue
		}
		if !seen[r.ShooterID] {
			seen[r.ShooterID] = true
			shooters = append(shooters, r.ShooterID)
		}
		k := key{r.ShooterID, r.ProgrammeID}
		attempts[k] = append(attempts[k], r.Entry)
		rules[r.ProgrammeID] = r.Rule
	}
	need := len(c.Parts)
	if c.Best > 0 {
		need = c.Best
	}
	var all []Standing
	for _, id := range shooters {
		st := Standing{ShooterID: id, Parts: make([]*float64, len(c.Parts)), Counted: make([]bool, len(c.Parts))}
		picked := make([][]RankingEntry, len(c.Parts))
		full := 0
		for i, p := range c.Parts {
			list := attempts[key{id, p.ProgrammeID}]
			if len(list) == 0 {
				continue
			}
			sum := 0.0
			for _, n := range Counting(rules[p.ProgrammeID], list, p.Count) {
				picked[i] = append(picked[i], list[n])
				sum += list[n].Total
			}
			sum = Round2(sum)
			st.Parts[i] = &sum
			if len(picked[i]) >= p.Count {
				full++
			}
		}
		order := make([]int, 0, len(c.Parts))
		for i := range c.Parts {
			if st.Parts[i] != nil {
				order = append(order, i)
			}
		}
		sort.SliceStable(order, func(a, b int) bool { return *st.Parts[order[a]] > *st.Parts[order[b]] })
		if c.Best > 0 && len(order) > c.Best {
			order = order[:c.Best]
		}
		slices.Sort(order)
		for _, i := range order {
			st.Counted[i] = true
			for _, e := range picked[i] {
				st.Entry.Total += e.Total
				st.Entry.Series = append(st.Entry.Series, e.Series...)
				st.Entry.Shots = append(st.Entry.Shots, e.Shots...)
				st.Entry.Mouches += e.Mouches
				st.Entry.Birthday, st.Entry.Category = e.Birthday, e.Category
			}
		}
		st.Entry.Total = Round2(st.Entry.Total)
		st.Complete = full >= need
		all = append(all, st)
	}

	classes := c.Classes
	if len(classes) == 0 {
		classes = []ClassificationClass{{Categories: AgeCategories}}
	}
	out := make([]ClassStanding, len(classes))
	for ci, cl := range classes {
		var ranked, unranked []Standing
		for _, st := range all {
			if !slices.Contains(cl.Categories, st.Entry.Category) {
				continue
			}
			if c.Complete && !st.Complete {
				unranked = append(unranked, st)
			} else {
				ranked = append(ranked, st)
			}
		}
		ranked = rankStandings(ranked, c.TieBreak)
		for i := range ranked {
			ranked[i].Rank = i + 1
		}
		out[ci] = ClassStanding{Name: cl.Name, Ranked: ranked, Unranked: rankStandings(unranked, c.TieBreak)}
	}
	return out
}

// rankStandings returns the standings in rank order by RankBy; never nil.
func rankStandings(list []Standing, criteria []TieBreak) []Standing {
	entries := make([]RankingEntry, len(list))
	for i, st := range list {
		entries[i] = st.Entry
	}
	out := make([]Standing, 0, len(list))
	for _, i := range RankBy(entries, criteria) {
		out = append(out, list[i])
	}
	return out
}
