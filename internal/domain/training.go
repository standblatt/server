package domain

import (
	"math"
	"slices"
	"sort"
	"time"
)

// Discipline is the scoring of a discipline: shots score 0 to MaxShot whole points, or 0.0 to 10.9 with Decimal.
type Discipline struct {
	SeriesCount    int
	ShotsPerSeries int
	MaxShot        int
	Decimal        bool
}

// ShotMax is the highest value of a single shot.
func (d Discipline) ShotMax() float64 {
	if d.Decimal {
		return 10.9
	}
	return float64(d.MaxShot)
}

// StandardShots is the number of shots of the discipline's usual match.
func (d Discipline) StandardShots() int { return d.SeriesCount * d.ShotsPerSeries }

// TrainingSeries is one series of a session: its total and optionally every shot, with or without its position.
type TrainingSeries struct {
	Shots  int       `json:"shots" minimum:"1" maximum:"100"`
	Total  float64   `json:"total" minimum:"0"`
	Values []float64 `json:"values,omitempty" doc:"Shot values in order; the total is their sum"`
	X      []float64 `json:"x,omitempty" doc:"Horizontal shot positions in mm from the centre, right is positive"`
	Y      []float64 `json:"y,omitempty" doc:"Vertical shot positions in mm from the centre, up is positive"`
}

const maxTrainingSeries = 30

// tenths rounds v to one decimal.
func tenths(v float64) float64 { return math.Round(v*10) / 10 }

// wholeOrTenths reports whether v is a whole number, or has at most one decimal when decimal is set.
func wholeOrTenths(v float64, decimal bool) bool {
	if decimal {
		return math.Abs(v*10-math.Round(v*10)) < 1e-6
	}
	return math.Abs(v-math.Round(v)) < 1e-6
}

// ValidateTraining checks series against the scoring of a discipline and returns them with totals
// recomputed from shot values.
func ValidateTraining(d Discipline, series []TrainingSeries) ([]TrainingSeries, []Problem) {
	var errs []Problem
	if len(series) == 0 || len(series) > maxTrainingSeries {
		return nil, []Problem{Problemf("A session has 1 to %d series", maxTrainingSeries)}
	}
	out := make([]TrainingSeries, len(series))
	max := d.ShotMax()
	for i, s := range series {
		n := i + 1
		if s.Shots < 1 || s.Shots > 100 {
			errs = append(errs, Problemf("Series %d: 1 to 100 shots", n))
			continue
		}
		if len(s.Values) > 0 {
			if len(s.Values) != s.Shots {
				errs = append(errs, Problemf("Series %d: %d shot values for %d shots", n, len(s.Values), s.Shots))
				continue
			}
			s.Total = 0
			for _, v := range s.Values {
				if v < 0 || v > max || !wholeOrTenths(v, d.Decimal) {
					errs = append(errs, Problemf("Series %d: shot value %g is not possible", n, v))
					break
				}
				s.Total += v
			}
			s.Total = tenths(s.Total)
		} else if s.Total < 0 || s.Total > float64(s.Shots)*max+1e-6 || !wholeOrTenths(s.Total, d.Decimal) {
			errs = append(errs, Problemf("Series %d: total %g is not possible with %d shots", n, s.Total, s.Shots))
		}
		if len(s.X) > 0 || len(s.Y) > 0 {
			if len(s.X) != len(s.Values) || len(s.Y) != len(s.Values) {
				errs = append(errs, Problemf("Series %d: every shot needs a value and a position", n))
			}
		}
		if len(s.Values) == 0 {
			s.Values = nil
		}
		if len(s.X) == 0 {
			s.X, s.Y = nil, nil
		}
		out[i] = s
	}
	return out, errs
}

// SessionFacts is what the statistics need from a session.
type SessionFacts struct {
	ID        int64
	StartedAt time.Time
	Kind      string
	Tags      []string
	WeaponID  int64
	Series    []TrainingSeries
}

// Totals returns the points and the number of shots of a session.
func (s SessionFacts) Totals() (float64, int) {
	var total float64
	var shots int
	for _, x := range s.Series {
		total += x.Total
		shots += x.Shots
	}
	return tenths(total), shots
}

// TrainingPoint is one session in the trend: total, points per shot and the moving average up to it.
type TrainingPoint struct {
	SessionID int64     `json:"sessionId"`
	At        time.Time `json:"at"`
	Kind      string    `json:"kind"`
	Total     float64   `json:"total"`
	Shots     int       `json:"shots"`
	PerShot   float64   `json:"perShot"`
	Moving    float64   `json:"moving" doc:"Average points per shot of this and the four sessions before"`
}

// SeriesStat is the average points per shot of the n-th series over all sessions.
type SeriesStat struct {
	Number  int     `json:"number"`
	PerShot float64 `json:"perShot"`
	Count   int     `json:"count"`
}

// TagStat compares the sessions with a tag to the overall average.
type TagStat struct {
	Tag      string  `json:"tag"`
	Sessions int     `json:"sessions"`
	PerShot  float64 `json:"perShot"`
	Delta    float64 `json:"delta"`
}

// WeaponStat sums up the sessions shot with one weapon.
type WeaponStat struct {
	WeaponID int64    `json:"weaponId"`
	Sessions int      `json:"sessions"`
	PerShot  float64  `json:"perShot"`
	Recent   *float64 `json:"recent" doc:"Points per shot in the last 30 days"`
	Best     *float64 `json:"best" doc:"Highest total among sessions with the standard number of shots"`
}

// CompareWeapons returns one line per weapon used in the sessions, highest points per shot first.
// Sessions without a weapon are left out.
func CompareWeapons(d Discipline, sessions []SessionFacts, now time.Time) []WeaponStat {
	type sum struct {
		all, recent acc
		sessions    int
		best        *float64
	}
	sums := map[int64]*sum{}
	recentFrom := now.AddDate(0, 0, -recentDays)
	for _, s := range sessions {
		total, shots := s.Totals()
		if s.WeaponID == 0 || shots == 0 {
			continue
		}
		w := sums[s.WeaponID]
		if w == nil {
			w = &sum{}
			sums[s.WeaponID] = w
		}
		w.sessions++
		w.all.add(total, shots)
		if !s.StartedAt.Before(recentFrom) {
			w.recent.add(total, shots)
		}
		if shots == d.StandardShots() && (w.best == nil || total > *w.best) {
			w.best = &total
		}
	}
	out := []WeaponStat{}
	for id, w := range sums {
		out = append(out, WeaponStat{WeaponID: id, Sessions: w.sessions, PerShot: *w.all.perShot(), Recent: w.recent.perShot(), Best: w.best})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].PerShot != out[j].PerShot {
			return out[i].PerShot > out[j].PerShot
		}
		return out[i].WeaponID < out[j].WeaponID
	})
	return out
}

// Group is the mean point of impact of positioned shots and their mean distance from it, in mm.
type Group struct {
	Shots  int     `json:"shots"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

// TrainingStats sums up the sessions of one discipline.
type TrainingStats struct {
	Sessions     int             `json:"sessions"`
	Shots        int             `json:"shots"`
	PerShot      float64         `json:"perShot"`
	Recent       *float64        `json:"recent" doc:"Points per shot in the last 30 days"`
	Previous     *float64        `json:"previous" doc:"Points per shot 30 to 60 days ago"`
	Training     *float64        `json:"training"`
	Competition  *float64        `json:"competition"`
	Best         *TrainingPoint  `json:"best" doc:"Highest total among sessions with the standard number of shots"`
	Points       []TrainingPoint `json:"points"`
	Series       []SeriesStat    `json:"series"`
	Distribution []int           `json:"distribution" doc:"Number of shots per whole value, index 0 = value 0"`
	Group        *Group          `json:"group" doc:"Positioned shots of the last 10 sessions that have positions"`
	Tags         []TagStat       `json:"tags"`
	Weapons      []WeaponStat    `json:"weapons" doc:"One line per weapon used in the discipline"`
}

type acc struct {
	points float64
	shots  int
}

// add adds points and shots to the accumulator.
func (a *acc) add(p float64, s int) { a.points += p; a.shots += s }

// perShot returns the points per shot rounded to two decimals, or nil without shots.
func (a acc) perShot() *float64 {
	if a.shots == 0 {
		return nil
	}
	v := Round2(a.points / float64(a.shots))
	return &v
}

// Round2 rounds v to two decimals.
func Round2(v float64) float64 { return math.Round(v*100) / 100 }

const (
	movingWindow  = 5
	groupSessions = 10
	recentDays    = 30
)

// Stats computes trend, averages, series curve, value distribution, group and tags for the sessions of one discipline.
func Stats(d Discipline, sessions []SessionFacts, now time.Time) TrainingStats {
	sessions = slices.Clone(sessions)
	sort.SliceStable(sessions, func(i, j int) bool { return sessions[i].StartedAt.Before(sessions[j].StartedAt) })
	st := TrainingStats{Points: []TrainingPoint{}, Series: []SeriesStat{}, Tags: []TagStat{}, Weapons: []WeaponStat{}, Distribution: make([]int, int(d.ShotMax())+1)}
	var all, recent, previous, training, competition acc
	var window []acc
	var seriesAcc []acc
	tags := map[string]*acc{}
	tagCount := map[string]int{}
	recentFrom, previousFrom := now.AddDate(0, 0, -recentDays), now.AddDate(0, 0, -2*recentDays)
	for _, s := range sessions {
		total, shots := s.Totals()
		if shots == 0 {
			continue
		}
		all.add(total, shots)
		switch {
		case !s.StartedAt.Before(recentFrom):
			recent.add(total, shots)
		case !s.StartedAt.Before(previousFrom):
			previous.add(total, shots)
		}
		if s.Kind == "competition" {
			competition.add(total, shots)
		} else {
			training.add(total, shots)
		}
		window = append(window, acc{total, shots})
		if len(window) > movingWindow {
			window = window[1:]
		}
		var w acc
		for _, a := range window {
			w.add(a.points, a.shots)
		}
		p := TrainingPoint{SessionID: s.ID, At: s.StartedAt, Kind: s.Kind, Total: total, Shots: shots,
			PerShot: Round2(total / float64(shots)), Moving: *w.perShot()}
		st.Points = append(st.Points, p)
		if shots == d.StandardShots() && (st.Best == nil || total > st.Best.Total) {
			best := p
			st.Best = &best
		}
		for i, x := range s.Series {
			if i == len(seriesAcc) {
				seriesAcc = append(seriesAcc, acc{})
			}
			seriesAcc[i].add(x.Total, x.Shots)
			for _, v := range x.Values {
				if b := int(math.Floor(v + 1e-9)); b >= 0 && b < len(st.Distribution) {
					st.Distribution[b]++
				}
			}
		}
		for _, t := range s.Tags {
			if tags[t] == nil {
				tags[t] = &acc{}
			}
			tags[t].add(total, shots)
			tagCount[t]++
		}
	}
	st.Sessions, st.Shots = len(st.Points), all.shots
	if p := all.perShot(); p != nil {
		st.PerShot = *p
	}
	st.Recent, st.Previous = recent.perShot(), previous.perShot()
	st.Training, st.Competition = training.perShot(), competition.perShot()
	for i, a := range seriesAcc {
		st.Series = append(st.Series, SeriesStat{Number: i + 1, PerShot: *a.perShot(), Count: countSeries(sessions, i)})
	}
	for t, a := range tags {
		if tagCount[t] < 2 {
			continue
		}
		p := *a.perShot()
		st.Tags = append(st.Tags, TagStat{Tag: t, Sessions: tagCount[t], PerShot: p, Delta: Round2(p - st.PerShot)})
	}
	sort.Slice(st.Tags, func(i, j int) bool { return st.Tags[i].Delta > st.Tags[j].Delta })
	st.Group = group(sessions)
	return st
}

// countSeries counts the sessions that have a series at index i.
func countSeries(sessions []SessionFacts, i int) int {
	n := 0
	for _, s := range sessions {
		if len(s.Series) > i {
			n++
		}
	}
	return n
}

// group returns the shot group over the last groupSessions sessions that have shot coordinates.
func group(sessions []SessionFacts) *Group {
	var xs, ys []float64
	used := 0
	for i := len(sessions) - 1; i >= 0 && used < groupSessions; i-- {
		found := false
		for _, s := range sessions[i].Series {
			xs = append(xs, s.X...)
			ys = append(ys, s.Y...)
			found = found || len(s.X) > 0
		}
		if found {
			used++
		}
	}
	return GroupOf(xs, ys)
}

// GroupOf returns the mean point of impact and the mean radius around it, or nil without shots.
func GroupOf(xs, ys []float64) *Group {
	if len(xs) == 0 || len(xs) != len(ys) {
		return nil
	}
	var g Group
	for i := range xs {
		g.X += xs[i]
		g.Y += ys[i]
	}
	n := float64(len(xs))
	g.X, g.Y = g.X/n, g.Y/n
	for i := range xs {
		g.Radius += math.Hypot(xs[i]-g.X, ys[i]-g.Y)
	}
	g.Shots, g.Radius = len(xs), Round2(g.Radius/n)
	g.X, g.Y = Round2(g.X), Round2(g.Y)
	return &g
}

// Feedback is what a shooter learns right after saving a session.
type Feedback struct {
	Rank         int      `json:"rank" doc:"Place of the session among sessions with the same number of shots; 0 if it has no shots"`
	Of           int      `json:"of"`
	PersonalBest bool     `json:"personalBest" doc:"Best of at least two comparable sessions"`
	Delta        *float64 `json:"delta" doc:"Points per shot compared with the 30 days before the session"`
}

// FeedbackFor rates the session with the given ID against the other sessions of its discipline.
func FeedbackFor(sessions []SessionFacts, id int64) Feedback {
	var cur *SessionFacts
	for i := range sessions {
		if sessions[i].ID == id {
			cur = &sessions[i]
		}
	}
	var f Feedback
	if cur == nil {
		return f
	}
	total, shots := cur.Totals()
	if shots == 0 {
		return f
	}
	f.Rank = 1
	var before acc
	from := cur.StartedAt.AddDate(0, 0, -recentDays)
	for _, s := range sessions {
		if s.ID == id {
			continue
		}
		t, n := s.Totals()
		if n == shots {
			f.Of++
			if t > total+1e-9 {
				f.Rank++
			}
		}
		if n > 0 && s.StartedAt.Before(cur.StartedAt) && !s.StartedAt.Before(from) {
			before.add(t, n)
		}
	}
	f.Of++
	f.PersonalBest = f.Rank == 1 && f.Of > 1
	if p := before.perShot(); p != nil {
		d := Round2(total/float64(shots) - *p)
		f.Delta = &d
	}
	return f
}

// WeekStreak counts the consecutive weeks (Monday to Sunday in loc) with at least one session, ending with
// the current week, or the week before if the current one has none yet.
func WeekStreak(times []time.Time, now time.Time, loc *time.Location) int {
	weeks := map[time.Time]bool{}
	for _, t := range times {
		weeks[weekStart(t, loc)] = true
	}
	w := weekStart(now, loc)
	if !weeks[w] {
		w = w.AddDate(0, 0, -7)
	}
	n := 0
	for weeks[w] {
		n++
		w = w.AddDate(0, 0, -7)
	}
	return n
}

// weekStart returns midnight of the Monday of t's week in loc.
func weekStart(t time.Time, loc *time.Location) time.Time {
	t = t.In(loc)
	day := (int(t.Weekday()) + 6) % 7
	return time.Date(t.Year(), t.Month(), t.Day()-day, 0, 0, 0, 0, loc)
}
