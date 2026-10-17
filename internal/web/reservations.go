package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// ReservationInput creates or changes a reservation; Lanes uses the "1, 5, 7-10" syntax.
type ReservationInput struct {
	Lanes       string    `json:"lanes" minLength:"1" doc:"e.g. 1, 5, 7-10"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
	LeadMinutes int       `json:"leadMinutes" minimum:"0" maximum:"1440" required:"false" doc:"Show the lanes as reserved this many minutes before the start"`
	Description string    `json:"description" minLength:"1" maxLength:"255"`
}

// ReservationView is a reservation with its lanes written compactly.
type ReservationView struct {
	store.Reservation
	LanesText string `json:"lanesText"`
}

// ReservationPath identifies a reservation of a tenant in the URL.
type ReservationPath struct {
	TenantID      int64 `path:"tenantId"`
	ReservationID int64 `path:"reservationId"`
}

// reservationView returns a reservation with its lanes as text.
func reservationView(r store.Reservation) ReservationView {
	return ReservationView{Reservation: r, LanesText: domain.FormatLanes(r.Lanes)}
}

// notifyRange announces a lane change to every event of the range.
func (s *Server) notifyRange(ctx context.Context, rangeID int64) {
	events, err := s.Store.EventsOfRange(ctx, rangeID)
	if err != nil {
		s.Log.Warn("events of range", "err", err)
		return
	}
	for _, ev := range events {
		s.notify(ctx, ev.TenantID, ev.EventID, realtime.TopicLanes)
	}
}

// saveReservation validates lanes and times, saves the reservation and names the overlapping ones on conflict.
func (s *Server) saveReservation(ctx context.Context, tenantID int64, r store.Reservation, in ReservationInput, isNew bool) (store.Reservation, error) {
	lanes, err := domain.ParseLanes(in.Lanes)
	if err != nil {
		return r, invalid(err.Error())
	}
	ranges, err := s.Store.ListRanges(ctx, tenantID)
	if err != nil {
		return r, err
	}
	laneCount := 0
	for _, rg := range ranges {
		if rg.ID == r.RangeID {
			laneCount = rg.LaneCount
		}
	}
	if laneCount == 0 {
		return r, missing("Range not found")
	}
	if lanes[len(lanes)-1] > laneCount {
		return r, invalid("The range has lanes 1–%d", laneCount)
	}
	if !in.EndsAt.After(in.StartsAt) {
		return r, invalid("The end must be after the start")
	}
	if isNew && !in.EndsAt.After(time.Now()) {
		return r, invalid("The reservation ends in the past")
	}
	r.Lanes, r.StartsAt, r.EndsAt, r.LeadMinutes, r.Description = lanes, in.StartsAt, in.EndsAt, in.LeadMinutes, strings.TrimSpace(in.Description)
	id, err := s.Store.SaveReservation(ctx, tenantID, r)
	if errors.Is(err, store.ErrOverlap) {
		others, oerr := s.Store.Overlaps(ctx, r)
		if oerr != nil {
			s.Log.Warn("reservation overlaps", "err", oerr)
			return r, conflict("Overlaps with another reservation")
		}
		var names []string
		for _, o := range others {
			names = append(names, fmt.Sprintf("%q (lanes %s)", o.Description, domain.FormatLanes(o.Lanes)))
		}
		return r, conflict("Overlaps with %s", strings.Join(names, ", "))
	}
	if err != nil {
		return r, notFound(err, "Reservation not found")
	}
	r.ID = id
	s.notifyRange(ctx, r.RangeID)
	return r, nil
}

// registerReservations registers the lane reservation endpoints.
func (s *Server) registerReservations(api huma.API) {
	type listIn struct {
		RangePrintPath
		All bool `query:"all" doc:"Include past reservations"`
	}
	type listOut struct{ Body []ReservationView }
	huma.Register(api, huma.Operation{
		OperationID: "list-reservations", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/reservations", Tags: []string{"reservations"},
	}, func(ctx context.Context, in *listIn) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListReservations(ctx, in.TenantID, in.RangeID, in.All)
		out := &listOut{Body: []ReservationView{}}
		for _, r := range list {
			out.Body = append(out.Body, reservationView(r))
		}
		return out, err
	})

	type viewOut struct{ Body ReservationView }
	type createIn struct {
		RangePrintPath
		Body ReservationInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-reservation", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/reservations", Tags: []string{"reservations"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*viewOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		r, err := s.saveReservation(ctx, in.TenantID, store.Reservation{RangeID: in.RangeID}, in.Body, true)
		return &viewOut{Body: reservationView(r)}, err
	})

	type updateIn struct {
		ReservationPath
		Body ReservationInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-reservation", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/reservations/{reservationId}", Tags: []string{"reservations"},
	}, func(ctx context.Context, in *updateIn) (*viewOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		current, err := s.Store.GetReservation(ctx, in.TenantID, in.ReservationID)
		if err != nil {
			return nil, notFound(err, "Reservation not found")
		}
		r, err := s.saveReservation(ctx, in.TenantID, current, in.Body, false)
		if err == nil && current.RangeID != r.RangeID {
			s.notifyRange(ctx, current.RangeID)
		}
		return &viewOut{Body: reservationView(r)}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-reservation", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/reservations/{reservationId}", Tags: []string{"reservations"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ReservationPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		rangeID, err := s.Store.DeleteReservation(ctx, in.TenantID, in.ReservationID)
		if err != nil {
			return nil, notFound(err, "Reservation not found")
		}
		s.notifyRange(ctx, rangeID)
		return nil, nil
	})
}

// RunReservationTicker announces lane changes caused by reservations starting or ending, every 30 seconds.
func (s *Server) RunReservationTicker(ctx context.Context) {
	tick := time.NewTicker(30 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-tick.C:
			ranges, err := s.Store.ReservationTick(ctx, now)
			if err != nil {
				s.Log.Warn("reservation tick", "err", err)
				continue
			}
			for _, r := range ranges {
				s.notifyRange(ctx, r)
			}
		}
	}
}
