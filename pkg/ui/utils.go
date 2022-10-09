package ui

import (
	"eth-toolkit/pkg/eth"
	"fmt"

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
		blurredStyle.Render("Press c to cancel"),
	))
}

func renderOutput(m UI) string {
	return docStyle.Render(fmt.Sprintf(
		"%s\n%s\n%s",
		titleStyle.Render(m.title),
		docStyle.Render(m.output),
		blurredStyle.Render("Press enter to continue"),
	))
}
