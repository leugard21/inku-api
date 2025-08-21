CREATE TABLE IF NOT EXISTS comic_genres (
    comic_id BIGINT NOT NULL REFERENCES comics (id) ON DELETE CASCADE,
    genre_id BIGINT NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (comic_id, genre_id)
);