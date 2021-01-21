package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"math/rand"
	"strconv"
)

func getNextBondName() string {
	return tokenPrefix + strconv.Itoa(totalBondCount+1)
}

func incrementBondCount() {
	totalBondCount += 1
}

func newSwapperBond(token string) {
	swapperBonds = append(swapperBonds, token)
}

func getRandomBondName(r *rand.Rand) (bondName string, ok bool) {
	if totalBondCount == 0 {
		return "", false
	} else if totalBondCount == 1 {
		return tokenPrefix + "1", true
	} else {
		randInt := simulation.RandIntBetween(r, 1, totalBondCount)
		return tokenPrefix + strconv.Itoa(randInt), true
	}
}

func getRandomBondNameExcept(r *rand.Rand, except string) (bondName string, ok bool) {
	if totalBondCount == 0 {
		return "", false
	} else if totalBondCount == 1 {
		token1 := tokenPrefix + "1"
		if except == token1 {
			return "", false // except string is the only token that exists
		} else {
			return tokenPrefix + "1", true // except string is not a token
		}
	} else {
		token := except
		for token == except {
			randInt := simulation.RandIntBetween(r, 1, totalBondCount+1)
			token = tokenPrefix + strconv.Itoa(randInt)
		}
		return token, true
	}
}

func getRandomNonEmptyString(r *rand.Rand) string {
	return simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 1, 100))
}

func getRandomFunctionType(r *rand.Rand) string {
	switch simulation.RandIntBetween(r, 0, 4) {
	case 0:
		return types.PowerFunction
	case 1:
		return types.SigmoidFunction
	case 2:
		return types.SwapperFunction
	case 3:
		return types.AugmentedFunction
	default:
		panic("function type integer out of bounds")
	}
}

func getRandomFunctionParameters(r *rand.Rand, functionType string, genesis bool) types.FunctionParams {
	switch functionType {
	case types.PowerFunction:
		m := simulation.RandIntBetween(r, 1, 100)
		n := simulation.RandIntBetween(r, 1, 5)
		c := simulation.RandIntBetween(r, 1, 1000)
		return types.FunctionParams{
			types.NewFunctionParam("m", sdk.NewDec(int64(m))),
			types.NewFunctionParam("n", sdk.NewDec(int64(n))),
			types.NewFunctionParam("c", sdk.NewDec(int64(c)))}
	case types.SigmoidFunction:
		a := simulation.RandIntBetween(r, 1, 10)
		b := simulation.RandIntBetween(r, 1, 10)
		c := simulation.RandIntBetween(r, 1, 10)
		return types.FunctionParams{
			types.NewFunctionParam("a", sdk.NewDec(int64(a))),
			types.NewFunctionParam("b", sdk.NewDec(int64(b))),
			types.NewFunctionParam("c", sdk.NewDec(int64(c)))}
	case types.AugmentedFunction:
		d0 := sdk.NewDec(int64(simulation.RandIntBetween(r, 1, 1000000)))
		p0 := simulation.RandomDecAmount(r, sdk.NewDec(10)).Add(sdk.SmallestDec())
		theta := simulation.RandomDecAmount(r, sdk.MustNewDecFromStr("0.9")).Add(sdk.SmallestDec())
		kappa := sdk.NewDec(int64(simulation.RandIntBetween(r, 1, 4)))
		functionParams := types.FunctionParams{
			types.NewFunctionParam("d0", d0),
			types.NewFunctionParam("p0", p0),
			types.NewFunctionParam("theta", theta),
			types.NewFunctionParam("kappa", kappa)}
		if genesis {
			R0 := d0.Mul(sdk.OneDec().Sub(theta))
			S0 := d0.Quo(p0)
			V0 := types.Invariant(R0, S0, kappa)

			functionParams = append(functionParams,
				types.FunctionParams{
					types.NewFunctionParam("R0", R0),
					types.NewFunctionParam("S0", S0),
					types.NewFunctionParam("V0", V0),
				}...)
		}
		return functionParams
	case types.SwapperFunction:
		return nil
	default:
		panic("unrecognized function type")
	}
}

func getRandomAllowSellsValue(r *rand.Rand) bool {
	if simulation.RandIntBetween(r, 1, 11) == 1 {
		return false
	} else { // 9 times out of 10, sells are allowed
		return true
	}
}

func getInitialBondState(functionType string) string {
	switch functionType {
	case types.AugmentedFunction:
		return types.HatchState
	default:
		return types.OpenState
	}
}

//noinspection GoNilness
func getDummyNonZeroReserve(reserveTokens []string) (reserve sdk.Coins) {
	for _, token := range reserveTokens {
		reserve = reserve.Add(sdk.NewCoin(token, sdk.OneInt()))
	}
	return reserve
}
