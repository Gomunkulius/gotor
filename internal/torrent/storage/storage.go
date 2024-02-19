package storage

import "github.com/anacrolix/torrent"

type Storage interface {
	Save(torrent *torrent.Torrent) (string, error)
	Get(hash string) (*torrent.Torrent, error)
	Update(hash string, update torrent.Torrent) error
	Delete(hash string) error
}