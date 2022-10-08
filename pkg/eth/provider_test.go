package eth

import (
	"math/big"
	"testing"
)

func TestProvider_GetBalance(t *testing.T) {

	provider := GetProvider("https://cloudflare-eth.com")
	addr := "0x000000000000000000000000000000000000dEaD"

	t.Run("Get balance at block", func(t *testing.T) {
		expected := "12567984693887489302095"
		got := provider.GetBalance(addr, 15705799).String()
		if got != expected {
			t.Errorf("Provider.GetBalance() at block = %v, want %v", got, expected)
		}
	})
	t.Run("Get last balance", func(t *testing.T) {
		got := provider.GetBalance(addr, 0)
		if got.Cmp(big.NewInt(0)) != 1 {
			t.Errorf("Provider.GetBalance() last balance failed")
		}
	})
}

func TestGetEthValue(t *testing.T) {
	wei := big.NewInt(1100000000000000000)
	expected := 1.1
	got := GetEthValue(wei)
	if got != expected {
		t.Errorf("GetEthValue() = %v, want %v", got, expected)
	}
}
