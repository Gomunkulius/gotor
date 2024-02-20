package torrent

import (
	"github.com/anacrolix/torrent"
	"gotor/internal/torrent/storage"
)

// DownloadTorrent Must be used as a goroutine
func DownloadTorrent(torrent *torrent.Torrent, cancel chan bool) {
	torrent.DownloadAll()
	<-cancel
}

func RemoveTorrent(s []*torrent.Torrent, cancel chan bool, index int, storage storage.Storage) []*torrent.Torrent {
	err := storage.Delete(s[index].InfoHash().String())
	if err != nil {
		return nil
	}
	cancel <- true
	return append(s[:index], s[index+1:]...)
}

func InitTorrents(torrents []*torrent.Torrent) (res []chan bool) {
	for _, t := range torrents {
		cancel := make(chan bool)
		go DownloadTorrent(t, cancel)
		res = append(res, cancel)
	}
	return res
}
