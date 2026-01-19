package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var counter uint64

func handler(w http.ResponseWriter, r *http.Request) {
	n := atomic.AddUint64(&counter, 1) - 1
	fmt.Fprintf(w, "pong %d\n", n)
}

func main() {

	fmt.Println("Ping-pong app started!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)

	http.HandleFunc("/pingpong", handler)
	http.ListenAndServe(addr, nil)
}
