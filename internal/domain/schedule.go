package domain

import (
	"cmp"
	"slices"
	"time"
)

// maxScheduleDays caps the expansion of a period so a typo in the year cannot produce thousands of days.
const maxScheduleDays = 400

// ScheduleDay is a single shooting day with optional shooting times.
type ScheduleDay struct {
	Date  string         `json:"date" format:"date"`
	Times []BookingHours `json:"times,omitempty" required:"false" doc:"Shooting times of the day; none = not announced"`
}

// WeeklyTimes are the shooting times of one weekday within the period.
type WeeklyTimes struct {
	Weekday int            `json:"weekday" minimum:"1" maximum:"7" doc:"1 = Monday … 7 = Sunday"`
	Times   []BookingHours `json:"times,omitempty" required:"false"`
}

// Schedule says when an event is shot: single days, a period (season) with weekly shooting times, or both.
type Schedule struct {
	Days   []ScheduleDay `json:"days,omitempty" required:"false" doc:"Single days; they override the weekly times on the same date"`
	From   string        `json:"from,omitempty" required:"false" doc:"First day of the period (YYYY-MM-DD)"`
	To     string        `json:"to,omitempty" required:"false" doc:"Last day of the period (YYYY-MM-DD)"`
	Weekly []WeeklyTimes `json:"weekly,omitempty" required:"false" doc:"Shooting weekdays of the period; none = every day"`
}

// isoWeekday returns 1 for Monday … 7 for Sunday.
func isoWeekday(t time.Time) int {
	return (int(t.Weekday())+6)%7 + 1
}

// ValidateSchedule checks dates, times and the period; the main day is added as a single day by NormalizeSchedule.
func ValidateSchedule(s Schedule) []Problem {
	var errs []Problem
	valid := func(ts []BookingHours) bool {
		for _, h := range ts {
			from, ok1 := clock(h.From)
			to, ok2 := clock(h.To)
			if !ok1 || !ok2 || from >= to {
				return false
			}
		}
		return true
	}
	for _, d := range s.Days {
		if _, err := time.Parse(time.DateOnly, d.Date); err != nil {
			errs = append(errs, Problemf("Date must be YYYY-MM-DD"))
			continue
		}
		if !valid(d.Times) {
			errs = append(errs, Problemf("Shooting times of %s: the end must be after the start", d.Date))
		}
	}
	if s.From != "" || s.To != "" {
		from, err1 := time.Parse(time.DateOnly, s.From)
		to, err2 := time.Parse(time.DateOnly, s.To)
		switch {
		case err1 != nil || err2 != nil:
			errs = append(errs, Problemf("The period needs a first and a last day"))
		case to.Before(from):
			errs = append(errs, Problemf("The period ends before it starts"))
		case to.Sub(from) > maxScheduleDays*24*time.Hour:
			errs = append(errs, Problemf("The period is longer than %d days", maxScheduleDays))
		}
	} else if len(s.Weekly) > 0 {
		errs = append(errs, Problemf("Weekly shooting times need a period"))
	}
	for _, w := range s.Weekly {
		if w.Weekday < 1 || w.Weekday > 7 {
			errs = append(errs, Problemf("Weekday must be 1–7"))
			continue
		}
		if !valid(w.Times) {
			errs = append(errs, Problemf("Weekly shooting times: the end must be after the start"))
		}
	}
	return errs
}

// NormalizeSchedule sorts and dedupes the single days, adds the main day unless the period covers it and drops
// duplicate weekdays (the first wins).
func NormalizeSchedule(s Schedule, mainDay string) Schedule {
	if !slices.ContainsFunc(s.Expand(), func(d ScheduleDay) bool { return d.Date == mainDay }) {
		s.Days = append(s.Days, ScheduleDay{Date: mainDay})
	}
	slices.SortStableFunc(s.Days, func(a, b ScheduleDay) int { return cmp.Compare(a.Date, b.Date) })
	s.Days = slices.CompactFunc(s.Days, func(a, b ScheduleDay) bool { return a.Date == b.Date })
	seen := map[int]bool{}
	s.Weekly = slices.DeleteFunc(s.Weekly, func(w WeeklyTimes) bool {
		dup := seen[w.Weekday]
		seen[w.Weekday] = true
		return dup
	})
	slices.SortFunc(s.Weekly, func(a, b WeeklyTimes) int { return a.Weekday - b.Weekday })
	return s
}

// Expand returns every shooting day, sorted: the days of the period on the weekly weekdays (every day without a
// weekly rule) with their times, plus the single days, which replace a period day of the same date.
func (s Schedule) Expand() []ScheduleDay {
	byDate := map[string]ScheduleDay{}
	from, err1 := time.Parse(time.DateOnly, s.From)
	to, err2 := time.Parse(time.DateOnly, s.To)
	if err1 == nil && err2 == nil && !to.Before(from) && to.Sub(from) <= maxScheduleDays*24*time.Hour {
		for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
			date := d.Format(time.DateOnly)
			if len(s.Weekly) == 0 {
				byDate[date] = ScheduleDay{Date: date}
				continue
			}
			for _, w := range s.Weekly {
				if w.Weekday == isoWeekday(d) {
					byDate[date] = ScheduleDay{Date: date, Times: w.Times}
				}
			}
		}
	}
	for _, d := range s.Days {
		byDate[d.Date] = d
	}
	out := make([]ScheduleDay, 0, len(byDate))
	for _, d := range byDate {
		out = append(out, d)
	}
	slices.SortFunc(out, func(a, b ScheduleDay) int { return cmp.Compare(a.Date, b.Date) })
	return out
}

// Dates returns the dates of Expand.
func (s Schedule) Dates() []string {
	out := []string{}
	for _, d := range s.Expand() {
		out = append(out, d.Date)
	}
	return out
}

// Times returns the shooting times of the expanded days that have some, by date.
func (s Schedule) Times() map[string][]BookingHours {
	out := map[string][]BookingHours{}
	for _, d := range s.Expand() {
		if len(d.Times) > 0 {
			out[d.Date] = d.Times
		}
	}
	return out
}

// Shift moves every date of the schedule by n days; weekly times stay on their weekdays.
func (s Schedule) Shift(n int) Schedule {
	move := func(date string) string {
		d, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return date
		}
		return d.AddDate(0, 0, n).Format(time.DateOnly)
	}
	days := make([]ScheduleDay, len(s.Days))
	for i, d := range s.Days {
		days[i] = ScheduleDay{Date: move(d.Date), Times: d.Times}
	}
	s.Days = days
	if s.From != "" {
		s.From, s.To = move(s.From), move(s.To)
	}
	return s
}
