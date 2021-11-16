package gql

import "github.com/hasura/go-graphql-client"

// applicationLinkQuery represents the GraphQL query to be used to get the application links list
// for a user having a specified Discord username.
type applicationLinkQuery struct {
	ApplicationLinks []struct {
		State       graphql.String `graphql:"state"`
		UserAddress graphql.String `graphql:"user_address"`
	} `graphql:"application_link(where:{username:{_ilike:$username}, application:{_ilike:\"discord\"}})"`
}

// validatorQuery represents the query to be used to get the validator information
// based on a self delegate address.
type validatorQuery struct {
	ValidatorsInfo []struct {
		ConsensusAddress graphql.String `graphql:"consensus_address"`
	} `graphql:"validator_info (where:{self_delegate_address:{_eq:$address}})"`
}
