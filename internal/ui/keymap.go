package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	AddTorrent    key.Binding
	DeleteTorrent key.Binding
	PauseTorrent  key.Binding
	Quit          key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.AddTorrent, k.DeleteTorrent, k.PauseTorrent}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit, k.AddTorrent, k.DeleteTorrent, k.PauseTorrent}}
}

var keys = KeyMap{
	AddTorrent: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "add torrent")),
	DeleteTorrent: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "delete torrent")),
	PauseTorrent: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("j", "pause torrent")),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit")),
}
