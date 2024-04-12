package torrent

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"math/rand/v2"
	"os"
	"runtime"
)

type Config struct {
	Port    int    `yaml:"port"`
	DataDir string `yaml:"data_dir"`
}

var DEFAULT_CONFIG_FILE_PATH = "/etc/gotor.yml"

func genConfigStub() string {
	return fmt.Sprintf(`port: %d
data_dir: "%s"`, rand.IntN(65535-20000)+20000, ".")
}

func NewConfig() (*Config, error) {
	if runtime.GOOS == "windows" {
		DEFAULT_CONFIG_FILE_PATH = os.Getenv("APPDATA") + "/gotor.yml"
	}
	f, err := os.OpenFile(DEFAULT_CONFIG_FILE_PATH, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	if stat, err := os.Stat(DEFAULT_CONFIG_FILE_PATH); err != nil || stat.Size() == 0 {
		_, err = f.WriteString(genConfigStub())
		if err != nil {
			return nil, err
		}
	}
	text, err := io.ReadAll(f)
	cfg := &Config{}
	err = yaml.Unmarshal(text, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}