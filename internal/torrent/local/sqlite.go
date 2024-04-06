package local

//
//import (
//	"github.com/anacrolix/torrent"
//	"github.com/anacrolix/torrent/metainfo"
//	"gorm.io/driver/sqlite"
//	"gorm.io/gorm"
//	torrent2 "gotor/internal/torrent"
//	"net/url"
//)
//
//// WARNING: DO NOT WORK! DO NOT USE!
//type storageSqlite struct {
//	db   *gorm.DB
//	conn *torrent.Client
//}
//
//func (s storageSqlite) Save(tf *torrent2.Torrent) (string, error) {
//	var m metainfo.Magnet
//	info := tf.Torrent.Metainfo()
//	m.Trackers = append(m.Trackers, info.UpvertedAnnounceList().DistinctValues()...)
//	if tf.Torrent.Info() != nil {
//		m.DisplayName = tf.Torrent.Info().BestName()
//	}
//	m.InfoHash = tf.Torrent.InfoHash()
//	m.Params = make(url.Values)
//	m.Params["ws"] = tf.Torrent.Metainfo().UrlList
//	hash := tf.Torrent.InfoHash().String()
//	model := torrent2.TorrentModel{
//		TorrentHash: hash,
//		Name:        tf.Torrent.Name(),
//		Magnet:      m.String(),
//	}
//	s.db.Create(&model)
//	if s.db.Error == nil {
//		return hash, s.db.Error
//	}
//	return hash, nil
//}
//
//func (s storageSqlite) Get(hash string) (*torrent2.Torrent, error) {
//	tor := torrent2.TorrentModel{}
//	s.db.Where("TorrentHash =?", hash).First(&tor)
//	if s.db.Error != nil {
//		return nil, s.db.Error
//	}
//	newTorrent, err := torrent2.NewTorrent(tor.Magnet, s.conn, torrent2.UP)
//	if err != nil {
//		return nil, err
//	}
//	<-newTorrent.Torrent.GotInfo()
//	return newTorrent, nil
//}
//
//func (s storageSqlite) GetAll() ([]*torrent2.Torrent, error) {
//	var torrents []*torrent2.TorrentModel
//	s.db.Find(torrents, nil)
//	if s.db.Error == nil {
//		return nil, s.db.Error
//	}
//	res := make([]*torrent2.Torrent, 0)
//	for _, tor := range torrents {
//
//		var magnet string
//		magnet = tor.Magnet
//		tor, err := torrent2.NewTorrent(magnet, s.conn, torrent2.UP)
//		if err != nil {
//			return nil, err
//		}
//		res = append(res, tor)
//	}
//	return res, nil
//}
//
//func (s storageSqlite) Delete(hash string) error {
//	s.db.Where("TorrentHash =?", hash).Delete(&torrent2.TorrentModel{TorrentHash: hash})
//	if s.db.Error != nil {
//		return s.db.Error
//	}
//	return nil
//}
//
//func NewStorage(path string, conn *torrent.Client) (torrent2.Storage, error) {
//
//	db, err := gorm.Open(sqlite.Open("gotor.db"), &gorm.Config{})
//
//	if err != nil {
//		return nil, err
//	}
//	err = db.AutoMigrate(&torrent2.TorrentModel{})
//	if err != nil {
//		return nil, err
//	}
//	return &storageSqlite{db: db, conn: conn}, nil
//}
