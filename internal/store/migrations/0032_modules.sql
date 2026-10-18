-- Building blocks per event (see internal/modules). Existing events keep everything they had: all modules, except
-- the lane sheet, control strip and ammunition desk where the old settings switched them off.
ALTER TABLE events ADD COLUMN modules text[] NOT NULL DEFAULT '{}';

SET LOCAL standblatt.merge = 'on';
UPDATE events SET
    modules = ARRAY['reservations', 'booking', 'attempts', 'awards', 'special_prizes', 'club_competition', 'classifications', 'kiosk', 'helpers']
        || CASE WHEN coalesce((settings->>'noLaneSheet')::boolean, false) THEN '{}'::text[] ELSE ARRAY['lane_sheet'] END
        || CASE WHEN coalesce((settings->>'noControlStrip')::boolean, false) THEN '{}'::text[] ELSE ARRAY['control_strip'] END
        || CASE WHEN coalesce((settings->>'ammoIssue')::boolean, false) THEN ARRAY['ammo'] ELSE '{}'::text[] END,
    settings = settings - 'noLaneSheet' - 'noControlStrip' - 'ammoIssue';
