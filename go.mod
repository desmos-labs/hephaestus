module github.com/desmos-labs/hephaestus

go 1.15

require (
	github.com/andersfylling/disgord v0.26.1
	github.com/cosmos/cosmos-sdk v0.44.4
	github.com/desmos-labs/cosmos-go-wallet v0.0.0-20211116103831-7f89c57b117e
	github.com/desmos-labs/desmos/v2 v2.3.1
	github.com/desmos-labs/themis/apis v0.0.0-20220111092734-c6a1b17a2b0c
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hasura/go-graphql-client v0.2.0
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/sys v0.0.0-20220110181412-a018aaa089fe // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/grpc v1.42.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/cosmos-sdk => github.com/desmos-labs/cosmos-sdk v0.43.0-alpha1.0.20211012153741-0450cc890f95

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2
