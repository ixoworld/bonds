package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// integer_units = 10**12 // account for decimal places to a token
// scale_units = 10**6 // millions of tokens, million of DAI
// mu = integer_units*scale_units

// value function for a given state (R,S)
func Invariant(R, S sdk.Dec, kappa int64) sdk.Dec {
	return Power(S, uint64(kappa)).Quo(R)
}

// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// return Supply S as a function of reserve R
func Supply(R sdk.Dec, kappa int64, V0 sdk.Dec) sdk.Dec {
	result, err := ApproxRoot(V0.Mul(R), uint64(kappa))
	if err != nil {
		panic(err)
	}
	return result
}

// This is the reverse of Supply(...) function
func Reserve(S sdk.Dec, kappa int64, V0 sdk.Dec) sdk.Dec {
	return Power(S, uint64(kappa)).Quo(V0)
}

// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// return a spot price P as a function of reserve R
func SpotPrice(R sdk.Dec, kappa int64, V0 sdk.Dec) sdk.Dec {
	kappaDec := sdk.NewDecFromInt(sdk.NewInt(kappa))

	temp1, err := ApproxRoot(V0, uint64(kappa))
	if err != nil {
		panic(err)
	}
	temp2, err := ApproxRoot(Power(R, uint64(kappa)-1), uint64(kappa))
	if err != nil {
		panic(err)
	}
	return (kappaDec.Mul(temp2)).Quo(temp1)
}

// for a given state (R,S)
// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// deposit deltaR to Mint deltaS
// with realized price deltaR/deltaS
func Mint(deltaR, R, S sdk.Dec, kappa int64, V0 sdk.Dec) (sdk.Dec, sdk.Dec) {
	temp, err := ApproxRoot(V0.Mul(R.Add(deltaR)), uint64(kappa))
	if err != nil {
		panic(err)
	}
	deltaS := temp.Sub(S)
	realizedPrice := deltaR.Quo(deltaS)
	return deltaS, realizedPrice
}

func MintAlt(deltaS, R, S sdk.Dec, kappa int64, V0 sdk.Dec) (sdk.Dec, sdk.Dec) {
	deltaR := (Power(deltaS.Add(S), uint64(kappa)).Quo(V0)).Sub(R)
	realizedPrice := deltaR.Quo(deltaS)
	return deltaR, realizedPrice
}

// for a given state (R,S)
// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// burn deltaS to Withdraw deltaR
// with realized price deltaR/deltaS
func Withdraw(deltaS, R, S sdk.Dec, kappa int64, V0 sdk.Dec) (sdk.Dec, sdk.Dec) {
	deltaR := R.Sub(Power(S.Sub(deltaS), uint64(kappa)).Quo(V0))
	realizedPrice := deltaR.Quo(deltaS)
	return deltaR, realizedPrice
}
