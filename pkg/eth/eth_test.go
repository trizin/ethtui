package eth

import (
	"testing"
)

func TestWalletData_SignTransaction(t *testing.T) {
	wallet := GenerateWallet()
	signed := wallet.SignTransaction(
		0,
		"0x000000000000000000000000000000000000dEaD",
		10,
		10000,
		100,
		"0x",
		1,
		2,
	)

	if signed == "" {
		t.Errorf("WalletData.SignTransaction() = %v", signed)
	}
}
