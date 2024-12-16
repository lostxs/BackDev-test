CREATE TABLE IF NOT EXISTS sessions (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL
);

CREATE INDEX idx_user_id ON sessions (user_id);
