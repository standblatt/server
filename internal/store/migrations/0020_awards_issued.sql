ALTER TYPE award_kind ADD VALUE 'cash';

ALTER TABLE registrations ADD COLUMN awards_issued_at timestamptz;
