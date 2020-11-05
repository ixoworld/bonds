#!/usr/bin/env bash

bondsd init local --chain-id bondschain-1

bondscli keys delete node --keyring-backend=test --force
bondscli keys delete alice --keyring-backend=test --force
bondscli keys delete bob --keyring-backend=test --force
bondscli keys delete charlie --keyring-backend=test --force
bondscli keys delete fee --keyring-backend=test --force

bondscli keys add node --keyring-backend=test
bondscli keys add alice --keyring-backend=test
bondscli keys add bob --keyring-backend=test
bondscli keys add charlie --keyring-backend=test
bondscli keys add fee --keyring-backend=test

# Note: important to add 'node' as a genesis-account since this is the validator
bondsd add-genesis-account $(bondscli keys show node --keyring-backend=test -a) 1000000000000uatom,1000000000000ubtc
bondsd add-genesis-account $(bondscli keys show alice --keyring-backend=test -a) 1000000000000uatom,1000000000000ubtc
bondsd add-genesis-account $(bondscli keys show bob --keyring-backend=test -a) 1000000000000uatom,1000000000000ubtc
bondsd add-genesis-account $(bondscli keys show charlie --keyring-backend=test -a) 1000000000000uatom,1000000000000ubtc

# Set staking token (both bond_denom and mint_denom)
STAKING_TOKEN="uatom"
FROM="\"bond_denom\": \"stake\""
TO="\"bond_denom\": \"$STAKING_TOKEN\""
sed -i "s/$FROM/$TO/" "$HOME"/.bondsd/config/genesis.json
FROM="\"mint_denom\": \"stake\""
TO="\"mint_denom\": \"$STAKING_TOKEN\""
sed -i "s/$FROM/$TO/" "$HOME"/.bondsd/config/genesis.json

bondscli config chain-id bondschain-1
bondscli config output json
bondscli config indent true
bondscli config trust-node true
bondscli config keyring-backend test

bondsd gentx --name node --amount 1000000uatom --keyring-backend=test

bondsd collect-gentxs
bondsd validate-genesis

bondsd start --pruning "everything" &
bondscli rest-server --chain-id bondschain-1 --trust-node && fg
