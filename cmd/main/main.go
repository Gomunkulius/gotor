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
		println("cant log into file")
		return
	}
	defer f.Close()
	cfg := torrent.NewDefaultClientConfig() // TODO: config
	c, err := torrent.NewClient(cfg)
	if err != nil || c == nil {
		println("cant connect")
		return
	}
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
	storage := local.NewStorageBbolt("bolt.db", c)
	if storage == nil {
		println("cant create storage")
		return
	}
	torrents, err := storage.GetAll()
	torrent2.InitTorrents(torrents)
	if err != nil {
		println("cant init torrents")
		return
	}
	torTable := ui.New(s, torrents)
	if err != nil {
		println("cant init ui")
		return
	}
	m := ui.NewModel(torTable, c, storage)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
