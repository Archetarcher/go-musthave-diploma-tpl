package route

import (
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewAuthRoute(router chi.Router, handler *handlers.AuthHandler) {
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)

}
