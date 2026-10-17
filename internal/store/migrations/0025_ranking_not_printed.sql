-- Rankings are no longer sent to a printer, only downloaded as PDF or CSV.
DELETE FROM print_settings WHERE document = 'ranking';
ALTER TABLE print_settings DROP CONSTRAINT print_settings_document_check;
ALTER TABLE print_settings ADD CONSTRAINT print_settings_document_check CHECK (document IN ('lane_sheet', 'control_strip', 'result_sheet', 'wreath_card'));
