package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	AddTorrent    key.Binding
	DeleteTorrent key.Binding
	PauseTorrent  key.Binding
	OpenTorrent   key.Binding
	Quit          key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.AddTorrent, k.DeleteTorrent, k.PauseTorrent, k.OpenTorrent}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit, k.AddTorrent, k.DeleteTorrent, k.PauseTorrent}}
}

var keys = KeyMap{
	AddTorrent: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "add torrent")),
	PauseTorrent: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "pause torrent")),
	DeleteTorrent: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "delete torrent")),
	OpenTorrent: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open torrent")),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit")),
}
