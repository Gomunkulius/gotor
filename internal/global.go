package internal

import (
	"github.com/charmbracelet/lipgloss"
	"gotor/internal/log"
	"log/slog"
)

var (
	Logger      *slog.Logger
	Port        int      = 6484
	SizePostfix []string = []string{"B", "KiB", "MiB", "GiB", "TiB"}
	HelpString  string   = "ctrl-c/q - quit, b - add torrent, j - remove torrent, e - pause torrent"
)

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("247"))

func InitGlobal() {
	Logger = log.GetLogger()
}
