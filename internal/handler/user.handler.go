package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
)

type userHandler struct {
    userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
    return &userHandler{userService: userService}
}


func (h *userHandler) GetUserByID(w http.ResponseWriter,r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user id")
        return
    }

	user, err := h.userService.GetUserByID(r.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

    writeJSON(w, http.StatusOK, user)
}

func (h *userHandler) GetUserByEmail(w http.ResponseWriter,r *http.Request) {
	
	var req struct {
		Email string `json:"email"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	user, err := h.userService.GetUserByEmail(r.Context(), req.Email)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "unexpected error")
        return
    }

    writeJSON(w, http.StatusOK, user)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter,r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user id")
        return
    }

	var params domain.UpdateUserParams
    if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	if err := h.userService.UpdateUser(r.Context(), id, params); err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
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

func (h *userHandler) DeleteUser(w http.ResponseWriter,r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user id")
        return
    }

	if err := h.userService.DeleteUser(r.Context(),id) ; err != nil {
		if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
		writeError(w, http.StatusInternalServerError, "unexpected error")
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}

