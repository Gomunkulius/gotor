package sqlite

import (
	"database/sql"
	"errors"
	"github.com/anacrolix/torrent"
	"os"
)

type StorageSqlite struct {
	db   *sql.DB
	conn *torrent.Client
}

func (s StorageSqlite) Save(tf *torrent.Torrent) (string, error) {
	q := "INSERT INTO torrent (torrent_hash, name, magnet) VALUES (?, ?)"
	hash := tf.InfoHash().String()
	_, err := s.db.Exec(q, hash, tf.Name())
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s StorageSqlite) Get(hash string) (*torrent.Torrent, error) {
	q := "SELECT magnet FROM torrents WHERE torrent_hash = ?"
	var magnet string
	err := s.db.QueryRow(q, hash).Scan(&hash)
	if err != nil {
		return nil, err
	}
	tor, err := s.conn.AddMagnet(magnet)
	if err != nil {
		return nil, err
	}
	<-tor.GotInfo()
	return tor, nil
}

func (s StorageSqlite) Update(hash string, update torrent.Torrent) error {

}

func (s StorageSqlite) Delete(hash string) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(path string) (*StorageSqlite, error) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// Creating a data.db
		_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	q := "create table torrents(\n    torrent_hash text not null unique,\n    name text not null,\n    magnet text not null,\n    constraint pk_torrent_hash primary key (torrent_hash)\n)"
	_, err = db.Exec(q)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &StorageSqlite{db: db}, nil
}
