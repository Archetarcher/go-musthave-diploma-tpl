package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

type AuthHandler struct {
	service AuthService
}

type AuthService interface {
	Register(ctx context.Context, user *domain.AuthRequest) (*domain.AuthResponse, *RestError)
	Login(ctx context.Context, user *domain.AuthRequest) (*domain.AuthResponse, *RestError)
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

	sendResponse(enc, response, http.StatusOK, writer)

}
func (h *AuthHandler) User(writer http.ResponseWriter, request *http.Request) {
	enc := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(request.Context())

	sendResponse(enc, claims, http.StatusOK, writer)

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
	fmt.Println(u)

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
