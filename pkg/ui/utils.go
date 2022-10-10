package ui

import (
	"eth-toolkit/pkg/eth"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

func displayWalletPublicKey(walletData eth.WalletData) string {
	return fmt.Sprintf(
		"%s\n%s",
		walletData.PublicKeyQR.ToSmallString(false),
		"Public Key: "+walletData.PublicKey,
	)
}

func displayWalletPrivateKey(walletData eth.WalletData) string {
	return fmt.Sprintf(
		"%s\n%s",
		walletData.PrivateKeyQR.ToSmallString(false),
		"Private Key: "+walletData.PrivateKey,
	)
}

func getText(placeHolder string) textinput.Model {
	ti := textinput.NewModel()
	ti.Placeholder = placeHolder
	ti.Focus()
	return ti
}

func renderInput(m UI) string {
	return docStyle.Render(fmt.Sprintf(
		"%s\n%s\n%s",
		titleStyle.Render(m.title),
		m.input.View(),
		blurredStyle.Render("Press alt+c to cancel"),
	))
}

func renderOutput(m UI) string {
	return docStyle.Render(fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		titleStyle.Render(m.title),
		docStyle.Render(m.output),
		blurredStyle.Render("Press enter to continue"),
		blurredStyle.Render("Press c to copy to clipboard"),
	))
}

func renderMultiInput(m UI) string {
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
}

func handleError(m *UI, err error) bool {
	if err != nil {
		setOutputState(m, "Error", err.Error())
		return true
	}
	return false
}
