package genqr

import (
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func GenerateQrCode(data string, logo *os.File) (*image.RGBA, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	return addLogoInQr(qr.Image(600), logo)
}
func addLogoInQr(qrCodeImage image.Image, logoFile *os.File) (*image.RGBA, error) {

	logo, err := png.Decode(logoFile)
	if err != nil {
		return nil, err
	}

	logoSize := 64
	resizedLogo := resize.Resize(uint(logoSize), uint(logoSize), logo, resize.Lanczos3)

	offset := image.Pt((qrCodeImage.Bounds().Dx()-logoSize)/2, (qrCodeImage.Bounds().Dy()-logoSize)/2)
	b := qrCodeImage.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, qrCodeImage, image.Point{}, draw.Src)
	draw.Draw(m, resizedLogo.Bounds().Add(offset), resizedLogo, image.Point{}, draw.Over)
	return m, nil
}
