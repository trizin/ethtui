package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const erc20ABI = `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

func TransferERC20Tokens(wallet WalletData, contractAddress string, toAddress string, amount *big.Int, provider *Provider) (string, error) {
	// Load the ERC20 token contract ABI
	contractABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return "", fmt.Errorf("failed to load ERC20 token contract ABI: %v", err)
	}

	senderAddress := wallet.PublicKey
	senderAddressEth := common.HexToAddress(senderAddress)

	// Get the nonce for the transaction
	nonce, err := provider.Client.PendingNonceAt(context.Background(), senderAddressEth)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// Create the ERC20 token transfer function call data
	data, err := contractABI.Pack("transfer", common.HexToAddress(toAddress), amount)
	if err != nil {
		return "", fmt.Errorf("failed to pack transfer function call data: %v", err)
	}

	// Create the transaction
	gasLimit, err := provider.GetEstimatedGasUsage([]byte(data))
	if err != nil {
		return "", fmt.Errorf("failed to get estimated gas usage: %v", err)
	}

	gasPrice, err := provider.GetGasPrice()
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}
	gasTipCap, err := provider.GetGasTipCap()
	if err != nil {
		return "", fmt.Errorf("failed to get gas tip cap: %v", err)
	}

	gasPriceFloat := GetGweiValue(gasPrice)
	gasTipCapFloat := GetGweiValue(gasTipCap)

	chainId, err := provider.GetChainId()
	if err != nil {
		return "", fmt.Errorf("failed to get chain id: %v", err)
	}

	dataStr := string(data) // Convert data to a string

	tx, err := wallet.SignTransaction(nonce, contractAddress, 0, int(gasLimit), gasPriceFloat, dataStr, chainId.Int64(), gasTipCapFloat)
	txhash, err := provider.SendSignedTransaction(tx)

	// Print the transaction hash
	log.Printf("Transaction sent: %s", txhash)

	return txhash, nil
}
