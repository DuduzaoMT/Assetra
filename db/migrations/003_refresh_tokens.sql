-- 003_refresh_tokens.sql
-- Refresh tokens table for secure authentication

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE
);

-- Index for fast lookup by token hash
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash) WHERE revoked = FALSE;

-- Index for fast lookup by user_id
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- Index for cleanup of expired tokens
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
