#!/usr/bin/env bash

PASSWORD="12345678"

bondsd init local --chain-id bondschain-1

yes $PASSWORD | bondscli keys delete miguel --force
yes $PASSWORD | bondscli keys delete francesco --force
yes $PASSWORD | bondscli keys delete shaun --force
yes $PASSWORD | bondscli keys delete reserve --force
yes $PASSWORD | bondscli keys delete reserve2 --force
yes $PASSWORD | bondscli keys delete reserve3 --force
yes $PASSWORD | bondscli keys delete reserve4 --force
yes $PASSWORD | bondscli keys delete reserve5 --force
yes $PASSWORD | bondscli keys delete fee --force
yes $PASSWORD | bondscli keys delete fee2 --force
yes $PASSWORD | bondscli keys delete fee3 --force
yes $PASSWORD | bondscli keys delete fee4 --force
yes $PASSWORD | bondscli keys delete fee5 --force

yes $PASSWORD | bondscli keys add miguel
yes $PASSWORD | bondscli keys add francesco
yes $PASSWORD | bondscli keys add shaun
yes $PASSWORD | bondscli keys add reserve
yes $PASSWORD | bondscli keys add reserve2
yes $PASSWORD | bondscli keys add reserve3
yes $PASSWORD | bondscli keys add reserve4
yes $PASSWORD | bondscli keys add reserve5
yes $PASSWORD | bondscli keys add fee
yes $PASSWORD | bondscli keys add fee2
yes $PASSWORD | bondscli keys add fee3
yes $PASSWORD | bondscli keys add fee4
yes $PASSWORD | bondscli keys add fee5

yes $PASSWORD | bondsd add-genesis-account "$(bondscli keys show miguel -a)" 100000000stake,1000000res,1000000rez
yes $PASSWORD | bondsd add-genesis-account "$(bondscli keys show francesco -a)" 100000000stake,1000000res,1000000rez
yes $PASSWORD | bondsd add-genesis-account "$(bondscli keys show shaun -a)" 100000000stake,1000000res,1000000rez

bondscli config chain-id bondschain-1
bondscli config output json
bondscli config indent true
bondscli config trust-node true

yes $PASSWORD | bondsd gentx --name miguel

bondsd collect-gentxs
bondsd validate-genesis

bondsd start --pruning "syncable" &
bondscli rest-server --chain-id bondschain-1 --trust-node && fg
