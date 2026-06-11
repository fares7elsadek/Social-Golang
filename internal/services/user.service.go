package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context,username, email, password string) error {
	
	user := &domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	existingUsername,_ := s.userRepo.GetUserByUsername(ctx,username)
	
	if existingUsername != nil {
		return fmt.Errorf("username already exists: %w", domain.ErrConflict)
	}

	existingEmail,_ := s.userRepo.GetUserByEmail(ctx,email);

	if  existingEmail != nil {
		return fmt.Errorf("email already exists: %w", domain.ErrConflict)
	}

	err := s.userRepo.CreateUser(ctx,user)

	if err!=nil{
		return fmt.Errorf("unexpected error: %w",domain.ErrServer)
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context,id int) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx,id)

	if err!= nil {
		if errors.Is(err,domain.ErrNotFound) {
			return nil,fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return user,nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx,email)

	if err!= nil {
		if errors.Is(err,domain.ErrNotFound) {
			return nil,fmt.Errorf("email not found: %w", domain.ErrNotFound)
		}
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return user,nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx,username)

	if err!= nil {
		if errors.Is(err,domain.ErrNotFound) {
			return nil,fmt.Errorf("username not found: %w", domain.ErrNotFound)
		}
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return user,nil
}




func (s *userService) UpdateUser(ctx context.Context, id int, userParams domain.UpdateUserParams) error {
    user, err := s.userRepo.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return fmt.Errorf("user not found: %w", domain.ErrNotFound)
        }
        return fmt.Errorf("unexpected error: %w", domain.ErrServer)
    }

    if userParams.Username != nil {
        existingUsername, _ := s.userRepo.GetUserByUsername(ctx, *userParams.Username)
        if existingUsername != nil && existingUsername.ID != id {
            return fmt.Errorf("username already exists: %w", domain.ErrConflict)
        }
        user.Username = *userParams.Username
    }

    if userParams.Email != nil {
        existingEmail, _ := s.userRepo.GetUserByEmail(ctx, *userParams.Email)
        if existingEmail != nil && existingEmail.ID != id {
            return fmt.Errorf("email already exists: %w", domain.ErrConflict)
        }
        user.Email = *userParams.Email
    }

    err = s.userRepo.UpdateUser(ctx, user)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return fmt.Errorf("user not found: %w", domain.ErrNotFound)
        }
        return fmt.Errorf("unexpected error: %w", domain.ErrServer)
    }
    return nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {

	err := s.userRepo.DeleteUser(ctx,id)

	if err != nil {
		if errors.Is(err,domain.ErrNotFound) {
			return fmt.Errorf("user not found: %w",domain.ErrNotFound)
		}
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}