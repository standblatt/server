-- Building blocks per club (see internal/modules); existing clubs keep print forms.
ALTER TABLE tenants ADD COLUMN modules text[] NOT NULL DEFAULT '{}';
UPDATE tenants SET modules = ARRAY['forms'];
