CREATE TABLE IF NOT EXISTS favorites (
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    comic_id BIGINT NOT NULL REFERENCES comics (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, comic_id)
);