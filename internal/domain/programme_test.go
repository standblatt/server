package domain

import (
	"strings"
	"testing"
)

func valid() ProgrammeDefinition {
	return ProgrammeDefinition{
		Name:   "Field match",
		Series: []Series{{Shots: 6, MaxPoints: 60}, {Shots: 6, MaxPoints: 60, TrackMouches: true}},
		Prices: []Price{{AgeCategory: Elite, Price: 20}, {AgeCategory: U21, Price: 0}},
		Awards: []AwardTier{{Name: "Wreath", Kind: AwardCard, Value: 25, Limits: []AwardLimit{{AgeCategory: Elite, Points: 100}}}},
	}
}

func TestProgrammeTotals(t *testing.T) {
	d := valid()
	if d.Shots() != 12 || d.MaxPoints() != 120 {
		t.Fatalf("totals: %d shots, %d points", d.Shots(), d.MaxPoints())
	}
	if errs := d.Validate(); len(errs) != 0 {
		t.Fatalf("valid definition rejected: %v", errs)
	}
}

func TestProgrammeValidate(t *testing.T) {
	cases := map[string]func(*ProgrammeDefinition){
		"Name is required":          func(d *ProgrammeDefinition) { d.Name = "" },
		"At least one series":       func(d *ProgrammeDefinition) { d.Series = nil },
		"Series 2":                  func(d *ProgrammeDefinition) { d.Series[1].MaxPoints = 0 },
		"Unknown age category":      func(d *ProgrammeDefinition) { d.Prices[0].AgeCategory = "kids" },
		"more than one price":       func(d *ProgrammeDefinition) { d.Prices[1].AgeCategory = Elite },
		"must not be negative":      func(d *ProgrammeDefinition) { d.Prices[0].Price = -1 },
		"kind must be":              func(d *ProgrammeDefinition) { d.Awards[0].Kind = "gold" },
		"must be between 0 and 120": func(d *ProgrammeDefinition) { d.Awards[0].Limits[0].Points = 121 },
		"more than one limit":       func(d *ProgrammeDefinition) { d.Awards[0].Limits = append(d.Awards[0].Limits, d.Awards[0].Limits[0]) },
		"Award 1: name is required": func(d *ProgrammeDefinition) { d.Awards[0].Name = "" },
	}
	for want, mutate := range cases {
		d := valid()
		mutate(&d)
		errs := d.Validate()
		if len(errs) == 0 || !strings.Contains(Join(errs, "|"), want) {
			t.Errorf("%s: got %v", want, errs)
		}
	}
}
