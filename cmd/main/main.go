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
	t, _ := c.AddMagnet("magnet:?xt=urn:btih:XF26SWOW4FRBWVDKYJAUQQCWJM3U2APZ&dn=debian-12.5.0-arm64-netinst.iso&xl=551858176&tr=http%3A%2F%2Fbttracker.debian.org%3A6969%2Fannounce")
	t1, _ := c.AddMagnet("magnet:?xt=urn:btih:FNTJQAETXQIYA35LKDFTZNAYGW4VUA3C&dn=debian-12.5.0-amd64-netinst.iso&xl=659554304&tr=http%3A%2F%2Fbttracker.debian.org%3A6969%2Fannounce")
	<-t.GotInfo()
	<-t1.GotInfo()
	go func() {
		t.DownloadAll()
		c.WaitAll()
	}()
	go func() {
		t1.DownloadAll()
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
	torTable := ui.New(s, []*torrent.Torrent{t, t1})

	m := ui.NewModel(torTable)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
