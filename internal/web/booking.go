package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

const (
	// bookingPrefix marks the secret of a club's manage link (bk_…).
	bookingPrefix   = "bk_"
	bookingWindow   = time.Hour
	bookingIPMax    = 20
	bookingEmailMax = 5
)

// BookingRange is a range clubs can book lanes on.
type BookingRange struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	LaneCount int    `json:"laneCount" doc:"Lanes in service"`
}

// SlotFree is how many lanes of a range are still free in a slot.
type SlotFree struct {
	RangeID int64 `json:"rangeId"`
	Free    int   `json:"free"`
}

// BookingSlot is a slot with the free lanes per range.
type BookingSlot struct {
	domain.Slot
	Free []SlotFree `json:"free"`
}

// BookingInfo is what the public booking page shows: the event, its ranges and slots and whether it is still open.
type BookingInfo struct {
	Event     string         `json:"event"`
	Organizer string         `json:"organizer"`
	Days      []string       `json:"days"`
	Timezone  string         `json:"timezone"`
	Locale    string         `json:"locale"`
	Ranges    []BookingRange `json:"ranges"`
	Slots     []BookingSlot  `json:"slots"`
	Deadline  time.Time      `json:"deadline" doc:"Requests and changes are accepted before this moment"`
	Open      bool           `json:"open"`
	MaxLanes  int            `json:"maxLanes" doc:"Lanes per slot and request; 0 = no limit"`
}

// BookingConfig is the online booking of an event as the admin sets it up.
type BookingConfig struct {
	Enabled  bool                   `json:"enabled"`
	Code     string                 `json:"code" doc:"Empty when off"`
	Link     string                 `json:"link" doc:"The public link for the invitation; empty when off"`
	QR       string                 `json:"qr" doc:"PNG data URL of a QR code of the link; empty when off"`
	Settings domain.BookingSettings `json:"settings"`
	Deadline time.Time              `json:"deadline" doc:"The effective deadline"`
	Slots    []domain.Slot          `json:"slots" doc:"The slots the settings give"`
}

// BookingManageAuth carries the secret of a club's manage link.
type BookingManageAuth struct {
	Authorization string `header:"Authorization" doc:"Bearer bk_…"`
}

// BookingCodePath carries the code of a public booking link.
type BookingCodePath struct {
	Code string `path:"code" maxLength:"64"`
}

// BookingRequestPath identifies a booking request of a tenant.
type BookingRequestPath struct {
	TenantID  int64 `path:"tenantId"`
	RequestID int64 `path:"requestId"`
}

// bookingLoc returns the time zone of an event, UTC when it is unknown.
func bookingLoc(ev store.Event) *time.Location {
	if loc, err := time.LoadLocation(ev.Timezone); err == nil {
		return loc
	}
	return time.UTC
}

// bookingRanges returns the ranges of an event with their lanes in service.
func bookingRanges(ev store.Event) []BookingRange {
	out := []BookingRange{}
	for _, r := range ev.Ranges {
		out = append(out, BookingRange{ID: r.ID, Name: r.Name, LaneCount: r.LaneCount - len(r.DisabledLanes)})
	}
	return out
}

// bookingInfo builds the public view of an event's booking: slots with the lanes still free per range, counting
// only confirmed reservations; those of exceptRequest are left out because changing it would free them.
func (s *Server) bookingInfo(ctx context.Context, t store.BookingTarget, exceptRequest int64) (BookingInfo, error) {
	ev := t.Event
	loc := bookingLoc(ev)
	slots := domain.BookingSlots(ev.Days, ev.Schedule.Times(), t.Settings, loc)
	deadline := domain.BookingDeadline(ev.Days, t.Settings, loc)
	info := BookingInfo{
		Event: ev.Name, Organizer: ev.Organizer, Days: ev.Days, Timezone: ev.Timezone, Locale: ev.Locale,
		Ranges: bookingRanges(ev), Slots: []BookingSlot{}, Deadline: deadline, Open: time.Now().Before(deadline), MaxLanes: t.Settings.MaxLanes,
	}
	if len(slots) == 0 {
		return info, nil
	}
	rangeIDs := []int64{}
	for _, r := range ev.Ranges {
		rangeIDs = append(rangeIDs, r.ID)
	}
	res, err := s.Store.ReservationsBetween(ctx, t.TenantID, rangeIDs, slots[0].Start, slots[len(slots)-1].End, exceptRequest)
	if err != nil {
		return info, err
	}
	for _, sl := range slots {
		bs := BookingSlot{Slot: sl, Free: []SlotFree{}}
		for _, r := range ev.Ranges {
			busy := []int{}
			for _, x := range res {
				if x.RangeID == r.ID && x.StartsAt.Before(sl.End) && x.EndsAt.After(sl.Start) {
					busy = append(busy, x.Lanes...)
				}
			}
			bs.Free = append(bs.Free, SlotFree{RangeID: r.ID, Free: len(domain.FreeLanes(r.LaneCount, r.DisabledLanes, busy))})
		}
		info.Slots = append(info.Slots, bs)
	}
	return info, nil
}

// checkBooking validates a request against the event's slots and the lanes free now (info leaves out the lanes the
// request already holds).
func checkBooking(info BookingInfo, in store.BookingRequestInput) error {
	if !info.Open {
		return conflict("The deadline for online booking has passed. Contact the organizer.")
	}
	counts := map[int64]int{}
	for _, r := range info.Ranges {
		counts[r.ID] = r.LaneCount
	}
	slots := make([]domain.Slot, len(info.Slots))
	for i, sl := range info.Slots {
		slots[i] = sl.Slot
	}
	if errs := domain.ValidateBooking(in.Items, in.Shooters, slots, counts, info.MaxLanes); len(errs) > 0 {
		return invalidList(errs)
	}
	for _, it := range in.Items {
		for _, sl := range info.Slots {
			if !sl.Start.Equal(it.Start) {
				continue
			}
			for _, f := range sl.Free {
				if f.RangeID == it.RangeID && it.Lanes > f.Free {
					loc, err := time.LoadLocation(info.Timezone)
					if err != nil {
						loc = time.UTC
					}
					return invalid("Only %d lanes are free at %s", f.Free, sl.Start.In(loc).Format("02.01.2006 15:04"))
				}
			}
		}
	}
	return nil
}

// trimBooking trims the text fields of a request and lower-cases the email.
func trimBooking(in *store.BookingRequestInput) {
	in.Club, in.Contact, in.Phone, in.Note = strings.TrimSpace(in.Club), strings.TrimSpace(in.Contact), strings.TrimSpace(in.Phone), strings.TrimSpace(in.Note)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if in.Shooters == nil {
		in.Shooters = []domain.BookingShooter{}
	}
}

// bookingTarget loads the event behind a public code, 404 when the code is unknown, switched off or the event closed.
func (s *Server) bookingTarget(ctx context.Context, code string) (store.BookingTarget, error) {
	t, err := s.Store.BookingByCode(ctx, code)
	return t, notFound(err, "This booking link is not valid.")
}

// bookingManaged loads the request behind a manage link and its event.
func (s *Server) bookingManaged(ctx context.Context, header string) (store.BookingTarget, store.BookingRequest, error) {
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || !strings.HasPrefix(token, bookingPrefix) {
		return store.BookingTarget{}, store.BookingRequest{}, unauthorized("This link is not valid.")
	}
	tenantID, req, err := s.Store.BookingRequestByHash(ctx, sessionID(token))
	if errors.Is(err, store.ErrNotFound) {
		return store.BookingTarget{}, req, unauthorized("This link is not valid.")
	}
	if err != nil {
		return store.BookingTarget{}, req, err
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, req.EventID)
	if err != nil {
		return store.BookingTarget{}, req, err
	}
	_, set, err := s.Store.BookingConfig(ctx, tenantID, req.EventID)
	return store.BookingTarget{TenantID: tenantID, Event: ev, Settings: set}, req, err
}

// slotText formats a slot for mails in the event's time zone, e.g. "13.03.2027 08:00–08:30".
func slotText(start, end time.Time, loc *time.Location) string {
	return start.In(loc).Format("02.01.2006 15:04") + "–" + end.In(loc).Format("15:04")
}

// bookingSummary lists the slots of a request for a mail, with the assigned lanes once confirmed.
func bookingSummary(l i18n.Lang, ev store.Event, req store.BookingRequest) []string {
	names := map[int64]string{}
	for _, r := range ev.Ranges {
		names[r.ID] = r.Name
	}
	loc := bookingLoc(ev)
	lines := []string{}
	for _, it := range req.Items {
		line := "- " + slotText(it.Start, it.End, loc) + ", " + names[it.RangeID] + ", " + i18n.T(l, "Lanes") + ": " + strconv.Itoa(it.Lanes)
		if len(it.Assigned) > 0 {
			line += " (" + domain.FormatLanes(it.Assigned) + ")"
		}
		lines = append(lines, line)
	}
	if n := len(req.Shooters); n > 0 {
		lines = append(lines, i18n.Tf(l, "Shooters announced: %d", n))
	}
	return lines
}

// mailClub sends a mail about a request to the club's contact in the event's language.
func (s *Server) mailClub(ctx context.Context, ev store.Event, req store.BookingRequest, subject string, intro []string, manageLink string) {
	l := i18n.Get(ev.Locale)
	body := append([]string{greeting(l, req.Contact), ""}, intro...)
	body = append(body, "")
	body = append(body, bookingSummary(l, ev, req)...)
	if manageLink != "" {
		body = append(body, "", i18n.T(l, "Check, change or cancel your request here:"), manageLink)
	}
	s.mail(mail.Message{To: req.Email, Subject: subject, Body: strings.Join(body, "\n")})
}

// mailAdmins tells every admin of the club with an email about a request that needs a look.
func (s *Server) mailAdmins(ctx context.Context, tenantID int64, ev store.Event, req store.BookingRequest, subject string) {
	members, err := s.Store.Members(ctx, tenantID)
	if err != nil {
		s.Log.Warn("booking admins", "err", err)
		return
	}
	l := i18n.Get(ev.Locale)
	link := s.link(ctx, fmt.Sprintf("/t/%d/events/%d/reservations", tenantID, ev.ID))
	for _, m := range members {
		if m.Role != store.RoleAdmin || m.Email == "" {
			continue
		}
		body := append([]string{greeting(l, ""), "", i18n.Tf(l, "%s (%s, %s) asks for lanes at %s:", req.Club, req.Contact, req.Email, ev.Name), ""},
			bookingSummary(l, ev, req)...)
		body = append(body, "", i18n.T(l, "Confirm or decline the request here:"), link)
		s.mail(mail.Message{To: m.Email, Subject: subject, Body: strings.Join(body, "\n")})
	}
}

// notifyBooking tells the event's pages that requests changed, and the lane grids when reservations did.
func (s *Server) notifyBooking(ctx context.Context, tenantID int64, ev store.Event, lanes bool) {
	s.notify(ctx, tenantID, ev.ID, realtime.TopicBooking)
	if lanes {
		for _, r := range ev.Ranges {
			s.notifyRange(ctx, r.ID)
		}
	}
}

// bookingErr maps the store errors of a request change to API errors.
func bookingErr(ev store.Event, err error) error {
	var full *store.FullSlotError
	switch {
	case errors.As(err, &full):
		return conflict("Only %d lanes are free at %s; move a reservation or decline the request.", full.Free,
			full.Item.Start.In(bookingLoc(ev)).Format("02.01.2006 15:04"))
	case errors.Is(err, store.ErrNotPending):
		return conflict("The request was already handled or cancelled.")
	}
	return notFound(err, "Request not found")
}

// registerBooking registers the public booking endpoints, the manage link and the admin endpoints.
func (s *Server) registerBooking(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-booking", Method: http.MethodGet, Path: "/api/booking/{code}", Tags: []string{"booking"},
		Summary: "The event, slots and free lanes behind a public booking link",
	}, func(ctx context.Context, in *BookingCodePath) (*struct{ Body BookingInfo }, error) {
		t, err := s.bookingTarget(ctx, in.Code)
		if err != nil {
			return nil, err
		}
		info, err := s.bookingInfo(ctx, t, 0)
		return &struct{ Body BookingInfo }{Body: info}, err
	})

	type createIn struct {
		BookingCodePath
		Body store.BookingRequestInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-booking-request", Method: http.MethodPost, Path: "/api/booking/{code}/requests", Tags: []string{"booking"},
		Summary: "Ask for lanes; the club gets a mail with a link to change or cancel the request", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *createIn) (*struct{}, error) {
		t, err := s.bookingTarget(ctx, in.Code)
		if err != nil {
			return nil, err
		}
		trimBooking(&in.Body)
		if err := s.rateLimited(ctx, "booking:ip:"+reqInfo(ctx).ip, bookingWindow, bookingIPMax); err != nil {
			return nil, err
		}
		if err := s.rateLimited(ctx, "booking:email:"+in.Body.Email, bookingWindow, bookingEmailMax); err != nil {
			return nil, err
		}
		info, err := s.bookingInfo(ctx, t, 0)
		if err != nil {
			return nil, err
		}
		if err := checkBooking(info, in.Body); err != nil {
			return nil, err
		}
		secret := bookingPrefix + randomToken()
		id, err := s.Store.CreateBookingRequest(ctx, t.TenantID, t.Event.ID, in.Body, sessionID(secret))
		if err != nil {
			return nil, err
		}
		req, err := s.Store.GetBookingRequest(ctx, t.TenantID, id)
		if err != nil {
			return nil, err
		}
		l := i18n.Get(t.Event.Locale)
		s.mailClub(ctx, t.Event, req, i18n.Tf(l, "Your lane request for %s", t.Event.Name),
			[]string{i18n.Tf(l, "we received your request for %s. The organizer confirms it by mail; until then the lanes are not reserved.", t.Event.Name)},
			s.link(ctx, "/b/manage/"+secret))
		s.mailAdmins(ctx, t.TenantID, t.Event, req, i18n.Tf(l, "New lane request from %s", req.Club))
		s.notifyBooking(ctx, t.TenantID, t.Event, false)
		return nil, nil
	})

	type manageOut struct {
		Body struct {
			Request store.BookingRequest `json:"request"`
			Info    BookingInfo          `json:"info"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "get-my-booking", Method: http.MethodGet, Path: "/api/booking-manage", Tags: []string{"booking"},
		Summary: "The request behind a manage link (Bearer bk_…) with the event's slots",
	}, func(ctx context.Context, in *BookingManageAuth) (*manageOut, error) {
		t, req, err := s.bookingManaged(ctx, in.Authorization)
		if err != nil {
			return nil, err
		}
		out := &manageOut{}
		out.Body.Request = req
		out.Body.Info, err = s.bookingInfo(ctx, t, req.ID)
		return out, err
	})

	type updateIn struct {
		BookingManageAuth
		Body store.BookingRequestInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-my-booking", Method: http.MethodPut, Path: "/api/booking-manage", Tags: []string{"booking"},
		Summary: "Change a request before the deadline; changed slots make a confirmed request pending again",
	}, func(ctx context.Context, in *updateIn) (*struct{ Body store.BookingRequest }, error) {
		t, req, err := s.bookingManaged(ctx, in.Authorization)
		if err != nil {
			return nil, err
		}
		if req.Status != store.BookingPending && req.Status != store.BookingConfirmed {
			return nil, bookingErr(t.Event, store.ErrNotPending)
		}
		trimBooking(&in.Body)
		info, err := s.bookingInfo(ctx, t, req.ID)
		if err != nil {
			return nil, err
		}
		if err := checkBooking(info, in.Body); err != nil {
			return nil, err
		}
		reopened, err := s.Store.UpdateBookingRequest(ctx, t.TenantID, req.ID, in.Body)
		if err != nil {
			return nil, bookingErr(t.Event, err)
		}
		updated, err := s.Store.GetBookingRequest(ctx, t.TenantID, req.ID)
		if err != nil {
			return nil, err
		}
		l := i18n.Get(t.Event.Locale)
		s.mailAdmins(ctx, t.TenantID, t.Event, updated, i18n.Tf(l, "%s changed the lane request", updated.Club))
		s.notifyBooking(ctx, t.TenantID, t.Event, reopened)
		return &struct{ Body store.BookingRequest }{Body: updated}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "cancel-my-booking", Method: http.MethodDelete, Path: "/api/booking-manage", Tags: []string{"booking"},
		Summary: "Cancel a request before the deadline; its lanes are free again", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *BookingManageAuth) (*struct{}, error) {
		t, req, err := s.bookingManaged(ctx, in.Authorization)
		if err != nil {
			return nil, err
		}
		if !time.Now().Before(domain.BookingDeadline(t.Event.Days, t.Settings, bookingLoc(t.Event))) {
			return nil, conflict("The deadline for online booking has passed. Contact the organizer.")
		}
		if err := s.Store.CancelBookingRequest(ctx, t.TenantID, req.ID); err != nil {
			return nil, bookingErr(t.Event, err)
		}
		l := i18n.Get(t.Event.Locale)
		s.mailAdmins(ctx, t.TenantID, t.Event, req, i18n.Tf(l, "%s cancelled the lane request", req.Club))
		s.notifyBooking(ctx, t.TenantID, t.Event, req.Status == store.BookingConfirmed)
		return nil, nil
	})

	s.registerBookingAdmin(api)
}

// registerBookingAdmin registers the endpoints with which a club sets up its online booking and handles requests.
func (s *Server) registerBookingAdmin(api huma.API) {
	config := func(ctx context.Context, tenantID, eventID int64) (*struct{ Body BookingConfig }, error) {
		ev, err := s.Store.GetEvent(ctx, tenantID, eventID)
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		code, set, err := s.Store.BookingConfig(ctx, tenantID, eventID)
		if err != nil {
			return nil, err
		}
		loc := bookingLoc(ev)
		c := BookingConfig{Enabled: code != "", Code: code, Settings: set,
			Deadline: domain.BookingDeadline(ev.Days, set, loc), Slots: domain.BookingSlots(ev.Days, ev.Schedule.Times(), set, loc)}
		if code != "" {
			c.Link = s.link(ctx, "/b/"+code)
			if c.QR, err = totpQR(c.Link); err != nil {
				return nil, err
			}
		}
		return &struct{ Body BookingConfig }{Body: c}, nil
	}

	huma.Register(api, huma.Operation{
		OperationID: "get-booking-config", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/booking", Tags: []string{"booking"},
	}, func(ctx context.Context, in *EventPath) (*struct{ Body BookingConfig }, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		return config(ctx, in.TenantID, in.EventID)
	})

	type settingsIn struct {
		EventPath
		Body domain.BookingSettings
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-booking-settings", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/events/{eventId}/booking", Tags: []string{"booking"},
	}, func(ctx context.Context, in *settingsIn) (*struct{ Body BookingConfig }, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		ev, err := s.Store.GetEvent(ctx, in.TenantID, in.EventID)
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		if errs := domain.ValidateBookingSettings(in.Body, ev.Days); len(errs) > 0 {
			return nil, invalidList(errs)
		}
		if err := s.Store.SetBookingSettings(ctx, in.TenantID, in.EventID, domain.CleanBookingSettings(in.Body, ev.Days)); err != nil {
			return nil, notFound(err, "Event not found")
		}
		s.notify(ctx, in.TenantID, in.EventID, realtime.TopicBooking)
		return config(ctx, in.TenantID, in.EventID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "new-booking-code", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/booking/code", Tags: []string{"booking"},
		Summary: "Switch the online booking on with a new link; an earlier link stops working",
	}, func(ctx context.Context, in *EventPath) (*struct{ Body BookingConfig }, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := s.Store.SetBookingCode(ctx, in.TenantID, in.EventID, randomToken()[:16]); err != nil {
			return nil, notFound(err, "Event not found")
		}
		s.notify(ctx, in.TenantID, in.EventID, realtime.TopicBooking)
		return config(ctx, in.TenantID, in.EventID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "disable-booking", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/events/{eventId}/booking/code", Tags: []string{"booking"},
		Summary: "Switch the online booking off; requests and reservations stay", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *EventPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := s.Store.SetBookingCode(ctx, in.TenantID, in.EventID, ""); err != nil {
			return nil, notFound(err, "Event not found")
		}
		s.notify(ctx, in.TenantID, in.EventID, realtime.TopicBooking)
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-booking-requests", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/booking-requests", Tags: []string{"booking"},
	}, func(ctx context.Context, in *EventPath) (*struct{ Body []store.BookingRequest }, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListBookingRequests(ctx, in.TenantID, in.EventID)
		return &struct{ Body []store.BookingRequest }{Body: list}, err
	})

	decide := func(op, suffix, summary string, run func(ctx context.Context, tenantID int64, req store.BookingRequest, ev store.Event, reason string) error) {
		type decideIn struct {
			BookingRequestPath
			Body *struct {
				Reason string `json:"reason,omitempty" required:"false" maxLength:"1000" doc:"Declining: told to the club"`
			}
		}
		huma.Register(api, huma.Operation{
			OperationID: op, Method: http.MethodPost, Path: "/api/tenants/{tenantId}/booking-requests/{requestId}/" + suffix, Tags: []string{"booking"},
			Summary: summary,
		}, func(ctx context.Context, in *decideIn) (*struct{ Body store.BookingRequest }, error) {
			if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
				return nil, err
			}
			req, err := s.Store.GetBookingRequest(ctx, in.TenantID, in.RequestID)
			if err != nil {
				return nil, notFound(err, "Request not found")
			}
			ev, err := s.Store.GetEvent(ctx, in.TenantID, req.EventID)
			if err != nil {
				return nil, err
			}
			reason := ""
			if in.Body != nil {
				reason = strings.TrimSpace(in.Body.Reason)
			}
			if err := run(ctx, in.TenantID, req, ev, reason); err != nil {
				return nil, bookingErr(ev, err)
			}
			req, err = s.Store.GetBookingRequest(ctx, in.TenantID, in.RequestID)
			return &struct{ Body store.BookingRequest }{Body: req}, err
		})
	}

	decide("confirm-booking-request", "confirm", "Reserve free lanes for every slot and take the shooters into the master data",
		func(ctx context.Context, tenantID int64, req store.BookingRequest, ev store.Event, _ string) error {
			if err := s.Store.ConfirmBooking(ctx, tenantID, req.ID); err != nil {
				return err
			}
			req, err := s.Store.GetBookingRequest(ctx, tenantID, req.ID)
			if err != nil {
				return err
			}
			l := i18n.Get(ev.Locale)
			s.mailClub(ctx, ev, req, i18n.Tf(l, "Lanes confirmed for %s", ev.Name),
				[]string{i18n.Tf(l, "your lanes for %s are reserved:", ev.Name)}, "")
			s.notifyBooking(ctx, tenantID, ev, true)
			return nil
		})

	decide("decline-booking-request", "decline", "Decline a request; lanes it held are free again",
		func(ctx context.Context, tenantID int64, req store.BookingRequest, ev store.Event, reason string) error {
			if err := s.Store.DeclineBookingRequest(ctx, tenantID, req.ID, reason); err != nil {
				return err
			}
			l := i18n.Get(ev.Locale)
			intro := []string{i18n.Tf(l, "unfortunately the organizer of %s cannot give you these lanes.", ev.Name)}
			if reason != "" {
				intro = append(intro, "", reason)
			}
			s.mailClub(ctx, ev, req, i18n.Tf(l, "Lane request for %s declined", ev.Name), intro, "")
			s.notifyBooking(ctx, tenantID, ev, req.Status == store.BookingConfirmed)
			return nil
		})

	decide("adopt-booking-shooters", "shooters", "Take shooters added after the confirmation into the master data",
		func(ctx context.Context, tenantID int64, req store.BookingRequest, ev store.Event, _ string) error {
			if err := s.Store.AdoptBookingShooters(ctx, tenantID, req.ID); err != nil {
				return err
			}
			s.notify(ctx, tenantID, ev.ID, realtime.TopicBooking)
			return nil
		})
}
