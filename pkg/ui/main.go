package ui

import (
	"eth-toolkit/pkg/eth"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func GetUI() UI {
	m := UI{title: "✨✨✨", list: list.NewModel(getMainItems(), list.NewDefaultDelegate(), 0, 0), input: getText(""), state: "main"}
	m.list.Title = "✨✨✨"
	return m
}

func (m UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		case "alt+c":
			if m.state == "input" || m.state == "sign_transaction" || m.state == "keystore_access" {
				m.setInState("")
				m.setState("main")
			}

		case "c":
			if m.state == "output" {
				msg := m.output
				// copy to clipboard
				clipboard.WriteAll(msg)
			}

		case "alt+e":
			if m.state == "sign_transaction" && m.provider != nil {
				data := m.multiInput[5].Value()
				if data == "" {
					data = "0x"
				}

				gasTipCap, _ := m.provider.GetGasTipCap()
				gasPrice, _ := m.provider.GetGasPrice()
				chainId, _ := m.provider.GetChainId()
				gasLimit, _ := m.provider.GetEstimatedGasUsage([]byte(data))
				nonce, _ := m.provider.GetNonce(m.walletData.PublicKey)

				// convert wei to gwei
				gasPriceFloat := eth.GetGweiValue(gasPrice)
				gasTipCapFloat := eth.GetGweiValue(gasTipCap)

				m.multiInput[7].SetValue(fmt.Sprintf("%f", gasTipCapFloat))
				m.multiInput[6].SetValue(fmt.Sprintf("%d", chainId))
				m.multiInput[4].SetValue(fmt.Sprintf("%f", gasPriceFloat))
				m.multiInput[3].SetValue(fmt.Sprintf("%d", gasLimit))
				m.multiInput[0].SetValue(fmt.Sprintf("%d", nonce))
			}

		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+p":
			if m.state == "main" {
				requestProvider(&m)
			}

		case "tab", "shift+tab", "up", "down":
			if m.state == "sign_transaction" || m.state == "keystore_access" {
				s := msg.String()
				m, cmds := moveIndex(m, s)
				return m, tea.Batch(cmds...)
			}

		case "enter":
			m, cmd = handleEnterPress(m)
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		docStyle.Width(msg.Width)
		m.input.Width = int(float64(msg.Width*5) / 6)
	}

	if m.state == "main" || m.state == "hdwallet" {
		m.list, cmd = m.list.Update(msg)
	}

	if m.state == "input" {
		m.input, cmd = m.input.Update(msg)
	}

	if m.state == "sign_transaction" || m.state == "keystore_access" {
		cmd = m.updateInputs(msg)
	}

	return m, cmd
}

func (m UI) View() string {

	if m.choice.title != "" {
		switch m.state {

		case "sign_transaction":
			b := renderMultiInput(m)

			return docStyle.Render(
				fmt.Sprintf(
					"%s\n\n%s\n%s\n%s",
					titleStyle.Render("Sign Transaction"),
					b,
					blurredStyle.Render("Press alt+c to cancel"),
					blurredStyle.Render("Press alt+e to estimate values (if connected to an RPC)"),
				))

		case "keystore_access":
			b := renderMultiInput(m)
			return docStyle.Render(
				fmt.Sprintf(
					"%s\n\n%s\n%s",
					titleStyle.Render("Access Keystore"),
					b,
					blurredStyle.Render("Press alt+c to cancel"),
				))

		case "input":
			return renderInput(m)

		case "output":
			return renderOutput(m)
		}
	}

	return docStyle.Render(m.list.View())
}
