package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
	"github.com/fares7elsadek/Social-Golang/internal/services/token"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken   = errors.New("email already registered")
	ErrInvalidCreds = errors.New("invalid email or password")
	ErrUserNotFound = errors.New("user not found")
	ErrUserInactive = errors.New("account is disabled")
)


type Service struct {
	users  repository.UserRepository
	tokens *token.Service
}

func NewAuthService(users repository.UserRepository, tokens *token.Service) *Service {
	return &Service{users: users, tokens: tokens}
}

func (s *Service) Register(ctx context.Context, email, password string, roles []domain.Role) (*domain.User, error) {
	existing, err := s.users.GetUserByEmail(ctx, email)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	if existing != nil {
		return nil, ErrEmailTaken
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
 
	if len(roles) == 0 {
		roles = []domain.Role{domain.RoleUser}
	}
 
	user := &domain.User{
		Email:        email,
		Password: string(hash),
		Roles:        roles,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
 
	if err := s.users.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
 
	return user, nil
}