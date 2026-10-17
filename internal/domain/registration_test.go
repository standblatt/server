package domain

import (
	"strings"
	"testing"
)

func TestAgeCategoryFor(t *testing.T) {
	cases := []struct {
		birthday string
		want     AgeCategory
	}{
		{"", Elite},
		{"2017-12-31", U10},
		{"2016-01-01", U13},
		{"2014-06-01", U13},
		{"2013-06-01", U15},
		{"2012-06-01", U15},
		{"2011-06-01", U17},
		{"2009-06-01", U19},
		{"2007-06-01", U21},
		{"2006-01-01", U21},
		{"2005-12-31", Elite},
		{"1981-06-01", Elite},
		{"1980-06-01", Senior},
		{"1967-06-01", Senior},
		{"1966-06-01", Veteran},
		{"1957-06-01", Veteran},
		{"1956-06-01", SeniorVeteran},
	}
	for _, c := range cases {
		if got := AgeCategoryFor(c.birthday, "2026-11-20"); got != c.want {
			t.Errorf("born %s: got %s, want %s", c.birthday, got, c.want)
		}
	}
}

func TestCheckProgrammes(t *testing.T) {
	target := int64(2)
	programmes := map[int64]ProgrammeInfo{
		1: {ID: 1, Name: "Field match", Variant: "Stgw 90", Active: true, Series: 2, Categories: []AgeCategory{Elite, U21}},
		2: {ID: 2, Name: "Main", Active: true, Series: 4, Categories: []AgeCategory{Elite}},
		3: {ID: 3, Name: "Field match", Variant: "Stgw 57", Active: true, Series: 2, Categories: []AgeCategory{Elite}},
		4: {ID: 4, Name: "Section", Active: true, Series: 2, Categories: []AgeCategory{Elite}, CombinableWithID: &target},
		5: {ID: 5, Name: "Old", Active: false, Series: 1, Categories: []AgeCategory{Elite}},
	}
	ok := []ProgrammeChoice{{ProgrammeID: 1}, {ProgrammeID: 2}, {ProgrammeID: 4, CombinedSeries: []int{1, 3}}}
	if errs := CheckProgrammes(ok, programmes, Elite, nil, nil, nil); len(errs) != 0 {
		t.Fatalf("valid choices rejected: %v", errs)
	}
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 5}}, programmes, Elite, map[int64]bool{5: true}, nil, nil); len(errs) != 0 {
		t.Fatalf("kept inactive programme rejected: %v", errs)
	}

	cases := map[string][]ProgrammeChoice{
		"at least one":                {},
		"does not belong":             {{ProgrammeID: 99}},
		"chosen twice":                {{ProgrammeID: 1}, {ProgrammeID: 1}},
		"not active":                  {{ProgrammeID: 5}},
		"not open to the age":         {{ProgrammeID: 2}},
		"variants of each other":      {{ProgrammeID: 1}, {ProgrammeID: 3}},
		"cannot be combined":          {{ProgrammeID: 1, CombinedSeries: []int{1, 2}}},
		"which is not chosen":         {{ProgrammeID: 4, CombinedSeries: []int{1, 2}}},
		"needs exactly 2 series":      {{ProgrammeID: 2}, {ProgrammeID: 4, CombinedSeries: []int{1}}},
		"invalid or chosen twice":     {{ProgrammeID: 2}, {ProgrammeID: 4, CombinedSeries: []int{1, 1}}},
		"series 5 of Main is invalid": {{ProgrammeID: 2}, {ProgrammeID: 4, CombinedSeries: []int{1, 5}}},
	}
	for want, choices := range cases {
		category := Elite
		if want == "not open to the age" {
			category = U21
		}
		errs := CheckProgrammes(choices, programmes, category, nil, nil, nil)
		if !strings.Contains(Join(errs, "|"), want) {
			t.Errorf("%s: got %v", want, errs)
		}
	}
}

func TestCheckProgrammesRanges(t *testing.T) {
	r1, r2 := int64(1), int64(2)
	programmes := map[int64]ProgrammeInfo{
		1: {ID: 1, Name: "Rifle", Active: true, Series: 2, Categories: []AgeCategory{Elite}, RangeID: &r1},
		2: {ID: 2, Name: "Pistol", Active: true, Series: 2, Categories: []AgeCategory{Elite}, RangeID: &r2},
		3: {ID: 3, Name: "Any", Active: true, Series: 2, Categories: []AgeCategory{Elite}},
	}
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 1}, {ProgrammeID: 3}}, programmes, Elite, nil, &r1, nil); len(errs) != 0 {
		t.Fatalf("rifle on range 1: %v", errs)
	}
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 1}}, programmes, Elite, nil, &r2, nil); !strings.Contains(Join(errs, "|"), "not shot on this range") {
		t.Fatalf("rifle on range 2: %v", errs)
	}
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 1}, {ProgrammeID: 2}}, programmes, Elite, nil, nil, nil); !strings.Contains(Join(errs, "|"), "different ranges") {
		t.Fatalf("rifle and pistol together: %v", errs)
	}
}

func TestCheckProgrammesLaneGroups(t *testing.T) {
	r, g25, g50 := int64(1), int64(10), int64(11)
	lanes25 := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	programmes := map[int64]ProgrammeInfo{
		1: {ID: 1, Name: "Precision", Active: true, Series: 2, Categories: []AgeCategory{Elite}, RangeID: &r, GroupID: &g50, GroupName: "50 m", GroupLanes: []int32{11, 12, 13, 14, 15}},
		2: {ID: 2, Name: "Rapid", Active: true, Series: 2, Categories: []AgeCategory{Elite}, RangeID: &r, GroupID: &g25, GroupName: "25 m", GroupLanes: lanes25},
	}
	lane := func(n int) *int { return &n }
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 1}}, programmes, Elite, nil, &r, lane(12)); len(errs) != 0 {
		t.Fatalf("50 m on lane 12: %v", errs)
	}
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 1}}, programmes, Elite, nil, &r, lane(3)); !strings.Contains(Join(errs, "|"), "shot on 50 m, not on lane 3") {
		t.Fatalf("50 m on lane 3: %v", errs)
	}
	if errs := CheckProgrammes([]ProgrammeChoice{{ProgrammeID: 1}, {ProgrammeID: 2}}, programmes, Elite, nil, nil, nil); !strings.Contains(Join(errs, "|"), "different lane groups") {
		t.Fatalf("25 m and 50 m together: %v", errs)
	}
}
