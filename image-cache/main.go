package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	imageURL  = "https://picsum.photos/id/237/400/400"
	cacheDir  = "cache"
	imageName = "picsum.jpg"
	interval  = 10 * time.Minute
)

var (
	imageBytes []byte
	mu         sync.RWMutex
)

// download and replace cache
func refreshImage() {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Println("Image download failed:", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Image read failed:", err)
		return
	}

	// Save to disk
	os.MkdirAll(cacheDir, 0755)
	path := filepath.Join(cacheDir, imageName)
	os.WriteFile(path, data, 0644)

	// Save to memory
	mu.Lock()
	imageBytes = data
	mu.Unlock()

	log.Println("Image cache refreshed")
}

// background scheduler
func startImageRefresher() {
	// Initial download
	refreshImage()

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			refreshImage()
		}
	}()
}

// serve cached image
func imageHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	data := imageBytes
	mu.RUnlock()

	if data == nil {
		http.Error(w, "Image not ready", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=600") // 10 minutes
	w.Write(data)
}

// serve HTML page
func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
	<title>Image + Text Submit</title>
	<style>
		body {
			font-family: sans-serif;
			max-width: 600px;
			margin: 40px auto;
		}
		img {
			max-width: 100%;
			border-radius: 8px;
			margin-bottom: 20px;
		}
		.counter {
			font-size: 0.9em;
			color: #555;
		}
	</style>
</head>
<body>

<h2>Cached Image (updates every 10 minutes)</h2>

<img src="/image" alt="Cached Picsum Image">

<h3>Enter text (max 140 characters)</h3>

<form method="POST" action="/submit">
	<textarea
		name="message"
		rows="4"
		style="width: 100%;"
		maxlength="140"
		oninput="updateCounter(this)"
		required
	></textarea>

	<div class="counter" id="counter">
		140 characters remaining
	</div>

	<br>
	<button type="submit">Submit</button>
</form>

<ul>
	<li>Learn Docker</li>
	<li>Learn Kubernetes</li>
</ul>
<script>
	function updateCounter(el) {
		const remaining = 140 - el.value.length;
		document.getElementById("counter").textContent =
			remaining + " characters remaining";
	}
</script>

</body>
</html>
`))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg := r.FormValue("message")

	if utf8.RuneCountInString(msg) > 140 {
		http.Error(w, "Message exceeds 140 characters", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h3>Message received:</h3><p>" +
		template.HTMLEscapeString(msg) + "</p><a href='/'>Back</a>"))
}

func main() {
	startImageRefresher()

	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/submit", submitHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
