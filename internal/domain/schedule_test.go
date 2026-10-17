package domain

import (
	"fmt"
	"testing"
)

func TestScheduleExpand(t *testing.T) {
	evening := []BookingHours{{From: "17:30", To: "19:30"}}
	s := Schedule{
		From: "2027-05-24", To: "2027-06-06",
		Weekly: []WeeklyTimes{{Weekday: 3, Times: evening}, {Weekday: 6}},
		Days:   []ScheduleDay{{Date: "2027-05-28", Times: []BookingHours{{From: "17:00", To: "20:00"}}}, {Date: "2027-05-29"}},
	}
	got := fmt.Sprint(s.Dates())
	if got != "[2027-05-26 2027-05-28 2027-05-29 2027-06-02 2027-06-05]" {
		t.Fatalf("dates %s", got)
	}
	times := s.Times()
	if len(times) != 3 || times["2027-05-28"][0].From != "17:00" || times["2027-06-02"][0].From != "17:30" {
		t.Fatalf("times %v", times)
	}
	if n := len(Schedule{From: "2027-04-01", To: "2027-08-31"}.Expand()); n != 153 {
		t.Fatalf("every day: %d", n)
	}
	if got := fmt.Sprint(s.Shift(364).Dates()); got != "[2028-05-24 2028-05-26 2028-05-27 2028-05-31 2028-06-03]" {
		t.Fatalf("shifted %s", got)
	}
}

func TestScheduleValidateNormalize(t *testing.T) {
	bad := Schedule{
		From: "2027-06-01", To: "2027-05-01",
		Weekly: []WeeklyTimes{{Weekday: 8}, {Weekday: 2, Times: []BookingHours{{From: "19:00", To: "18:00"}}}},
		Days:   []ScheduleDay{{Date: "soon"}},
	}
	if errs := ValidateSchedule(bad); len(errs) != 4 {
		t.Fatalf("errors %v", errs)
	}
	if errs := ValidateSchedule(Schedule{Weekly: []WeeklyTimes{{Weekday: 1}}}); len(errs) != 1 {
		t.Fatalf("weekly without period %v", errs)
	}
	n := NormalizeSchedule(Schedule{Days: []ScheduleDay{{Date: "2027-01-10"}, {Date: "2027-01-08"}, {Date: "2027-01-10"}}}, "2027-01-09")
	if got := fmt.Sprint(n.Dates()); got != "[2027-01-08 2027-01-09 2027-01-10]" || len(n.Days) != 3 {
		t.Fatalf("normalized %s", got)
	}
	n = NormalizeSchedule(Schedule{From: "2027-04-01", To: "2027-08-31"}, "2027-05-29")
	if len(n.Days) != 0 {
		t.Fatalf("main day inside the period added: %v", n.Days)
	}
}
