DROP TABLE club_seats;
DROP TABLE club_plans;

ALTER TABLE users
    DROP COLUMN trial_ends_at,
    DROP COLUMN paid_until,
    DROP COLUMN stripe_customer_id,
    DROP COLUMN stripe_subscription_id;
