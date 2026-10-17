package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Printable documents.
const (
	DocLaneSheet    = "lane_sheet"
	DocControlStrip = "control_strip"
	DocResultSheet  = "result_sheet"
	DocWreathCard   = "wreath_card"
	DocForm         = "form"
)

// Documents lists the printable documents.
var Documents = []string{DocLaneSheet, DocControlStrip, DocResultSheet, DocWreathCard, DocForm}

// AgentOnline is how recently a print agent must have called in to count as online.
const AgentOnline = 60 * time.Second

// PrinterOption is a CUPS option of a printer with its possible values, as reported by `lpoptions -l`.
type PrinterOption struct {
	Name    string   `json:"name"`
	Label   string   `json:"label"`
	Values  []string `json:"values"`
	Default string   `json:"default"`
}

// Printer is a CUPS printer reported by a print agent.
type Printer struct {
	Name    string          `json:"name"`
	Options []PrinterOption `json:"options"`
}

// PrintAgent is a print agent registered for a range.
type PrintAgent struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	Printers   []Printer  `json:"printers"`
}

// AgentRef identifies the agent behind a token.
type AgentRef struct {
	ID       int64
	TenantID int64
	RangeID  int64
}

// PrintSetting routes one document of a range to a printer with lp options.
type PrintSetting struct {
	Document string            `json:"document" enum:"lane_sheet,control_strip,result_sheet,wreath_card,form"`
	Printer  string            `json:"printer"`
	Options  map[string]string `json:"options"`
}

// PrintJob is a queued or finished print job.
type PrintJob struct {
	ID         int64      `json:"id"`
	EventID    *int64     `json:"eventId,omitempty"`
	Document   string     `json:"document"`
	Title      string     `json:"title"`
	Printer    string     `json:"printer"`
	Status     string     `json:"status" enum:"queued,printing,done,failed"`
	Error      string     `json:"error"`
	Attempts   int        `json:"attempts"`
	CreatedAt  time.Time  `json:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

// ClaimedJob is a job handed to a print agent.
type ClaimedJob struct {
	ID      int64             `json:"id"`
	Title   string            `json:"title"`
	Printer string            `json:"printer"`
	Options map[string]string `json:"options"`
	PDF     []byte            `json:"pdf"`
}

// rangeInTenant returns ErrNotFound unless the range belongs to the tenant.
func rangeInTenant(ctx context.Context, q querier, tenantID, rangeID int64) error {
	var ok bool
	if err := q.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM ranges WHERE tenant_id = $1 AND id = $2)`, tenantID, rangeID).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// CreatePrintAgent registers an agent for a range by the hash of its token.
func (s *Store) CreatePrintAgent(ctx context.Context, tenantID, rangeID int64, name, tokenHash string) (int64, error) {
	if err := rangeInTenant(ctx, s.pool, tenantID, rangeID); err != nil {
		return 0, err
	}
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO print_agents (tenant_id, range_id, name, token_hash) VALUES ($1, $2, $3, $4) RETURNING id`,
		tenantID, rangeID, name, tokenHash).Scan(&id)
	return id, mapErr(err)
}

// ListPrintAgents lists the agents of a range.
func (s *Store) ListPrintAgents(ctx context.Context, tenantID, rangeID int64) ([]PrintAgent, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, created_at, last_seen_at, printers FROM print_agents
		WHERE tenant_id = $1 AND range_id = $2 ORDER BY created_at`, tenantID, rangeID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[PrintAgent])
}

// DeletePrintAgent revokes an agent.
func (s *Store) DeletePrintAgent(ctx context.Context, tenantID, agentID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM print_agents WHERE tenant_id = $1 AND id = $2`, tenantID, agentID))
}

// AgentByToken resolves an agent token hash and records the call.
func (s *Store) AgentByToken(ctx context.Context, tokenHash string) (AgentRef, error) {
	var a AgentRef
	err := s.pool.QueryRow(ctx, `UPDATE print_agents SET last_seen_at = now() WHERE token_hash = $1
		RETURNING id, tenant_id, range_id`, tokenHash).Scan(&a.ID, &a.TenantID, &a.RangeID)
	return a, mapErr(err)
}

// SetAgentPrinters stores the printers an agent reported.
func (s *Store) SetAgentPrinters(ctx context.Context, agentID int64, printers []Printer) error {
	if printers == nil {
		printers = []Printer{}
	}
	_, err := s.pool.Exec(ctx, `UPDATE print_agents SET printers = $2 WHERE id = $1`, agentID, printers)
	return err
}

// ListPrintSettings returns the document routing of a range.
func (s *Store) ListPrintSettings(ctx context.Context, tenantID, rangeID int64) ([]PrintSetting, error) {
	rows, err := s.pool.Query(ctx, `SELECT document, printer, options FROM print_settings
		WHERE tenant_id = $1 AND range_id = $2 ORDER BY document`, tenantID, rangeID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[PrintSetting])
}

// SetPrintSetting routes a document of a range to a printer, or removes the routing when printer is empty.
func (s *Store) SetPrintSetting(ctx context.Context, tenantID, rangeID int64, p PrintSetting) error {
	if err := rangeInTenant(ctx, s.pool, tenantID, rangeID); err != nil {
		return err
	}
	if p.Printer == "" {
		_, err := s.pool.Exec(ctx, `DELETE FROM print_settings WHERE range_id = $1 AND document = $2`, rangeID, p.Document)
		return err
	}
	if p.Options == nil {
		p.Options = map[string]string{}
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO print_settings (tenant_id, range_id, document, printer, options) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (range_id, document) DO UPDATE SET printer = EXCLUDED.printer, options = EXCLUDED.options`,
		tenantID, rangeID, p.Document, p.Printer, p.Options)
	return err
}

// ErrNoPrinter means a document cannot be printed on the range right now.
var ErrNoPrinter = errors.New("store: no printer configured or no print agent online")

// EnqueuePrint queues a PDF on the printer configured for the document, if an agent of the range is online.
func (s *Store) EnqueuePrint(ctx context.Context, tenantID, rangeID int64, eventID *int64, document, title string, pdf []byte) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO print_jobs (tenant_id, range_id, event_id, document, title, pdf, printer, options)
		SELECT $1, $2, $3, $4, $5, $6, ps.printer, ps.options FROM print_settings ps
		WHERE ps.tenant_id = $1 AND ps.range_id = $2 AND ps.document = $4
			AND EXISTS (SELECT 1 FROM print_agents a WHERE a.range_id = $2 AND a.last_seen_at > now() - $7::interval)
		RETURNING id`, tenantID, rangeID, eventID, document, title, pdf, AgentOnline.String()).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNoPrinter
	}
	return id, err
}

// ClaimPrintJob hands the oldest queued job of the agent's range to the agent.
func (s *Store) ClaimPrintJob(ctx context.Context, a AgentRef) (ClaimedJob, error) {
	var j ClaimedJob
	err := s.pool.QueryRow(ctx, `UPDATE print_jobs SET status = 'printing', claimed_at = now(), attempts = attempts + 1
		WHERE id = (SELECT id FROM print_jobs WHERE tenant_id = $1 AND range_id = $2 AND status = 'queued'
			ORDER BY created_at LIMIT 1 FOR UPDATE SKIP LOCKED)
		RETURNING id, title, printer, options, pdf`, a.TenantID, a.RangeID).Scan(&j.ID, &j.Title, &j.Printer, &j.Options, &j.PDF)
	return j, mapErr(err)
}

// FinishPrintJob records the outcome the agent reported and returns the job's event, if any.
func (s *Store) FinishPrintJob(ctx context.Context, a AgentRef, jobID int64, ok bool, message string) (*int64, error) {
	status := "done"
	if !ok {
		status = "failed"
	}
	var eventID *int64
	err := s.pool.QueryRow(ctx, `UPDATE print_jobs SET status = $4::print_status, error = $5, finished_at = now()
		WHERE tenant_id = $1 AND range_id = $2 AND id = $3 AND status = 'printing' RETURNING event_id`,
		a.TenantID, a.RangeID, jobID, status, message).Scan(&eventID)
	return eventID, mapErr(err)
}

// RequeueStalePrintJobs puts jobs back in the queue whose agent did not report within the timeout; after
// three attempts they fail.
func (s *Store) RequeueStalePrintJobs(ctx context.Context, timeout time.Duration) error {
	_, err := s.pool.Exec(ctx, `UPDATE print_jobs SET
			status = CASE WHEN attempts >= 3 THEN 'failed'::print_status ELSE 'queued'::print_status END,
			error = CASE WHEN attempts >= 3 THEN 'The print agent did not confirm the job' ELSE error END,
			finished_at = CASE WHEN attempts >= 3 THEN now() END
		WHERE status = 'printing' AND claimed_at < now() - $1::interval`, timeout.String())
	return err
}

// DeleteOldPrintJobs removes finished jobs older than a week; the PDFs are only needed while printing.
func (s *Store) DeleteOldPrintJobs(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM print_jobs WHERE created_at < now() - interval '7 days'`)
	return err
}

// ListPrintJobs lists the latest jobs of a range.
func (s *Store) ListPrintJobs(ctx context.Context, tenantID, rangeID int64, limit int) ([]PrintJob, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, event_id, document, title, printer, status, error, attempts, created_at, finished_at
		FROM print_jobs WHERE tenant_id = $1 AND range_id = $2 ORDER BY created_at DESC LIMIT $3`, tenantID, rangeID, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[PrintJob])
}
