package domain

import (
	"time"
)

type User struct {
	ID        int      "json:\"id\""
	Username  string   "json:\"username\""
	Email     string   "json:\"email\""
	CreatedAt time.Time "json:\"created_at\""
	UpdatedAt time.Time "json:\"updated_at\""
	Password  string   "json:\"-\""
}

type Post struct {
	ID        int      "json:\"id\""
	Title     string   "json:\"title\""
	Content   string   "json:\"content\""
	AuthorID  int      "json:\"author_id\""
	CreatedAt time.Time "json:\"created_at\""
	UpdatedAt time.Time "json:\"updated_at\""
}

type Comment struct {
	ID        int      "json:\"id\""
	Content   string   "json:\"content\""
	AuthorID  int      "json:\"author_id\""
	PostID    int      "json:\"post_id\""
	CreatedAt time.Time "json:\"created_at\""
	UpdatedAt time.Time "json:\"updated_at\""
}

