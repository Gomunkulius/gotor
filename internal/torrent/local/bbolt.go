package local

import (
	"encoding/json"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	torrent2 "gotor/internal/torrent"
	"gotor/pkg/bbolt/client"
	"net/url"
)

type storageBbolt struct {
	client client.Client
	conn   *torrent.Client
}

func (s storageBbolt) Save(tf *torrent2.Torrent) (string, error) {
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
	model := torrent2.TorrentModel{
		TorrentHash: hash,
		Name:        tf.Torrent.Name(),
		Magnet:      m.String(),
	}
	tfjson, err := json.Marshal(model)
	if err != nil {
		return hash, err
	}
	err = s.client.Put([]byte(model.TorrentHash), tfjson)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s storageBbolt) Get(hash string) (*torrent2.Torrent, error) {
	tor := torrent2.TorrentModel{}
	buf := make([]byte, 0)
	buf = s.client.Get([]byte(hash))
	err := json.Unmarshal(buf, &tor)
	if err != nil {
		return nil, err
	}
	newTorrent, err := torrent2.NewTorrent(tor.Magnet, s.conn, torrent2.UP)
	if err != nil {
		return nil, err
	}
	<-newTorrent.Torrent.GotInfo()
	return newTorrent, nil
}

func (s storageBbolt) GetAll() ([]*torrent2.Torrent, error) {
	var res []*torrent2.Torrent
	err := s.client.ForEach(func(key, value []byte) error {
		tor := torrent2.TorrentModel{}

		err := json.Unmarshal(value, &tor)
		if err != nil {
			panic(err)
			return err
		}
		newTorrent, err := torrent2.NewTorrent(tor.Magnet, s.conn, torrent2.UP)
		if err != nil {
			return err
		}
		<-newTorrent.Torrent.GotInfo()
		res = append(res, newTorrent)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s storageBbolt) Delete(hash string) error {
	return s.Delete(hash)
}

func NewStorageBbolt(path string, conn *torrent.Client) torrent2.Storage {
	cl, err := client.NewClient(path)
	if err != nil {
		return nil
	}
	return &storageBbolt{conn: conn, client: cl}
}
