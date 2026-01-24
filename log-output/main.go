package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// file, err := os.Open("/logs/output.log")
		// if err != nil {
		// 	http.Error(w, "Failed to read file", http.StatusInternalServerError)
		// 	return
		// }
		// defer file.Close()
		// scanner := bufio.NewScanner(file)
		// line := ""
		// if scanner.Scan() {
		// 	line = scanner.Text()
		// 	fmt.Println(line)
		// } else {
		// 	http.Error(w, "File is empty", http.StatusInternalServerError)
		// 	return
		// }

		// if err := scanner.Err(); err != nil {
		// 	panic(err)
		// }

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		resp, err := http.Get("http://ping-pong-app-svc:80/pings")
		if err != nil {
			http.Error(w, "Failed to fetch data:"+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read data", http.StatusInternalServerError)
			return
		}

		line := string(data)

		currentTimestamp := time.Now().Format(time.RFC3339)
		message := os.Getenv("MESSAGE")

		file, err := os.ReadFile("/config/information.txt")
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		output := fmt.Sprintf("file content: %s\nenv variable: MESSAGE=%s\n%s: %s\n%s", file, message, currentTimestamp, appRandomString, line)
		w.Write([]byte(output))
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
