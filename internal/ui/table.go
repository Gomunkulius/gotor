package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/dustin/go-humanize"
	torrent2 "gotor/internal/torrent"
	"strconv"
)

type TorrentTable struct {
	Torrents []*torrent2.Torrent
	Table    table.Model
	styles   table.Styles
}

func New(style table.Styles, torrents []*torrent2.Torrent) TorrentTable {
	columns := []table.Column{
		{Title: "ℹ️Name", Width: 10},
		{Title: "📊Size", Width: 10},
		{Title: "📈Progress", Width: 10},
		{Title: "✈️Status", Width: 10},
		{Title: "🧩Peers", Width: 10},
		{Title: "⬆️Up speed", Width: 10},
	}
	rows := []table.Row{}
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

	return TorrentTable{
		Torrents: torrents,
		Table:    t,
		styles:   style,
	}
}

func (t *TorrentTable) Update() {

	columns := []table.Column{
		{Title: "ℹ️Name", Width: 10},
		{Title: "📊Size", Width: 10},
		{Title: "📈Progress", Width: 10},
		{Title: "✈️Status", Width: 10},
		{Title: "🧩Peers", Width: 10},
		{Title: "⬆️Up speed", Width: 10},
	}
	var rows []table.Row
	for _, tor := range t.Torrents {

		percentage := (float32(tor.Torrent.Stats().PiecesComplete) / float32(tor.Torrent.NumPieces())) * 100.0
		written := tor.Torrent.Stats().BytesWritten
		row := table.Row{
			tor.Torrent.Name(),
			fmt.Sprintf("%s", humanize.Bytes(uint64(tor.Torrent.Length()))),
			fmt.Sprintf("%.2f%%", percentage),
			"Up",
			strconv.Itoa(tor.Torrent.Stats().ActivePeers),
			fmt.Sprintf("%s/s", humanize.Bytes(uint64(written.Int64()))),
		}

		rows = append(rows, row)
	}
	tab := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithStyles(t.styles),
	)
	tab.SetCursor(t.Table.Cursor())
	t.Table = tab
	return
}
