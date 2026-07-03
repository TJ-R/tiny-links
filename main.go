package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
	base62_map := make_base62_map()
	fmt.Println(base62_encoding(123, base62_map))

	http.HandleFunc("/build-link", buildLinkHandler(db, base62_map))

	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "Listening on port 8080")
	fmt.Println()
	log.Fatal(s.ListenAndServe())
}

type BuildLinkReq struct {
	Url string `json:"url"`
}

// Upper case so that it is exported and not private
type BuildLinkResp struct {
	Url string `json:"url"`
}

type ErrorResp struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Link struct {
	ID      int64
	Url     string
	TinyUrl string
}

func buildLinkHandler(db *sql.DB, base62_map map[int8]rune) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			buildLinkReq := BuildLinkReq{}

			err := decoder.Decode(&buildLinkReq)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				http.Error(w, "Error when decoding request", http.StatusInternalServerError)
				return
			}
			// Insert
			res, err := db.Exec("INSERT INTO link (url) VALUES (?)", buildLinkReq.Url)

			if err != nil {
				log.Fatal(err)
			}

			id, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			// Need to build the shortened url
			tiny_url := base62_encoding(int(id), base62_map)

			// Update db with tiny url
			link := &Link{}
			err = db.QueryRow("UPDATE link SET tiny_url = ? WHERE id = ? RETURNING id, url, tiny_url", tiny_url, id).Scan(&link.ID, &link.Url, &link.TinyUrl)
			if err != nil {
				log.Fatal(err)
			}

			// Return it
			builder := strings.Builder{}
			builder.WriteString("localhost:8080/")
			builder.WriteString(link.TinyUrl)

			resp := BuildLinkResp{Url: builder.String()}
			dat, err := json.Marshal(resp)

			if err != nil {
				fmt.Fprintf(os.Stderr, "[Error]: Encountered an error when marshalling json reponse.")
				errorRes := ErrorResp{Message: "Internal Server Error", Code: http.StatusInternalServerError}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(errorRes.Code)

				// TODO Determine if this actually writes the err res
				json.NewEncoder(w).Encode(errorRes)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(dat)
		},
	)
}

func base62_encoding(id int, base62_map map[int8]rune) string {
	stringBuilder := strings.Builder{}
	quotient := id
	remainder := -1
	for {
		remainder = quotient % 62
		quotient /= 62

		// if remainder is 0 we are done
		if remainder == 0 {
			break
		} else {
			stringBuilder.WriteRune(base62_map[int8(remainder)])
		}
	}

	return reverseString(stringBuilder.String())
}

func reverseString(s string) string {
	stringBuilder := strings.Builder{}

	for i := len(s) - 1; i >= 0; i-- {
		stringBuilder.WriteByte(s[i])
	}

	return stringBuilder.String()
}

func make_base62_map() map[int8]rune {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	m := make(map[int8]rune)

	for i := range int8(len(alphabet)) {
		m[i] = rune(alphabet[i])
	}

	return m
}
