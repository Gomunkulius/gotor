package main

import (
	"fmt"
	log2 "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"gotor/internal/torrent/local"
	"gotor/internal/ui"
	"math/rand/v2"
	"os"
)

func main() {
	internal.InitGlobal()

	// anacrolix dumbass!
	if len(log2.Default.Handlers) != 0 {
		log2.Default.Handlers[0] = internal.MyHandler{}
	}

	f, err := tea.LogToFile("LOG.log", "debug")
	if err != nil {
		println("cant log into file")
		return
	}
	defer f.Close()
	gcfg, err := torrent2.NewConfig()
	if err != nil || gcfg == nil {
		fmt.Printf("cant create config err %v\n", err)
		return
	}
	cfg := torrent.NewDefaultClientConfig() // TODO: config
	cfg.DataDir = gcfg.DataDir
	cfg.ListenPort = gcfg.Port
	cfg.Logger = log2.Logger{}
	cfg.Debug = false
	c, err := torrent.NewClient(cfg)
	defer c.Close()
	if err != nil || c == nil {
		fmt.Printf("cant connect err %v\n", err)
		cfg.ListenPort = rand.IntN(65535-20000) + 20000
		crand, err := torrent.NewClient(cfg)
		if err != nil || crand == nil {
			fmt.Printf("cant connect on port %d err %v\n", cfg.ListenPort, err)
			return
		}
		c = crand
		fmt.Printf("new non default client created, listening on port %d\n", cfg.ListenPort)
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
	files := torrent2.InitTorrents(torrents, c)
	if err != nil {
		println("cant init torrents")
		return
	}
	torTable := ui.NewTorrentTable(s, files)
	go torTable.CountSpeed()
	m := ui.NewModel(torTable, c, storage)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
