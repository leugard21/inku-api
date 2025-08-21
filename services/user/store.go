package user

import (
	"database/sql"
	"errors"
	"time"

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
		RETURNING id, role, created_at`,
		user.Username, user.Avatar, user.Email, user.Password,
	).Scan(&user.ID, &user.Role, &user.CreatedAt)
}

func (s *Store) GetUserByIdentifier(identifier string) (*types.User, error) {
	row := s.db.QueryRow(`
	SELECT id, username, avatar, email, password, role, created_at
	FROM users
	WHERE email = $1 OR username = $1`, identifier)

	var u types.User
	if err := row.Scan(&u.ID, &u.Username, &u.Avatar, &u.Email, &u.Password, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (s *Store) GetUserByID(id int64) (*types.User, error) {
	row := s.db.QueryRow(`
		SELECT id, username, avatar, email, password, role, created_at
		FROM users
		WHERE id = $1`, id,
	)

	var u types.User
	if err := row.Scan(&u.ID, &u.Username, &u.Avatar, &u.Email, &u.Password, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *Store) SaveRefreshToken(userID int64, token string, expiresAt time.Time) error {
	_, err := s.db.Exec(`
		INSERT INTO refresh_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)`,
		token, userID, expiresAt,
	)
	return err
}

func (s *Store) GetRefreshToken(token string) (*types.RefreshToken, error) {
	row := s.db.QueryRow(`
		SELECT token, user_id, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1`, token,
	)
	var rt types.RefreshToken
	if err := row.Scan(&rt.Token, &rt.UserID, &rt.ExpiresAt, &rt.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (s *Store) DeleteRefreshToken(token string) error {
	_, err := s.db.Exec(`DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}
