#!/usr/bin/env bash

ALICE=$(bondscli keys show alice --keyring-backend=test -a)
BOB=$(bondscli keys show bob --keyring-backend=test -a)
FEE=$(bondscli keys show fee --keyring-backend=test -a)

# Visual representation of bonding curve:
# https://www.desmos.com/calculator/2ael5dojku

# Create a Power Function bonding curve
bondscli tx bonds create-bond \
  --token=mytoken \
  --name="My Token" \
  --description="My first continuous token" \
  --function-type=power_function \
  --function-parameters="m:12,n:2,c:100" \
  --reserve-tokens=uatom \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0 \
  --fee-address="$FEE" \
  --max-supply=1000000mytoken \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells \
  --signers="$ALICE" \
  --batch-blocks=3 \
  --from alice --keyring-backend=test --broadcast-mode block -y
# Query the created bond
bondscli q bonds bond mytoken
# We can keep an eye on the batch
watch -n 1 bondscli q bonds batch mytoken

# Query the price of buying 10mytoken
bondscli q bonds buy-price 10mytoken
# Query the token price at supply=10mytoken
bondscli q bonds price 10mytoken

# Buy 10mytoken from alice with max spend of 1000000uatom
bondscli tx bonds buy 10mytoken 1000000uatom --from alice --keyring-backend=test --broadcast-mode block -y
# Wait for order to get processed
sleep 21
# Query alice's account
bondscli q account "$ALICE"

# Query the price of buying 10mytoken
bondscli q bonds buy-price 10mytoken
# Query the token price at supply=20mytoken
bondscli q bonds price 20mytoken

# Buy 10mytoken from bob with max spend of 1000000uatom
bondscli tx bonds buy 10mytoken 1000000uatom --from bob --keyring-backend=test --broadcast-mode block -y
# Wait for order to get processed
sleep 21
# Query bob's account
bondscli q account "$BOB"

# Sell 10mytoken from alice at a profit :]
bondscli tx bonds sell 10mytoken --from alice --keyring-backend=test --broadcast-mode block -y
# Wait for order to get processed
sleep 21
# Query alice's account
bondscli q account "$ALICE"

# Sell 10mytoken from bob at a loss :[
bondscli tx bonds sell 10mytoken --from bob --keyring-backend=test --broadcast-mode block -y
# Wait for order to get processed
sleep 21
# Query bob's account
bondscli q account "$BOB"
