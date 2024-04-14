package torrent

import (
	"gotor/internal"
	"testing"
)

func TestTorrent_ToDTO(t *testing.T) {
	c, err := internal.MockClient()
	defer c.Close()
	if err != nil {
		t.Fatalf("error creating mock client: %s", err.Error())
	}
	// debian magnet for test
	test_magnet := "magnet:?xt=urn:btih:FNTJQAETXQIYA35LKDFTZNAYGW4VUA3C&dn=debian-12.5.0-amd64-netinst.iso&xl=659554304&tr=http%3A%2F%2Fbttracker.debian.org%3A6969%2Fannounce"
	newTorrent, err := NewTorrent(test_magnet, c, UP)
	if err != nil {
		t.Fatalf("failed to create torrent: %v", err)
	}
	<-newTorrent.Torrent.GotInfo()
	torrentDTO := newTorrent.ToDTO()
	if torrentDTO.TorrentHash != newTorrent.Torrent.InfoHash().String() {
		t.Fatalf("expected torrent hash to be %s but got %s", newTorrent.Torrent.InfoHash().String(), torrentDTO.TorrentHash)
	}
	if torrentDTO.Name != newTorrent.Torrent.Name() {
		t.Fatalf("expected torrent name to be %s but got %s", newTorrent.Torrent.Name(), torrentDTO.Name)
	}
	return
}

func TestTorrentModel_ToTorrent(t *testing.T) {
	c, err := internal.MockClient()
	defer c.Close()
	if err != nil {
		t.Fatalf("error creating mock client: %s", err.Error())
	}
	test_magnet := "magnet:?xt=urn:btih:FNTJQAETXQIYA35LKDFTZNAYGW4VUA3C&dn=debian-12.5.0-amd64-netinst.iso&xl=659554304&tr=http%3A%2F%2Fbttracker.debian.org%3A6969%2Fannounce"
	oldTorrent, err := NewTorrent(test_magnet, c, UP)
	if err != nil {
		t.Fatalf("failed to create torrent: %v", err)
	}
	dto := TorrentModel{
		TorrentHash: oldTorrent.Torrent.InfoHash().String(),
		Status:      0,
		Name:        oldTorrent.Torrent.Name(),
		Magnet:      test_magnet,
	}
	newTorrent, err := dto.ToTorrent(c)
	if err != nil {
		t.Fatalf("failed to convert torrent to Torrent: %v", err)
	}
	if newTorrent.Torrent.InfoHash() != oldTorrent.Torrent.InfoHash() {
		t.Fatalf("expected torrent hash to be %s but got %s", oldTorrent.Torrent.InfoHash().String(), newTorrent.Torrent.InfoHash().String())
	}
	return
}
