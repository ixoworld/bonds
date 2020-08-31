#!/usr/bin/env bash

PASSWORD="12345678"
GAS_PRICES="0.025stake"

bondsd init local --chain-id bondschain-1

bondscli keys add miguel --ledger
yes $PASSWORD | bondscli keys add francesco
yes $PASSWORD | bondscli keys add shaun
yes $PASSWORD | bondscli keys add reserve
yes $PASSWORD | bondscli keys add fee

bondsd add-genesis-account "$(bondscli keys show miguel -a)" 100000000stake,1000000res,1000000rez
bondsd add-genesis-account "$(bondscli keys show francesco -a)" 100000000stake,1000000res,1000000rez
bondsd add-genesis-account "$(bondscli keys show shaun -a)" 100000000stake,1000000res,1000000rez

bondscli config chain-id bondschain-1
bondscli config output json
bondscli config indent true
bondscli config trust-node true

echo "$PASSWORD" | bondsd gentx --name miguel

bondsd collect-gentxs
bondsd validate-genesis

bondsd start --pruning "everything" &
bondscli rest-server --chain-id bondschain-1 --trust-node && fg
