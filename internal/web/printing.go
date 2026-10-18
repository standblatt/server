package web

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/docs"
	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// AgentWait is how long GET /api/agent/jobs/next waits for a job before it answers 204.
var AgentWait = 25 * time.Second

// AgentAuth carries the print agent token of a request.
type AgentAuth struct {
	Authorization string `header:"Authorization" doc:"Bearer <agent token>"`
}

// agentAuth returns the print agent of a Bearer token, 401 when it is missing or unknown.
func (s *Server) agentAuth(ctx context.Context, header string) (store.AgentRef, error) {
	token := strings.TrimPrefix(header, "Bearer ")
	if token == "" || token == header {
		return store.AgentRef{}, unauthorized("Agent token missing")
	}
	a, err := s.Store.AgentByToken(ctx, sessionID(token))
	if errors.Is(err, store.ErrNotFound) {
		return a, unauthorized("Unknown or revoked agent token")
	}
	return a, err
}

// RangePrintPath identifies a range of a tenant for print management.
type RangePrintPath struct {
	TenantID int64 `path:"tenantId"`
	RangeID  int64 `path:"rangeId"`
}

// registerPrinting registers the print agent endpoints and the print agent, setting and job management endpoints.
func (s *Server) registerPrinting(api huma.API) {
	type helloIn struct {
		AgentAuth
		Body struct {
			Printers []store.Printer `json:"printers"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "agent-hello", Method: http.MethodPost, Path: "/api/agent/hello", Tags: []string{"agent"},
		Summary: "Report the printers of the agent", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *helloIn) (*struct{}, error) {
		a, err := s.agentAuth(ctx, in.Authorization)
		if err != nil {
			return nil, err
		}
		return nil, s.Store.SetAgentPrinters(ctx, a.ID, in.Body.Printers)
	})

	type jobOut struct {
		Status int
		Body   *store.ClaimedJob
	}
	huma.Register(api, huma.Operation{
		OperationID: "agent-next-job", Method: http.MethodGet, Path: "/api/agent/jobs/next", Tags: []string{"agent"},
		Summary: "Wait up to 25 s for the next job of the agent's range; 204 when there is none",
	}, func(ctx context.Context, in *AgentAuth) (*jobOut, error) {
		a, err := s.agentAuth(ctx, in.Authorization)
		if err != nil {
			return nil, err
		}
		deadline := time.After(AgentWait)
		// ponytail: polls the queue once a second; a NOTIFY wake-up would cut latency if that ever matters
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			j, err := s.Store.ClaimPrintJob(ctx, a)
			if err == nil {
				return &jobOut{Status: http.StatusOK, Body: &j}, nil
			}
			if !errors.Is(err, store.ErrNotFound) {
				return nil, err
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-deadline:
				return &jobOut{Status: http.StatusNoContent}, nil
			case <-tick.C:
			}
		}
	})

	type doneIn struct {
		AgentAuth
		JobID int64 `path:"jobId"`
		Body  struct {
			OK    bool   `json:"ok"`
			Error string `json:"error,omitempty" required:"false"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "agent-finish-job", Method: http.MethodPost, Path: "/api/agent/jobs/{jobId}", Tags: []string{"agent"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *doneIn) (*struct{}, error) {
		a, err := s.agentAuth(ctx, in.Authorization)
		if err != nil {
			return nil, err
		}
		eventID, err := s.Store.FinishPrintJob(ctx, a, in.JobID, in.Body.OK, in.Body.Error)
		if err != nil {
			return nil, notFound(err, "Job not found or not being printed")
		}
		if eventID != nil {
			s.notify(ctx, a.TenantID, *eventID, realtime.TopicPrint)
		}
		return nil, nil
	})

	type agentsOut struct{ Body []store.PrintAgent }
	huma.Register(api, huma.Operation{
		OperationID: "list-print-agents", Metadata: needs("", permissions.RangesManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-agents", Tags: []string{"printing"},
	}, func(ctx context.Context, in *RangePrintPath) (*agentsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListPrintAgents(ctx, in.TenantID, in.RangeID)
		return &agentsOut{Body: list}, err
	})

	type createAgentIn struct {
		RangePrintPath
		Body struct {
			Name string `json:"name" minLength:"1" maxLength:"100"`
		}
	}
	type createAgentOut struct {
		Body struct {
			Agent  store.PrintAgent `json:"agent"`
			Secret string           `json:"secret" doc:"Shown once; set it as STANDBLATT_AGENT_TOKEN on the agent"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-print-agent", Metadata: needs("", permissions.RangesManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-agents", Tags: []string{"printing"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createAgentIn) (*createAgentOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		secret := randomToken()
		name := strings.TrimSpace(in.Body.Name)
		id, err := s.Store.CreatePrintAgent(ctx, in.TenantID, in.RangeID, name, sessionID(secret))
		if err != nil {
			return nil, notFound(err, "Range not found")
		}
		out := &createAgentOut{}
		out.Body.Agent = store.PrintAgent{ID: id, Name: name, Printers: []store.Printer{}}
		out.Body.Secret = secret
		return out, nil
	})

	type agentPath struct {
		TenantID int64 `path:"tenantId"`
		AgentID  int64 `path:"agentId"`
	}
	huma.Register(api, huma.Operation{
		OperationID: "delete-print-agent", Metadata: needs("", permissions.RangesManage), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/print-agents/{agentId}", Tags: []string{"printing"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *agentPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeletePrintAgent(ctx, in.TenantID, in.AgentID), "Print agent not found")
	})

	type settingsOut struct{ Body []store.PrintSetting }
	huma.Register(api, huma.Operation{
		OperationID: "list-print-settings", Metadata: needs("", permissions.EventsOperate, permissions.RangesManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-settings", Tags: []string{"printing"},
	}, func(ctx context.Context, in *RangePrintPath) (*settingsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListPrintSettings(ctx, in.TenantID, in.RangeID)
		return &settingsOut{Body: list}, err
	})

	type settingIn struct {
		RangePrintPath
		Body store.PrintSetting
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-print-setting", Metadata: needs("", permissions.RangesManage), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-settings", Tags: []string{"printing"},
		Summary: "Route a document to a printer; an empty printer turns printing off for it", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *settingIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.SetPrintSetting(ctx, in.TenantID, in.RangeID, in.Body), "Range not found")
	})

	type testIn struct {
		RangePrintPath
		Body struct {
			Printer string            `json:"printer" minLength:"1"`
			Options map[string]string `json:"options,omitempty" required:"false"`
		}
	}
	type testOut struct{ Body PrintResult }
	huma.Register(api, huma.Operation{
		OperationID: "print-test-page", Metadata: needs("", permissions.RangesManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-test", Tags: []string{"printing"},
		Summary: "Send a test page to a printer of the range with the given options; 409 when no agent of the range is online",
	}, func(ctx context.Context, in *testIn) (*testOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		pdf, title, err := s.renderTestPage(ctx, in.TenantID, in.RangeID, in.Body.Printer, in.Body.Options)
		if err != nil {
			return nil, err
		}
		id, err := s.Store.EnqueueTestPage(ctx, in.TenantID, in.RangeID, title, pdf, in.Body.Printer, in.Body.Options)
		if errors.Is(err, store.ErrNoPrinter) {
			return nil, conflict("No print agent of this range is online")
		}
		if err != nil {
			return nil, notFound(err, "Range not found")
		}
		return &testOut{Body: PrintResult{Queued: true, JobID: id}}, nil
	})

	type jobsOut struct{ Body []store.PrintJob }
	huma.Register(api, huma.Operation{
		OperationID: "list-print-jobs", Metadata: needs("", permissions.EventsOperate, permissions.RangesManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/ranges/{rangeId}/print-jobs", Tags: []string{"printing"},
	}, func(ctx context.Context, in *RangePrintPath) (*jobsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListPrintJobs(ctx, in.TenantID, in.RangeID, 50)
		return &jobsOut{Body: list}, err
	})
}

// renderTestPage renders the test page for a printer of a range in the language of the request and returns it with
// its job title.
func (s *Server) renderTestPage(ctx context.Context, tenantID, rangeID int64, printer string, options map[string]string) ([]byte, string, error) {
	t, err := s.Store.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, "", err
	}
	ranges, err := s.Store.ListRanges(ctx, tenantID)
	if err != nil {
		return nil, "", err
	}
	i := slices.IndexFunc(ranges, func(r store.Range) bool { return r.ID == rangeID })
	if i < 0 {
		return nil, "", notFound(store.ErrNotFound, "Range not found")
	}
	l := i18n.Get(reqInfo(ctx).lang)
	pdf, err := docs.RenderTestPage(docs.TestPage{Club: t.Name, Range: ranges[i].Name, Printer: printer, Options: options,
		Time: time.Now().In(zone("")).Format("02.01.2006 15:04"), Lang: l})
	return pdf, i18n.Tf(l, "Test page %s", printer), err
}
