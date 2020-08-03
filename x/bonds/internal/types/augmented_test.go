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
	d0 := sdk.MustNewDecFromStr("500.0")  // initial raise (reserve)
	p0 := sdk.MustNewDecFromStr("0.01")   // initial price (reserve per token)
	theta := sdk.MustNewDecFromStr("0.4") // funding fee fraction

	R0 := d0.Mul(sdk.OneDec().Sub(theta)) // initial reserve (raise minus funding)
	S0 := d0.Quo(p0)                      // initial supply

	kappa := int64(3)              // price exponent
	V0 := Invariant(R0, S0, kappa) // invariant

	expectedR0 := sdk.MustNewDecFromStr("300.0")
	expectedS0 := sdk.MustNewDecFromStr("50000.0")
	expectedV0 := sdk.MustNewDecFromStr("416666666666.666666666666666667")

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

func TestReserve(t *testing.T) {
	decimals := sdk.NewDec(100000) // 10^5
	testCases := []struct {
		reserve sdk.Dec
		kappa   int64
		V0      sdk.Dec
	}{
		{sdk.MustNewDecFromStr("0.05"), 1, sdk.MustNewDecFromStr("12345678.12345678")},
		{sdk.MustNewDecFromStr("5"), 2, sdk.MustNewDecFromStr("123456.123456")},
		{sdk.MustNewDecFromStr("500.500"), 3, sdk.MustNewDecFromStr("50000.50000")},
		{sdk.MustNewDecFromStr("50000.50000"), 4, sdk.MustNewDecFromStr("500.500")},
		{sdk.MustNewDecFromStr("123456.123456"), 5, sdk.MustNewDecFromStr("5")},
		{sdk.MustNewDecFromStr("12345678.12345678"), 6, sdk.MustNewDecFromStr("0.05")},
	}
	for _, tc := range testCases {
		calculatedSupply := Supply(tc.reserve, tc.kappa, tc.V0)
		calculatedReserve := Reserve(calculatedSupply, tc.kappa, tc.V0)

		tc.reserve = tc.reserve.Mul(decimals).TruncateDec()
		calculatedReserve = calculatedReserve.Mul(decimals).TruncateDec()

		require.Equal(t, tc.reserve, calculatedReserve)
	}
}
