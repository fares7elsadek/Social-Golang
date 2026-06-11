package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
)

type commentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *commentHandler {
	return &commentHandler{commentService: commentService}
}


func (h *commentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.PathValue("postId")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
        writeError(w, http.StatusBadRequest, "invalid post id")
        return
    }

	authorIdStr := r.PathValue("authorId")
	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
        writeError(w, http.StatusBadRequest, "invalid author id")
        return
    }

	var req struct {
        Content string `json:"content"`
    }

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	if req.Content == ""   {
        writeError(w, http.StatusBadRequest, "content is required")
        return
    }

	if err := h.commentService.CreateComment(r.Context(),authorId ,postId, req.Content); err != nil {
        if errors.Is(err, domain.ErrConflict) {
            writeError(w, http.StatusConflict, err.Error())
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

	writeJSON(w, http.StatusCreated, map[string]string{"message": "comment created"})
}

func (h *commentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("commentId")
	id,err := strconv.Atoi(idStr)
	if err != nil {
        writeError(w, http.StatusBadRequest, "invalid comment id")
        return
    }

	comment ,err := h.commentService.GetCommentByID(r.Context(),id)

	if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "comment not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

	writeJSON(w, http.StatusOK, comment)
}

func (h *commentHandler) GetCommentsByPostId(w http.ResponseWriter,r *http.Request){
	idStr := r.PathValue("postId")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid postId id")
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

    comments, err := h.commentService.GetCommentByPostID(
        r.Context(),
        id,
        limit,
        offset,
    )

    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "comments not found")
            return
        }

        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

    writeJSON(w, http.StatusOK, comments)
}

func (h *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("commentId")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid commentId id")
        return
    }

	var req struct {
		Content string `json:"cotent"`
	}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	if err := h.commentService.UpdateComment(r.Context(), id, req.Content); err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "comment not found")
            return
        }
        if errors.Is(err, domain.ErrConflict) {
            writeError(w, http.StatusConflict, err.Error())
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"message": "comment updated"})
}

func (h *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("commentId")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid comment id")
        return
    }

	if err := h.commentService.DeleteComment(r.Context(),id) ; err != nil {
		if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "comment not found")
            return
        }
		writeError(w, http.StatusInternalServerError, "unexpected error")
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "comment deleted"})
}