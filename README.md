# Bonds Module
The Bonds module is a custom [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) module that provides universal token bonding curve functions to mint, burn or swap any token in a Cosmos blockchain.

In the future, once the Cosmos Inter-Blockchain Communication (IBC) protocol is available, this should enable cross-network exchanges of tokens at algorithmically-determined prices.

The Bonds module can be deployed through Cosmos Hubs and Zones to deliver applications such as:
* Automated market-makers (see [Uniswap](https://uniswap.io))
* Decentralised exchanges (see [Bancor](https://bancor.network))
* Curation markets (see [Relevant](https://github.com/relevant-community/contracts/tree/bondingCurves/contracts))
* Development Impact Bonds (see ixo alpha-Bonds)
* Continuous organisations (see [Moloch DAO](https://molochdao.com/))

> [Hayek famously said](https://books.google.co.uk/books?id=Udi_BwAAQBAJ&pg=PA32&lpg=PA32&dq=%22prices+are+an+instrument+of+communication+and+guidance+which+embody+more+information+than+we+directly+have%22&source=bl&ots=LMFRhcW0QS&sig=ACfU3U0I6_J3_uBI96ZFKAxCo-p6yh_eNg&hl=en&sa=X&ved=2ahUKEwimguWHpOjmAhWFTBUIHQCYASYQ6AEwAnoECAkQAQ#v=onepage&q=%22prices%20are%20an%20instrument%20of%20communication%20and%20guidance%20which%20embody%20more%20information%20than%20we%20directly%20have%22&f=false) that "...prices are an instrument of communication and guidance which embody more information than we directly have".

## Module functions

Any Cosmos application chain that implements the Bonds module is able to perform functions such as:
* Issue a new token with custom parameters.
* Pool liquidity for reserves.
* Provide continuous funding.
* Automatically mint and burn tokens at deterministic prices.
* Swap tokens atomically within the same network.
* Exchange tokens across networks, with the IBC protocol.
* Batch token transactions to prevent front-running
* Launch a decentralised autonomous initial coin offerings ([DAICO](https://ethresear.ch/t/explanation-of-daicos/465))
* ...and other **DeFi**ant innovations.

## Pricing algorithm libraries
The Bonds module framework supports libraries for all types of pricing algorithms, such as:
* Exponential
* Logarithmic
* Negative exponential
* Constant product
* Positive initial price
* Quasi-polynomial
* Reserved Supply (Augmented)

Each formula is specified within the module library. 
This includes:
* Derived Mint equation
* Derived Burn equation

Updates to the module pricing functions must pass through a network governance process to update the module on all nodes, for changes to be made. This is an important security feature.

## Parameters
Each bond has an initial set of constant state (invariant) parameters that cannot be updated once these have been initialised. Parameters include:
* Pricing function (the algorithm that will be used)
* Issuer
* Token name
* Token symbol
* Reserve wallet address
* Transaction fee rate
* Exit tax rate
* Fee wallet address
* Maximum token supply
* Order quantity limits
* Sanity rates

When a bond transaction (such as buy, sell, swap) is submitted, this includes the variable parameters:
* Order quantity
* Maximum price
* Wallet address

Some of the parameters of the bond may be edited:
* Token name
* Sanity rates
* Order quantity limits

## Building and Running

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

This repository is set up as a Cosmos SDK application and contains the Bonds module under ```./x/bonds/```.

To build and run the application:
```bash
make run
```

Alternatively, to run with one of the users set up to use Ledger:
```bash
make run_ledger
```

To build and run the application with some preset accounts:
```bash
make run_with_data
```

### Demos

To run a demo (requires application to be run using `run_with_data`):
```bash
make demo
```

The demo consists of:
- Bond creation
- Bond querying
- A mix of buys and sells

To run a more specific demo, check out the `scripts/more_demos/` folder.

## Tutorials

Guided tutorials are also provided and can be found in the tutorials folder [here](tutorials/README.md)!
