package main

import (
	"log"
	"net/http"
)

func main() {
	// API endpoint for QR code generation
	http.HandleFunc("/api/qr", generateQRCode)

	// Serve frontend files
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	// Serve static assets
	staticFs := http.FileServer(http.Dir("../static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFs))

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}