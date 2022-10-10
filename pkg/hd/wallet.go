package hd

import (
	"eth-toolkit/pkg/eth"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type HDWallet struct {
	wallet *hdwallet.Wallet
}

type HDAccount struct {
	Index   int
	Address string
}

func (w *HDWallet) getPathByIndex(index int) accounts.Account {
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", index))
	account, err := w.wallet.Derive(path, true)
	if err != nil {
		panic(err)
	}
	return account
}

func (w *HDWallet) GetAccount(index int) eth.WalletData {
	account := w.getPathByIndex(index)

	privateKeyHex, err := w.wallet.PrivateKeyHex(account)
	if err != nil {
		panic(err)
	}

	wallet := eth.GetWalletFromPK(privateKeyHex)
	return wallet
}

func (w *HDWallet) GetAddresses(
	startIndex int,
	endIndex int,
) []HDAccount {
	addresses := []HDAccount{}

	for i := startIndex; i < endIndex; i++ {
		addresses = append(addresses, HDAccount{
			Index:   i,
			Address: w.getPathByIndex(i).Address.Hex(),
		})
	}
	return addresses
}

func NewHDWallet(mnemonic string) (*HDWallet, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return &HDWallet{wallet}, nil
}
