package handlers

import (
	"net/http"

	"namlhn/qrgenerateweb/internal/services"
)

// GenerateQRCode handles QR code generation requests
func GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	qrService := services.NewQRService()
	qrService.GenerateQRCode(w, r)
}
