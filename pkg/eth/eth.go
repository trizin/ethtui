package eth

import (
	"bytes"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"

	"eth-toolkit/pkg/qr"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/skip2/go-qrcode"
)

type WalletData struct {
	PrivateKey   string
	PublicKey    string
	PrivateKeyQR *qrcode.QRCode
	PublicKeyQR  *qrcode.QRCode
}

func GetWalletFromPK(pk string) WalletData {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(err)
	}
	return GetWalletDataFromPKECDSA(privateKey)
}

func GenerateWallet() WalletData {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return GetWalletDataFromPKECDSA(privateKey)
}

func (w WalletData) PrivateKeyECDSA() *ecdsa.PrivateKey {
	data, _ := crypto.HexToECDSA(w.PrivateKey)
	return data
}

func (w WalletData) CreateKeystore(password string) string {
	fileName := "./" + w.PublicKey + ".keystore"
	ks := keystore.NewKeyStore(fileName, keystore.StandardScryptN, keystore.StandardScryptP)
	_, err := ks.ImportECDSA(w.PrivateKeyECDSA(), password)
	if err != nil {
		panic(err)
	}
	return fileName
}

func LoadKeystore(path string, password string) WalletData {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := ks.Accounts()
	if len(accounts) == 0 {
		panic("No accounts found in keystore")
	}
	account := accounts[0]

	keyjson, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		panic(err)
	}

	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		panic(err)
	}

	if key.Address != account.Address {
		panic("Key address mismatch")
	}

	return GetWalletDataFromPKECDSA(key.PrivateKey)

}

func (w WalletData) SignTransaction(
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
	privateKey := w.PrivateKeyECDSA()

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

func (w WalletData) SignMessage(dataString string) string {
	data := []byte(dataString)
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), w.PrivateKeyECDSA())
	if err != nil {
		panic(err)
	}
	return hexutil.Encode(signature)
}

func GetWalletDataFromPKECDSA(privateKey *ecdsa.PrivateKey) WalletData {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := (hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyHex := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	walletData := WalletData{
		PrivateKey:   privateKeyHex,
		PublicKey:    publicKeyHex,
		PrivateKeyQR: qr.GenerateQr(privateKeyHex),
		PublicKeyQR:  qr.GenerateQr(publicKeyHex),
	}

	return walletData
}
