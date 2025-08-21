package types

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Avatar    *string   `json:"avatar,omitempty"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comic struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	CoverURL    string    `json:"coverUrl"`
	Status      string    `json:"status" validate:"oneof=ongoing completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Chapter struct {
	ID            int64     `json:"id"`
	ComicID       int64     `json:"comicId"`
	Title         string    `json:"title"`
	ChapterNumber int       `json:"chapterNumber"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Page struct {
	ID         int64     `json:"id"`
	ChapterID  int64     `json:"chapterId"`
	PageNumber int       `json:"pageNumber"`
	ImageURL   string    `json:"imageUrl"`
	CreatedAt  time.Time `json:"createdAt"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByIdentifier(identifier string) (*User, error)
	GetUserByID(id int64) (*User, error)
	SaveRefreshToken(userID int64, token string, expiresAt time.Time) error
	GetRefreshToken(token string) (*RefreshToken, error)
	DeleteRefreshToken(token string) error
}

type ComicStore interface {
	CreateComic(comic *Comic) error
	GetAllComics() ([]*Comic, error)
	GetComicByID(id int64) (*Comic, error)
}

type ChapterStore interface {
	CreateChapter(comicID int64, ch *Chapter) error
	GetChaptersByComic(comicID int64) ([]*Chapter, error)
}

type PageStore interface {
	CreatePage(chapterID int64, p *Page)
	errorGetPagesByChapter(chapterID int64) ([]*Page, error)
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

type CreateComicPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
	CoverURL    string `json:"coverUrl"`
	Status      string `json:"status" validate:"oneof=ongoing completed"`
}

type CreateChapterPayload struct {
	Title         string `json:"title"`
	ChapterNumber int    `json:"chapterNumber" validate:"required"`
}

type CreatePagePayload struct {
	PageNumber int `json:"pageNumber" validate:"required"`
}
