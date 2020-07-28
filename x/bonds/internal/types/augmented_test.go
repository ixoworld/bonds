package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
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

func printLines(title string, values []float64) {
	print(title + " = [")
	for i, value := range values {
		if i == len(values)-1 {
			fmt.Print(fmt.Sprintf("%f", value))
		} else {
			fmt.Print(fmt.Sprintf("%f, ", value))
		}
	}
	println("]")
}

func TestExample1(t *testing.T) {
	d0 := 5.  // million DAI
	p0 := .01 // DAI per tokens
	theta := .4

	R0 := d0 * (1 - theta) // million DAI
	S0 := d0 / p0

	kappa := 6.
	V0 := Invariant(R0, S0, kappa)

	expectedR0 := 3.0
	expectedS0 := 500.0
	expectedV0 := 5208333333333333.0

	require.Equal(t, expectedR0, R0)
	require.Equal(t, expectedS0, S0)
	require.Equal(t, expectedV0, V0)

	reserve := arange(0, 100, .01)

	supp := make([]float64, len(reserve))
	for i, r := range reserve {
		supp[i] = Supply(r, kappa, V0)
	}

	price := make([]float64, len(reserve))
	for i, r := range reserve {
		price[i] = SpotPrice(r, kappa, V0)
	}

	printLines("reserve", reserve)
	printLines("supp", supp)
	printLines("price", price)
}

func TestExample2(t *testing.T) {
	d0 := 5.  // million DAI
	p0 := .01 // DAI per tokens
	theta := .4

	R0 := d0 * (1 - theta) // million DAI
	S0 := d0 / p0

	kappa := 6.
	V0 := Invariant(R0, S0, kappa)

	expectedR0 := 3.0
	expectedS0 := 500.0
	expectedV0 := 5208333333333333.0

	require.Equal(t, expectedR0, R0)
	require.Equal(t, expectedS0, S0)
	require.Equal(t, expectedV0, V0)

	// given V0 and kappa
	// sweep the reserve
	reserve := arange(.01, 100, .01)

	price := make([]float64, len(reserve))
	for i, r := range reserve {
		price[i] = SpotPrice(r, kappa, V0)
	}

	// realized price for withdrawing burning .1% of tokens
	withdrawPrice := make([]float64, len(reserve))
	for i, r := range reserve {
		_, withdrawPrice[i] = Withdraw(Supply(r, kappa, V0)/1000, r, Supply(r, kappa, V0), kappa, V0)
	}

	// realized price for depositing .1% more Xdai into the reserve
	mintPrice := make([]float64, len(reserve))
	for i, r := range reserve {
		_, mintPrice[i] = Mint(r/1000, r, Supply(r, kappa, V0), kappa, V0)
	}

	printLines("reserve", reserve)
	printLines("price", price)
	printLines("withdrawPrice", withdrawPrice)
	printLines("mintPrice", mintPrice)
}

//func TestExample3(t *testing.T) {
//	d0 := 5.  // million DAI
//	p0 := .01 // DAI per tokens
//	theta := .4
//
//	R0 := d0 * (1 - theta) // million DAI
//	S0 := d0 / p0
//
//	kappa := 6.
//	V0 := Invariant(R0, S0, kappa)
//
//	// given V0 and kappa
//	R := 20.
//	S := Supply(R, kappa, V0)
//	p := SpotPrice(R, kappa, V0)
//	// sweep the transaction fraction
//	TXF := logspace(-6, 0, num = 1000)
//
//	// realized price for withdrawing burning .1% of tokens
//	withdrawPrice := make([]float64, len(TXF))
//	for i, txf := range TXF {
//		_, withdrawPrice[i] = Withdraw(S*txf, R, S, kappa, V0)
//	}
//
//	// realized price for depositing .1% more Xdai into the reserve
//	mintPrice := make([]float64, len(TXF))
//	for i, txf := range TXF {
//		_, mintPrice[i] = Mint(R*txf, R, S, kappa, V0)
//	}
//
//	printLines("TXF", TXF)
//	printLines("withdrawPrice", withdrawPrice)
//	printLines("mintPrice", mintPrice)
//}
