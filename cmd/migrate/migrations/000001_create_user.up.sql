CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email citext UNIQUE NOT NULL
);

CREATE INDEX idx_email ON users (email);
