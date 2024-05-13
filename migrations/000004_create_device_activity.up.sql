CREATE TABLE IF NOT EXISTS device_activity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    package_name TEXT NOT NULL,
    device_id  TEXT NOT NULL,
    author_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (device_id) REFERENCES device_info(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);
