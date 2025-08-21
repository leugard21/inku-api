CREATE TABLE IF NOT EXISTS chapters (
    id BIGSERIAL PRIMARY KEY,
    comic_id BIGINT NOT NULL REFERENCES comics (id) ON DELETE CASCADE,
    title TEXT,
    chapter_number INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (comic_id, chapter_number)
);