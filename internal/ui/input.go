package ui

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
)

type InputModel struct {
	width      int
	height     int
	table      *TorrentTable
	storage    torrent2.Storage
	inputField textinput.Model
	conn       *torrent.Client
}

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.inputField.Focus()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "esc":
				return m, ExitCmd(Main)
			case "enter":
				magnet := m.inputField.Value()
				t, err := torrent2.NewTorrent(magnet, m.conn, torrent2.UP)
				if err != nil {
					return m, nil
				}
				<-t.Torrent.GotInfo()
				if err != nil {
					return m, nil
				}
				for _, t2 := range m.table.Torrents {
					if t.Torrent.InfoHash() == t2.Torrent.InfoHash() {
						return m, ExitCmd(Main)
					}
				}
				m.table.Torrents = append(m.table.Torrents, t)
				go torrent2.DownloadTorrent(t)
				return m, tea.Batch(SaveExitCmd(m.storage, t), tickEvery())
			}
		}
	}
	var cmd tea.Cmd
	m.inputField, cmd = m.inputField.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, internal.InputStyle.Render("Enter magnet"),
			internal.InputStyle.Render(m.inputField.View())),
	)
}

func NewInputModel(width int, height int, table *TorrentTable, conn *torrent.Client, storage torrent2.Storage) *InputModel {
	return &InputModel{
		storage:    storage,
		width:      width,
		height:     height,
		table:      table,
		inputField: textinput.New(),
		conn:       conn,
	}
}
