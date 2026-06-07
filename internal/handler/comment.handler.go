package handler

import (
	"net/http"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
)

type commentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *commentHandler {
	return &commentHandler{commentService: commentService}
}


func (h *commentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {

}

func (h *commentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {

}

func (h *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {

}

func (h *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {

}