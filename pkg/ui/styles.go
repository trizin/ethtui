package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	titleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#9763e6")).Bold(true)
	errorTitleStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#DC143C")).Foreground(lipgloss.Color("#ffffff"))
	focusedButton       = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
