package main

import (
	"crypto/rand"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Pre-generated random buffer — filled once at startup, reused for all downloads.
// Avoids per-request crypto/rand overhead that pins CPU at 100%.
var downloadBuf [65536]byte

func init() {
	rand.Read(downloadBuf[:])
}

//go:embed web/*
var webAssets embed.FS

func main() {
	port := 8080

	// Simple arg parsing: speedtest [port]
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "-h", "--help", "help":
			fmt.Fprintf(os.Stderr, "Usage: speedtest [port]\n\n  port  Port to listen on (default: 8080)\n\nExamples:\n  speedtest        # listen on :8080\n  speedtest 3000   # listen on :3000\n")
			os.Exit(0)
		case "-v", "--version", "version":
			fmt.Println("speedtest v0.1.0")
			os.Exit(0)
		default:
			p, err := strconv.Atoi(arg)
			if err != nil || p < 1 || p > 65535 {
				fmt.Fprintf(os.Stderr, "Error: invalid port %q\n\nUsage: speedtest [port]\n", arg)
				os.Exit(1)
			}
			port = p
		}
	}

	mux := http.NewServeMux()

	// Serve embedded web UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := webAssets.ReadFile("web/index.html")
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// Download endpoint - streams random data for speed measurement
	mux.HandleFunc("/api/download", func(w http.ResponseWriter, r *http.Request) {
		sizeStr := r.URL.Query().Get("size")
		size := 25 * 1024 * 1024
		if sizeStr != "" {
			if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100*1024*1024 {
				size = s
			}
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(size))
		w.Header().Set("Cache-Control", "no-store")

		// Reuse pre-generated buffer — no per-request randomness needed
		remaining := size
		for remaining > 0 {
			n := len(downloadBuf)
			if n > remaining {
				n = remaining
			}
			written, err := w.Write(downloadBuf[:n])
			if err != nil {
				return
			}
			remaining -= written
		}
	})

	// Upload endpoint - consumes posted data
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		start := time.Now()
		n, _ := io.Copy(io.Discard, r.Body)
		elapsed := time.Since(start).Seconds()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"bytes":%d,"elapsed":%.4f}`, n, elapsed)
	})

	// Ping endpoint for latency measurement
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprintf(w, `{"time":%d}`, time.Now().UnixMilli())
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Speedtest server running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v\nTry a different port: speedtest <port>", err)
	}
}
