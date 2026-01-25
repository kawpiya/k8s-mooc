package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	fmt.Println("Ping-pong app started!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)

	pingPongHandler := func(w http.ResponseWriter, r *http.Request) {
		var count int64

		err := db.QueryRow(`
		UPDATE pingpong
		SET count = count + 1
		RETURNING count
	`).Scan(&count)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Ping / Pongs: %d\n", count)
	}

	pingsHandler := func(w http.ResponseWriter, r *http.Request) {
		var count int64

		err := db.QueryRow(`
			SELECT count FROM pingpong LIMIT 1
		`).Scan(&count)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Ping / Pongs: %d\n", count)
	}

	http.HandleFunc("/pingpong", pingPongHandler)
	http.HandleFunc("/pings", pingsHandler)

	log.Fatal(http.ListenAndServe(addr, nil))
}
