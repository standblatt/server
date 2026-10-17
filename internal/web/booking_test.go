package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

var manageRe = regexp.MustCompile(`/b/manage/(bk_[A-Za-z0-9_-]+)`)

// mailTo waits for the next mails until one goes to addr.
func (m mailbox) mailTo(t *testing.T, addr string) mail.Message {
	t.Helper()
	for range 10 {
		if msg := m.next(t); msg.To == addr {
			return msg
		}
	}
	t.Fatalf("no mail to %s", addr)
	return mail.Message{}
}

func TestOnlineBooking(t *testing.T) {
	e, box := personalEnv(t, nil)
	ctx := context.Background()
	tid := e.tenant("Host club")
	bossID := e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	if err := e.st.SetUserEmail(ctx, bossID, "boss@example.org", true); err != nil {
		t.Fatal(err)
	}
	admin, staff, anon := e.login("boss"), e.login("helper"), e.client()
	rng, _ := e.st.CreateRange(ctx, tid, store.RangeInput{Name: "10m", LaneCount: 4})
	day := time.Now().AddDate(0, 2, 0).Format(time.DateOnly)
	ev, err := e.st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "Lupi", StartsOn: day, Timezone: "Europe/Zurich", Locale: "en"})
	if err != nil {
		t.Fatal(err)
	}
	base := fmt.Sprintf("/api/tenants/%d/events/%d/booking", tid, ev)

	code, b := e.do(staff, "POST", base+"/code", nil)
	expect(t, "staff switches booking on", code, 403, b)
	code, b = e.do(admin, "PUT", base, domain.BookingSettings{From: "09:00", To: "08:00"})
	expect(t, "invalid hours", code, 422, b)
	code, b = e.do(admin, "PUT", base, domain.BookingSettings{From: "08:00", To: "10:00", SlotMinutes: 60})
	expect(t, "settings", code, 200, b)
	code, b = e.do(admin, "POST", base+"/code", nil)
	expect(t, "switch on", code, 200, b)
	var cfg web.BookingConfig
	json.Unmarshal(b, &cfg)
	if !cfg.Enabled || !strings.HasSuffix(cfg.Link, "/b/"+cfg.Code) || len(cfg.Slots) != 2 {
		t.Fatalf("config: %s", b)
	}

	public := "/api/booking/" + cfg.Code
	code, b = e.do(anon, "GET", "/api/booking/nope", nil)
	expect(t, "unknown code", code, 404, b)
	code, b = e.do(anon, "GET", public, nil)
	expect(t, "public info", code, 200, b)
	var info web.BookingInfo
	json.Unmarshal(b, &info)
	if !info.Open || len(info.Slots) != 2 || info.Slots[0].Free[0].Free != 4 || info.Event != "Lupi" {
		t.Fatalf("info: %s", b)
	}
	first, second := info.Slots[0].Slot, info.Slots[1].Slot

	request := store.BookingRequestInput{
		Club: "SV Gäste", Contact: "Ida Muster", Email: "Ida@Example.org",
		Items:    []domain.BookingItem{{RangeID: rng.ID, Start: first.Start, End: first.End, Lanes: 3}},
		Shooters: []domain.BookingShooter{{Firstname: "Ida", Lastname: "Muster", Birthday: "1990-05-01"}},
	}
	tooMany := request
	tooMany.Items = []domain.BookingItem{{RangeID: rng.ID, Start: first.Start, End: first.End, Lanes: 5}}
	code, b = e.do(anon, "POST", public+"/requests", tooMany)
	expect(t, "more lanes than the range has", code, 422, b)
	code, b = e.do(anon, "POST", public+"/requests", request)
	expect(t, "request", code, 204, b)

	clubMail := box.mailTo(t, "ida@example.org")
	m := manageRe.FindStringSubmatch(clubMail.Body)
	if m == nil || !strings.Contains(clubMail.Body, "Lanes: 3") {
		t.Fatalf("club mail: %+v", clubMail)
	}
	secret := m[1]
	adminMail := box.mailTo(t, "boss@example.org")
	if !strings.Contains(adminMail.Body, "SV Gäste") || !strings.Contains(adminMail.Body, fmt.Sprintf("/t/%d/events/%d/reservations", tid, ev)) {
		t.Fatalf("admin mail: %+v", adminMail)
	}

	requests := fmt.Sprintf("/api/tenants/%d/events/%d/booking-requests", tid, ev)
	code, b = e.do(staff, "GET", requests, nil)
	expect(t, "staff lists requests", code, 200, b)
	var list []store.BookingRequest
	json.Unmarshal(b, &list)
	if len(list) != 1 || list[0].Status != store.BookingPending || len(list[0].Shooters) != 1 {
		t.Fatalf("requests: %s", b)
	}
	reqPath := fmt.Sprintf("/api/tenants/%d/booking-requests/%d", tid, list[0].ID)
	code, b = e.do(staff, "POST", reqPath+"/confirm", nil)
	expect(t, "staff confirms", code, 403, b)

	other := e.tenant("Other club")
	e.user("stranger", other, store.RoleAdmin)
	code, b = e.do(e.login("stranger"), "POST", fmt.Sprintf("/api/tenants/%d/booking-requests/%d/confirm", other, list[0].ID), nil)
	expect(t, "other tenant confirms", code, 404, b)

	foreign, _ := e.st.EnsureClub(context.Background(), other, "Foreign")
	code, b = e.do(admin, "POST", reqPath+"/confirm", map[string]any{"clubId": foreign})
	expect(t, "confirm with a club of another tenant", code, 404, b)

	home, _ := e.st.EnsureClub(context.Background(), tid, "Home shooters")
	code, b = e.do(admin, "POST", reqPath+"/confirm", map[string]any{"clubId": home})
	expect(t, "confirm", code, 200, b)
	var confirmed store.BookingRequest
	json.Unmarshal(b, &confirmed)
	if confirmed.Status != store.BookingConfirmed || fmt.Sprint(confirmed.Items[0].Assigned) != "[1 2 3]" || confirmed.Shooters[0].ShooterID == nil ||
		confirmed.ClubID == nil || *confirmed.ClubID != home {
		t.Fatalf("confirmed: %s", b)
	}
	if sh, _ := e.st.GetShooter(context.Background(), tid, *confirmed.Shooters[0].ShooterID); len(sh.Clubs) != 1 || sh.Clubs[0].ID != home {
		t.Fatalf("shooter not in the chosen club: %+v", sh.Clubs)
	}
	if msg := box.mailTo(t, "ida@example.org"); !strings.Contains(msg.Body, "Lanes: 3 (1-3)") {
		t.Fatalf("confirmation mail: %+v", msg)
	}
	code, b = e.do(anon, "GET", public, nil)
	json.Unmarshal(b, &info)
	if info.Slots[0].Free[0].Free != 1 || info.Slots[1].Free[0].Free != 4 {
		t.Fatalf("free after confirmation: %s", b)
	}

	// The manage link: read, change shooters only (stays confirmed), change slots (pending again, lanes freed).
	code, b = e.helper("GET", "/api/booking-manage", "bk_nope", nil)
	expect(t, "unknown manage link", code, 401, b)
	code, b = e.helper("GET", "/api/booking-manage", secret, nil)
	expect(t, "manage", code, 200, b)
	if !strings.Contains(string(b), `"status":"confirmed"`) {
		t.Fatalf("manage: %s", b)
	}
	request.Shooters = append(request.Shooters, domain.BookingShooter{Firstname: "Max", Lastname: "Neu", Birthday: "2001-01-01"})
	code, b = e.helper("PUT", "/api/booking-manage", secret, request)
	expect(t, "add shooter", code, 200, b)
	if !strings.Contains(string(b), `"status":"confirmed"`) {
		t.Fatalf("shooter change reopened: %s", b)
	}
	box.mailTo(t, "boss@example.org")
	code, b = e.do(admin, "POST", reqPath+"/shooters", nil)
	expect(t, "adopt shooters", code, 200, b)
	if found, _ := e.st.SearchShooters(ctx, tid, "Max Neu", 0, 5, 0); len(found) != 1 {
		t.Fatalf("adopted shooter missing: %+v", found)
	}

	request.Items = []domain.BookingItem{{RangeID: rng.ID, Start: second.Start, End: second.End, Lanes: 4}}
	code, b = e.helper("PUT", "/api/booking-manage", secret, request)
	expect(t, "move to the second slot", code, 200, b)
	if !strings.Contains(string(b), `"status":"pending"`) {
		t.Fatalf("slot change kept confirmation: %s", b)
	}
	box.mailTo(t, "boss@example.org")
	code, b = e.do(anon, "GET", public, nil)
	json.Unmarshal(b, &info)
	if info.Slots[0].Free[0].Free != 4 {
		t.Fatalf("lanes not freed: %s", b)
	}

	// Lanes taken by hand in the meantime: confirming finds too few and says so.
	if _, err := e.st.SaveReservation(ctx, tid, store.Reservation{RangeID: rng.ID, Lanes: []int{2}, StartsAt: second.Start, EndsAt: second.End, Description: "Members"}); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(admin, "POST", reqPath+"/confirm", nil)
	expect(t, "confirm a full slot", code, 409, b)
	code, b = e.do(admin, "POST", reqPath+"/decline", map[string]string{"reason": "Fully booked"})
	expect(t, "decline", code, 200, b)
	if msg := box.mailTo(t, "ida@example.org"); !strings.Contains(msg.Body, "Fully booked") {
		t.Fatalf("decline mail: %+v", msg)
	}
	code, b = e.helper("PUT", "/api/booking-manage", secret, request)
	expect(t, "change a declined request", code, 409, b)
	code, b = e.helper("DELETE", "/api/booking-manage", secret, nil)
	expect(t, "cancel a declined request", code, 409, b)

	// After the deadline nothing can be asked or changed.
	yesterday := time.Now().AddDate(0, 0, -2).Format(time.DateOnly)
	code, b = e.do(admin, "PUT", base, domain.BookingSettings{From: "08:00", To: "10:00", SlotMinutes: 60, Deadline: yesterday})
	expect(t, "deadline in the past", code, 200, b)
	code, b = e.do(anon, "POST", public+"/requests", request)
	expect(t, "request after the deadline", code, 409, b)

	code, b = e.do(admin, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/reports/bookings?format=csv", tid, ev), nil)
	expect(t, "bookings report", code, 200, b)
	if strings.Contains(string(b), "SV Gäste") {
		t.Fatalf("declined request in the report: %s", b)
	}

	code, b = e.do(admin, "DELETE", base+"/code", nil)
	expect(t, "switch off", code, 204, b)
	code, b = e.do(anon, "GET", public, nil)
	expect(t, "switched off", code, http.StatusNotFound, b)
}

func TestOnlineBookingRateLimit(t *testing.T) {
	e, _ := personalEnv(t, nil)
	ctx := context.Background()
	tid := e.tenant("Host club")
	e.user("boss", tid, store.RoleAdmin)
	admin, anon := e.login("boss"), e.client()
	rng, _ := e.st.CreateRange(ctx, tid, store.RangeInput{Name: "10m", LaneCount: 10})
	day := time.Now().AddDate(0, 2, 0).Format(time.DateOnly)
	ev, _ := e.st.CreateEvent(ctx, tid, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "Lupi", StartsOn: day, Timezone: "Europe/Zurich"})
	code, b := e.do(admin, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/booking/code", tid, ev), nil)
	expect(t, "switch on", code, 200, b)
	var cfg web.BookingConfig
	json.Unmarshal(b, &cfg)
	in := store.BookingRequestInput{Club: "Spam", Contact: "Spam", Email: "spam@example.org",
		Items: []domain.BookingItem{{RangeID: rng.ID, Start: cfg.Slots[0].Start, End: cfg.Slots[0].End, Lanes: 1}}}
	for i := range 5 {
		code, b = e.do(anon, "POST", "/api/booking/"+cfg.Code+"/requests", in)
		expect(t, fmt.Sprintf("request %d", i+1), code, 204, b)
	}
	code, b = e.do(anon, "POST", "/api/booking/"+cfg.Code+"/requests", in)
	expect(t, "sixth request of the same email", code, 429, b)
}
