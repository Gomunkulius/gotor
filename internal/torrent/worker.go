package torrent

import "github.com/anacrolix/torrent"

// DownloadTorrent Must be used as a goroutine
func DownloadTorrent(torrent *Torrent) {
	torrent.Torrent.DownloadAll()
	<-torrent.cancel
}

func RemoveTorrent(s []*Torrent, index int, storage Storage) []*Torrent {
	tor := s[index]
	err := storage.Delete(tor.Torrent.InfoHash().String())
	if err != nil {
		return nil
	}
	tor.cancel <- true
	return append(s[:index], s[index+1:]...)
}

func TogglePauseTorrent(s []*Torrent, index int) []*Torrent {
	tor := s[index]
	switch tor.Status {
	case UP:
		tor.Torrent.DisallowDataDownload()
		tor.Status = PAUSE
	case PAUSE:
		tor.Torrent.AllowDataDownload()
		tor.Status = UP
	}
	return s
}

func InitTorrents(torrents []*TorrentModel, conn *torrent.Client) []*Torrent {
	var res []*Torrent
	for _, tm := range torrents {
		t, _ := tm.ToTorrent(conn)
		res = append(res, t)
		if t.Status == PAUSE {
			t.Torrent.DisallowDataDownload()
		}
		go DownloadTorrent(t)
	}
	return res
}
