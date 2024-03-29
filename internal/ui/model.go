package ui

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"time"
)

type MainModel struct {
	keys       KeyMap
	help       help.Model
	width      int
	height     int
	table      TorrentTable
	inputFlag  bool
	inputField textinput.Model
	conn       *torrent.Client
	storage    torrent2.Storage
}

type TickMsg time.Time

func NewModel(table TorrentTable, conn *torrent.Client, storage torrent2.Storage) MainModel {
	return MainModel{
		keys: keys,

		help:       help.New(),
		storage:    storage,
		inputField: textinput.New(),
		table:      table,
		conn:       conn,
	}
}

// Send a message every second.
func tickEvery() tea.Cmd {
	return tea.Every(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TorrentInfoUpdate string

func (m MainModel) Init() tea.Cmd {
	m.help.ShowAll = true
	return tickEvery()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd

	switch msg := msg.(type) {
	case TickMsg:
		m.table.Table, _ = m.table.Table.Update(msg)
		m.table.Update()
		m.table.Table.UpdateViewport()
		return m, tea.Batch(tickEvery())
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.inputFlag {
				m.inputFlag = false
			}
		case "b":
			if !m.inputFlag {
				m.inputFlag = true
				m.inputField.Focus()
				return m, nil
			}

		case "enter":
			if m.inputFlag {
				magnet := m.inputField.Value()
				t, err := torrent2.NewTorrent(magnet, m.conn, torrent2.UP)
				<-t.Torrent.GotInfo()
				if err != nil {
					return m, nil
				}
				for _, t2 := range m.table.Torrents {
					if t == t2 {
						m.inputFlag = false
						m.inputField.SetValue("")
						return m, nil
					}
				}
				m.table.Torrents = append(m.table.Torrents, t)
				go torrent2.DownloadTorrent(t)
				m.inputFlag = false
				_, err = m.storage.Save(t)
				if err != nil {
					return nil, nil
				}
				return m, nil
			}
		case "j":
			if !m.inputFlag {
				if len(m.table.Torrents) == 0 {
					return m, nil
				}
				m.table.Torrents = torrent2.RemoveTorrent(
					m.table.Torrents,
					m.table.Table.Cursor(),
					m.storage)

			}

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.inputField, _ = m.inputField.Update(msg)

	m.table.Table, _ = m.table.Table.Update(msg)
	m.table.Update()
	return m, nil

}

func (m MainModel) View() string {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	if m.inputFlag {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, internal.InputStyle.Render("Enter magnet"),
				internal.InputStyle.Render(m.inputField.View())),
		)
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, internal.BaseStyle.Render(m.table.Table.View()), textStyle.Render(m.help.View(m.keys))))
}
