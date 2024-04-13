package internal

import (
	"github.com/charmbracelet/lipgloss"
	"log/slog"
)

var (
	Logger     *slog.Logger
	HelpString string = "ctrl-c/q - quit, b - add torrent, j - remove torrent, e - pause torrent"
)

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("247"))

var InputStyle = lipgloss.NewStyle().
	Width(80).Align(lipgloss.Center).Foreground(lipgloss.Color("229"))

func InitGlobal() {
	Logger = GetLogger()

}
