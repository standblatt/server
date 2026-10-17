-- Imported competition results keep the shooter's rank in the programme at the time of the import.
ALTER TABLE training_sessions
    ADD COLUMN rank         integer CHECK (rank > 0),
    ADD COLUMN participants integer CHECK (participants > 0);
