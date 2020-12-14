---
description: Tutorial
---

# Swapper Function AMM

## Contents

* [Bond Configuration](02_swapper.md#bond-configuration)
* [Bond Creation](02_swapper.md#bond-creation)
* [Supply Liquidity](02_swapper.md#supply-liquidity)
* [Make a Swap](02_swapper.md#make-a-swap)

## Bond Configuration

Configuration steps that were covered in previous tutorials will not be described in this tutorial, unless they take on new meaning.

### Curve Function

In this tutorial, a swapper function bond will be created. The swapper function implemented by the Bonds module is shown below, where `x` and `y` represent the balances of two distinct reserves \(i.e. two reserve tokens\). The constant `k` is not user-decided but is simply the product of the two balances, which is expected to remain the same for any number of swaps and only increases/decreases when liquidity increases/decreases \(i.e. a buy or sell, respectively\).

![swapper function](../.gitbook/assets/swapper.png)

This function allows us to calculate the price of buying a number of bond tokens or the returns when selling bond tokens, as well as the returns for swapping a number of `x` tokens to `y` tokens and vice-versa:

* Buying and selling is considered adding/removing liquidity to/from the swap AMM, and thus buying requires depositing both `x` and `y` tokens, and similarly returns when selling are of both `x` and `y` tokens. Internal calculations that determine prices/returns will not be discussed in this tutorial.
* Swap amounts are determined by using the above function directly. If someone increases the value of `x` \(by depositing `x` tokens\), the value of `y` has to change such that the product remains `k`. This change \(decrease\) in `y` is precisely what the user gets in return. The same applies the other way round.

Note that in the case of the swapper function, we do not have any extra constants \(referred to as function parameters\) that we need to figure out a value for.

### Reserve Token/s

The reserve tokens are the tokens that accounts will send to the bond in order to mint bond tokens. A swapper function has a pair of reserve tokens, which are in fact the tokens that can be swapped for each other. In this tutorial, we will use `res` and `rez`.

Note that in the case of the swapper function, the amount of reserve required to add liquidity to the AMM \(i.e. buy a certain number of `demo` tokens\) will depend on the ratio of the current reserves. If the `x` balance is greater, than more `x` tokens will be required, and vice-versa.

### Fees

In the case of a swapper function, transaction fees also apply to swaps, and not just buys and sells.

### Order Quantity Limits

In the case of a swapper function, order quantity limits also apply to swap amounts. In this tutorial, we will set an order limit of `100demo`, `5000res`, and `6000rez`. We could have decided to not put any limits at all.

### Sanity Rate and Sanity Margin Percentage

The sanity values \(sanity rate and sanity margin percentage\) are used in the case of a swapper function to set a range of valid exchange rate \(`x/y`\) between the two reserve tokens, such that if a swap order causes the exchange rate to go outside of the valid range, the swap is cancelled.

The valid exchange rate range is defined by `sanity rate ± sanity margin percentage`. In other words, between `(100 - sanity margin percentage) x sanity rate` and `(100 + sanity margin percentage) x sanity rate`.

In this tutorial, we will go with a `0.5` sanity rate and `20%` sanity margin percentage. This means that the reserve balance of `x` is expected to be half that of `y`, with a 20 percent error. If `x=500`, then `y` can be between `833.33` and `1250`.

### Other Customisation

Other customisation options that this tutorial will not go into is the ability to disable sells \(burns\), the ability to have multiple signers as the creators/editors of the bond, and the ability to add an outcome payment. In this tutorial, sells will be enabled, the signer will be set to the address underlying the `shaun` account \(created when running `make run_with_data`\), and there will be no outcome payment \(discussed in the augmented function tutorial\).

## Bond Creation

The bond, with the above configurations, can be created as follows. Note that sanity rate and sanity margin percentage only apply to swapper functions and were thus just set to `0`.

```bash
SHAUNADDR="$(bondscli keys show shaun -a)"
FEEADDR="$(bondscli keys show fee -a)"

bondscli tx bonds create-bond \
  --token=demo \
  --name="My Bond" \
  --description="Description about my bond" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens=res,rez \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEEADDR" \
  --max-supply=1000000demo \
  --order-quantity-limits="10abc,5000res,6000rez" \
  --sanity-rate="0.5" \
  --sanity-margin-percentage="20" \
  --allow-sells \
  --signers="$SHAUNADDR" \
  --batch-blocks=2 \
  --outcome-payment="" \
  --from shaun \
  --keyring-backend=test \
  --broadcast-mode block \
  --gas-prices=0.025stake \
  -y
```

The created bond can be queried using `bondscli q bonds bond demo`, which should return the following, but with different addresses:

```bash
{
  "type": "bonds/Bond",
  "value": {
    "token": "demo",
    "name": "My Bond",
    "description": "Description about my bond",
    "creator": "cosmos1km5cj5yq5c4757ksmpe88sx4snyfgd6wx8nfzx",
    "function_type": "swapper_function",
    "function_parameters": null,
    "reserve_tokens": [
      "res",
      "rez"
    ],
    "tx_fee_percentage": "0.500000000000000000",
    "exit_fee_percentage": "0.100000000000000000",
    "fee_address": "cosmos19d7xn0e6sr9l2yq8pls37vwu4ursj974k3e4sf",
    "max_supply": {
      "denom": "demo",
      "amount": "1000000"
    },
    "order_quantity_limits": [
      {
        "denom": "abc",
        "amount": "10"
      },
      {
        "denom": "res",
        "amount": "5000"
      },
      {
        "denom": "rez",
        "amount": "5000"
      }
    ],
    "sanity_rate": "0.500000000000000000",
    "sanity_margin_percentage": "20.000000000000000000",
    "current_supply": {
      "denom": "demo",
      "amount": "0"
    },
    "current_reserve": [],
    "allow_sells": true,
    "signers": [
      "cosmos1km5cj5yq5c4757ksmpe88sx4snyfgd6wx8nfzx"
    ],
    "batch_blocks": "2",
    "outcome_payment": [],
    "state": "OPEN"
  }
}
```

Note that some extra fields that we did not input ourselves are present. These were discussed in previous tutorials.

## Supply Liquidity

Liquidity can be added to the AMM by performing buys \(mint-to-deposit\), which mints bond tokens. Similarly, liquidity would be removed by performing sells \(burn-to-withdraw\) which burns bond tokens. Once liquidity \(i.e. reserve tokens\) is added by liquidity providers, swaps can take place using the reserve that is in place. If a liquidity provider burns-to-withdraw, they get a proportional share of each of the two reserve pools \(`x` and `y`\) in exchange for bond tokens.

In general, when adding liquidity to a swapper function, the current exchange rate \(based on the `x` and `y` balances\) is used to determine how much of each reserve token makes up the price. The first buy is special and plays a very important role in specifying the price of the bond token. Since we have no price reference for the first buy in a swapper function, the `MaxPrices` specified are used as the actual price, with no extra fees charged.

This effectively means that if the user requested `n` bond tokens with max prices `x1` and `y2` \(for reserve tokens `x` and `y`\), the next buyers will have to pay `(x1/n)` and `(y1/n)` tokens per bond token requested. Specifying high `x1` and `y1` prices for a small `n` \(say `n=1`\) means that the next buyers will have to pay at most `x1` and `y1` per bond token. **In summary, it is important that the first buy is well-calculated and performed carefully.**

In this tutorial, we will perform a buy of `1demo` with max prices `500res` and `1000rez`.

```bash
bondscli tx bonds buy 1demo 500res,1000rez \
  --from miguel \
  --keyring-backend=test \
  --broadcast-mode block \
  --gas-prices=0.025stake \
  -y
```

We can query the `miguel` account to confirm that the `1demo` has reached the account by using `bondscli q account $(bondscli keys show miguel -a)`, where we can also see that exactly `500res` and `1000rez` were taken out of the account.

```bash
...
"coins": [
  {
    "denom": "demo",
    "amount": "1"
  },
  {
    "denom": "res",
    "amount": "999500"
  },
  {
    "denom": "rez",
    "amount": "999000"
  },
...
```

At this point, if we query the buy price for an additional `1demo` using `bondscli q bonds buy-price 1demo`, as expected, this shows a `500res` and `1000rez` price, excluding fees:

```text
...
"prices": [
  {
    "denom": "res",
    "amount": "500"
  },
  {
    "denom": "rez",
    "amount": "1000"
  }
],
...
```

## Make a Swap

Before performing a swap, we can query the current returns for swapping the tokens. Say we want to swap `10res` to `rez`, we can perform the query `bondscli q bonds swap-return demo 10res rez`, which gives:

```bash
{
  "total_returns": [
    {
      "denom": "rez",
      "amount": "17"
    }
  ],
  "total_fees": [
    {
      "denom": "res",
      "amount": "1"
    }
  ]
}
```

Since the current reserve balances are `500res` and `1000rez`, then the value of constant `k` is expected to be `500000`. If the new balance of `res` becomes `509` \(`1res` fee deducted\), then the new value of `rez` will be `500000/509 = 982.32`. The decrease of `17.68` rounded down is what the user gets in return `=17rez`.

We can perform a swap of `10res` as shown below. The account used is the `miguel` account \(created when running `make run_with_data`\).

```bash
bondscli tx bonds swap demo 10 res rez \
  --from miguel \
  --keyring-backend=test \
  --broadcast-mode block \
  --gas-prices=0.025stake \
  -y
```

We can query the `miguel` account to confirm that the demo tokens are no longer in the account by using `bondscli q account $(bondscli keys show miguel -a)`. A maximum of 2 blocks-worth of time might need to pass for the order in the batch to get processed.

```bash
...
"coins": [
  ...
  {
    "denom": "res",
    "amount": "999490"
  },
  {
    "denom": "rez",
    "amount": "999017"
  },
  ...
],
...
```

Note how the account now has `999490res` and `999017rez` \(a `10res` decrease and `17rez` increase\). Note also that we have not violated the sanity values \(sanity rate and sanity margin percentage\) with this swap.

