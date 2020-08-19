#!/bin/bash

rm -rf ~/.bonds*

# Initialize the genesis.json file that will help you to bootstrap the network
bondsd init MyValidator --chain-id=bonds-localnet

bondscli config chain-id bonds-localnet
bondscli config output json
bondscli config indent true
bondscli config trust-node true
bondscli config keyring-backend test

# Change default bond token genesis.json
sed -i 's/"leveldb"/"goleveldb"/g' ~/.bondsd/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bondsd/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bondsd/config/config.toml

# Create a key to hold your validator account
bondscli keys add jack
bondscli keys add alice

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some coins
# with the genesis.app_state.staking.params.bond_denom denom, the default is staking
bondsd add-genesis-account jack 100000000000ubonds --keyring-backend test
bondsd add-genesis-account alice 100000000000ubonds --keyring-backend test

# Generate the transaction that creates your validator
bondsd gentx --name jack --amount=10000000ubonds --keyring-backend test

# Add the generated bonding transaction to the genesis file
bondsd collect-gentxs
bondsd validate-genesis

# Now its safe to start `bondsd`
bondsd start
