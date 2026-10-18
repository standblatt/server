package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func TestAttempts(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, ev := eventFixture(t, st, "Club")
	def := fieldMatch()
	repeat := 5.0
	def.Attempts = 3
	def.Prices[1].RepeatPrice = &repeat
	pid, err := st.CreateProgramme(ctx, tid, ev, def)
	if err != nil {
		t.Fatal(err)
	}
	anna, _ := st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Muster"}, nil)
	register := func(count int) (int64, error) {
		return st.CreateRegistration(ctx, tid, store.NewRegistration{EventID: ev, ShooterID: anna, AgeCategory: domain.Elite,
			Programmes: []domain.ProgrammeChoice{{ProgrammeID: pid, Count: count}}})
	}
	attempts := func(regID int64) []int {
		t.Helper()
		r, err := st.GetRegistration(ctx, tid, regID)
		if err != nil {
			t.Fatal(err)
		}
		var out []int
		for _, e := range r.Entries {
			out = append(out, e.Attempt)
		}
		return out
	}

	first, err := register(2)
	if err != nil {
		t.Fatal(err)
	}
	if got := attempts(first); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("first registration: %v", got)
	}
	if _, err := register(2); !errors.Is(err, store.ErrProgrammeTaken) {
		t.Fatalf("beyond the limit: %v", err)
	}
	second, err := register(1)
	if err != nil {
		t.Fatal(err)
	}
	if got := attempts(second); len(got) != 1 || got[0] != 3 {
		t.Fatalf("second registration: %v", got)
	}

	reg, _ := st.GetRegistration(ctx, tid, first)
	if _, err := st.SaveResults(ctx, tid, store.EntryRef{ID: reg.Entries[0].ID, RegistrationID: first, EventID: ev, ProgrammeID: pid},
		[]domain.SeriesResult{{Points: 50}, {Points: 50}}); err != nil {
		t.Fatal(err)
	}
	if err := st.UpdateRegistration(ctx, tid, reg, nil, nil, []domain.ProgrammeChoice{{ProgrammeID: pid, Count: 1}}); err != nil {
		t.Fatal(err)
	}
	if got := attempts(first); len(got) != 1 || got[0] != 1 {
		t.Fatalf("after counting down, the scored attempt stays: %v", got)
	}
	reg, _ = st.GetRegistration(ctx, tid, first)
	if err := st.UpdateRegistration(ctx, tid, reg, nil, nil, []domain.ProgrammeChoice{{ProgrammeID: pid, Count: 2}}); err != nil {
		t.Fatal(err)
	}
	if got := attempts(first); len(got) != 2 || got[1] != 4 {
		t.Fatalf("counting up again numbers after the last attempt: %v", got)
	}
	reg, _ = st.GetRegistration(ctx, tid, first)
	if err := st.UpdateRegistration(ctx, tid, reg, nil, nil, []domain.ProgrammeChoice{}); !errors.Is(err, store.ErrHasResults) {
		t.Fatalf("dropping the scored attempt: %v", err)
	}

	facts, err := st.ListEntryFacts(ctx, tid, ev)
	if err != nil {
		t.Fatal(err)
	}
	prices := map[int]float64{}
	for _, f := range facts {
		prices[f.Attempt] = f.Price
	}
	if prices[1] != 20.5 || prices[3] != 5 || prices[4] != 5 {
		t.Fatalf("prices by attempt: %v", prices)
	}

	def.Normalize()
	def.Attempts = 2
	if err := st.UpdateProgramme(ctx, tid, pid, def); !errors.Is(err, store.ErrAttemptsUsed) {
		t.Fatalf("limit below the attempts shot: %v", err)
	}
	def.Attempts, def.AttemptRule = 5, domain.AttemptFirst
	if err := st.UpdateProgramme(ctx, tid, pid, def); err != nil {
		t.Fatalf("raising the limit while in use: %v", err)
	}
}

func TestCountingRows(t *testing.T) {
	rows := []store.RankingRow{
		{ProgrammeID: 1, ShooterID: 1, EntryID: 11, Attempt: 1, Rule: domain.AttemptBest, Series: []float64{90}},
		{ProgrammeID: 1, ShooterID: 1, EntryID: 12, Attempt: 2, Rule: domain.AttemptBest, Series: []float64{95}},
		{ProgrammeID: 1, ShooterID: 2, EntryID: 21, Attempt: 1, Rule: domain.AttemptBest, Series: []float64{80}},
		{ProgrammeID: 2, ShooterID: 1, EntryID: 31, Attempt: 1, Rule: domain.AttemptLast, Series: []float64{99}},
		{ProgrammeID: 2, ShooterID: 1, EntryID: 32, Attempt: 2, Rule: domain.AttemptLast, Series: []float64{70}},
	}
	got := store.CountingRows(rows)
	want := map[int64]int{12: 2, 21: 1, 32: 2}
	if len(got) != len(want) {
		t.Fatalf("got %d rows: %+v", len(got), got)
	}
	for _, r := range got {
		if n, ok := want[r.EntryID]; !ok || r.Attempts != n {
			t.Errorf("entry %d with %d attempts not expected", r.EntryID, r.Attempts)
		}
	}
}

func TestClassificationsCopyAndCleanup(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, ev := eventFixture(t, st, "Club")
	keep, _ := st.CreateProgramme(ctx, tid, ev, fieldMatch())
	other := fieldMatch()
	other.Name = "Other"
	drop, _ := st.CreateProgramme(ctx, tid, ev, other)
	c := domain.ClassificationDefinition{Name: "Section", TieBreak: domain.SSVTieBreak,
		Parts: []domain.ClassificationPart{{ProgrammeID: keep, Count: 1}, {ProgrammeID: drop, Count: 2}}}
	if _, err := st.CreateClassification(ctx, tid, ev, c); err != nil {
		t.Fatal(err)
	}
	copied, err := st.CopyEvent(ctx, tid, ev, "Next year", "2027-06-01", false)
	if err != nil {
		t.Fatal(err)
	}
	list, err := st.ListClassifications(ctx, tid, copied)
	if err != nil || len(list) != 1 || len(list[0].Parts) != 2 || list[0].Parts[0].ProgrammeID == keep || list[0].Parts[1].Count != 2 {
		t.Fatalf("copied classification: %+v %v", list, err)
	}
	if err := st.DeleteProgramme(ctx, tid, drop); err != nil {
		t.Fatal(err)
	}
	list, _ = st.ListClassifications(ctx, tid, ev)
	if len(list[0].Parts) != 1 || list[0].Parts[0].ProgrammeID != keep {
		t.Fatalf("a deleted programme leaves the classification: %+v", list[0].Parts)
	}
}
