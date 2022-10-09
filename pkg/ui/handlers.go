package ui

import (
	"eth-toolkit/pkg/eth"
	"eth-toolkit/pkg/hd"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func handleEnterPress(m UI) (UI, tea.Cmd) {
	if m.state == "new_wallet" || m.state == "get_info_wallet" || m.state == "output" {
		if m.getInState() == "new_hd_wallet_output" {
			mnm := m.output
			m.hdWallet = hd.NewHDWallet(mnm)
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
			walletData := eth.GetWalletFromPK(privateKey)
			loadWalletState(&m, walletData)
		case "mnemonic":
			mnm := m.getInputValue()
			m.hdWallet = hd.NewHDWallet(mnm)
			m.loadHDWallet()
		case "sign_message":
			message := m.getInputValue()
			signedMessage := m.walletData.SignMessage(message)
			setOutputState(
				&m, "Signed Message", signedMessage,
			)
		case "send_tx":
			signedTx := m.getInputValue()
			txHash, err := m.provider.SendSignedTransaction(signedTx)
			var output string
			if err != nil {
				output = fmt.Sprintf("Error: %s", err)
			} else {
				output = fmt.Sprintf("Transaction hash: %s", txHash)
			}
			setOutputState(&m, "Send Transaction", output)
		case "query_bal":
			addr := m.getInputValue()
			balance := m.provider.GetBalance(addr, 0)
			eth_value := eth.GetEthValue(balance)
			output := fmt.Sprintf("Balance is: %v", eth_value)
			setOutputState(&m, "Account Balance", output)
		case "save_keystore":
			password := m.getInputValue()
			keystoreFile := m.walletData.CreateKeystore(password)
			setOutputState(&m, "Keystore file saved", "Path: "+keystoreFile)
		case "update_provider":
			m.provider = eth.GetProvider(m.getInputValue())
			loadWalletState(&m, m.walletData)
		}

	} else if m.state == "sign_transaction" {
		if m.focusIndex == len(m.multiInput) {
			signedTransaction := signTransaction(m)
			setOutputState(&m, "Signed Transaction Hash", signedTransaction)
			m.setMultiInputView()
		}
	} else if m.state == "keystore_access" {
		path := m.multiInput[0].Value()
		password := m.multiInput[1].Value()
		walletData := eth.LoadKeystore(path, password)
		loadWalletState(&m, walletData)
	} else if m.state == "hdwallet" {
		item, ok := m.list.SelectedItem().(ListItem)
		if ok {
			if item.id == "quit" {
				quitToMainMenu(&m)
				m.hdWallet = nil
			} else {
				index, _ := strconv.Atoi(item.id)
				privateKey := m.hdWallet.GetAccount(index).PrivateKey
				loadWalletState(&m, eth.GetWalletFromPK(privateKey))
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
				output, _ := hdwallet.NewMnemonic(128)
				setOutputState(&m, "Mnemonic Words (seperated by space), SAVE IT somewhere safe", output)
				m.setInState("new_hd_wallet_output")
			case "pk":
				setInputState(&m, "Private Key", "Private key", item.id)
			case "sign_message":
				setInputState(&m, "Sign Message", "Message to sign", item.id)
			case "save_keystore":
				setInputState(&m, "Save Keystore", "Password", item.id)
			case "query_bal":
				setInputState(&m, "Query Balance", "Address", item.id)
			case "send_tx":
				setInputState(&m, "Send Transaction", "Signed Transaction Hash", item.id)
			case "provider_options":
				m.loadListItems(getProviderItems(m), "Query Chain")
			case "account_bal":
				balance := m.provider.GetBalance(m.walletData.PublicKey, 0)
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
