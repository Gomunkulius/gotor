package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/dustin/go-humanize"
	torrent2 "gotor/internal/torrent"
	"runtime"
	"strconv"
	"time"
)

type TorrentTable struct {
	Torrents []*torrent2.Torrent
	Table    table.Model
	styles   table.Styles
}

func NewTorrentTable(style table.Styles, torrents []*torrent2.Torrent) *TorrentTable {
	columns := []table.Column{
		{Title: "📛Name", Width: 11},
		{Title: "📊Size", Width: 11},
		{Title: "📈Progress", Width: 11},
		{Title: "🔄Status", Width: 11},
		{Title: "🧩Peers", Width: 11},
		{Title: "💨Speed", Width: 11},
	}
	var rows []table.Row
	for _, tor := range torrents {

		speed := tor.Torrent.Info().Length
		percentage := (float32(tor.Torrent.Stats().PiecesComplete) / float32(tor.Torrent.NumPieces())) * 100.0

		written := tor.Torrent.Stats().BytesWritten
		row := table.Row{
			tor.Torrent.Name(),
			fmt.Sprintf("%s", humanize.Bytes(uint64(speed))),
			fmt.Sprintf("%.2f%%", percentage),
			"Up",
			strconv.Itoa(tor.Torrent.Stats().ActivePeers),
			humanize.Bytes(uint64(written.Int64())),
		}

		rows = append(rows, row)
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithStyles(style),
	)

	return &TorrentTable{
		Torrents: torrents,
		Table:    t,
		styles:   style,
	}
}

func (t *TorrentTable) Update(width, height int) {
	columnWidth := (width * 3) / 4
	columns := []table.Column{
		{Title: "📛Name", Width: int(columnWidth / 3)},
		{Title: "📊Size", Width: int(columnWidth / 6)},
		{Title: "📈Progress", Width: (int(columnWidth/6) * 3) / 4},
		{Title: "🔄Status", Width: int(columnWidth/6) / 2},
		{Title: "🧩Peers", Width: int(columnWidth/6) / 2},
		{Title: "💨Speed", Width: int(columnWidth / 6)},
	}
	if runtime.GOOS == "windows" {
		columns = []table.Column{
			{Title: "Name", Width: int(columnWidth / 3)},
			{Title: "Size", Width: int(columnWidth / 6)},
			{Title: "Progress", Width: (int(columnWidth/6) * 3) / 4},
			{Title: "Status", Width: int(columnWidth/6) / 2},
			{Title: "Peers", Width: int(columnWidth/6) / 2},
			{Title: "Speed", Width: int(columnWidth / 6)},
		}
	}
	var rows []table.Row
	for _, tor := range t.Torrents {

		percentage := (float32(tor.Torrent.Stats().PiecesComplete) / float32(tor.Torrent.NumPieces())) * 100.0
		written := tor.Speed1s
		status := "Up"
		if tor.Status == torrent2.PAUSE {
			status = "Pause"
		}
		row := table.Row{
			tor.Torrent.Name(),
			fmt.Sprintf("%s", humanize.Bytes(uint64(tor.Torrent.Length()))),
			fmt.Sprintf("%.2f%%", percentage),
			status,
			strconv.Itoa(tor.Torrent.Stats().ActivePeers),
			fmt.Sprintf("%s/s", humanize.Bytes(uint64(written))),
		}

		rows = append(rows, row)
	}
	tableHeight := (height * 2) / 4
	tab := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
		table.WithStyles(t.styles),
	)
	tab.SetCursor(t.Table.Cursor())
	t.Table = tab
	return
}

func (t *TorrentTable) CountSpeed() {
	for {
		for _, tor := range t.Torrents {
			if tor.Status == torrent2.PAUSE {
				tor.Speed1s = 0
				continue
			}
			before := tor.Torrent.Stats().BytesRead
			time.Sleep(1 * time.Second)
			after := tor.Torrent.Stats().BytesRead
			diff := after.Int64() - before.Int64()
			tor.Speed1s = diff
		}
	}
}
