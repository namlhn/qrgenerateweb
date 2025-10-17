package qr

import (
	"image"
	"image/color"
	"image/draw"

	go_qr "github.com/piglig/go-qr"
	"github.com/nfnt/resize"
)

// QRGenerator provides QR code generation functionality
type QRGenerator struct{}

// EncodeText encodes text into a QR code
func (g *QRGenerator) EncodeText(text string, level go_qr.Ecc) (*go_qr.QrCode, error) {
	return go_qr.EncodeText(text, level)
}

// CreateImage creates an image from a QR code with custom shapes and colors
func (g *QRGenerator) CreateImage(qr *go_qr.QrCode, dotColor, bgColor, eyeColor color.Color, eyeShape, dotShape string, padding int) image.Image {
	size := qr.GetSize()
	imgSize := size*10 + 2*padding*10
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if qr.GetModule(x, y) {
				isEye, isEyeBall := g.isEyeModule(x, y, size)

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
					g.drawCircle(img, rectX+5, rectY+5, 5, currentColor)
				case "rounded":
					g.drawCircle(img, rectX+5, rectY+5, 5, currentColor)
					rect := image.Rect(rectX+2, rectY+2, rectX+8, rectY+8)
					draw.Draw(img, rect, &image.Uniform{currentColor}, image.Point{}, draw.Src)
				case "diamond":
					// A simple diamond shape
					for i := 0; i < 5; i++ {
						g.drawLine(img, rectX+5, rectY+i, rectX+5, rectY+10-i, currentColor)
						g.drawLine(img, rectX+i, rectY+5, rectX+10-i, rectY+5, currentColor)
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
func (g *QRGenerator) drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	// Basic line drawing, for more complex cases a proper algorithm like Bresenham's should be used
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			img.Set(x, y, c)
		}
	}
}

// isEyeModule checks if a module is part of a finder pattern (an "eye")
func (g *QRGenerator) isEyeModule(x, y, size int) (isEye bool, isEyeBall bool) {
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

// drawCircle draws a filled circle
func (g *QRGenerator) drawCircle(img *image.RGBA, x0, y0, r int, c color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(x0+x, y0+y, c)
			}
		}
	}
}

// OverlayLogo places a logo in the center of a QR code
func (g *QRGenerator) OverlayLogo(qrCode image.Image, logo image.Image) image.Image {
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
