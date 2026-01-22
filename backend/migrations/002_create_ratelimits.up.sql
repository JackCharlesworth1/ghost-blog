CREATE TABLE IF NOT EXISTS ratelimits (
    ip_address VARCHAR(45) PRIMARY KEY,
    post_count INTEGER DEFAULT 0,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ratelimits_window_start ON ratelimits(window_start);
