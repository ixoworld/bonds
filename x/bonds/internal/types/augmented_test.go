package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"math"
	"strings"
	"testing"
)

func arange(start, stop, step float64) []float64 {
	N := int(math.Ceil((stop - start) / step))
	rnge := make([]float64, N, N)
	i := 0
	for x := start; x < stop; x += step {
		rnge[i] = x
		i += 1
	}
	return rnge
}

func printLines(title string, values []sdk.Dec) {
	print(title + " = [")
	for i, value := range values {
		index := strings.Index(value.String(), ".") + 7
		if i == len(values)-1 {
			fmt.Print(value.String()[:index])
		} else {
			fmt.Print(value.String()[:index] + ", ")
		}
	}
	println("]")
}

func TestExample1(t *testing.T) {
	d0 := sdk.MustNewDecFromStr("5.0")  // million DAI
	p0 := sdk.MustNewDecFromStr("0.01") // DAI per tokens
	theta := sdk.MustNewDecFromStr("0.4")

	R0 := d0.Mul(sdk.OneDec().Sub(theta)) // million DAI
	S0 := d0.Quo(p0)

	kappa := int64(3)
	V0 := Invariant(R0, S0, kappa)

	expectedR0 := sdk.MustNewDecFromStr("3.0")
	expectedS0 := sdk.MustNewDecFromStr("500.0")
	expectedV0 := sdk.MustNewDecFromStr("41666666.666666666666666667")

	require.Equal(t, expectedR0, R0)
	require.Equal(t, expectedS0, S0)
	require.Equal(t, expectedV0, V0)

	reserveF64 := arange(0, 100, .01)
	reserve := make([]sdk.Dec, len(reserveF64))
	for i, r := range reserveF64 {
		reserve[i] = sdk.MustNewDecFromStr(fmt.Sprintf("%f", r))
	}

	supp := make([]sdk.Dec, len(reserve))
	for i, r := range reserve {
		supp[i] = Supply(r, kappa, V0)
	}

	price := make([]sdk.Dec, len(reserve))
	for i, r := range reserve {
		price[i] = SpotPrice(r, kappa, V0)
	}

	printLines("reserve", reserve)
	printLines("supp", supp)
	printLines("price", price)
}
