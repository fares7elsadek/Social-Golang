package handler

import (
	"net/http"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter,r *http.Request)
	GetUserByID(w http.ResponseWriter,r *http.Request)
	GetUserByEmail(w http.ResponseWriter,r *http.Request)
	UpdateUser(w http.ResponseWriter,r *http.Request)
	DeleteUser(w http.ResponseWriter,r *http.Request)
}

type PostHandler interface {
	CreatePost(w http.ResponseWriter,r *http.Request)
	GetPostByID(w http.ResponseWriter,r *http.Request)
	GetPostsByAuthorId(w http.ResponseWriter,r *http.Request)
	UpdatePost(w http.ResponseWriter,r *http.Request)
	DeletePost(w http.ResponseWriter,r *http.Request)
}


type CommentHandler interface {
	CreateComment(w http.ResponseWriter,r *http.Request)
	GetCommentByID(w http.ResponseWriter,r *http.Request)
	GetCommentsByPostId(w http.ResponseWriter,r *http.Request)
	UpdateComment(w http.ResponseWriter,r *http.Request)
	DeleteComment(w http.ResponseWriter,r *http.Request)
}
