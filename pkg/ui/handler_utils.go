package ui

import (
	"eth-toolkit/pkg/eth"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func requestProvider(m *UI) {
	setInputState(m, "Set Provider", "Enter provider URL", "update_provider")
}

func moveIndex(m UI, s string) (UI, []tea.Cmd) {
	// Cycle indexes
	if s == "up" || s == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > len(m.multiInput) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.multiInput)
	}

	cmds := make([]tea.Cmd, len(m.multiInput))
	for i := 0; i <= len(m.multiInput)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.multiInput[i].Focus()
			m.multiInput[i].PromptStyle = focusedStyle
			m.multiInput[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.multiInput[i].Blur()
		m.multiInput[i].PromptStyle = noStyle
		m.multiInput[i].TextStyle = noStyle
	}

	return m, cmds
}

func signTransaction(m UI) string {
	nonce, _ := strconv.Atoi(m.multiInput[0].Value())
	toAddress := m.multiInput[1].Value()
	value, _ := strconv.ParseFloat(m.multiInput[2].Value(), 64)
	gasLimit, _ := strconv.Atoi(m.multiInput[3].Value())
	gasPrice, _ := strconv.ParseFloat(m.multiInput[4].Value(), 64)
	data := m.multiInput[5].Value()
	chainId, _ := strconv.Atoi(m.multiInput[6].Value())
	gasTipCap, _ := strconv.ParseFloat(m.multiInput[7].Value(), 64)

	signedTransaction := m.walletData.SignTransaction(uint64(nonce), toAddress, value, gasLimit, gasPrice, data, int64(chainId), gasTipCap)
	return signedTransaction
}

func setOutputState(m *UI, title string, output string) {
	m.setState("output")
	m.title = title
	m.output = output
}

func setInputState(m *UI, title string, placeholder string, instate string) {
	m.setState("input")
	m.setInState(instate)
	m.title = title
	m.input = getText(placeholder)
}

func loadWalletState(m *UI, walletData eth.WalletData) {
	m.walletData = walletData
	m.loadListItems(
		getControlWalletItems(*m),
		m.walletData.PublicKey,
	)
}

func quitToMainMenu(m *UI) {
	m.list.SetItems(getMainItems())
	m.resetListCursor()
	m.setState("main")
	m.setListTitle("✨✨✨")
}
