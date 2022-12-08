package themis

import (
	"crypto/rsa"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/desmos-labs/hephaestus/types"
)

// Client represents a Themis client
type Client struct {
	host    string
	privKey *rsa.PrivateKey
}

// NewClient returns a new Client instance
func NewClient(cfg *types.ThemisConfig) (*Client, error) {
	bz, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	parseResult, err := ssh.ParseRawPrivateKey(bz)
	if err != nil {
		return nil, err
	}

	return &Client{
		host:    cfg.Host,
		privKey: parseResult.(*rsa.PrivateKey),
	}, nil
}
