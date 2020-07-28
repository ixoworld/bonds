package types

import "math"

// integer_units = 10**12 // account for decimal places to a token
// scale_units = 10**6 // millions of tokens, million of DAI
// mu = integer_units*scale_units

// value function for a given state (R,S)
func Invariant(R, S, kappa float64) float64 {
	return math.Pow(S, kappa) / R
}

// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// return Supply S as a function of reserve R
func Supply(R, kappa, V0 float64) float64 {
	return math.Pow(V0*R, 1/kappa)
}

func Reserve(S, kappa, V0 float64) float64 {
	return math.Pow(S, kappa) / V0
}

// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// return a spot price P as a function of reserve R
func SpotPrice(R, kappa, V0 float64) float64 {
	return (kappa * math.Pow(R, (kappa-1)/kappa)) / math.Pow(V0, 1/kappa)
}

// for a given state (R,S)
// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// deposit deltaR to Mint deltaS
// with realized price deltaR/deltaS
func Mint(deltaR, R, S, kappa, V0 float64) (float64, float64) {
	deltaS := math.Pow(V0*(R+deltaR), 1/kappa) - S
	realizedPrice := deltaR / deltaS
	return deltaS, realizedPrice
}

func MintAlt(deltaS, R, S, kappa, V0 float64) (float64, float64) {
	deltaR := (math.Pow(deltaS+S, kappa) / V0) - R
	realizedPrice := deltaR / deltaS
	return deltaR, realizedPrice
}

// for a given state (R,S)
// given a value function (parameterized by kappa)
// and an invariant coeficient V0
// burn deltaS to Withdraw deltaR
// with realized price deltaR/deltaS
func Withdraw(deltaS, R, S, kappa, V0 float64) (float64, float64) {
	deltaR := R - math.Pow(S-deltaS, kappa)/V0
	realizedPrice := deltaR / deltaS
	return deltaR, realizedPrice
}
