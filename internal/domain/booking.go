package domain

import (
	"slices"
	"strings"
	"time"
)

// Defaults of the online booking; an empty setting falls back to them.
const (
	BookingFrom         = "08:00"
	BookingTo           = "17:00"
	BookingSlotMinutes  = 30
	BookingDeadlineDays = 14
	MaxBookingShooters  = 200
	// BookingLeadTime is how long before its start a slot can still be asked for.
	BookingLeadTime = 2 * time.Hour
)

// BookingHours are the shooting hours of one day as HH:MM.
type BookingHours struct {
	From string `json:"from" pattern:"^[0-2][0-9]:[0-5][0-9]$"`
	To   string `json:"to" pattern:"^[0-2][0-9]:[0-5][0-9]$"`
}

// BookingSettings are the options of the online booking of an event; empty fields mean the default.
type BookingSettings struct {
	From        string                  `json:"from,omitempty" required:"false" pattern:"^([0-2][0-9]:[0-5][0-9])?$" doc:"Shooting hours of every day; default 08:00–17:00"`
	To          string                  `json:"to,omitempty" required:"false" pattern:"^([0-2][0-9]:[0-5][0-9])?$"`
	BreakFrom   string                  `json:"breakFrom,omitempty" required:"false" pattern:"^([0-2][0-9]:[0-5][0-9])?$" doc:"Optional break of every day"`
	BreakTo     string                  `json:"breakTo,omitempty" required:"false" pattern:"^([0-2][0-9]:[0-5][0-9])?$"`
	SlotMinutes int                     `json:"slotMinutes,omitempty" required:"false" minimum:"0" maximum:"480" doc:"Length of a slot; default 30"`
	Opens       string                  `json:"opens,omitempty" required:"false" doc:"First day clubs can ask (YYYY-MM-DD); empty = right away"`
	Deadline    string                  `json:"deadline,omitempty" required:"false" doc:"Last day clubs can ask or change (YYYY-MM-DD), at the latest the last day of the event; default 14 days before the first day"`
	MaxLanes    int                     `json:"maxLanes,omitempty" required:"false" minimum:"0" doc:"Lanes per slot and request; 0 = no limit"`
	Days        map[string]BookingHours `json:"days,omitempty" required:"false" doc:"Days with other hours, by date"`
}

// Slot is one bookable time window.
type Slot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// clock parses HH:MM into minutes after midnight.
func clock(s string) (int, bool) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, false
	}
	return t.Hour()*60 + t.Minute(), true
}

// orDefault returns v, or def when v is empty or zero.
func orDefault[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}

// ValidateBookingSettings checks the times of the settings and that the other days are days of the event.
func ValidateBookingSettings(s BookingSettings, days []string) []Problem {
	var errs []Problem
	valid := func(h BookingHours) bool {
		from, ok1 := clock(h.From)
		to, ok2 := clock(h.To)
		return ok1 && ok2 && from < to
	}
	if !valid(BookingHours{orDefault(s.From, BookingFrom), orDefault(s.To, BookingTo)}) {
		errs = append(errs, Problemf("Shooting hours: the end must be after the start"))
	}
	if (s.BreakFrom != "" || s.BreakTo != "") && !valid(BookingHours{s.BreakFrom, s.BreakTo}) {
		errs = append(errs, Problemf("Break: the end must be after the start"))
	}
	if s.Deadline != "" {
		if _, err := time.Parse(time.DateOnly, s.Deadline); err != nil {
			errs = append(errs, Problemf("Deadline: not a date"))
		} else if len(days) > 0 && s.Deadline > slices.Max(days) {
			errs = append(errs, Problemf("The deadline is after the last shooting day"))
		}
	}
	if s.Opens != "" {
		if opens, err := time.Parse(time.DateOnly, s.Opens); err != nil {
			errs = append(errs, Problemf("Booking opens: not a date"))
		} else if !opens.Before(BookingDeadline(days, s, time.UTC)) {
			errs = append(errs, Problemf("Booking opens after the deadline"))
		}
	}
	for day, h := range s.Days {
		if !slices.Contains(days, day) {
			continue
		}
		if !valid(h) {
			errs = append(errs, Problemf("Hours of %s: the end must be after the start", day))
		}
	}
	return errs
}

// CleanBookingSettings drops the other hours of days that are no longer days of the event.
func CleanBookingSettings(s BookingSettings, days []string) BookingSettings {
	for day := range s.Days {
		if !slices.Contains(days, day) {
			delete(s.Days, day)
		}
	}
	if len(s.Days) == 0 {
		s.Days = nil
	}
	return s
}

// BookingSlots returns the slots of every day of the event in loc: the day's hours cut into slots of the slot
// length, skipping the break; a rest shorter than a slot before the end or the break is dropped. A day's hours are
// its other hours of the settings, else its shooting times (one or more windows), else the default hours.
func BookingSlots(days []string, times map[string][]BookingHours, s BookingSettings, loc *time.Location) []Slot {
	step := orDefault(s.SlotMinutes, BookingSlotMinutes)
	breakFrom, hasBreak := clock(s.BreakFrom)
	breakTo, ok := clock(s.BreakTo)
	hasBreak = hasBreak && ok && breakFrom < breakTo
	out := []Slot{}
	for _, day := range days {
		d, err := time.ParseInLocation(time.DateOnly, day, loc)
		if err != nil {
			continue
		}
		windows := times[day]
		if hours, other := s.Days[day]; other {
			windows = []BookingHours{hours}
		} else if len(windows) == 0 {
			windows = []BookingHours{{orDefault(s.From, BookingFrom), orDefault(s.To, BookingTo)}}
		}
		for _, hours := range windows {
			from, ok1 := clock(hours.From)
			to, ok2 := clock(hours.To)
			if !ok1 || !ok2 {
				continue
			}
			for m := from; m+step <= to; m += step {
				if hasBreak && m < breakTo && m+step > breakFrom {
					m = breakTo - step
					continue
				}
				out = append(out, Slot{
					Start: time.Date(d.Year(), d.Month(), d.Day(), 0, m, 0, 0, loc),
					End:   time.Date(d.Year(), d.Month(), d.Day(), 0, m+step, 0, 0, loc),
				})
			}
		}
	}
	return out
}

// BookingDeadline returns the end of the last day clubs can ask or change their request: the deadline of the
// settings, else BookingDeadlineDays before the first day of the event.
func BookingDeadline(days []string, s BookingSettings, loc *time.Location) time.Time {
	day, err := time.ParseInLocation(time.DateOnly, s.Deadline, loc)
	if err != nil && len(days) > 0 {
		first := slices.Min(days)
		if f, err := time.ParseInLocation(time.DateOnly, first, loc); err == nil {
			day = f.AddDate(0, 0, -BookingDeadlineDays)
		}
	}
	return day.AddDate(0, 0, 1)
}

// BookingOpens returns the start of the first day clubs can ask, the zero time when booking is open right away.
func BookingOpens(s BookingSettings, loc *time.Location) time.Time {
	day, err := time.ParseInLocation(time.DateOnly, s.Opens, loc)
	if err != nil {
		return time.Time{}
	}
	return day
}

// FreeLanes returns the lanes of a range that are neither disabled nor busy, in order.
func FreeLanes(laneCount int, disabled []int32, busy []int) []int {
	out := []int{}
	for l := 1; l <= laneCount; l++ {
		if !slices.Contains(disabled, int32(l)) && !slices.Contains(busy, l) {
			out = append(out, l)
		}
	}
	return out
}

// AssignLanes picks n of the free lanes (ascending): the first run of n neighbours if there is one, else the
// lowest n; false when fewer than n are free.
func AssignLanes(free []int, n int) ([]int, bool) {
	if n <= 0 || len(free) < n {
		return nil, false
	}
	for i := 0; i+n <= len(free); i++ {
		if free[i+n-1]-free[i] == n-1 {
			return slices.Clone(free[i : i+n]), true
		}
	}
	return slices.Clone(free[:n]), true
}

// PickLanes picks n of the free lanes (ascending): the wanted ones that are free, the rest by AssignLanes from the
// other free lanes; false when fewer than n are free.
func PickLanes(free, wanted []int, n int) ([]int, bool) {
	if n <= 0 || len(free) < n {
		return nil, false
	}
	var keep, other []int
	for _, l := range free {
		if slices.Contains(wanted, l) && len(keep) < n {
			keep = append(keep, l)
		} else {
			other = append(other, l)
		}
	}
	if len(keep) < n {
		rest, _ := AssignLanes(other, n-len(keep))
		keep = append(keep, rest...)
	}
	slices.Sort(keep)
	return keep, true
}

// BookingItem asks for a number of lanes of a range in one slot.
type BookingItem struct {
	RangeID int64     `json:"rangeId"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Lanes   int       `json:"lanes" minimum:"1"`
	Wanted  []int     `json:"wanted,omitempty" required:"false" maxItems:"500" doc:"Lane numbers the club picked, as many as lanes; empty = any"`
}

// bookingKey identifies the slot of an item on its range, with its lanes when withLanes is set.
type bookingKey struct {
	RangeID    int64
	Start, End time.Time
	Lanes      int
	Wanted     string
}

// key returns the comparable key of an item; withLanes adds the lane count and the wanted lanes.
func (it BookingItem) key(withLanes bool) bookingKey {
	k := bookingKey{RangeID: it.RangeID, Start: it.Start.UTC(), End: it.End.UTC()}
	if withLanes {
		w := slices.Sorted(slices.Values(it.Wanted))
		k.Lanes, k.Wanted = it.Lanes, FormatLanes(w)
	}
	return k
}

// BookingShooter is a shooter a club announces with its request.
type BookingShooter struct {
	Firstname     string `json:"firstname" maxLength:"100"`
	Lastname      string `json:"lastname" maxLength:"100"`
	Birthday      string `json:"birthday" doc:"YYYY-MM-DD"`
	LicenseNumber string `json:"licenseNumber,omitempty" required:"false" maxLength:"20"`
}

// ValidateBooking checks a request against the slots and the lanes in service (range ID → lane numbers) of the
// event: every item a slot of a range of the event, each at most once, between 1 and the limit of lanes, wanted lanes
// (if any) as many as lanes, distinct and in service; shooters with names and a birthday, at most MaxBookingShooters.
// Lanes per item are capped by maxLanes when it is above 0.
func ValidateBooking(items []BookingItem, shooters []BookingShooter, slots []Slot, rangeLanes map[int64][]int, maxLanes int) []Problem {
	var errs []Problem
	if len(items) == 0 {
		errs = append(errs, Problemf("Choose at least one slot"))
	}
	seen := map[bookingKey]bool{}
	for _, it := range items {
		key := it.key(false)
		inService, ok := rangeLanes[it.RangeID]
		count := len(inService)
		inSlot := slices.ContainsFunc(slots, func(s Slot) bool { return s.Start.Equal(it.Start) && s.End.Equal(it.End) })
		switch {
		case !ok || !inSlot:
			errs = append(errs, Problemf("A chosen slot is not offered"))
		case seen[key]:
			errs = append(errs, Problemf("A slot is chosen twice"))
		case it.Lanes < 1 || it.Lanes > count:
			errs = append(errs, Problemf("Lanes: between 1 and %d", count))
		case maxLanes > 0 && it.Lanes > maxLanes:
			errs = append(errs, Problemf("At most %d lanes per slot", maxLanes))
		case len(it.Wanted) > 0 && !wantedOK(it, inService):
			errs = append(errs, Problemf("A chosen lane is not offered"))
		}
		seen[key] = true
	}
	if len(shooters) > MaxBookingShooters {
		errs = append(errs, Problemf("At most %d shooters", MaxBookingShooters))
	}
	for i, s := range shooters {
		if strings.TrimSpace(s.Firstname) == "" || strings.TrimSpace(s.Lastname) == "" {
			errs = append(errs, Problemf("Shooter %d: first and last name are required", i+1))
		}
		if _, err := time.Parse(time.DateOnly, s.Birthday); err != nil {
			errs = append(errs, Problemf("Shooter %d: birthday missing or not a date", i+1))
		}
	}
	return errs
}

// wantedOK reports whether the wanted lanes of an item are as many as its lanes, distinct and in service.
func wantedOK(it BookingItem, inService []int) bool {
	w := slices.Sorted(slices.Values(it.Wanted))
	if len(w) != it.Lanes || len(slices.Compact(slices.Clone(w))) != len(w) {
		return false
	}
	return !slices.ContainsFunc(w, func(l int) bool { return !slices.Contains(inService, l) })
}

// SameBookingItems reports whether two requests ask for the same lanes (count and wanted numbers) in the same slots,
// in any order.
func SameBookingItems(a, b []BookingItem) bool {
	if len(a) != len(b) {
		return false
	}
	have := map[bookingKey]int{}
	for _, it := range a {
		have[it.key(true)]++
	}
	for _, it := range b {
		k := it.key(true)
		if have[k] == 0 {
			return false
		}
		have[k]--
	}
	return true
}
