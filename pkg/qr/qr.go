package qr

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/google/uuid"
)

func GenerateQrCode(url string, dimension int, base64OverlayLogo string) (string, error) {
	// Generate QR code
	qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, dimension, dimension)

	// Create image with white background
	img := image.NewRGBA(image.Rect(0, 0, dimension, dimension))
	white := color.RGBA{uint8(dimension), uint8(dimension), uint8(dimension), uint8(dimension)}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	// Draw QR code onto the image
	qrBounds := qrCode.Bounds()
	qrOffset := image.Pt((img.Bounds().Dx()-qrBounds.Dx())/2, (img.Bounds().Dy()-qrBounds.Dy())/2)
	draw.Draw(img, qrCode.Bounds().Add(qrOffset), qrCode, qrBounds.Min, draw.Over)

	// Add logo image if provided
	if base64OverlayLogo != "" {
		// Decode base64 logo image
		decodedImg, err := base64.StdEncoding.DecodeString(base64OverlayLogo)
		if err != nil {
			log.Println(err)
			return "", err
		}

		// Decode the logo image
		logoImg, _, err := image.Decode(bytes.NewReader(decodedImg))
		if err != nil {
			log.Println(err)
			return "", err
		}

		// Overlay the logo on the QR code
		logoBounds := logoImg.Bounds()
		logoOffset := image.Pt((img.Bounds().Dx()-logoBounds.Dx())/2, (img.Bounds().Dy()-logoBounds.Dy())/2)
		draw.Draw(img, logoImg.Bounds().Add(logoOffset), logoImg, logoBounds.Min, draw.Over)
	}

	// Save the final image to a file
	fileName := uuid.New().String() + ".png"
	outFile, _ := os.Create(fileName)
	defer outFile.Close()
	if err := png.Encode(outFile, img); err != nil {
		log.Println(err)
		return "", err
	}

	return fileName, nil
}
