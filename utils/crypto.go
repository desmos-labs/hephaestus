package utils

import (
	"crypto/rsa"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// ReadPrivateKeyFromFile reads an RSA private key from the file located at the given path
func ReadPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	bz, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parseResult, err := ssh.ParseRawPrivateKey(bz)
	if err != nil {
		return nil, err
	}

	return parseResult.(*rsa.PrivateKey), nil
}
