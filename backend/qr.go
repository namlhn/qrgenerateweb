package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"

	go_qr "github.com/piglig/go-qr"
	"github.com/nfnt/resize"
)

func generateQRCode(w http.ResponseWriter, r *http.Request) {
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
	eyeShape := r.FormValue("eyeShape")
	dotShape := r.FormValue("dotShape")

	// Default colors
	dotColor, _ := parseColor("#000000")
	bgColor, _ := parseColor("#FFFFFF")

	if c, err := parseColor(dotColorStr); err == nil {
		dotColor = c
	}
	if c, err := parseColor(bgColorStr); err == nil {
		bgColor = c
	}

	// Generate QR code
	qr, err := go_qr.EncodeText(text, go_qr.Medium)
	if err != nil {
		log.Printf("Error encoding QR code: %v", err)
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	// Create image from QR code
	qrImg := createImage(qr, dotColor, bgColor, eyeShape, dotShape)

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
		qrImg = overlayLogo(qrImg, logo)
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

// createImage creates an image from a QR code with custom shapes and colors.
func createImage(qr *go_qr.QrCode, dotColor, bgColor color.Color, eyeShape, dotShape string) image.Image {
	size := qr.GetSize()
	img := image.NewRGBA(image.Rect(0, 0, size*10, size*10))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if qr.GetModule(x, y) {
				isEye, isEyeBall := isEyeModule(x, y, size)

				currentShape := dotShape
				if isEye {
					currentShape = eyeShape
				}

				// For eye balls, always draw a solid square
				if isEyeBall {
					currentShape = "square"
				}

				switch currentShape {
				case "circle":
					drawCircle(img, x*10+5, y*10+5, 5, dotColor)
				case "rounded":
					// Simple rounded square - draw a circle and a smaller square inside
					drawCircle(img, x*10+5, y*10+5, 5, dotColor)
					rect := image.Rect(x*10+2, y*10+2, (x+1)*10-2, (y+1)*10-2)
					draw.Draw(img, rect, &image.Uniform{dotColor}, image.Point{}, draw.Src)
				default: // "square"
					rect := image.Rect(x*10, y*10, (x+1)*10, (y+1)*10)
					draw.Draw(img, rect, &image.Uniform{dotColor}, image.Point{}, draw.Src)
				}
			}
		}
	}

	return img
}

// isEyeModule checks if a module is part of a finder pattern (an "eye").
func isEyeModule(x, y, size int) (isEye bool, isEyeBall bool) {
	// Check for top-left eye
	if x >= 0 && x < 7 && y >= 0 && y < 7 {
		isEye = true
		if x >= 2 && x < 5 && y >= 2 && y < 5 {
			isEyeBall = true
		}
		return
	}
	// Check for top-right eye
	if x >= size-7 && x < size && y >= 0 && y < 7 {
		isEye = true
		if x >= size-5 && x < size-2 && y >= 2 && y < 5 {
			isEyeBall = true
		}
		return
	}
	// Check for bottom-left eye
	if x >= 0 && x < 7 && y >= size-7 && y < size {
		isEye = true
		if x >= 2 && x < 5 && y >= size-5 && y < size-2 {
			isEyeBall = true
		}
		return
	}
	return false, false
}


// drawCircle draws a filled circle.
func drawCircle(img *image.RGBA, x0, y0, r int, c color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(x0+x, y0+y, c)
			}
		}
	}
}


// overlayLogo places a logo in the center of a QR code.
func overlayLogo(qrCode, logo image.Image) image.Image {
	// Resize the logo to be 1/4 of the QR code size
	logo = resize.Resize(uint(qrCode.Bounds().Dx()/4), 0, logo, resize.Lanczos3)

	// Create a new image with the same dimensions as the QR code
	b := qrCode.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, qrCode, image.Point{}, draw.Src)

	// Calculate the position to place the logo (center)
	offset := image.Pt((b.Dx()-logo.Bounds().Dx())/2, (b.Dy()-logo.Bounds().Dy())/2)

	// Draw the logo over the QR code
	draw.Draw(m, logo.Bounds().Add(offset), logo, image.Point{}, draw.Over)

	return m
}