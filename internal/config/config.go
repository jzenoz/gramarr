package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/tommy647/gramarr/internal/auth"
	"github.com/tommy647/gramarr/internal/bot"
	"github.com/tommy647/gramarr/internal/bot/telegram"
	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/sonarr"
)

type Config struct {
	Telegram telegram.Config `json:"telegram"`
	Auth     auth.Config     `json:"auth"`
	Bot      bot.Config      `json:"bot"`
	Radarr   *radarr.Config  `json:"radarr"`
	Sonarr   *sonarr.Config  `json:"sonarr"`
}

func LoadConfig(configDir string) (*Config, error) {
	configPath := filepath.Join(configDir, "config.json")
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var c = &Config{}
	err = json.NewDecoder(bytes.NewBuffer(file)).Decode(c)
	return c, err
}

// ValidateConfig @todo: implement?
func ValidateConfig(c *Config) error { return nil }
