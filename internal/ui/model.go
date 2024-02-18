package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotor/internal"
	"time"
)

type MainModel struct {
	keys   table.KeyMap
	width  int
	height int
	table  TorrentTable
	cnt    int
}

type TickMsg time.Time

func NewModel(table TorrentTable) MainModel {
	return MainModel{
		width:  0,
		height: 0,
		table:  table,
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
		return m, tickEvery()
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
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.cnt++
	m.table.Table, _ = m.table.Table.Update(msg)
	m.table.Update()
	return m, nil

}

func (m MainModel) View() string {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, internal.BaseStyle.Render(m.table.Table.View()), textStyle.Render(internal.HelpString)),
	)
}
