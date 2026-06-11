package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
)

type postHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *postHandler {
	return &postHandler{postService: postService}
}

func (h *postHandler) CreatePost(w http.ResponseWriter,r *http.Request) {
	authorIdStr := r.PathValue("authorId")
	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user id")
        return
    }

	var req struct {
        Title    string `json:"title"`
        Content string `json:"content"`
    }

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	if req.Title == "" || req.Content == ""  {
        writeError(w, http.StatusBadRequest, "title and content are required")
        return
    }

	if err := h.postService.CreatePost(r.Context(), authorId, req.Title, req.Content); err != nil {
        if errors.Is(err, domain.ErrConflict) {
            writeError(w, http.StatusConflict, err.Error())
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

	writeJSON(w, http.StatusCreated, map[string]string{"message": "post created"})
}

func (h *postHandler) GetPostByID(w http.ResponseWriter,r *http.Request) {
	idStr := r.PathValue("postId")
	id,err := strconv.Atoi(idStr)
	if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user id")
        return
    }

	post ,err := h.postService.GetPostByID(r.Context(),id)

	if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "post not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

	writeJSON(w, http.StatusOK, post)
}

func (h *postHandler) GetPostsByAuthorId(
    w http.ResponseWriter,
    r *http.Request,
) {
    idStr := r.PathValue("authorId")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid author id")
        return
    }

    // Default values
    limit := 10
    offset := 0

    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        limit, err = strconv.Atoi(limitStr)
        if err != nil || limit <= 0 {
            writeError(w, http.StatusBadRequest, "invalid limit")
            return
        }
    }

    if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
        offset, err = strconv.Atoi(offsetStr)
        if err != nil || offset < 0 {
            writeError(w, http.StatusBadRequest, "invalid offset")
            return
        }
    }

    posts, err := h.postService.GetPostsByAuthorID(
        r.Context(),
        id,
        limit,
        offset,
    )

    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "posts not found")
            return
        }

        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

    writeJSON(w, http.StatusOK, posts)
}

func (h *postHandler) UpdatePost(w http.ResponseWriter,r *http.Request) {
	idStr := r.PathValue("postId")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid post id")
        return
    }

	var params domain.UpdatePostParams
    if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	if err := h.postService.UpdatePost(r.Context(), id, params); err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "post not found")
            return
        }
        if errors.Is(err, domain.ErrConflict) {
            writeError(w, http.StatusConflict, err.Error())
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}	

func (h *postHandler) DeletePost(w http.ResponseWriter,r *http.Request) {
	idStr := r.PathValue("postId")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid post id")
        return
    }

	if err := h.postService.DeletePost(r.Context(),id) ; err != nil {
		if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "post not found")
            return
        }
		writeError(w, http.StatusInternalServerError, "unexpected error")
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
}
