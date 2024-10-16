package handlers

import (
	"context"
	"encoding/json"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type AuthHandler struct {
	service AuthService
}

type AuthService interface {
	Register(ctx context.Context, user *domain.AuthRequest) (*domain.AuthResponse, *domain.Error)
	Login(ctx context.Context, user *domain.AuthRequest) (*domain.AuthResponse, *domain.Error)
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	user, err := validateAuthRequest(request)
	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	response, sErr := h.service.Register(request.Context(), user)
	if sErr != nil {
		sendResponse(enc, sErr, sErr.Code, writer)
		return
	}

	writer.Header().Set("Authorization", spew.Sprintf("Bearer %s", response.Token))
	sendResponse(enc, response, http.StatusOK, writer)

}

func (h *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")

	user, err := validateAuthRequest(request)
	if err != nil {
		sendResponse(enc, err, err.Code, writer)
		return
	}

	response, sErr := h.service.Login(request.Context(), user)
	if sErr != nil {
		sendResponse(enc, sErr, sErr.Code, writer)
		return
	}
	writer.Header().Set("Authorization", spew.Sprintf("Bearer %s", response.Token))
	sendResponse(enc, response, http.StatusOK, writer)

}

func validateAuthRequest(request *http.Request) (*domain.AuthRequest, *RestError) {
	var u domain.AuthRequest

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}
	logger.Log.Info("user info", zap.Any("user", u))

	v := validator.New()
	err = v.Struct(u)

	if err != nil {
		return nil, &RestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}
	}

	return &u, nil
}
