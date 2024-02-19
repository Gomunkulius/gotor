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
	torTable := ui.New(s, []*torrent.Torrent{})

	m := ui.NewModel(torTable, c)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
