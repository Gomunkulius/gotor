package ui

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ncruces/zenity"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"io"
	"log"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := internal.InputStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return internal.SelChooseStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type ChooseModel struct {
	width   int
	height  int
	listmod list.Model
	table   *TorrentTable
	storage torrent2.Storage
	conn    *torrent.Client
}

func (m ChooseModel) Init() tea.Cmd {
	return nil
}

func (m ChooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.listmod.Cursor() {
			case 0:
				return m, ExitCmd(Input)
			case 1:
				path, err := zenity.SelectFile(
					zenity.Filename("."),
					zenity.FileFilters{
						{"Torrent files", []string{"*.torrent"}, false},
					})
				if err != nil {
					return m, nil
				}
				t, err := torrent2.NewTorrentFromFile(path, m.conn, torrent2.UP)
				if err != nil {
					return m, nil
				}
				<-t.Torrent.GotInfo()
				if err != nil {
					return m, nil
				}
				for _, t2 := range m.table.Torrents {
					if t.Torrent.InfoHash() == t2.Torrent.InfoHash() {
						return m, ExitCmd(Main)
					}
				}
				m.table.Torrents = append(m.table.Torrents, t)
				go torrent2.DownloadTorrent(t)
				return m, tea.Batch(SaveExitCmd(m.storage, t), tickEvery())

			}
		case "esc":
			log.Printf("esc pressed")
			return m, ExitCmd(Main)
		}

	}
	var cmd tea.Cmd
	m.listmod, cmd = m.listmod.Update(msg)
	return m, cmd
}

func (m ChooseModel) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, internal.InputStyle.Render("How you want add torrent?"+"\n"),
			internal.InputStyle.Render(m.listmod.View())),
	)
}

func NewChooseModel(width, height int, table *TorrentTable, storage torrent2.Storage, conn *torrent.Client) *ChooseModel {
	items := []list.Item{
		item("Enter magnet link"),
		item("Select file"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, 14)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	return &ChooseModel{
		width:   width,
		height:  height,
		listmod: l,
		table:   table,
		storage: storage,
		conn:    conn,
	}
}
