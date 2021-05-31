package types

import (
	"encoding/hex"
	"fmt"
)

type ConnectionData struct {
	Address   string
	PubKey    string
	Username  string
	Signature string
}

func NewConnectionData(address, pubKey, value, signature string) *ConnectionData {
	return &ConnectionData{
		Address:   address,
		PubKey:    pubKey,
		Username:  value,
		Signature: signature,
	}
}

// Validate validates the given ConnectionData
func (c *ConnectionData) Validate() error {
	if _, err := hex.DecodeString(c.Address); err != nil {
		return NewWarnErr(fmt.Sprintf("Invalid address format: %s", c.Address))
	}

	if _, err := hex.DecodeString(c.PubKey); err != nil {
		return NewWarnErr(fmt.Sprintf("Invalid public key format: %s", c.PubKey))
	}

	if _, err := hex.DecodeString(c.Signature); err != nil {
		return NewWarnErr(fmt.Sprintf("Invalid signature format: %s", c.Signature))
	}

	return nil
}
