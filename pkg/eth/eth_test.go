package eth

import (
	"testing"
)

func TestWalletData_SignTransaction(t *testing.T) {
	wallet := GenerateWallet()
	signed, err := wallet.SignTransaction(
		0,
		"0x000000000000000000000000000000000000dEaD",
		10,
		10000,
		100,
		"0x",
		1,
		2,
	)

	if err != nil {
		t.Errorf("WalletData.SignTransaction() error = %v", err)
		return
	}

	if signed == "" {
		t.Errorf("WalletData.SignTransaction() = %v", signed)
	}
}
