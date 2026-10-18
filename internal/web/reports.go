package web

import (
	"context"
	"maps"
	"math"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/docs"
	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
)

// Statistics summarises an event.
type Statistics struct {
	Registered     int             `json:"registered"`
	Shooting       int             `json:"shooting"`
	Completed      int             `json:"completed"`
	CompletionRate float64         `json:"completionRate" doc:"Completed registrations in percent"`
	Entries        int             `json:"entries" doc:"Programmes sold"`
	Revenue        float64         `json:"revenue"`
	Programmes     []ProgrammeStat `json:"programmes"`
	Categories     []CategoryStat  `json:"categories"`
	Days           []DayStat       `json:"days" doc:"Registrations per shooting day, from the time of registration"`
	Awards         []AwardStat     `json:"awards"`
	AwardsValue    float64         `json:"awardsValue"`
}

// DayStat is the participation and revenue of one day.
type DayStat struct {
	Date          string  `json:"date" format:"date"`
	Registrations int     `json:"registrations"`
	Revenue       float64 `json:"revenue"`
}

// ProgrammeStat is the participation and revenue of one programme.
type ProgrammeStat struct {
	ProgrammeID int64   `json:"programmeId"`
	Title       string  `json:"title"`
	Entries     int     `json:"entries"`
	Completed   int     `json:"completed"`
	Revenue     float64 `json:"revenue"`
}

// CategoryStat counts registrations of an age category.
type CategoryStat struct {
	Category domain.AgeCategory `json:"category" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
	Count    int                `json:"count"`
}

// AwardStat counts how often an award tier was given.
type AwardStat struct {
	Programme string  `json:"programme"`
	Name      string  `json:"name"`
	Kind      string  `json:"kind" enum:"natural,card"`
	Value     float64 `json:"value"`
	Count     int     `json:"count"`
	Total     float64 `json:"total"`
}

// ClubCompetitionLine is a club in the club competition view.
type ClubCompetitionLine struct {
	domain.ClubResult
	Name      string `json:"name"`
	SSVNumber string `json:"ssvNumber"`
}

// ClubCompetitionView is the club competition of an event.
type ClubCompetitionView struct {
	Programmes []string              `json:"programmes"`
	Ranked     []ClubCompetitionLine `json:"ranked"`
	Unranked   []ClubCompetitionLine `json:"unranked"`
}

type eventData struct {
	event      store.Event
	programmes map[int64]store.Programme
	order      []store.Programme
	facts      []store.EntryFact
	counting   map[int64]bool
	rankings   []ProgrammeRanking
	ranked     []ClassificationRanking
	bookings   []store.BookingRequest
}

// loadEventData loads the event, its programmes and its entry facts for statistics and reports.
func (s *Server) loadEventData(ctx context.Context, tenantID, eventID int64) (eventData, error) {
	var d eventData
	var err error
	if d.event, err = s.Store.GetEvent(ctx, tenantID, eventID); err != nil {
		return d, notFound(err, "Event not found")
	}
	if d.order, err = s.Store.ListProgrammes(ctx, tenantID, eventID); err != nil {
		return d, err
	}
	d.programmes = map[int64]store.Programme{}
	for _, p := range d.order {
		d.programmes[p.ID] = p
	}
	if d.facts, err = s.Store.ListEntryFacts(ctx, tenantID, eventID); err != nil {
		return d, err
	}
	d.counting, err = s.Store.CountingEntries(ctx, tenantID, eventID)
	return d, err
}

// title returns the title of a programme of the event.
func (d eventData) title(id int64) string {
	p := d.programmes[id]
	return programmeTitle(p.Name, p.Variant)
}

// statistics derives the event statistics from the entry facts: revenue, entries, days, categories and awards.
func (d eventData) statistics() Statistics {
	st := Statistics{Programmes: []ProgrammeStat{}, Categories: []CategoryStat{}, Days: []DayStat{}, Awards: []AwardStat{}}
	loc := eventLoc(d.event.Timezone)
	days := map[string]*DayStat{}
	seenDay := map[int64]bool{}
	regs := map[int64]store.EntryFact{}
	perProgramme := map[int64]*ProgrammeStat{}
	for _, p := range d.order {
		perProgramme[p.ID] = &ProgrammeStat{ProgrammeID: p.ID, Title: programmeTitle(p.Name, p.Variant)}
	}
	type awardKey struct {
		programme int64
		tier      int64
	}
	awards := map[awardKey]*AwardStat{}
	var awardOrder []awardKey
	for _, f := range d.facts {
		regs[f.RegistrationID] = f
		st.Entries++
		st.Revenue += f.Price
		day := f.RegisteredAt.In(loc).Format(time.DateOnly)
		if days[day] == nil {
			days[day] = &DayStat{Date: day}
		}
		days[day].Revenue += f.Price
		if !seenDay[f.RegistrationID] {
			seenDay[f.RegistrationID] = true
			days[day].Registrations++
		}
		ps := perProgramme[f.ProgrammeID]
		ps.Entries++
		ps.Revenue += f.Price
		if f.EntryStatus == "completed" {
			ps.Completed++
		}
		if f.AwardTierID != nil {
			k := awardKey{f.ProgrammeID, *f.AwardTierID}
			if awards[k] == nil {
				awards[k] = &AwardStat{Programme: d.title(f.ProgrammeID), Name: f.AwardName, Kind: f.AwardKind, Value: f.AwardValue}
				awardOrder = append(awardOrder, k)
			}
			awards[k].Count++
			awards[k].Total += f.AwardValue
			st.AwardsValue += f.AwardValue
		}
	}
	categories := map[domain.AgeCategory]int{}
	for _, f := range regs {
		switch f.RegistrationStatus {
		case store.StatusRegistered:
			st.Registered++
		case store.StatusShooting:
			st.Shooting++
		case store.StatusCompleted:
			st.Completed++
		}
		categories[f.Category]++
	}
	if n := len(regs); n > 0 {
		st.CompletionRate = math.Round(float64(st.Completed)/float64(n)*1000) / 10
	}
	for _, c := range domain.AgeCategories {
		if categories[c] > 0 {
			st.Categories = append(st.Categories, CategoryStat{Category: c, Count: categories[c]})
		}
	}
	for _, p := range d.order {
		ps := perProgramme[p.ID]
		ps.Revenue = domain.Round2(ps.Revenue)
		st.Programmes = append(st.Programmes, *ps)
	}
	for _, day := range slices.Sorted(maps.Keys(days)) {
		ds := days[day]
		ds.Revenue = domain.Round2(ds.Revenue)
		st.Days = append(st.Days, *ds)
	}
	for _, k := range awardOrder {
		a := awards[k]
		a.Total = domain.Round2(a.Total)
		st.Awards = append(st.Awards, *a)
	}
	st.Revenue, st.AwardsValue = domain.Round2(st.Revenue), domain.Round2(st.AwardsValue)
	return st
}

// clubCompetition ranks the clubs by the first two series of the counting attempts in programmes that count for it.
func (d eventData) clubCompetition() ClubCompetitionView {
	v := ClubCompetitionView{Programmes: []string{}, Ranked: []ClubCompetitionLine{}, Unranked: []ClubCompetitionLine{}}
	for _, p := range d.order {
		if p.CountsForClubCompetition {
			v.Programmes = append(v.Programmes, programmeTitle(p.Name, p.Variant))
		}
	}
	names := map[int64][2]string{}
	var scores []domain.ClubScore
	for _, f := range d.facts {
		if f.ClubID == nil || !d.counting[f.EntryID] || !d.programmes[f.ProgrammeID].CountsForClubCompetition {
			continue
		}
		names[*f.ClubID] = [2]string{f.ClubName, f.ClubSSVNumber}
		scores = append(scores, domain.ClubScore{ClubID: *f.ClubID, Score: f.FirstTwoSeries})
	}
	excluded := int64(0)
	if d.event.Settings.ExcludedClubID != nil {
		excluded = *d.event.Settings.ExcludedClubID
	}
	ranked, unranked := domain.ClubCompetition(scores, excluded)
	for _, r := range ranked {
		v.Ranked = append(v.Ranked, ClubCompetitionLine{ClubResult: r, Name: names[r.ClubID][0], SSVNumber: names[r.ClubID][1]})
	}
	for _, r := range unranked {
		v.Unranked = append(v.Unranked, ClubCompetitionLine{ClubResult: r, Name: names[r.ClubID][0], SSVNumber: names[r.ClubID][1]})
	}
	return v
}

// chf formats money the Swiss way: CHF 1'234.50.
func chf(v float64) string {
	s := strconv.FormatFloat(math.Abs(v), 'f', 2, 64)
	whole, frac, _ := strings.Cut(s, ".")
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte('\'')
		}
		b.WriteRune(r)
	}
	sign := ""
	if v < 0 {
		sign = "-"
	}
	return "CHF " + sign + b.String() + "." + frac
}

// kindLabel returns the translated label of an award kind.
func kindLabel(l i18n.Lang, kind string) string {
	switch kind {
	case string(domain.AwardCard):
		return i18n.T(l, "Wreath card")
	case string(domain.AwardCash):
		return i18n.T(l, "Cash prize")
	}
	return i18n.T(l, "Prize in kind")
}

// Reports that can be exported.
const (
	ReportFinancial        = "financial"
	ReportRevenue          = "revenue"
	ReportAwards           = "awards"
	ReportClubCompetition  = "club_competition"
	ReportClubParticipants = "club_participants"
	ReportSpecialPrizes    = "special_prizes"
	ReportAwardsOpen       = "awards_open"
	ReportBookings         = "bookings"
	ReportClassifications  = "classifications"
)

// report builds the report of kind and its file name, 404 for unknown kinds.
func (d eventData) report(kind string, now time.Time) (docs.Report, string, error) {
	l := i18n.Get(d.event.Locale)
	r := docs.Report{Lang: l, Subtitle: d.event.Name + " · " + swissDate(d.event.StartsOn), Generated: now}
	switch kind {
	case ReportFinancial:
		st := d.statistics()
		r.Title = i18n.T(l, "Financial report")
		r.Summary = [][2]string{
			{i18n.T(l, "Revenue"), chf(st.Revenue)}, {i18n.T(l, "Programmes sold"), strconv.Itoa(st.Entries)},
			{i18n.T(l, "Awards given"), strconv.Itoa(len(awardFacts(d.facts)))}, {i18n.T(l, "Value of awards"), chf(st.AwardsValue)},
		}
		rev := docs.Section{Heading: i18n.T(l, "Revenue by programme"), Columns: []docs.Column{col(i18n.T(l, "Programme"), 5, "L"), col(i18n.T(l, "Participants"), 1.5, "C"), col(i18n.T(l, "Revenue"), 2, "R")}}
		for _, p := range st.Programmes {
			rev.Rows = append(rev.Rows, []string{p.Title, strconv.Itoa(p.Entries), chf(p.Revenue)})
		}
		rev.Total = []string{i18n.T(l, "Total"), strconv.Itoa(st.Entries), chf(st.Revenue)}
		aw := docs.Section{Heading: i18n.T(l, "Awards"), Columns: []docs.Column{col(i18n.T(l, "Programme"), 3, "L"), col(i18n.T(l, "Award"), 3, "L"), col(i18n.T(l, "Kind"), 2, "L"), col(i18n.T(l, "Value"), 1.5, "R"), col(i18n.T(l, "Count"), 1, "C"), col(i18n.T(l, "Total"), 1.5, "R")}}
		count := 0
		for _, a := range st.Awards {
			aw.Rows = append(aw.Rows, []string{a.Programme, a.Name, kindLabel(l, a.Kind), chf(a.Value), strconv.Itoa(a.Count), chf(a.Total)})
			count += a.Count
		}
		aw.Total = []string{i18n.T(l, "Total"), "", "", "", strconv.Itoa(count), chf(st.AwardsValue)}
		r.Sections = []docs.Section{rev, aw}
	case ReportRevenue:
		r.Title = i18n.T(l, "Revenue by programme and age category")
		type key struct {
			programme int64
			category  domain.AgeCategory
		}
		agg := map[key][2]float64{}
		for _, f := range d.facts {
			k := key{f.ProgrammeID, f.Category}
			v := agg[k]
			agg[k] = [2]float64{v[0] + 1, v[1] + f.Price}
		}
		sec := docs.Section{Columns: []docs.Column{col(i18n.T(l, "Programme"), 4, "L"), col(i18n.T(l, "Category"), 2, "L"), col(i18n.T(l, "Participants"), 1.5, "C"), col(i18n.T(l, "Price"), 1.5, "R"), col(i18n.T(l, "Subtotal"), 2, "R")}}
		total, n := 0.0, 0
		for _, p := range d.order {
			for _, c := range domain.AgeCategories {
				v, ok := agg[key{p.ID, c}]
				if !ok {
					continue
				}
				sec.Rows = append(sec.Rows, []string{d.title(p.ID), categoryLabel(l, c), strconv.Itoa(int(v[0])), chf(v[1] / v[0]), chf(v[1])})
				total += v[1]
				n += int(v[0])
			}
		}
		sec.Total = []string{i18n.T(l, "Total"), "", strconv.Itoa(n), "", chf(total)}
		r.Sections = []docs.Section{sec}
	case ReportAwards:
		r.Title, r.Landscape = i18n.T(l, "Awards"), true
		sec := docs.Section{Columns: []docs.Column{
			col(i18n.T(l, "Last name"), 2, "L"), col(i18n.T(l, "First name"), 2, "L"), col(i18n.T(l, "Category"), 1.5, "L"), col(i18n.T(l, "Club"), 3, "L"), col(i18n.T(l, "Programme"), 3, "L"),
			col(i18n.T(l, "Kind"), 1.5, "L"), col(i18n.T(l, "Award"), 2, "L"), col(i18n.T(l, "Value"), 1.3, "R"), col(i18n.T(l, "Result"), 1, "R"), col(i18n.T(l, "Mouches"), 1, "R"),
		}}
		facts := awardFacts(d.facts)
		sort.SliceStable(facts, func(i, j int) bool {
			if facts[i].AwardKind != facts[j].AwardKind {
				return facts[i].AwardKind < facts[j].AwardKind
			}
			return facts[i].AwardName < facts[j].AwardName
		})
		for _, f := range facts {
			sec.Rows = append(sec.Rows, []string{f.Lastname, f.Firstname, categoryLabel(l, f.Category), f.ClubName, d.title(f.ProgrammeID),
				kindLabel(l, f.AwardKind), f.AwardName, chf(f.AwardValue), docs.Points(f.Total), strconv.Itoa(f.Mouches)})
		}
		r.Sections = []docs.Section{sec}
	case ReportClubCompetition:
		v := d.clubCompetition()
		r.Title = i18n.T(l, "Club competition")
		if len(v.Programmes) > 0 {
			r.Summary = [][2]string{{i18n.T(l, "Programmes counted"), strings.Join(v.Programmes, ", ")}}
		}
		ranked := docs.Section{Heading: i18n.T(l, "Ranking"), Columns: []docs.Column{col(i18n.T(l, "Rank"), 1, "C"), col(i18n.T(l, "SSV no."), 1.5, "C"), col(i18n.T(l, "Club"), 5, "L"), col(i18n.T(l, "Result"), 1.5, "R"), col(i18n.T(l, "Participants"), 1.5, "C")}}
		outside := docs.Section{Heading: i18n.T(l, "Out of competition"), Columns: ranked.Columns}
		for _, c := range v.Ranked {
			row := []string{strconv.Itoa(c.Rank), c.SSVNumber, c.Name, strconv.FormatFloat(c.Result, 'f', 3, 64), strconv.Itoa(c.Participants)}
			if c.Excluded {
				row[0] = "–"
				outside.Rows = append(outside.Rows, row)
			} else {
				ranked.Rows = append(ranked.Rows, row)
			}
		}
		unranked := docs.Section{
			Heading: i18n.T(l, "Not ranked"), Note: i18n.Tf(l, "At least %d participants are required for a ranking.", domain.ClubMinParticipants),
			Columns: []docs.Column{col(i18n.T(l, "SSV no."), 1.5, "C"), col(i18n.T(l, "Club"), 6.5, "L"), col(i18n.T(l, "Participants"), 1.5, "C")},
		}
		for _, c := range v.Unranked {
			unranked.Rows = append(unranked.Rows, []string{c.SSVNumber, c.Name, strconv.Itoa(c.Participants)})
		}
		r.Sections = []docs.Section{ranked}
		if len(outside.Rows) > 0 {
			r.Sections = append(r.Sections, outside)
		}
		r.Sections = append(r.Sections, unranked)
	case ReportClubParticipants:
		r.Title = i18n.T(l, "Club participants by programme and age category")
		type key struct {
			programme int64
			club      string
			ssv       string
			category  domain.AgeCategory
		}
		counts := map[key]int{}
		var keys []key
		for _, f := range d.facts {
			k := key{f.ProgrammeID, f.ClubName, f.ClubSSVNumber, f.Category}
			if counts[k] == 0 {
				keys = append(keys, k)
			}
			counts[k]++
		}
		sort.SliceStable(keys, func(i, j int) bool {
			if keys[i].club != keys[j].club {
				return keys[i].club < keys[j].club
			}
			return keys[i].category < keys[j].category
		})
		for _, p := range d.order {
			sec := docs.Section{Heading: d.title(p.ID), Columns: []docs.Column{col(i18n.T(l, "Club"), 5, "L"), col(i18n.T(l, "SSV no."), 1.5, "C"), col(i18n.T(l, "Category"), 2, "L"), col(i18n.T(l, "Count"), 1, "C")}}
			n := 0
			for _, k := range keys {
				if k.programme != p.ID {
					continue
				}
				club := k.club
				if club == "" {
					club = i18n.T(l, "(no club)")
				}
				sec.Rows = append(sec.Rows, []string{club, k.ssv, categoryLabel(l, k.category), strconv.Itoa(counts[k])})
				n += counts[k]
			}
			if n == 0 {
				continue
			}
			sec.Total = []string{i18n.T(l, "Total"), "", "", strconv.Itoa(n)}
			r.Sections = append(r.Sections, sec)
		}
	case ReportAwardsOpen:
		r.Title = i18n.T(l, "Awards to hand out")
		sec := docs.Section{Columns: []docs.Column{
			col(i18n.T(l, "Last name"), 2, "L"), col(i18n.T(l, "First name"), 2, "L"), col(i18n.T(l, "Club"), 3, "L"), col(i18n.T(l, "Programme"), 3, "L"),
			col(i18n.T(l, "Kind"), 1.5, "L"), col(i18n.T(l, "Award"), 2, "L"), col(i18n.T(l, "Value"), 1.3, "R"),
		}}
		shooters := map[int64]bool{}
		cash, value := 0.0, 0.0
		for _, f := range awardFacts(d.facts) {
			if f.AwardsIssuedAt != nil {
				continue
			}
			shooters[f.RegistrationID] = true
			value += f.AwardValue
			if f.AwardKind == string(domain.AwardCash) {
				cash += f.AwardValue
			}
			sec.Rows = append(sec.Rows, []string{f.Lastname, f.Firstname, f.ClubName, d.title(f.ProgrammeID), kindLabel(l, f.AwardKind), f.AwardName, chf(f.AwardValue)})
		}
		sec.Total = []string{i18n.T(l, "Total"), "", "", "", "", strconv.Itoa(len(sec.Rows)), chf(value)}
		r.Summary = [][2]string{{i18n.T(l, "Shooters with awards to collect"), strconv.Itoa(len(shooters))}, {i18n.T(l, "Cash to pay out"), chf(cash)}}
		r.Sections = []docs.Section{sec}
	case ReportSpecialPrizes:
		r.Title = i18n.T(l, "Special prizes")
		for _, pr := range d.rankings {
			rules := d.programmes[pr.ProgrammeID].SpecialPrizes
			if rules.Ranks == 0 {
				continue
			}
			candidates := make([]domain.PrizeCandidate, len(pr.Lines))
			clubs := map[string]int{}
			for i, line := range pr.Lines {
				candidates[i] = domain.PrizeCandidate{Category: line.Category, Club: line.Club}
				if line.Club != "" {
					clubs[line.Club]++
				}
			}
			won := docs.Section{
				Heading: pr.Title,
				Note:    i18n.Tf(l, "Ranks 1–%d per category with at least %d participants, club with at least %d in this programme", rules.Ranks, rules.MinCategory, rules.MinClub),
				Columns: []docs.Column{col(i18n.T(l, "Category"), 2, "L"), col(i18n.T(l, "Prize"), 1, "C"), col(i18n.T(l, "Rank"), 1, "C"),
					col(i18n.T(l, "Name"), 3, "L"), col(i18n.T(l, "Year"), 1, "C"), col(i18n.T(l, "Club"), 4, "L"), col(i18n.T(l, "Total"), 1.2, "R")},
			}
			none := docs.Section{
				Heading: i18n.Tf(l, "%s: no special prize", pr.Title),
				Columns: []docs.Column{col(i18n.T(l, "Category"), 2, "L"), col(i18n.T(l, "Rank"), 1, "C"), col(i18n.T(l, "Name"), 3, "L"),
					col(i18n.T(l, "Club"), 4, "L"), col(i18n.T(l, "Reason"), 4, "L")},
			}
			for _, cp := range domain.SpecialPrizes(candidates, rules) {
				cat := categoryLabel(l, cp.Category)
				if cp.Participants < rules.MinCategory {
					none.Rows = append(none.Rows, []string{cat, "", "", "", i18n.Tf(l, "Participants in the category: %d", cp.Participants)})
					continue
				}
				for n, i := range cp.Winners {
					line := pr.Lines[i]
					won.Rows = append(won.Rows, []string{cat, strconv.Itoa(n + 1), strconv.Itoa(line.Rank), line.Lastname + " " + line.Firstname,
						line.BirthYear, line.Club, docs.Points(line.Total)})
				}
				for _, i := range cp.Skipped {
					line := pr.Lines[i]
					reason := i18n.T(l, "No club")
					if line.Club != "" {
						reason = i18n.Tf(l, "Participants of the club: %d", clubs[line.Club])
					}
					none.Rows = append(none.Rows, []string{cat, strconv.Itoa(line.Rank), line.Lastname + " " + line.Firstname, line.Club, reason})
				}
			}
			r.Sections = append(r.Sections, won)
			if len(none.Rows) > 0 {
				r.Sections = append(r.Sections, none)
			}
		}
		if len(r.Sections) == 0 {
			r.Summary = [][2]string{{i18n.T(l, "No programme awards special prizes."), ""}}
		}
	case ReportBookings:
		r.Title = i18n.T(l, "Lane bookings")
		names := map[int64]string{}
		for _, rg := range d.event.Ranges {
			names[rg.ID] = rg.Name
		}
		loc := eventLoc(d.event.Timezone)
		clubs, pending := 0, 0
		for _, b := range d.bookings {
			if b.Status != store.BookingConfirmed && b.Status != store.BookingPending {
				continue
			}
			clubs++
			status := i18n.T(l, "confirmed")
			if b.Status == store.BookingPending {
				status = i18n.T(l, "not confirmed yet")
				pending++
			}
			slots := docs.Section{
				Heading: b.Club + " (" + status + ")",
				Note:    strings.Join(slices.DeleteFunc([]string{b.Contact, b.Email, b.Phone, b.Note}, func(x string) bool { return x == "" }), " · "),
				Columns: []docs.Column{col(i18n.T(l, "Time"), 3, "L"), col(i18n.T(l, "Range"), 3, "L"), col(i18n.T(l, "Lanes"), 1, "C"), col(i18n.T(l, "Lane numbers"), 3, "L")},
			}
			for _, it := range b.Items {
				slots.Rows = append(slots.Rows, []string{slotText(it.Start, it.End, loc), names[it.RangeID], strconv.Itoa(it.Lanes), domain.FormatLanes(it.Assigned)})
			}
			r.Sections = append(r.Sections, slots)
			if len(b.Shooters) > 0 {
				sh := docs.Section{Heading: i18n.Tf(l, "Shooters of %s", b.Club),
					Columns: []docs.Column{col(i18n.T(l, "Name"), 4, "L"), col(i18n.T(l, "Year"), 1, "C"), col(i18n.T(l, "Licence"), 2, "L")}}
				for _, x := range b.Shooters {
					sh.Rows = append(sh.Rows, []string{x.Lastname + " " + x.Firstname, x.Birthday[:min(4, len(x.Birthday))], x.LicenseNumber})
				}
				r.Sections = append(r.Sections, sh)
			}
		}
		r.Summary = [][2]string{{i18n.T(l, "Clubs"), strconv.Itoa(clubs)}, {i18n.T(l, "Not confirmed yet"), strconv.Itoa(pending)}}
	case ReportClassifications:
		r.Title = i18n.T(l, "Classifications")
		for _, c := range d.ranked {
			cols := []docs.Column{col(i18n.T(l, "Rank"), 1, "C"), col(i18n.T(l, "Name"), 4, "L"), col(i18n.T(l, "Year"), 1, "C"),
				col(i18n.T(l, "Category"), 2, "L"), col(i18n.T(l, "Club"), 4, "L")}
			for _, p := range c.Parts {
				cols = append(cols, col(p.Title, 1.5, "R"))
			}
			cols = append(cols, col(i18n.T(l, "Total"), 1.5, "R"))
			for _, cl := range c.Classes {
				sec := docs.Section{Heading: c.Name, Columns: cols}
				if cl.Name != "" {
					sec.Heading = c.Name + " – " + cl.Name
				}
				for _, group := range [][]ClassificationLine{cl.Lines, cl.Unranked} {
					for _, line := range group {
						rank := ""
						if line.Rank > 0 {
							rank = strconv.Itoa(line.Rank)
						}
						row := []string{rank, line.Lastname + " " + line.Firstname, line.BirthYear, categoryLabel(l, line.Category), line.Club}
						for _, p := range line.Parts {
							v := "–"
							if p.Points != nil {
								v = docs.Points(*p.Points)
								if !p.Counted {
									v = "(" + v + ")"
								}
							}
							row = append(row, v)
						}
						sec.Rows = append(sec.Rows, append(row, docs.Points(line.Total)))
					}
				}
				r.Sections = append(r.Sections, sec)
			}
		}
		if len(r.Sections) == 0 {
			r.Summary = [][2]string{{i18n.T(l, "The event has no classifications."), ""}}
		}
	default:
		return r, "", missing("Unknown report")
	}
	return r, filename(kind, d.event.Name), nil
}

// col returns a report column.
func col(title string, width float64, align string) docs.Column {
	return docs.Column{Title: title, Width: width, Align: align}
}

// awardFacts returns the facts of entries that have an award.
func awardFacts(facts []store.EntryFact) []store.EntryFact {
	var out []store.EntryFact
	for _, f := range facts {
		if f.AwardTierID != nil {
			out = append(out, f)
		}
	}
	return out
}

// FileOutput is a downloaded document.
type FileOutput struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}

// fileOut renders a report as CSV or PDF.
func fileOut(r docs.Report, name, format string) (*FileOutput, error) {
	if format == "csv" {
		b, err := docs.ReportCSV(r)
		return &FileOutput{ContentType: "text/csv; charset=utf-8", ContentDisposition: `attachment; filename="` + strings.TrimSuffix(name, ".pdf") + `.csv"`, Body: b}, err
	}
	b, err := docs.RenderReport(r)
	return &FileOutput{ContentType: "application/pdf", ContentDisposition: `inline; filename="` + name + `"`, Body: b}, err
}

// eventNow returns the current time in the event's time zone.
func eventNow(ev store.Event) time.Time {
	return time.Now().In(eventLoc(ev.Timezone))
}

// registerReports registers the statistics, club competition, report and ranking CSV endpoints.
func (s *Server) registerReports(api huma.API) {
	type statsOut struct{ Body Statistics }
	huma.Register(api, huma.Operation{
		OperationID: "get-statistics", Metadata: needs("", permissions.ReportsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/statistics", Tags: []string{"reports"},
	}, func(ctx context.Context, in *EventPath) (*statsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		d, err := s.loadEventData(ctx, in.TenantID, in.EventID)
		if err != nil {
			return nil, err
		}
		return &statsOut{Body: d.statistics()}, nil
	})

	type clubOut struct{ Body ClubCompetitionView }
	huma.Register(api, huma.Operation{
		OperationID: "get-club-competition", Metadata: needs(modules.ClubCompetition, permissions.ReportsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/club-competition", Tags: []string{"reports"},
	}, func(ctx context.Context, in *EventPath) (*clubOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		d, err := s.loadEventData(ctx, in.TenantID, in.EventID)
		if err != nil {
			return nil, err
		}
		return &clubOut{Body: d.clubCompetition()}, nil
	})

	type reportIn struct {
		EventPath
		Report string `path:"report" enum:"financial,revenue,awards,club_competition,club_participants,special_prizes,awards_open,bookings,classifications"`
		Format string `query:"format" enum:"pdf,csv" default:"pdf"`
	}
	fileResponses := map[string]*huma.Response{"200": {Description: "PDF or CSV", Content: map[string]*huma.MediaType{"application/pdf": {}, "text/csv": {}}}}
	huma.Register(api, huma.Operation{
		OperationID: "get-report", Metadata: needs("", permissions.ReportsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/reports/{report}", Tags: []string{"reports"},
		Responses: fileResponses,
	}, func(ctx context.Context, in *reportIn) (*FileOutput, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		d, err := s.loadEventData(ctx, in.TenantID, in.EventID)
		if err != nil {
			return nil, err
		}
		if k := reportModules[in.Report]; k != "" {
			if err := requireModule(d.event, k); err != nil {
				return nil, err
			}
		}
		if in.Report == ReportBookings {
			if d.bookings, err = s.Store.ListBookingRequests(ctx, in.TenantID, in.EventID); err != nil {
				return nil, err
			}
		}
		if in.Report == ReportClassifications {
			if d.ranked, err = s.classificationRankings(ctx, in.TenantID, in.EventID, false); err != nil {
				return nil, err
			}
		}
		if in.Report == ReportSpecialPrizes {
			if d.rankings, err = s.rankings(ctx, in.TenantID, in.EventID, false); err != nil {
				return nil, err
			}
		}
		r, name, err := d.report(in.Report, eventNow(d.event))
		if err != nil {
			return nil, err
		}
		return fileOut(r, name, in.Format)
	})

	type rankingCSVIn struct {
		ProgrammePath
	}
	huma.Register(api, huma.Operation{
		OperationID: "ranking-csv", Metadata: needs("", permissions.ReportsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/ranking.csv", Tags: []string{"reports"},
		Responses: fileResponses,
	}, func(ctx context.Context, in *rankingCSVIn) (*FileOutput, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		p, err := s.Store.GetProgramme(ctx, in.TenantID, in.ProgrammeID)
		if err != nil {
			return nil, notFound(err, "Programme not found")
		}
		ev, err := s.Store.GetEvent(ctx, in.TenantID, p.EventID)
		if err != nil {
			return nil, err
		}
		rankings, err := s.rankings(ctx, in.TenantID, p.EventID, false)
		if err != nil {
			return nil, err
		}
		l := i18n.Get(ev.Locale)
		cols := []docs.Column{col(i18n.T(l, "Rank"), 1, "C"), col(i18n.T(l, "Last name"), 2, "L"), col(i18n.T(l, "First name"), 2, "L"), col(i18n.T(l, "Year"), 1, "C"), col(i18n.T(l, "Category"), 1, "L"), col(i18n.T(l, "Club"), 3, "L")}
		for i := range p.Series {
			cols = append(cols, col(i18n.Tf(l, "S%d", i+1), 1, "R"))
		}
		cols = append(cols, col(i18n.T(l, "Total"), 1, "R"), col(i18n.T(l, "Mouches"), 1, "R"))
		repeatable := p.AttemptLimit() != 1
		if repeatable {
			cols = append(cols, col(i18n.T(l, "Attempt"), 1, "C"), col(i18n.T(l, "Attempts"), 1, "C"))
		}
		sec := docs.Section{Columns: cols}
		for _, r := range rankings {
			if r.ProgrammeID != p.ID {
				continue
			}
			for _, line := range r.Lines {
				row := []string{strconv.Itoa(line.Rank), line.Lastname, line.Firstname, line.BirthYear, categoryLabel(l, line.Category), line.Club}
				for i := range p.Series {
					v := ""
					if i < len(line.Series) {
						v = docs.Points(line.Series[i])
					}
					row = append(row, v)
				}
				row = append(row, docs.Points(line.Total), strconv.Itoa(line.Mouches))
				if repeatable {
					row = append(row, strconv.Itoa(line.Attempt), strconv.Itoa(line.Attempts))
				}
				sec.Rows = append(sec.Rows, row)
			}
		}
		title := programmeTitle(p.Name, p.Variant)
		r := docs.Report{Lang: l, Title: i18n.Tf(l, "Ranking %s", title), Subtitle: ev.Name + " · " + swissDate(ev.StartsOn), Sections: []docs.Section{sec}}
		return fileOut(r, filename("Ranking", p.Name, p.Variant), "csv")
	})
}
