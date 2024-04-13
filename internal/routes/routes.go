package routes

import (
	"net/http"

	"github.com/bingKegeta/Knight-Link/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

func TokenAuthMiddleware(tokenAuth *jwtauth.JWTAuth, handler func(*jwtauth.JWTAuth, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(tokenAuth, w, r)
	}
}

func Routes(tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}),
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/users", UserRoutes(tokenAuth))
		r.Mount("/api/auth", AuthRoutes(tokenAuth))
		r.Mount("/api/events", EventRoutes(tokenAuth))
		r.Mount("/api/rsos", RSORoutes())
		r.Mount("/api/unis", UniRoutes())
		r.Mount("/api/locations", LocationRoutes())
		// Add new route groups here
	})

	return router
}

func UserRoutes(tokenAuth *jwtauth.JWTAuth) http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/{userId}", handlers.GetUser)
	})

	// Routes without token need.
	// router.Get("/{userId}", handlers.GetUser)
	router.Post("/", handlers.CreateUser)
	router.Get("/get_students", handlers.GetAllStudents)

	// Add new user-related endpoints here (e.g., update profile picture)
	return router
}

func AuthRoutes(tokenAuth *jwtauth.JWTAuth) http.Handler {
	router := chi.NewRouter()
	router.Post("/login", TokenAuthMiddleware(tokenAuth, handlers.Login))
	router.Post("/logout", handlers.Logout)

	// Add new auth-related endpoints here (e.g., refresh token)
	return router
}

func EventRoutes(tokenAuth *jwtauth.JWTAuth) http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", handlers.CreateEvent)
	})
	router.Get("/", handlers.GetAllEvents)
	router.Delete("/{eventId}", handlers.DeleteEvent)
	router.Put("/{eventId}", handlers.UpdateEvent)

	// router.Post("/", handlers.CreateEvent) commented in favor of auth version

	// Add new event-related endpoints here (e.g., attend/unattend event, submit feedback)
	router.Post("/attend", handlers.AttendEvent)      // Example for attending an event
	router.Delete("/attend", handlers.UnattendEvent)  // Example for unattending an event
	router.Post("/feedback", handlers.CreateFeedback) // Example for submitting feedback
	router.Get("/feedback", handlers.GetFeedback)
	return router
}

func RSORoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllRSOs)
	router.Post("/", handlers.CreateRSO)
	router.Get("/{rsoId}", handlers.GetRSO)
	router.Delete("/{rsoId}", handlers.DeleteRSO)
	router.Put("/{rsoId}", handlers.UpdateRSO)

	// Add new RSO-related endpoints here (e.g., join/leave RSO)
	router.Post("/join", handlers.JoinRSO)            // Example for joining an RSO
	router.Delete("/{rsoId}/join", handlers.LeaveRSO) // Example for leaving an RSO
	return router
}

func UniRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllUnis)
	router.Put("/{uni_id}", handlers.UpdateUniDetails)
	// Add new Uni-related endpoints here (e.g. join/leave Uni)

	return router
}

func LocationRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", handlers.GetAllLocations) // Assuming all visible events in a given radius is shown
	router.Post("/create", handlers.CreateLocation)
	// Add other routes as required (e.g. add/delete locations)
	return router
}
