package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	testing2 "gotor/internal"
	torrent2 "gotor/internal/torrent"
	"gotor/internal/torrent/local"
	"log"
	"testing"
)

func createTestModel() *MainModel {
	c, err := testing2.MockClient()
	if err != nil || c == nil {
		log.Fatalf("cant create client %v", err)
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
		log.Fatalf("cant create storage %v", err)
		return nil
	}
	torrents, err := storage.GetAll()
	files := torrent2.InitTorrents(torrents, c)
	if err != nil {
		log.Fatalf("cant init torrents %v", err)
		return nil
	}
	torTable := NewTorrentTable(s, files)

	if err != nil {
		log.Fatalf("cant create torrent table %v", err)
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
