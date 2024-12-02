package likes

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AddLikeList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("list_id")
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	list_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `INSERT INTO "ListLike" (user_id, list_id) VALUES ($1, $2)`
	_, err = conn.Exec(context.Background(), query, user_id, list_id)
	if err != nil {
		log.Println("Error executing query", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query = `UPDATE "List" SET total_likes = total_likes + 1 WHERE id = $1`
	_, err = conn.Exec(context.Background(), query, list_id)
	if err != nil {
		log.Println("Error executing query", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Like added to review"}`))
}

func RemoveLikeList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("list_id")
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	list_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `DELETE FROM "ListLike" WHERE user_id = $1 AND list_id = $2`
	_, err = conn.Exec(context.Background(), query, user_id, list_id)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error removing values from table ListLike", http.StatusInternalServerError)
		return
	}

	query = `UPDATE "List" SET total_likes = total_likes - 1  WHERE id = $1`
	_, err = conn.Exec(context.Background(), query, list_id)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error updating table List", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Like removed from list"}`))
}
