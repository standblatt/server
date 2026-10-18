package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

// agentCall sends a request of the print agent protocol (github.com/standblatt/print-agent) with token and returns
// the status code and body.
func agentCall(t *testing.T, url, token, method, path string, body any) (int, []byte) {
	t.Helper()
	var r io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		r = bytes.NewReader(raw)
	}
	req, _ := http.NewRequest(method, url+path, r)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, out
}

// printJob is a job as the server hands it to an agent.
type printJob struct {
	ID      int64             `json:"id"`
	Title   string            `json:"title"`
	Printer string            `json:"printer"`
	Options map[string]string `json:"options"`
	PDF     []byte            `json:"pdf"`
}

func TestPrinting(t *testing.T) {
	web.AgentWait = 2 * time.Second
	e := newEnv(t)
	ctx := context.Background()
	f := e.registrationFixture("Club")
	e.user("boss", f.tenant, store.RoleAdmin)
	e.user("helper", f.tenant, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	ev, _ := e.st.GetEvent(ctx, f.tenant, f.event)
	rangePath := fmt.Sprintf("/api/tenants/%d/ranges/%d", f.tenant, ev.Ranges[0].ID)

	lane := 1
	code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", f.tenant, f.event),
		web.RegistrationInput{ShooterID: f.shooter, Lane: &lane, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	expect(t, "register", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)

	for _, doc := range []string{"lane_sheet", "control_strip"} {
		resp, err := staff.Get(fmt.Sprintf("%s/api/tenants/%d/registrations/%d/documents/%s", e.srv.URL, f.tenant, reg.ID, doc))
		if err != nil {
			t.Fatal(err)
		}
		pdf, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 || resp.Header.Get("Content-Type") != "application/pdf" || !strings.HasPrefix(string(pdf), "%PDF") {
			t.Fatalf("%s: %d %s", doc, resp.StatusCode, resp.Header.Get("Content-Type"))
		}
	}
	code, b = e.do(staff, "GET", fmt.Sprintf("/api/tenants/%d/registrations/%d/documents/result_sheet", f.tenant, reg.ID), nil)
	expect(t, "result sheet before results", code, 409, b)

	printReq := map[string]any{"document": "lane_sheet", "registrationId": reg.ID}
	code, b = e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/print", f.tenant), printReq)
	expect(t, "print without printer", code, 200, b)
	if !strings.Contains(string(b), `"queued":false`) {
		t.Fatalf("print without printer should not queue: %s", b)
	}

	code, b = e.do(staff, "POST", rangePath+"/print-agents", map[string]any{"name": "Pi"})
	expect(t, "staff creates agent", code, 403, b)
	code, b = e.do(admin, "POST", rangePath+"/print-agents", map[string]any{"name": "Pi"})
	expect(t, "create agent", code, 201, b)
	var created struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(b, &created)
	testPath := rangePath + "/print-test"
	code, b = e.do(admin, "POST", testPath, map[string]any{"printer": "Office_A5"})
	expect(t, "test page without an online agent", code, 409, b)

	agentDo := func(token, method, path string, body any) (int, []byte) {
		return agentCall(t, e.srv.URL, token, method, path, body)
	}
	hello := func(body any) (int, []byte) { return agentDo(created.Secret, "POST", "/api/agent/hello", body) }
	cupsHello := map[string]any{"version": "v0.1.0", "os": "linux", "printers": []map[string]any{
		{"name": "Office_A5", "options": []map[string]any{{"name": "InputSlot", "label": "Media Source", "values": []string{"Auto", "Tray1", "Tray2"}, "default": "Auto"}}},
		{"name": "Brother_QL", "options": []map[string]any{{"name": "PageSize", "label": "Media Size", "values": []string{"29x90mm"}, "default": "29x90mm"}}},
	}}
	code, b = hello(cupsHello)
	expect(t, "hello", code, 200, b)
	next := func() *printJob {
		t.Helper()
		code, b := agentDo(created.Secret, "GET", "/api/agent/jobs/next", nil)
		if code == http.StatusNoContent {
			return nil
		}
		expect(t, "next job", code, 200, b)
		var j printJob
		json.Unmarshal(b, &j)
		return &j
	}
	finish := func(id int64) {
		t.Helper()
		code, b := agentDo(created.Secret, "POST", fmt.Sprintf("/api/agent/jobs/%d", id), map[string]any{"ok": true})
		expect(t, "finish job", code, 204, b)
	}
	code, b = e.do(admin, "GET", rangePath+"/print-agents", nil)
	expect(t, "get", code, 200, b)
	var agents []store.PrintAgent
	json.Unmarshal(b, &agents)
	if len(agents) != 1 || len(agents[0].Printers) != 2 || agents[0].Printers[0].Options[0].Values[2] != "Tray2" || agents[0].LastSeenAt == nil ||
		agents[0].Version != "v0.1.0" || agents[0].OS != "linux" {
		t.Fatalf("agents: %s", b)
	}

	code, b = e.do(admin, "PUT", rangePath+"/print-settings", map[string]any{
		"document": "lane_sheet", "printer": "Office_A5", "options": map[string]string{"InputSlot": "Tray2", "PageSize": "A5", "fit-to-page": ""},
	})
	expect(t, "route lane sheet", code, 204, b)

	code, b = hello(map[string]any{"printers": []map[string]any{{"name": "Brother_QL", "options": []any{}}}, "version": "v0.1.0", "os": "windows"})
	expect(t, "hello without the routed printer", code, 200, b)
	var h store.AgentHello
	json.Unmarshal(b, &h)
	if h.Club != "Club" || h.Range != ev.Ranges[0].Name || len(h.MissingPrinters) != 1 || h.MissingPrinters[0] != "Office_A5" {
		t.Fatalf("hello answer: %s", b)
	}
	code, b = e.do(admin, "GET", rangePath+"/print-agents", nil)
	expect(t, "get", code, 200, b)
	json.Unmarshal(b, &agents)
	if agents[0].Version != "v0.1.0" || agents[0].OS != "windows" {
		t.Fatalf("agent info: %s", b)
	}
	code, b = hello(map[string]any{"printers": []any{}, "os": "windows"})
	expect(t, "hello without a version", code, 422, b)
	code, b = hello(map[string]any{"printers": []any{}, "version": "v0.1.0", "os": "plan9"})
	expect(t, "hello from an unknown OS", code, 422, b)
	code, b = hello(cupsHello)
	expect(t, "hello again", code, 200, b)

	code, b = e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/print", f.tenant), printReq)
	expect(t, "print", code, 200, b)
	if !strings.Contains(string(b), `"queued":true`) {
		t.Fatalf("print should queue: %s", b)
	}

	j := next()
	if j == nil || j.Printer != "Office_A5" || j.Title != "Lane sheet Adult Anna" || !strings.HasPrefix(string(j.PDF), "%PDF") ||
		len(j.Options) != 3 || j.Options["InputSlot"] != "Tray2" || j.Options["PageSize"] != "A5" || j.Options["fit-to-page"] != "" {
		t.Fatalf("lane sheet job: %+v", j)
	}
	finish(j.ID)
	if j := next(); j != nil {
		t.Fatalf("queue should be empty: %+v", j)
	}

	code, b = e.do(staff, "GET", rangePath+"/print-jobs", nil)
	expect(t, "get", code, 200, b)
	var jobs []store.PrintJob
	json.Unmarshal(b, &jobs)
	if len(jobs) != 1 || jobs[0].Status != "done" || jobs[0].Printer != "Office_A5" {
		t.Fatalf("jobs: %s", b)
	}

	code, b = agentDo("wrong", "POST", "/api/agent/hello", cupsHello)
	expect(t, "wrong agent token", code, 401, b)

	code, b = e.do(staff, "POST", testPath, map[string]any{"printer": "Brother_QL"})
	expect(t, "staff test page", code, 403, b)
	code, b = e.do(admin, "POST", testPath, map[string]any{"printer": "Brother_QL", "options": map[string]string{"PageSize": "29x90mm"}})
	expect(t, "test page", code, 200, b)
	j = next()
	if j == nil || j.Printer != "Brother_QL" || j.Options["PageSize"] != "29x90mm" || j.Title != "Test page Brother_QL" || !strings.HasPrefix(string(j.PDF), "%PDF") {
		t.Fatalf("test page job: %+v", j)
	}
	finish(j.ID)
	code, b = e.do(admin, "POST", fmt.Sprintf("/api/tenants/%d/ranges/%d/print-test", f.tenant, 999999), map[string]any{"printer": "Brother_QL"})
	expect(t, "test page on an unknown range", code, 404, b)

	agentPath := fmt.Sprintf("/api/tenants/%d/print-agents/%d", f.tenant, agents[0].ID)
	code, b = e.do(staff, "DELETE", agentPath, nil)
	expect(t, "staff deletes agent", code, 403, b)
	code, b = e.do(admin, "DELETE", agentPath, nil)
	expect(t, "delete agent", code, 204, b)
	code, b = hello(cupsHello)
	expect(t, "deleted agent", code, 401, b)
}
