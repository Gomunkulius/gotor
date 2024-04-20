package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"os/exec"
	"runtime"
)

type ProgramModel struct {
	height  int
	width   int
	storage torrent2.Storage
	table   *TorrentTable
	keys    KeyMap
	help    help.Model
	cfg     *torrent2.Config
}

func (m ProgramModel) Init() tea.Cmd {
	return nil
}

func (m ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.table.Table, _ = m.table.Table.Update(msg)
		m.table.Update(m.width, m.height)
		m.table.Table.UpdateViewport()
		return m, tea.Batch(tickEvery())
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
			_, err := m.storage.Save(m.table.Torrents[m.table.Table.Cursor()])
			if err != nil {
				return m, tea.Batch(tickEvery())
			}
			return m, tickEvery()
		case "j":
			if len(m.table.Torrents) == 0 {
				return m, nil
			}
			m.table.Torrents = torrent2.RemoveTorrent(
				m.table.Torrents,
				m.table.Table.Cursor(),
				m.storage)
		case "b":
			return m, ExitCmd(Choose)
		case "enter":
			if runtime.GOOS == "windows" {
				cmd := exec.Command("cmd", "/c", fmt.Sprintf("explorer.exe /select %s",
					m.cfg.DataDir))
				cmd.Run()
			} else {
				cmd := exec.Command("open", m.cfg.DataDir)
				cmd.Run()
			}
		}
	}
	m.table.Table, _ = m.table.Table.Update(msg)
	m.table.Update(m.width, m.height)
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
func NewProgramModel(height int, width int, storage torrent2.Storage, table *TorrentTable, keys KeyMap, cfg *torrent2.Config) *ProgramModel {
	return &ProgramModel{
		cfg:     cfg,
		height:  height,
		width:   width,
		storage: storage,
		table:   table,
		keys:    keys,
		help:    help.New(),
	}
}
