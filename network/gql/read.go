package gql

import (
	"context"
	"fmt"
	"strings"

	"github.com/hasura/go-graphql-client"

	"github.com/desmos-labs/hephaestus/types"
)

// CheckIsValidator checks whether the user having the given username is a validator or not
// based on the data present on the specific GraphQL endpoint.
func (c *Client) CheckIsValidator(username string) (bool, error) {
	// Build the query and the arguments
	var linkQuery applicationLinkQuery
	variables := map[string]interface{}{
		"username": graphql.String(fmt.Sprintf("%%%s%%", strings.ToLower(username))),
	}

	err := c.desmosClient.Query(context.Background(), &linkQuery, variables)
	if err != nil {
		return false, types.NewWarnErr("Error while querying the server: %s", err)
	}

	if len(linkQuery.ApplicationLinks) == 0 {
		return false, types.NewWarnErr(`No link found for your account. 
Please make sure you create a Desmos profile and connect your Discord account first.
Use the `+"`!%s`"+`command to know more.`, types.CmdConnect)
	}

	applicationLink := linkQuery.ApplicationLinks[0]
	if applicationLink.State != "APPLICATION_LINK_STATE_VERIFICATION_SUCCESS" {
		return false, types.NewWarnErr("Invalid link status: %s. Try reconnecting your Discord account.",
			applicationLink.State)
	}

	var qry validatorQuery
	variables = map[string]interface{}{
		"address": applicationLink.UserAddress,
	}

	err = c.chainClient.Query(context.Background(), &qry, variables)
	if err != nil {
		return false, types.NewWarnErr("Error while querying the validator info: %s", err)
	}

	return len(qry.ValidatorsInfo) > 0, nil
}
