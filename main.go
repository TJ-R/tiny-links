package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
	"tiny-links/internal/handlers"
	"tiny-links/internal/utils"
)

func main() {
	db, err := sql.Open("sqlite3", "./tiny-links.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables(db)

	// Create server and listen for requests
	base62_map := utils.Make_base62_map()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.WelcomeHandler())
	mux.HandleFunc("POST /build-link", handlers.BuildLinkHandler(db, base62_map))
	mux.HandleFunc("GET /{tiny_url}", handlers.RedirectHandler(db))
	mux.HandleFunc("POST /user", handlers.HandleCreateUser(db))
	mux.HandleFunc("GET /user", handlers.HandleGetUsers(db))

	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.NewCrossOriginProtection().Handler(mux),
	}

	fmt.Fprintf(os.Stdout, "Listening on port 8080")
	fmt.Println()
	log.Fatal(s.ListenAndServe())
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS link (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				url TEXT NOT NULL,
				tiny_url TEXT
			)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}
}
