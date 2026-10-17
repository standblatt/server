package domain

import (
	"fmt"
	"testing"
	"time"
)

// slotText formats slots as "2027-03-13 08:00-08:30" lines for comparison.
func slotText(slots []Slot) []string {
	out := []string{}
	for _, s := range slots {
		out = append(out, s.Start.Format("2006-01-02 15:04")+"-"+s.End.Format("15:04"))
	}
	return out
}

func TestBookingSlots(t *testing.T) {
	zurich, _ := time.LoadLocation("Europe/Zurich")
	days := []string{"2027-03-12", "2027-03-13"}

	all := BookingSlots(days, nil, BookingSettings{}, zurich)
	if len(all) != 2*18 || all[0].Start.Format("15:04") != "08:00" || all[17].End.Format("15:04") != "17:00" {
		t.Fatalf("defaults: %v", slotText(all))
	}

	s := BookingSettings{
		From: "08:00", To: "12:00", BreakFrom: "10:00", BreakTo: "10:45", SlotMinutes: 60,
		Days: map[string]BookingHours{"2027-03-12": {From: "17:00", To: "19:30"}, "2027-04-01": {From: "01:00", To: "02:00"}},
	}
	got := fmt.Sprint(slotText(BookingSlots(days, nil, s, zurich)))
	want := "[2027-03-12 17:00-18:00 2027-03-12 18:00-19:00 2027-03-13 08:00-09:00 2027-03-13 09:00-10:00 2027-03-13 10:45-11:45]"
	if got != want {
		t.Fatalf("slots\n got %s\nwant %s", got, want)
	}

	// Summer time starts on 2027-03-28 at 02:00; slots keep their wall-clock times.
	dst := BookingSlots([]string{"2027-03-27", "2027-03-28"}, nil, BookingSettings{From: "08:00", To: "09:00"}, zurich)
	if len(dst) != 4 || dst[2].Start.Format("15:04 MST") != "08:00 CEST" || dst[0].Start.Format("MST") != "CET" {
		t.Fatalf("dst: %v", slotText(dst))
	}

	times := map[string][]BookingHours{"2027-03-12": {{From: "09:00", To: "10:00"}, {From: "17:00", To: "18:00"}}}
	got = fmt.Sprint(slotText(BookingSlots(days[:1], times, BookingSettings{SlotMinutes: 60}, zurich)))
	if got != "[2027-03-12 09:00-10:00 2027-03-12 17:00-18:00]" {
		t.Fatalf("shooting times: %s", got)
	}
}

func TestBookingSettingsAndDeadline(t *testing.T) {
	zurich, _ := time.LoadLocation("Europe/Zurich")
	days := []string{"2027-03-13", "2027-03-12"}
	if errs := ValidateBookingSettings(BookingSettings{}, days); len(errs) != 0 {
		t.Fatalf("defaults invalid: %v", errs)
	}
	bad := BookingSettings{From: "12:00", To: "08:00", BreakFrom: "10:00", Deadline: "soon",
		Days: map[string]BookingHours{"2027-03-12": {From: "18:00", To: "17:00"}, "2020-01-01": {From: "x"}}}
	if errs := ValidateBookingSettings(bad, days); len(errs) != 4 {
		t.Fatalf("want 4 problems, got %v", errs)
	}
	clean := CleanBookingSettings(bad, days)
	if _, ok := clean.Days["2020-01-01"]; ok || len(clean.Days) != 1 {
		t.Fatalf("clean: %v", clean.Days)
	}

	d := BookingDeadline(days, BookingSettings{}, zurich)
	if d.Format("2006-01-02 15:04") != "2027-02-27 00:00" {
		t.Fatalf("default deadline %v", d)
	}
	d = BookingDeadline(days, BookingSettings{Deadline: "2027-03-01"}, zurich)
	if d.Format("2006-01-02 15:04") != "2027-03-02 00:00" {
		t.Fatalf("deadline %v", d)
	}
}

func TestLanes(t *testing.T) {
	free := FreeLanes(10, []int32{3}, []int{1, 6})
	if fmt.Sprint(free) != "[2 4 5 7 8 9 10]" {
		t.Fatalf("free %v", free)
	}
	for _, c := range []struct {
		n    int
		want string
		ok   bool
	}{{2, "[4 5]", true}, {4, "[7 8 9 10]", true}, {5, "[2 4 5 7 8]", true}, {8, "[]", false}, {0, "[]", false}} {
		got, ok := AssignLanes(free, c.n)
		if ok != c.ok || (ok && fmt.Sprint(got) != c.want) {
			t.Errorf("assign %d: %v %v", c.n, got, ok)
		}
	}
}

func TestValidateBooking(t *testing.T) {
	zurich, _ := time.LoadLocation("Europe/Zurich")
	slots := BookingSlots([]string{"2027-03-13"}, nil, BookingSettings{}, zurich)
	counts := map[int64]int{7: 10, 8: 4}
	item := func(rangeID int64, slot, lanes int) BookingItem {
		return BookingItem{RangeID: rangeID, Start: slots[slot].Start, End: slots[slot].End, Lanes: lanes}
	}
	ok := []BookingItem{item(7, 0, 4), item(8, 0, 4), item(7, 1, 2)}
	shooter := BookingShooter{Firstname: "Ida", Lastname: "Muster", Birthday: "1990-05-01"}
	if errs := ValidateBooking(ok, []BookingShooter{shooter}, slots, counts, 0); len(errs) != 0 {
		t.Fatalf("valid request: %v", errs)
	}
	off := item(7, 0, 1)
	off.Start = off.Start.Add(time.Minute)
	bad := []BookingItem{item(9, 0, 1), off, item(7, 2, 1), item(7, 2, 1), item(8, 3, 5), item(7, 4, 0)}
	errs := ValidateBooking(bad, []BookingShooter{{Firstname: "Ida", Birthday: "1.5.1990"}}, slots, counts, 0)
	if len(errs) != 7 {
		t.Fatalf("want 7 problems, got %v", errs)
	}
	if errs := ValidateBooking(ok, nil, slots, counts, 3); len(errs) != 2 {
		t.Fatalf("max lanes: %v", errs)
	}
	if errs := ValidateBooking(nil, nil, slots, counts, 0); len(errs) != 1 {
		t.Fatalf("empty: %v", errs)
	}
	if !SameBookingItems(ok, []BookingItem{ok[2], ok[0], ok[1]}) || SameBookingItems(ok, ok[:2]) || SameBookingItems(ok, []BookingItem{ok[0], ok[1], item(7, 1, 3)}) {
		t.Fatal("SameBookingItems")
	}
}
