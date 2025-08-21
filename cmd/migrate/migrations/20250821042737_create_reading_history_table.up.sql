CREATE TABLE IF NOT EXISTS reading_history (
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    comic_id BIGINT NOT NULL REFERENCES comics (id) ON DELETE CASCADE,
    chapter_id BIGINT NOT NULL REFERENCES chapters (id) ON DELETE CASCADE,
    page_number INT DEFAULT 1,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, comic_id)
);