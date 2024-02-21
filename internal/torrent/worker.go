package torrent

import (
	"gotor/internal/torrent/storage"
)

// DownloadTorrent Must be used as a goroutine
func DownloadTorrent(torrent *Torrent) {
	torrent.Torrent.DownloadAll()
	<-torrent.cancel
}

func RemoveTorrent(s []*Torrent, index int, storage storage.Storage) []*Torrent {
	tor := s[index]
	err := storage.Delete(tor.Torrent.InfoHash().String())
	if err != nil {
		return nil
	}
	tor.cancel <- true
	return append(s[:index], s[index+1:]...)
}

func InitTorrents(torrents []*Torrent) {
	for _, t := range torrents {
		go DownloadTorrent(t)
	}
}
