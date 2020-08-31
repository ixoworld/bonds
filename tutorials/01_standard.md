# Creating a Bond for Continuous Funding using Standard Bonds Module Functions

## Contents

- [Bond Configuration](#bond-configuration)
- [Bond Creation](#bond-creation)
- [Mint to Deposit](#mint-to-deposit)
- [Burn to Withdraw](#burn-to-withdraw)

## Bond Configuration

### Curve Function

In this tutorial, a power function bond will be created. The power function implemented by the Bonds module is shown below, where `y` represents the price per bond token, in reserve tokens, for a specific supply `x` of bond tokens:

<img alt="power function price" src="./img/power1.png" height="20"/>
 
The remaining values `m`, `n`, and `c` are constants that we need to come up with ourselves. In this tutorial, the values picked will be: `m=12`, `n=2`, `c=100`, which gives us the below curve.

- Increasing/decreasing `m` and `n` makes the incline steeper or more gradual, respectively, meaning a quicker increase in price or a more gradual increase. `n` has a greater effect on this, given that it is a power.
- Increasing/decreasing `c` lifts/lowers the curve, respectively, which means greater prices throughout the curve but maintaining the steepness of the curve.

<img alt="power function graph" src="img/power3.png"/>

From the above power function, a reserve function is deduced by integrating the power curve. This is shown below and includes the same `m`, `n`, `c` constants that were in the original function, but also includes `r`, which is the reserve balance that is required to be in place for the bond token's supply to be `x`. 

<img alt="power function reserve" src="./img/power2.png" height="40"/>

This function allows us to calculate the price of buying a number of bond tokens or the returns when selling bond tokens.

### Bond Token and Name

Next we can think about what denomination to use for the bond token and what to name our bond. We will go with a simple `demo` token and name the bond `"My Bond"`. This means that `demo` tokens will be minted whenever an account buys into the curve (i.e. mints-to-deposit).

We also need to decide what the maximum supply will be for the bond token. This places a hard limit on the number of bond tokens that can ever exist. In this tutorial, we will set a maximum supply of `1000000demo`.

### Reserve Token/s

The reserve tokens are the tokens that accounts will send to the bond in order to mint bond tokens. This means that the accounts will need to have these reserve tokens. For this reason, we will simply use `stake` as the reserve token. We could have easily decided to instead use another reserve token.

### Fees

We can decide to charge a transaction and/or exit fee. Transaction fees are charged whenever an account buys/sells (i.e. mints/burns) into/from the curve, whereas an exit fee is charged only when the account sells from the curve. Setting an exit fee disincentivises accounts from selling and is a front-running deterrent.

In this tutorial, we will set the the transaction fee to `0.5%` and the exit fee to `0.1%`. This will mean that when an account sells, a total of `0.6%` will be taken as fees. Fees are sent to a fee address that is picked by the bond creator. In this tutorial, the address underlying the `fee` account (created when running `make run_with_data`) will be used.

### Order Quantity Limits

We can also decide to place limits on order quantity. This serves to limit the amount of bond tokens that can be bought/sold (i.e. minted/burned) at once in a single order. In this tutorial, we will set an order limit of `100demo`. We could have decided to not put any limits at all.

### Batch Blocks

Since the Bonds module implements batched orders, we need to decide on the number of blocks that each batch is valid for. In this tutorial we will go with 2 blocks, meaning that a batch is only valid for 2 block and will get cleared out at the end of every second block.

### Other Customisation

Other customisation options that this tutorial will not go into is the ability to disable sells (burns), the ability to have multiple signers as the creators/editors of the bond, and the ability to add an outcome payment. In this tutorial, sells will be enabled, the signer will be set to the address underlying the `shaun` account (created when running `make run_with_data`), and there will be no outcome payment (discussed in the augmented function tutorial).

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
  --function-type=power_function \
  --function-parameters="m:12,n:2,c:100" \
  --reserve-tokens=stake \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEEADDR" \
  --max-supply=1000000demo \
  --order-quantity-limits=100demo \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells \
  --signers="$SHAUNADDR" \
  --batch-blocks=2 \
  --outcome_payment="" \
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
    "creator": "cosmos1grafevmzch5xv4909x5uet0t7nfdll3n97jqjp",
    "function_type": "power_function",
    "function_parameters": [
      {
        "param": "m",
        "value": "12.000000000000000000"
      },
      {
        "param": "n",
        "value": "2.000000000000000000"
      },
      {
        "param": "c",
        "value": "100.000000000000000000"
      }
    ],
    "reserve_tokens": [
      "stake"
    ],
    "tx_fee_percentage": "0.500000000000000000",
    "exit_fee_percentage": "0.100000000000000000",
    "fee_address": "cosmos169rg7xj80dc6qf036q6ahpvvalc5zqp0qvze6z",
    "max_supply": {
      "denom": "demo",
      "amount": "1000000"
    },
    "order_quantity_limits": [
      {
        "denom": "demo",
        "amount": "100"
      }
    ],
    "sanity_rate": "0.000000000000000000",
    "sanity_margin_percentage": "0.000000000000000000",
    "current_supply": {
      "denom": "demo",
      "amount": "0"
    },
    "current_reserve": [],
    "allow_sells": true,
    "signers": [
      "cosmos1grafevmzch5xv4909x5uet0t7nfdll3n97jqjp"
    ],
    "batch_blocks": "1",
    "outcome_payment": [],
    "state": "OPEN"
  }
}
```

Note that some extra fields that we did not input ourselves are present:
- `current_supply`: stores the current bond token supply and increases/decreases whenever a buy/sell is performed
- `current_reserve`: stores the current reserve that has been sent to the bond as a result of buys and increases/decreases whenever a buy/sell is performed
- `state`: stores the current state of the bond, which throughout this tutorial will remain `OPEN`

We are also able to query the bond's current batch using `bondscli q bonds batch demo`, which should return the below. Since we have not performed any buys/sells/swaps, the associated fields are all zero or null. The blocks remaining starts at the `batch-blocks` value that we had picked, decreases by 1 at the end of each block, and is reset to `batch-blocks` as soon as it reaches 0.

In the case of this tutorial, since we set `batch-blocks` to 2, the `blocks_remaining` value will start at 2, go to 1, and back to 2 (since it will have reached 0). We will never see `blocks_remaining` reach 0.

```bash
{
  "type": "bonds/Batch",
  "value": {
    "token": "demo",
    "blocks_remaining": "2",
    "total_buy_amount": {
      "denom": "demo",
      "amount": "0"
    },
    "total_sell_amount": {
      "denom": "demo",
      "amount": "0"
    },
    "buy_prices": null,
    "sell_prices": null,
    "buys": null,
    "sells": null,
    "swaps": null
  }
}
```

## Mint to Deposit

Before performing a buy (mint-to-deposit), we can query the current price to buy the tokens. Say we want to buy `10demo`, we can perform the query `bondscli q bonds buy-price 10demo`, which gives:

```bash
{
  "adjusted_supply": {
    "denom": "demo",
    "amount": "0"
  },
  "prices": [
    {
      "denom": "stake",
      "amount": "5000"
    }
  ],
  "tx_fees": [
    {
      "denom": "stake",
      "amount": "25"
    }
  ],
  "total_prices": [
    {
      "denom": "stake",
      "amount": "5025"
    }
  ],
  "total_fees": [
    {
      "denom": "stake",
      "amount": "25"
    }
  ]
}
```

This query returns not just a single final price, but the fee-less price and fees separately, and the total prices and total fees. The most important value is the total prices, which includes the total fees. Note that these fees are separate from gas fees that need to be paid to make use of the blockchain.

These values can be verified by using the reserve function presented in the [Curve Function](#curve-function) section and the transaction fee percentage set during [bond creation](#bond-creation-and-querying) (0.5%):
- The price is: `r = (m/(n+1))x^(n+1) + cx = 4(x^3) + 100x = 4(10^3) + 100(10) = 5000`
- The fee is based on the price: `fee = 0.5% of 5000 = 25`
- The total prices is thus: `5000 + 25 = 5025`

Note that in the above working, `x` was set to 10 because that is the supply that will be reached if the buy order goes through. If another buy order of `10demo` comes in, calculating the buy price will require calculating the reserve required for `20demo` and subtracting the current reserve balance, since the buyer only needs to pay for `10demo`, not `20demo`.

Now that we know the buy price, we can perform a buy of `10demo` with a maximum spend of `5025stake`. However, since the price can change as more buyers and sellers interact with the bond, it is a good idea to set the maximum spend higher than the buy price. Let's go with `5100stake`. The account used is the `miguel` account (created when running `make run_with_data`).

```bash
bondscli tx bonds buy 10demo 5100stake \
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
        "amount": "10"
      },
      ...
      {
        "denom": "stake",
        "amount": "99989975"
      }
    ],
...
```

Note how the account now has `10demo` and `99989975stake` (a `10025stake` decrease!). The decrease in stake includes the actual buy price charged `5025stake` and the blockchain gas fees `5000stake`.

## Burn to Withdraw

Before performing a sell (burn-to-withdraw), we can query the current returns for selling the tokens. Say we want to sell `10demo`, we can perform the query `bondscli q bonds sell-return 10demo`, which gives:

```bash
{
  "adjusted_supply": {
    "denom": "demo",
    "amount": "10"
  },
  "returns": [
    {
      "denom": "stake",
      "amount": "5000"
    }
  ],
  "tx_fees": [
    {
      "denom": "stake",
      "amount": "25"
    }
  ],
  "exit_fees": [
    {
      "denom": "stake",
      "amount": "5"
    }
  ],
  "total_returns": [
    {
      "denom": "stake",
      "amount": "4970"
    }
  ],
  "total_fees": [
    {
      "denom": "stake",
      "amount": "30"
    }
  ]
}
```

Note how in this case, both transaction (0.5%) and exit fees (0.5%) are included. Note also that the fee-less returns value (`5000`) matches the fee-less price value that we queried when [minting-to-deposit](#mint-to-deposit). In this case, the fee gets charged by the account getting less tokens in return, rather than the account having to pay more.

We can perform a sell of `10demo` as shown below. The account used is the `miguel` account (created when running `make run_with_data`).

```bash
bondscli tx bonds sell 10demo \
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
        "denom": "stake",
        "amount": "99989945"
      }
    ],
...
```

Note how the account now has `99989945stake` (a `30stake` decrease!). The decrease in stake is due to the blockchain gas fees (`5000stake`) being greater than the returns `4970stake` that the account got by selling the demo tokens (`4970 - 5000 = 30`). In a more typical and happier scenario, the returns from selling bond tokens would be much greater than the gas fees!
