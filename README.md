# Hephaestus

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/desmos-labs/hephaestus)](https://github.com/desmos-labs/hephaestus/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/desmos-labs/hephaestus/.svg)](https://pkg.go.dev/github.com/desmos-labs/hephaestus/)
[![Go Report](https://goreportcard.com/badge/github.com/desmos-labs/hephaestus)](https://goreportcard.com/report/github.com/desmos-labs/hephaestus)
![License](https://img.shields.io/github/license/desmos-labs/hephaestus.svg)
[![Codecov](https://codecov.io/gh/desmos-labs/hephaestus/branch/main/graph/badge.svg)](https://codecov.io/gh/desmos-labs/hephaestus/branch/main)
[![Tests status](https://github.com/desmos-labs/hephaestus/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/desmos-labs/hephaestus/actions/workflows/tests.yml?query=branch%3Amain)
[![Lint status](https://github.com/desmos-labs/hephaestus/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/desmos-labs/hephaestus/actions/workflows/lint.yml?query=branch%3Amain)

Hephaestus is the official Discord bot for Desmos.

## Configuration
```yaml
networks:
  testnet:
    chain:
      id: "morpheus-apollo-2"
      bech32_prefix: "desmos"
      rpc_addr: "https://rpc.morpheus.desmos.network:443"
      grpc_addr: "https://grpc.morpheus.desmos.network:443"
      chain_graphql_addr: "https://gql-morpheus.desmos.forbole.com/v1/graphql"
      djuno_graphql_addr: "https://gql.morpheus.desmos.network/v1/graphql"
      gas_price: "0.1udaric"

    account:
      mnemonic: "<Faucet mnemonic account>"
      hd_path: "m/44'/852'/0'/0/0"

    themis:
      host: "https://themis.morpheus.desmos.network"
      private_key_path: "<Themis private key file path>"

    discord:
      verified_user_role_id: "<Discord verified user role ID>"
      verified_validator_role_id: "<Discord verified validator role ID>"


# ---------------------------------------------------------------------------------------------------------------

bot:
  token: "<Discord bot token>"
  prefix: "!"

  limitations:
    - command: "help"
      duration: "72h" # 3 days

    - command: "docs"
      duration: "15m" # 15 minutes

    - command: "send"
      duration: "0m" # 7 Days
```