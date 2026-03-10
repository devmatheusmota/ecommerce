DROP INDEX IF EXISTS idx_addresses_user_id;
ALTER TABLE addresses DROP CONSTRAINT IF EXISTS addresses_type_check;
ALTER TABLE addresses
    DROP COLUMN IF EXISTS type,
    DROP COLUMN IF EXISTS is_default_billing,
    DROP COLUMN IF EXISTS is_default_shipping;
