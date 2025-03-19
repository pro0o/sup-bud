package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define flags
	port := flag.Int("port", 8080, "Port to serve on")
	directory := flag.String("dir", "./web", "Directory to serve")
	flag.Parse()

	// Get absolute path
	dir, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatal(err)
	}

	// Configure server
	fileServer := http.FileServer(http.Dir(dir))

	// Add WASM MIME type for proper handling
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if filepath.Ext(r.URL.Path) == ".wasm" {
			w.Header().Set("Content-Type", "application/wasm")
		}
		fileServer.ServeHTTP(w, r)
	})

	// Start server
	address := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting server on http://localhost%s serving from %s\n", address, dir)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}
