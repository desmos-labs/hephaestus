package chain

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
)

type Client struct {
	bankClient banktypes.QueryClient
}

func NewClient(conn *grpc.ClientConn) (*Client, error) {
	return &Client{
		bankClient: banktypes.NewQueryClient(conn),
	}, nil
}
