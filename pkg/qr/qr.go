package qr

import (
	"github.com/skip2/go-qrcode"
)

func GenerateQr(text string) *qrcode.QRCode {
	qr, _ := qrcode.New(text, qrcode.Medium)
	return qr
}
