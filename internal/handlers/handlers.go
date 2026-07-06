package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"tiny-links/internal/utils"
)

type BuildLinkReq struct {
	Url string `json:"url"`
}

// Upper case so that it is exported and public
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

func WelcomeHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to the tiny-links api"))
		})
}

// Need to take tiny url and map it back to its long url
// and redirect user
func RedirectHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			tiny_url := r.PathValue("tiny_url")
			link := &Link{}
			// Find original url in db
			err := db.QueryRow("SELECT * FROM link WHERE tiny_url = (?)", tiny_url).Scan(&link.ID, &link.Url, &link.TinyUrl)
			if err != nil {
				log.Fatal(err)
			}
			// Redirect user & doubles as sending response

			builder := strings.Builder{}
			// This is to make it "absolute"
			// Seems like things could become an issue if I don't
			// validate links on building part.
			// TODO Build Link valdiating middleware? or just validate the link
			// TODO in buildLinkHandler. But this will do for now
			builder.WriteString("//")
			builder.WriteString(link.Url)
			http.Redirect(w, r, builder.String(), http.StatusMovedPermanently)
		})
}
func BuildLinkHandler(db *sql.DB, base62_map map[int8]rune) http.HandlerFunc {
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
			tiny_url := utils.Base62_encoding(int(id), base62_map)

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
