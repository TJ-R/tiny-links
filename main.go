package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tiny-links/internal/handlers"
	"tiny-links/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./tiny-links.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS link (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			tiny_url TEXT
		)`)

	if err != nil {
		log.Fatal(err)
	}

	// Create server and listen for requests
	base62_map := utils.Make_base62_map()

	http.HandleFunc("/", handlers.WelcomeHandler())
	http.HandleFunc("/build-link", handlers.BuildLinkHandler(db, base62_map))
	http.HandleFunc("/{tiny_url}", handlers.RedirectHandler(db))

	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "Listening on port 8080")
	fmt.Println()
	log.Fatal(s.ListenAndServe())
}
