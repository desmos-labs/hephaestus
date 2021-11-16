module github.com/desmos-labs/hephaestus

go 1.15

require (
	github.com/andersfylling/disgord v0.26.1
	github.com/cosmos/cosmos-sdk v0.44.2
	github.com/desmos-labs/cosmos-go-wallet v0.0.0-20211116103831-7f89c57b117e
	github.com/desmos-labs/desmos/v2 v2.2.0-testnet
	github.com/desmos-labs/themis/apis v0.0.0-20210531132313-0b7c43eb5978
	github.com/gin-gonic/gin v1.7.2 // indirect
	github.com/go-playground/validator/v10 v10.7.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hasura/go-graphql-client v0.2.0
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	google.golang.org/grpc v1.41.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/cosmos-sdk => github.com/desmos-labs/cosmos-sdk v0.43.0-alpha1.0.20211012153741-0450cc890f95

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2
