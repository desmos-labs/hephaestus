package config

import (
	"io/ioutil"

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
	BotConfig   *BotConfig     `toml:"bot"`
	Account     *AccountConfig `toml:"account"`
	ChainConfig *ChainConfig   `toml:"chain"`
}

type BotConfig struct {
	Token  string `toml:"token"`
	Prefix string `toml:"prefix"`
}

type AccountConfig struct {
	Mnemonic string `toml:"mnemonic"`
	HDPath   string `toml:"hd_path"`
}

type ChainConfig struct {
	NodeURI      string `toml:"node_uri"`
	Bech32Prefix string `toml:"bech32_prefix"`
	ChainID      string `toml:"id"`
	Fees         string `toml:"fees"`
}
