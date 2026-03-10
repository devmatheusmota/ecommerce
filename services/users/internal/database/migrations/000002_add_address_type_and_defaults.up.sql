ALTER TABLE addresses
    ADD COLUMN type VARCHAR(32) NOT NULL DEFAULT 'shipping',
    ADD COLUMN is_default_billing BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_default_shipping BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE addresses DROP CONSTRAINT IF EXISTS addresses_type_check;
ALTER TABLE addresses ADD CONSTRAINT addresses_type_check CHECK (type IN ('billing', 'shipping'));

CREATE INDEX idx_addresses_user_id ON addresses(user_id);
