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

	functionParametersPower = FunctionParams{
		NewFunctionParam("m", sdk.NewInt(12)),
		NewFunctionParam("n", sdk.NewInt(2)),
		NewFunctionParam("c", sdk.NewInt(100))}
	functionParametersSigmoid = FunctionParams{
		NewFunctionParam("a", sdk.NewInt(3)),
		NewFunctionParam("b", sdk.NewInt(5)),
		NewFunctionParam("c", sdk.NewInt(1))}

	functionParametersPowerHuge = FunctionParams{
		NewFunctionParam("m", sdk.NewInt(1)),
		NewFunctionParam("n", sdk.NewInt(100)),
		NewFunctionParam("c", sdk.NewInt(0))}
	functionParametersSigmoidHuge = FunctionParams{
		NewFunctionParam("a", sdk.NewInt(int64(^uint64(0)>>1))),
		NewFunctionParam("b", sdk.NewInt(0)),
		NewFunctionParam("c", sdk.NewInt(1))}

	powerReserves     = []string{reserveToken}
	multitokenReserve = []string{reserveToken, reserveToken2}
	swapperReserves   = []string{reserveToken, reserveToken2}

	initToken                  = token
	initName                   = "test token"
	initDescription            = "this is a test token"
	initCreator                = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	initReserveAddress         = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	initFeeAddress             = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	initTxFeePercentage        = sdk.MustNewDecFromStr("0.1")
	initExitFeePercentage      = sdk.MustNewDecFromStr("0.1")
	initMaxSupply              = sdk.NewInt64Coin(initToken, 10000)
	initOrderQuantityLimits    = sdk.Coins(nil)
	initSanityRate             = sdk.MustNewDecFromStr(blankSanityRate)
	initSanityMarginPercentage = sdk.MustNewDecFromStr(blankSanityMarginPercentage)
	initAllowSell              = "true"
	initSigners                = []sdk.AccAddress{initCreator}
	initBatchBlocks            = sdk.NewUint(10)

	maxInt64 = sdk.NewInt(int64(^uint64(0) >> 1))
)

func getValidPowerFunctionBond() Bond {
	functionType := PowerFunction
	functionParams := functionParametersPower
	reserveTokens := powerReserves
	return NewBond(initToken, initName, initDescription,
		initCreator, functionType, functionParams,
		reserveTokens, initReserveAddress, initTxFeePercentage,
		initExitFeePercentage, initFeeAddress, initMaxSupply,
		initOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks)
}

func getValidBond() Bond {
	return getValidPowerFunctionBond()
}

// New Reserve Balances

func NewDecMultitokenReserveFromDec(value sdk.Dec) sdk.DecCoins {
	return sdk.DecCoins{
		sdk.NewDecCoinFromDec(reserveToken, value),
		sdk.NewDecCoinFromDec(reserveToken2, value),
	}.Sort()
}

func NewDecMultitokenReserveFromInt(value int64) sdk.DecCoins {
	return sdk.NewDecCoins(sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, value),
		sdk.NewInt64Coin(reserveToken2, value),
	))
}

// Messages

func NewValidMsgCreateBond() MsgCreateBond {
	functionType := PowerFunction
	functionParams := functionParametersPower
	reserveTokens := powerReserves
	return NewMsgCreateBond(initToken, initName, initDescription,
		initCreator, functionType, functionParams,
		reserveTokens, initTxFeePercentage, initExitFeePercentage,
		initFeeAddress, initMaxSupply, initOrderQuantityLimits, initSanityRate,
		initSanityMarginPercentage, initAllowSell, initSigners, initBatchBlocks)
}

func NewEmptyStringsMsgEditBond() MsgEditBond {
	return NewMsgEditBond(initToken, "", "", "", "", "",
		initCreator, initSigners)
}

func NewValidMsgEditBond() MsgEditBond {
	return NewMsgEditBond(initToken, "newName", "newDescription", "", "0", "0",
		initCreator, initSigners)
}

func NewValidMsgBuy() MsgBuy {
	buyer := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount, _ := sdk.ParseCoin("10" + initToken)
	maxPrices, _ := sdk.ParseCoins("50" + initToken)
	return NewMsgBuy(buyer, amount, maxPrices)
}

func NewValidMsgSell() MsgSell {
	seller := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount, _ := sdk.ParseCoin("10" + initToken)
	return NewMsgSell(seller, amount)
}

func NewValidMsgSwap() MsgSwap {
	swapper := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	from := sdk.NewInt64Coin(reserveToken, 10)
	return NewMsgSwap(swapper, initToken, from, reserveToken2)
}
