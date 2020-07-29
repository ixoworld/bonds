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
  yes $PASSWORD | bondscli tx bonds "$cmd" --from miguel -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
}

tx_from_f() {
  cmd=$1
  shift
  yes $PASSWORD | bondscli tx bonds "$cmd" --from francesco -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
}

RET=$(bondscli status 2>&1)
if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
  wait
fi

PASSWORD="12345678"
GAS_PRICES="0.025stake"
MIGUEL=$(yes $PASSWORD | bondscli keys show miguel -a)
FRANCESCO=$(yes $PASSWORD | bondscli keys show francesco -a)
FEE=$(yes $PASSWORD | bondscli keys show fee -a)

# d0 := 500.0   // initial raise (reserve)
# p0 := 0.01    // initial price (reserve per token)
# theta := 0.4  // initial allocation (percentage)
# kappa := 3.0  // degrees of polynomial (i.e. x^2, x^4, x^6)

# R0 = 300              // initial reserve (1-theta)*d0
# S0 = 50000            // initial supply
# V0 = 416666666666.667 // invariant

echo "Creating bond..."
tx_from_m create-bond \
  --token=abc \
  --name="A B C" \
  --description="Description about A B C" \
  --function-type=augmented_function \
  --function-parameters="d0:500.0,p0:0.01,theta:0.4,kappa:3.0" \
  --reserve-tokens=res \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0 \
  --fee-address="$FEE" \
  --max-supply=1000000abc \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells=true \
  --signers="$MIGUEL" \
  --batch-blocks=1
echo "Created bond..."
bondscli query bonds bond abc

echo "Miguel buys 20000abc..."
tx_from_m buy 20000abc 100000res
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco buys 20000abc..."
tx_from_f buy 20000abc 100000res
echo "Francesco's account..."
bondscli query auth account "$FRANCESCO"

echo "Miguel cannot buy 10001abc..."
tx_from_m buy 10001abc 100000res
echo "Miguel cannot sell anything..."
tx_from_m sell 10000abc

echo "Miguel can buy 10000abc..."
tx_from_m buy 10000abc 100000res
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Bond state is now open..."
bondscli query bonds bond abc

echo "Miguel sells 30000abc..."
tx_from_m sell 30000abc
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco sells 30000abc..."
tx_from_f sell 30000abc
echo "Francesco's account..."
bondscli query auth account "$FRANCESCO"
