package cosmos

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/discord-bot/config"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Client struct {
	cliCtx client.Context

	privKey cryptotypes.PrivKey

	fees string
}

func NewClient(accCfg *config.AccountConfig, chainCfg *config.ChainConfig) (*Client, error) {
	// Get the private keys
	algo := hd.Secp256k1
	derivedPriv, err := algo.Derive()(accCfg.Mnemonic, "", accCfg.HDPath)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(derivedPriv)

	// Build the RPC client
	rpcClient, err := rpchttp.New(chainCfg.NodeURI, "/websocket")
	if err != nil {
		return nil, err
	}

	// Build the config
	prefix := chainCfg.Bech32Prefix
	sdkCfg := sdk.GetConfig()
	sdkCfg.SetBech32PrefixForAccount(prefix, prefix+sdk.PrefixPublic)
	sdkCfg.SetBech32PrefixForValidator(
		prefix+sdk.PrefixValidator+sdk.PrefixOperator,
		prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	sdkCfg.SetBech32PrefixForConsensusNode(
		prefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		prefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)
	sdkCfg.Seal()

	// Build the context
	encodingConfig := app.MakeTestEncodingConfig()
	cliCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithClient(rpcClient).
		WithChainID(chainCfg.ChainID).
		WithFromAddress(sdk.AccAddress(privKey.PubKey().Address()))

	return &Client{
		cliCtx:  cliCtx,
		privKey: privKey,
		fees:    chainCfg.Fees,
	}, nil
}

func (client *Client) Address() string {
	return sdk.AccAddress(client.privKey.PubKey().Address()).String()
}
