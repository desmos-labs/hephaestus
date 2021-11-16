package network

import (
	"github.com/andersfylling/disgord"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/desmos-labs/cosmos-go-wallet/cosmos"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/network/chain"
	"github.com/desmos-labs/hephaestus/network/gql"
	"github.com/desmos-labs/hephaestus/network/themis"
	"github.com/desmos-labs/hephaestus/types"
)

type Client struct {
	discordCfg *types.DiscordConfig

	themis  *themis.Client
	graphQL *gql.Client
	network *cosmos.Client
	wallet  *cosmos.Wallet
	chain   *chain.Client
}

func NewClient(cfg *types.NetworkConfig, encodingConfig params.EncodingConfig) (*Client, error) {
	cosmosClient, err := cosmos.NewClient(cfg.Chain.ChainConfig, encodingConfig.Marshaler)
	if err != nil {
		return nil, err
	}

	wallet, err := cosmos.NewWallet(cfg.Account, cosmosClient, encodingConfig.TxConfig)
	if err != nil {
		return nil, err
	}

	themisClient, err := themis.NewClient(cfg.Themis)
	if err != nil {
		return nil, err
	}

	gqlClient, err := gql.NewClient(cfg.Chain.GRPCAddr)
	if err != nil {
		return nil, err
	}

	chainClient, err := chain.NewClient(cosmosClient.GRPCConn)
	if err != nil {
		return nil, err
	}

	return &Client{
		themis:  themisClient,
		graphQL: gqlClient,
		wallet:  wallet,
		network: cosmosClient,
		chain:   chainClient,
	}, nil
}

// ParseAddress parses the given address as a sdk.AccAddress instance
func (n *Client) ParseAddress(address string) (sdk.AccAddress, error) {
	return n.network.ParseAddress(address)
}

// GetBalance returns the balance of the given user
func (n *Client) GetBalance(user string) (sdk.Coins, error) {
	return n.chain.GetBalance(user)
}

// SendTokens sends the specified amount of tokens to the provided user
func (n *Client) SendTokens(user string, amount int64) (*sdk.TxResponse, error) {
	txMsg := &banktypes.MsgSend{
		FromAddress: n.wallet.AccAddress(),
		ToAddress:   user,
		Amount:      sdk.NewCoins(sdk.NewCoin(n.network.GetFeeDenom(), sdk.NewInt(amount))),
	}

	// Send the transaction
	log.Debug().Str(types.LogCommand, types.CmdSend).Str(types.LogRecipient, user).Msg("sending tokens")
	return n.wallet.BroadCastTx(txMsg)
}

// UploadDataToThemis uploads the given data to Themis
func (n *Client) UploadDataToThemis(data *types.ConnectionData) error {
	return n.themis.UploadData(data)
}

// GetDiscordRole returns the role that should be assigned to the Discord user having the given username,
// based on the fact that they have connected their Discord account to a validator or user Desmos Profile
func (n *Client) GetDiscordRole(username string) (disgord.Snowflake, error) {
	isValidator, err := n.graphQL.CheckIsValidator(username)
	if err != nil {
		return 0, err
	}

	if isValidator {
		return disgord.Snowflake(n.discordCfg.VerifiedValidatorRole), nil
	}

	return disgord.Snowflake(n.discordCfg.VerifiedUserRole), nil
}
