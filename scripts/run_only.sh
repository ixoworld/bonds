#!/usr/bin/env bash

# Uncomment the below to broadcast REST endpoint
# Do not forget to comment the bottom lines !!
# bondsd start --pruning "everything" &
# bondscli rest-server --chain-id bondschain-1 --laddr="tcp://0.0.0.0:1317" --trust-node && fg

bondsd start --pruning "everything" &
bondscli rest-server --chain-id bondschain-1 --trust-node && fg
