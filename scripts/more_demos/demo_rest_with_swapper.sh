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

rest_from_m() {
  curl -s -X POST localhost:1317/"$1" --data-binary "$2" >tx.json                   # generate
  yes $PASSWORD | bondscli tx sign tx.json --from=miguel --output-document=tx.json  # sign
  python3 demo_rest_fix_tx_format.py                                                # reformat
  curl -s -X POST localhost:1317/txs --data-binary "$(<tx.json)"                    # broadcast
  rm tx.json
}

rest_from_f() {
  curl -s -X POST localhost:1317/"$1" --data-binary "$2" >tx.json                      # generate
  yes $PASSWORD | bondscli tx sign tx.json --from=francesco --output-document=tx.json  # sign
  python3 demo_rest_fix_tx_format.py                                                   # reformat
  curl -s -X POST localhost:1317/txs --data-binary "$(<tx.json)"                       # broadcast
  rm tx.json
}

query_bond() {
  curl -X GET localhost:1317/bonds/"$1"
}

query_account() {
  curl -X GET localhost:1317/auth/accounts/"$1"
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
# shellcheck disable=SC2046
rest_from_m bonds/create_bond '{
                                "base_req":{
                                  "from":"'$MIGUEL'",
                                  "chain_id":"bondschain-1",
                                  "gas_prices":[{"denom":"stake","amount":"0.025"}]
                                },
                                "token":"abc",
                                "name":"A B C",
                                "description":"Description about A B C",
                                "function_type":"swapper_function",
                                "function_parameters":"",
                                "reserve_tokens":"res,rez",
                                "tx_fee_percentage":"0.5",
                                "exit_fee_percentage":"0.1",
                                "fee_address":"'$FEE'",
                                "max_supply":"1000000abc",
                                "order_quantity_limits":"",
                                "sanity_rate":"0",
                                "sanity_margin_percentage":"0",
                                "allow_sells":"true",
                                "signers":"'$MIGUEL'",
                                "batch_blocks":"1"
                              }'
echo "Created bond..."
query_bond abc

echo "Editing bond..."
# shellcheck disable=SC2046
rest_from_m bonds/edit_bond '{
                              "base_req":{
                                "from":"'$MIGUEL'",
                                "chain_id":"bondschain-1",
                                "gas_prices":[{"denom":"stake","amount":"0.025"}]
                              },
                              "token":"abc",
                              "name":"New A B C",
                              "description":"New description about A B C",
                              "sanity_rate":"[do-not-modify]",
                              "sanity_margin_percentage":"[do-not-modify]",
                              "order_quantity_limits":"[do-not-modify]",
                              "signers":"'$MIGUEL'"
                            }'
echo "Edited bond..."
query_bond abc

echo "Miguel buys 1abc..."
# shellcheck disable=SC2046
rest_from_m bonds/buy '{
                        "base_req":{
                          "from":"'$MIGUEL'",
                          "chain_id":"bondschain-1",
                          "gas_prices":[{"denom":"stake","amount":"0.025"}]
                        },
                        "bond_token":"abc",
                        "bond_amount":"1",
                        "max_prices":"500res,1000rez"
                      }'
echo "Miguel's account..."
query_account "$MIGUEL"

echo "Francesco buys 10abc..."
# shellcheck disable=SC2046
rest_from_f bonds/buy '{
                        "base_req":{
                          "from":"'$FRANCESCO'",
                          "chain_id":"bondschain-1",
                          "gas_prices":[{"denom":"stake","amount":"0.025"}]
                        },
                        "bond_token":"abc",
                        "bond_amount":"10",
                        "max_prices":"10100res,10100rez"
                      }'
echo "Francesco's account..."
query_account "$FRANCESCO"

echo "Miguel swap 5000 res to rez..."
# shellcheck disable=SC2046
rest_from_m bonds/swap '{
                        "base_req":{
                          "from":"'$MIGUEL'",
                          "chain_id":"bondschain-1",
                          "gas_prices":[{"denom":"stake","amount":"0.025"}]
                        },
                        "bond_token":"abc",
                        "from_amount":"5000",
                        "from_token":"res",
                        "to_token":"rez"
                      }'
echo "Miguel's account..."
bondscli q auth account "$MIGUEL"

echo "Francesco swap 5000 rez to res..."
# shellcheck disable=SC2046
rest_from_f bonds/swap '{
                          "base_req":{
                            "from":"'$FRANCESCO'",
                            "chain_id":"bondschain-1",
                            "gas_prices":[{"denom":"stake","amount":"0.025"}]
                          },
                          "bond_token":"abc",
                          "from_amount":"5000",
                          "from_token":"rez",
                          "to_token":"res"
                        }'
echo "Francesco's account..."
bondscli q auth account "$FRANCESCO"

echo "Miguel sells 1abc..."
# shellcheck disable=SC2046
rest_from_m bonds/sell '{
                          "base_req":{
                            "from":"'$MIGUEL'",
                            "chain_id":"bondschain-1",
                            "gas_prices":[{"denom":"stake","amount":"0.025"}]
                          },
                          "bond_token":"abc",
                          "bond_amount":"1"
                        }'
echo "Miguel's account..."
query_account "$MIGUEL"

echo "Francesco sells 10abc..."
# shellcheck disable=SC2046
rest_from_f bonds/sell '{
                          "base_req":{
                            "from":"'$FRANCESCO'",
                            "chain_id":"bondschain-1",
                            "gas_prices":[{"denom":"stake","amount":"0.025"}]
                          },
                          "bond_token":"abc",
                          "bond_amount":"10"
                        }'
echo "Francesco's account..."
query_account "$FRANCESCO"
