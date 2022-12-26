package hd

import (
	"testing"
)

func TestNewHDWallet(t *testing.T) {
	t.Run("Test new HDWallet", func(t *testing.T) {
		mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
		wallet, _ := NewHDWallet(mnemonic)
		if wallet == nil {
			t.Errorf("NewHDWallet() = %v, want %v", wallet, nil)
		}
	})

}

func TestGetAccount(t *testing.T) {
	t.Run("Test get accounts", func(t *testing.T) {
		mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
		wallet, _ := NewHDWallet(mnemonic)
		if wallet == nil {
			t.Errorf("NewHDWallet() = %v, want %v", wallet, nil)
		}

		accounts := wallet.GetAddresses(0, 10)
		if len(accounts) != 10 {
			t.Errorf("GetAddresses() = %v, want %v", len(accounts), 10)
		}
		if accounts[0].Address != "0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947" {
			t.Errorf("GetAddresses() = %v, want %v", accounts[0].Address, "0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947")
		}
		if accounts[1].Address != "0x6F8f46D4b86A623fD5d12A07847008e8Fc7a9A53" {
			t.Errorf("GetAddresses() = %v, want %v", accounts[1].Address, "0x8230645aC28A4EdD1b0B53E7Cd8019744E9dD559")
		}
	})
}

func TestGetAccountByIndex(t *testing.T) {
	t.Run("Test get account by index", func(t *testing.T) {
		mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
		wallet, _ := NewHDWallet(mnemonic)
		if wallet == nil {
			t.Errorf("NewHDWallet() = %v, want %v", wallet, nil)
		}

		accounts := wallet.GetAddresses(0, 10)
		if len(accounts) != 10 {
			t.Errorf("GetAddresses() = %v, want %v", len(accounts), 10)
		}
		expectedAddress := accounts[0].Address
		acc := wallet.GetAccount(0)

		if acc.PublicKey != expectedAddress {
			t.Errorf("GetAccount() = %v, want %v", acc.PublicKey, expectedAddress)
		}

	})
}
