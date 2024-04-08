package ui

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	"gotor/internal/ui"
)

type InputModel struct {
	width      int
	height     int
	table      ui.TorrentTable
	inputField textinput.Model
	conn       *torrent.Client
}

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.TickMsg:
		m.table.Table, _ = m.table.Table.Update(msg)
		m.table.Update()
		m.table.Table.UpdateViewport()
		return m, tea.Batch(tickEvery())
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {

		}
	}

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

func NewInputModel(width int, height int, table ui.TorrentTable, conn *torrent.Client) *InputModel {
	return &InputModel{
		width:      width,
		height:     height,
		table:      table,
		inputField: textinput.New(),
		conn:       conn,
	}
}
