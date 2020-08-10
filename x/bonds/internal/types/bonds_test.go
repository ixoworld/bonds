package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestExtraParameterRestrictions_Power(t *testing.T) {
	paramRestrictions := ExtraParameterRestrictions[PowerFunction]

	testCases := []struct {
		m           string
		n           string
		c           string
		expectError bool
	}{
		{"10", "10", "10", false},       // integers allowed for all
		{"0", "0", "0", false},          // zeroes allowed for all
		{"10.10", "10", "10.10", false}, // float m and c allowed
		{"10", "10.10", "10", true},     // float n not allowed
	}

	for _, tc := range testCases {
		mDec := sdk.MustNewDecFromStr(tc.m)
		nDec := sdk.MustNewDecFromStr(tc.n)
		cDec := sdk.MustNewDecFromStr(tc.c)
		err := paramRestrictions(FunctionParams{
			NewFunctionParam("m", mDec),
			NewFunctionParam("n", nDec),
			NewFunctionParam("c", cDec),
		}.AsMap())

		if tc.expectError {
			require.Error(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestExtraParameterRestrictions_Sigmoid(t *testing.T) {
	paramRestrictions := ExtraParameterRestrictions[SigmoidFunction]

	testCases := []struct {
		a           string
		b           string
		c           string
		expectError bool
	}{
		{"10", "10", "10", false},          // integers allowed for all
		{"0", "0", "10", false},            // zeroes allowed for a and b
		{"10", "10", "0", true},            // zero not allowed for c
		{"10.10", "10.10", "10.10", false}, // floats allowed for all
	}

	for _, tc := range testCases {
		aDec := sdk.MustNewDecFromStr(tc.a)
		bDec := sdk.MustNewDecFromStr(tc.b)
		cDec := sdk.MustNewDecFromStr(tc.c)
		err := paramRestrictions(FunctionParams{
			NewFunctionParam("a", aDec),
			NewFunctionParam("b", bDec),
			NewFunctionParam("c", cDec),
		}.AsMap())

		if tc.expectError {
			require.Error(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestExtraParameterRestrictions_Augmented(t *testing.T) {
	paramRestrictions := ExtraParameterRestrictions[AugmentedFunction]

	testCases := []struct {
		d0          string
		p0          string
		theta       string
		kappa       string
		expectError bool
	}{
		{"10", "10", "0.5", "10", false},      // valid values
		{"0", "10", "0.5", "10", true},        // d0 can NOT be 0
		{"10", "0", "0.5", "10", true},        // p0 can NOT be 0
		{"10", "10", "0", "10", false},        // theta can be 0
		{"10", "10", "0.5", "0", true},        // kappa can NOT be 0
		{"10", "10", "1", "10", true},         // theta can NOT be 1
		{"10", "10", "1.1", "10", true},       // theta can NOT be >1
		{"10", "10.10", "0.5", "10", false},   // p0 and theta can be floats
		{"10.10", "10.10", "0.5", "10", true}, // d0 can NOT be a float
		{"10", "10.10", "0.5", "10.10", true}, // kappa can NOT be a float
	}

	for _, tc := range testCases {
		d0Dec := sdk.MustNewDecFromStr(tc.d0)
		p0Dec := sdk.MustNewDecFromStr(tc.p0)
		thetaDec := sdk.MustNewDecFromStr(tc.theta)
		kappaDec := sdk.MustNewDecFromStr(tc.kappa)
		err := paramRestrictions(FunctionParams{
			NewFunctionParam("d0", d0Dec),
			NewFunctionParam("p0", p0Dec),
			NewFunctionParam("theta", thetaDec),
			NewFunctionParam("kappa", kappaDec),
		}.AsMap())

		if tc.expectError {
			require.Error(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestFunctionParamsAsMap(t *testing.T) {
	actualResult := functionParametersPower().AsMap()
	expectedResult := map[string]sdk.Dec{
		"m": sdk.NewDec(12),
		"n": sdk.NewDec(2),
		"c": sdk.NewDec(100),
	}
	require.Equal(t, expectedResult, actualResult)
}

func TestFunctionParamsString(t *testing.T) {
	testCases := []struct {
		params   FunctionParams
		expected string
	}{
		// Note: square brackets added below in the require.Equal()
		{FunctionParams{}, ""},
		{FunctionParams{NewFunctionParam("a", sdk.OneDec())},
			"{\"param\":\"a\",\"value\":\"1.000000000000000000\"}"},
		{functionParametersPower(), "" +
			"{\"param\":\"m\",\"value\":\"12.000000000000000000\"}," +
			"{\"param\":\"n\",\"value\":\"2.000000000000000000\"}," +
			"{\"param\":\"c\",\"value\":\"100.000000000000000000\"}"},
		{functionParametersSigmoid(), "" +
			"{\"param\":\"a\",\"value\":\"3.000000000000000000\"}," +
			"{\"param\":\"b\",\"value\":\"5.000000000000000000\"}," +
			"{\"param\":\"c\",\"value\":\"1.000000000000000000\"}"},
	}
	for _, tc := range testCases {
		require.Equal(t, "["+tc.expected+"]", tc.params.String())
	}
}

func TestFunctionParamsAsMapReturnIsAsExpected(t *testing.T) {
	actualResult := functionParametersPower().AsMap()
	expectedResult := map[string]sdk.Dec{
		"m": sdk.NewDec(12), "n": sdk.NewDec(2), "c": sdk.NewDec(100)}
	require.Equal(t, expectedResult, actualResult)
}

func TestNewBondDefaultValuesAndSorting(t *testing.T) {
	customReserveTokens := []string{"b", "a"}
	customOrderQuantityLimits, _ := sdk.ParseCoins("100bbb,100aaa")
	sortedReserveTokens := []string{"a", "b"}
	sortedOrderQuantityLimits, _ := sdk.ParseCoins("100aaa,100bbb")

	bond := NewBond(initToken, initName, initDescription, initCreator,
		PowerFunction, functionParametersPower(), customReserveTokens,
		initTxFeePercentage, initExitFeePercentage, initFeeAddress, initMaxSupply,
		customOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks, initOutcomePayment, initState)

	expectedCurrentSupply := sdk.NewInt64Coin(bond.Token, 0)

	require.Equal(t, expectedCurrentSupply, bond.CurrentSupply)
	require.Equal(t, sortedReserveTokens, bond.ReserveTokens)
	require.Equal(t, sortedOrderQuantityLimits, bond.OrderQuantityLimits)
}

func TestGetNewReserveCoinReturnPasses(t *testing.T) {
	bond := getValidBond()

	require.Equal(t, sdk.NewInt64Coin(bond.Token, 0), bond.CurrentSupply)
}

func TestGetNewReserveDecCoins(t *testing.T) {
	bond := getValidBond()
	bond.ReserveTokens = []string{"aaa", "bbb"}

	amount := sdk.MustNewDecFromStr("10")
	actualResult := bond.GetNewReserveDecCoins(amount)

	expectedResult := sdk.NewDecCoins(sdk.NewCoins(
		sdk.NewInt64Coin("aaa", 10),
		sdk.NewInt64Coin("bbb", 10),
	))

	require.Equal(t, expectedResult, actualResult)
}

func TestGetPricesAtSupply(t *testing.T) {
	bond := getValidBond()
	// TODO: add more test cases

	// For augmented function tests
	baseMap := functionParametersAugmented().AsMap()
	//R0 := baseMap["d0"].Mul(sdk.OneDec().Sub(baseMap["theta"]))
	S0 := baseMap["d0"].Quo(baseMap["p0"])
	//V0 := Invariant(R0, S0, baseMap["kappa"].TruncateInt64())

	testCases := []struct {
		functionType      string
		functionParams    FunctionParams
		reserveTokens     []string
		supply            sdk.Int
		state             string
		expected          string
		functionAvailable bool
	}{
		// Power
		{PowerFunction, functionParametersPower(), multitokenReserve(),
			sdk.NewInt(0), OpenState, "100", true},
		{PowerFunction, functionParametersPower(), multitokenReserve(),
			sdk.NewInt(1000), OpenState, "12000100", true},
		// Sigmoid
		{SigmoidFunction, functionParametersSigmoid(), multitokenReserve(),
			sdk.NewInt(1000), OpenState, "5.999998484887893066", true},
		// Augmented
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			sdk.NewInt(0), HatchState, "0.01", true},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			sdk.NewInt(0), OpenState, "0", true},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			S0.TruncateInt(), HatchState, "0.01", true}, // p0=0.01
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			S0.TruncateInt(), OpenState, "0.018", true},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			S0.MulInt64(2).TruncateInt(), HatchState, "0.01", true}, // p0=0.01
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			S0.MulInt64(2).TruncateInt(), OpenState, "0.072", true},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			sdk.NewInt(10000000), OpenState, "720", true},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			sdk.NewInt(12345678), OpenState, "1097.3935100137248", true},
		// Swapper
		{SwapperFunction, nil, swapperReserves(),
			sdk.NewInt(100), OpenState, "100", false},
	}
	for _, tc := range testCases {
		bond.FunctionType = tc.functionType
		bond.FunctionParameters = tc.functionParams
		bond.ReserveTokens = tc.reserveTokens
		bond.State = tc.state

		actualResult, err := bond.GetPricesAtSupply(tc.supply)
		if tc.functionAvailable {
			require.Nil(t, err)
			expectedDec := sdk.MustNewDecFromStr(tc.expected)
			expectedResult := newDecMultitokenReserveFromDec(expectedDec).Add(nil)
			// __.Add(nil) is added so that zeroes are detected and removed
			// For example "0.00000res,0.00000rez" becomes ""
			require.Equal(t, expectedResult, actualResult)
		} else {
			require.Error(t, err)
		}
	}
}

func TestGetCurrentPrices(t *testing.T) {
	bond := getValidBond()
	// TODO: add more test cases

	swapperReserveBalances := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)

	augmentedP0 := functionParametersAugmented().AsMap()["p0"].String()

	testCases := []struct {
		functionType    string
		functionParams  FunctionParams
		reserveTokens   []string
		currentSupply   sdk.Int
		reserveBalances sdk.Coins
		state           string
		expected        string
	}{
		// Power
		{PowerFunction, functionParametersPower(), multitokenReserve(),
			sdk.NewInt(100), nil, OpenState, "120100"},
		// Sigmoid
		{SigmoidFunction, functionParametersSigmoid(), multitokenReserve(),
			sdk.NewInt(100), nil, OpenState, "5.999833808824623900"},
		// Augmented
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			sdk.NewInt(12345678), nil, HatchState, augmentedP0},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			sdk.NewInt(12345678), nil, OpenState, "1097.3935100137248"},
		// Swapper
		{SwapperFunction, nil, swapperReserves(),
			sdk.NewInt(100), swapperReserveBalances, OpenState, "100"},
	}
	for _, tc := range testCases {
		bond.FunctionType = tc.functionType
		bond.FunctionParameters = tc.functionParams
		bond.ReserveTokens = tc.reserveTokens
		bond.State = tc.state
		bond.CurrentSupply = sdk.NewCoin(bond.Token, tc.currentSupply)

		actualResult, _ := bond.GetCurrentPricesPT(tc.reserveBalances)
		expectedDec := sdk.MustNewDecFromStr(tc.expected)
		expectedResult := newDecMultitokenReserveFromDec(expectedDec)
		require.Equal(t, expectedResult, actualResult)
	}
}

func TestReserveAtSupply(t *testing.T) {
	bond := getValidBond()

	testCases := []struct {
		functionType   string
		functionParams FunctionParams
		supply         sdk.Int
		expected       string
	}{
		// Power
		{PowerFunction, functionParametersPower(), sdk.NewInt(100),
			"4010000"},
		{PowerFunction, functionParametersPower(), maxInt64,
			"3138550867693340380897047610841017818694071568064447512472.0"},
		{PowerFunction, functionParametersPowerHuge(), sdk.NewInt(5),
			"390525200604461289807786418456824866174854670846050992460534124091120.049504950495049505"},
		////{PowerFunction, functionParametersPowerHuge, sdk.NewInt(6),
		////	""}, // causes integer overflow
		// Sigmoid
		{SigmoidFunction, functionParametersSigmoid(), sdk.NewInt(100),
			"569.718730495548543525"},
		{SigmoidFunction, functionParametersSigmoid(), maxInt64,
			"55340232221128654811.702941459221645510"},
		{SigmoidFunction, functionParametersSigmoidHuge(), sdk.NewInt(1),
			"13043817825332782212.764456919596679543"},
		{SigmoidFunction, functionParametersSigmoidHuge(), maxInt64,
			"170141183460469231685570443531610226691.0"},
		// Augmented
		{AugmentedFunction, functionParametersAugmentedFull(), sdk.NewInt(1),
			"0.0000000000024"},
		{AugmentedFunction, functionParametersAugmentedFull(), sdk.NewInt(50000),
			"300"},
		{AugmentedFunction, functionParametersAugmentedFull(), maxInt64,
			"1883130520616004228538228566503104186246547815.244551376209997517"},
	}
	for _, tc := range testCases {
		bond.FunctionType = tc.functionType
		bond.FunctionParameters = tc.functionParams

		actualResult := bond.ReserveAtSupply(tc.supply)
		expectedResult := sdk.MustNewDecFromStr(tc.expected)
		require.Equal(t, expectedResult, actualResult)
	}
}

func TestGetReserveDeltaForLiquidityDelta(t *testing.T) {
	bond := getValidBond()
	bond.FunctionType = SwapperFunction
	bond.ReserveTokens = swapperReserves()
	// TODO: add more test cases

	reserveBalances := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)

	testCases := []struct {
		currentSupply  sdk.Int
		liquidityDelta sdk.Int
	}{
		{sdk.NewInt(2), sdk.NewInt(10)},
	}
	for _, tc := range testCases {
		bond.CurrentSupply = sdk.NewCoin(bond.Token, tc.currentSupply)

		actualResult := bond.GetReserveDeltaForLiquidityDelta(
			tc.liquidityDelta, reserveBalances)
		expectedResult := newDecMultitokenReserveFromInt(50000)
		require.Equal(t, expectedResult, actualResult)
	}
}

func TestGetPricesToMint(t *testing.T) {
	bond := getValidBond()
	// TODO: add more test cases

	tenK := int64(10000)
	reserveBalances10000 := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, tenK),
		sdk.NewInt64Coin(reserveToken2, tenK),
	)
	reserveBalances10 := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10),
		sdk.NewInt64Coin(reserveToken2, 10),
	)

	// For augmented function tests
	baseMap := functionParametersAugmented().AsMap()
	R0 := baseMap["d0"].Mul(sdk.OneDec().Sub(baseMap["theta"]))
	S0 := baseMap["d0"].Quo(baseMap["p0"])
	kappa := baseMap["kappa"].TruncateInt64()
	V0 := Invariant(R0, S0, kappa)
	augmentedSupplyForReserve10000 :=
		Supply(sdk.NewDec(tenK), kappa, V0).Ceil().TruncateInt()

	testCases := []struct {
		functionType    string
		functionParams  FunctionParams
		reserveTokens   []string
		reserveBalances sdk.Coins
		currentSupply   sdk.Int
		amount          sdk.Int
		state           string
		expectedPrice   string
		fails           bool
	}{
		// Power
		{PowerFunction, functionParametersPower(), multitokenReserve(),
			reserveBalances10000, sdk.ZeroInt(), sdk.NewInt(100), OpenState, "4000000", false},
		{PowerFunction, functionParametersPower(), multitokenReserve(),
			nil, sdk.ZeroInt(), sdk.NewInt(100), OpenState, "4010000", false},
		// Sigmoid
		{SigmoidFunction, functionParametersSigmoid(), multitokenReserve(),
			nil, sdk.ZeroInt(), sdk.NewInt(100), OpenState, "569.718730495548543525", false},
		{SigmoidFunction, functionParametersSigmoid(), multitokenReserve(),
			reserveBalances10, sdk.ZeroInt(), sdk.NewInt(100), OpenState, "559.718730495548543525", false},
		// Augmented
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			nil, sdk.ZeroInt(), sdk.NewInt(5000), HatchState, "50", false}, // p0=0.01; 0.01*5000 = 50
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			reserveBalances10000, augmentedSupplyForReserve10000, sdk.NewInt(5000), HatchState, "50", false}, // p0=0.01; 0.01*5000 = 50
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			nil, sdk.ZeroInt(), sdk.NewInt(5000), OpenState, "0.3", false},
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			reserveBalances10000, augmentedSupplyForReserve10000, sdk.NewInt(5000), OpenState, "961.4547618461", false},
		// Swapper
		{SwapperFunction, FunctionParams{}, swapperReserves(),
			reserveBalances10000, sdk.NewInt(2), sdk.NewInt(10), OpenState, "50000", false},
		{SwapperFunction, FunctionParams{}, swapperReserves(),
			nil, sdk.NewInt(2), sdk.NewInt(10), OpenState, "0", false}, // impossible scenario
		{SwapperFunction, FunctionParams{}, swapperReserves(),
			nil, sdk.ZeroInt(), sdk.NewInt(10), OpenState, "0", true},
	}
	for _, tc := range testCases {
		bond.FunctionType = tc.functionType
		bond.FunctionParameters = tc.functionParams
		bond.ReserveTokens = tc.reserveTokens
		bond.State = tc.state
		bond.CurrentSupply = sdk.NewCoin(bond.Token, tc.currentSupply)

		actualResult, err := bond.GetPricesToMint(tc.amount, tc.reserveBalances)
		if tc.fails {
			require.Error(t, err)
		} else {
			require.Nil(t, err)
			expectedDec := sdk.MustNewDecFromStr(tc.expectedPrice)
			expectedResult := newDecMultitokenReserveFromDec(expectedDec)
			require.Equal(t, expectedResult, actualResult)
		}
	}
}

func TestGetReturnsForBurn(t *testing.T) {
	bond := getValidBond()
	// TODO: add more test cases

	reserveBalances232 := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 232),
		sdk.NewInt64Coin(reserveToken2, 232),
	)

	tenK := int64(10000)
	reserveBalances10000 := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, tenK),
		sdk.NewInt64Coin(reserveToken2, tenK),
	)

	swapperReserveBalances := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)

	// For augmented function tests
	baseMap := functionParametersAugmented().AsMap()
	R0 := baseMap["d0"].Mul(sdk.OneDec().Sub(baseMap["theta"]))
	S0 := baseMap["d0"].Quo(baseMap["p0"])
	kappa := baseMap["kappa"].TruncateInt64()
	V0 := Invariant(R0, S0, kappa)
	augmentedSupplyForReserve10000 :=
		Supply(sdk.NewDec(tenK), kappa, V0).Ceil().TruncateInt()

	testCases := []struct {
		functionType    string
		functionParams  FunctionParams
		reserveTokens   []string
		reserveBalances sdk.Coins
		currentSupply   sdk.Int
		amount          sdk.Int
		expectedReturn  string
	}{
		// Power
		{PowerFunction, functionParametersPower(), multitokenReserve(),
			reserveBalances232, sdk.NewInt(2), sdk.OneInt(), "128"},
		// Sigmoid
		{SigmoidFunction, functionParametersSigmoid(), multitokenReserve(),
			reserveBalances232, sdk.NewInt(2), sdk.OneInt(), "231.927741663925372840"},
		// Augmented (note: unlike in minting, state not taken into consideration when
		// burning since burning only possible in open phase, so state cannot be hatch)
		{AugmentedFunction, functionParametersAugmentedFull(), multitokenReserve(),
			reserveBalances10000, augmentedSupplyForReserve10000, sdk.NewInt(5000), "903.4871183539"},
		// Swapper
		{SwapperFunction, FunctionParams{}, swapperReserves(),
			swapperReserveBalances, sdk.NewInt(2), sdk.OneInt(), "5000"},
	}
	for _, tc := range testCases {
		bond.FunctionType = tc.functionType
		bond.FunctionParameters = tc.functionParams
		bond.ReserveTokens = tc.reserveTokens
		bond.CurrentSupply = sdk.NewCoin(bond.Token, tc.currentSupply)

		actualResult := bond.GetReturnsForBurn(tc.amount, tc.reserveBalances)
		expectedDec := sdk.MustNewDecFromStr(tc.expectedReturn)
		expectedResult := newDecMultitokenReserveFromDec(expectedDec)
		require.Equal(t, expectedResult, actualResult)
	}
}

func TestGetReturnsForSwap(t *testing.T) {
	bond := getValidBond()
	bond.FunctionType = SwapperFunction
	bond.FunctionParameters = nil
	bond.ReserveTokens = swapperReserves()

	reserveBalances := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)

	zeroPoint1Percent := sdk.MustNewDecFromStr("0.001")
	largeInput := maxInt64
	largeFee := largeInput.ToDec().Mul(
		zeroPoint1Percent).Ceil().TruncateInt()
	smallInput := sdk.NewInt(3) // but not too small
	smallFee := smallInput.ToDec().Mul(
		zeroPoint1Percent).Ceil().TruncateInt()

	testCases := []struct {
		bondTxFee           string
		from                string
		to                  string
		amount              sdk.Int
		expectedReturn      sdk.Int
		expectedFee         sdk.Int
		amountInvalid       bool // too large or too small
		invalidReserveToken bool
	}{
		{"0.1", reserveToken, reserveToken2, smallInput, sdk.OneInt(),
			smallFee, false, false},
		{"0.1", reserveToken, reserveToken2, sdk.NewInt(2), sdk.OneInt(),
			sdk.OneInt(), true, false},
		{"0.1", reserveToken, reserveToken2, sdk.NewInt(1), sdk.OneInt(),
			sdk.OneInt(), true, false},
		{"0.1", reserveToken, reserveToken2, sdk.NewInt(0), sdk.OneInt(),
			sdk.OneInt(), true, false},
		{"0.1", reserveToken, reserveToken2, largeInput, sdk.NewInt(9999),
			largeFee, false, false},
		{"0.1", reserveToken, "dummytoken", sdk.NewInt(3), sdk.OneInt(),
			sdk.OneInt(), false, true}, // identical to first case but dummytoken
		{"0.1", "dummytoken", reserveToken2, sdk.NewInt(3), sdk.OneInt(),
			sdk.OneInt(), false, true}, // identical to first case but dummytoken
	}
	for _, tc := range testCases {
		bond.TxFeePercentage = sdk.MustNewDecFromStr(tc.bondTxFee)
		fromAmount := sdk.NewCoin(tc.from, tc.amount)
		actualResult, actualFee, err := bond.GetReturnsForSwap(
			fromAmount, tc.to, reserveBalances)
		if tc.amountInvalid {
			require.Error(t, err)
			require.Equal(t, err.Code(), CodeSwapAmountInvalid)
		} else if tc.invalidReserveToken {
			require.Error(t, err)
			require.Equal(t, err.Code(), CodeReserveTokenInvalid)
		} else {
			require.Nil(t, err)
			expectedResult := sdk.NewCoins(sdk.NewCoin(tc.to, tc.expectedReturn))
			expectedFee := sdk.NewCoin(tc.from, tc.expectedFee)
			require.Equal(t, expectedResult, actualResult)
			require.Equal(t, expectedFee, actualFee)
		}
	}
}

func TestGetReturnsForSwapNonSwapperFunctionFails(t *testing.T) {
	bond := getValidBond()
	testCases := []string{PowerFunction, SigmoidFunction}

	for _, tc := range testCases {
		bond.FunctionType = tc

		dummyCoin := sdk.NewCoin(reserveToken, sdk.OneInt()) // to avoid panic

		_, _, err := bond.GetReturnsForSwap(dummyCoin, "", sdk.Coins{})
		require.Error(t, err)
		require.False(t, err.Result().IsOK())
		require.Equal(t, err.Code(), CodeFunctionNotAvailableForFunctionType)
	}
}

func TestBondGetTxFee(t *testing.T) {
	bond := Bond{}
	zeroPointOne := sdk.MustNewDecFromStr("0.1")

	// Fee is always rounded to ceiling, so for any input N > 0, fee(N) > 0

	testCases := []struct {
		input           string
		txFeePercentage sdk.Dec
		expected        int64
	}{

		{"2000000000000", zeroPointOne, 2000000000},
		{"2000", zeroPointOne, 2},
		{"200", zeroPointOne, 1},      // 200 * 0.1 = 0.2 = 1 (rounded)
		{"20", zeroPointOne, 1},       // 20 * 0.1 = 00.2 = 1 (rounded)
		{"0.000002", zeroPointOne, 1}, // 0.000002 * 0.1 = small number = 1 (rounded)
		{"0", zeroPointOne, 0},
		{"2000", sdk.ZeroDec(), 0},
		{"0.000002", sdk.ZeroDec(), 0},
	}
	for _, tc := range testCases {
		inputToken := sdk.NewDecCoinFromDec(reserveToken,
			sdk.MustNewDecFromStr(tc.input))
		expected := sdk.NewInt64Coin(reserveToken, tc.expected)

		bond.TxFeePercentage = tc.txFeePercentage
		require.Equal(t, expected, bond.GetTxFee(inputToken))
	}
}

func TestBondGetExitFee(t *testing.T) {
	bond := Bond{}
	zeroPointOne := sdk.MustNewDecFromStr("0.1")

	// Fee is always rounded to ceiling, so for any input N > 0, fee(N) > 0

	testCases := []struct {
		input             string
		exitFeePercentage sdk.Dec
		expected          int64
	}{
		{"2000000000000", zeroPointOne, 2000000000},
		{"2000", zeroPointOne, 2},
		{"200", zeroPointOne, 1},      // 200 * 0.1 = 0.2 = 1 (rounded)
		{"20", zeroPointOne, 1},       // 20 * 0.1 = 00.2 = 1 (rounded)
		{"0.000002", zeroPointOne, 1}, // 0.000002 * 0.1 = small number = 1 (rounded)
		{"0", zeroPointOne, 0},
		{"2000", sdk.ZeroDec(), 0},
		{"0.000002", sdk.ZeroDec(), 0},
	}
	for _, tc := range testCases {
		inputToken := sdk.NewDecCoinFromDec(reserveToken,
			sdk.MustNewDecFromStr(tc.input))
		expected := sdk.NewInt64Coin(reserveToken, tc.expected)

		bond.ExitFeePercentage = tc.exitFeePercentage
		require.Equal(t, expected, bond.GetExitFee(inputToken))
	}
}

func TestBondGetTxFees(t *testing.T) {
	bond := Bond{}
	bond.TxFeePercentage = sdk.MustNewDecFromStr("0.1")

	// Fee is always rounded to ceiling, so for any input N > 0, fee(N) > 0

	inputTokens, err := sdk.ParseDecCoins("" +
		"200000000.0aaa," +
		"2000.0bbb," +
		"200.0ccc," +
		"20.0ddd," +
		"0.000002eee")
	require.Nil(t, err)

	expected, err := sdk.ParseCoins("" +
		"200000aaa," +
		"2bbb," +
		"1ccc," +
		"1ddd," +
		"1eee")
	require.Nil(t, err)

	require.Equal(t, expected, bond.GetTxFees(inputTokens))
}

func TestBondGetExitFees(t *testing.T) {
	bond := Bond{}
	bond.ExitFeePercentage = sdk.MustNewDecFromStr("0.1")

	// Fee is always rounded to ceiling, so for any input N > 0, fee(N) > 0

	inputTokens, err := sdk.ParseDecCoins("" +
		"200000000.0aaa," +
		"2000.0bbb," +
		"200.0ccc," +
		"20.0ddd," +
		"0.000002eee")
	require.Nil(t, err)

	expected, err := sdk.ParseCoins("" +
		"200000aaa," +
		"2bbb," +
		"1ccc," +
		"1ddd," +
		"1eee")
	require.Nil(t, err)

	require.Equal(t, expected, bond.GetExitFees(inputTokens))
}

func TestSignersEqualTo(t *testing.T) {
	bond := getValidBond()

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	bond.Signers = []sdk.AccAddress{addr1, addr2}

	testCases := []struct {
		toCompareTo   []sdk.AccAddress
		expectedEqual bool
	}{
		{[]sdk.AccAddress{addr1}, false},               // One missing
		{[]sdk.AccAddress{addr1, addr2, addr3}, false}, // One extra
		{[]sdk.AccAddress{addr1, addr3}, false},        // One different
		{[]sdk.AccAddress{addr2, addr1}, false},        // Incorrect order
		{[]sdk.AccAddress{addr1, addr2}, true},         // Equal
	}
	for _, tc := range testCases {
		require.Equal(t, tc.expectedEqual, bond.SignersEqualTo(tc.toCompareTo))
	}
}

func TestReserveDenomsEqualTo(t *testing.T) {
	bond := getValidBond()

	denom1 := reserveToken
	denom2 := reserveToken2
	denom3 := reserveToken3
	bond.ReserveTokens = []string{denom1, denom2}

	testCases := []struct {
		toCompareTo   []string
		expectedEqual bool
	}{
		{[]string{denom1}, false},                 // One missing
		{[]string{denom1, denom2, denom3}, false}, // One extra
		{[]string{denom1, denom3}, false},         // One different
		{[]string{denom2, denom1}, true},          // Incorrect order (allowed)
		{[]string{denom1, denom2}, true},          // Equal
	}
	for _, tc := range testCases {
		coins := sdk.Coins{}
		for _, res := range tc.toCompareTo {
			coins = coins.Add(sdk.Coins{sdk.NewCoin(res, sdk.OneInt())})
		}
		require.Equal(t, tc.expectedEqual, bond.ReserveDenomsEqualTo(coins))
	}
}

func TestAnyOrderQuantityLimitsExceeded(t *testing.T) {
	bond := getValidBond()
	bond.OrderQuantityLimits, _ = sdk.ParseCoins("100aaa,200bbb")

	testCases := []struct {
		amounts         string
		exceedsAnyLimit bool
	}{
		{"99aaa", false},         // aaa <= 100
		{"100aaa", false},        // aaa <= 100
		{"101aaa", true},         // aaa >  100
		{"101bbb", false},        // bbb <= 200
		{"100aaa,200bbb", false}, // aaa <= 100, bbb <= 200
		{"101aaa,200bbb", true},  // aaa >  100, bbb <= 200
		{"100aaa,201bbb", true},  // aaa <= 100, bbb >  200
		{"101aaa,201bbb", true},  // aaa >  100, bbb >  200
	}
	for _, tc := range testCases {
		amounts, _ := sdk.ParseCoins(tc.amounts)
		require.Equal(t, tc.exceedsAnyLimit,
			bond.AnyOrderQuantityLimitsExceeded(amounts))
	}
}

func TestReservesViolateSanityRateReturnsFalseWhenSanityRateIsZero(t *testing.T) {
	bond := getValidBond()

	r1 := reserveToken
	r2 := reserveToken2
	bond.ReserveTokens = []string{r1, r2}

	testCases := []struct {
		reserves               string
		sanityRate             string
		sanityMarginPercentage string
		violates               bool
	}{
		{fmt.Sprintf(" 500%s,1000%s", r1, r2),
			"0", "0", false}, // no sanity checks
		{fmt.Sprintf("1000%s,1000%s", r1, r2),
			"0", "0", false}, // no sanity checks
		{fmt.Sprintf(" 500%s,1000%s", r1, r2),
			"0.5", "0", false}, //  500/1000 == 0.5
		{fmt.Sprintf("1000%s,1000%s", r1, r2),
			"0.5", "0", true}, // 1000/1000 != 0.5
		{fmt.Sprintf(" 100%s,1000%s", r1, r2),
			"0.5", "0", true}, //  100/1000 != 0.5
		{fmt.Sprintf(" 100%s,1000%s", r1, r2),
			"0.5", "79", true}, // 0.5+-79% => 0.105 to 0.895, and 100/1000 is in not this range
		{fmt.Sprintf(" 100%s,1000%s", r1, r2),
			"0.5", "80", false}, // 0.5+-80% => 0.100 to 0.900, and 100/1000 is in this range
		{fmt.Sprintf(" 100%s,1000%s", r1, r2),
			"0.5", "81", false}, // 0.5+-81% => 0.095 to 0.905, and 100/1000 is in this range
		{fmt.Sprintf(" 100%s,1000%s", r1, r2),
			"0.5", "101", false}, // identical to above but negative lower limit gets rounded to 0
	}
	for _, tc := range testCases {
		reserves, _ := sdk.ParseCoins(tc.reserves)
		srDec := sdk.MustNewDecFromStr(tc.sanityRate)
		smpDec := sdk.MustNewDecFromStr(tc.sanityMarginPercentage)

		bond.SanityRate = srDec
		bond.SanityMarginPercentage = smpDec

		actualResult := bond.ReservesViolateSanityRate(reserves)
		require.Equal(t, tc.violates, actualResult)
	}
}
