package service

import "github.com/fares7elsadek/Social-Golang/internal/domain"

type UserService interface {
	CreateUser(username, email, password string) error
	GetUserByID(id int) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(id int, username, email, password string) error
	DeleteUser(id int) error
}

type PostService interface {
	CreatePost(userID int, content string) error
	GetPostByID(id int) (*domain.Post, error)
	UpdatePost(id int, content string) error
	DeletePost(id int) error
}

type CommentService interface {
	CreateComment(userID, postID int, content string) error
	GetCommentByID(id int) (*domain.Comment, error)
	UpdateComment(id int, content string) error
	DeleteComment(id int) error
}