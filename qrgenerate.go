package thaiqr

import (
	"bytes"
	"fmt"
	qr "github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func GenerateQR(payload string) (*[]byte, error) {
	qrBytes, err := qr.Encode(payload, qr.Highest, 512)

	if err != nil {
		fmt.Println("Failed to encode QR:", err)
		return nil, err
	}
	return &qrBytes, nil
}

func GenerateQRWithThaiQRLogo(payload string) (*[]byte, error) {
	qrCode, err := EncodeThaiQRLogo(payload)
	if err != nil {
		fmt.Println("Failed to encode QR:", err)
		return nil, err
	}

	qrBytes := qrCode.Bytes()

	return &qrBytes, nil
}

func EncodeThaiQRLogo(content string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	file, err := os.Open("assets/thaiqr.png")
	if err != nil {
		fmt.Println("Failed to open logo:", err)
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	logo, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode PNG with logo:", err)
		return nil, err
	}

	code, err := qr.New(content, qr.Highest)
	if err != nil {
		return nil, err
	}

	img := code.Image(512)
	p, ok := img.(*image.Paletted)
	if !ok {
		return nil, fmt.Errorf("undefined qr-code type provided: %T", img)
	}

	if err = png.Encode(&buf, overlayLogo(p, logo)); err != nil {
		return nil, err
	}

	return &buf, nil
}

// overlayLogo - blends logo to the center of the QR code.
func overlayLogo(dst *image.Paletted, logo image.Image) *image.NRGBA {
	res := image.NewNRGBA(dst.Rect)
	draw.Draw(res, res.Bounds(), dst, image.Point{0, 0}, draw.Src)
	offsetX := dst.Bounds().Max.X/2 - logo.Bounds().Max.Y/2
	offsetY := dst.Bounds().Max.Y/2 - logo.Bounds().Max.Y/2
	draw.Draw(res, res.Bounds(), logo, image.Point{-offsetX, -offsetY}, draw.Over)
	return res
}
