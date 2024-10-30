package qr

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/google/uuid"
)

func GenerateQrCode(url string, dimension int, base64OverlayLogo string) (string, error) {
	// Generate QR code
	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	qrCode, err = barcode.Scale(qrCode, dimension, dimension)
	if err != nil {
		return "", err
	}

	// Create image with white background
	img := image.NewRGBA(image.Rect(0, 0, dimension, dimension))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	// Draw QR code onto the image
	qrBounds := qrCode.Bounds()
	qrOffset := image.Pt((img.Bounds().Dx()-qrBounds.Dx())/2, (img.Bounds().Dy()-qrBounds.Dy())/2)
	draw.Draw(img, qrCode.Bounds().Add(qrOffset), qrCode, qrBounds.Min, draw.Over)

	// Add logo image if provided
	if base64OverlayLogo != "" {
		decodedImg, err := base64.StdEncoding.DecodeString(base64OverlayLogo)
		if err != nil {
			return "", err
		}
		logoImg, _, err := image.Decode(bytes.NewReader(decodedImg))
		if err != nil {
			return "", err
		}

		// Overlay the logo on the QR code
		logoBounds := logoImg.Bounds()
		logoOffset := image.Pt((img.Bounds().Dx()-logoBounds.Dx())/2, (img.Bounds().Dy()-logoBounds.Dy())/2)
		draw.Draw(img, logoImg.Bounds().Add(logoOffset), logoImg, logoBounds.Min, draw.Over)
	}

	// Save the final image to a file or return as base64
	fileName := uuid.New().String() + ".png"
	outFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return "", err
	}

	return fileName, nil
}