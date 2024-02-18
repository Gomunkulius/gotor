package ui

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	"gotor/internal"
	"math"
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
		index := int(math.Round(math.Log(float64(tor.Info().Length))/math.Log(1000))) - 1

		postfix := internal.SizePostfix[index]
		speed := tor.Info().Length
		percentage := (float32(tor.Stats().PiecesComplete) / float32(tor.NumPieces())) * 100.0
		row := table.Row{
			tor.Name(),
			fmt.Sprintf("%.2f%s", float32(speed)/(float32(math.Pow(1000, float64(index)))), postfix),
			fmt.Sprintf("%.2f%%", percentage),
			"Up",
			strconv.Itoa(tor.Stats().ActivePeers),
			"DEBUG",
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
		index := int(math.Round(math.Log(float64(tor.Info().Length))/math.Log(1000))) - 1

		postfix := internal.SizePostfix[index]
		speed := tor.Info().Length
		percentage := (float32(tor.Stats().PiecesComplete) / float32(tor.NumPieces())) * 100.0
		row := table.Row{
			tor.Name(),
			fmt.Sprintf("%.2f%s", float32(speed)/(float32(math.Pow(1000, float64(index)))), postfix),
			fmt.Sprintf("%.2f%%", percentage),
			"Up",
			strconv.Itoa(tor.Stats().ActivePeers),
			"DEBUG",
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
