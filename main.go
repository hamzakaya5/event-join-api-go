package main

import (
	"gamegos_case/database"
	"gamegos_case/routes"
	"log"
	"net/http"
)

func main() {
	// PostgreSQL configuration

	database.InitRedis()
	err := database.ConnectToPostgre()
	if err != nil {
		panic(err)
	}

	mux := routes.RegisterRoutes()

	// Start server
	log.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
