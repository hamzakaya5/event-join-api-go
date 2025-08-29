package handlers

import (
	"context"
	"database/sql"
	"gamegos_case/models"
	"strconv"

	"github.com/jackc/pgx/v5"
	// "encoding/json"
	// "fmt"
	// "gamegos_case/database"
	// "gamegos_case/models"
	// "net/http"
	// "strconv"
	// "time"
)

func GetPlayersEvent(tx pgx.Tx, ctx context.Context, usr int) (int, error) {
	var isEnrolledLocal sql.NullInt32
	err := tx.QueryRow(ctx,
		"SELECT event_number FROM players WHERE id=$1",
		usr,
	).Scan(&isEnrolledLocal)
	if err != nil {
		return 0, err
	}
	return int(isEnrolledLocal.Int32), nil
}

func IsEventExist(tx pgx.Tx, ctx context.Context, eventID int) (exists int, err error) {
	var isExist int

	err = tx.QueryRow(ctx,
		"SELECT id FROM events WHERE id=$1",
		eventID,
	).Scan(&isExist)

	if err != nil {
		if err == pgx.ErrNoRows {
			return isExist, nil
		}
		return isExist, err
	}
	return isExist, nil

}

func GetEventHistoryID(tx pgx.Tx, ctx context.Context, playerID int, eventID int) (prevEventId int, err error) {

	err = tx.QueryRow(ctx,
		"SELECT id FROM event_history WHERE player_id=$1 and event_id=$2",
		playerID, eventID,
	).Scan(&prevEventId)
	return prevEventId, err
}

func GetCategoryForLevel(tx pgx.Tx, ctx context.Context, lvl int, event_id int) (groupId int, categoryId int, err error) {
	var count int
	err = tx.QueryRow(ctx,
		"SELECT gr.id, ctg.id, gr.group_count FROM group_name gr JOIN categories ctg ON ctg.id = gr.category_id WHERE  ctg.min_level <=  $1 and $1 <= ctg.max_level and gr.group_count < 10",
		lvl,
	).Scan(&groupId, &categoryId, &count)
	if err != nil {
		return 0, 0, err
	}

	key := strconv.Itoa(event_id) + "|" + strconv.Itoa(groupId)
	val := models.GroupProcessHolder[key]

	if (count + val) < 10 {
		//we are increasing group process number so other request would see how many players a group includes actually.
		IncreaseGroupProcess(key)
		return groupId, categoryId, nil
	}

	return 0, categoryId, nil
}

func GetCategoryForPlayerLevel(tx pgx.Tx, ctx context.Context, lvl int) (categoryId int, err error) {
	err = tx.QueryRow(ctx,
		"SELECT id FROM categories WHERE $1 BETWEEN min_level and max_level",
		lvl,
	).Scan(&categoryId)
	if err != nil {
		return 0, err
	}
	return categoryId, nil
}

func InsertGroup(tx pgx.Tx, ctx context.Context, categoryId int, eventNo int) (groupId int, err error) {
	err = tx.QueryRow(ctx,
		`INSERT INTO group_name (category_id, event_id, group_count)
			VALUES ($1, $2, $3	)
			RETURNING id`,
		categoryId, eventNo, 1,
	).Scan(&groupId)

	if err != nil {
		return 0, err
	}
	return groupId, nil
}

func IncreaseGroupCount(tx pgx.Tx, ctx context.Context, groupId int) (err error) {
	_, err = tx.Exec(ctx,
		`UPDATE group_name
		SET group_count = group_count + 1
		WHERE id = $1`, groupId,
	)
	return err
}

func UpdatePlayer(tx pgx.Tx, ctx context.Context, groupId int, evnt int, usr int) (err error) {
	_, err = tx.Exec(ctx,
		`UPDATE players
		SET "group" = $1, event_number = $2
		WHERE id = $3`, groupId, evnt, usr,
	)
	return err
}

func UpdatePlayerEventHistory(tx pgx.Tx, ctx context.Context, eventNo int, usr int) (err error) {
	_, err = tx.Exec(ctx,
		`INSERT INTO event_history (player_id, event_id )
			VALUES ($1, $2)`, usr, eventNo,
	)
	return err
}

func IncreaseGroupProcess(key string) {
	l := GetLock(key)
	l.Lock()
	defer l.Unlock()

	models.GroupProcessHolder[key]++
}

func DecreaseGroupProcess(key string) {
	models.GroupProcessHolder[key]++
}
