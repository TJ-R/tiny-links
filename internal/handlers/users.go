package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type UserCreationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleCreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		userCreationReq := UserCreationRequest{}

		err := decoder.Decode(&userCreationReq)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
			userCreationReq.Username,
			userCreationReq.Password,
		)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func HandleGetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
