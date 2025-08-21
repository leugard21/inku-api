CREATE TABLE IF NOT EXISTS comicS (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    author TEXT,
    cover_url TEXT,
    status TEXT DEFAULT 'ongoing',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);