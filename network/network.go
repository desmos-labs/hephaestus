package network

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/desmos-labs/cosmos-go-wallet/client"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"
	"github.com/desmos-labs/cosmos-go-wallet/wallet"
	"github.com/desmos-labs/desmos/v4/app/desmos/cmd/sign"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/utils"

	"github.com/desmos-labs/hephaestus/network/chain"
	"github.com/desmos-labs/hephaestus/network/gql"
	"github.com/desmos-labs/hephaestus/network/themis"
	"github.com/desmos-labs/hephaestus/types"
)

type Client struct {
	discordCfg *types.DiscordConfig

	themis  *themis.Client
	graphQL *gql.Client
	network *client.Client
	wallet  *wallet.Wallet
	chain   *chain.Client
}

func NewClient(cfg *types.NetworkConfig, encodingConfig params.EncodingConfig) (*Client, error) {
	cosmosClient, err := client.NewClient(cfg.Chain.ChainConfig, encodingConfig.Marshaler)
	if err != nil {
		return nil, err
	}

	wallet, err := wallet.NewWallet(cfg.Account, cosmosClient, encodingConfig.TxConfig)
	if err != nil {
		return nil, err
	}

	themisClient, err := themis.NewClient(cfg.Themis)
	if err != nil {
		return nil, err
	}

	gqlClient, err := gql.NewClient(cfg.Chain.ChainGraphQL, cfg.Chain.DesmosGraphQL)
	if err != nil {
		return nil, err
	}

	chainClient, err := chain.NewClient(cosmosClient.GRPCConn)
	if err != nil {
		return nil, err
	}

	return &Client{
		discordCfg: cfg.Discord,

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
	sender, err := sdk.AccAddressFromBech32(n.wallet.AccAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid sender address: %s", err)
	}

	receiver, err := sdk.AccAddressFromBech32(user)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %s", err)
	}

	data := wallettypes.NewTransactionData(banktypes.NewMsgSend(
		sender,
		receiver,
		sdk.NewCoins(sdk.NewCoin(n.network.GetFeeDenom(), sdk.NewInt(amount))),
	)).WithGasAuto().WithFeeAuto().WithMemo("Sent from Hephaestus")

	// Send the transaction
	log.Debug().Str(types.LogCommand, types.CmdSend).Str(types.LogRecipient, user).Msg("sending tokens")
	return n.wallet.BroadcastTxSync(data)
}

// UploadDataToThemis uploads the given data to Themis
func (n *Client) UploadDataToThemis(username string, data *sign.SignatureData) error {
	return n.themis.UploadData(username, data)
}

// GetDiscordRole returns the role that should be assigned to the Discord user having the given username,
// based on the fact that they have connected their Discord account to a validator or user Desmos Profile
func (n *Client) GetDiscordRole(username string) (disgord.Snowflake, error) {
	discordLink, err := n.graphQL.GetDiscordLink(username)
	if err != nil {
		return disgord.Snowflake(0), err
	}

	if discordLink == nil {
		return disgord.Snowflake(0), types.NewWarnErr(`No link found for your account. 
Please make sure you create a Desmos profile and connect your Discord account first.
Use the `+"`!%s`"+`command to know more.`, types.CmdConnect)
	}

	if !discordLink.IsValid() {
		return disgord.Snowflake(0), types.NewWarnErr(
			"Invalid link status: %s. Try reconnecting your Discord account.", discordLink.State)
	}

	isValidator, err := n.graphQL.CheckIsValidator(discordLink)
	if err != nil {
		return disgord.Snowflake(0), err
	}

	if isValidator {
		return disgord.Snowflake(n.discordCfg.VerifiedValidatorRole), nil
	}

	return disgord.Snowflake(n.discordCfg.VerifiedUserRole), nil
}

// CleanRoles cleans the roles of all the members that are not verified anymore
func (n *Client) CleanRoles(s disgord.Session) error {
	for _, serverID := range s.GetConnectedGuilds() {
		// Get all the members
		members, err := s.Guild(serverID).GetMembers(&disgord.GetMembersParams{})
		if err != nil {
			return types.NewWarnErr("error while getting members: %s", err)
		}

		for _, member := range members {
			// Check if the user is verified
			isVerified, role := n.isVerified(member)
			if !isVerified {
				continue
			}

			// Get the application link
			discordLink, err := n.graphQL.GetDiscordLink(utils.GetUserUsername(member.User))
			if err != nil {
				return err
			}

			// If the link does not exist anymore, or it's invalid, remove it
			if discordLink == nil || !discordLink.IsValid() {
				log.Debug().Str("user", utils.GetUserUsername(member.User)).Msg("removing verified role")
				err = s.Guild(serverID).Member(member.UserID).RemoveRole(role)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// isVerified tells whether the given member is verified or not.
// If it's verified, it returns the verification role that they have assigned as well.
func (n *Client) isVerified(member *disgord.Member) (bool, disgord.Snowflake) {
	verifiedRoles := []disgord.Snowflake{
		disgord.Snowflake(n.discordCfg.VerifiedUserRole),
		disgord.Snowflake(n.discordCfg.VerifiedValidatorRole),
	}

	for _, assignedRole := range member.Roles {
		for _, verifyingRole := range verifiedRoles {
			if assignedRole == verifyingRole {
				return true, verifyingRole
			}
		}
	}

	return false, disgord.Snowflake(0)
}
