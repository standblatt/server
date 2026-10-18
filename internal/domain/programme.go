package domain

import (
	"strings"
)

// AwardKind is what a shooter receives for reaching an award tier.
type AwardKind string

// Award kinds: a prize in kind, or a printed wreath card worth a cash value.
const (
	AwardNatural AwardKind = "natural"
	AwardCard    AwardKind = "card"
	// AwardCash is a prize paid out in money ("Auszahler"); Value is the amount.
	AwardCash AwardKind = "cash"
)

// Series is one series (Passe) of a programme.
type Series struct {
	Shots        int    `json:"shots" minimum:"1" maximum:"100"`
	MaxPoints    int    `json:"maxPoints" minimum:"1"`
	Name         string `json:"name,omitempty" required:"false"`
	TrackMouches bool   `json:"trackMouches,omitempty" required:"false"`
	TrackShots   bool   `json:"trackShots,omitempty" required:"false" doc:"Every shot can be recorded; the series points are their sum"`
}

// Price is what a shooter of an age category pays; categories without a price may not shoot the programme.
type Price struct {
	AgeCategory AgeCategory `json:"ageCategory" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
	Price       float64     `json:"price" minimum:"0"`
	RepeatPrice *float64    `json:"repeatPrice,omitempty" required:"false" minimum:"0" doc:"Price of every further attempt (Nachdoppel); left out = the same price, 0 = paid once"`
}

// AttemptRule decides which attempt counts when a shooter shoots a programme several times.
type AttemptRule string

// Attempt rules: the best attempt by the ranking rules, the first or the last one.
const (
	AttemptBest  AttemptRule = "best"
	AttemptFirst AttemptRule = "first"
	AttemptLast  AttemptRule = "last"
)

// AwardLimit is the total a shooter of an age category needs for an award tier.
type AwardLimit struct {
	AgeCategory AgeCategory `json:"ageCategory" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
	Points      float64     `json:"points" minimum:"0"`
}

// AwardTier is one award level (Kranz) of a programme.
type AwardTier struct {
	Name   string       `json:"name" minLength:"1" maxLength:"100"`
	Kind   AwardKind    `json:"kind" enum:"natural,card,cash"`
	Value  float64      `json:"value" minimum:"0" doc:"Value in CHF"`
	Limits []AwardLimit `json:"limits"`
}

// ProgrammeDefinition is everything that defines a programme. It is the editor's input, the export format
// and the content of library templates. Programmes with the same name in one event are variants of each
// other: a shooter registers for at most one of them.
type ProgrammeDefinition struct {
	Name                     string      `json:"name" minLength:"1" maxLength:"100"`
	Variant                  string      `json:"variant" maxLength:"100" required:"false"`
	Description              string      `json:"description" required:"false"`
	Barcode                  string      `json:"barcode" maxLength:"50" required:"false"`
	CountsForClubCompetition bool        `json:"countsForClubCompetition" required:"false"`
	ShowOnRankings           bool        `json:"showOnRankings" required:"false"`
	Attempts                 int         `json:"attempts,omitempty" required:"false" minimum:"0" maximum:"1000" doc:"How often a shooter may shoot the programme at the event; left out or 0 = once"`
	UnlimitedAttempts        bool        `json:"unlimitedAttempts,omitempty" required:"false" doc:"Any number of attempts; Attempts is ignored"`
	AttemptRule              AttemptRule `json:"attemptRule,omitempty" required:"false" enum:"best,first,last" doc:"Which attempt counts in the ranking; default best"`
	AwardCountingOnly        bool        `json:"awardCountingOnly,omitempty" required:"false" doc:"Only the counting attempt receives an award, not every attempt"`
	Series                   []Series    `json:"series"`
	Prices                   []Price     `json:"prices"`
	Awards                   []AwardTier `json:"awards"`
}

// Shots is the number of shots over all series.
func (d ProgrammeDefinition) Shots() int {
	n := 0
	for _, s := range d.Series {
		n += s.Shots
	}
	return n
}

// MaxPoints is the highest possible total over all series.
func (d ProgrammeDefinition) MaxPoints() int {
	n := 0
	for _, s := range d.Series {
		n += s.MaxPoints
	}
	return n
}

// AttemptLimit is how many attempts a shooter may shoot; 0 means unlimited.
func (d ProgrammeDefinition) AttemptLimit() int {
	if d.UnlimitedAttempts {
		return 0
	}
	return max(d.Attempts, 1)
}

// Rule is the attempt rule, best when none is set.
func (d ProgrammeDefinition) Rule() AttemptRule {
	if d.AttemptRule == "" {
		return AttemptBest
	}
	return d.AttemptRule
}

// Normalize trims text fields, fills the attempt defaults and turns nil lists into empty ones.
func (d *ProgrammeDefinition) Normalize() {
	if d.UnlimitedAttempts {
		d.Attempts = 0
	} else {
		d.Attempts = max(d.Attempts, 1)
	}
	d.AttemptRule = d.Rule()
	d.Name = strings.TrimSpace(d.Name)
	d.Variant = strings.TrimSpace(d.Variant)
	d.Barcode = strings.TrimSpace(d.Barcode)
	if d.Series == nil {
		d.Series = []Series{}
	}
	if d.Prices == nil {
		d.Prices = []Price{}
	}
	if d.Awards == nil {
		d.Awards = []AwardTier{}
	}
	for i := range d.Awards {
		d.Awards[i].Name = strings.TrimSpace(d.Awards[i].Name)
		if d.Awards[i].Limits == nil {
			d.Awards[i].Limits = []AwardLimit{}
		}
	}
}

// Validate returns every problem with the definition; an empty result means it is usable.
func (d ProgrammeDefinition) Validate() []Problem {
	var errs []Problem
	if d.Name == "" {
		errs = append(errs, Problemf("Name is required"))
	}
	if len(d.Series) == 0 {
		errs = append(errs, Problemf("At least one series is required"))
	}
	for i, s := range d.Series {
		if s.Shots < 1 || s.MaxPoints < 1 {
			errs = append(errs, Problemf("Series %d: shots and maximum points must be positive", i+1))
		}
	}
	if d.AttemptRule != "" && d.AttemptRule != AttemptBest && d.AttemptRule != AttemptFirst && d.AttemptRule != AttemptLast {
		errs = append(errs, Problemf("Unknown attempt rule %q", d.AttemptRule))
	}
	if d.Attempts < 0 || d.Attempts > 1000 {
		errs = append(errs, Problemf("Attempts must be between 1 and 1000"))
	}
	seen := map[AgeCategory]bool{}
	for _, p := range d.Prices {
		switch {
		case !p.AgeCategory.Valid():
			errs = append(errs, Problemf("Unknown age category %q", p.AgeCategory))
		case seen[p.AgeCategory]:
			errs = append(errs, Problemf("Age category %s has more than one price", p.AgeCategory))
		case p.Price < 0 || (p.RepeatPrice != nil && *p.RepeatPrice < 0):
			errs = append(errs, Problemf("Price for %s must not be negative", p.AgeCategory))
		}
		seen[p.AgeCategory] = true
	}
	max := float64(d.MaxPoints())
	for i, a := range d.Awards {
		if a.Name == "" {
			errs = append(errs, Problemf("Award %d: name is required", i+1))
		}
		if a.Kind != AwardNatural && a.Kind != AwardCard && a.Kind != AwardCash {
			errs = append(errs, Problemf("Award %d: kind must be natural, card or cash", i+1))
		}
		if a.Value < 0 {
			errs = append(errs, Problemf("Award %d: value must not be negative", i+1))
		}
		limits := map[AgeCategory]bool{}
		for _, l := range a.Limits {
			switch {
			case !l.AgeCategory.Valid():
				errs = append(errs, Problemf("Award %d: unknown age category %q", i+1, l.AgeCategory))
			case limits[l.AgeCategory]:
				errs = append(errs, Problemf("Award %d: age category %s has more than one limit", i+1, l.AgeCategory))
			case l.Points < 0 || (max > 0 && l.Points > max):
				errs = append(errs, Problemf("Award %d: limit for %s must be between 0 and %.0f", i+1, l.AgeCategory, max))
			}
			limits[l.AgeCategory] = true
		}
	}
	return errs
}
