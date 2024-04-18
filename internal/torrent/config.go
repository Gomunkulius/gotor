package torrent

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Port    int    `yaml:"port"`
	DataDir string `yaml:"data_dir"`
}

var DEFAULT_CONFIG_FILE_PATH = ".config/gotor.yml"

func genConfigStub() string {
	var configPath string
	namePrompt := &survey.Input{
		Message: "Enter path where you want to save torrents: (Enter for default)",
	}
	survey.Ask([]*survey.Question{{
		Name:     "Config",
		Prompt:   namePrompt,
		Validate: nil,
	}}, &configPath)
	if configPath == "" {
		if runtime.GOOS == "windows" {
			configPath = os.Getenv("APPDATA")
		} else {
			configPath = filepath.Join(os.Getenv("HOME"))
		}
	}
	return fmt.Sprintf(`port: %d
data_dir: "%s"`, rand.IntN(65535-20000)+20000, configPath)
}

func NewConfig() (*Config, error) {
	if runtime.GOOS == "windows" {
		DEFAULT_CONFIG_FILE_PATH = os.Getenv("APPDATA") + "/gotor.yml"
	} else {
		DEFAULT_CONFIG_FILE_PATH = filepath.Join(os.Getenv("HOME"), DEFAULT_CONFIG_FILE_PATH)
	}
	if stat, err := os.Stat(DEFAULT_CONFIG_FILE_PATH); err != nil || stat.Size() == 0 {
		f, err := os.Create(DEFAULT_CONFIG_FILE_PATH)
		if err != nil {
			return nil, err
		}
		_, err = f.WriteString(genConfigStub())
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(DEFAULT_CONFIG_FILE_PATH, os.O_RDWR|os.O_CREATE, 0777)
	text, err := io.ReadAll(f)
	cfg := &Config{}
	err = yaml.Unmarshal(text, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
