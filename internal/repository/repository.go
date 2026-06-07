package repository

import "github.com/fares7elsadek/Social-Golang/internal/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id int) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int) error
}

type PostRepository interface {
	CreatePost(post *domain.Post) error
	GetPostByID(id int) (*domain.Post, error)
	UpdatePost(post *domain.Post) error
	DeletePost(id int) error
}

type CommentRepository interface {
	CreateComment(comment *domain.Comment) error
	GetCommentByID(id int) (*domain.Comment, error)
	UpdateComment(comment *domain.Comment) error
	DeleteComment(id int) error
}