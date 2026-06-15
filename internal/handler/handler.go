package handler

import (
	"net/http"
)

type UserHandler interface {
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

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Me(w http.ResponseWriter, r *http.Request)
}
