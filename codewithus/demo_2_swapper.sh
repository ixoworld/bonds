#!/usr/bin/env bash

ALICE=$(bondscli keys show alice --keyring-backend=test -a)
BOB=$(bondscli keys show bob --keyring-backend=test -a)
FEE=$(bondscli keys show fee --keyring-backend=test -a)

# Visual representation of liquidity pool:
# https://www.desmos.com/calculator/plmzg2hboo

# Create liquidity pool
bondscli tx bonds create-bond \
  --token=atombtcpool \
  --name="ATOM-BTC Pool" \
  --description="An ATOM-BTC Liquidity Pool" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens=uatom,ubtc \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEE" \
  --max-supply=1000000atombtcpool \
  --order-quantity-limits="30000000000uatom,100000000ubtc" \
  --sanity-rate="3000" \
  --sanity-margin-percentage="10" \
  --allow-sells \
  --signers="$ALICE" \
  --batch-blocks=1 \
  --from alice --keyring-backend=test --broadcast-mode block -y
# Query the created bond
bondscli q bonds bond atombtcpool

# Alice adds liquidity (+10atombtcpool -3000ATOM -1BTC)
bondscli tx bonds buy 10atombtcpool 3000000000uatom,1000000ubtc --from alice --keyring-backend=test --broadcast-mode block -y
# Query alice's account
bondscli q account "$ALICE"

# Bob adds liquidity (+100atombtcpool with max prices 310,000ATOM and 110BTC)
bondscli tx bonds buy 1000atombtcpool 310000000000uatom,110000000ubtc --from bob --keyring-backend=test --broadcast-mode block -y
# Query bob's account
bondscli q account "$BOB"

# Query the swap returns from 3000ATOM to BTC (expected to be ~1BTC before fees)
bondscli q bonds swap-return atombtcpool 3000000000uatom ubtc
# Query the swap returns from 1BTC to ATOM (expected to be ~3000ATOM before fees)
bondscli q bonds swap-return atombtcpool 1000000ubtc uatom

# Alice swaps 3000ATOM to BTC at the current rate
bondscli tx bonds swap atombtcpool 3000000000 uatom ubtc --from alice --keyring-backend=test --broadcast-mode block -y
# Query alice's account
bondscli q account "$ALICE"

# Bob swaps 1BTC to ATOM at the current rate
bondscli tx bonds swap atombtcpool 1000000 ubtc uatom --from bob --keyring-backend=test --broadcast-mode block -y
# Query bob's account
bondscli q account "$BOB"

# Alice tries to swap above the order limit
bondscli tx bonds swap atombtcpool 30000000001 uatom ubtc --from alice --keyring-backend=test --broadcast-mode block -y
# Query alice's account [no changes]
bondscli q account "$ALICE"

# Bob tries to swap but violates sanity rates (tx successful but order will fail)
bondscli tx bonds swap atombtcpool 100000000 ubtc uatom --from bob --keyring-backend=test --broadcast-mode block -y
# Query bob's account [no changes]
bondscli q account "$BOB"

# Alice removes liquidity (-10atombtcpool +uatom +ubtc)
bondscli tx bonds sell 10atombtcpool --from alice --keyring-backend=test --broadcast-mode block -y
# Query alice's account
bondscli q account "$ALICE"

# Bob removes liquidity (-1000atombtcpool +uatom +ubtc)
bondscli tx bonds sell 1000atombtcpool --from bob --keyring-backend=test --broadcast-mode block -y
# Query bob's account
bondscli q account "$BOB"
