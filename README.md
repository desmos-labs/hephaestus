# Hephaestus
Hephaestus is the official Discord bot for Desmos.

## Configuration
```toml
[bot]
prefix = "<Command prefix (optional - default '!')"
token = "<Discord bot token>"

[account]
mnemonic = "<Mnemonic phrase>"
hd_path = "<Account HD path>"

[chain]
id = "<Chain ID>"
node_uri = "<Node URI>"
bech32_prefix = "<Bech32 prefix>"
fees =  "<Fees to be paid each transaction (optional)>"
```