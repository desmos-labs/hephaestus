package types

import (
	"os"
	"time"

	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"
	"gopkg.in/yaml.v3"
)

// Parse allows to parse the file at the provided path into a Config object
func Parse(filePath string) (*Config, error) {
	bz, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	return &config, yaml.Unmarshal(bz, &config)
}

// ------------------------------------------------------------------------------------------------------------------

// Config contains all the configuration data
type Config struct {
	Networks  *NetworksConfig `yaml:"networks"`
	BotConfig *BotConfig      `yaml:"bot"`
}

type NetworksConfig struct {
	Testnet *NetworkConfig `yaml:"testnet"`
	Mainnet *NetworkConfig `yaml:"mainnet"`
}

// NetworkConfig contains the configuration about the Desmos network to be used
type NetworkConfig struct {
	Chain   *ChainConfig               `yaml:"chain"`
	Account *wallettypes.AccountConfig `yaml:"account"`
	Themis  *ThemisConfig              `yaml:"themis"`
	Discord *DiscordConfig             `yaml:"discord"`
}

// ChainConfig wraps the wallet ChainConfig structure adding the GraphQL endpoint
type ChainConfig struct {
	*wallettypes.ChainConfig `yaml:"-,inline"`
	ChainGraphQL             string `yaml:"chain_graphql_addr"`
	DesmosGraphQL            string `yaml:"djuno_graphql_addr"`
}

// ThemisConfig contains the configuration of the Themis APIs endpoint
type ThemisConfig struct {
	// Themis host URL to where to upload data
	Host string `yaml:"host"`

	// Path to the PKCS#8-encoded private key of the bot
	PrivateKeyPath string `yaml:"private_key_path"`
}

// DiscordConfig contains the configuration about the Discord server
type DiscordConfig struct {
	VerifiedUserRole      uint64 `yaml:"verified_user_role_id"`
	VerifiedValidatorRole uint64 `yaml:"verified_validator_role_id"`
}

// BotConfig contains the configuration about the bot
type BotConfig struct {
	Token       string              `yaml:"token"`
	Prefix      string              `yaml:"prefix"`
	Limitations []*LimitationConfig `yaml:"limitations"`
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
	Command  string        `yaml:"command"`
	Duration time.Duration `yaml:"duration"`
}
