package local

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/chaisql/chai"
	torrent2 "gotor/internal/torrent"
	"net/url"
	"os"
)

type StorageSqlite struct {
	db   *chai.DB
	conn *torrent.Client
}

func (s StorageSqlite) Save(tf *torrent2.Torrent) (string, error) {
	q := "INSERT INTO torrents (torrent_hash, name, magnet) VALUES (?, ?, ?)"
	var m metainfo.Magnet
	info := tf.Torrent.Metainfo()
	m.Trackers = append(m.Trackers, info.UpvertedAnnounceList().DistinctValues()...)
	if tf.Torrent.Info() != nil {
		m.DisplayName = tf.Torrent.Info().BestName()
	}
	m.InfoHash = tf.Torrent.InfoHash()
	m.Params = make(url.Values)
	m.Params["ws"] = tf.Torrent.Metainfo().UrlList
	hash := tf.Torrent.InfoHash().String()
	err := s.db.Exec(q, hash, tf.Torrent.Name(), m.String())
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s StorageSqlite) Get(hash string) (*torrent2.Torrent, error) {
	q := "SELECT magnet FROM torrents WHERE torrent_hash = ?"
	var magnet string
	row, err := s.db.QueryRow(q, hash)
	err = row.Scan(&magnet)
	if err != nil {
		return nil, err
	}
	newTorrent, err := torrent2.NewTorrent(magnet, s.conn, torrent2.UP)
	if err != nil {
		return nil, err
	}
	<-newTorrent.Torrent.GotInfo()
	return newTorrent, nil
}

func (s StorageSqlite) GetAll() ([]*torrent2.Torrent, error) {
	q := "SELECT magnet FROM torrents"
	row, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	res := make([]*torrent2.Torrent, 0)
	defer row.Close()

	err = row.Iterate(func(r *chai.Row) error {
		var magnet string
		err := r.Scan(&magnet)
		if err != nil {
			return err
		}
		tor, err := torrent2.NewTorrent(magnet, s.conn, torrent2.UP)
		if err != nil {
			return err
		}
		res = append(res, tor)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s StorageSqlite) Delete(hash string) error {
	q := "delete from torrents where torrent_hash = %s"
	err := s.db.Exec(fmt.Sprintf(q, hash))
	if err != nil {
		return err
	}
	return nil
}

func NewStorage(path string, conn *torrent.Client) (*StorageSqlite, error) {

	if _, err := os.Stat(path); err != nil {
		// Creating a data.db
		_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
	}
	db, err := chai.Open(path)

	if err != nil {
		return nil, err
	}
	q := "create table torrents(\n    torrent_hash text not null unique,\n    name text not null,\n    magnet text not null,\n    constraint pk_torrent_hash primary key (torrent_hash)\n)"
	err = db.Exec(q)
	if err != nil {
		return nil, err
	}
	return &StorageSqlite{db: db, conn: conn}, nil
}
