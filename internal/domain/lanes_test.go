package domain

import (
	"reflect"
	"testing"
)

func TestParseLanes(t *testing.T) {
	cases := map[string][]int{
		"5":            {5},
		"5, 7, 9":      {5, 7, 9},
		"1-4":          {1, 2, 3, 4},
		"1, 5, 7-10":   {1, 5, 7, 8, 9, 10},
		"5, 3, 5, 1,3": {1, 3, 5},
		" 2 - 3 ,":     {2, 3},
	}
	for in, want := range cases {
		got, err := ParseLanes(in)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("%q: got %v %v, want %v", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "abc", "0", "5-2", "-3", "1-x"} {
		if _, err := ParseLanes(bad); err == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
}

func TestFormatLanes(t *testing.T) {
	if got := FormatLanes([]int{1, 2, 3, 5, 7, 8, 10, 11, 12, 13}); got != "1-3, 5, 7, 8, 10-13" {
		t.Fatal(got)
	}
	if got := FormatLanes([]int{4}); got != "4" {
		t.Fatal(got)
	}
}
