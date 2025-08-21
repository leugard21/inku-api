package comic

import (
	"database/sql"
	"fmt"

	"github.com/leugard21/inku-api/types"
	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateComic(c *types.Comic) error {
	return s.db.QueryRow(`
        INSERT INTO comics (title, description, author, cover_url, status, genres, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        RETURNING id, created_at, updated_at`,
		c.Title, c.Description, c.Author, c.CoverURL, c.Status, pq.Array(c.Genres),
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (s *Store) GetAllComics() ([]*types.Comic, error) {
	rows, err := s.db.Query(`
		SELECT id, title, description, author, cover_url, status, genres, created_at, updated_at
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
	var c types.Comic
	err := s.db.QueryRow(`
        SELECT id, title, description, author, cover_url, status, genres, created_at, updated_at
        FROM comics WHERE id=$1`, id).
		Scan(&c.ID, &c.Title, &c.Description, &c.Author, &c.CoverURL, &c.Status, pq.Array(&c.Genres), &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) SearchComicsAdvanced(q, genre, status, sort string) ([]*types.Comic, error) {
	baseQuery := `
        SELECT c.id, c.title, c.description, c.author, c.cover_url, c.status, c.genres, c.created_at, c.updated_at
        FROM comics c
        WHERE 1=1
    `
	args := []interface{}{}
	idx := 1

	if q != "" {
		baseQuery += fmt.Sprintf(` AND (c.title ILIKE '%%' || $%d || '%%'
                             OR c.author ILIKE '%%' || $%d || '%%'
                             OR c.description ILIKE '%%' || $%d || '%%')`, idx, idx, idx)
		args = append(args, q)
		idx++
	}

	if genre != "" {
		baseQuery += fmt.Sprintf(" AND $%d = ANY(c.genres)", idx)
		args = append(args, genre)
		idx++
	}

	if status != "" {
		baseQuery += fmt.Sprintf(" AND c.status = $%d", idx)
		args = append(args, status)
		idx++
	}

	switch sort {
	case "newest":
		baseQuery += " ORDER BY c.created_at DESC"
	case "oldest":
		baseQuery += " ORDER BY c.created_at ASC"
	case "title":
		baseQuery += " ORDER BY c.title ASC"
	default:
		baseQuery += " ORDER BY c.created_at DESC"
	}

	baseQuery += " LIMIT 50"

	rows, err := s.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comics []*types.Comic
	for rows.Next() {
		var c types.Comic
		if err := rows.Scan(
			&c.ID, &c.Title, &c.Description, &c.Author,
			&c.CoverURL, &c.Status, pq.Array(&c.Genres),
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comics = append(comics, &c)
	}

	return comics, nil
}
