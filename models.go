package main

import "github.com/skip2/go-qrcode"

type WalletData struct {
	PrivateKey   string
	PublicKey    string
	PrivateKeyQR *qrcode.QRCode
	PublicKeyQR  *qrcode.QRCode
}
