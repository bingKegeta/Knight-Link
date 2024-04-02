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
		r.Mount("/api/unis", UniRoutes())
		r.Mount("/api/locations", LocationRoutes())
		// Add new route groups here
	})

	return router
}

func UserRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/{userId}", handlers.GetUser)
	router.Delete("/{userId}", handlers.DeleteUser)
	router.Put("/{userId}", handlers.UpdateUser)
	router.Post("/", handlers.CreateUser)

	// Add new user-related endpoints here (e.g., update profile picture)
	return router
}

func AuthRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/login", handlers.Login)
	router.Post("/logout", handlers.Logout)

	// Add new auth-related endpoints here (e.g., refresh token)
	return router
}

func EventRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllEvents) // Assuming filtering is implemented in handler
	router.Get("/{eventId}", handlers.GetEvent)
	router.Delete("/{eventId}", handlers.DeleteEvent)
	router.Put("/{eventId}", handlers.UpdateEvent)
	router.Post("/", handlers.CreateEvent)

	// Add new event-related endpoints here (e.g., attend/unattend event, submit feedback)
	router.Post("/{eventId}/attend", handlers.AttendEvent)      // Example for attending an event
	router.Delete("/{eventId}/attend", handlers.UnattendEvent)  // Example for unattending an event
	router.Post("/{eventId}/feedback", handlers.CreateFeedback) // Example for submitting feedback
	return router
}

func RSORoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllRSOs) // Assuming filtering is implemented in handler
	router.Get("/{rsoId}", handlers.GetRSO)
	router.Delete("/{rsoId}", handlers.DeleteRSO)
	router.Put("/{rsoId}", handlers.UpdateRSO)
	router.Post("/", handlers.CreateRSO)

	// Add new RSO-related endpoints here (e.g., join/leave RSO)
	router.Post("/{rsoId}/join", handlers.JoinRSO)    // Example for joining an RSO
	router.Delete("/{rsoId}/join", handlers.LeaveRSO) // Example for leaving an RSO
	return router
}

func UniRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllUnis) // Assuming filtering is implemented in handler
	router.Get("/{uni_id}", handlers.GetUni)
	router.Delete("/{uni_id}", handlers.DeleteUni)
	router.Put("/{uni_id}", handlers.UpdateUniDetails)
	router.Post("/{uni_id}", handlers.CreateUni)

	// Add new Uni-related endpoints here (e.g. join/leave Uni)
	router.Post("/{uni_id}/join", handlers.JoinUni)
	router.Delete("/{uni_id}/join", handlers.LeaveUni)
	return router
}

func LocationRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetLocations) // Assuming all visible events in a given radius is shown

	// Add other routes as required (e.g. add/delete locations)
	return router
}
