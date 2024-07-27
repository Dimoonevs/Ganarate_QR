package genqr

import (
	"github.com/skip2/go-qrcode"
)

func GenerateQrCode(data string) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	return qr.PNG(600)
}
