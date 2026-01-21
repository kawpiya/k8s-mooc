package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var counter uint64

func main() {
	fmt.Println("Ping-pong app started!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddUint64(&counter, 1)

		fmt.Fprintf(w, "Ping / Pongs: %d\n", n)
	})

	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.LoadUint64(&counter)
		fmt.Fprintf(w, "Ping / Pongs: %d\n", n)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
