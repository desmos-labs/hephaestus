package types

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
	BotConfig          *BotConfig          `toml:"bot"`
	ThemisConfig       *ThemisConfig       `toml:"themis"`
	VerificationConfig *VerificationConfig `toml:"verification"`
	ChainConfig        *ChainConfig        `toml:"chain"`
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
	Token  string `toml:"token"`
	Prefix string `toml:"prefix"`

	// Path to the PKCS#8-encoded private key of the bot
	PrivateKeyPath string              `toml:"private_key_path"`
	Limitations    []*LimitationConfig `toml:"limitations"`
}

type VerificationConfig struct {
	GraphQLEndpoint       string `toml:"graphql_endpoint"`
	VerifiedUserRole      uint64 `toml:"verified_user_role_id"`
	VerifiedValidatorRole uint64 `toml:"verified_validator_role_id"`
}

// ThemisConfig contains the configuration of the Themis APIs endpoint
type ThemisConfig struct {
	Host string `toml:"host"`
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
