# Hephaestus
Hephaestus is the official Discord bot for Desmos.

## Configuration
```toml
[chain]
id = "<Chain ID>"
node_uri = "<Node URI>"
bech32_prefix = "<Bech32 prefix>"
fees =  "<Fees to be paid each transaction (optional)>"

[chain.account]
mnemonic = "<Mnemonic phrase>"
hd_path = "<Account HD path>"


[bot]
prefix = "<Command prefix (optional - default '!')"
token = "<Discord bot token>"

[[bot.limitations]]
command = "<Command to limit>"
duration = "<Time for which to limit a command after success (eg. 5m, 8h, 3d)>"
```