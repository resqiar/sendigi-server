CREATE TABLE IF NOT EXISTS app_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    package_name TEXT NOT NULL,
    lock_status BOOLEAN DEFAULT false,
    icon TEXT,
    time_usage BIGINT default 0,
    author_id UUID NOT NULL,
    device_id TEXT NOT NULL,

    date_locked TEXT,
    time_start_locked TEXT,
    time_end_locked TEXT,
    recurring TEXT,

    FOREIGN KEY (author_id) REFERENCES users(id),
    FOREIGN KEY (device_id) REFERENCES device_info(id)
);
