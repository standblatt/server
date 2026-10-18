-- Lane assignment is a module; every event keeps it (events without a range lose it through modules.Effective).
SET LOCAL standblatt.merge = 'on';
UPDATE events SET modules = ARRAY['lanes'] || modules WHERE NOT 'lanes' = ANY(modules);
