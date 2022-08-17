package gql

import graphql "github.com/hasura/go-graphql-client"

// Client represents a GraphQL client
type Client struct {
	chainClient  *graphql.Client
	desmosClient *graphql.Client
}

// NewClient returns a new Client instance
func NewClient(chainEndpoint string, desmosEndpoint string) (*Client, error) {
	return &Client{
		chainClient:  graphql.NewClient(chainEndpoint, nil),
		desmosClient: graphql.NewClient(desmosEndpoint, nil),
	}, nil
}
