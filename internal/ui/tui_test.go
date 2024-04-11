package ui

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	torrent2 "gotor/internal/torrent"
	"gotor/internal/torrent/local"
	"testing"
)

func createTestModel() *MainModel {
	cfg := torrent.NewDefaultClientConfig()
	c, err := torrent.NewClient(cfg)
	if err != nil || c == nil {
		println("cant connect")
		return nil
	}
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("247")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	storage := local.NewStorageBbolt("bolt.db", c)
	if storage == nil {
		println("cant create storage")
		return nil
	}
	torrents, err := storage.GetAll()
	torrent2.InitTorrents(torrents)
	if err != nil {
		println("cant init torrents")
		return nil
	}
	torTable := NewTorrentTable(s, torrents)
	if err != nil {
		println("cant init ui")
		return nil
	}
	m := NewModel(torTable, c, storage)
	return &m
}

func TestMainModel_Update(t *testing.T) {
	m := createTestModel()
	mod, _ := m.Update(ExitCmd(Input)())
	main := mod.(MainModel)
	m = &main
	if m.state != Input {
		t.Errorf("expected diff state")
	}
	if m.View() != m.inputModel.View() {
		t.Errorf("expected diff view")
	}
	mod, _ = m.Update(ExitCmd(Main)())
	main = mod.(MainModel)
	m = &main
	if m.state != Main {
		t.Errorf("expected diff state")
	}
	if m.View() != m.programModel.View() {
		t.Errorf("expected diff view")
	}
}
