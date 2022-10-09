package ui

import (
	"fmt"
	"strconv"
	"strings"

	"eth-toolkit/pkg/eth"
	"eth-toolkit/pkg/hd"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func GetUI() UI {
	m := UI{title: "✨✨✨", list: list.NewModel(getMainItems(), list.NewDefaultDelegate(), 0, 0), input: getText(""), state: "main"}
	m.list.Title = "✨✨✨"
	return m
}

func (m UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+p":
			if m.state == "main" {
				requestProvider(&m)
			}

		case "tab", "shift+tab", "up", "down":
			if m.state == "sign_transaction" || m.state == "keystore_access" || m.state == "mnemonic" {
				s := msg.String()
				m, cmds := moveIndex(m, s)
				return m, tea.Batch(cmds...)
			}

		case "enter":
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
					m.setState("hdwallet")
					m.setListTitle("HD Wallet Addresses")
					m.list.SetItems(getHdWalletItems(m.hdWallet))
					m.resetListCursor()
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
					m.title = "Mnemonic Words (seperated by space)"
					setInputState(&m, "Mnemonic Words (seperated by space)", "airport loud mixture")
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
					setInputState(&m, "Private Key", "Private key")
				case "sign_message":
					setInputState(&m, "Sign Message", "Message to sign")
				case "save_keystore":
					setInputState(&m, "Save Keystore", "Password")
				case "provider_options":
					m.loadListItems(getProviderItems(m), "Query Chain")
				case "account_bal":
					m.title = "Account Balance"
					balance := m.provider.GetBalance(m.walletData.PublicKey, 0)
					eth_value := eth.GetEthValue(balance)
					output := fmt.Sprintf("Balance is: %v", eth_value)
					setOutputState(&m, "Account Balance", output)
				case "query_bal":
					setInputState(&m, "Query Balance", "Address")
				case "send_tx":
					setInputState(&m, "Send Transaction", "Signed Transaction Hash")
				case "back":
					loadWalletState(&m, m.walletData)
				}

				if ok {
					m.choice = item
				}
				m.setInState(item.id)
			}
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		docStyle.Width(msg.Width)
		m.input.Width = int(float64(msg.Width*5) / 6)
	}

	var cmd tea.Cmd

	if m.state == "main" || m.state == "hdwallet" {
		m.list, cmd = m.list.Update(msg)
	}

	if m.state == "input" {
		m.input, cmd = m.input.Update(msg)
	}

	if m.state == "sign_transaction" || m.state == "keystore_access" || m.state == "mnemonic" {
		cmd = m.updateInputs(msg)
	}

	return m, cmd
}

func (m UI) View() string {

	if m.choice.title != "" {
		switch m.state {

		case "sign_transaction":
			var b strings.Builder
			for i := range m.multiInput {
				b.WriteString(m.multiInput[i].View())
				if i < len(m.multiInput)-1 {
					b.WriteRune('\n')
				}
			}

			button := &blurredButton
			if m.focusIndex == len(m.multiInput) {
				button = &focusedButton
			}
			fmt.Fprintf(&b, "\n\n%s\n\n", *button)

			return docStyle.Render(
				fmt.Sprintf(
					"%s\n\n%s",
					"Sign Transaction",
					b.String(),
				))

		case "keystore_access":
			var b strings.Builder
			for i := range m.multiInput {
				b.WriteString(m.multiInput[i].View())
				if i < len(m.multiInput)-1 {
					b.WriteRune('\n')
				}
			}

			button := &blurredButton
			if m.focusIndex == len(m.multiInput) {
				button = &focusedButton
			}
			fmt.Fprintf(&b, "\n\n%s\n\n", *button)

			return b.String()

		case "input":
			return docStyle.Render(fmt.Sprintf(
				"%s\n%s\n%s",
				titleStyle.Render(m.title),
				m.input.View(),
				blurredStyle.Render("Press ctrl+c to quit"),
			))

		case "output":
			in := fmt.Sprintf(
				"%s\n%s\n%s",
				titleStyle.Render(m.title),
				docStyle.Render(m.output),
				blurredStyle.Render("Press enter to continue"),
			)

			return docStyle.Render(in)
		}
	}

	return docStyle.Render(m.list.View())
}
