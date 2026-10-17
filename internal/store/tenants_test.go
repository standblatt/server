package store_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func TestTenantIsolation(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()

	a, err := st.CreateTenant(ctx, "Club A")
	if err != nil {
		t.Fatal(err)
	}
	b, err := st.CreateTenant(ctx, "Club B")
	if err != nil {
		t.Fatal(err)
	}
	rangeA, err := st.CreateRange(ctx, a, store.RangeInput{Name: "300m", LaneCount: 10})
	if err != nil {
		t.Fatal(err)
	}
	evA, err := st.CreateEvent(ctx, a, store.EventInput{RangeIDs: []int64{rangeA.ID}, Name: "Field match", StartsOn: "2026-06-01", Timezone: "Europe/Zurich"})
	if err != nil {
		t.Fatal(err)
	}

	if list, err := st.ListRanges(ctx, b); err != nil || len(list) != 0 {
		t.Fatalf("tenant B sees ranges of A: %v %v", list, err)
	}
	if list, err := st.ListEvents(ctx, b); err != nil || len(list) != 0 {
		t.Fatalf("tenant B sees events of A: %v %v", list, err)
	}
	if _, err := st.GetEvent(ctx, b, evA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("GetEvent across tenants: %v", err)
	}
	if err := st.UpdateRange(ctx, b, rangeA.ID, store.RangeInput{Name: "x", LaneCount: 1}); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("UpdateRange across tenants: %v", err)
	}
	if err := st.DeleteEvent(ctx, b, evA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("DeleteEvent across tenants: %v", err)
	}
	if _, err := st.CreateEvent(ctx, b, store.EventInput{RangeIDs: []int64{rangeA.ID}, Name: "x", StartsOn: "2026-06-01", Timezone: "UTC"}); !errors.Is(err, store.ErrInUse) {
		t.Fatalf("event of B on range of A: %v", err)
	}
	if err := st.DeleteRange(ctx, a, rangeA.ID); !errors.Is(err, store.ErrInUse) {
		t.Fatalf("delete range used by event: %v", err)
	}
	codeA, err := st.CreateHelperCode(ctx, a, evA, store.HelperCode{Name: "Office", Role: store.HelperResults}, "hash-a")
	if err != nil {
		t.Fatal(err)
	}
	if list, err := st.ListHelperCodes(ctx, b, evA); err != nil || len(list) != 0 {
		t.Fatalf("tenant B sees helper codes of A: %v %v", list, err)
	}
	if _, err := st.DeleteHelperCode(ctx, b, codeA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("DeleteHelperCode across tenants: %v", err)
	}
	if _, err := st.CreateHelperCode(ctx, b, evA, store.HelperCode{Name: "x", Role: store.HelperResults}, "hash-b"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("helper code of B on event of A: %v", err)
	}
	slot := domain.BookingItem{RangeID: rangeA.ID, Start: time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC), End: time.Date(2026, 6, 1, 8, 30, 0, 0, time.UTC), Lanes: 2}
	booking := store.BookingRequestInput{Club: "Guests", Contact: "Ida", Email: "ida@example.org", Items: []domain.BookingItem{slot}}
	reqA, err := st.CreateBookingRequest(ctx, a, evA, booking, "booking-a")
	if err != nil {
		t.Fatal(err)
	}
	if list, err := st.ListBookingRequests(ctx, b, evA); err != nil || len(list) != 0 {
		t.Fatalf("tenant B sees booking requests of A: %v %v", list, err)
	}
	if _, err := st.GetBookingRequest(ctx, b, reqA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("GetBookingRequest across tenants: %v", err)
	}
	if err := st.ConfirmBooking(ctx, b, reqA, nil); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("ConfirmBooking across tenants: %v", err)
	}
	clubB, _ := st.EnsureClub(ctx, b, "Club of B")
	if err := st.ConfirmBooking(ctx, a, reqA, &clubB); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("ConfirmBooking with a club of another tenant: %v", err)
	}
	if err := st.DeclineBookingRequest(ctx, b, reqA, ""); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("DeclineBookingRequest across tenants: %v", err)
	}
	if _, err := st.UpdateBookingRequest(ctx, b, reqA, booking); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("UpdateBookingRequest across tenants: %v", err)
	}
	if err := st.AdoptBookingShooters(ctx, b, reqA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("AdoptBookingShooters across tenants: %v", err)
	}
	if _, _, err := st.BookingConfig(ctx, b, evA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("BookingConfig across tenants: %v", err)
	}
	if err := st.SetBookingCode(ctx, b, evA, "code-b"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("SetBookingCode across tenants: %v", err)
	}
	if _, err := st.CreateBookingRequest(ctx, b, evA, booking, "booking-b"); !errors.Is(err, store.ErrInUse) {
		t.Fatalf("booking request of B on event of A: %v", err)
	}
	if err := st.DeleteTenant(ctx, a); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetEvent(ctx, a, evA); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("event survived tenant deletion: %v", err)
	}
}

func TestEventRoundTrip(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")
	rng, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "10m", LaneCount: 12, DisabledLanes: []int32{3, 7}})
	club := int64(42)
	id, err := st.CreateEvent(ctx, tid, store.EventInput{
		RangeIDs: []int64{rng.ID}, Name: "Indoor 2026", StartsOn: "2026-11-20", Organizer: "Club", Timezone: "Europe/Zurich",
		Settings: store.EventSettings{ExcludedClubID: &club, ShowLaneStatus: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	ev, err := st.GetEvent(ctx, tid, id)
	if err != nil {
		t.Fatal(err)
	}
	if ev.StartsOn != "2026-11-20" || ev.Ranges[0].Name != "10m" || ev.Settings.ExcludedClubID == nil || *ev.Settings.ExcludedClubID != 42 || !ev.Settings.ShowLaneStatus {
		t.Fatalf("unexpected event: %+v", ev)
	}
	ranges, _ := st.ListRanges(ctx, tid)
	if len(ranges) != 1 || len(ranges[0].DisabledLanes) != 2 {
		t.Fatalf("unexpected ranges: %+v", ranges)
	}
}

func TestProgrammeRange(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")
	other, _ := st.CreateTenant(ctx, "Other")
	r25, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "Pistol 25 m", LaneCount: 10})
	r50, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "Pistol 50 m", LaneCount: 5})
	spare, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "Spare", LaneCount: 2})
	foreign, _ := st.CreateRange(ctx, other, store.RangeInput{Name: "Other", LaneCount: 2})
	ev, _ := st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{r25.ID, r50.ID}, Name: "Pistol 2027", StartsOn: "2027-05-01", Timezone: "UTC"})
	p, _ := st.CreateProgramme(ctx, tid, ev, fieldMatch())
	if err := st.SetProgrammeRange(ctx, tid, p, &r50.ID); err != nil {
		t.Fatal(err)
	}
	if progs, _ := st.ListProgrammes(ctx, tid, ev); *progs[0].RangeID != r50.ID {
		t.Fatalf("programme on 50 m: %+v", progs[0])
	}
	if rules, _ := st.LoadEventRules(ctx, tid, ev); *rules.Programmes[p].RangeID != r50.ID || rules.Ranges[r50.ID].LaneCount != 5 {
		t.Fatalf("rules: %+v", rules)
	}
	if err := st.SetProgrammeRange(ctx, tid, p, &spare.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("range outside the event: %v", err)
	}
	if err := st.SetProgrammeRange(ctx, tid, p, &foreign.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("range of another tenant: %v", err)
	}
	if err := st.SetProgrammeRange(ctx, other, p, nil); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("programme of another tenant: %v", err)
	}
}
