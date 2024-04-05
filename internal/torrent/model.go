package torrent

import "github.com/anacrolix/torrent"

type Status int

const (
	UP Status = iota
	PAUSE
)

type Torrent struct {
	Torrent *torrent.Torrent

	// cancel is the channel used to stop the torrent download
	cancel chan bool
	Status Status
}

type TorrentModel struct {
	TorrentHash string `gorm:"primary_key"`
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
