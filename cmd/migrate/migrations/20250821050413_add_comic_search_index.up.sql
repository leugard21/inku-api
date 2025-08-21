CREATE INDEX IF NOT EXISTS idx_comics_title_tsv ON comics USING gin (to_tsvector('english', title));

CREATE INDEX IF NOT EXISTS idx_comics_author_tsv ON comics USING gin (
    to_tsvector('english', author)
);

CREATE INDEX IF NOT EXISTS idx_comics_description_tsv ON comics USING gin (
    to_tsvector('english', description)
);