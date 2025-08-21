package comic

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

func (s *Store) CreateComic(comic *types.Comic) error {
	return s.db.QueryRow(`
        INSERT INTO comics (title, description, author, cover_url, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at`,
		comic.Title, comic.Description, comic.Author, comic.CoverURL, comic.Status,
	).Scan(&comic.ID, &comic.CreatedAt, &comic.UpdatedAt)
}

func (s *Store) GetAllComics() ([]*types.Comic, error) {
	rows, err := s.db.Query(`
		SELECT id, title, description, author, cover_url, status, created_at, updated_at
		FROM comics
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comics []*types.Comic
	for rows.Next() {
		var c types.Comic
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.CoverURL, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comics = append(comics, &c)
	}
	return comics, nil
}

func (s *Store) GetComicByID(id int64) (*types.Comic, error) {
	row := s.db.QueryRow(`
		SELECT id, title, description, author, cover_url, status, created_at, updated_at
		FROM comics WHERE id = $1`, id)

	var c types.Comic
	if err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.CoverURL, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
