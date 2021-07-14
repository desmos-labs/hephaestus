module github.com/desmos-labs/hephaestus

go 1.15

require (
	github.com/andersfylling/disgord v0.26.1
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/desmos-labs/desmos v0.17.2
	github.com/desmos-labs/themis/apis v0.0.0-20210531132313-0b7c43eb5978
	github.com/gin-gonic/gin v1.7.2 // indirect
	github.com/go-playground/validator/v10 v10.7.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/pelletier/go-toml v1.8.1
	github.com/rs/zerolog v1.21.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/cosmos-sdk => github.com/desmos-labs/cosmos-sdk v0.42.5-0.20210712073217-87acd62da7d7

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
