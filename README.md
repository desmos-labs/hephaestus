# Hephaestus
Hephaestus is the official Discord bot for Desmos.

## Configuration
```toml
[chain]
id = "<Chain ID>"
rpc_addr = "<Address of the RPC>"
grpc_addr = "<Address of the gRPC>"
gas_price =  "<Minimum gas price to be used>"

[account]
bech32_prefix = "<Bech32 prefix>"
mnemonic = "<Mnemonic phrase>"
hd_path = "<Account HD path>"


[bot]
prefix = "<Command prefix (optional - default '!')"
token = "<Discord bot token>"

[[bot.limitations]]
command = "<Command to limit>"
duration = "<Time for which to limit a command after success (eg. 5m, 8h, 3d)>"
```