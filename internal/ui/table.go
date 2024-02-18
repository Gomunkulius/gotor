package ui

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	"github.com/dustin/go-humanize"
	"strconv"
)

type TorrentTable struct {
	Torrents []*torrent.Torrent
	Table    table.Model
	styles   table.Styles
}

func New(style table.Styles, torrents []*torrent.Torrent) TorrentTable {
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

		speed := tor.Info().Length
		percentage := (float32(tor.Stats().PiecesComplete) / float32(tor.NumPieces())) * 100.0

		written := tor.Stats().BytesWritten
		row := table.Row{
			tor.Name(),
			fmt.Sprintf("%s", humanize.Bytes(uint64(speed))),
			fmt.Sprintf("%.2f%%", percentage),
			"Up",
			strconv.Itoa(tor.Stats().ActivePeers),
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
		percentage := (float32(tor.Stats().PiecesComplete) / float32(tor.NumPieces())) * 100.0
		written := tor.Stats().BytesWritten
		row := table.Row{
			tor.Name(),
			fmt.Sprintf("%s", humanize.Bytes(uint64(tor.Length()))),
			fmt.Sprintf("%.2f%%", percentage),
			"Up",
			strconv.Itoa(tor.Stats().ActivePeers),
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
