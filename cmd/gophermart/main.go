package main

import (
	"context"
	"fmt"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/api/rest"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/store/pgsql"
)

func main() {

	c := config.NewConfig()

	ctx := context.Background()

	storage, err := pgsql.NewStore(ctx, &pgsql.Config{DatabaseUri: c.DatabaseUri, MigrationsPath: c.MigrationsPath})
	if err != nil {
		fmt.Println(err)
		return
	}

	app := rest.NewApplication(c, storage)

	app.MountMiddleware()
	app.MountHandlers(ctx)

	fmt.Println("started server")

	if err := app.Run(c); err != nil {
		fmt.Println(err)
		return
	}
}
