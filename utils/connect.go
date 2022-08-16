package utils

import (
	"encoding/hex"
	"encoding/json"

	signcmd "github.com/desmos-labs/desmos/v2/app/desmos/cmd/sign"
	"github.com/desmos-labs/hephaestus/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func GetSignatureData(jsonData string) (*signcmd.SignatureData, error) {
	var signatureData signcmd.SignatureData
	err := json.Unmarshal([]byte(jsonData), &signatureData)
	if err != nil {
		return nil, types.NewWarnErr("Invalid data provided: %s", err)
	}

	// Verify the signature
	pubKeyBz, err := hex.DecodeString(signatureData.PubKey)
	if err != nil {
		return nil, types.NewWarnErr("Error while reading public key: %s", err)
	}

	valueBz, err := hex.DecodeString(signatureData.Value)
	if err != nil {
		return nil, types.NewWarnErr("Error while reading value: %s", err)
	}

	sigBz, err := hex.DecodeString(signatureData.Signature)
	if err != nil {
		return nil, types.NewWarnErr("Error while reading signature: %s", err)
	}

	pubKey := secp256k1.PubKey(pubKeyBz)
	if !pubKey.VerifySignature(valueBz, sigBz) {
		return nil, types.NewWarnErr("Invalid signature. Make sure you have signed the message using the correct account")
	}

	return &signatureData, nil
}
