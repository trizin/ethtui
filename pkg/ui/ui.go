package ui

import (
	"ethtui/pkg/eth"
	"ethtui/pkg/hd"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type UI struct {
	list  list.Model
	input textinput.Model

	choice        ListItem
	state         string
	previousState string
	instate       string
	walletData    eth.WalletData
	output        string
	title         string
	hdWallet      *hd.HDWallet
	provider      *eth.Provider

	multiInput []textinput.Model
	focusIndex int
}

func (m UI) Init() tea.Cmd {
	return nil
}

func (m *UI) setListTitle(title string) {
	m.list.Title = title
}

func (m *UI) resetListCursor() {
	m.list.Select(0)
}

func (m *UI) setState(state string) {
	if state == "output" {
		m.input = getText("")
	}
	m.previousState = m.state
	m.state = state
}

func (m *UI) setInState(state string) {
	m.instate = state
}

func (m UI) getInState() string {
	return m.instate
}

func (m *UI) setSendERC20TokensView() {
	m.multiInput = make([]textinput.Model, 3)

	var t textinput.Model
	for i := range m.multiInput {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle

		switch i {
		case 0:
			t.Prompt = "To Address: "
			t.Placeholder = "0x"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "0.01"
			t.Prompt = "Amount: "
		case 2:
			t.Prompt = "Contract Address: "
			t.Placeholder = "0x"
		}

		m.multiInput[i] = t
	}
}

func (m *UI) setMultiInputView() {
	m.multiInput = make([]textinput.Model, 8)

	var t textinput.Model
	for i := range m.multiInput {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle

		switch i {
		case 0:
			t.Prompt = "Nonce: "
			t.Placeholder = "5"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "0x"
			t.CharLimit = 42
			t.Prompt = "To Address: "
		case 2:
			t.Prompt = "Value (ETH): "
			t.Placeholder = "0.01"
			t.CharLimit = 20
		case 3:
			t.Prompt = "Gas Limit: "
			t.Placeholder = "70000"
			t.CharLimit = 20
		case 4:
			t.Prompt = "Gas Price (GWEI): "
			t.Placeholder = "120"
			t.CharLimit = 20
		case 5:
			t.Prompt = "Data: "
			t.Placeholder = "0x"
		case 6:
			t.Prompt = "Chain ID: "
			t.Placeholder = "1"
		case 7:
			t.Prompt = "Priority fee (GWEI): "
			t.Placeholder = "1"
		}

		m.multiInput[i] = t
	}
}

func (m *UI) setMultiInputViewKeystoreFile() {
	m.multiInput = make([]textinput.Model, 2)

	var t textinput.Model
	for i := range m.multiInput {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle

		switch i {
		case 0:
			t.Prompt = "Keystore File Path: "
			t.Placeholder = "./0x.keystore"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "password"
			t.Prompt = "Password: "
			t.EchoCharacter = '⚬'
			t.EchoMode = textinput.EchoPassword
		}

		m.multiInput[i] = t
	}
}

func (m *UI) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.multiInput))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.multiInput {
		m.multiInput[i], cmds[i] = m.multiInput[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *UI) getInputValue() string {
	s := m.input.Value()
	m.input.SetValue("")
	return s
}

func (m *UI) loadListItems(items []list.Item, title string) {
	m.list.SetItems(items)
	m.resetListCursor()
	m.setListTitle(title)
	m.title = title
	m.setState("main")
}

func (m *UI) loadHDWallet() {
	m.loadListItems(getHdWalletItems(m.hdWallet), "HD Wallet Addresses")
	m.setState("hdwallet")
}
