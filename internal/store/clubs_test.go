package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func TestClubsMatchByKey(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")
	hill, err := st.CreateClub(ctx, tid, store.ClubInput{Name: "SG Hilltown"})
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := st.EnsureClub(ctx, tid, " sg  HILLTOWN"); id != hill.ID {
		t.Fatalf("EnsureClub ignores case and spaces: %d", id)
	}
	if _, err := st.CreateClub(ctx, tid, store.ClubInput{Name: "sg hilltown"}); !errors.Is(err, store.ErrExists) {
		t.Fatalf("create a club differing only in case: %v", err)
	}
	other, _ := st.CreateClub(ctx, tid, store.ClubInput{Name: "PC Lakeside"})
	if err := st.UpdateClub(ctx, tid, other.ID, store.ClubInput{Name: "SG HILLTOWN"}); !errors.Is(err, store.ErrExists) {
		t.Fatalf("rename to an existing club: %v", err)
	}
	if err := st.UpdateClub(ctx, tid, hill.ID, store.ClubInput{Name: "SG Hilltown 1890"}); err != nil {
		t.Fatalf("rename itself: %v", err)
	}
}

func TestMergeClubs(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")
	other, _ := st.CreateTenant(ctx, "Other")
	into, _ := st.CreateClub(ctx, tid, store.ClubInput{Name: "SG Hilltown"})
	from, _ := st.CreateClub(ctx, tid, store.ClubInput{Name: "S.G. Hiltown", SSVNumber: "1234"})
	foreign, _ := st.CreateClub(ctx, other, store.ClubInput{Name: "Foreign"})
	both, _ := st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Both"}, []int64{into.ID, from.ID})
	only, _ := st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Bert", Lastname: "Only"}, []int64{from.ID})
	rng, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "10m", LaneCount: 4})
	ev, err := st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "Old", StartsOn: "2025-11-20", Timezone: "UTC",
		Settings: store.EventSettings{ExcludedClubID: &from.ID}})
	if err != nil {
		t.Fatal(err)
	}
	p, _ := st.CreateProgramme(ctx, tid, ev, fieldMatch())
	reg, err := st.CreateRegistration(ctx, tid, store.NewRegistration{EventID: ev, ShooterID: only, ClubID: &from.ID, AgeCategory: domain.Elite,
		Programmes: []domain.ProgrammeChoice{{ProgrammeID: p}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.SetEventArchived(ctx, tid, ev, true); err != nil {
		t.Fatal(err)
	}

	if err := st.MergeClubs(ctx, tid, from.ID, foreign.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("merge into a club of another tenant: %v", err)
	}
	if err := st.MergeClubs(ctx, other, from.ID, into.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("merge in another tenant: %v", err)
	}
	if err := st.MergeClubs(ctx, tid, from.ID, from.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("merge into itself: %v", err)
	}
	if err := st.MergeClubs(ctx, tid, from.ID, into.ID); err != nil {
		t.Fatal(err)
	}

	clubs, _ := st.ListClubs(ctx, tid)
	if len(clubs) != 1 || clubs[0].Name != "SG Hilltown" || clubs[0].SSVNumber != "1234" || clubs[0].ShooterCount != 2 {
		t.Fatalf("clubs after merge: %+v", clubs)
	}
	if r, _ := st.GetRegistration(ctx, tid, reg); r.ClubID == nil || *r.ClubID != into.ID {
		t.Fatalf("registration of the closed event: %+v", r.ClubID)
	}
	if e, _ := st.GetEvent(ctx, tid, ev); e.Settings.ExcludedClubID == nil || *e.Settings.ExcludedClubID != into.ID || e.ArchivedAt == nil {
		t.Fatalf("event after merge: %+v", e)
	}
	if sh, _ := st.GetShooter(ctx, tid, both); len(sh.Clubs) != 1 {
		t.Fatalf("shooter in both clubs: %+v", sh.Clubs)
	}
	if _, err := st.CreateRegistration(ctx, tid, store.NewRegistration{EventID: ev, ShooterID: both, AgeCategory: domain.Elite,
		Programmes: []domain.ProgrammeChoice{{ProgrammeID: p}}}); !errors.Is(err, store.ErrArchived) {
		t.Fatalf("the closed event is open again after a merge: %v", err)
	}
}
