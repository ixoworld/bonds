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
SHAUN=$(yes $PASSWORD | bondscli keys show shaun -a)
FEE=$(yes $PASSWORD | bondscli keys show fee -a)

echo "Creating bond..."
tx_from_m create-bond \
  --token=abc \
  --name="A B C" \
  --description="Description about A B C" \
  --function-type=power_function \
  --function-parameters="m:12,n:2,c:100" \
  --reserve-tokens=res \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
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

echo "Editing name..."
tx_from_m edit-bond \
  --token=abc \
  --name="New A B C" \
  --signers="$MIGUEL"
echo "Edited name..."
bondscli query bonds bond abc

echo "Editing description..."
tx_from_m edit-bond \
  --token=abc \
  --description="New description about A B C" \
  --signers="$MIGUEL"
echo "Edited description..."
bondscli query bonds bond abc

echo "Editing order quantity limits..."
tx_from_m edit-bond \
  --token=abc \
  --order-quantity-limits=100abc,200res1,300res2 \
  --signers="$MIGUEL"
echo "Edited description..."
bondscli query bonds bond abc

echo "Editing sanity rate and margin..."
tx_from_m edit-bond \
  --token=abc \
  --sanity-rate=100000 \
  --sanity-margin-percentage=10 \
  --signers="$MIGUEL"
echo "Edited description..."
bondscli query bonds bond abc

echo "Editing nothing..."
tx_from_m edit-bond \
  --token=abc \
  --signers="$MIGUEL"
echo "...this caused an error"
