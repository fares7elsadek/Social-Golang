package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
	ErrTokenRevoked = errors.New("token revoked")
)


type RefreshStore interface {
	Save(ctx context.Context, userID, tokenID string, ttl time.Duration) error
	Exists(ctx context.Context, userID, tokenID string) (bool, error)
	Delete(ctx context.Context, userID, tokenID string) error
	DeleteAll(ctx context.Context, userID string) error
}

type jwtClaims struct {
	UserID string       `json:"sub"`
	Email  string       `json:"email"`
	Roles  []domain.Role `json:"roles"`
	JTI    string       `json:"jti"`
	jwt.RegisteredClaims
}

type Config struct {
	AccessSecret       []byte
	RefreshSecret      []byte
	AccessTokenTTL     time.Duration // e.g. 15 * time.Minute
	RefreshTokenTTL    time.Duration // e.g. 7 * 24 * time.Hour
}

func DefaultConfig(accessSecret, refreshSecret []byte) Config {
	return Config{
		AccessSecret:    accessSecret,
		RefreshSecret:   refreshSecret,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
}

type Service struct {
	cfg   Config
	store RefreshStore
}

func New(cfg Config, store RefreshStore) *Service {
	return &Service{cfg: cfg, store: store}
}


func (s *Service) Issue(ctx context.Context, user *domain.User) (*domain.TokenPair , error){
	accessToken , err := s.makeAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("make access token: %w", err)
	}

	refreshToken, tokenID, err := s.makeRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("make refresh token: %w", err)
	}
	
	if err := s.store.Save(ctx, strconv.Itoa(user.ID), tokenID, s.cfg.RefreshTokenTTL); err != nil {
		return nil, fmt.Errorf("persist refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *Service) VerifyAccess(tokenStr string) (*domain.Claims , error){
	claims := &jwtClaims{}
	_,err := jwt.ParseWithClaims(tokenStr,claims,func(t *jwt.Token) (any, error) {
		if _,ok := t.Method.(*jwt.SigningMethodHMAC) ; !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.cfg.AccessSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	return &domain.Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Roles:  claims.Roles,
	}, nil
}


func (s *Service) Refresh(ctx context.Context, refreshTokenStr string, user *domain.User) (*domain.TokenPair, error){
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.cfg.RefreshSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	ok, err := s.store.Exists(ctx, strconv.Itoa(user.ID), claims.JTI)

	if err != nil {
		return nil, fmt.Errorf("check refresh store: %w", err)
	}

	if !ok {
		_ = s.store.DeleteAll(ctx, strconv.Itoa(user.ID))
		return nil, ErrTokenRevoked
	}

	if err := s.store.Delete(ctx, strconv.Itoa(user.ID), claims.JTI); err != nil {
		return nil, fmt.Errorf("rotate refresh token: %w", err)
	}

	return s.Issue(ctx, user)
}



func (s *Service) makeAccessToken(user *domain.User) (string , error){
	now := time.Now()
	claims := jwtClaims{
		UserID: strconv.Itoa(user.ID),
		Email: user.Email,
		Roles: user.Roles,
		JTI: newTokenID(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
			Issuer: "go-auth",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodES256,claims).SignedString(s.cfg.AccessSecret)
}

func (s *Service) RevokeAll(ctx context.Context, userID string) error {
	return s.store.DeleteAll(ctx, userID)
}


func (s *Service) makeRefreshToken(user *domain.User) (tokenStr , tokenId string ,err error){
	tokenId = newTokenID()
	now := time.Now()
	claims := jwtClaims{
		UserID: strconv.Itoa(user.ID),
		JTI: newTokenID(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
			Issuer: "go-auth",
		},
	}
	tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodES256,claims).SignedString(s.cfg.AccessSecret)
	return 
}

func newTokenID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}