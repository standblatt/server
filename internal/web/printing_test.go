package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/agent"
	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func fakeCUPS(t *testing.T) (dir string, lpLog string) {
	dir = t.TempDir()
	lpLog = filepath.Join(dir, "lp.log")
	scripts := map[string]string{
		"lpstat":    "#!/bin/sh\necho 'Office_A5 Brother_QL'\n",
		"lpoptions": "#!/bin/sh\ncase \"$2\" in Office_A5) echo 'InputSlot/Media Source: *Auto Tray1 Tray2';; *) echo 'PageSize/Media Size: *29x90mm';; esac\n",
		"lp":        "#!/bin/sh\necho \"$@\" >> " + lpLog + "\nhead -c 5 \"$(eval echo \\${$#})\" >> " + lpLog + "\necho >> " + lpLog + "\n",
	}
	for name, body := range scripts {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return dir, lpLog
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

	dir, lpLog := fakeCUPS(t)
	a := &agent.Agent{
		URL: e.srv.URL, Token: created.Secret, Client: &http.Client{Timeout: 10 * time.Second},
		LP: filepath.Join(dir, "lp"), LPStat: filepath.Join(dir, "lpstat"), LPOptions: filepath.Join(dir, "lpoptions"),
		Log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if err := a.Hello(ctx); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(admin, "GET", rangePath+"/print-agents", nil)
	expect(t, "get", code, 200, b)
	var agents []store.PrintAgent
	json.Unmarshal(b, &agents)
	if len(agents) != 1 || len(agents[0].Printers) != 2 || agents[0].Printers[0].Options[0].Values[2] != "Tray2" || agents[0].LastSeenAt == nil {
		t.Fatalf("agents: %s", b)
	}

	code, b = e.do(admin, "PUT", rangePath+"/print-settings", map[string]any{
		"document": "lane_sheet", "printer": "Office_A5", "options": map[string]string{"InputSlot": "Tray2", "PageSize": "A5", "fit-to-page": ""},
	})
	expect(t, "route lane sheet", code, 204, b)

	code, b = e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/print", f.tenant), printReq)
	expect(t, "print", code, 200, b)
	if !strings.Contains(string(b), `"queued":true`) {
		t.Fatalf("print should queue: %s", b)
	}

	runCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	j, err := a.Next(runCtx)
	if err != nil || j == nil {
		t.Fatalf("agent got no job: %v", err)
	}
	if err := a.Finish(runCtx, j.ID, a.Print(runCtx, j)); err != nil {
		t.Fatal(err)
	}
	log, _ := os.ReadFile(lpLog)
	if !strings.Contains(string(log), "-d Office_A5 -t Lane sheet Adult Anna -o InputSlot=Tray2 -o PageSize=A5 -o fit-to-page") || !strings.Contains(string(log), "%PDF") {
		t.Fatalf("lp was called with: %s", log)
	}
	if j, err := a.Next(runCtx); err != nil || j != nil {
		t.Fatalf("queue should be empty: %+v %v", j, err)
	}

	code, b = e.do(staff, "GET", rangePath+"/print-jobs", nil)
	expect(t, "get", code, 200, b)
	var jobs []store.PrintJob
	json.Unmarshal(b, &jobs)
	if len(jobs) != 1 || jobs[0].Status != "done" || jobs[0].Printer != "Office_A5" {
		t.Fatalf("jobs: %s", b)
	}

	bad := &agent.Agent{URL: e.srv.URL, Token: "wrong", Client: http.DefaultClient, Log: a.Log, LPStat: a.LPStat, LPOptions: a.LPOptions}
	if err := bad.Hello(ctx); err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("wrong agent token: %v", err)
	}

	code, b = e.do(staff, "POST", testPath, map[string]any{"printer": "Brother_QL"})
	expect(t, "staff test page", code, 403, b)
	code, b = e.do(admin, "POST", testPath, map[string]any{"printer": "Brother_QL", "options": map[string]string{"PageSize": "29x90mm"}})
	expect(t, "test page", code, 200, b)
	j, err = a.Next(runCtx)
	if err != nil || j == nil || j.Printer != "Brother_QL" || j.Options["PageSize"] != "29x90mm" || j.Title != "Test page Brother_QL" || !strings.HasPrefix(string(j.PDF), "%PDF") {
		t.Fatalf("test page job: %+v %v", j, err)
	}
	if err := a.Finish(runCtx, j.ID, nil); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(admin, "POST", fmt.Sprintf("/api/tenants/%d/ranges/%d/print-test", f.tenant, 999999), map[string]any{"printer": "Brother_QL"})
	expect(t, "test page on an unknown range", code, 404, b)

	agentPath := fmt.Sprintf("/api/tenants/%d/print-agents/%d", f.tenant, agents[0].ID)
	code, b = e.do(staff, "DELETE", agentPath, nil)
	expect(t, "staff deletes agent", code, 403, b)
	code, b = e.do(admin, "DELETE", agentPath, nil)
	expect(t, "delete agent", code, 204, b)
	if err := a.Hello(ctx); err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("deleted agent: %v", err)
	}
}
