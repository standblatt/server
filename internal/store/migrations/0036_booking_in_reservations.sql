-- Online booking is part of the reservations module: events with booking keep reservations.
SET LOCAL standblatt.merge = 'on';
UPDATE events SET modules = array_remove(modules, 'booking') || CASE WHEN 'reservations' = ANY(modules) THEN '{}'::text[] ELSE ARRAY['reservations'] END
WHERE 'booking' = ANY(modules);
