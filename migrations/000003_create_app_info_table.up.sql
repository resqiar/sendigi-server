CREATE TABLE IF NOT EXISTS app_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    package_name TEXT UNIQUE NOT NULL,
    lock_status BOOLEAN DEFAULT false,
    icon TEXT,
    time_usage BIGINT default 0,
    author_id UUID NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(id)
);
