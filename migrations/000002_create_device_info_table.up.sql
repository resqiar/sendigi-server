CREATE TABLE IF NOT EXISTS device_info (
    id TEXT UNIQUE PRIMARY KEY,
    device_name TEXT,
    device_brand TEXT,
    api_level INTEGER,
    android_version TEXT,
    manufacturer TEXT,
    product_name TEXT,
    battery_level INTEGER DEFAULT 0,
    is_charging BOOLEAN DEFAULT false,
    author_id UUID NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(id)
);
