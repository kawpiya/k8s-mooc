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

	// Open file in append mode
	file, err := os.OpenFile(
		"/logs/output.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		line := fmt.Sprintf("[%s] %s\n", t.Format(time.RFC3339), appRandomString)

		// Write to stdout
		fmt.Print("Written line: " + line)

		// Append to file
		if _, err := file.WriteString(line); err != nil {
			panic(err)
		}
	}

	log.Fatal(http.ListenAndServe(addr, nil))
}
