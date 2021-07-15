package gql_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/hephaestus/gql"
)

func TestCheckIsValidator(t *testing.T) {
	result, err := gql.CheckIsValidator("https://gql.morpheus.desmos.network/v1/graphql", "Lucag__#5237")
	require.NoError(t, err)
	require.False(t, result)
}
