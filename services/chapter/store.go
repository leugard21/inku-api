package chapter

import (
	"database/sql"

	"github.com/leugard21/inku-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateChapter(comicID int64, ch *types.Chapter) error {
	return s.db.QueryRow(`
        INSERT INTO chapters (comic_id, title, chapter_number, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at`,
		comicID, ch.Title, ch.ChapterNumber,
	).Scan(&ch.ID, &ch.CreatedAt, &ch.UpdatedAt)
}

func (s *Store) GetChaptersByComic(comicID int64) ([]*types.Chapter, error) {
	rows, err := s.db.Query(`
        SELECT id, comic_id, title, chapter_number, created_at, updated_at
        FROM chapters
        WHERE comic_id = $1
        ORDER BY chapter_number ASC`, comicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*types.Chapter
	for rows.Next() {
		var ch types.Chapter
		if err := rows.Scan(&ch.ID, &ch.ComicID, &ch.Title, &ch.ChapterNumber, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, err
		}
		chapters = append(chapters, &ch)
	}
	return chapters, nil
}
