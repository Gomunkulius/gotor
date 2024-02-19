package ui

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"time"
)

type MainModel struct {
	keys       table.KeyMap
	width      int
	height     int
	table      TorrentTable
	inputFlag  bool
	inputField textinput.Model
	cancels    []chan bool
	conn       *torrent.Client
}

type TickMsg time.Time

func NewModel(table TorrentTable, conn *torrent.Client) MainModel {
	return MainModel{
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
			if m.table.Table.Focused() {
				m.table.Table.Blur()
			} else {
				m.table.Table.Focus()
			}
		case "b":
			if !m.inputFlag {
				m.inputFlag = true
				m.inputField.Focus()
			}

		case "enter":
			if m.inputFlag {
				magnet := m.inputField.Value()
				t, _ := m.conn.AddMagnet(magnet)
				m.table.Torrents = append(m.table.Torrents, t)
				<-t.GotInfo()
				canc := make(chan bool)
				go torrent2.DownloadTorrent(t, canc)
				m.cancels = append(m.cancels, canc)
				m.inputFlag = false
				return m, nil
			}
		case "j":
			if !m.inputFlag {
				if len(m.table.Torrents) == 0 {
					return m, nil
				}
				m.table.Torrents = torrent2.RemoveTorrent(m.table.Torrents, m.cancels[m.table.Table.Cursor()], m.table.Table.Cursor())
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
		lipgloss.JoinVertical(lipgloss.Center, internal.BaseStyle.Render(m.table.Table.View()), textStyle.Render(internal.HelpString)),
	)
}
