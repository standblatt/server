package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// ImportRow is one line of a shooter import: a shooter and optionally one of their clubs, either an existing one
// (ClubID) or a new one to create by name (Club).
type ImportRow struct {
	ShooterInput
	Club   string
	ClubID int64
}

// ImportResult counts what an import changed.
type ImportResult struct {
	Created      int `json:"created"`
	Updated      int `json:"updated"`
	ClubsCreated int `json:"clubsCreated"`
	Links        int `json:"links"`
}

// ImportShooters merges rows into a tenant's master data in one transaction. Shooters are matched by
// license number, or by name and birthday when the row has none; empty fields keep the stored value.
// Clubs named without an ID are found regardless of case and spaces or created. Nothing is deleted.
func (s *Store) ImportShooters(ctx context.Context, tenantID int64, rows []ImportRow) (ImportResult, error) {
	var res ImportResult
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		clubs := map[string]int64{}
		seen := map[int64]bool{}
		for _, r := range rows {
			id, created, err := upsertShooter(ctx, tx, tenantID, r.ShooterInput)
			if err != nil {
				return err
			}
			switch {
			case created:
				res.Created++
				seen[id] = true
			case !seen[id]:
				res.Updated++
				seen[id] = true
			}
			if r.Club == "" && r.ClubID == 0 {
				continue
			}
			clubID, ok := r.ClubID, r.ClubID != 0
			if !ok {
				clubID, ok = clubs[r.Club]
			}
			if !ok {
				var newClub bool
				if clubID, newClub, err = ensureClub(ctx, tx, tenantID, r.Club); err != nil {
					return err
				}
				if newClub {
					res.ClubsCreated++
				}
				clubs[r.Club] = clubID
			}
			tag, err := tx.Exec(ctx, `INSERT INTO shooter_clubs (tenant_id, shooter_id, club_id) VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`, tenantID, id, clubID)
			if err != nil {
				return err
			}
			res.Links += int(tag.RowsAffected())
		}
		return nil
	})
	return res, err
}

// upsertShooter merges a shooter by licence number, or without one by name and birthday, and inserts it otherwise;
// it returns the ID and whether the shooter was created.
func upsertShooter(ctx context.Context, tx pgx.Tx, tenantID int64, in ShooterInput) (int64, bool, error) {
	var id int64
	var created bool
	if in.LicenseNumber != "" {
		err := tx.QueryRow(ctx, `INSERT INTO shooters (tenant_id, license_number, firstname, lastname, birthday, postalcode, city)
			VALUES ($1, $2, $3, $4, nullif($5, '')::date, $6, $7)
			ON CONFLICT (tenant_id, license_number) DO UPDATE SET
				firstname = EXCLUDED.firstname, lastname = EXCLUDED.lastname,
				birthday = coalesce(EXCLUDED.birthday, shooters.birthday),
				postalcode = coalesce(nullif(EXCLUDED.postalcode, ''), shooters.postalcode),
				city = coalesce(nullif(EXCLUDED.city, ''), shooters.city),
				updated_at = now()
			RETURNING id, xmax = 0`,
			tenantID, in.LicenseNumber, in.Firstname, in.Lastname, in.Birthday, in.Postalcode, in.City).Scan(&id, &created)
		return id, created, err
	}
	err := tx.QueryRow(ctx, `UPDATE shooters SET
			postalcode = coalesce(nullif($5, ''), postalcode), city = coalesce(nullif($6, ''), city), updated_at = now()
		WHERE id = (SELECT id FROM shooters WHERE tenant_id = $1 AND license_number IS NULL
			AND lower(firstname) = lower($2) AND lower(lastname) = lower($3)
			AND birthday IS NOT DISTINCT FROM nullif($4, '')::date ORDER BY id LIMIT 1)
		RETURNING id`, tenantID, in.Firstname, in.Lastname, in.Birthday, in.Postalcode, in.City).Scan(&id)
	if !errors.Is(err, pgx.ErrNoRows) {
		return id, false, err
	}
	err = tx.QueryRow(ctx, `INSERT INTO shooters (tenant_id, firstname, lastname, birthday, postalcode, city)
		VALUES ($1, $2, $3, nullif($4, '')::date, $5, $6) RETURNING id`,
		tenantID, in.Firstname, in.Lastname, in.Birthday, in.Postalcode, in.City).Scan(&id)
	return id, true, err
}
