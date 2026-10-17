package store

import (
	"context"
	"math"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// Profile is the personal data of a shooter using their own training log.
type Profile struct {
	Firstname     string `json:"firstname" maxLength:"100"`
	Lastname      string `json:"lastname" maxLength:"100"`
	Birthday      string `json:"birthday" required:"false" doc:"YYYY-MM-DD or empty"`
	LicenseNumber string `json:"licenseNumber" required:"false" maxLength:"20" doc:"Used with the birthday to find competition results"`
}

// GetProfile returns the profile of a user; ErrNotFound before onboarding.
func (s *Store) GetProfile(ctx context.Context, userID int64) (Profile, error) {
	var p Profile
	err := s.pool.QueryRow(ctx, `SELECT firstname, lastname, coalesce(birthday::text, ''), license_number FROM profiles WHERE user_id = $1`,
		userID).Scan(&p.Firstname, &p.Lastname, &p.Birthday, &p.LicenseNumber)
	return p, mapErr(err)
}

// SaveProfile creates or replaces the profile of a user.
func (s *Store) SaveProfile(ctx context.Context, userID int64, p Profile) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO profiles (user_id, firstname, lastname, birthday, license_number)
		VALUES ($1, $2, $3, nullif($4, '')::date, $5)
		ON CONFLICT (user_id) DO UPDATE SET firstname = EXCLUDED.firstname, lastname = EXCLUDED.lastname,
			birthday = EXCLUDED.birthday, license_number = EXCLUDED.license_number`,
		userID, p.Firstname, p.Lastname, p.Birthday, p.LicenseNumber)
	return mapErr(err)
}

// Discipline is a built-in discipline (Key set) or one a user defined.
type Discipline struct {
	ID             int64  `json:"id"`
	Key            string `json:"key" doc:"Identifier of a built-in discipline, empty for own ones"`
	Name           string `json:"name"`
	SeriesCount    int    `json:"seriesCount"`
	ShotsPerSeries int    `json:"shotsPerSeries"`
	MaxShot        int    `json:"maxShot"`
	Decimal        bool   `json:"decimal"`
	Target         string `json:"target" doc:"Target face for tapping shots, empty if none"`
}

// Rules returns the scoring of the discipline.
func (d Discipline) Rules() domain.Discipline {
	return domain.Discipline{SeriesCount: d.SeriesCount, ShotsPerSeries: d.ShotsPerSeries, MaxShot: d.MaxShot, Decimal: d.Decimal}
}

// DisciplineInput holds the fields of an own discipline.
type DisciplineInput struct {
	Name           string `json:"name" minLength:"1" maxLength:"60"`
	SeriesCount    int    `json:"seriesCount" minimum:"1" maximum:"20"`
	ShotsPerSeries int    `json:"shotsPerSeries" minimum:"1" maximum:"100"`
	MaxShot        int    `json:"maxShot" minimum:"1" maximum:"100"`
	Decimal        bool   `json:"decimal" doc:"Shots score 0.0 to 10.9; needs maxShot 10"`
	Target         string `json:"target" required:"false" enum:",air_pistol,air_rifle,pistol_precision,rifle_50,rifle_300"`
}

const disciplineCols = `id, coalesce(key, ''), name, series_count, shots_per_series, max_shot, decimal, target`

// ListDisciplines returns the built-in disciplines and the user's own ones.
func (s *Store) ListDisciplines(ctx context.Context, userID int64) ([]Discipline, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+disciplineCols+` FROM disciplines WHERE user_id IS NULL OR user_id = $1
		ORDER BY position, lower(name), id`, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Discipline])
}

// GetDiscipline returns a discipline the user may use.
func (s *Store) GetDiscipline(ctx context.Context, userID, id int64) (Discipline, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+disciplineCols+` FROM disciplines WHERE id = $2 AND (user_id IS NULL OR user_id = $1)`, userID, id)
	if err != nil {
		return Discipline{}, err
	}
	d, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[Discipline])
	return d, mapErr(err)
}

// CreateDiscipline adds an own discipline.
func (s *Store) CreateDiscipline(ctx context.Context, userID int64, in DisciplineInput) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO disciplines (user_id, name, series_count, shots_per_series, max_shot, decimal, target)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		userID, in.Name, in.SeriesCount, in.ShotsPerSeries, in.MaxShot, in.Decimal, in.Target).Scan(&id)
	return id, mapErr(err)
}

// UpdateDiscipline changes an own discipline. Sessions keep their series.
func (s *Store) UpdateDiscipline(ctx context.Context, userID, id int64, in DisciplineInput) error {
	return affected(s.pool.Exec(ctx, `UPDATE disciplines SET name = $3, series_count = $4, shots_per_series = $5, max_shot = $6,
		decimal = $7, target = $8 WHERE id = $2 AND user_id = $1`,
		userID, id, in.Name, in.SeriesCount, in.ShotsPerSeries, in.MaxShot, in.Decimal, in.Target))
}

// DeleteDiscipline removes an own discipline; ErrInUse while sessions use it.
func (s *Store) DeleteDiscipline(ctx context.Context, userID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM disciplines WHERE id = $2 AND user_id = $1`, userID, id))
}

// Weapon is one of the user's own guns, to compare results between them.
type Weapon struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// WeaponInput holds the fields of a weapon.
type WeaponInput struct {
	Name string `json:"name" minLength:"1" maxLength:"60"`
}

// ListWeapons returns the user's weapons by name.
func (s *Store) ListWeapons(ctx context.Context, userID int64) ([]Weapon, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name FROM weapons WHERE user_id = $1 ORDER BY lower(name), id`, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Weapon])
}

// GetWeapon returns a weapon of the user.
func (s *Store) GetWeapon(ctx context.Context, userID, id int64) (Weapon, error) {
	w := Weapon{ID: id}
	err := s.pool.QueryRow(ctx, `SELECT name FROM weapons WHERE user_id = $1 AND id = $2`, userID, id).Scan(&w.Name)
	return w, mapErr(err)
}

// CreateWeapon adds a weapon; ErrExists if the user has one with the same name.
func (s *Store) CreateWeapon(ctx context.Context, userID int64, in WeaponInput) (Weapon, error) {
	w := Weapon{Name: in.Name}
	err := s.pool.QueryRow(ctx, `INSERT INTO weapons (user_id, name) VALUES ($1, $2) RETURNING id`, userID, in.Name).Scan(&w.ID)
	return w, mapErr(err)
}

// RenameWeapon changes the name of a weapon of the user.
func (s *Store) RenameWeapon(ctx context.Context, userID, id int64, in WeaponInput) error {
	return affected(s.pool.Exec(ctx, `UPDATE weapons SET name = $3 WHERE user_id = $1 AND id = $2`, userID, id, in.Name))
}

// DeleteWeapon removes a weapon of the user; its sessions stay without a weapon.
func (s *Store) DeleteWeapon(ctx context.Context, userID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM weapons WHERE user_id = $1 AND id = $2`, userID, id))
}

// TrainingSession is a training or competition in a user's log.
type TrainingSession struct {
	ID           int64                   `json:"id"`
	DisciplineID int64                   `json:"disciplineId"`
	Kind         string                  `json:"kind" enum:"training,competition"`
	StartedAt    time.Time               `json:"startedAt"`
	Title        string                  `json:"title"`
	Place        string                  `json:"place"`
	Notes        string                  `json:"notes"`
	Tags         []string                `json:"tags"`
	WeaponID     *int64                  `json:"weaponId"`
	Rating       *int                    `json:"rating" doc:"How it felt, 1 to 5"`
	Imported     bool                    `json:"imported" doc:"Taken over from a competition run with Standblatt"`
	Rank         *int                    `json:"rank,omitempty" doc:"Imported competitions: place in the programme's ranking"`
	Participants *int                    `json:"participants,omitempty" doc:"Imported competitions: shooters ranked in the programme"`
	Series       []domain.TrainingSeries `json:"series"`
	Total        float64                 `json:"total"`
	Shots        int                     `json:"shots"`
}

// TrainingInput holds the fields of a session.
type TrainingInput struct {
	DisciplineID int64                   `json:"disciplineId"`
	Kind         string                  `json:"kind" enum:"training,competition"`
	StartedAt    time.Time               `json:"startedAt"`
	Title        string                  `json:"title" required:"false" maxLength:"200"`
	Place        string                  `json:"place" required:"false" maxLength:"200"`
	Notes        string                  `json:"notes" required:"false" maxLength:"5000"`
	Tags         []string                `json:"tags" required:"false" maxItems:"20"`
	WeaponID     *int64                  `json:"weaponId,omitempty" required:"false"`
	Rating       *int                    `json:"rating,omitempty" required:"false" minimum:"1" maximum:"5"`
	ClientID     string                  `json:"clientId,omitempty" required:"false" format:"uuid" doc:"Set by the client so a retried upload does not create a second session"`
	Series       []domain.TrainingSeries `json:"series" minItems:"1" maxItems:"30"`
}

// facts returns the session as input for the statistics.
func (t *TrainingSession) facts() domain.SessionFacts {
	f := domain.SessionFacts{ID: t.ID, StartedAt: t.StartedAt, Kind: t.Kind, Tags: t.Tags, Series: t.Series}
	if t.WeaponID != nil {
		f.WeaponID = *t.WeaponID
	}
	return f
}

const sessionCols = `ts.id, ts.discipline_id, ts.kind::text, ts.started_at, ts.title, ts.place, ts.notes, ts.tags, ts.weapon_id, ts.rating,
	ts.entry_id IS NOT NULL, ts.rank, ts.participants,
	coalesce((SELECT jsonb_agg(jsonb_build_object('shots', ss.shots, 'total', ss.total, 'values', ss.shot_values, 'x', ss.xs, 'y', ss.ys)
		ORDER BY ss.number) FROM session_series ss WHERE ss.session_id = ts.id), '[]')`

// scanSession scans a row of sessionCols and sums the total and shots of its series.
func scanSession(row pgx.CollectableRow) (TrainingSession, error) {
	var t TrainingSession
	err := row.Scan(&t.ID, &t.DisciplineID, &t.Kind, &t.StartedAt, &t.Title, &t.Place, &t.Notes, &t.Tags, &t.WeaponID, &t.Rating, &t.Imported, &t.Rank, &t.Participants, &t.Series)
	for _, s := range t.Series {
		t.Total += s.Total
		t.Shots += s.Shots
	}
	t.Total = math.Round(t.Total*10) / 10
	return t, err
}

// TrainingFilter narrows the session list; zero values mean no filter.
type TrainingFilter struct {
	DisciplineID int64
	Kind         string
	Before       *time.Time
	Limit        int
}

// ListTrainings returns the user's sessions, newest first.
func (s *Store) ListTrainings(ctx context.Context, userID int64, f TrainingFilter) ([]TrainingSession, error) {
	if f.Limit == 0 {
		f.Limit = 100000
	}
	rows, err := s.pool.Query(ctx, `SELECT `+sessionCols+` FROM training_sessions ts
		WHERE ts.user_id = $1 AND ($2 = 0 OR ts.discipline_id = $2) AND ($3 = '' OR ts.kind::text = $3)
			AND ($4::timestamptz IS NULL OR ts.started_at < $4)
		ORDER BY ts.started_at DESC, ts.id DESC LIMIT $5`, userID, f.DisciplineID, f.Kind, f.Before, f.Limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanSession)
}

// GetTraining returns one session of the user.
func (s *Store) GetTraining(ctx context.Context, userID, id int64) (TrainingSession, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+sessionCols+` FROM training_sessions ts WHERE ts.user_id = $1 AND ts.id = $2`, userID, id)
	if err != nil {
		return TrainingSession{}, err
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanSession)
	return t, mapErr(err)
}

// TrainingFacts returns the sessions of one discipline for the statistics.
func (s *Store) TrainingFacts(ctx context.Context, userID, disciplineID int64) ([]domain.SessionFacts, error) {
	list, err := s.ListTrainings(ctx, userID, TrainingFilter{DisciplineID: disciplineID})
	facts := make([]domain.SessionFacts, len(list))
	for i := range list {
		facts[i] = list[i].facts()
	}
	return facts, err
}

// saveSeries replaces the series of a session.
func saveSeries(ctx context.Context, tx pgx.Tx, sessionID int64, series []domain.TrainingSeries) error {
	if _, err := tx.Exec(ctx, `DELETE FROM session_series WHERE session_id = $1`, sessionID); err != nil {
		return err
	}
	for i, x := range series {
		if _, err := tx.Exec(ctx, `INSERT INTO session_series (session_id, number, shots, total, shot_values, xs, ys)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, sessionID, i+1, x.Shots, x.Total, x.Values, x.X, x.Y); err != nil {
			return err
		}
	}
	return nil
}

// SaveTraining creates a session (id 0) or replaces one of the user. A create with a client ID that the user
// already sent updates that session instead. ErrNotFound if the session is not the user's.
func (s *Store) SaveTraining(ctx context.Context, userID, id int64, in TrainingInput) (int64, error) {
	tags := in.Tags
	if tags == nil {
		tags = []string{}
	}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if id == 0 {
			err := tx.QueryRow(ctx, `INSERT INTO training_sessions (user_id, discipline_id, kind, started_at, title, place, notes, tags, rating, client_id, weapon_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, nullif($10, '')::uuid, $11)
				ON CONFLICT (user_id, client_id) DO UPDATE SET discipline_id = EXCLUDED.discipline_id, kind = EXCLUDED.kind,
					started_at = EXCLUDED.started_at, title = EXCLUDED.title, place = EXCLUDED.place, notes = EXCLUDED.notes,
					tags = EXCLUDED.tags, rating = EXCLUDED.rating, weapon_id = EXCLUDED.weapon_id, updated_at = now()
				RETURNING id`, userID, in.DisciplineID, in.Kind, in.StartedAt, in.Title, in.Place, in.Notes, tags, in.Rating, in.ClientID, in.WeaponID).Scan(&id)
			if err != nil {
				return err
			}
		} else {
			tag, err := tx.Exec(ctx, `UPDATE training_sessions SET discipline_id = $3, kind = $4, started_at = $5, title = $6, place = $7,
				notes = $8, tags = $9, rating = $10, weapon_id = $11, updated_at = now() WHERE user_id = $1 AND id = $2`,
				userID, id, in.DisciplineID, in.Kind, in.StartedAt, in.Title, in.Place, in.Notes, tags, in.Rating, in.WeaponID)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return pgx.ErrNoRows
			}
		}
		return saveSeries(ctx, tx, id, in.Series)
	})
	return id, mapErr(err)
}

// DeleteTraining removes a session of the user.
func (s *Store) DeleteTraining(ctx context.Context, userID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM training_sessions WHERE user_id = $1 AND id = $2`, userID, id))
}

// Tags returns the tags the user has used, most used first.
func (s *Store) Tags(ctx context.Context, userID int64) ([]string, error) {
	rows, err := s.pool.Query(ctx, `SELECT tag FROM training_sessions, unnest(tags) tag WHERE user_id = $1
		GROUP BY tag ORDER BY count(*) DESC, tag LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

// Goal is a target total for the standard match of a discipline.
type Goal struct {
	DisciplineID int64   `json:"disciplineId"`
	TargetTotal  float64 `json:"targetTotal" exclusiveMinimum:"0"`
	DueOn        string  `json:"dueOn" required:"false" doc:"YYYY-MM-DD or empty"`
}

// ListGoals returns the user's goals.
func (s *Store) ListGoals(ctx context.Context, userID int64) ([]Goal, error) {
	rows, err := s.pool.Query(ctx, `SELECT discipline_id, target_total::float8, coalesce(due_on::text, '') FROM goals WHERE user_id = $1
		ORDER BY discipline_id`, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Goal])
}

// SetGoal creates or replaces the goal of a discipline.
func (s *Store) SetGoal(ctx context.Context, userID int64, g Goal) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO goals (user_id, discipline_id, target_total, due_on) VALUES ($1, $2, $3, nullif($4, '')::date)
		ON CONFLICT (user_id, discipline_id) DO UPDATE SET target_total = EXCLUDED.target_total, due_on = EXCLUDED.due_on`,
		userID, g.DisciplineID, g.TargetTotal, g.DueOn)
	return mapErr(err)
}

// DeleteGoal removes the goal of a discipline.
func (s *Store) DeleteGoal(ctx context.Context, userID, disciplineID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM goals WHERE user_id = $1 AND discipline_id = $2`, userID, disciplineID))
}

// CompetitionResult is a completed entry of a Standblatt event whose shooter has the user's license number and birthday.
type CompetitionResult struct {
	EntryID      int64     `json:"entryId"`
	TenantID     int64     `json:"-"`
	EventID      int64     `json:"-"`
	Rank         int       `json:"rank" doc:"Place in the programme's ranking"`
	Participants int       `json:"participants" doc:"Shooters ranked in the programme"`
	Event        string    `json:"event"`
	Date         string    `json:"date"`
	Club         string    `json:"club" doc:"The club that ran the event"`
	Programme    string    `json:"programme"`
	Points       []float64 `json:"points"`
	Shots        []int     `json:"shots"`
	Total        float64   `json:"total"`
	MaxPerShot   float64   `json:"maxPerShot"`
}

const competitionQuery = `SELECT e.id, e.tenant_id, ev.id, ev.name, ev.starts_on::text, t.name, trim(p.name || ' ' || p.variant),
		array_agg(r.points::float8 ORDER BY r.series_number), array_agg(ps.shots ORDER BY r.series_number),
		sum(r.points)::float8, max(ps.max_points::float8 / ps.shots)
	FROM profiles pr
	JOIN shooters sh ON pr.license_number <> '' AND pr.birthday IS NOT NULL
		AND ltrim(sh.license_number, '0') = ltrim(pr.license_number, '0') AND sh.birthday = pr.birthday
	JOIN entries e ON e.shooter_id = sh.id AND e.tenant_id = sh.tenant_id AND e.status = 'completed'
	JOIN events ev ON ev.id = e.event_id
	JOIN tenants t ON t.id = e.tenant_id
	JOIN programmes p ON p.id = e.programme_id
	JOIN results r ON r.entry_id = e.id
	JOIN programme_series ps ON ps.programme_id = e.programme_id AND ps.number = r.series_number
	WHERE pr.user_id = $1 AND ($2 = 0 OR e.id = $2) AND ev.archived_at IS NOT NULL
		AND NOT EXISTS (SELECT 1 FROM training_sessions ts WHERE ts.user_id = $1 AND ts.entry_id = e.id)
	GROUP BY e.id, ev.id, ev.name, ev.starts_on, t.name, p.name, p.variant
	ORDER BY ev.starts_on DESC, e.id`

// CompetitionResults lists the user's competition results of closed events that are not in the log yet, with their rank.
func (s *Store) CompetitionResults(ctx context.Context, userID int64) ([]CompetitionResult, error) {
	return s.competitionResults(ctx, userID, 0)
}

// competitionResults lists the user's competition results not in the log yet, only the given entry unless entryID is 0.
// Only closed (archived) events count: their rankings are final.
func (s *Store) competitionResults(ctx context.Context, userID, entryID int64) ([]CompetitionResult, error) {
	rows, err := s.pool.Query(ctx, competitionQuery, userID, entryID)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (CompetitionResult, error) {
		var c CompetitionResult
		err := row.Scan(&c.EntryID, &c.TenantID, &c.EventID, &c.Event, &c.Date, &c.Club, &c.Programme, &c.Points, &c.Shots, &c.Total, &c.MaxPerShot)
		return c, err
	})
	if err != nil {
		return nil, err
	}
	placings := map[int64]map[int64]Placing{}
	for i, c := range list {
		if placings[c.EventID] == nil {
			if placings[c.EventID], err = s.Placings(ctx, c.TenantID, c.EventID); err != nil {
				return nil, err
			}
		}
		list[i].Rank, list[i].Participants = placings[c.EventID][c.EntryID].Rank, placings[c.EventID][c.EntryID].Participants
	}
	return list, nil
}

// ImportCompetition adds a competition result to the user's log as a competition session of the given discipline.
// ErrNotFound if the result is not the user's or already imported.
func (s *Store) ImportCompetition(ctx context.Context, userID, entryID, disciplineID int64, loc *time.Location) (int64, error) {
	list, err := s.competitionResults(ctx, userID, entryID)
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, ErrNotFound
	}
	c := list[0]
	date, err := time.ParseInLocation("2006-01-02", c.Date, loc)
	if err != nil {
		return 0, err
	}
	var id int64
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `INSERT INTO training_sessions (user_id, discipline_id, kind, started_at, title, place, entry_id, rank, participants)
			VALUES ($1, $2, 'competition', $3, $4, $5, $6, nullif($7, 0), nullif($8, 0)) RETURNING id`,
			userID, disciplineID, date.Add(12*time.Hour), c.Event+" · "+c.Programme, c.Club, entryID, c.Rank, c.Participants).Scan(&id)
		if err != nil {
			return err
		}
		series := make([]domain.TrainingSeries, len(c.Points))
		for i := range c.Points {
			series[i] = domain.TrainingSeries{Shots: c.Shots[i], Total: math.Round(c.Points[i]*10) / 10}
		}
		return saveSeries(ctx, tx, id, series)
	})
	return id, mapErr(err)
}

// Activity is the start and the number of shots of a session.
type Activity struct {
	At    time.Time
	Shots int
}

// TrainingActivity returns start and shots of the user's sessions of the last two years.
func (s *Store) TrainingActivity(ctx context.Context, userID int64) ([]Activity, error) {
	rows, err := s.pool.Query(ctx, `SELECT ts.started_at, coalesce((SELECT sum(shots) FROM session_series ss WHERE ss.session_id = ts.id), 0)
		FROM training_sessions ts WHERE ts.user_id = $1 AND ts.started_at > now() - interval '2 years'`, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Activity])
}
