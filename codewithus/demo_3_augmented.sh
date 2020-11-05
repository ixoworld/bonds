#!/usr/bin/env bash

ALICE=$(bondscli keys show alice --keyring-backend=test -a)  # organiser
FEE=$(bondscli keys show fee --keyring-backend=test -a)      # funding pool

bondscli keys add p1 --keyring-backend=test
bondscli keys add p2 --keyring-backend=test
bondscli keys add p3 --keyring-backend=test

P1=$(bondscli keys show p1 --keyring-backend=test -a)  # participant 1 (new account)
P2=$(bondscli keys show p2 --keyring-backend=test -a)  # participant 2 (new account)
P3=$(bondscli keys show p3 --keyring-backend=test -a)  # participant 3 (new account)

# Creator-specified values:
#   d0 := 1000000       // initial raise (reserve [uatom])
#   p0 := 1             // initial price (reserve per token [uatom/ufit])
#   theta := 0          // initial allocation (percentage)
#   kappa := 3.0        // degrees of polynomial (i.e. x^2, x^4, x^6)

# Calculated values:
#   R0 = 1000000        // initial reserve (calculated by (1-theta)*d0)
#   S0 = 1000000        // initial supply (equal to R0 since p0=1)
#   V0 = 1000000000000  // invariant which relates supply and reserve

# Create an Augmented Bonding Curve
bondscli tx bonds create-bond \
  --token=ufit \
  --name="FIT Initiative" \
  --description="An incentivised fitness initiative" \
  --function-type=augmented_function \
  --function-parameters="d0:1000000,p0:1,theta:0,kappa:3.0" \
  --reserve-tokens=uatom \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0 \
  --fee-address="$FEE" \
  --max-supply=20000000ufit \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells \
  --signers="$ALICE" \
  --batch-blocks=1 \
  --outcome-payment="300000000uatom" \
  --from alice --keyring-backend=test --broadcast-mode block -y
# Query the created bond
bondscli q bonds bond ufit

# DAY 1 :: P1, P2, P3 all submit a valid claim, so organiser sends 1ATOM each
bondscli tx send "$ALICE" "$P1" 1000000uatom --from alice --keyring-backend=test --broadcast-mode block -y
bondscli tx send "$ALICE" "$P2" 1000000uatom --from alice --keyring-backend=test --broadcast-mode block -y
bondscli tx send "$ALICE" "$P3" 1000000uatom --from alice --keyring-backend=test --broadcast-mode block -y

# Each account now has 1ATOM
bondscli q account "$P1"
bondscli q account "$P2"
bondscli q account "$P3"

# P1 buys 1ATOM worth of FIT (=1FIT)      [note: 1ATOM=1FIT because p0=1]
bondscli tx bonds buy 1000000ufit 1000000uatom --from p1 --broadcast-mode block -y
# P2 buys 1ATOM worth of FIT (=0.26FIT)   [note: 1ATOM=0.26FIT because we are now in OPEN phase]
bondscli tx bonds buy 259921ufit 1000000uatom --from p2 --broadcast-mode block -y
# P3 keeps the ATOM and does not buy FIT

# Bond is now in OPEN phase
bondscli q bonds bond ufit

# Note that P1 and P2 are both able to sell (and re-buy) their FIT tokens during the OPEN phase

# A secondary market of FIT tokens (e.g. FIT<->USD liquidity pool) might be set up as well!

# Assume some days pass, participants exercise, more ATOM is exchanged, and project is coming to an end

# Organiser deems the project a success and makes the outcome payment of 300ATOM
# Note: amount not specified here since it was written into the bond at creation
bondscli tx bonds make-outcome-payment ufit --from alice --broadcast-mode block -y

# Bond is now in SETTLE phase (no buying/selling allowed)
bondscli q bonds bond ufit

# P1 withdraws their share
bondscli tx bonds withdraw-share ufit --from p1 --keyring-backend=test --broadcast-mode block -y
# P2 withdraws their share
bondscli tx bonds withdraw-share ufit --from p2 --keyring-backend=test --broadcast-mode block -y

# P1 balance is equal to that of P2 (each gets a proportional share of the outcome payment)
bondscli q account "$P1"
bondscli q account "$P2"

# Of course P3 still only has 1ATOM
bondscli q account "$P3"
