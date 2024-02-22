package torrent

import "github.com/anacrolix/torrent"

type Status int

const (
	UP Status = iota
	PAUSE
)

type Torrent struct {
	Torrent *torrent.Torrent
	cancel  chan bool
	Status  Status
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
