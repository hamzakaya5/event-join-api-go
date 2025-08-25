package main

import (
	"gamegos_case/database"
	"gamegos_case/routes"
	"log"
	"net/http"
	"time"
)

func main() {
	// PostgreSQL configuration

	database.InitRedis()

	for i := range 5 {
		log.Printf("Attempting to connect to PostgreSQL, attempt %d", i+1)

		err := database.ConnectToPostgre()
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
	}

	mux := routes.RegisterRoutes()

	// Start server
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
