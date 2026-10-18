package store_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func fieldMatch() domain.ProgrammeDefinition {
	return domain.ProgrammeDefinition{
		Name: "Field match", Variant: "Stgw 90", Barcode: "123", CountsForClubCompetition: true, ShowOnRankings: true,
		Series: []domain.Series{{Shots: 6, MaxPoints: 60, Name: "A"}, {Shots: 6, MaxPoints: 60, TrackMouches: true}},
		Prices: []domain.Price{{AgeCategory: domain.U21, Price: 0}, {AgeCategory: domain.Elite, Price: 20.5}},
		Awards: []domain.AwardTier{
			{Name: "Wreath", Kind: domain.AwardCard, Value: 25, Limits: []domain.AwardLimit{{AgeCategory: domain.U21, Points: 95}, {AgeCategory: domain.Elite, Points: 100}}},
			{Name: "Badge", Kind: domain.AwardNatural, Value: 0, Limits: []domain.AwardLimit{}},
		},
	}
}

func eventFixture(t *testing.T, st *store.Store, tenant string) (int64, int64) {
	t.Helper()
	ctx := context.Background()
	tid, err := st.CreateTenant(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	rng, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "300m", LaneCount: 10})
	ev, err := st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "Event", StartsOn: "2026-06-01", Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	return tid, ev
}

func TestProgrammeRoundTrip(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, ev := eventFixture(t, st, "Club")
	def := fieldMatch()
	repeat := 12.0
	def.Attempts, def.AttemptRule, def.AwardCountingOnly = 3, domain.AttemptLast, true
	def.Prices[1].RepeatPrice = &repeat

	id, err := st.CreateProgramme(ctx, tid, ev, def)
	if err != nil {
		t.Fatal(err)
	}
	p, err := st.GetProgramme(ctx, tid, id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.ProgrammeDefinition, def) || p.Shots != 12 || p.MaxPoints != 120 || !p.Active || p.EventID != ev {
		t.Fatalf("round trip:\n got %+v\nwant %+v", p, def)
	}
	if _, err := st.CreateProgramme(ctx, tid, ev, def); !errors.Is(err, store.ErrExists) {
		t.Fatalf("duplicate name and variant: %v", err)
	}

	def.Series = def.Series[:1]
	def.Awards = nil
	def.Prices = []domain.Price{{AgeCategory: domain.Senior, Price: 15}}
	if err := st.UpdateProgramme(ctx, tid, id, def); err != nil {
		t.Fatal(err)
	}
	p, _ = st.GetProgramme(ctx, tid, id)
	if len(p.Series) != 1 || len(p.Awards) != 0 || len(p.Prices) != 1 || p.Prices[0].AgeCategory != domain.Senior {
		t.Fatalf("after update: %+v", p)
	}
}

func TestProgrammeIsolationCopyAndCombination(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, ev := eventFixture(t, st, "Club")
	other, otherEv := eventFixture(t, st, "Other")

	a, _ := st.CreateProgramme(ctx, tid, ev, fieldMatch())
	sec := fieldMatch()
	sec.Variant = "Stgw 57"
	b, _ := st.CreateProgramme(ctx, tid, ev, sec)
	foreign, _ := st.CreateProgramme(ctx, other, otherEv, fieldMatch())

	if _, err := st.CreateProgramme(ctx, other, ev, fieldMatch()); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("programme in foreign event: %v", err)
	}
	if _, err := st.GetProgramme(ctx, other, a); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("read across tenants: %v", err)
	}
	if err := st.SetCombination(ctx, tid, a, &foreign); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("combine with foreign programme: %v", err)
	}
	if err := st.SetCombination(ctx, tid, a, &a); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("combine with itself: %v", err)
	}
	if err := st.SetCombination(ctx, tid, a, &b); err != nil {
		t.Fatal(err)
	}

	ranges, _ := st.ListRanges(ctx, tid)
	next, _ := st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{ranges[0].ID}, Name: "Next year", StartsOn: "2027-06-01", Timezone: "UTC"})
	n, err := st.CopyProgrammes(ctx, tid, ev, next)
	if err != nil || n != 2 {
		t.Fatalf("copy: %d %v", n, err)
	}
	if n, _ = st.CopyProgrammes(ctx, tid, ev, next); n != 0 {
		t.Fatalf("second copy should skip existing programmes, copied %d", n)
	}
	if _, err := st.CopyProgrammes(ctx, tid, ev, otherEv); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("copy into foreign event: %v", err)
	}
	copies, _ := st.ListProgrammes(ctx, tid, next)
	byVariant := map[string]store.Programme{}
	for _, p := range copies {
		byVariant[p.Variant] = p
	}
	c90, c57 := byVariant["Stgw 90"], byVariant["Stgw 57"]
	if c90.CombinableWithID == nil || *c90.CombinableWithID != c57.ID || len(c90.Awards) != 2 {
		t.Fatalf("copied combination: %+v", c90)
	}
	if err := st.DeleteProgramme(ctx, tid, c57.ID); err != nil {
		t.Fatal(err)
	}
	if p, _ := st.GetProgramme(ctx, tid, c90.ID); p.CombinableWithID != nil {
		t.Fatalf("combination survived deletion of its target: %+v", p.CombinableWithID)
	}
}

func TestTemplates(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	id, err := st.SaveTemplate(ctx, fieldMatch())
	if err != nil {
		t.Fatal(err)
	}
	changed := fieldMatch()
	changed.Description = "updated"
	if again, _ := st.SaveTemplate(ctx, changed); again != id {
		t.Fatalf("same name and variant should replace the template: %d vs %d", again, id)
	}
	tpl, err := st.GetTemplate(ctx, id)
	if err != nil || tpl.Definition.Description != "updated" || len(tpl.Definition.Awards) != 2 {
		t.Fatalf("template: %+v %v", tpl, err)
	}
	list, _ := st.ListTemplates(ctx)
	if len(list) != 1 {
		t.Fatalf("templates: %+v", list)
	}
}
