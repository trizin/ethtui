package main

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

func (w WalletData) signTransaction(
	nonce int,
	toAddress string,
	value float64,
	gasLimit int,
	gasPrice float64,
	data string,
) string {
	tx := types.NewTransaction(
		uint64(nonce),
		common.HexToAddress(toAddress),
		big.NewInt(int64(value*1e18)),
		uint64(gasLimit),
		big.NewInt(int64(gasPrice*1e9)),
		[]byte(data),
	)
	privateKey := w.privateKeyECDSA()

	signedTx, err := types.SignTx(
		tx,
		types.NewEIP155Signer(nil), privateKey,
	)
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	signedTx.EncodeRLP(&buff)

	return hexutil.Encode(buff.Bytes())
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
