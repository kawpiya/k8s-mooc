package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func generateRandomString() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func main() {
	appRandomString, err := generateRandomString()
	if err != nil {
		panic(err)
	}

	fmt.Println("Application started")
	fmt.Println("Random string:", appRandomString)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		currentTimestamp := time.Now().Format(time.RFC3339)
		w.Write([]byte(currentTimestamp + ": " + appRandomString))
	})
	log.Fatal(http.ListenAndServe(addr, nil))

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Printf("%s: %s\n", t.Format(time.RFC3339), appRandomString)
	}
}
