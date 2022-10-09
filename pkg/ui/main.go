package ui

import (
	"fmt"

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
					"%s\n\n%s\n%s",
					titleStyle.Render("Sign Transaction"),
					b,
					blurredStyle.Render("Press c to cancel"),
				))

		case "keystore_access":
			b := renderMultiInput(m)
			return docStyle.Render(
				fmt.Sprintf(
					"%s\n\n%s\n%s",
					titleStyle.Render("Access Keystore"),
					b,
					blurredStyle.Render("Press c to cancel"),
				))

		case "input":
			return renderInput(m)

		case "output":
			return renderOutput(m)
		}
	}

	return docStyle.Render(m.list.View())
}
