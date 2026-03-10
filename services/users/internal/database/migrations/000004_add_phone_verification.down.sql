DROP INDEX IF EXISTS idx_phone_verification_codes_expires_at;
DROP INDEX IF EXISTS idx_phone_verification_codes_user_id;
DROP TABLE IF EXISTS phone_verification_codes;
ALTER TABLE users DROP COLUMN IF EXISTS phone_verified;
