package cosmos

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/spf13/pflag"
)

// signTx signs the transaction that can be obtained from the given factory and builder,
// and returns the signature that is created
func (client *Client) signTx(factory tx.Factory, builder client.TxBuilder) (signing.SignatureV2, error) {
	// Set an empty signature first
	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   client.privKey.PubKey(),
		Data:     &sigData,
		Sequence: factory.Sequence(),
	}

	err := builder.SetSignatures(sig)
	if err != nil {
		return signing.SignatureV2{}, err
	}

	// Sign the transaction with the private key
	return tx.SignWithPrivKey(
		factory.SignMode(),
		authsigning.SignerData{
			ChainID:       factory.ChainID(),
			AccountNumber: factory.AccountNumber(),
			Sequence:      factory.Sequence(),
		},
		builder,
		client.privKey,
		client.cliCtx.TxConfig,
		factory.Sequence(),
	)
}

// createSignedTx builds and signs a transaction containing the given messages
func (client *Client) createSignedTx(msgs ...sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	factoryCLI := tx.NewFactoryCLI(client.cliCtx, &pflag.FlagSet{}).
		WithFees(client.fees).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// Create the factory
	factory, err := tx.PrepareFactory(client.cliCtx, factoryCLI)
	if err != nil {
		return nil, err
	}

	// Build an unsigned transaction
	builder, err := tx.BuildUnsignedTx(factory, msgs...)
	if err != nil {
		return nil, err
	}

	sig, err := client.signTx(factory, builder)
	if err != nil {
		return nil, err
	}

	// Set the signatures
	err = builder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	// Encode the transaction
	return client.cliCtx.TxConfig.TxEncoder()(builder.GetTx())
}

// BroadcastTx allows to broadcast a transaction containing the given messages
func (client *Client) BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	// Build the transaction
	txBytes, err := client.createSignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	// Broadcast the transaction to a Tendermint node
	return client.cliCtx.BroadcastTx(txBytes)
}
