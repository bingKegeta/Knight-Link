package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"

	"github.com/bingKegeta/Knight-Link/internal/routes"
)

type App struct {
	router http.Handler
}

var tokenAuth *jwtauth.JWTAuth

func init() {
	// Initialization logic...
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)
}

func New() *App {
	app := &App{
		router: routes.Routes(tokenAuth),
	}
	fmt.Println("Server started.")

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
