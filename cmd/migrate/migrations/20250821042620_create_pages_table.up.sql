CREATE TABLE IF NOT EXISTS pages (
    id BIGSERIAL PRIMARY KEY,
    chapter_id BIGINT NOT NULL REFERENCES chapters (id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (chapter_id, page_number)
);