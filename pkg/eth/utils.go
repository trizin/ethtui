package eth

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

func GetEthValue(wei *big.Int) float64 {
	fbalance := new(big.Float)
	fbalance.SetString(wei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	val, _ := ethValue.Float64()
	return val
}

func GetGweiValue(wei *big.Int) float64 {
	fbalance := new(big.Float)
	fbalance.SetString(wei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(9)))
	val, _ := ethValue.Float64()
	return val
}

func GetTransactionInfoString(provider *Provider, txHash string) (string, error) {
	tx, pending, err := provider.GetTransactionInfo(txHash)
	if err != nil {
		return "", err
	}
	output := fmt.Sprintf("Transaction hash: %s\n", txHash)
	output += fmt.Sprintf("To: %s\n", tx.To().String())
	output += fmt.Sprintf("Gas Limit: %d\n", tx.Gas())
	output += fmt.Sprintf("Gas Price: %f GWEI\n", GetGweiValue(tx.GasPrice()))
	output += fmt.Sprintf("Nonce: %d\n", tx.Nonce())
	output += fmt.Sprintf("Value: %f ETH\n", GetEthValue(tx.Value()))
	output += fmt.Sprintf("Pending: %t\n", pending)

	if !pending {
		receipt, err := provider.GetTransactionReceipt(txHash)
		if err != nil {
			return "", err
		}

		output += fmt.Sprintf("Block Hash: %s\n", receipt.BlockHash.String())
		output += fmt.Sprintf("Block Number: %d\n", receipt.BlockNumber)
		output += fmt.Sprintf("Cumulative Gas Used: %d\n", receipt.CumulativeGasUsed)
		output += fmt.Sprintf("Gas Used: %d\n", receipt.GasUsed)
		output += fmt.Sprintf("Status: %d\n", receipt.Status)
		output += fmt.Sprintf("Transaction Index: %d\n", receipt.TransactionIndex)
	}

	return output, nil
}

func GetBlockInfoString(provider *Provider, blockNumber uint64) (string, error) {
	block, err := provider.GetBlockInfo(blockNumber)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Block Number: %d\n", block.Number())
	output += fmt.Sprintf("Block Hash: %s\n", block.Hash().String())
	output += fmt.Sprintf("Block Parent Hash: %s\n", block.ParentHash().String())
	output += fmt.Sprintf("Block Gas Limit: %d\n", block.GasLimit())
	output += fmt.Sprintf("Block Gas Used: %d\n", block.GasUsed())
	output += fmt.Sprintf("Block Size: %s\n", block.Size().String())
	output += fmt.Sprintf("Block Timestamp: %d\n", block.Time())
	t := time.Unix(int64(block.Time()), 0)
	// format time to string 13th Agu 2020 12:00:00
	output += fmt.Sprintf("Block Date: %s\n", t.Format("2nd Jan 2006 15:04:05"))
	output += fmt.Sprintf("Block Transactions: %d\n", len(block.Transactions()))
	return output, nil
}
