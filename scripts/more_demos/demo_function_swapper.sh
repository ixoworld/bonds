#!/usr/bin/env bash

wait() {
  echo "Waiting for chain to start..."
  while :; do
    RET=$(bondscli status 2>&1)
    if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
      sleep 1
    else
      echo "A few more seconds..."
      sleep 6
      break
    fi
  done
}

tx_from_m() {
  cmd=$1
  shift
  yes $PASSWORD | bondscli tx bonds "$cmd" --from miguel --keyring-backend=test -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
}

tx_from_f() {
  cmd=$1
  shift
  yes $PASSWORD | bondscli tx bonds "$cmd" --from francesco --keyring-backend=test -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
}

RET=$(bondscli status 2>&1)
if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
  wait
fi

PASSWORD="12345678"
GAS_PRICES="0.025stake"
MIGUEL=$(yes $PASSWORD | bondscli keys show miguel --keyring-backend=test -a)
FRANCESCO=$(yes $PASSWORD | bondscli keys show francesco --keyring-backend=test -a)
SHAUN=$(yes $PASSWORD | bondscli keys show shaun --keyring-backend=test -a)
FEE=$(yes $PASSWORD | bondscli keys show fee --keyring-backend=test -a)

echo "Creating bond..."
tx_from_m create-bond \
  --token=abc \
  --name="A B C" \
  --description="Description about A B C" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens=res,rez \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEE" \
  --max-supply=1000000abc \
  --order-quantity-limits="10abc,5000res,5000rez" \
  --sanity-rate="0.5" \
  --sanity-margin-percentage="20" \
  --allow-sells \
  --signers="$MIGUEL" \
  --batch-blocks=1
echo "Created bond..."
bondscli q bonds bond abc

echo "Miguel buys 1abc..."
tx_from_m buy 1abc 500res,1000rez
echo "Miguel's account..."
bondscli q auth account "$MIGUEL"

echo "Francesco buys 10abc..."
tx_from_f buy 10abc 10100res,10100rez
echo "Francesco's account..."
bondscli q auth account "$FRANCESCO"

echo "Miguel swap 500 res to rez..."
tx_from_m swap abc 500 res rez
echo "Miguel's account..."
bondscli q auth account "$MIGUEL"

echo "Francesco swap 500 rez to res..."
tx_from_f swap abc 500 rez res
echo "Francesco's account..."
bondscli q auth account "$FRANCESCO"

echo "Miguel swaps above order limit (tx will fail)..."
tx_from_m swap abc 5001 res rez
echo "Miguel's account (no  changes)..."
bondscli q auth account "$MIGUEL"

echo "Francesco swaps to violate sanity (tx will be successful but order will fail)..."
tx_from_f swap abc 5000 rez res
echo "Francesco's account (no changes)..."
bondscli q auth account "$FRANCESCO"

echo "Miguel sells 1abc..."
tx_from_m sell 1abc
echo "Miguel's account..."
bondscli q auth account "$MIGUEL"

echo "Francesco sells 10abc..."
tx_from_f sell 10abc
echo "Francesco's account..."
bondscli q auth account "$FRANCESCO"
