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
					m.output = ""
					m.hdWallet = hd.NewHDWallet(mnm)
					m.setState("hdwallet")
					m.setListTitle("HD Wallet Addresses")
					m.list.SetItems(getHdWalletItems(m.hdWallet))
					m.resetListCursor()
				} else {
					m.setState("main")
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

			} else if m.state == "pk" {
				privateKey := m.input.Value()
				m.input.SetValue("")
				walletData := eth.GetWalletFromPK(privateKey)
				loadWalletState(&m, walletData)
			} else if m.state == "mnemonic" {
				mnm := m.input.Value()
				m.input.SetValue("")
				m.hdWallet = hd.NewHDWallet(mnm)
				m.setState("hdwallet")
				m.setListTitle("HD Wallet Addresses")
				m.list.SetItems(getHdWalletItems(m.hdWallet))
				m.resetListCursor()
			} else if m.state == "sign_message" {
				message := m.input.Value()
				signedMessage := m.walletData.SignMessage(message)
				setOutputState(
					&m, "Signed Message", signedMessage,
				)
			} else if m.state == "send_tx" {
				signedTx := m.input.Value()
				txHash, err := m.provider.SendSignedTransaction(signedTx)
				var output string
				if err != nil {
					output = fmt.Sprintf("Error: %s", err)
				} else {
					output = fmt.Sprintf("Transaction hash: %s", txHash)
				}
				setOutputState(&m, "Send Transaction", output)
			} else if m.state == "query_bal" {
				addr := m.input.Value()
				balance := m.provider.GetBalance(addr, 0)
				eth_value := eth.GetEthValue(balance)
				output := fmt.Sprintf("Balance is: %v", eth_value)
				setOutputState(&m, "Account Balance", output)

			} else if m.state == "save_keystore" {
				password := m.input.Value()
				keystoreFile := m.walletData.CreateKeystore(password)
				setOutputState(&m, "Keystore file saved", "Path: "+keystoreFile)
			} else if m.state == "hdwallet" {
				item, ok := m.list.SelectedItem().(ListItem)
				if ok {
					if item.id == "quit" {
						m.list.SetItems(getMainItems())
						m.resetListCursor()
						m.setState("main")
						m.setListTitle("✨✨✨")
						m.hdWallet = nil
					} else {
						index, _ := strconv.Atoi(item.id)
						privateKey := m.hdWallet.GetAccount(index).PrivateKey
						m.walletData = eth.GetWalletFromPK(privateKey)
						m.setState("main")
						m.list.SetItems(getControlWalletItems(m))
						m.resetListCursor()
						m.setListTitle(m.walletData.PublicKey)
					}
				}

			} else if m.state == "main" || m.state == "access_wallet" {
				item, ok := m.list.SelectedItem().(ListItem)

				m.setState(item.id)
				switch item.id {
				case "sign_transaction":
					m.setMultiInputView()
				case "keystore_access":
					m.setMultiInputViewKeystoreFile()
				case "mnemonic":
					m.title = "Mnemonic Words (seperated by space)"
				case "access_wallet":
					m.list.SetItems(getAccessWalletItems())
					m.resetListCursor()
					m.setListTitle("Access Wallet")
				case "new_wallet":
					walletData := eth.GenerateWallet()
					m.walletData = walletData
					m.setState("main")
					m.list.SetItems(getControlWalletItems(m))
					m.resetListCursor()
					m.input = getText("")
					m.setListTitle(m.walletData.PublicKey)
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
					m.setInState(item.id)
				case "sign_message":
					setInputState(&m, "Sign Message", "Message to sign")
					m.setInState(item.id)
				case "save_keystore":
					setInputState(&m, "Save Keystore", "Password")
					m.setInState(item.id)
				case "provider_options":
					m.title = "Query Chain"
					m.list.SetItems(getProviderItems(m))
					m.resetListCursor()
					m.setState("main")
				case "account_bal":
					m.title = "Account Balance"
					balance := m.provider.GetBalance(m.walletData.PublicKey, 0)
					eth_value := eth.GetEthValue(balance)
					output := fmt.Sprintf("Balance is: %v", eth_value)
					setOutputState(&m, "Account Balance", output)
				case "query_bal":
					setInputState(&m, "Query Balance", "Address")
					m.setInState(item.id)
				case "send_tx":
					setInputState(&m, "Send Transaction", "Signed Transaction Hash")
					m.setInState(item.id)
				case "back":
					m.setState("main")
					m.list.SetItems(getControlWalletItems(m))
					m.resetListCursor()
				}

				if m.state == "quit" {
					m.list.SetItems(getMainItems())
					m.resetListCursor()
					m.setState("main")
					m.setListTitle("✨✨✨")
				}

				if ok {
					m.choice = item
				}
			} else if m.state == "update_provider" {
				m.provider = eth.GetProvider(m.input.Value())
				m.setState("main")
				m.input.SetValue("")
				m.list.SetItems(getControlWalletItems(m))
				m.resetListCursor()
			}
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		docStyle.Width(msg.Width)
		m.input.Width = int(float64(msg.Width*5) / 6)
	}

	var cmd tea.Cmd

	if m.state == "main" || m.state == "access_wallet" || m.state == "hdwallet" {
		m.list, cmd = m.list.Update(msg)
	}

	if m.state == "pk" || m.state == "sign_message" || m.state == "save_keystore" || m.state == "keystore_access" || m.state == "mnemonic" || m.state == "update_provider" || m.state == "send_tx" || m.state == "query_bal" {
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

		case "save_keystore", "pk", "sign_message", "mnemonic", "update_provider", "send_tx", "query_bal":
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
