package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	middleware "github.com/fares7elsadek/Social-Golang/internal/middlewares"
	"github.com/fares7elsadek/Social-Golang/internal/services/auth"
	"github.com/fares7elsadek/Social-Golang/internal/services/token"
)

type authHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

// POST /auth/register

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, nil)
	if err != nil {
		if errors.Is(err, auth.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

// POST /auth/login

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pair, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCreds):
			writeError(w, http.StatusUnauthorized, "invalid email or password")
		case errors.Is(err, auth.ErrUserInactive):
			writeError(w, http.StatusForbidden, "account is disabled")
		default:
			writeError(w, http.StatusInternalServerError, "login failed")
		}
		return
	}

	writeJSON(w, http.StatusOK, pair)
}

// POST /auth/refresh

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userIdNumber, err := strconv.Atoi(req.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "unexpected error")
	}

	pair, err := h.authService.RefreshSession(r.Context(), req.RefreshToken, userIdNumber)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrTokenExpired):
			writeError(w, http.StatusUnauthorized, "refresh token expired")
		case errors.Is(err, token.ErrTokenRevoked):
			writeError(w, http.StatusUnauthorized, "refresh token revoked")
		default:
			writeError(w, http.StatusUnauthorized, "invalid refresh token")
		}
		return
	}

	writeJSON(w, http.StatusOK, pair)
}

// POST /auth/logout

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	userIdNumber, err := strconv.Atoi(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "unexpected error")
	}

	if err := h.authService.Logout(r.Context(), userIdNumber); err != nil {
		writeError(w, http.StatusInternalServerError, "logout failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /me  (protected)

func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromCtx(r.Context())
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"roles":   claims.Roles,
	})
}