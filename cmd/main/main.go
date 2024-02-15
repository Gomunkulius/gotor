package main

import (
	"gotor/internal"
	"gotor/internal/parsing"
	"os"
	"syscall"
)

func main() {
	internal.InitGlobal()
	f, _ := os.OpenFile("debian.torrent", syscall.O_RDWR, 0644)

	torrent, _ := parsing.Open(f)
	internal.Logger.Info("torrent", torrent)
}
