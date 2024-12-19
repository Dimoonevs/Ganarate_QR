package genqr

import (
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func GenerateQrCode(data string, logo *os.File) (*image.RGBA, error) {
	qr, err := qrcode.New(data, qrcode.Highest)
	if err != nil {
		return nil, err
	}

	return addLogoInQr(qr.Image(600), logo)
}
func addLogoInQr(qrCodeImage image.Image, logoFile *os.File) (*image.RGBA, error) {
	logoWidth := 250
	logoHeight := 95

	// Декодирование логотипа
	logo, err := jpeg.Decode(logoFile)
	if err != nil {
		return nil, err
	}

	// Масштабирование логотипа до заданных размеров
	resizedLogo := resize.Resize(uint(logoWidth), uint(logoHeight), logo, resize.Lanczos3)

	// Определение смещения для центровки логотипа
	offset := image.Pt((qrCodeImage.Bounds().Dx()-logoWidth)/2, (qrCodeImage.Bounds().Dy()-logoHeight)/2)

	// Создание нового изображения с учетом размеров QR-кода
	b := qrCodeImage.Bounds()
	m := image.NewRGBA(b)

	// Рисование QR-кода и логотипа на новом изображении
	draw.Draw(m, b, qrCodeImage, image.Point{}, draw.Src)
	draw.Draw(m, resizedLogo.Bounds().Add(offset), resizedLogo, image.Point{}, draw.Over)

	return m, nil
}
