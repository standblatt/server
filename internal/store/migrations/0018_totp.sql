ALTER TABLE users
    ADD COLUMN totp_secret  text,
    ADD COLUMN totp_counter bigint NOT NULL DEFAULT 0;
