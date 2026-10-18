package domain

import (
	"slices"
	"strconv"
)

// AgeCategoryFor derives the SSV age category from the birth year and the year the event takes place
// (SSV: competition year minus year of birth, the day of the competition never matters).
// Both dates are YYYY-MM-DD; an unknown birthday counts as Elite.
func AgeCategoryFor(birthday, eventDate string) AgeCategory {
	if len(birthday) < 4 || len(eventDate) < 4 {
		return Elite
	}
	born, err1 := strconv.Atoi(birthday[:4])
	year, err2 := strconv.Atoi(eventDate[:4])
	if err1 != nil || err2 != nil {
		return Elite
	}
	switch age := year - born; {
	case age <= 9:
		return U10
	case age <= 12:
		return U13
	case age <= 14:
		return U15
	case age <= 16:
		return U17
	case age <= 18:
		return U19
	case age <= 20:
		return U21
	case age <= 45:
		return Elite
	case age <= 59:
		return Senior
	case age <= 69:
		return Veteran
	default:
		return SeniorVeteran
	}
}

// ProgrammeChoice is one programme a shooter registers for. CombinedSeries is set when the programme takes
// these series (numbers, 1-based) of the programme it is combinable with. Count is the number of attempts
// shot in this registration (default 1).
type ProgrammeChoice struct {
	ProgrammeID    int64 `json:"programmeId"`
	CombinedSeries []int `json:"combinedSeries,omitempty" required:"false"`
	Count          int   `json:"count,omitempty" required:"false" minimum:"0" maximum:"100" doc:"Attempts of the programme in this registration; left out = 1"`
}

// Attempts is the number of attempts of the choice, at least one.
func (c ProgrammeChoice) Attempts() int {
	return max(c.Count, 1)
}

// ProgrammeInfo is what the registration rules need to know about a programme of the event.
type ProgrammeInfo struct {
	ID               int64
	Name             string
	Variant          string
	Active           bool
	Series           int
	Categories       []AgeCategory
	CombinableWithID *int64
	RangeID          *int64
	AttemptLimit     int
}

// title returns the programme name with its variant in parentheses, if any.
func (p ProgrammeInfo) title() string {
	if p.Variant == "" {
		return p.Name
	}
	return p.Name + " (" + p.Variant + ")"
}

// CheckProgrammes validates the programme choices of a registration and returns every problem.
// keep lists programmes the registration already has; they stay allowed even if they were deactivated since.
// rangeID is the range of the registration's lane; programmes tied to a range must match it and each other.
func CheckProgrammes(choices []ProgrammeChoice, programmes map[int64]ProgrammeInfo, category AgeCategory, keep map[int64]bool, rangeID *int64) []Problem {
	var errs []Problem
	if len(choices) == 0 {
		return []Problem{Problemf("Choose at least one programme")}
	}
	chosen := map[int64]bool{}
	byName := map[string]string{}
	var onRange *ProgrammeInfo
	for _, c := range choices {
		p, ok := programmes[c.ProgrammeID]
		if !ok {
			errs = append(errs, Problemf("Programme %d does not belong to this event", c.ProgrammeID))
			continue
		}
		if chosen[c.ProgrammeID] {
			errs = append(errs, Problemf("%s is chosen twice", p.title()))
			continue
		}
		chosen[c.ProgrammeID] = true
		if c.Attempts() > 1 && p.AttemptLimit == 1 {
			errs = append(errs, Problemf("%s can be shot only once", p.title()))
		} else if p.AttemptLimit > 0 && c.Attempts() > p.AttemptLimit {
			errs = append(errs, Problemf("%s can be shot at most %d times", p.title(), p.AttemptLimit))
		}
		if !p.Active && !keep[p.ID] {
			errs = append(errs, Problemf("%s is not active", p.title()))
		}
		if !slices.Contains(p.Categories, category) {
			errs = append(errs, Problemf("%s is not open to the age category %s", p.title(), category))
		}
		if other, dup := byName[p.Name]; dup {
			errs = append(errs, Problemf("%s and %s are variants of each other; choose one", other, p.title()))
		}
		byName[p.Name] = p.title()
		if p.RangeID != nil {
			if rangeID != nil && *rangeID != *p.RangeID {
				errs = append(errs, Problemf("%s is not shot on this range", p.title()))
			}
			if onRange != nil && *onRange.RangeID != *p.RangeID {
				errs = append(errs, Problemf("%s and %s are shot on different ranges", onRange.title(), p.title()))
			}
			q := p
			onRange = &q
		}
	}
	for _, c := range choices {
		p, ok := programmes[c.ProgrammeID]
		if !ok || len(c.CombinedSeries) == 0 {
			continue
		}
		if p.CombinableWithID == nil {
			errs = append(errs, Problemf("%s cannot be combined", p.title()))
			continue
		}
		if c.Attempts() > 1 {
			errs = append(errs, Problemf("%s takes its series from another programme and can be shot only once", p.title()))
			continue
		}
		target := programmes[*p.CombinableWithID]
		if !chosen[target.ID] {
			errs = append(errs, Problemf("%s takes its series from %s, which is not chosen", p.title(), target.title()))
			continue
		}
		if len(c.CombinedSeries) != p.Series {
			errs = append(errs, Problemf("%s needs exactly %d series of %s", p.title(), p.Series, target.title()))
		}
		seen := map[int]bool{}
		for _, n := range c.CombinedSeries {
			if n < 1 || n > target.Series || seen[n] {
				errs = append(errs, Problemf("%s: series %d of %s is invalid or chosen twice", p.title(), n, target.title()))
			}
			seen[n] = true
		}
	}
	return errs
}
