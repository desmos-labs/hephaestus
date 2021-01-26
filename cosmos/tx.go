package cosmos

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/spf13/pflag"
)

// BroadcastTx allows to broadcast a transaction containing the given messages
func (client *Client) BroadcastTx(msgs ...sdk.Msg) error {
	// Build the factory CLI
	factoryCLI := tx.NewFactoryCLI(client.cliCtx, &pflag.FlagSet{}).
		WithFees(client.fees).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// Create the factory
	factory, err := tx.PrepareFactory(client.cliCtx, factoryCLI)
	if err != nil {
		return err
	}

	// Build an unsigned transaction
	builder, err := tx.BuildUnsignedTx(factory, msgs...)
	if err != nil {
		return err
	}

	// Sign the transaction with the private key
	sigs, err := tx.SignWithPrivKey(
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
	if err != nil {
		return err
	}

	// Set the signatures
	err = builder.SetSignatures(sigs)
	if err != nil {
		return err
	}

	// Encode the transaction
	txBytes, err := client.cliCtx.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return err
	}

	data := sigs.Data.(*signing.SingleSignatureData)
	fmt.Println(hex.EncodeToString(data.Signature))

	// Broadcast the transaction to a Tendermint node
	res, err := client.cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return client.cliCtx.PrintProto(res)
}
