ALTER TABLE users ADD COLUMN phone_verified BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE phone_verification_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(10) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_phone_verification_codes_user_id ON phone_verification_codes(user_id);
CREATE INDEX idx_phone_verification_codes_expires_at ON phone_verification_codes(expires_at);
