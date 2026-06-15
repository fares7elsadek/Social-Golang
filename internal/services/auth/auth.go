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

type AuthService interface {
	Register(ctx context.Context, email, password string, roles []domain.Role) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.TokenPair, error)
	RefreshSession(ctx context.Context, refreshToken string, userID int) (*domain.TokenPair, error)
	Logout(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
}


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
 
	roles = []domain.Role{domain.RoleUser}
 
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


func (s *Service) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// constant-time behaviour — prevents email enumeration via timing.
			_ = bcrypt.CompareHashAndPassword([]byte("$2a$10$dummy"), []byte(password))
			return nil, ErrInvalidCreds
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
 
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCreds
	}
 
	if !user.IsActive {
		return nil, ErrUserInactive
	}
 
	pair, err := s.tokens.Issue(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("issue tokens: %w", err)
	}
 
	return pair, nil
}

func (s *Service) RefreshSession(ctx context.Context, refreshToken string, userID int) (*domain.TokenPair, error) {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
 
	if !user.IsActive {
		return nil, ErrUserInactive
	}
 
	return s.tokens.Refresh(ctx, refreshToken, user)
}
 
func (s *Service) Logout(ctx context.Context, userID int) error {
	return s.tokens.RevokeAll(ctx, userID)
}

func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
 
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCreds
	}
 
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
 
	user.Password = string(hash)
	user.UpdatedAt = time.Now().UTC()
 
	if err := s.users.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
 
	return s.tokens.RevokeAll(ctx, userID)
}