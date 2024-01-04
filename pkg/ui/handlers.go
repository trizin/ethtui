package ui

import (
	"ethtui/pkg/eth"
	"ethtui/pkg/hd"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func handleEnterPress(m UI) (UI, tea.Cmd) {
	if m.state == "new_wallet" || m.state == "get_info_wallet" || m.state == "output" {
		if m.getInState() == "new_hd_wallet_output" {
			mnm := m.output
			hdwallet, err := hd.NewHDWallet(mnm)
			if handleError(&m, err) {
				return m, nil
			}
			m.hdWallet = hdwallet
			m.loadHDWallet()
			m.setInState("")
			return m, nil
		}
		m.setState("main")
	} else if m.state == "input" {
		instate := m.getInState()
		m.setInState("")
		switch instate {
		case "pk":
			privateKey := m.getInputValue()
			walletData, err := eth.GetWalletFromPK(privateKey)
			if handleError(&m, err) {
				return m, nil
			}
			loadWalletState(&m, walletData)
		case "mnemonic":
			mnm := m.getInputValue()
			hdwallet, err := hd.NewHDWallet(mnm)
			if handleError(&m, err) {
				return m, nil
			}
			m.hdWallet = hdwallet
			m.loadHDWallet()
		case "sign_message":
			message := m.getInputValue()
			signedMessage, err := m.walletData.SignMessage(message)
			if handleError(&m, err) {
				return m, nil
			}
			setOutputState(
				&m, "Signed Message", signedMessage,
			)
		case "send_tx":
			signedTx := m.getInputValue()
			txHash, err := m.provider.SendSignedTransaction(signedTx)
			if handleError(&m, err) {
				return m, nil
			}
			setOutputState(&m, "Transaction hash", txHash)
		case "query_bal":
			addr := m.getInputValue()
			balance, err := m.provider.GetBalance(addr, 0)
			if handleError(&m, err) {
				return m, nil
			}
			eth_value := eth.GetEthValue(balance)
			output := fmt.Sprintf("Balance is: %v", eth_value)
			setOutputState(&m, "Account Balance", output)
		case "query_tx":
			txHash := m.getInputValue()
			output, err := eth.GetTransactionInfoString(m.provider, txHash)
			if handleError(&m, err) {
				return m, nil
			}
			setOutputState(&m, "Transaction Info", output)

		case "query_block":
			blockNumber, err := strconv.ParseInt(m.getInputValue(), 10, 64)
			if handleError(&m, err) {
				return m, nil
			}
			output, err := eth.GetBlockInfoString(m.provider, uint64(blockNumber))
			if handleError(&m, err) {
				return m, nil
			}

			setOutputState(&m, "Block Info", output)

		case "save_keystore":
			password := m.getInputValue()
			keystoreFile, err := m.walletData.CreateKeystore(password)
			if handleError(&m, err) {
				return m, nil
			}
			setOutputState(&m, "Keystore file saved", "Path: "+keystoreFile)
		case "update_provider":
			p, err := eth.GetProvider(m.getInputValue())
			if handleError(&m, err) {
				return m, nil
			}
			m.provider = p
			loadWalletState(&m, m.walletData)
		}
	} else if m.state == "sign_transaction" {
		if m.focusIndex == len(m.multiInput) {
			signedTransaction, err := signTransaction(m)
			if handleError(&m, err) {
				return m, nil
			}
			setOutputState(&m, "Signed Transaction Hash", signedTransaction)
			m.setMultiInputView()
		}
	} else if m.state == "send_erc20_tokens" {
		if m.focusIndex == len(m.multiInput) {
			txhash, err := sendERC20Tokens(m)
			if handleError(&m, err) {
				return m, nil
			}
			setOutputState(&m, "Transaction Sent", txhash)
			m.setMultiInputView()
		}
	} else if m.state == "keystore_access" {
		path := m.multiInput[0].Value()
		password := m.multiInput[1].Value()
		walletData, err := eth.LoadKeystore(path, password)
		if handleError(&m, err) {
			return m, nil
		}
		loadWalletState(&m, walletData)
	} else if m.state == "hdwallet" {
		item, ok := m.list.SelectedItem().(ListItem)
		if ok {
			if item.id == "quit" {
				quitToMainMenu(&m)
				m.hdWallet = nil
			} else {
				index, _ := strconv.Atoi(item.id)
				walletData := m.hdWallet.GetAccount(index)
				loadWalletState(&m, walletData)
			}
		}

	} else if m.state == "main" {
		item, ok := m.list.SelectedItem().(ListItem)
		if ok {
			m.setState(item.id)
			if m.state == "quit" {
				quitToMainMenu(&m)
			}

			switch item.id {
			case "send_erc20_tokens":
				m.setSendERC20TokensView()
			case "sign_transaction":
				m.setMultiInputView()
			case "keystore_access":
				m.setMultiInputViewKeystoreFile()
			case "mnemonic":
				setInputState(&m, "Mnemonic Words (seperated by space)", "airport loud mixture", "mnemonic")
			case "access_wallet":
				m.loadListItems(getAccessWalletItems(), "Access Wallet")
			case "new_wallet":
				walletData := eth.GenerateWallet()
				loadWalletState(&m, walletData)
			case "public_key":
				output := displayWalletPublicKey(m.walletData)
				setOutputState(&m, "Public Key", output)
			case "private_key":
				output := displayWalletPrivateKey(m.walletData)
				setOutputState(&m, "Private Key", output)
			case "new_hd_wallet":
				size, _ := strconv.Atoi(item.tmp)
				output, _ := hdwallet.NewMnemonic(size)
				setOutputState(&m, "Mnemonic Words (seperated by space), SAVE IT somewhere safe", output)
				m.setInState("new_hd_wallet_output")
			case "new_hd_wallet_pick":
				m.loadListItems(getHDWalletChooseItems(m), "Choose Mnemonic Size")
			case "pk":
				setInputState(&m, "Private Key", "Private key", item.id)
			case "sign_message":
				setInputState(&m, "Sign Message", "Message to sign", item.id)
			case "save_keystore":
				setInputState(&m, "Save Keystore", "Password", item.id)
			case "query_bal":
				setInputState(&m, "Query Balance", "Address", item.id)
			case "query_tx":
				setInputState(&m, "Query Transaction", "Transaction Hash", item.id)
			case "query_block":
				setInputState(&m, "Query Block", "Block Number", item.id)
			case "send_tx":
				setInputState(&m, "Send Transaction", "Signed Transaction Hash", item.id)
			case "provider_options":
				m.loadListItems(getProviderItems(m), "Query Chain")
			case "account_bal":
				balance, err := m.provider.GetBalance(m.walletData.PublicKey, 0)
				if handleError(&m, err) {
					return m, nil
				}
				eth_value := eth.GetEthValue(balance)
				output := fmt.Sprintf("Balance is: %v", eth_value)
				setOutputState(&m, "Account Balance", output)
			case "back":
				loadWalletState(&m, m.walletData)
			}
			m.choice = item
		}
	}

	return m, nil
}
