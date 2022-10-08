package eth

import (
	"testing"
)

func TestProvider_GetBalance(t *testing.T) {

	provider := GetProvider("https://cloudflare-eth.com")
	addr := "0x000000000000000000000000000000000000dEaD"

	t.Run("Get balance at block", func(t *testing.T) {
		expected := "12567984693887489302095"
		got := provider.GetBalance(addr, 15705799).String()
		if got != expected {
			t.Errorf("Provider.GetBalance() = %v, want %v", got, expected)
		}
	})

}
