package user

import (
	"database/sql"
	"errors"

	"github.com/leugard21/inku-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user *types.User) error {
	return s.db.QueryRow(`
		INSERT INTO users (username, avatar, email, password, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at`,
		user.Username, user.Avatar, user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
}

func (s *Store) GetUserByIdentifier(identifier string) (*types.User, error) {
	row := s.db.QueryRow(`
	SELECT id, username, avatar, email, password, created_at
	FROM users
	WHERE email = $1 OR username = $1`, identifier)

	var u types.User
	if err := row.Scan(&u.ID, &u.Username, &u.Avatar, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
