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

func getRandomFunctionParameters(r *rand.Rand, functionType string) types.FunctionParams {
	switch functionType {
	case types.PowerFunction:
		m := simulation.RandIntBetween(r, 1, 100)  // 12
		n := simulation.RandIntBetween(r, 1, 5)    // 5
		c := simulation.RandIntBetween(r, 1, 1000) // 100
		return types.FunctionParams{
			types.NewFunctionParam("m", sdk.NewDec(int64(m))),
			types.NewFunctionParam("n", sdk.NewDec(int64(n))),
			types.NewFunctionParam("c", sdk.NewDec(int64(c)))}
	case types.SigmoidFunction:
		a := simulation.RandIntBetween(r, 1, 10) // 3
		b := simulation.RandIntBetween(r, 1, 10) // 5
		c := simulation.RandIntBetween(r, 1, 10) // 1
		return types.FunctionParams{
			types.NewFunctionParam("a", sdk.NewDec(int64(a))),
			types.NewFunctionParam("b", sdk.NewDec(int64(b))),
			types.NewFunctionParam("c", sdk.NewDec(int64(c)))}
	case types.SwapperFunction:
		return nil
	default:
		panic("unrecognized function type")
	}
}

func getRandomAllowSellsValue(r *rand.Rand) string {
	if simulation.RandIntBetween(r, 1, 11) == 1 {
		return types.FALSE
	} else { // 9 times out of 10, sells are allowed
		return types.TRUE
	}
}

//noinspection GoNilness
func getDummyNonZeroReserve(reserveTokens []string) (reserve sdk.Coins) {
	for _, token := range reserveTokens {
		reserve = reserve.Add(sdk.Coins{sdk.NewCoin(token, sdk.OneInt())})
	}
	return reserve
}
