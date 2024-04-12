package testing

import (
	"github.com/anacrolix/torrent"
	"math/rand/v2"
)

// MockClient Almost mock client :)
// This is real client but with random port
func MockClient() (*torrent.Client, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.ListenPort = rand.IntN(65535-20000) + 20000
	c, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}
