CREATE TABLE IF NOT EXISTS notification_config (
    id SERIAL PRIMARY KEY,

    email TEXT,
    email_status BOOLEAN DEFAULT false,
    whatsapp TEXT,
    whatsapp_status BOOLEAN DEFAULT false,

    strategy TEXT DEFAULT 'LOCKED',
    
    user_id UUID UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
