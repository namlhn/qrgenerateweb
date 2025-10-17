package main

import (
	"log"
	"net/http"

	"namlhn/qrgenerateweb/internal/handlers"
)

func main() {
	// API endpoint for QR code generation
	http.HandleFunc("/api/qr", handlers.GenerateQRCode)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve frontend files
	frontendFs := http.FileServer(http.Dir("./web/templates/"))
	http.Handle("/", frontendFs)

	log.Println("Server starting on port 8099...")
	if err := http.ListenAndServe(":8099", nil); err != nil {
		log.Fatal(err)
	}
}
