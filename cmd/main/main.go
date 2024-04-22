package main

import (
	"fmt"
	log2 "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"gotor/internal"
	torrent2 "gotor/internal/torrent"
	"gotor/internal/torrent/local"
	"gotor/internal/ui"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func main() {
	app := &cli.App{
		Name:        "gotor",
		Usage:       "Launch a torrent client",
		Description: "The least functional and simplest torrent client you've ever seen",
		Action:      setup,
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Opens configuration file in editor",
				Action:  openEditor,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "editor",
						Aliases: []string{"e"},
						Usage:   "Set editor to open config file",
					},
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Shows version of the program",
				Action: func(context *cli.Context) error {
					fmt.Printf("GOTOR version: %s\n", internal.VERSION)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "set-dir",
				Usage:  "set data directory",
				Action: setDirectory,
			},
			&cli.StringFlag{
				Name:   "set-port",
				Usage:  "set port for input connections",
				Action: setPort,
			},
		},
		EnableBashCompletion: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setup(ctx *cli.Context) error {
	internal.InitGlobal()

	// anacrolix dumbass!
	if len(log2.Default.Handlers) != 0 {
		log2.Default.Handlers[0] = internal.MyHandler{}
	}
	log.SetOutput(io.Discard)

	gcfg, err := torrent2.NewConfig()
	if err != nil || gcfg == nil {
		return fmt.Errorf("cant create config err %v", err)
	}
	cfg := torrent.NewDefaultClientConfig() // TODO: config
	cfg.DataDir = gcfg.DataDir
	cfg.ListenPort = gcfg.Port
	cfg.Logger = log2.Logger{}
	cfg.Debug = false
	c, err := torrent.NewClient(cfg)
	defer c.Close()
	if err != nil || c == nil {
		fmt.Printf("cant connect err %v\n", err)
		cfg.ListenPort = rand.IntN(65535-20000) + 20000
		crand, err := torrent.NewClient(cfg)
		if err != nil || crand == nil {
			return fmt.Errorf("cant connect on port %d err %v", cfg.ListenPort, err)
		}
		c = crand
		fmt.Printf("new non default client created, listening on port %d\n", cfg.ListenPort)
	}
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("247")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	storage := local.NewStorageBbolt("bolt.db", c)
	if storage == nil {
		return fmt.Errorf("Can't create storage")
	}
	torrents, err := storage.GetAll()
	files := torrent2.InitTorrents(torrents, c)
	if err != nil {
		println("cant init torrents")
		return fmt.Errorf("Can't init torrents")
	}
	torTable := ui.NewTorrentTable(s, files)
	go torTable.CountSpeed()
	m := ui.NewModel(torTable, c, storage, gcfg)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}
	fmt.Println("cant start")
	return nil
}

func openEditor(ctx *cli.Context) error {
	editor := "vim"
	if ctx.String("editor") != "" {
		editor = ctx.String("editor")
	}
	_, err := torrent2.NewConfig()
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		editor = "notepad.exe"
	}
	cmd := exec.Command(editor, torrent2.DEFAULT_CONFIG_FILE_PATH)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func setDirectory(context *cli.Context, s string) error {
	cfg, err := torrent2.NewConfig()
	if err != nil {
		return err
	}
	if s != "" {
		cfg.DataDir = s
	}
	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(torrent2.DEFAULT_CONFIG_FILE_PATH, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Write(yml)
	if err != nil {
		return err
	}
	return nil
}

func setPort(_ *cli.Context, s string) error {
	cfg, err := torrent2.NewConfig()
	if err != nil {
		return err
	}
	if s != "" {
		port, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Port must be an integer (0...65535)")
			return err
		}
		cfg.Port = port
	}
	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(torrent2.DEFAULT_CONFIG_FILE_PATH, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Write(yml)
	if err != nil {
		return err
	}
	return nil
}
