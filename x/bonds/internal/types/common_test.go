package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	token = "testtoken"

	blankSanityRate             = "0"
	blankSanityMarginPercentage = "0"
	reserveToken                = "res"
	reserveToken2               = "rez"
	reserveToken3               = "rec"

	initToken                  = token
	initName                   = "test token"
	initDescription            = "this is a test token"
	initCreator                = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	initFeeAddress             = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	initTxFeePercentage        = sdk.MustNewDecFromStr("0.1")
	initExitFeePercentage      = sdk.MustNewDecFromStr("0.1")
	initMaxSupply              = sdk.NewInt64Coin(initToken, 10000)
	initOrderQuantityLimits    = sdk.Coins(nil)
	initSanityRate             = sdk.MustNewDecFromStr(blankSanityRate)
	initSanityMarginPercentage = sdk.MustNewDecFromStr(blankSanityMarginPercentage)
	initAllowSell              = true
	initSigners                = []sdk.AccAddress{initCreator}
	initBatchBlocks            = sdk.NewUint(10)
	initOutcomePayment         = sdk.Coins(nil)
	initState                  = OpenState

	// 9223372036854775807
	maxInt64 = sdk.NewInt(int64(^uint64(0) >> 1))
)

func functionParametersPower() FunctionParams {
	return FunctionParams{
		NewFunctionParam("m", sdk.NewDec(12)),
		NewFunctionParam("n", sdk.NewDec(2)),
		NewFunctionParam("c", sdk.NewDec(100))}
}

func functionParametersPower2() FunctionParams {
	return FunctionParams{
		NewFunctionParam("m", sdk.NewDec(12)),
		NewFunctionParam("n", sdk.NewDecWithPrec(25, 1)),
		NewFunctionParam("c", sdk.NewDec(100))}
}

func functionParametersSigmoid() FunctionParams {
	return FunctionParams{
		NewFunctionParam("a", sdk.NewDec(3)),
		NewFunctionParam("b", sdk.NewDec(5)),
		NewFunctionParam("c", sdk.NewDec(1))}
}

func functionParametersAugmented() FunctionParams {
	return FunctionParams{
		NewFunctionParam("d0", sdk.MustNewDecFromStr("500.0")),
		NewFunctionParam("p0", sdk.MustNewDecFromStr("0.01")),
		NewFunctionParam("theta", sdk.MustNewDecFromStr("0.4")),
		NewFunctionParam("kappa", sdk.MustNewDecFromStr("3.0"))}
}

func functionParametersAugmentedFull() FunctionParams {
	base := functionParametersAugmented()
	baseMap := base.AsMap()

	R0 := baseMap["d0"].Mul(sdk.OneDec().Sub(baseMap["theta"]))
	S0 := baseMap["d0"].Quo(baseMap["p0"])
	V0 := Invariant(R0, S0, baseMap["kappa"])
	extras := FunctionParams{
		NewFunctionParam("R0", R0),
		NewFunctionParam("S0", S0),
		NewFunctionParam("V0", V0)}

	return append(base, extras...)
}

func functionParametersPowerHuge() FunctionParams {
	return FunctionParams{
		NewFunctionParam("m", sdk.NewDec(1)),
		NewFunctionParam("n", sdk.NewDec(100)),
		NewFunctionParam("c", sdk.NewDec(0))}
}

func functionParametersSigmoidHuge() FunctionParams {
	return FunctionParams{
		NewFunctionParam("a", sdk.NewDec(int64(^uint64(0)>>1))),
		NewFunctionParam("b", sdk.NewDec(0)),
		NewFunctionParam("c", sdk.NewDec(1))}
}

func powerReserves() []string     { return []string{reserveToken} }
func multitokenReserve() []string { return []string{reserveToken, reserveToken2} }
func swapperReserves() []string   { return []string{reserveToken, reserveToken2} }

func getValidPowerFunctionBond() Bond {
	functionType := PowerFunction
	functionParams := functionParametersPower()
	reserveTokens := powerReserves()
	return NewBond(initToken, initName, initDescription, initCreator,
		functionType, functionParams, reserveTokens, initTxFeePercentage,
		initExitFeePercentage, initFeeAddress, initMaxSupply,
		initOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks, initOutcomePayment, initState)
}

func getValidBond() Bond {
	return getValidPowerFunctionBond()
}

// New Reserve Balances

func newDecMultitokenReserveFromDec(value sdk.Dec) sdk.DecCoins {
	return sdk.DecCoins{
		sdk.NewDecCoinFromDec(reserveToken, value),
		sdk.NewDecCoinFromDec(reserveToken2, value),
	}.Sort()
}

func newDecMultitokenReserveFromInt(value int64) sdk.DecCoins {
	return sdk.NewDecCoinsFromCoins(
		sdk.NewInt64Coin(reserveToken, value),
		sdk.NewInt64Coin(reserveToken2, value),
	)
}

// Messages

func newValidMsgCreateBond() MsgCreateBond {
	functionType := PowerFunction
	functionParams := functionParametersPower()
	reserveTokens := powerReserves()
	return NewMsgCreateBond(initToken, initName, initDescription, initCreator,
		functionType, functionParams, reserveTokens, initTxFeePercentage,
		initExitFeePercentage, initFeeAddress, initMaxSupply,
		initOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks, initOutcomePayment)
}

func newValidMsgCreateSwapperBond() MsgCreateBond {
	validMsg := newValidMsgCreateBond()
	validMsg.FunctionType = SwapperFunction
	validMsg.FunctionParameters = nil
	validMsg.ReserveTokens = swapperReserves()
	return validMsg
}

func newEmptyStringsMsgEditBond() MsgEditBond {
	return NewMsgEditBond(initToken, "", "", "", "", "",
		initCreator, initSigners)
}

func newValidMsgEditBond() MsgEditBond {
	return NewMsgEditBond(initToken, "newName", "newDescription", "", "0", "0",
		initCreator, initSigners)
}

func newValidMsgBuy() MsgBuy {
	buyer := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount, _ := sdk.ParseCoin("10" + initToken)
	maxPrices, _ := sdk.ParseCoins("50" + initToken)
	return NewMsgBuy(buyer, amount, maxPrices)
}

func newValidMsgSell() MsgSell {
	seller := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount, _ := sdk.ParseCoin("10" + initToken)
	return NewMsgSell(seller, amount)
}

func newValidMsgSwap() MsgSwap {
	swapper := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	from := sdk.NewInt64Coin(reserveToken, 10)
	return NewMsgSwap(swapper, initToken, from, reserveToken2)
}
