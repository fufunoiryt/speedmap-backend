package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// Configuration CORS - Headers explicites pour LibreSpeed
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, Content-Encoding, Content-Length, Cache-Control, Pragma, Origin, X-Requested-With")
	(*w).Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Encoding")
	(*w).Header().Set("Access-Control-Max-Age", "86400")
	(*w).Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	(*w).Header().Set("Pragma", "no-cache")
}

func main() {
	// --- 1. Empty (Ping / Upload) ---
	http.HandleFunc("/backend/empty.php", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}
		// Upload handling: read body but discard it
		if r.Method == "POST" {
			// Just read headers, body is discarded automatically by Go if not read, 
			// but for speedtest accuracy we might want to drain it? 
			// Standard http server reads request. 
			// Just respond 200 OK.
		}
		w.WriteHeader(http.StatusOK)
	})

	// --- 2. Garbage (Download) ---
	http.HandleFunc("/backend/garbage.php", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}

		// Default chunks
		chunks := 4
		ckSize := 1048576 // 1 MB default

		// Parse params
		if ckSizeStr := r.URL.Query().Get("ckSize"); ckSizeStr != "" {
			if val, err := strconv.Atoi(ckSizeStr); err == nil {
				ckSize = val * 1048576 // MB to bytes
			}
		}
		// LibreSpeed requests often don't send params for default garbage, 
		// but sometimes send 'chunks' for duration control?
		// The JS usually requests garbage.php without params for random data.
		
		// Create buffer of random data
		data := make([]byte, ckSize)
		rand.Read(data)

		w.Header().Set("Content-Description", "File Transfer")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=random.dat")
		w.Header().Set("Content-Transfer-Encoding", "binary")

		for i := 0; i < chunks; i++ {
			w.Write(data)
		}
	})

	// --- 3. Get IP ---
	http.HandleFunc("/backend/getIP.php", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}
		
		// Try to get real IP from headers (Cloud Run / Proxy)
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		fmt.Fprintf(w, "{\"processedString\": \"%s\", \"rawIspInfo\": \"\"}", ip)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("LibreSpeed Backend listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
