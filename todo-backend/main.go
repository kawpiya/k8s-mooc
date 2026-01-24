package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	todos = make([]string, 0)
	mu    sync.Mutex
)

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("todo")
	if text == "" || len(text) > 140 {
		http.Error(w, "Todo must be 1–140 characters", http.StatusBadRequest)
		return
	}

	mu.Lock()
	todos = append(todos, text)
	mu.Unlock()

	// Redirect back to page (prevents duplicate submit on refresh)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodosHandler(w, r)
		case http.MethodPost:
			createTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
