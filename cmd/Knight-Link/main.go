package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bingKegeta/Knight-Link/internal/routes"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: routes.LoadRoutes(),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8000",
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func main() {
	app := New()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
