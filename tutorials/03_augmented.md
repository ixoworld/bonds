# Creating an augmented bonding curve

Throughout this tutorial, some knowledge around [Augmented Bonding Curves](https://medium.com/giveth/deep-dive-augmented-bonding-curves-3f1f7c1fa751) will be assumed.

## Contents

* [Bond Configuration](03_augmented.md#bond-configuration)
* [Bond Creation](03_augmented.md#bond-creation)
* [Mint to Deposit \(Hatch Phase\)](03_augmented.md#mint-to-deposit-hatch-phase)
* [Mint to Deposit and Burn to Withdraw \(Open Phase\)](03_augmented.md#mint-to-deposit-and-burn-to-withdraw-open-phase)
* [Outcome Payment, Settlement, Share Withdrawal](03_augmented.md#outcome-payment-settlement-share-withdrawal)

## Bond Configuration

### Curve Function

In this tutorial, an augmented function bond will be created. The augmented function implemented by the Bonds module can be represented by the below formulas [\[ref\]](https://medium.com/giveth/deep-dive-augmented-bonding-curves-3f1f7c1fa751):

* Initial reserve:

  ![initial reserve](../.gitbook/assets/augmented1.png)

* Initial supply:

  ![initial supply](../.gitbook/assets/augmented2.png)

* Constant power function invariant:

  ![constant power](../.gitbook/assets/augmented3.png)

* Invariant function:

  ![invariant function](../.gitbook/assets/augmented4.png)

* Pricing function:

  ![pricing function](../.gitbook/assets/augmented5.png)

* Reserve function:

  ![reserve function](../.gitbook/assets/augmented6.png)

From all of the above formulas, the four constants that we need to come up with ourselves are:

* `d0`: the total initial raise of reserve tokens, which will be split between the initial reserve `R0` and the initial funding `theta * d0`
* `p0`: the fixed price per token during the hatch phase, used to determine the initial supply `S0`
* `theta`: the initial allocation \(as a percentage of initial raise `d0`\), i.e. the percentage allocated directly to the funding pool \(a.k.a fee address\)
* `kappa`: a polynomial degree representing the steepness of the price curve

In this tutorial, the values picked will be: `d0=500.0`, `p0=0.01`, `theta=0.4`, `kappa=3.0`, which gives us the below curve.

![augmented graph](../.gitbook/assets/augmented7.png)

Generated using: [https://github.com/BlockScience/cadCAD-Tutorials/tree/master/00-Reference-Mechanisms/01-augmented-bonding-curve](https://github.com/BlockScience/cadCAD-Tutorials/tree/master/00-Reference-Mechanisms/01-augmented-bonding-curve)

### Fees

In this tutorial, the transaction fee percentage will be set to `0%`, the main reason being that the augmented function has an integrated fee percentage `theta` for the hatch phase. The exit fee percentage will be `0.1%`

### Outcome Payment

The outcome payment is an optional \(non-enforceable\) promise that the bond creator makes to investors specified in terms of a token amount. This is the amount that will need to be deposited into the bond to transition it from the `OPEN` state to the `SETTLE` state, once the goals of the bond have been reached as a result of the reserve token funding that the bond received.

This new reserve will be immediately available to all bond token holders, and the amount available to each holder depends on the amount of bond tokens that they hold. In the case of this tutorial, this will be set to `100000stake`.

### Other Customisation

Other customisation options that this tutorial will not go into is the ability to disable sells \(burns\), and the ability to have multiple signers as the creators/editors of the bond. In this tutorial, sells will be enabled, and the signer will be set to the address underlying the `shaun` account \(created when running `make run_with_data`\).

Additionally, sanity rate and sanity margin percentage only apply to swapper functions and so they will not be discussed. In this tutorial, these were thus just set to `0`.

## Bond Creation

The bond, with the above configurations, can be created as follows:

```bash
SHAUNADDR="$(bondscli keys show shaun -a)"
FEEADDR="$(bondscli keys show fee -a)"

bondscli tx bonds create-bond \
  --token=demo \
  --name="My Bond" \
  --description="Description about my bond" \
  --function-type=augmented_function \
  --function-parameters="d0:500.0,p0:0.01,theta:0.4,kappa:3.0" \
  --reserve-tokens=stake \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEEADDR" \
  --max-supply=1000000demo \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells \
  --signers="$SHAUNADDR" \
  --batch-blocks=1 \
  --outcome-payment="100000stake" \
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
    "creator": "cosmos1ew5xdhcpamfex0qqu0h8usmek4e3nnc8sm27jc",
    "function_type": "augmented_function",
    "function_parameters": [
      {
        "param": "theta",
        "value": "0.400000000000000000"
      },
      {
        "param": "kappa",
        "value": "3.000000000000000000"
      },
      {
        "param": "d0",
        "value": "500.000000000000000000"
      },
      {
        "param": "p0",
        "value": "0.010000000000000000"
      },
      {
        "param": "R0",
        "value": "300.000000000000000000"
      },
      {
        "param": "S0",
        "value": "50000.000000000000000000"
      },
      {
        "param": "V0",
        "value": "416666666666.666666666666666667"
      }
    ],
    "reserve_tokens": [
      "stake"
    ],
    "tx_fee_percentage": "0.000000000000000000",
    "exit_fee_percentage": "0.100000000000000000",
    "fee_address": "cosmos1g9caxdsgqj2060lp4fjx533svax57v9alv8qc4",
    "max_supply": {
      "denom": "demo",
      "amount": "1000000"
    },
    "order_quantity_limits": [],
    "sanity_rate": "0.000000000000000000",
    "sanity_margin_percentage": "0.000000000000000000",
    "current_supply": {
      "denom": "demo",
      "amount": "0"
    },
    "current_reserve": [],
    "allow_sells": false,
    "signers": [
      "cosmos1ew5xdhcpamfex0qqu0h8usmek4e3nnc8sm27jc"
    ],
    "batch_blocks": "1",
    "outcome_payment": [
      {
        "denom": "res",
        "amount": "100000"
      }
    ],
    "state": "HATCH"
  }
}
```

Note that some extra fields that we did not input ourselves are present. Some of these were discussed in previous tutorials. However, note that in the case of the augmented function:

* the function parameters were extended to include newly calculated values `R0`, `S0`, and `V0`, which are based on the formulae presented in the [Curve Function](03_augmented.md#curve-function)
* the intial state is actually `HATCH`, representing the hatch phase, rather than `OPEN`

## Mint to Deposit \(Hatch Phase\)

During the hatch phase, one can only perform a buy \(mint-to-deposit\). The buying price will be `p0`, and we can query this for confirmation using `bondscli q bonds current-price demo`, which gives:

```bash
[
  {
    "denom": "stake",
    "amount": "0.010000000000000000"
  }
]
```

Given that the initial supply `S0` is `50000`, it will require `50000demo` in order for the augmented curve to transition to the `OPEN` phase. We can go ahead and just perform a buy of `50000demo` in a single buy. The expected price will be `0.01 x 50000 = 500stake`. We can confirm this using `bondscli q bonds buy-price 50000demo`, which gives:

```bash
...
"total_prices": [
{
  "denom": "stake",
  "amount": "500"
}
],
...
```

Note that this matches the initial raise `d0`. Also note that since we are not charging any transaction fees, the total price quoted by the query matches exactly our calculations above.

We can perform the buy as follows, with `500stake` as the max spend. The account used is the `miguel` account \(created when running `make run_with_data`\).

```bash
bondscli tx bonds buy 50000demo 500stake \
  --from miguel \
  --keyring-backend=test \
  --broadcast-mode block \
  --gas-prices=0.025stake \
  -y
```

We can query the `miguel` account to confirm that the demo tokens have reached the account by using `bondscli q account $(bondscli keys show miguel -a)`. A maximum of 2 blocks-worth of time might need to pass for the order in the batch to get processed.

```bash
...
"coins": [
  {
    "denom": "demo",
    "amount": "50000"
  },
  ...
  {
    "denom": "stake",
    "amount": "99994500"
  }
],
...
```

Note how the account now has `50000demo` and `99994500stake` \(a `5500stake` decrease!\). The decrease in stake includes the buying price charged `500stake` and the blockchain gas fees `5000stake`.

We can also confirm the supply and reserve values and that the bond has transitioned to the `OPEN` phase by querying the bond using `bondscli q bonds bond demo` which gives the below result. Note how the current reserve matches `R0=300`.

```bash
...
"current_supply": {
  "denom": "demo",
  "amount": "50000"
},
"current_reserve": [
  {
    "denom": "stake",
    "amount": "300"
  }
],
...
"state": "OPEN"
...
```

The remaining `200` out of the `500stake` deposited were sent to the funding pool \(i.e. fee address\) and can be queried using `bondscli q account "$FEEADDR"`. Note that the bond creator is expected to have access to this funding pool and will be able to safely use any funds send to it.

## Mint to Deposit and Burn to Withdraw \(Open Phase\)

Now that the `OPEN` phase has been reached, we can query the price again using `bondscli q bonds current-price demo`, which gives:

```bash
[
  {
    "denom": "stake",
    "amount": "0.018000000000000000"
  }
]
```

This can be matched up with the pricing function presented in the [Curve Function](03_augmented.md#curve-function) section:

![pricing function](../.gitbook/assets/augmented8.png)

Given that mint-to-deposit and burn-to-withdraw during the `OPEN` phase have been covered in previous tutorials, these will not be covered again in this tutorial.

## Outcome Payment, Settlement, Share Withdrawal

As a refresher, the outcome payment is an amount of tokens that the bond creator had indicated would be paid to the bond once certain goals were reached. In the case of this tutorial, this is `100000stake`.

Let's assume that those goals were reached and the bond creator wants to make the outcome payment. Note that anyone with enough tokens is able to make the outcome payment, not just the bond creator.

Before making the payment, it is interesting to query the returns from selling before the outcome payment reaches the reserves, using `bondscli q bonds sell-return 50000demo`, which gives:

```bash
...
"total_returns": [
  {
    "denom": "stake",
    "amount": "299"
  }
  ],
  "total_fees": [
  {
    "denom": "stake",
    "amount": "1"
  }
]
...
```

Note that the maximum that the user can get back at the moment is the exact amount that was initially invested, `300stake`, minus an exit fee of `1stake`.

Now let's make the outcome payment from the bond creator. The account used is the `shaun` account \(created when running `make run_with_data`\).

```bash
bondscli tx bonds make-outcome-payment demo \
  --from shaun \
  --keyring-backend=test \
  --broadcast-mode block \
  --gas-prices=0.025stake \
  -y
```

This causes a state transition from `OPEN` to `SETTLE` and adds `100000stake` to the reserve. Both of these can be confirmed by querying the bond using `bondscli q bonds bond demo`, which gives:

```bash
...
"current_reserve": [
  {
    "denom": "stake",
    "amount": "100300"
  }
],
...
"state": "SETTLE"
...
```

At this stage, both buys and sells have been disabled and the only action that is possible is for bond token holders to withdraw their share of the reserve pool by performing a share withdrawal. The account used is the `miguel` account \(created when running `make run_with_data`\).

```text
bondscli tx bonds withdraw-share demo \
  --from miguel \
  --keyring-backend=test \
  --broadcast-mode block \
  --gas-prices=0.025stake \
  -y
```

Since the `miguel` account held 100% of the bond token supply, this share withdrawal sends all of the bond reserve to `miguel` and burns the entire bond token supply \(sent by `miguel`, which held all of these\). In fact if we query the bond one last time, this information can be confirmed, using `bondscli q bonds bond demo`, which gives:

```bash
...
"current_supply": {
  "denom": "demo",
  "amount": "0"
},
"current_reserve": [],
...
```

Querying the `miguel` account reveals that the account no longer holds any `demo` tokens and has received `100300stake`.

