package gql

import (
	"context"
	"fmt"
	"strings"

	graphql "github.com/hasura/go-graphql-client"

	"github.com/desmos-labs/hephaestus/types"
)

// GetDiscordLink returns the ApplicationLink representing the connection between the Desmos Profile and Discord.
// It returns an error if the link was not found, or it's not in the correct state.
func (c *Client) GetDiscordLink(username string) (*ApplicationLink, error) {
	// Build the query and the arguments
	var linkQuery applicationLinkQuery
	variables := map[string]interface{}{
		"username": graphql.String(fmt.Sprintf("%%%s%%", strings.ToLower(username))),
	}

	err := c.desmosClient.Query(context.Background(), &linkQuery, variables)
	if err != nil {
		return nil, types.NewWarnErr("Error while querying the server: %s", err)
	}

	if len(linkQuery.ApplicationLinks) == 0 {
		return nil, nil
	}

	return linkQuery.ApplicationLinks[0], nil
}

// CheckIsValidator checks whether the user having the given username is a validator or not
// based on the data present on the specific GraphQL endpoint.
func (c *Client) CheckIsValidator(appLink *ApplicationLink) (bool, error) {
	var qry validatorQuery
	variables := map[string]interface{}{
		"address": appLink.UserAddress,
	}

	err := c.chainClient.Query(context.Background(), &qry, variables)
	if err != nil {
		return false, types.NewWarnErr("Error while querying the validator info: %s", err)
	}

	return len(qry.ValidatorsInfo) > 0, nil
}
