#!/usr/bin/env bash

PASSWORD="12345678"
MIGUEL=$(yes $PASSWORD | bondscli keys show miguel -a)
FRANCESCO=$(yes $PASSWORD | bondscli keys show francesco -a)
SHAUN=$(yes $PASSWORD | bondscli keys show shaun -a)
RESERVE=$(yes $PASSWORD | bondscli keys show reserve -a)
FEE=$(yes $PASSWORD | bondscli keys show fee -a)

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
  yes $PASSWORD | bondscli tx bonds "$cmd" --from miguel -y --broadcast-mode block "$@"
}

tx_from_f() {
  cmd=$1
  shift
  yes $PASSWORD | bondscli tx bonds "$cmd" --from francesco -y --broadcast-mode block "$@"
}

RET=$(bondscli status 2>&1)
if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
  wait
fi

echo "Creating bond..."
tx_from_m create-bond \
  --token=abc \
  --name="A B C" \
  --description="Description about A B C" \
  --function-type=sigmoid_function \
  --function-parameters="a:3,b:5,c:1" \
  --reserve-tokens=res \
  --reserve-address="$RESERVE" \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEE" \
  --max-supply=1000000abc \
  --order-quantity-limits="" \
  --sanity-rate="" \
  --sanity-margin-percentage="" \
  --allow-sells=true \
  --signers="$MIGUEL" \
  --batch-blocks=1
echo "Created bond..."
bondscli query bonds bond abc

echo "Miguel buys 50abc..."
tx_from_m buy 50abc 1000000res
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco buys 50abc..."
tx_from_f buy 50abc 1000000res
echo "Francesco's account..."
bondscli query auth account "$FRANCESCO"

echo "Miguel sells 50abc..."
tx_from_m sell 50abc
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco sells 50abc..."
tx_from_f sell 50abc
echo "Francesco's account..."
bondscli query auth account "$FRANCESCO"
