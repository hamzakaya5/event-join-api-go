package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gamegos_case/database"
	"gamegos_case/models"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
	//save to db the password after hashing
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

	// Set response header and send JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func PrivateHandler(w http.ResponseWriter, r *http.Request) {

	eventNo := r.Header.Get("eventNo")
	userId := r.Header.Get("userId")
	level := r.Header.Get("level")
	// level := 61

	var groupId int
	var prevEventId int
	var categoryId int

	usr, nil1 := strconv.Atoi(userId)
	evnt, nil2 := strconv.Atoi(eventNo)
	lvl, nil3 := strconv.Atoi(level)
	if nil1 != nil || nil2 != nil || nil3 != nil {
		http.Error(w, `{"error":"invalid userId or eventNo or level"}`, http.StatusBadRequest)
		return
	}

	// fmt.Println(usr, evnt)

	var isEnrolled int
	// var lvl int

	ctx := r.Context()

	tx, err := database.DBConn.Pool.Begin(ctx)
	if err != nil {
		http.Error(w, `{"error":"could not start tx"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	// err = tx.QueryRow(ctx,
	// 	"SELECT event_number FROM players WHERE id=$1",
	// 	usr,
	// ).Scan(&isEnrolled)
	isEnrolled, err = GetPlayersEvent(tx, ctx, usr)

	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("Player not found:", usr)
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		fmt.Println("Error querying player:", err)
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}

	// fmt.Println("pgrow", lvl)
	// fmt.Println(isEnrolled, "isenrolled")
	if isEnrolled != 0 {
		http.Error(w, "User is enrolled in any event", http.StatusAlreadyReported)
		return
	}

	// Check if event exists
	var isExist int
	isExist, err = IsEventExist(tx, ctx, evnt)
	// err = tx.QueryRow(ctx,
	// 	"SELECT id FROM events WHERE id=$1",
	// 	evnt,
	// ).Scan(&isExist)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
			return
		}
		fmt.Println("Error querying event:", err)
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}

	if isExist == 0 {
		http.Error(w, "No such event exists", http.StatusNotFound)
		return
	}

	// err = tx.QueryRow(ctx,
	// 	"SELECT id FROM event_history WHERE player_id=$1 and event_id=$2",
	// 	usr, evnt,
	// ).Scan(&prevEventId)
	prevEventId, err = GetEventHistoryID(tx, ctx, usr, evnt)

	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error querying event_history:", err)
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}

	if prevEventId != 0 {
		http.Error(w, "Event already joined", http.StatusConflict)
		return
	}

	// err = tx.QueryRow(ctx,
	// 	"SELECT gr.id, ctg.id FROM group_name gr JOIN categories ctg ON ctg.id = gr.category_id WHERE  ctg.min_level <=  $1 and $1 <= ctg.max_level and gr.group_count < 10",
	// 	lvl,
	// ).Scan(&groupId, &categoryId)
	groupId, categoryId, err = GetCategoryForLevel(tx, ctx, lvl)

	if err != nil && err != pgx.ErrNoRows {
		fmt.Println("Error selecting group:", err)
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	fmt.Println(groupId, " groupid here", categoryId, lvl)

	if groupId == 0 {
		fmt.Println("no existing group found, creating new group")
		// var newCategoryId int

		// err := tx.QueryRow(ctx,
		// 	"SELECT id FROM categories WHERE $1 BETWEEN min_level and max_level",
		// 	lvl,
		// ).Scan(&categoryId)
		categoryId, err = GetCategoryForPlayerLevel(tx, ctx, lvl)
		if err != nil {
			fmt.Println("Error selecting category:", err)
			http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
			return
		}

		// err = tx.QueryRow(ctx,
		// 	`INSERT INTO group_name (category_id, event_id, group_count)
		// 	VALUES ($1, $2, $3	)
		// 	RETURNING id`,
		// 	categoryId, eventNo, 1,
		// ).Scan(&groupId)
		groupId, err = InsertGroup(tx, ctx, categoryId, evnt)
		// fmt.Print(groupId, "after added newgroupid")

		if err != nil {
			http.Error(w, "Error creating new group", http.StatusInternalServerError)
			return
		}

	} else {
		//increment group count
		fmt.Println(groupId, "existinggroupid")

		// _, err := tx.Exec(ctx,
		// 	`UPDATE group_name
		// 	SET group_count = group_count + 1
		// 	WHERE id = $1`, groupId,
		// )
		err = IncreaseGroupCount(tx, ctx, groupId)
		if err != nil {
			http.Error(w, `{"error":"could not update group count"}`, http.StatusInternalServerError)
			return
		}
	}

	fmt.Println(groupId, "finalgroupid", eventNo, userId, categoryId)

	// _, err = tx.Exec(ctx,
	// 	`UPDATE players
	// 	SET "group" = $1, event_number = $2
	// 	WHERE id = $3`, groupId, evnt, usr,
	// )
	err = UpdatePlayer(tx, ctx, groupId, evnt, usr)
	if err != nil {
		// fmt.Fprintf(w, "Error updating player: %v", err)
		http.Error(w, `{"error":"could not update group count"}`, http.StatusInternalServerError)
		return
	}

	// _, err = tx.Exec(ctx,
	// 	`INSERT INTO event_history (player_id, event_id )
	// 		VALUES ($1, $2)`, usr, evnt,
	// )
	err = UpdatePlayerEventHistory(tx, ctx, evnt, usr)
	if err != nil {
		fmt.Println("Error inserting event_history:", err)
		http.Error(w, `{"error":"could not insert event history"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Println("Transaction commit failed:", err) // <-- print/log the error
		http.Error(w, `{"error":"could not commit transaction"}`, http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Transaction committed successfully")
	}

	http.Error(w, "Player Registered Successfully!", http.StatusOK)

}
