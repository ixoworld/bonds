#!/usr/bin/env bash

PASSWORD="12345678"
MIGUEL=$(yes $PASSWORD | bondscli keys show miguel -a)
FRANCESCO=$(yes $PASSWORD | bondscli keys show francesco -a)
SHAUN=$(yes $PASSWORD | bondscli keys show shaun -a)
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
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens=res,rez \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0 \
  --max-supply=1000000abc \
  --order-quantity-limits="10abc,5000res,5000rez" \
  --sanity-rate="" \
  --sanity-margin-percentage="" \
  --allow-sells=true \
  --signers="$MIGUEL" \
  --batch-blocks=1
echo "Created bond..."
bondscli query bonds bond abc

echo "Miguel buys 1abc..."
tx_from_m buy 1abc 5000res,10000rez
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco swaps...(1/8)"
tx_from_f swap abc 5000 res rez
echo "Francesco swaps...(2/8)"
tx_from_f swap abc 5000 rez res
echo "Francesco swaps...(3/8)"
tx_from_f swap abc 5000 res rez
echo "Francesco swaps...(4/8)"
tx_from_f swap abc 5000 rez res
echo "Francesco swaps...(5/8)"
tx_from_f swap abc 5000 res rez
echo "Francesco swaps...(6/8)"
tx_from_f swap abc 5000 rez res
echo "Francesco swaps...(7/8)"
tx_from_f swap abc 5000 res rez
echo "Francesco swaps...(8/8)"
tx_from_f swap abc 5000 rez res

echo "Miguel sells 1abc..."
tx_from_m sell 1abc
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"
echo "Miguel made profit."
