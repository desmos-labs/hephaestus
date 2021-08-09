package gql

import (
	"context"
	"fmt"
	"strings"

	"github.com/hasura/go-graphql-client"

	"github.com/desmos-labs/hephaestus/types"
)

// ApplicationLinkQuery represents the GraphQL query to be used to get the application links list
// for a user having a specified Discord username.
type ApplicationLinkQuery struct {
	ApplicationLinks []struct {
		State       graphql.String `graphql:"state"`
		UserAddress graphql.String `graphql:"user_address"`
	} `graphql:"application_link(where:{username:{_ilike:$username}, application:{_ilike:\"discord\"}})"`
}

// ValidatorQuery represents the query to be used to get the validator information
// based on a self delegate address.
type ValidatorQuery struct {
	ValidatorsInfo []struct {
		ConsensusAddress graphql.String `graphql:"consensus_address"`
	} `graphql:"validator_info (where:{self_delegate_address:{_eq:$address}})"`
}

// CheckIsValidator checks whether the user having the given username is a validator or not
// based on the data present on the specific GraphQL endpoint.
func CheckIsValidator(endpoint string, username string) (bool, error) {
	// Build the query and the arguments
	var linkQuery ApplicationLinkQuery
	variables := map[string]interface{}{
		"username": graphql.String(fmt.Sprintf("%%%s%%", strings.ToLower(username))),
	}

	client := graphql.NewClient(endpoint, nil)
	err := client.Query(context.Background(), &linkQuery, variables)
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

	var validatorQuery ValidatorQuery
	variables = map[string]interface{}{
		"address": applicationLink.UserAddress,
	}

	err = client.Query(context.Background(), &validatorQuery, variables)
	if err != nil {
		return false, types.NewWarnErr("Error while querying the validator info: %s", err)
	}

	return len(validatorQuery.ValidatorsInfo) > 0, nil
}
