package config

import (
	"io/ioutil"
	"time"

	"github.com/pelletier/go-toml"
)

// Parse allows to parse the file at the provided path into a Config object
func Parse(filePath string) (*Config, error) {
	bz, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	return &config, toml.Unmarshal(bz, &config)
}

// ------------------------------------------------------------------------------------------------------------------

// Config contains all the configuration data
type Config struct {
	BotConfig   *BotConfig   `toml:"bot"`
	ChainConfig *ChainConfig `toml:"chain"`
}

type ChainConfig struct {
	NodeURI      string         `toml:"node_uri"`
	Bech32Prefix string         `toml:"bech32_prefix"`
	ChainID      string         `toml:"id"`
	Fees         string         `toml:"fees"`
	Account      *AccountConfig `toml:"account"`
}

type AccountConfig struct {
	Mnemonic string `toml:"mnemonic"`
	HDPath   string `toml:"hd_path"`
}

type BotConfig struct {
	Token       string              `toml:"token"`
	Prefix      string              `toml:"prefix"`
	Limitations []*LimitationConfig `toml:"limitations"`
}

func (cfg *BotConfig) FindLimitationByCommand(command string) *LimitationConfig {
	for _, limitation := range cfg.Limitations {
		if limitation.Command == command {
			return limitation
		}
	}
	return nil
}

type LimitationConfig struct {
	Command  string        `toml:"command"`
	Duration time.Duration `toml:"duration"`
}
