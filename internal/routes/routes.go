package routes

import (
	"net/http"

	"github.com/bingKegeta/Knight-Link/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/users", UserRoutes())
		r.Mount("/api/auth", AuthRoutes())
		r.Mount("/api/events", EventRoutes())
		r.Mount("/api/rsos", RSORoutes())
	})

	return router
}

func UserRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/{userId}", handlers.GetUser)
	router.Delete("/{userId}", handlers.DeleteUser)
	router.Put("/{userId}", handlers.UpdateUser)
	router.Post("/", handlers.CreateUser)
	return router
}

func AuthRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/login", handlers.Login)
	router.Post("/logout", handlers.Logout)
	return router
}

func EventRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllEvents)
	router.Get("/{eventId}", handlers.GetEvent)
	router.Delete("/{eventId}", handlers.DeleteEvent)
	router.Put("/{eventId}", handlers.UpdateEvent)
	router.Post("/", handlers.CreateEvent)
	return router
}

func RSORoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllRSOs)
	router.Get("/{rsoId}", handlers.GetRSO)
	router.Delete("/{rsoId}", handlers.DeleteRSO)
	router.Put("/{rsoId}", handlers.UpdateRSO)
	router.Post("/", handlers.CreateRSO)
	return router
}
