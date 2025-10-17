package services

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"

	go_qr "github.com/piglig/go-qr"
	"namlhn/qrgenerateweb/pkg/qr"
)

// QRService handles QR code generation business logic
type QRService struct{}

// NewQRService creates a new QR service instance
func NewQRService() *QRService {
	return &QRService{}
}

// GenerateQRCode generates a QR code based on the request parameters
func (s *QRService) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text parameter is missing", http.StatusBadRequest)
		return
	}

	// Get customization options from the form
	dotColorStr := r.FormValue("dotColor")
	bgColorStr := r.FormValue("bgColor")
	eyeColorStr := r.FormValue("eyeColor")
	eyeShape := r.FormValue("eyeShape")
	dotShape := r.FormValue("dotShape")
	paddingStr := r.FormValue("padding")

	// Default colors and padding
	dotColor, _ := parseColor("#000000")
	bgColor, _ := parseColor("#FFFFFF")

	eyeColor := dotColor // Default eye color is the same as dot color
	if c, err := parseColor(dotColorStr); err == nil {
		dotColor = c
	}
	if c, err := parseColor(bgColorStr); err == nil {
		bgColor = c
	}
	if c, err := parseColor(eyeColorStr); err == nil {
		eyeColor = c
	}
	padding, err := strconv.Atoi(paddingStr)
	if err != nil || padding < 0 {
		padding = 1
	}

	// Generate QR code using the pkg/qr generator
	generator := &qr.QRGenerator{}
	qrCode, err := generator.EncodeText(text, go_qr.Medium)
	if err != nil {
		log.Printf("Error encoding QR code: %v", err)
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	// Create image from QR code
	qrImg := generator.CreateImage(qrCode, dotColor, bgColor, eyeColor, eyeShape, dotShape, padding)

	// Handle logo
	file, _, err := r.FormFile("logo")
	if err == nil {
		defer file.Close()

		logo, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Error decoding logo: %v", err)
			http.Error(w, "Error processing logo", http.StatusInternalServerError)
			return
		}

		// Combine QR code and logo
		qrImg = generator.OverlayLogo(qrImg, logo)
	}

	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, qrImg)
	if err != nil {
		log.Printf("Error encoding QR code to png: %v", err)
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}
}

// parseColor converts a hex color string (e.g., "#RRGGBB") to a color.Color object.
func parseColor(s string) (color.Color, error) {
	c := color.RGBA{A: 255}
	var err error
	if s[0] != '#' {
		return c, fmt.Errorf("invalid color format")
	}
	_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	if err != nil {
		return color.Black, err
	}
	return c, nil
}
