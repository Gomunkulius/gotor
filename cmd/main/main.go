package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
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
	defer c.Close()
	t, _ := c.AddMagnet("magnet:?xt=urn:btih:HJVLTRV6UJEL7TJJQFT25QZYOAPK3N22&dn=Crusader%20Kings%203%20by%20Igruha&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce")
	<-t.GotInfo()

	go func() {
		t.DownloadAll()
		c.WaitAll()
	}()

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
	torTable := ui.New(s, []*torrent.Torrent{t})

	m := ui.NewModel(torTable)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
