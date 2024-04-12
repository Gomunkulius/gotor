package torrent

type Storage interface {
	Save(torrent *Torrent) (string, error)
	Get(hash string) (*TorrentModel, error)
	GetAll() ([]*TorrentModel, error)
	Delete(hash string) error
}
