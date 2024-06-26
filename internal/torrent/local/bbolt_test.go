package local

import (
	testing2 "gotor/internal"
	torrent2 "gotor/internal/torrent"
	"testing"
)

func TestStorageBbolt_Get(t *testing.T) {
	c, err := testing2.MockClient()
	if err != nil || c == nil {
		t.Errorf("cannot create client: %v", err)
		return
	}
	gcfg, err := torrent2.NewConfig()
	if err != nil || gcfg == nil {
		t.Fatalf("cannot create config: %v", err)
	}
	storage := NewStorageBbolt("bolt.db", c)
	if storage == nil {
		t.Errorf("cannot create storage: %v", err)
		return
	}
	model := torrent2.TorrentModel{
		TorrentHash: "test",
		Name:        "testt",
		// Debian magnet link for testing
		Magnet: "magnet:?xt=urn:btih:FNTJQAETXQIYA35LKDFTZNAYGW4VUA3C&dn=debian-12.5.0-amd64-netinst.iso&xl=659554304&tr=http%3A%2F%2Fbttracker.debian.org%3A6969%2Fannounce",
	}
	tor, err := torrent2.NewTorrent(model.Magnet, c, torrent2.PAUSE)
	if err != nil {
		t.Errorf("cannot create torrent: %v", err)
	}
	<-tor.Torrent.GotInfo()
	_, err = storage.Save(tor)
	if err != nil {
		t.Errorf("cannot save torrent: %v", err)
	}
	tor1, err := storage.Get(tor.Torrent.InfoHash().String())
	if err != nil {
		t.Errorf("cannot get torrent: %v", err)
	}
	if tor.Torrent.InfoHash().String() != tor1.TorrentHash {
		t.Errorf("got torrent %v, want %v", tor, tor1)
	}
}
