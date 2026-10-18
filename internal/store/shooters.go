package store

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/standblatt/server/internal/domain"
)

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// Shooter is a person in a tenant's master data. Empty LicenseNumber and Birthday mean unknown.
type Shooter struct {
	ID            int64     `json:"id"`
	LicenseNumber string    `json:"licenseNumber"`
	Firstname     string    `json:"firstname"`
	Lastname      string    `json:"lastname"`
	Birthday      string    `json:"birthday" doc:"YYYY-MM-DD or empty"`
	Postalcode    string    `json:"postalcode"`
	City          string    `json:"city"`
	Clubs         []ClubRef `json:"clubs"`
}

// ShooterInput holds the editable fields of a shooter.
type ShooterInput struct {
	LicenseNumber string `json:"licenseNumber" maxLength:"20" required:"false"`
	Firstname     string `json:"firstname" minLength:"1" maxLength:"100"`
	Lastname      string `json:"lastname" minLength:"1" maxLength:"100"`
	Birthday      string `json:"birthday" required:"false" doc:"YYYY-MM-DD or empty"`
	Postalcode    string `json:"postalcode" maxLength:"10" required:"false"`
	City          string `json:"city" maxLength:"100" required:"false"`
}

const shooterCols = `s.id, coalesce(s.license_number, ''), s.firstname, s.lastname, coalesce(s.birthday::text, ''),
	s.postalcode, s.city,
	coalesce((SELECT json_agg(json_build_object('id', c.id, 'name', c.name) ORDER BY c.name)
		FROM shooter_clubs sc JOIN clubs c ON c.id = sc.club_id WHERE sc.shooter_id = s.id), '[]')`

var digitsRe = regexp.MustCompile(`^[0-9]+$`)

// likeEscape escapes the LIKE wildcards and the backslash in s.
func likeEscape(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

// searchFilter turns a search term into a WHERE fragment and its arguments, starting at placeholder $n.
// Digits only: exact license number, zero-padded to six digits. Two words: first and last name in
// either order. Otherwise: substring of license number, first name or last name.
func searchFilter(q string, n int) (string, []any) {
	q = strings.TrimSpace(q)
	p := func(i int) string { return fmt.Sprintf("$%d", n+i) }
	switch {
	case q == "":
		return "true", nil
	case digitsRe.MatchString(q):
		return "s.license_number = " + p(0), []any{domain.NormalizeLicense(q)}
	case strings.Contains(q, " "):
		a, b, _ := strings.Cut(q, " ")
		a, b = likeEscape(a)+"%", likeEscape(strings.TrimSpace(b))+"%"
		return fmt.Sprintf(`((s.firstname ILIKE %[1]s AND s.lastname ILIKE %[2]s) OR (s.firstname ILIKE %[2]s AND s.lastname ILIKE %[1]s))`, p(0), p(1)),
			[]any{a, b}
	default:
		like := "%" + likeEscape(q) + "%"
		return fmt.Sprintf(`(s.license_number ILIKE %[1]s OR s.firstname ILIKE %[1]s OR s.lastname ILIKE %[1]s)`, p(0)), []any{like}
	}
}

// SearchShooters finds shooters of a tenant by license number or name, optionally within one club, one page of
// limit rows from offset, ordered by name.
func (s *Store) SearchShooters(ctx context.Context, tenantID int64, q string, clubID int64, limit, offset int) ([]Shooter, error) {
	where, args := shooterFilter(tenantID, q, clubID)
	args = append(args, limit, offset)
	rows, err := s.pool.Query(ctx, `SELECT `+shooterCols+` FROM shooters s WHERE `+where+
		fmt.Sprintf(` ORDER BY lower(s.lastname), lower(s.firstname), s.id LIMIT $%d OFFSET $%d`, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Shooter])
}

// CountShooters counts the shooters SearchShooters finds without limit.
func (s *Store) CountShooters(ctx context.Context, tenantID int64, q string, clubID int64) (int, error) {
	where, args := shooterFilter(tenantID, q, clubID)
	var n int
	err := s.pool.QueryRow(ctx, `SELECT count(*) FROM shooters s WHERE `+where, args...).Scan(&n)
	return n, err
}

// shooterFilter returns the WHERE clause and arguments of a shooter search in a tenant, optionally in one club.
func shooterFilter(tenantID int64, q string, clubID int64) (string, []any) {
	where, args := searchFilter(q, 2)
	args = append([]any{tenantID}, args...)
	if clubID != 0 {
		args = append(args, clubID)
		where += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM shooter_clubs sc WHERE sc.shooter_id = s.id AND sc.club_id = $%d)", len(args))
	}
	return "s.tenant_id = $1 AND " + where, args
}

// GetShooter returns one shooter of a tenant.
func (s *Store) GetShooter(ctx context.Context, tenantID, shooterID int64) (Shooter, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+shooterCols+` FROM shooters s WHERE s.tenant_id = $1 AND s.id = $2`, tenantID, shooterID)
	if err != nil {
		return Shooter{}, err
	}
	sh, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[Shooter])
	return sh, mapErr(err)
}

// CreateShooter adds a shooter and assigns the given clubs; ErrExists if the license number is taken,
// ErrInUse if a club does not belong to the tenant.
func (s *Store) CreateShooter(ctx context.Context, tenantID int64, in ShooterInput, clubIDs []int64) (int64, error) {
	var id int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `INSERT INTO shooters (tenant_id, license_number, firstname, lastname, birthday, postalcode, city)
			VALUES ($1, nullif($2, ''), $3, $4, nullif($5, '')::date, $6, $7) RETURNING id`,
			tenantID, in.LicenseNumber, in.Firstname, in.Lastname, in.Birthday, in.Postalcode, in.City).Scan(&id)
		if err != nil {
			return err
		}
		for _, c := range clubIDs {
			if _, err := tx.Exec(ctx, `INSERT INTO shooter_clubs (tenant_id, shooter_id, club_id) VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`, tenantID, id, c); err != nil {
				return err
			}
		}
		return nil
	})
	return id, mapErr(err)
}

// UpdateShooter changes the personal data of a shooter.
func (s *Store) UpdateShooter(ctx context.Context, tenantID, shooterID int64, in ShooterInput) error {
	return affected(s.pool.Exec(ctx, `UPDATE shooters SET license_number = nullif($3, ''), firstname = $4, lastname = $5,
		birthday = nullif($6, '')::date, postalcode = $7, city = $8, updated_at = now()
		WHERE tenant_id = $1 AND id = $2`,
		tenantID, shooterID, in.LicenseNumber, in.Firstname, in.Lastname, in.Birthday, in.Postalcode, in.City))
}

// DeleteShooter removes a shooter of a tenant.
func (s *Store) DeleteShooter(ctx context.Context, tenantID, shooterID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM shooters WHERE tenant_id = $1 AND id = $2`, tenantID, shooterID))
}

// AddShooterClub assigns a club to a shooter; ErrNotFound if either does not belong to the tenant.
func (s *Store) AddShooterClub(ctx context.Context, tenantID, shooterID, clubID int64) error {
	tag, err := s.pool.Exec(ctx, `INSERT INTO shooter_clubs (tenant_id, shooter_id, club_id)
		SELECT $1, sh.id, c.id FROM shooters sh, clubs c
		WHERE sh.tenant_id = $1 AND sh.id = $2 AND c.tenant_id = $1 AND c.id = $3
		ON CONFLICT DO NOTHING`, tenantID, shooterID, clubID)
	if err != nil || tag.RowsAffected() == 1 {
		return mapErr(err)
	}
	var linked bool
	err = s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM shooter_clubs
		WHERE tenant_id = $1 AND shooter_id = $2 AND club_id = $3)`, tenantID, shooterID, clubID).Scan(&linked)
	if err == nil && !linked {
		return ErrNotFound
	}
	return err
}

// RemoveShooterClub removes a club assignment.
func (s *Store) RemoveShooterClub(ctx context.Context, tenantID, shooterID, clubID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM shooter_clubs WHERE tenant_id = $1 AND shooter_id = $2 AND club_id = $3`,
		tenantID, shooterID, clubID))
}
