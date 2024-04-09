package ui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
)

type ProgramModel struct {
	height  int
	width   int
	storage torrent2.Storage
	table   TorrentTable
	keys    KeyMap
	help    help.Model
}

func (m ProgramModel) Init() tea.Cmd {
	return nil
}

func (m ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k":
			if len(m.table.Torrents) == 0 {
				return m, nil
			}
			m.table.Torrents = torrent2.TogglePauseTorrent(
				m.table.Torrents,
				m.table.Table.Cursor())
		case "j":
			if len(m.table.Torrents) == 0 {
				return m, nil
			}
			m.table.Torrents = torrent2.RemoveTorrent(
				m.table.Torrents,
				m.table.Table.Cursor(),
				m.storage)
		case "b":
			return m, ExitCmd(Input)
		}
	}
	m.table.Table, _ = m.table.Table.Update(msg)
	m.table.Update()
	return m, nil
}

func (m ProgramModel) View() string {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, internal.BaseStyle.Render(m.table.Table.View()), textStyle.Render(m.help.View(m.keys))))
}
func NewProgramModel(height int, width int, storage torrent2.Storage, table TorrentTable, keys KeyMap) *ProgramModel {
	return &ProgramModel{
		height:  height,
		width:   width,
		storage: storage,
		table:   table,
		keys:    keys,
		help:    help.New(),
	}
}
