package main

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func getWalletFromPK(pk string) WalletData {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(err)
	}
	return getWalletDataFromPKECDSA(privateKey)
}

func generateWallet() WalletData {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return getWalletDataFromPKECDSA(privateKey)
}

func (w WalletData) privateKeyECDSA() *ecdsa.PrivateKey {
	data, _ := crypto.HexToECDSA(w.PrivateKey)
	return data
}

func (w WalletData) signMessage(dataString string) string {
	data := []byte(dataString)
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), w.privateKeyECDSA())
	if err != nil {
		panic(err)
	}
	return hexutil.Encode(signature)
}

func getWalletDataFromPKECDSA(privateKey *ecdsa.PrivateKey) WalletData {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := (hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyHex := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	walletData := WalletData{
		PrivateKey:   privateKeyHex,
		PublicKey:    publicKeyHex,
		PrivateKeyQR: generateQr(privateKeyHex),
		PublicKeyQR:  generateQr(publicKeyHex),
	}

	return walletData
}
