// Package agent is the print agent that runs next to the printers at a range. It only makes outgoing
// requests: it reports the CUPS printers, waits for jobs and prints them with lp.
package agent

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Printer is a CUPS printer with its options, as sent to the server.
type Printer struct {
	Name    string   `json:"name"`
	Options []Option `json:"options"`
}

// Option is one `lpoptions -l` line: name, label, possible values and the default.
type Option struct {
	Name    string   `json:"name"`
	Label   string   `json:"label"`
	Values  []string `json:"values"`
	Default string   `json:"default"`
}

// Job is a print job handed out by the server.
type Job struct {
	ID      int64             `json:"id"`
	Title   string            `json:"title"`
	Printer string            `json:"printer"`
	Options map[string]string `json:"options"`
	PDF     []byte            `json:"pdf"`
}

// Agent talks to one Standblatt server.
type Agent struct {
	URL       string
	Token     string
	LP        string
	LPStat    string
	LPOptions string
	Version   string
	Client    *http.Client
	Log       *slog.Logger
}

// ParseOptions reads the output of `lpoptions -p <printer> -l`, e.g. "InputSlot/Media Source: Auto *Tray1 Tray2".
func ParseOptions(out string) []Option {
	var opts []Option
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		head, values, ok := strings.Cut(sc.Text(), ":")
		if !ok {
			continue
		}
		name, label, _ := strings.Cut(head, "/")
		o := Option{Name: strings.TrimSpace(name), Label: strings.TrimSpace(label), Values: []string{}}
		for _, v := range strings.Fields(values) {
			if strings.HasPrefix(v, "*") {
				v = v[1:]
				o.Default = v
			}
			o.Values = append(o.Values, v)
		}
		if o.Name != "" {
			opts = append(opts, o)
		}
	}
	return opts
}

// Discover lists the local CUPS printers with their options.
func (a *Agent) Discover(ctx context.Context) ([]Printer, error) {
	out, err := exec.CommandContext(ctx, a.LPStat, "-e").Output()
	if err != nil {
		return nil, fmt.Errorf("lpstat: %w", err)
	}
	printers := []Printer{}
	for _, name := range strings.Fields(string(out)) {
		opts, err := exec.CommandContext(ctx, a.LPOptions, "-p", name, "-l").Output()
		if err != nil {
			a.Log.Warn("lpoptions", "printer", name, "err", err)
		}
		printers = append(printers, Printer{Name: name, Options: ParseOptions(string(opts))})
	}
	return printers, nil
}

// request sends an authenticated request with an optional JSON body; status codes of 300 and above become errors.
func (a *Agent) request(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimSuffix(a.URL, "/")+path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%s %s: %s %s", method, path, resp.Status, strings.TrimSpace(string(b)))
	}
	return resp, nil
}

// Hello reports the printers, the version (dev when unset) and the operating system to the server.
func (a *Agent) Hello(ctx context.Context) error {
	printers, err := a.Discover(ctx)
	if err != nil {
		return err
	}
	resp, err := a.request(ctx, http.MethodPost, "/api/agent/hello", map[string]any{"printers": printers, "version": cmp.Or(a.Version, "dev"), "os": runtime.GOOS})
	if err != nil {
		return err
	}
	resp.Body.Close()
	a.Log.Info("printers reported", "count", len(printers))
	return nil
}

// Next waits for the next job; nil means none arrived within the server's wait time.
func (a *Agent) Next(ctx context.Context) (*Job, error) {
	resp, err := a.request(ctx, http.MethodGet, "/api/agent/jobs/next", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	var j Job
	return &j, json.NewDecoder(resp.Body).Decode(&j)
}

// Print sends a job to CUPS with lp.
func (a *Agent) Print(ctx context.Context, j *Job) error {
	f, err := os.CreateTemp("", "standblatt-*.pdf")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(j.PDF); err != nil {
		f.Close()
		return err
	}
	f.Close()
	args := []string{"-d", j.Printer, "-t", j.Title}
	keys := make([]string, 0, len(j.Options))
	for k := range j.Options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v := j.Options[k]; v == "" {
			args = append(args, "-o", k)
		} else {
			args = append(args, "-o", k+"="+v)
		}
	}
	args = append(args, f.Name())
	out, err := exec.CommandContext(ctx, a.LP, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("lp: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Finish reports the outcome of a job.
func (a *Agent) Finish(ctx context.Context, id int64, printErr error) error {
	body := map[string]any{"ok": printErr == nil}
	if printErr != nil {
		body["error"] = printErr.Error()
	}
	resp, err := a.request(ctx, http.MethodPost, fmt.Sprintf("/api/agent/jobs/%d", id), body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Run reports the printers every five minutes and prints jobs until ctx ends.
func (a *Agent) Run(ctx context.Context) {
	lastHello := time.Time{}
	backoff := time.Second
	for ctx.Err() == nil {
		if time.Since(lastHello) > 5*time.Minute {
			if err := a.Hello(ctx); err != nil {
				a.Log.Warn("report printers", "err", err)
			} else {
				lastHello = time.Now()
			}
		}
		j, err := a.Next(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			a.Log.Warn("waiting for jobs", "err", err, "retry", backoff)
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			backoff = min(backoff*2, time.Minute)
			continue
		}
		backoff = time.Second
		if j == nil {
			continue
		}
		printErr := a.Print(ctx, j)
		if printErr != nil {
			a.Log.Warn("print failed", "job", j.ID, "printer", j.Printer, "err", printErr)
		} else {
			a.Log.Info("printed", "job", j.ID, "printer", j.Printer, "title", j.Title)
		}
		if err := a.Finish(ctx, j.ID, printErr); err != nil {
			a.Log.Warn("report job", "job", j.ID, "err", err)
		}
	}
}
