package types

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Avatar    *string   `json:"avatar,omitempty"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByIdentifier(identifier string) (*User, error)
}

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginPayload struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
