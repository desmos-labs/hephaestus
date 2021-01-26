package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
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
	Token       string         `toml:"token"`
	Account     *AccountConfig `toml:"account"`
	ChainConfig *ChainConfig   `toml:"chain"`
}

type AccountConfig struct {
	Mnemonic string `toml:"mnemonic"`
	HDPath   string `toml:"hd_path"`
}

type ChainConfig struct {
	NodeURI      string `toml:"node_uri"`
	Bech32Prefix string `toml:"bech32_prefix"`
	ChainID      string `toml:"chain_id"`
	Fees         string `toml:"fees"`
}
