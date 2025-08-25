package routes

import (
	"gamegos_case/handlers"
	"gamegos_case/middleware"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Public route
	mux.HandleFunc("/login", handlers.LoginHandler)

	mux.HandleFunc("/register", handlers.RegisterHandler)

	// Private route with middleware
	mux.Handle("/join", middleware.AuthMiddleware(http.HandlerFunc(handlers.PrivateHandler)))

	return mux
}
