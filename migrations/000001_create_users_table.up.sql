CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    provider   VARCHAR(20),
    username   VARCHAR(100) UNIQUE NOT NULL,
    email   VARCHAR(100) UNIQUE NOT NULL,
    bio        TEXT,
    picture_url TEXT
);
