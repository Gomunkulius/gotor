package torrent

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"net/url"
)

type Status int

const (
	UP Status = iota
	PAUSE
)

type Torrent struct {
	Torrent *torrent.Torrent

	// cancel is the channel used to stop the torrent download
	cancel  chan bool
	Status  Status
	Speed1s int64
}

type TorrentModel struct {
	TorrentHash string `gorm:"primary_key"`
	Status      int
	Name        string
	Magnet      string
}

func NewTorrent(magnet string, conn *torrent.Client, status Status) (*Torrent, error) {
	tor, err := conn.AddMagnet(magnet)

	if err != nil {
		return nil, err
	}

	return &Torrent{
		Torrent: tor,
		cancel:  make(chan bool),
		Status:  status,
	}, nil
}

func NewTorrentFromFile(path string, conn *torrent.Client, status Status) (*Torrent, error) {
	tmptor, err := conn.AddTorrentFromFile(path)
	if err != nil {
		return nil, err
	}
	var m metainfo.Magnet
	info := tmptor.Metainfo()
	m.Trackers = append(m.Trackers, info.UpvertedAnnounceList().DistinctValues()...)
	if tmptor.Info() != nil {
		m.DisplayName = tmptor.Info().BestName()
	}
	m.InfoHash = tmptor.InfoHash()
	m.Params = make(url.Values)
	m.Params["ws"] = tmptor.Metainfo().UrlList
	tor, err := conn.AddMagnet(m.String())
	if err != nil {
		return nil, err
	}
	return &Torrent{
		Torrent: tor,
		cancel:  make(chan bool),
		Status:  status,
	}, nil
}

func (t *TorrentModel) ToTorrent(conn *torrent.Client) (*Torrent, error) {
	newTorrent, err := NewTorrent(t.Magnet, conn, Status(t.Status))
	if err != nil {
		return nil, err
	}
	<-newTorrent.Torrent.GotInfo()
	return newTorrent, nil
}

func (tf *Torrent) ToDTO() *TorrentModel {
	var m metainfo.Magnet
	info := tf.Torrent.Metainfo()
	m.Trackers = append(m.Trackers, info.UpvertedAnnounceList().DistinctValues()...)
	if tf.Torrent.Info() != nil {
		m.DisplayName = tf.Torrent.Info().BestName()
	}
	m.InfoHash = tf.Torrent.InfoHash()
	m.Params = make(url.Values)
	m.Params["ws"] = tf.Torrent.Metainfo().UrlList
	model := TorrentModel{
		TorrentHash: tf.Torrent.InfoHash().String(),
		Name:        tf.Torrent.Name(),
		Magnet:      m.String(),
	}
	return &model
}
