package handlers

import (
	"context"
	"encoding/json"
	"gamegos_case/database"
	"gamegos_case/models"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	email := r.Header.Get("email")
	password := r.Header.Get("password")
	username := r.Header.Get("username")

	if email == "" || password == "" || username == "" {
		http.Error(w, "Missing header fields", http.StatusBadRequest)
		return
	}

	hash, error := HashPassword(password)
	if nil != error {
		resp := models.ErrorResponse{
			Status:  401,
			Message: "Missing header fields",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	// fmt.Fprintln(w, "Hashed password", hash, username, email)

	num := Rand1to70()

	//produced a random number for representing level
	db_resp, _ := database.DBConn.Pool.Exec(ctx, "INSERT INTO players (username, email, password_hash, level) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", username, email, hash, num)
	// fmt.Print(db_resp.RowsAffected())

	if db_resp.RowsAffected() == 0 {
		resp := models.RegisterResponse{
			Status:  200,
			Message: username + " already registered",
			Level:   num,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return

	}
	// fmt.Fprintln(w, "User registered!", ctx, email, password)

	resp := models.RegisterResponse{
		Status:  200,
		Message: username + "registered successfully",
		Level:   num,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
