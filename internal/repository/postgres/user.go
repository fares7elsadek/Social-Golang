package postgres

import (
    "github.com/fares7elsadek/Social-Golang/internal/domain"
)

type userRepository struct {
    
}

func NewUserRepository() *userRepository {
    return &userRepository{}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	
    return nil
}

func (r *userRepository) GetUserByID(id int) (*domain.User, error) {
    return nil, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
    return nil, nil
}

func (r *userRepository) UpdateUser(user *domain.User) error {
    return nil
}

func (r *userRepository) DeleteUser(id int) error {
    return nil
}