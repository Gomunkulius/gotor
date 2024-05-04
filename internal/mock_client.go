package internal

import (
	"fmt"
	"math/rand/v2"

	"github.com/anacrolix/torrent"
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
	if c == nil {
		return nil, fmt.Errorf("Shit")
	}
	return c, nil
}
