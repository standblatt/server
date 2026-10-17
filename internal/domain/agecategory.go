// Package domain holds the competition rules that do not depend on storage or HTTP.
package domain

// AgeCategory is an SSV age category; the values match the age_category database enum.
type AgeCategory string

// Age categories from youngest to oldest.
const (
	U10           AgeCategory = "u10"
	U13           AgeCategory = "u13"
	U15           AgeCategory = "u15"
	U17           AgeCategory = "u17"
	U19           AgeCategory = "u19"
	U21           AgeCategory = "u21"
	Elite         AgeCategory = "elite"
	Senior        AgeCategory = "senior"
	Veteran       AgeCategory = "veteran"
	SeniorVeteran AgeCategory = "senior_veteran"
)

// AgeCategories lists all categories from youngest to oldest.
var AgeCategories = []AgeCategory{U10, U13, U15, U17, U19, U21, Elite, Senior, Veteran, SeniorVeteran}

// Valid reports whether c is a known category.
func (c AgeCategory) Valid() bool {
	for _, k := range AgeCategories {
		if c == k {
			return true
		}
	}
	return false
}

// Junior reports whether c is one of the youth categories U10 to U21.
func (c AgeCategory) Junior() bool {
	switch c {
	case U10, U13, U15, U17, U19, U21:
		return true
	}
	return false
}
