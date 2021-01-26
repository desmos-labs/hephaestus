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

func BroadcastTx(cosmosClient *Client, msg sdk.Msg) error {
	// Build the factory CLI
	factoryCLI := tx.NewFactoryCLI(cosmosClient.cliCtx, &pflag.FlagSet{}).
		WithFees(cosmosClient.fees).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// Create the factory
	factory, err := tx.PrepareFactory(cosmosClient.cliCtx, factoryCLI)
	if err != nil {
		return err
	}

	// Build an unsigned transaction
	builder, err := tx.BuildUnsignedTx(factory, msg)
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
		cosmosClient.privKey,
		cosmosClient.cliCtx.TxConfig,
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
	txBytes, err := cosmosClient.cliCtx.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return err
	}

	data := sigs.Data.(*signing.SingleSignatureData)
	fmt.Println(hex.EncodeToString(data.Signature))

	// Broadcast the transaction to a Tendermint node
	res, err := cosmosClient.cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cosmosClient.cliCtx.PrintProto(res)
}
