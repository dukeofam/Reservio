package routes

import (
	"net/http"
	"reservio/controllers"
	"reservio/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter returns a *mux.Router with all routes and middleware configured
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Global middlewares
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.MetricsMiddleware)
	r.Use(middleware.SecurityHeadersMiddleware)
	r.Use(middleware.RateLimitMiddleware)

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// Health check endpoint for Docker
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", controllers.Register).Methods("POST")
	auth.HandleFunc("/login", controllers.Login).Methods("POST")
	auth.HandleFunc("/logout", controllers.Logout).Methods("POST")
	auth.HandleFunc("/request-reset", controllers.RequestPasswordReset).Methods("POST")
	auth.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST")

	parent := api.PathPrefix("/parent").Subrouter()
	parent.Use(middleware.Protected)
	parent.Use(middleware.CSRFMiddleware)
	parent.HandleFunc("/children", controllers.AddChild).Methods("POST")
	parent.HandleFunc("/children", controllers.GetChildren).Methods("GET")
	parent.HandleFunc("/children/{id}", controllers.GetChild).Methods("GET")
	parent.HandleFunc("/reserve", controllers.MakeReservation).Methods("POST")
	parent.HandleFunc("/reservations", controllers.GetMyReservations).Methods("GET")
	parent.HandleFunc("/reservations/{id}", controllers.CancelReservation).Methods("DELETE")
	parent.HandleFunc("/children/{id}", controllers.EditChild).Methods("PUT")
	parent.HandleFunc("/children/{id}", controllers.DeleteChild).Methods("DELETE")

	// Admin routes: all require Protected + AdminOnly middleware
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.Protected)
	admin.Use(middleware.AdminOnly)
	admin.Use(middleware.CSRFMiddleware)
	admin.HandleFunc("/slots", controllers.CreateSlot).Methods("POST")
	admin.HandleFunc("/approve/{id}", controllers.ApproveReservation).Methods("PUT")
	admin.HandleFunc("/reject/{id}", controllers.RejectReservation).Methods("PUT")
	admin.HandleFunc("/reservations", controllers.GetReservationsByStatus).Methods("GET")
	admin.HandleFunc("/users", controllers.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/users/{id}/role", controllers.UpdateUserRole).Methods("PUT")

	user := api.PathPrefix("/user").Subrouter()
	user.Use(middleware.Protected)
	user.Use(middleware.CSRFMiddleware)
	user.HandleFunc("/profile", controllers.GetProfile).Methods("GET")
	user.HandleFunc("/profile", controllers.UpdateProfile).Methods("PUT")

	api.HandleFunc("/slots", controllers.ListSlots).Methods("GET")
	api.HandleFunc("/slots/{id}", controllers.GetSlot).Methods("GET")

	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"version": "1.0.0", "commit": "dev"}`))
	}).Methods("GET")

	return r
}
