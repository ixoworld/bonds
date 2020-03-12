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
  yes $PASSWORD | bondscli tx bonds "$cmd" --from miguel -y --broadcast-mode block "$@"
}

tx_from_f() {
  cmd=$1
  shift
  yes $PASSWORD | bondscli tx bonds "$cmd" --from francesco -y --broadcast-mode block "$@"
}

create_bond_multisig() {
  bondscli tx bonds create-bond \
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
    --sanity-rate="" \
    --sanity-margin-percentage="" \
    --allow-sells=true \
    --signers="$(bondscli keys show francesco -a),$(bondscli keys show shaun -a)" \
    --batch-blocks=1 \
    --from="$MIGUEL" -y --broadcast-mode block --generate-only >multisig.json
  yes $PASSWORD | bondscli tx sign multisig.json --from=francesco --output-document=multisig.json
  yes $PASSWORD | bondscli tx sign multisig.json --from=shaun --output-document=multisig.json
  bondscli tx broadcast multisig.json
  rm multisig.json
}

RET=$(bondscli status 2>&1)
if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
  wait
fi

PASSWORD="12345678"
MIGUEL=$(yes $PASSWORD | bondscli keys show miguel -a)
FRANCESCO=$(yes $PASSWORD | bondscli keys show francesco -a)
SHAUN=$(yes $PASSWORD | bondscli keys show shaun -a)
FEE=$(yes $PASSWORD | bondscli keys show fee -a)

echo "Creating bond..."
create_bond_multisig
echo "Waiting a bit..."
sleep 5
echo "Created bond..."
bondscli query bonds bond abc

echo "Miguel buys 10abc..."
tx_from_m buy 10abc 1000000res
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco buys 10abc..."
tx_from_f buy 10abc 1000000res
echo "Francesco's account..."
bondscli query auth account "$FRANCESCO"

echo "Miguel sells 10abc..."
tx_from_m sell 10abc
echo "Miguel's account..."
bondscli query auth account "$MIGUEL"

echo "Francesco sells 10abc..."
tx_from_f sell 10abc
echo "Francesco's account..."
bondscli query auth account "$FRANCESCO"
