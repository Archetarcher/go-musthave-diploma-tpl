package rest

import (
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/api/route"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
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

func (app *Application) MountHandlers(orderHandler *handlers.OrderHandler, authHandler *handlers.AuthHandler) {
	route.MountRoutes(app.router, app.config.Token, orderHandler, authHandler)

}

func (app *Application) Run(config *config.AppConfig) error {
	logger.Log.Info("Running server ", zap.String("address", config.RunAddr))

	err := http.ListenAndServe(config.RunAddr, app.router)
	if err != nil {
		return err
	}
	return nil
}
