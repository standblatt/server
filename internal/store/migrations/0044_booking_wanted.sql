-- Online booking: the lane numbers a club picked for an item; empty = any free lanes.
ALTER TABLE booking_items ADD COLUMN wanted integer[] NOT NULL DEFAULT '{}';
-- Confirmed items want the lanes they got, so saving an unchanged request from the manage link keeps it confirmed.
UPDATE booking_items i SET wanted = r.lanes
    FROM reservations r WHERE r.booking_item_id = i.id;
