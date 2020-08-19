package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestPowerDecDec(t *testing.T) {
	testCases := []struct {
		inX string
		inY uint64
		out string
	}{{"0", 0, "1"}, {"0.5", 2, "0.25"}}
	for _, tc := range testCases {
		inXDec := sdk.MustNewDecFromStr(tc.inX)
		outDec := sdk.MustNewDecFromStr(tc.out)
		require.Equal(t, outDec, Power(inXDec, tc.inY))
	}
}

func TestRootDec(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{{"0", "0"}, {"0.25", "0.5"}, {"3", "1.732050807568877294"}, {"9", "3"}}
	for _, tc := range testCases {
		inDec := sdk.MustNewDecFromStr(tc.in)
		outDec := sdk.MustNewDecFromStr(tc.out)
		result, err := ApproxRoot(inDec, 2)
		require.Nil(t, err)
		require.Equal(t, outDec, result)
	}
}

func TestRoundReservePrice(t *testing.T) {
	token := "token"

	// In general, RoundReservePrice rounds up

	testCases := []struct {
		in  string
		out int64
	}{{"9", 9}, {"1.6", 2}, {"0.5", 1}, {"0.4", 1}, {"0", 0}}
	for _, tc := range testCases {
		inDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.in))
		outInt := sdk.NewCoin(token, sdk.NewInt(tc.out))
		require.Equal(t, outInt, RoundReservePrice(inDec))
	}
}

func TestRoundReservePricesRoundsAllValues(t *testing.T) {
	tokens := []string{"token1", "token2", "token3"}
	ins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokens[0], sdk.MustNewDecFromStr("0.4")),
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("1.6")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("3")),
	}
	outs := sdk.Coins{
		sdk.NewInt64Coin(tokens[0], 1),
		sdk.NewInt64Coin(tokens[1], 2),
		sdk.NewInt64Coin(tokens[2], 3),
	}
	require.True(t, RoundReservePrices(ins).IsEqual(outs))
}

func TestRoundReserveReturn(t *testing.T) {
	token := "token"

	// In general, RoundReserveReturn rounds down

	testCases := []struct {
		in  string
		out int64
	}{{"5", 5}, {"1.4", 1}, {"1.9", 1}, {"0.5", 0}, {"0", 0}}
	for _, tc := range testCases {
		inDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.in))
		outInt := sdk.NewCoin(token, sdk.NewInt(tc.out))
		require.True(t, outInt.IsEqual(RoundReserveReturn(inDec)))
	}
}

func TestRoundReserveReturnsRoundsAllValues(t *testing.T) {
	tokens := []string{"token1", "token2", "token3"}
	ins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokens[0], sdk.MustNewDecFromStr("0.4")),
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("1.6")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("3")),
	}
	outs := sdk.Coins{
		// 0token1
		sdk.NewInt64Coin(tokens[1], 1),
		sdk.NewInt64Coin(tokens[2], 3),
	}
	require.True(t, RoundReserveReturns(ins).IsEqual(outs))
}

func TestMultiplyDecCoinByDec(t *testing.T) {
	token := "token"
	testCases := []struct {
		inCoin string
		scale  string
		out    string
	}{
		{"2", "2", "4"},        // all integers
		{"0.5", "2", "1"},      // result is integer
		{"2", "0.5", "1"},      // result is integer
		{"1.5", "0.5", "0.75"}, // all numbers decimal
		{"5", "1", "5"},        // N x 1 = N
		{"5", "0", "0"},        // N x 0 = 0
	}
	for _, tc := range testCases {
		inDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.inCoin))
		scaleDec := sdk.MustNewDecFromStr(tc.scale)
		outDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.out))
		require.True(t, outDec.IsEqual(MultiplyDecCoinByDec(inDec, scaleDec)))
	}
}

func TestMultiplyDecCoinsByDecMultipliesAllValues(t *testing.T) {
	tokens := []string{"token1", "token2", "token3", "token4"}
	ins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokens[0], sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("0.5")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("1")),
		sdk.NewDecCoinFromDec(tokens[3], sdk.MustNewDecFromStr("2")),
	}.Sort()
	scaleDec := sdk.MustNewDecFromStr("0.5")
	outs := sdk.DecCoins{
		// 0token1
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("0.25")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("0.5")),
		sdk.NewDecCoinFromDec(tokens[3], sdk.MustNewDecFromStr("1")),
	}.Sort()
	require.True(t, MultiplyDecCoinsByDec(ins, scaleDec).IsEqual(outs))
}

func TestMultiplyDecCoinByInt(t *testing.T) {
	token := "token"
	testCases := []struct {
		inCoin string
		scale  int64
		out    string
	}{
		{"2", 2, "4"},   // all integers
		{"0.5", 2, "1"}, // result is integer
		{"5", 1, "5"},   // N x 1 = N
		{"5", 0, "0"},   // N x 0 = 0
	}
	for _, tc := range testCases {
		inDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.inCoin))
		scaleInt := sdk.NewInt(tc.scale)
		outDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.out))
		require.True(t, outDec.IsEqual(MultiplyDecCoinByInt(inDec, scaleInt)))
	}
}

func TestMultiplyDecCoinsByIntMultipliesAllValues(t *testing.T) {
	tokens := []string{"token1", "token2", "token3", "token4"}
	ins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokens[0], sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("0.25")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("0.5")),
		sdk.NewDecCoinFromDec(tokens[3], sdk.MustNewDecFromStr("1")),
	}.Sort()
	scaleInt := sdk.NewInt(2)
	outs := sdk.DecCoins{
		// 0token1
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("0.5")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("1")),
		sdk.NewDecCoinFromDec(tokens[3], sdk.MustNewDecFromStr("2")),
	}.Sort()
	require.True(t, MultiplyDecCoinsByInt(ins, scaleInt).IsEqual(outs))
}

func TestDivideDecCoinByDec(t *testing.T) {
	token := "token"
	testCases := []struct {
		inCoin string
		scale  string
		out    string
	}{
		{"4", "2", "2"},       // all integers
		{"1", "2", "0.5"},     // result is decimal
		{"1.5", "0.5", "3"},   // result is integer
		{"1.5", "2.5", "0.6"}, // all numbers decimal
		{"5", "5", "1"},       // N / N = 1
		//{"5", "0", "0"},       // N / 0 = ? (panics)
	}
	for _, tc := range testCases {
		inDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.inCoin))
		scaleDec := sdk.MustNewDecFromStr(tc.scale)
		outDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.out))
		require.True(t, outDec.IsEqual(DivideDecCoinByDec(inDec, scaleDec)))
	}
}

func TestMultiplyDecCoinsByDecDividesAllValues(t *testing.T) {
	tokens := []string{"token1", "token2", "token3", "token4"}
	ins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(tokens[0], sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("0.5")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("1")),
		sdk.NewDecCoinFromDec(tokens[3], sdk.MustNewDecFromStr("2")),
	}.Sort()
	scaleDec := sdk.MustNewDecFromStr("2")
	outs := sdk.DecCoins{
		// 0token1
		sdk.NewDecCoinFromDec(tokens[1], sdk.MustNewDecFromStr("0.25")),
		sdk.NewDecCoinFromDec(tokens[2], sdk.MustNewDecFromStr("0.5")),
		sdk.NewDecCoinFromDec(tokens[3], sdk.MustNewDecFromStr("1")),
	}.Sort()
	require.True(t, DivideDecCoinsByDec(ins, scaleDec).IsEqual(outs))
}

func TestRoundFee(t *testing.T) {
	token := "token"

	// In general, RoundFee rounds up

	testCases := []struct {
		in  string
		out int64
	}{{"7", 7}, {"0.4", 1}, {"67.7", 68}, {"96.5", 97}, {"0", 0}}
	for _, tc := range testCases {
		inDec := sdk.NewDecCoinFromDec(token, sdk.MustNewDecFromStr(tc.in))
		outInt := sdk.NewCoin(token, sdk.NewInt(tc.out))
		require.True(t, outInt.IsEqual(RoundFee(inDec)))
	}
}

func TestAdjustFees(t *testing.T) {
	bond := getValidBond()
	bond.ExitFeePercentage = sdk.MustNewDecFromStr("0.1")

	fees, err := sdk.ParseCoins("" +
		"2000bbb," +
		"2000ccc," +
		"2000ddd," +
		"2000eee")
	require.Nil(t, err)

	maxFees, err := sdk.ParseCoins("" +
		"2001bbb," + // greater
		"2000ccc," + // equal
		"1999ddd," + // less
		"2000fff") // extra max (token 'fff' not in fees)
	require.Nil(t, err)

	expected, err := sdk.ParseCoins("" +
		"2000bbb," + // max > value => no rounding
		"2000ccc," + // max = value => no rounding
		"1999ddd") // max < value => rounded
	// for token eee, no max => assumes max = 0 => rounded to 0
	// for token fff, no fees => not in adjusted fees
	require.Nil(t, err)

	require.Equal(t, expected, AdjustFees(fees, maxFees))
}

func TestStringsToString(t *testing.T) {
	testCases := []struct {
		in  []string
		out string
	}{
		{[]string{}, "[]"},
		{[]string{"str1"}, "[str1]"},
		{[]string{"str1", "str2", "str3"}, "[str1,str2,str3]"},
	}
	for _, tc := range testCases {
		require.Equal(t, tc.out, StringsToString(tc.in))
	}
}

func TestAccAddressesToString(t *testing.T) {

	oneIn := []sdk.AccAddress{
		sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
	}
	threeIn := []sdk.AccAddress{
		sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
	}

	testCases := []struct {
		in  []sdk.AccAddress
		out string
	}{
		{[]sdk.AccAddress{}, "[]"},
		{oneIn, fmt.Sprintf("[%s]", oneIn[0])},
		{threeIn, fmt.Sprintf("[%s,%s,%s]", threeIn[0], threeIn[1], threeIn[2])},
	}
	for _, tc := range testCases {
		require.Equal(t, tc.out, AccAddressesToString(tc.in))
	}
}

func TestPower(t *testing.T) {
	testCases := []struct {
		input    sdk.Dec
		power    uint64
		expected sdk.Dec
	}{
		{sdk.OneDec(), 10, sdk.OneDec()},                                        // 1.0 ^ (10) => 1.0
		{sdk.NewDecWithPrec(5, 1), 2, sdk.NewDecWithPrec(25, 2)},                // 0.5 ^ 2 => 0.25
		{sdk.NewDecWithPrec(2, 1), 2, sdk.NewDecWithPrec(4, 2)},                 // 0.2 ^ 2 => 0.04
		{sdk.NewInt(3).ToDec(), 3, sdk.NewInt(27).ToDec()},                      // 3 ^ 3 => 27
		{sdk.NewInt(-3).ToDec(), 4, sdk.NewInt(81).ToDec()},                     // -3 ^ 4 = 81
		{sdk.NewDecWithPrec(1414213562373095049, 18), 2, sdk.NewInt(2).ToDec()}, // 1.414213562373095049 ^ 2 = 2
	}

	for i, tc := range testCases {
		res := Power(tc.input, tc.power)
		require.True(t, tc.expected.Sub(res).Abs().LTE(sdk.SmallestDec()), "unexpected result for test case %d, input: %v", i, tc.input)
	}
}

func TestApproxRoot(t *testing.T) {
	testCases := []struct {
		input    sdk.Dec
		root     uint64
		expected sdk.Dec
	}{
		{sdk.OneDec(), 10, sdk.OneDec()},                                        // 1.0 ^ (0.1) => 1.0
		{sdk.NewDecWithPrec(25, 2), 2, sdk.NewDecWithPrec(5, 1)},                // 0.25 ^ (0.5) => 0.5
		{sdk.NewDecWithPrec(4, 2), 2, sdk.NewDecWithPrec(2, 1)},                 // 0.04 => 0.2
		{sdk.NewInt(27).ToDec(), 3, sdk.NewInt(3).ToDec()},                      // 27 ^ (1/3) => 3
		{sdk.NewInt(-81).ToDec(), 4, sdk.NewInt(-3).ToDec()},                    // -81 ^ (0.25) => -3
		{sdk.NewInt(2).ToDec(), 2, sdk.NewDecWithPrec(1414213562373095049, 18)}, // 2 ^ (0.5) => 1.414213562373095049
		{sdk.NewDecWithPrec(1005, 3), 31536000, sdk.MustNewDecFromStr("1.000000000158153904")},
	}

	for i, tc := range testCases {
		res, err := ApproxRoot(tc.input, tc.root)
		require.NoError(t, err)
		require.True(t, tc.expected.Sub(res).Abs().LTE(sdk.SmallestDec()), "unexpected result for test case %d, input: %v", i, tc.input)
	}
}
