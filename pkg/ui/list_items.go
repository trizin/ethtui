package ui

import (
	"ethtui/pkg/hd"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
)

func getMainItems() []list.Item {
	items := []list.Item{
		ListItem{title: "New Wallet", desc: "Create a new wallet", id: "new_wallet"},
		ListItem{title: "New HD Wallet", desc: "Create a new HD wallet", id: "new_hd_wallet"},
		ListItem{title: "Access Wallet", desc: "Access an existing wallet", id: "access_wallet"},
	}
	return items
}

func getProviderItems(m UI) []list.Item {
	items := []list.Item{
		ListItem{title: "Wallet Balance", desc: "Query the balance of active wallet", id: "account_bal"},
		ListItem{title: "Send Transaction", desc: "Send a transaction", id: "send_tx"},
		ListItem{title: "Query Balance", desc: "Query balance of an address", id: "query_bal"},
		ListItem{title: "Query Transaction", desc: "Get transaction details", id: "query_tx"},
		ListItem{title: "Query Block", desc: "Get block details", id: "query_block"},
		ListItem{title: "Go Back", desc: "Go back to wallet management", id: "back"},
	}
	return items
}

func getControlWalletItems(m UI) []list.Item {
	items := []list.Item{}

	if m.provider != nil {
		items = append(items, ListItem{title: "Provider", desc: "Query chain", id: "provider_options"})
	}

	items = append(items, ListItem{title: "Public Key", desc: "Display public key and QR", id: "public_key"},
		ListItem{title: "Private Key", desc: "Display private key and QR", id: "private_key"},
		ListItem{title: "Save Keystore", desc: "Save the wallet to a keystore file", id: "save_keystore"},
		ListItem{title: "Sign Message", desc: "Sign a message with the private key", id: "sign_message"},
		ListItem{title: "Sign Transaction", desc: "Sign a transaction with the private key", id: "sign_transaction"},
		ListItem{title: "Quit", desc: "Quit to main menu", id: "quit"})
	return items
}

func getAccessWalletItems() []list.Item {
	items := []list.Item{
		ListItem{title: "Private Key", desc: "Access your wallet using a private key", id: "pk"},
		ListItem{title: "Keystore File", desc: "Access your wallet using a keystore file", id: "keystore_access"},
		ListItem{title: "Mnemonic", desc: "Access your wallet using mnemonic words", id: "mnemonic"},
		ListItem{title: "Quit", desc: "Quit to main menu", id: "quit"},
	}
	return items
}

func getHdWalletItems(wallet *hd.HDWallet) []list.Item {
	accounts := wallet.GetAddresses(0, 1000)
	items := []list.Item{ListItem{title: "Quit", desc: "Quit to main menu", id: "quit"}}
	for i := 0; i <= len(accounts)-1; i++ {
		acindex := strconv.Itoa(accounts[i].Index)
		items = append(items, ListItem{title: fmt.Sprintf("%s. %s", acindex, accounts[i].Address), id: acindex})
	}
	return items
}
