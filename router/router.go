package router

import (
	"cms-project/internal/navLinks"
	"cms-project/internal/users"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// New creates and configures a new router
func New() *mux.Router {
	r := mux.NewRouter()

	// Add middleware
	r.Use(corsMiddleware)

	// Redirect root to Swagger documentation
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// Add a catch-all route for OPTIONS requests
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Register routes
	usersRouter := r.PathPrefix("/api/users").Subrouter()
	users.RegisterUsersRoutes(usersRouter)

	// Nav Links routes
	navLinksRouter := r.PathPrefix("/api/nav_links").Subrouter()
	navLinks.RegisterNavLinksRoutes(navLinksRouter)

	// Swagger route
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Add a catch-all handler for 404s
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "Route not found"}`))
}

// func registerMenuRoutes(r *mux.Router) {
// 	r.HandleFunc("/menus", menu.GetMenusHandler).Methods("GET")
// 	r.HandleFunc("/menus", menu.CreateMenuHandler).Methods("POST")
// 	r.HandleFunc("/menus/{id}", menu.GetMenuByIDHandler).Methods("GET")
// 	r.HandleFunc("/menus/{id}", menu.DeleteMenuHandler).Methods("DELETE")
// }
