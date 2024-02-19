package torrent

import "github.com/anacrolix/torrent"

// DownloadTorrent Must be used as a goroutine
func DownloadTorrent(torrent *torrent.Torrent, cancel chan bool) {
	torrent.DownloadAll()
	<-cancel
}

func RemoveTorrent(s []*torrent.Torrent, cancel chan bool, index int) []*torrent.Torrent {
	cancel <- true
	return append(s[:index], s[index+1:]...)
}
