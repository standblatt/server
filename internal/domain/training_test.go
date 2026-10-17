package domain

import (
	"testing"
	"time"
)

var airPistol = Discipline{SeriesCount: 6, ShotsPerSeries: 10, MaxShot: 10}
var airRifle = Discipline{SeriesCount: 6, ShotsPerSeries: 10, MaxShot: 10, Decimal: true}

func TestValidateTraining(t *testing.T) {
	out, errs := ValidateTraining(airPistol, []TrainingSeries{
		{Shots: 10, Total: 93},
		{Shots: 3, Total: 999, Values: []float64{10, 9, 8}, X: []float64{1, 2, 3}, Y: []float64{0, 0, 0}},
	})
	if len(errs) > 0 || out[1].Total != 27 {
		t.Fatalf("got %v %v", out, errs)
	}
	cases := []struct {
		d      Discipline
		series []TrainingSeries
	}{
		{airPistol, nil},
		{airPistol, []TrainingSeries{{Shots: 10, Total: 101}}},
		{airPistol, []TrainingSeries{{Shots: 10, Total: 93.5}}},
		{airPistol, []TrainingSeries{{Shots: 2, Values: []float64{10}}}},
		{airPistol, []TrainingSeries{{Shots: 1, Values: []float64{9.5}}}},
		{airRifle, []TrainingSeries{{Shots: 1, Values: []float64{11}}}},
		{airRifle, []TrainingSeries{{Shots: 1, Values: []float64{10.45}}}},
		{airRifle, []TrainingSeries{{Shots: 1, Values: []float64{10.4}, X: []float64{1}}}},
		{airPistol, []TrainingSeries{{Shots: 0, Total: 0}}},
	}
	for i, c := range cases {
		if _, errs := ValidateTraining(c.d, c.series); len(errs) == 0 {
			t.Errorf("case %d: expected an error", i)
		}
	}
	if out, errs := ValidateTraining(airRifle, []TrainingSeries{{Shots: 2, Values: []float64{10.4, 9.8}}}); len(errs) > 0 || out[0].Total != 20.2 {
		t.Fatalf("decimal: %v %v", out, errs)
	}
}

func day(d int) time.Time { return time.Date(2026, 10, d, 18, 0, 0, 0, time.UTC) }

func session(id int64, d int, kind string, tags []string, totals ...float64) SessionFacts {
	s := SessionFacts{ID: id, StartedAt: day(d), Kind: kind, Tags: tags}
	for _, t := range totals {
		s.Series = append(s.Series, TrainingSeries{Shots: 10, Total: t})
	}
	return s
}

func TestStats(t *testing.T) {
	d := Discipline{SeriesCount: 2, ShotsPerSeries: 10, MaxShot: 10}
	sessions := []SessionFacts{
		session(3, 12, "competition", nil, 95, 91),
		session(1, 1, "training", []string{"tired"}, 90, 88),
		session(2, 5, "training", []string{"tired", "new grip"}, 92, 90),
		session(4, 13, "training", []string{"new grip"}, 94),
	}
	sessions[3].Series[0].Values = []float64{10, 10, 10, 10, 9, 9, 9, 9, 9, 9}
	sessions[3].Series[0].X = []float64{1, 3, 1, 3, 1, 3, 1, 3, 1, 3}
	sessions[3].Series[0].Y = []float64{-2, -2, 0, 0, -2, -2, 0, 0, -2, -2}
	st := Stats(d, sessions, day(14))
	if st.Sessions != 4 || st.Shots != 70 || st.PerShot != 9.14 {
		t.Fatalf("totals: %+v", st)
	}
	if st.Points[0].SessionID != 1 || st.Points[2].Moving != 9.1 {
		t.Fatalf("points: %+v", st.Points)
	}
	if st.Best == nil || st.Best.SessionID != 3 {
		t.Fatalf("best: %+v", st.Best)
	}
	if *st.Competition != 9.3 || *st.Training != 9.08 {
		t.Fatalf("kinds: %v %v", *st.Competition, *st.Training)
	}
	if len(st.Series) != 2 || st.Series[0].PerShot != 9.28 || st.Series[1].Count != 3 {
		t.Fatalf("series: %+v", st.Series)
	}
	if st.Distribution[10] != 4 || st.Distribution[9] != 6 {
		t.Fatalf("distribution: %v", st.Distribution)
	}
	if g := st.Group; g == nil || g.X != 2 || g.Y != -1.2 || g.Shots != 10 {
		t.Fatalf("group: %+v", st.Group)
	}
	if len(st.Tags) != 2 || st.Tags[0].Tag != "new grip" || st.Tags[1].Delta >= 0 {
		t.Fatalf("tags: %+v", st.Tags)
	}
	if st.Recent == nil || st.Previous != nil {
		t.Fatalf("recent: %v %v", st.Recent, st.Previous)
	}
	empty := Stats(d, nil, day(14))
	if empty.Sessions != 0 || empty.Recent != nil || empty.Best != nil || empty.Group != nil {
		t.Fatalf("empty: %+v", empty)
	}
}

func TestCompareWeapons(t *testing.T) {
	d := Discipline{SeriesCount: 2, ShotsPerSeries: 10, MaxShot: 10}
	a, b, c, none := session(1, 1, "training", nil, 90, 90), session(2, 14, "training", nil, 95, 93), session(3, 14, "training", nil, 80), session(4, 2, "training", nil, 99, 99)
	a.WeaponID, b.WeaponID, c.WeaponID = 7, 7, 8
	got := CompareWeapons(d, []SessionFacts{a, b, c, none}, day(1).AddDate(0, 0, 35))
	if len(got) != 2 || got[0].WeaponID != 7 || got[0].Sessions != 2 || got[0].PerShot != 9.2 || *got[0].Best != 188 || *got[0].Recent != 9.4 {
		t.Fatalf("first: %+v", got)
	}
	if got[1].WeaponID != 8 || got[1].PerShot != 8 || got[1].Best != nil {
		t.Fatalf("second: %+v", got[1])
	}
}

func TestFeedbackFor(t *testing.T) {
	sessions := []SessionFacts{
		session(1, 1, "training", nil, 90, 88),
		session(2, 5, "training", nil, 92, 90),
		session(3, 12, "training", nil, 95, 91),
		session(4, 13, "training", nil, 50),
	}
	f := FeedbackFor(sessions, 3)
	if f.Rank != 1 || f.Of != 3 || !f.PersonalBest || f.Delta == nil || *f.Delta != 0.3 {
		t.Fatalf("best: %+v %v", f, *f.Delta)
	}
	f = FeedbackFor(sessions, 1)
	if f.Rank != 3 || f.PersonalBest || f.Delta != nil {
		t.Fatalf("first: %+v", f)
	}
	if f := FeedbackFor(sessions, 4); f.Of != 1 || f.PersonalBest {
		t.Fatalf("alone: %+v", f)
	}
}

func TestWeekStreak(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 10, 14, 12, 0, 0, 0, loc)
	times := []time.Time{now.AddDate(0, 0, -7), now.AddDate(0, 0, -14), now.AddDate(0, 0, -28)}
	if n := WeekStreak(times, now, loc); n != 2 {
		t.Fatalf("got %d", n)
	}
	if n := WeekStreak(append(times, now), now, loc); n != 3 {
		t.Fatalf("got %d", n)
	}
	if n := WeekStreak(nil, now, loc); n != 0 {
		t.Fatalf("got %d", n)
	}
}
