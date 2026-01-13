package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
	randomString, err := generateRandomString()
	if err != nil {
		panic(err)
	}

	fmt.Println("Application started")
	fmt.Println("Random string:", randomString)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Printf("%s: %s\n", t.Format(time.RFC3339), randomString)
	}
}
