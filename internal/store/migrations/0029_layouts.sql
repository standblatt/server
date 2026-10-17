-- Print forms: values placed onto pre-printed sheets such as the army's lane sheets. The background is a scan of
-- the empty form for the preview only; it is never printed.
CREATE TABLE layouts (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id       bigint NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name            text NOT NULL,
    definition      jsonb NOT NULL,
    background      bytea,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, name),
    UNIQUE (id, tenant_id)
);

ALTER TABLE programmes ADD COLUMN layout_id bigint,
    ADD FOREIGN KEY (layout_id, tenant_id) REFERENCES layouts(id, tenant_id) ON DELETE SET NULL (layout_id);

ALTER TABLE print_settings DROP CONSTRAINT print_settings_document_check;
ALTER TABLE print_settings ADD CONSTRAINT print_settings_document_check
    CHECK (document IN ('lane_sheet', 'control_strip', 'result_sheet', 'wreath_card', 'form'));
