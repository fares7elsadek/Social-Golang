package service

import (
	"context"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context,username, email, password string) error
	GetUserByID(ctx context.Context,id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context,email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context,id int,userParams domain.UpdateUserParams) error
	DeleteUser(ctx context.Context,id int) error
}

type PostService interface {
	CreatePost(ctx context.Context,userID int,title string ,content string) error
	GetPostByID(ctx context.Context,id int) (*domain.Post, error)
	GetPostsByAuthorID(ctx context.Context,authorId,limit,offset int) ([]*domain.Post, error)
	UpdatePost(ctx context.Context,id int, postParams domain.UpdatePostParams) error
	DeletePost(ctx context.Context,id int) error
}

type CommentService interface {
	CreateComment(ctx context.Context,userID, postID int, content string) error
	GetCommentByID(ctx context.Context,id int) (*domain.Comment, error)
	GetCommentByPostID(ctx context.Context,postId,limit,offset int) ([]*domain.Comment, error)
	UpdateComment(ctx context.Context,id int, content string) error
	DeleteComment(ctx context.Context,id int) error
}