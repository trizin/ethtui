package main

import (
	"github.com/skip2/go-qrcode"
)

func generateQr(text string) *qrcode.QRCode {
	qr, _ := qrcode.New(text, qrcode.Medium)
	return qr
}
