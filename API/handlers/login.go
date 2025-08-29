package handlers

import (
	"context"
	"encoding/json"
	"gamegos_case/database"
	"gamegos_case/models"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	email := r.Header.Get("email")
	password := r.Header.Get("password")

	// hash, err := HashPassword(password)
	// if err != nil {
	// 	fmt.Fprintln(w, "Error while hashing password", err)
	// 	return
	// }

	var storedHash string
	var lvl int
	var user_id int

	err := database.DBConn.Pool.QueryRow(ctx,
		"SELECT email, password_hash, level, id FROM players WHERE email=$1",
		email,
	).Scan(&email, &storedHash, &lvl, &user_id)

	if err != nil {
		http.Error(w, `{"error":"invalid email or user not found"}`, http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(password, storedHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(email)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	err = database.RedisClient.Set(ctx, token, email, 24*time.Hour).Err()

	if err != nil {
		http.Error(w, "Error storing token in Redis", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Status:  200,
		Message: email + "logged in successfully",
		Token:   token,
		Level:   lvl,
		Id:      user_id,
	}

	// Set response header and send JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
