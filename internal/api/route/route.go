package route

import (
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountRoutes(router chi.Router, tokenConfig config.Token, orderHandler *handlers.OrderHandler, authHandler *handlers.AuthHandler) {
	router.Route("/api/user", func(r chi.Router) {
		NewAuthRoute(r, authHandler)
		NewOrderRoute(r, orderHandler, tokenConfig)
	})
}
