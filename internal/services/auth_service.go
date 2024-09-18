package services

import (
	"context"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthService struct {
	repo        UserRepository
	tokenConfig config.Token
}
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	GetUserByID(ctx context.Context, user int) (*domain.User, error)
	UpdateUserBalance(ctx context.Context, user domain.User) (*domain.User, error)
}

func NewAuthService(repo UserRepository, tokenConfig config.Token) *AuthService {
	return &AuthService{repo: repo, tokenConfig: tokenConfig}
}

func (s *AuthService) Register(ctx context.Context, request *domain.AuthRequest) (*domain.AuthResponse, *handlers.RestError) {

	user, err := s.repo.GetUserByLogin(ctx, request.Login)
	if err != nil {
		return nil, &handlers.RestError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	if user != nil {
		return nil, &handlers.RestError{Code: http.StatusConflict, Message: "user already exists with this login", Err: err}
	}

	hash, err := getPasswordHash(request.Password)
	if err != nil {
		return nil, &handlers.RestError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	user, err = s.repo.Create(ctx, domain.User{
		Login: request.Login,
		Hash:  hash,
	})

	if err != nil {
		return nil, &handlers.RestError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	token, err := util.CreateToken(user, s.tokenConfig)
	if err != nil {
		return nil, &handlers.RestError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	return &domain.AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.tokenConfig.ExpiresInMinutes)).Format(time.RFC3339),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, request *domain.AuthRequest) (*domain.AuthResponse, *handlers.RestError) {
	user, err := s.repo.GetUserByLogin(ctx, request.Login)
	if err != nil {
		return nil, &handlers.RestError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	if user == nil {
		return nil, &handlers.RestError{Code: http.StatusUnauthorized, Message: "bad credentials user not found", Err: err}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(request.Password)) != nil {
		return nil, &handlers.RestError{Code: http.StatusUnauthorized, Message: "bad credentials, login password pare are not valid", Err: err}
	}

	token, err := util.CreateToken(user, s.tokenConfig)
	if err != nil {
		return nil, &handlers.RestError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	return &domain.AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.tokenConfig.ExpiresInMinutes)).Format(time.RFC3339),
	}, nil

}

func getPasswordHash(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
