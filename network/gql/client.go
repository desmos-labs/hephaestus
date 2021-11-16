package gql

import "github.com/hasura/go-graphql-client"

// Client represents a GraphQL client
type Client struct {
	client *graphql.Client
}

// NewClient returns a new Client instance
func NewClient(endpoint string) (*Client, error) {
	client := graphql.NewClient(endpoint, nil)

	return &Client{
		client: client,
	}, nil
}
