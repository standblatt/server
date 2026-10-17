package domain

import (
	"reflect"
	"testing"
)

func TestRank(t *testing.T) {
	entries := []RankingEntry{
		0: {Total: 370, Series: []float64{95, 90, 90, 95}, Birthday: "1980-01-01", Category: Elite},
		1: {Total: 380, Series: []float64{95, 95, 95, 95}, Birthday: "1980-01-01", Category: Elite},
		2: {Total: 370, Series: []float64{90, 95, 95, 90}, Birthday: "1980-01-01", Category: Elite},
		3: {Total: 370.001, Series: []float64{95, 90, 90, 95}, Mouches: 3, Birthday: "1990-01-01", Category: Elite},
		4: {Total: 360, Series: []float64{90, 90, 90, 90}, Birthday: "2008-05-05", Category: U21},
		5: {Total: 360, Series: []float64{90, 90, 90, 90}, Birthday: "2010-05-05", Category: U17},
		6: {Total: 360, Series: []float64{90, 90, 90, 90}, Birthday: "1950-01-01", Category: SeniorVeteran},
		7: {Total: 360, Series: []float64{90, 90, 90, 90}, Birthday: "1970-01-01", Category: Senior},
		8: {Total: 350, Series: []float64{80, 90}, Birthday: "", Category: Elite},
		9: {Total: 350, Series: []float64{90, 80}, Birthday: "1960-01-01", Category: Veteran},
	}
	want := []int{
		1,    // highest total
		3, 0, // equal within tolerance and same series; more mouches first
		2,    // last series 90 < 95
		6, 7, // non-juniors before juniors, older first
		5, 4, // juniors: younger first
		8, 9, // last series 90 vs 80 decides before the unknown birthday matters
	}
	if got := Rank(entries); !reflect.DeepEqual(got, want) {
		t.Fatalf("rank order: got %v, want %v", got, want)
	}
}

func TestRankComparesSeriesFromTheLast(t *testing.T) {
	entries := []RankingEntry{
		{Total: 100, Series: []float64{50, 50}},
		{Total: 100, Series: []float64{100}},
	}
	if got := Rank(entries); !reflect.DeepEqual(got, []int{1, 0}) {
		t.Fatalf("got %v", got)
	}
}
