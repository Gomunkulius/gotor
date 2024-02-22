package torrent

type Storage interface {
	Save(torrent *Torrent) (string, error)
	Get(hash string) (*Torrent, error)
	GetAll() ([]*Torrent, error)
	Delete(hash string) error
}
