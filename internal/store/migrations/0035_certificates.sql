-- Certificates printed from the rankings get a printer of their own.
ALTER TABLE print_settings DROP CONSTRAINT print_settings_document_check;
ALTER TABLE print_settings ADD CONSTRAINT print_settings_document_check
    CHECK (document IN ('lane_sheet', 'control_strip', 'result_sheet', 'wreath_card', 'form', 'certificate'));
