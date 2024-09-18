package route

import (
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func NewAuthRoute(router chi.Router, handler *handlers.AuthHandler, tokenConfig config.Token) {
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)

	router.Group(func(userRouter chi.Router) {
		userRouter.Use(jwtauth.Verifier(tokenConfig.AuthToken))
		userRouter.Use(jwtauth.Authenticator(tokenConfig.AuthToken))

		userRouter.Get("/id", handler.User)
	})
}
