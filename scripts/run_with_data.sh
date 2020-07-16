#!/usr/bin/env bash

PASSWORD="12345678"
GAS_PRICES="0.025stake"

bondsd init local --chain-id bondschain-1

yes $PASSWORD | bondscli keys delete miguel --force
yes $PASSWORD | bondscli keys delete francesco --force
yes $PASSWORD | bondscli keys delete shaun --force
yes $PASSWORD | bondscli keys delete fee --force
yes $PASSWORD | bondscli keys delete fee2 --force
yes $PASSWORD | bondscli keys delete fee3 --force
yes $PASSWORD | bondscli keys delete fee4 --force
yes $PASSWORD | bondscli keys delete fee5 --force

yes $PASSWORD | bondscli keys add miguel
yes $PASSWORD | bondscli keys add francesco
yes $PASSWORD | bondscli keys add shaun
yes $PASSWORD | bondscli keys add fee
yes $PASSWORD | bondscli keys add fee2
yes $PASSWORD | bondscli keys add fee3
yes $PASSWORD | bondscli keys add fee4
yes $PASSWORD | bondscli keys add fee5

# Note: important to add 'miguel' as a genesis-account since this is the chain's validator
yes $PASSWORD | bondsd add-genesis-account "$(bondscli keys show miguel -a)" 200000000stake,1000000res,1000000rez
yes $PASSWORD | bondsd add-genesis-account "$(bondscli keys show francesco -a)" 100000000stake,1000000res,1000000rez
yes $PASSWORD | bondsd add-genesis-account "$(bondscli keys show shaun -a)" 100000000stake,1000000res,1000000rez

# Set min-gas-prices
FROM="minimum-gas-prices = \"\""
TO="minimum-gas-prices = \"0.025stake\""
sed -i "s/$FROM/$TO/" "$HOME"/.bondsd/config/app.toml

bondscli config chain-id bondschain-1
bondscli config output json
bondscli config indent true
bondscli config trust-node true

yes $PASSWORD | bondsd gentx --name miguel

bondsd collect-gentxs
bondsd validate-genesis

# Uncomment the below to broadcast node RPC endpoint
#FROM="laddr = \"tcp:\/\/127.0.0.1:26657\""
#TO="laddr = \"tcp:\/\/0.0.0.0:26657\""
#sed -i "s/$FROM/$TO/" "$HOME"/.bondsd/config/config.toml

# Uncomment the below to broadcast REST endpoint
# Do not forget to comment the bottom lines !!
# bondsd start --pruning "syncable" &
# bondscli rest-server --chain-id bondschain-1 --laddr="tcp://0.0.0.0:1317" --trust-node && fg

bondsd start --pruning "syncable" &
bondscli rest-server --chain-id bondschain-1 --trust-node && fg
