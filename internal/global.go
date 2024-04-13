package internal

import (
	"github.com/charmbracelet/lipgloss"
)

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("247"))

var InputStyle = lipgloss.NewStyle().
	Width(80).Align(lipgloss.Center).Foreground(lipgloss.Color("229"))

var SelChooseStyle = InputStyle.Copy().Foreground(lipgloss.Color("226"))

func InitGlobal() {}
