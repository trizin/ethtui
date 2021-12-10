package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListItem struct {
	title string
	desc  string
	id    string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }

type UI struct {
	list  list.Model
	input textinput.Model

	choice     ListItem
	state      string
	inputText  string
	walletData WalletData
	output     string
}

func (m UI) Init() tea.Cmd {
	return nil
}

func (m UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":

			if m.state == "new_wallet" || m.state == "get_info_wallet" || m.state == "output" {
				m.state = "main"
			} else if m.state == "pk" {
				privateKey := m.input.Value()
				m.walletData = getWalletFromPK(privateKey)
				m.state = "main"
				m.list.SetItems(getControlWalletItems())
				m.input = getText()

			} else if m.state == "sign_message" {
				message := m.input.Value()
				signedMessage := m.walletData.signMessage(message)
				m.output = signedMessage
				m.state = "output"
				m.input = getText()
			} else if m.state == "main" || m.state == "access_wallet" {
				item, ok := m.list.SelectedItem().(ListItem)

				if item.id == "access_wallet" {
					m.list.SetItems(getAccessWalletItems())
				}

				m.state = item.id

				if m.state == "quit" {
					return m, tea.Quit
				}

				if ok {
					m.choice = item
				}
			}
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd

	if m.state == "main" || m.state == "access_wallet" {
		m.list, cmd = m.list.Update(msg)
	}

	if m.state == "pk" || m.state == "sign_message" {
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func getMainItems() []list.Item {
	items := []list.Item{
		ListItem{title: "New Wallet", desc: "Create a new wallet", id: "new_wallet"},
		ListItem{title: "Access Wallet", desc: "Access an existing wallet", id: "access_wallet"},
	}
	return items
}

func getAccessWalletItems() []list.Item {
	items := []list.Item{
		ListItem{title: "Private Key", desc: "Access your wallet using your private key", id: "pk"},
		ListItem{title: "JSON File", desc: "Access a wallet using your keystore file", id: "json"},
	}
	return items
}

func getControlWalletItems() []list.Item {
	items := []list.Item{
		ListItem{title: "Get Info", desc: "Display addresses, keys and qr codes", id: "get_info_wallet"},
		ListItem{title: "Sign Message", desc: "Sign a message with the private key", id: "sign_message"},
		ListItem{title: "Sign Transaction", desc: "Sign a transaction with the private key", id: "sign_transaction"},
		ListItem{title: "Quit", desc: "Quit to main menu", id: "quit"},
	}
	return items
}

func getText() textinput.Model {
	ti := textinput.NewModel()
	ti.BlinkSpeed = 1
	ti.Placeholder = "Private Key"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	return ti
}

func getUI() UI {
	m := UI{list: list.NewModel(getMainItems(), list.NewDefaultDelegate(), 0, 0), input: getText(), state: "main"}
	return m
}

func (m UI) setState(state string) {
	m.state = state
}

func displayWalletInfo(walletData WalletData) string {
	return fmt.Sprintf(
		"New ETH Wallet\n%s%s\n%s\n%s\n%s",
		walletData.PublicKeyQR.ToSmallString(true), walletData.PrivateKeyQR.ToSmallString(true),
		"Private Key: "+walletData.PrivateKey,
		"Public Key: "+walletData.PublicKey,
		"Press enter to go back",
	) + "\n"
}

func (m UI) View() string {

	if m.choice.title != "" {
		switch m.state {
		case "new_wallet":
			walletData := generateWallet()
			return docStyle.Render(displayWalletInfo(walletData))

		case "get_info_wallet":
			return docStyle.Render(displayWalletInfo(m.walletData))

		case "pk":
			return docStyle.Render(fmt.Sprintf(
				"Private Key\n%s\n%s",
				m.input.View(),
				"Press ctrl+c to quit",
			))
		case "sign_message":
			return docStyle.Render(fmt.Sprintf(
				"Message to sign: \n%s\n%s",
				m.input.View(),
				"Press ctrl+c to quit",
			))
		case "output":
			return docStyle.Render(fmt.Sprintf(
				"Output \n%s\n%s",
				m.output,
				"Press enter to continue",
			))
		}
	}

	return docStyle.Render(m.list.View())
}
