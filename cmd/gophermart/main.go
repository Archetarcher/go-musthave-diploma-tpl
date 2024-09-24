package main

import (
	"context"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/api/rest"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/handlers"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/provider"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/repositories"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/services"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
	"go.uber.org/zap"
	"log"
	"sync"
)

func main() {

	c := config.NewConfig()
	c.ParseConfig()

	if err := logger.Initialize("info"); err != nil {
		log.Fatal("failed to init logger")
	}

	ctx := context.Background()

	storage, err := pgsql.NewStore(ctx, &pgsql.Config{DatabaseURI: c.DatabaseURI, MigrationsPath: c.MigrationsPath})
	if err != nil {
		log.Fatal("failed to init storage with error", err)

		return
	}

	app := rest.NewApplication(c, storage)

	userRepository := repositories.NewPGUserRepository(storage)
	orderAccrualRepository := repositories.NewPGOrderAccrualRepository(storage)
	orderWithdrawalRepository := repositories.NewPGOrderWithdrawalRepository(storage)

	orderService := services.NewOrderService(orderAccrualRepository, orderWithdrawalRepository, userRepository)
	authService := services.NewAuthService(userRepository, c.Token)

	orderHandler := handlers.NewOrderHandler(orderService)
	authHandler := handlers.NewAuthHandler(authService)

	app.MountMiddleware()
	app.MountHandlers(orderHandler, authHandler)

	go func() {
		//start workers
		orders := make(chan domain.OrderAccrual, c.Worker.Count*2)

		var wg sync.WaitGroup
		wg.Add(1)

		w := provider.CreateNewAccrualProvider(userRepository, orderAccrualRepository, c)
		w.CreateWorkers(ctx, orders)
		err = w.Process(ctx, &wg, orders)
		if err != nil {
			logger.Log.Error("accrual provider failed with error", zap.Error(err))
		}
	}()

	logger.Log.Info("starting server")
	if err := app.Run(c); err != nil {
		logger.Log.Error("application failed with error", zap.Error(err))
		return
	}

}
