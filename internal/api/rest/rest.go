package rest

import (
	"context"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/api/route"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/provider"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/repositories"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/services"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"sync"
)

type Application struct {
	router   *chi.Mux
	config   *config.AppConfig
	pgxStore *pgsql.Store
}

func NewApplication(appConfig *config.AppConfig, store *pgsql.Store) *Application {
	util.GenerateAuthToken(appConfig)
	return &Application{
		router:   chi.NewRouter(),
		pgxStore: store,
		config:   appConfig,
	}
}

func (app *Application) MountMiddleware() {
	app.router.Use(middleware.Logger)
}

func (app *Application) MountHandlers(ctx context.Context) {
	userRepository := repositories.NewPGUserRepository(app.pgxStore)
	orderAccrualRepository := repositories.NewPGOrderAccrualRepository(app.pgxStore)
	orderWithdrawalRepository := repositories.NewPGOrderWithdrawalRepository(app.pgxStore)

	orderService := services.NewOrderService(orderAccrualRepository, orderWithdrawalRepository, userRepository)
	authService := services.NewAuthService(userRepository, app.config.Token)

	orderHandler := handlers.NewOrderHandler(orderService)
	authHandler := handlers.NewAuthHandler(authService)

	route.MountRoutes(app.router, app.config.Token, orderHandler, authHandler)

	// start workers
	orders := make(chan domain.OrderAccrual, app.config.Worker.Count*2)

	var wg sync.WaitGroup
	wg.Add(1)

	w := provider.CreateNewAccrualProvider(userRepository, orderAccrualRepository, app.config)
	w.CreateWorkers(ctx, orders)
	w.Process(ctx, &wg, orders)
}

func (app *Application) Run(config *config.AppConfig) error {

	err := http.ListenAndServe(config.RunAddr, app.router)
	if err != nil {
		return err
	}
	return nil
}
