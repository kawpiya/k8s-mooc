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
		n := atomic.AddUint64(&counter, 1) - 1
		fmt.Fprintf(w, "pong %d\n", n)

		file, err := os.OpenFile(
			"/logs/output.log",
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			0644,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		if _, err := file.WriteString("Ping / Pongs: " + fmt.Sprint(n)); err != nil {
			panic(err)
		}
	})

	http.ListenAndServe(addr, nil)
}
