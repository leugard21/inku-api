package page

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

func (s *Store) CreatePage(chapterID int64, p *types.Page) error {
	return s.db.QueryRow(`
        INSERT INTO pages (chapter_id, page_number, image_url, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, created_at`,
		chapterID, p.PageNumber, p.ImageURL,
	).Scan(&p.ID, &p.CreatedAt)
}

func (s *Store) GetPagesByChapter(chapterID int64) ([]*types.Page, error) {
	rows, err := s.db.Query(`
        SELECT id, chapter_id, page_number, image_url, created_at
        FROM pages
        WHERE chapter_id = $1
        ORDER BY page_number ASC`, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*types.Page
	for rows.Next() {
		var p types.Page
		if err := rows.Scan(&p.ID, &p.ChapterID, &p.PageNumber, &p.ImageURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}
	return pages, nil
}
