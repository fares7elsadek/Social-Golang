package handler

import (
	"net/http"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) CreateUser(w http.ResponseWriter,r *http.Request) {

}

func (h *userHandler) GetUserByID(w http.ResponseWriter,r *http.Request) {

}

func (h *userHandler) GetUserByEmail(w http.ResponseWriter,r *http.Request) {

}

func (h *userHandler) UpdateUser(w http.ResponseWriter,r *http.Request) {

}

func (h *userHandler) DeleteUser(w http.ResponseWriter,r *http.Request) {

}

