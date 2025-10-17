package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strconv"

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

	// Generate QR code
	qr, err := go_qr.EncodeText(text, go_qr.Medium)
	if err != nil {
		log.Printf("Error encoding QR code: %v", err)
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	// Create image from QR code
	qrImg := createImage(qr, dotColor, bgColor, eyeColor, eyeShape, dotShape, padding)

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
func createImage(qr *go_qr.QrCode, dotColor, bgColor, eyeColor color.Color, eyeShape, dotShape string, padding int) image.Image {
	size := qr.GetSize()
	imgSize := size*10 + 2*padding*10
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if qr.GetModule(x, y) {
				isEye, isEyeBall := isEyeModule(x, y, size)

				currentColor := dotColor
				currentShape := dotShape

				if isEye {
					currentColor = eyeColor
					currentShape = eyeShape
				}
				if isEyeBall {
					currentShape = "square" // Eyeballs are always square for clarity
				}

				rectX := (x + padding) * 10
				rectY := (y + padding) * 10

				switch currentShape {
				case "circle":
					drawCircle(img, rectX+5, rectY+5, 5, currentColor)
				case "rounded":
					drawCircle(img, rectX+5, rectY+5, 5, currentColor)
					rect := image.Rect(rectX+2, rectY+2, rectX+8, rectY+8)
					draw.Draw(img, rect, &image.Uniform{currentColor}, image.Point{}, draw.Src)
				case "diamond":
					// A simple diamond shape
					for i := 0; i < 5; i++ {
						drawLine(img, rectX+5, rectY+i, rectX+5, rectY+10-i, currentColor)
						drawLine(img, rectX+i, rectY+5, rectX+10-i, rectY+5, currentColor)
					}
				default: // "square"
					rect := image.Rect(rectX, rectY, rectX+10, rectY+10)
					draw.Draw(img, rect, &image.Uniform{currentColor}, image.Point{}, draw.Src)
				}
			}
		}
	}

	return img
}

// drawLine draws a line (used for diamond shape)
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
    // Basic line drawing, for more complex cases a proper algorithm like Bresenham's should be used
    for x := x1; x <= x2; x++ {
        for y := y1; y <= y2; y++ {
            img.Set(x, y, c)
        }
    }
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
func overlayLogo(qrCode image.Image, logo image.Image) image.Image {
	// Resize the logo to be 1/5 of the QR code size for better proportions
	logoSize := uint(qrCode.Bounds().Dx() / 5)
	logo = resize.Resize(logoSize, 0, logo, resize.Lanczos3)

	// Create a new drawable image from the qrCode
	b := qrCode.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, qrCode, image.Point{}, draw.Src)

	// Calculate the position to place the logo (center)
	offset := image.Pt((b.Dx()-logo.Bounds().Dx())/2, (b.Dy()-logo.Bounds().Dy())/2)
	logoBounds := logo.Bounds().Add(offset)

	// Clear the area behind the logo by drawing a white rectangle
	// This uses the background color of the QR code, assuming it's mostly white
	// A more robust solution might pass the bgColor from the main function
	draw.Draw(m, logoBounds, &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw the logo over the cleared area
	draw.Draw(m, logoBounds, logo, image.Point{}, draw.Over)

	return m
}