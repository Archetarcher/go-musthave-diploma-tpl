package route

import (
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func NewOrderRoute(router chi.Router, handler *handlers.OrderHandler, token config.Token) {
	router.Group(func(userRouter chi.Router) {
		userRouter.Use(jwtauth.Verifier(token.AuthToken))
		userRouter.Use(jwtauth.Authenticator(token.AuthToken))

		userRouter.Post("/orders", handler.RegisterAccrualOrder)
		userRouter.Get("/orders", handler.GetAllAccrual)

		userRouter.Get("/balance", handler.GetUserBalance)
		userRouter.Post("/balance/withdraw", handler.RegisterWithdrawalOrder)
		userRouter.Get("/withdrawals", handler.GetAllWithdrawal)

	})
}
