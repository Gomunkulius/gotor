package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"gotor/internal/torrent/local"
	"gotor/internal/ui"
	"os"
)

func main() {
	internal.InitGlobal()
	f, err := tea.LogToFile("LOG.log", "debug")
	if err != nil {
		return
	}
	defer f.Close()
	c, _ := torrent.NewClient(nil)
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("247")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	storage, err := local.NewStorage(".", c)
	torrents, err := storage.GetAll()
	cancels := torrent2.InitTorrents(torrents)
	if err != nil {
		return
	}
	torTable := ui.New(s, torrents)
	if err != nil {
		return
	}
	m := ui.NewModel(torTable, c, storage, cancels)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
