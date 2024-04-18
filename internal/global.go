package internal

import (
	log2 "github.com/anacrolix/log"
	"github.com/charmbracelet/lipgloss"
)

type MyHandler struct{}

func (m MyHandler) Handle(r log2.Record) {
	return
}

const VERSION = "0.3.2"

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("247"))

var InputStyle = lipgloss.NewStyle().
	Width(80).Align(lipgloss.Center).Foreground(lipgloss.Color("229"))

var SelChooseStyle = InputStyle.Copy().Foreground(lipgloss.Color("226"))

func InitGlobal() {}
