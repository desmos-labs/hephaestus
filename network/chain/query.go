package chain

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// GetBalance returns the balance for the given user and denom
func (c *Client) GetBalance(user string) ([]sdk.Coin, error) {
	res, err := c.bankClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{
		Address: user,
	})
	if err != nil {
		return nil, err
	}

	return res.Balances, nil
}
