package local

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/chaisql/chai"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	torrent2 "gotor/internal/torrent"
	"net/url"
	"os"
)

type storageSqlite struct {
	db   *gorm.DB
	conn *torrent.Client
}

func (s storageSqlite) Save(tf *torrent2.Torrent) (string, error) {
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

func (s storageSqlite) Get(hash string) (*torrent2.Torrent, error) {
	q := "SELECT magnet FROM torrents WHERE torrent_hash = ?;"
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

func (s storageSqlite) GetAll() ([]*torrent2.Torrent, error) {
	q := "SELECT magnet FROM torrents;"
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

func (s storageSqlite) Delete(hash string) error {
	s.db.Where("TorrentHash =?", hash).Delete(&torrent2.TorrentModel{TorrentHash: hash})
	if s.db.Error != nil {
		return s.db.Error
	}
	return nil
}

func NewStorage(path string, conn *torrent.Client) (torrent2.Storage, error) {

	if _, err := os.Stat(path); err != nil {
		// Creating a data.db
		_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open("gotor.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&torrent2.TorrentModel{})
	if err != nil {
		return nil, err
	}
	return &storageSqlite{db: db, conn: conn}, nil
}
