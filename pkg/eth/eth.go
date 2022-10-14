package eth

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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

func GetWalletFromPK(pk string) (WalletData, error) {
	// check if starts with 0x
	if pk[:2] == "0x" {
		// remove 0x
		pk = pk[2:]
	}

	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return WalletData{}, err
	}
	return GetWalletDataFromPKECDSA(privateKey), nil
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

func (w WalletData) CreateKeystore(password string) (string, error) {
	fileName := "./" + w.PublicKey + ".keystore"
	ks := keystore.NewKeyStore(fileName, keystore.StandardScryptN, keystore.StandardScryptP)
	_, err := ks.ImportECDSA(w.PrivateKeyECDSA(), password)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func LoadKeystore(path string, password string) (WalletData, error) {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := ks.Accounts()
	if len(accounts) == 0 {
		return WalletData{}, errors.New("no accounts found in keystore")
	}
	account := accounts[0]

	keyjson, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		return WalletData{}, err
	}

	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		return WalletData{}, err
	}

	if key.Address != account.Address {
		return WalletData{}, errors.New("address mismatch")
	}

	return GetWalletDataFromPKECDSA(key.PrivateKey), nil

}

func (w WalletData) SignTransaction(
	nonce uint64,
	toAddress string,
	value float64,
	gasLimit int,
	gasPrice float64,
	data string,
	chainId int64,
	gasTipCap float64,
) (string, error) {
	addr := common.HexToAddress(toAddress)
	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   big.NewInt(chainId),
			Nonce:     uint64(nonce),
			To:        &addr,
			Value:     big.NewInt(int64(value * 1e18)),
			GasFeeCap: big.NewInt(int64(gasPrice * 1e9)),
			GasTipCap: big.NewInt(int64(gasTipCap * 1e9)),
			Gas:       uint64(gasLimit),
			Data:      []byte(data),
		},
	)
	privateKey := w.PrivateKeyECDSA()

	signedTx, err := types.SignTx(
		tx,
		types.LatestSignerForChainID(big.NewInt(chainId)), privateKey,
	)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	signedTx.EncodeRLP(&buff)

	return hexutil.Encode(buff.Bytes()), nil
}

func (w WalletData) SignMessage(dataString string) (string, error) {
	data := []byte(dataString)
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), w.PrivateKeyECDSA())
	if err != nil {
		return "", err
	}
	return hexutil.Encode(signature), nil
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
