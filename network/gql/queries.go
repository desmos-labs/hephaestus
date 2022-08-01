package gql

import (
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/hasura/go-graphql-client"
)

// applicationLinkQuery represents the GraphQL query to be used to get the application links list
// for a user having a specified Discord username.
type applicationLinkQuery struct {
	ApplicationLinks []*ApplicationLink `graphql:"application_link(where:{username:{_ilike:$username}, application:{_ilike:\"discord\"}})"`
}

type ApplicationLink struct {
	State       graphql.String `graphql:"state"`
	UserAddress graphql.String `graphql:"user_address"`
}

func (a *ApplicationLink) IsValid() bool {
	return a.State == graphql.String(profilestypes.AppLinkStateVerificationSuccess)
}

// validatorQuery represents the query to be used to get the validator information
// based on a self delegate address.
type validatorQuery struct {
	ValidatorsInfo []struct {
		ConsensusAddress graphql.String `graphql:"consensus_address"`
	} `graphql:"validator_info (where:{self_delegate_address:{_eq:$address}})"`
}
