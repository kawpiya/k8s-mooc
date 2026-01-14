package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World"))
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
