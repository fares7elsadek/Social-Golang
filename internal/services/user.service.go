package service

import (
	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(username, email, password string) error {
	user := &domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	return s.userRepo.CreateUser(user)
}

func (s *userService) GetUserByID(id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) UpdateUser(id int, username, email, password string) error {
	user := &domain.User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	}
	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepo.DeleteUser(id)
}