package store_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func TestBookingConfirm(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Host")
	rng, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "10m", LaneCount: 6, DisabledLanes: []int32{2}})
	ev, err := st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "Lupi", StartsOn: "2027-03-13", Timezone: "Europe/Zurich"})
	if err != nil {
		t.Fatal(err)
	}
	at := func(h int) time.Time { return time.Date(2027, 3, 13, h, 0, 0, 0, time.UTC) }
	if _, err := st.SaveReservation(ctx, tid, store.Reservation{RangeID: rng.ID, Lanes: []int{1}, StartsAt: at(8), EndsAt: at(9), Description: "Members"}); err != nil {
		t.Fatal(err)
	}
	slot := func(h, lanes int) domain.BookingItem {
		return domain.BookingItem{RangeID: rng.ID, Start: at(h), End: at(h).Add(30 * time.Minute), Lanes: lanes}
	}

	// Known shooters: one by licence (stored with leading zeros and another name), one by name and birthday.
	known, _ := st.CreateShooter(ctx, tid, store.ShooterInput{LicenseNumber: "012345", Firstname: "Anna", Lastname: "Stored", Birthday: ""}, nil)
	byName, _ := st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Bert", Lastname: "Muster", Birthday: "1980-02-02"}, nil)
	in := store.BookingRequestInput{
		Club: "sv gäste", Contact: "Ida", Email: "ida@example.org",
		Items: []domain.BookingItem{slot(8, 3)},
		Shooters: []domain.BookingShooter{
			{Firstname: "Hacked", Lastname: "Name", Birthday: "1990-01-01", LicenseNumber: "12345"},
			{Firstname: "bert", Lastname: "MUSTER", Birthday: "1980-02-02", LicenseNumber: "777"},
			{Firstname: "Cleo", Lastname: "Neu", Birthday: "2005-05-05"},
		},
	}
	guests, _ := st.EnsureClub(ctx, tid, "SV Gäste")
	id, err := st.CreateBookingRequest(ctx, tid, ev, in, "hash-1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := st.GetBookingRequest(ctx, tid, id)
	knownFlags := ""
	for _, s := range req.Shooters {
		knownFlags += fmt.Sprintf("%s:%v ", s.Firstname, s.Known)
	}
	if req.Status != store.BookingPending || len(req.Items) != 1 || knownFlags != "bert:true Hacked:true Cleo:false " || req.Shooters[1].LicenseNumber != "012345" {
		t.Fatalf("request: %+v / %s", req, knownFlags)
	}

	if err := st.ConfirmBooking(ctx, tid, id); err != nil {
		t.Fatal(err)
	}
	if err := st.ConfirmBooking(ctx, tid, id); !errors.Is(err, store.ErrNotPending) {
		t.Fatalf("confirm twice: %v", err)
	}
	req, _ = st.GetBookingRequest(ctx, tid, id)
	if req.Status != store.BookingConfirmed || fmt.Sprint(req.Items[0].Assigned) != "[3 4 5]" {
		t.Fatalf("confirmed: %+v", req.Items)
	}
	anna, _ := st.GetShooter(ctx, tid, known)
	bert, _ := st.GetShooter(ctx, tid, byName)
	if anna.Lastname != "Stored" || anna.Birthday != "1990-01-01" || bert.LicenseNumber != "000777" {
		t.Fatalf("master data overwritten or not filled: %+v %+v", anna, bert)
	}
	found, _ := st.SearchShooters(ctx, tid, "Cleo Neu", guests, 10, 0)
	if len(found) != 1 {
		t.Fatalf("new shooter not created in the existing club (case-insensitive): %+v", found)
	}

	// Only shooters added: still confirmed, the new one waits for adoption.
	in.Shooters = append(in.Shooters, domain.BookingShooter{Firstname: "Dora", Lastname: "Spät", Birthday: "1970-07-07"})
	if reopened, err := st.UpdateBookingRequest(ctx, tid, id, in); err != nil || reopened {
		t.Fatalf("update shooters: %v %v", reopened, err)
	}
	req, _ = st.GetBookingRequest(ctx, tid, id)
	if req.Status != store.BookingConfirmed || req.Shooters[0].ShooterID == nil || req.Shooters[3].ShooterID != nil {
		t.Fatalf("after adding a shooter: %+v", req)
	}
	if err := st.AdoptBookingShooters(ctx, tid, id); err != nil {
		t.Fatal(err)
	}

	// A second club fills the slot: confirming it fails with the full item, nothing is stored.
	other, _ := st.CreateBookingRequest(ctx, tid, ev, store.BookingRequestInput{Club: "Other", Contact: "Max", Email: "max@example.org",
		Items: []domain.BookingItem{slot(9, 1), slot(8, 2)}}, "hash-2")
	var full *store.FullSlotError
	if err := st.ConfirmBooking(ctx, tid, other); !errors.As(err, &full) || full.Free != 1 || !full.Item.Start.Equal(at(8)) {
		t.Fatalf("full slot: %v", err)
	}
	if req, _ := st.GetBookingRequest(ctx, tid, other); req.Status != store.BookingPending || len(req.Items[1].Assigned) != 0 {
		t.Fatalf("partial confirmation stored: %+v", req)
	}

	// Changing the slots of a confirmed request frees its lanes and reopens it.
	in.Items = []domain.BookingItem{slot(8, 1)}
	if reopened, err := st.UpdateBookingRequest(ctx, tid, id, in); err != nil || !reopened {
		t.Fatalf("update items: %v %v", reopened, err)
	}
	if err := st.ConfirmBooking(ctx, tid, other); err != nil {
		t.Fatalf("lanes not freed: %v", err)
	}
	if err := st.CancelBookingRequest(ctx, tid, other); err != nil {
		t.Fatal(err)
	}
	if list, _ := st.ReservationsBetween(ctx, tid, []int64{rng.ID}, at(0), at(23), 0); len(list) != 1 {
		t.Fatalf("cancel kept reservations: %+v", list)
	}
	if err := st.CancelBookingRequest(ctx, tid, other); !errors.Is(err, store.ErrNotPending) {
		t.Fatalf("cancel twice: %v", err)
	}
	if _, err := st.UpdateBookingRequest(ctx, tid, other, in); !errors.Is(err, store.ErrNotPending) {
		t.Fatalf("update cancelled: %v", err)
	}
}

func TestBookingSettingsCopy(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Host")
	rng, _ := st.CreateRange(ctx, tid, store.RangeInput{Name: "300m", LaneCount: 10})
	ev, _ := st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "Feldschiessen", StartsOn: "2027-05-29",
		Schedule: domain.Schedule{From: "2027-04-01", To: "2027-08-31", Weekly: []domain.WeeklyTimes{{Weekday: 3}},
			Days: []domain.ScheduleDay{{Date: "2027-05-28", Times: []domain.BookingHours{{From: "17:00", To: "19:00"}}}}}, Timezone: "Europe/Zurich"})
	set := domain.BookingSettings{From: "09:00", SlotMinutes: 45, Deadline: "2027-05-01",
		Days: map[string]domain.BookingHours{"2027-05-28": {From: "17:00", To: "20:00"}}}
	if err := st.SetBookingSettings(ctx, tid, ev, set); err != nil {
		t.Fatal(err)
	}
	if err := st.SetBookingCode(ctx, tid, ev, "abc"); err != nil {
		t.Fatal(err)
	}
	if target, err := st.BookingByCode(ctx, "abc"); err != nil || target.Event.ID != ev || target.Settings.SlotMinutes != 45 {
		t.Fatalf("by code: %+v %v", target, err)
	}
	copied, err := st.CopyEvent(ctx, tid, ev, "Feldschiessen 2028", "2028-06-03")
	if err != nil {
		t.Fatal(err)
	}
	code, got, err := st.BookingConfig(ctx, tid, copied)
	if err != nil || code != "" || got.Deadline != "" || got.From != "09:00" || got.SlotMinutes != 45 || got.Days["2028-06-02"].From != "17:00" || len(got.Days) != 1 {
		t.Fatalf("copied booking: %q %+v %v", code, got, err)
	}
	cev, err := st.GetEvent(ctx, tid, copied)
	if err != nil || cev.Schedule.From != "2028-04-06" || cev.Days[0] != "2028-04-12" || !slices.Contains(cev.Days, "2028-06-02") ||
		!slices.Contains(cev.Days, "2028-06-03") || cev.Schedule.Times()["2028-06-02"][0].From != "17:00" {
		t.Fatalf("copied schedule: %+v %v", cev, err)
	}
	if err := st.SetEventArchived(ctx, tid, ev, true); err != nil {
		t.Fatal(err)
	}
	if _, err := st.BookingByCode(ctx, "abc"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("archived event still bookable: %v", err)
	}
}
