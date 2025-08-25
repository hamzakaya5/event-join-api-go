package routes

import (
	"gamegos_case/handlers"
	"gamegos_case/middleware"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	mux.Handle("/join", middleware.AuthMiddleware(http.HandlerFunc(handlers.PrivateHandler)))

	return mux
}
