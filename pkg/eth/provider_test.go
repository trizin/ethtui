package eth

import (
	"math/big"
	"testing"
)

var rpcUrl = "https://rpc.ankr.com/eth"

func TestProvider_GetBalance(t *testing.T) {

	provider := GetProvider(rpcUrl)
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

func TestGetTransactionReceipt(t *testing.T) {
	provider := GetProvider(rpcUrl)
	txHash := "0x82237a9d319cbb9a46d1bbbdbac870918e70ae9f0350db24dc578c1a5cf4d859"

	receipt, err := provider.GetTransactionReceipt(txHash)

	if err != nil {
		t.Errorf("Provider.GetTransactionReceipt() error = %v", err)
	}

	if receipt == nil {
		t.Errorf("Provider.GetTransactionReceipt() receipt = %v", receipt)
	}

	if receipt.Status != 1 {
		t.Errorf("Provider.GetTransactionReceipt() status = %v", receipt.Status)
	}

	if receipt.BlockNumber.Cmp(big.NewInt(15701530)) != 0 {
		t.Errorf("Provider.GetTransactionReceipt() block number = %v", receipt.BlockNumber)
	}

	if receipt.TransactionIndex != 179 {
		t.Errorf("Provider.GetTransactionReceipt() transaction index = %v", receipt.TransactionIndex)
	}

	if receipt.GasUsed != 21000 {
		t.Errorf("Provider.GetTransactionReceipt() gas used = %v", receipt.GasUsed)
	}
}

func TestGetTransactionInfo(t *testing.T) {
	provider := GetProvider(rpcUrl)
	txHash := "0x82237a9d319cbb9a46d1bbbdbac870918e70ae9f0350db24dc578c1a5cf4d859"

	tx, pending, err := provider.GetTransactionInfo(txHash)

	if err != nil {
		t.Errorf("Provider.GetTransactionInfo() error = %v", err)
	}

	if pending != false {
		t.Errorf("Provider.GetTransactionInfo() success = %v", pending)
	}

	if tx == nil {
		t.Errorf("Provider.GetTransactionInfo() tx = %v", tx)
	}

	if tx.Hash().String() != txHash {
		t.Errorf("Provider.GetTransactionInfo() hash = %v", tx.Hash().String())
	}

	if tx.Value().Cmp(big.NewInt(53146227074366)) != 0 {
		t.Errorf("Provider.GetTransactionInfo() value = %v", tx.Value())
	}

	if tx.Gas() != 21000 {
		t.Errorf("Provider.GetTransactionInfo() gas = %v", tx.Gas())
	}

	if tx.GasPrice().Cmp(big.NewInt(5664994449)) != 0 {
		t.Errorf("Provider.GetTransactionInfo() gas price = %v", tx.GasPrice())
	}

	if tx.Nonce() != 21 {
		t.Errorf("Provider.GetTransactionInfo() nonce = %v", tx.Nonce())
	}

	if tx.To().String() != "0x000000000000000000000000000000000000dEaD" {
		t.Errorf("Provider.GetTransactionInfo() to = %v", tx.To().String())
	}

	if tx.Data() == nil {
		t.Errorf("Provider.GetTransactionInfo() data = %v", tx.Data())
	}
}

func TestGetNonce(t *testing.T) {
	provider := GetProvider(rpcUrl)
	addr := "0x000000000000000000000000000000000000dEaD"

	expected := uint64(0)
	got := provider.GetNonce(addr)
	if got != expected {
		t.Errorf("Provider.GetNonce() = %v, want %v", got, expected)
	}
}

func TestSignAndSendTransaction(t *testing.T) {
	pk := "0x7477652b0d4f24e0b5cfdc60f49e4f58deb7c8781cdf92079b5cb17515615de7"
	wallet := GetWalletFromPK(pk)
	provider := GetProvider("http://localhost:8545")
	addr := "0x000000000000000000000000000000000000dEaD"
	sender := wallet.PublicKey

	signedTx := wallet.SignTransaction(
		int(provider.GetNonce(sender)),
		addr,
		1.0,
		90000,
		20.0,
		"0x",
		provider.GetChainId().Int64(),
	)

	balSender := provider.GetBalance(sender, 0)
	if balSender.Cmp(big.NewInt(0)) != 1 {
		t.Errorf("Provider.GetBalance() not enough balance")
		return
	}

	beforebal := provider.GetBalance(addr, 0)
	txHash, err := provider.SendSignedTransaction(signedTx)

	if err != nil {
		t.Errorf("Provider.SendTransaction() error = %v", err)
		return
	}

	afterbal := provider.GetBalance(addr, 0)

	if txHash == "" {
		t.Errorf("Provider.SendTransaction() tx hash = %v", txHash)
	}

	if beforebal.Cmp(afterbal) != -1 {
		t.Errorf("Provider.SendTransaction() balance = %v", afterbal)
	}

}
