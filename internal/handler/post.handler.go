package handler

import (
	"net/http"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
)

type postHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *postHandler {
	return &postHandler{postService: postService}
}

func (h *postHandler) CreatePost(w http.ResponseWriter,r *http.Request) {

}

func (h *postHandler) GetPostByID(w http.ResponseWriter,r *http.Request) {

}

func (h *postHandler) UpdatePost(w http.ResponseWriter,r *http.Request) {

}

func (h *postHandler) DeletePost(w http.ResponseWriter,r *http.Request) {

}
