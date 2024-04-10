package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"gotor/internal/torrent"
)

func ExitCmd(state ProgramState) tea.Cmd {
	return func() tea.Msg {
		return ChangeStateMsg(state)
	}
}

func SaveExitCmd(storage torrent.Storage, tor *torrent.Torrent) tea.Cmd {
	return func() tea.Msg {
		storage.Save(tor)
		return ChangeStateMsg(Main)
	}
}
