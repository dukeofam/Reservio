package routes

import (
	"net/http"
	"reservio/controllers"
	"reservio/middleware"

	"github.com/gorilla/mux"
)

// SetupRouter returns a *mux.Router with all routes and middleware configured
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.CSRFMiddleware)

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", controllers.Register).Methods("POST")
	auth.HandleFunc("/login", controllers.Login).Methods("POST")
	auth.HandleFunc("/logout", controllers.Logout).Methods("POST")
	auth.HandleFunc("/request-reset", controllers.RequestPasswordReset).Methods("POST")
	auth.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST")

	parent := api.PathPrefix("/parent").Subrouter()
	parent.Use(middleware.Protected)
	parent.HandleFunc("/children", controllers.AddChild).Methods("POST")
	parent.HandleFunc("/children", controllers.GetChildren).Methods("GET")
	parent.HandleFunc("/reserve", controllers.MakeReservation).Methods("POST")
	parent.HandleFunc("/reservations", controllers.GetMyReservations).Methods("GET")
	parent.HandleFunc("/reservations/{id}", controllers.CancelReservation).Methods("DELETE")
	parent.HandleFunc("/children/{id}", controllers.EditChild).Methods("PUT")
	parent.HandleFunc("/children/{id}", controllers.DeleteChild).Methods("DELETE")

	// Admin routes: all require Protected + AdminOnly middleware
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.Protected)
	admin.Use(middleware.AdminOnly)
	admin.HandleFunc("/slots", controllers.CreateSlot).Methods("POST")
	admin.HandleFunc("/approve/{id}", controllers.ApproveReservation).Methods("PUT")
	admin.HandleFunc("/reject/{id}", controllers.RejectReservation).Methods("PUT")
	admin.HandleFunc("/reservations", controllers.GetReservationsByStatus).Methods("GET")
	admin.HandleFunc("/users", controllers.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/users/{id}/role", controllers.UpdateUserRole).Methods("PUT")

	user := api.PathPrefix("/user").Subrouter()
	user.Use(middleware.Protected)
	user.HandleFunc("/profile", controllers.GetProfile).Methods("GET")
	user.HandleFunc("/profile", controllers.UpdateProfile).Methods("PUT")

	api.HandleFunc("/slots", controllers.ListSlots).Methods("GET")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version": "1.0.0", "commit": "dev"}`))
	}).Methods("GET")

	return r
}
