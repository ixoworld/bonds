package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/ixoworld/bonds/x/bonds/app"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	token = "testtoken"

	token1 = "testtoken1"
	token2 = "testtoken2"
	token3 = "testtoken3"

	blankSanityRate             = "0"
	blankSanityMarginPercentage = "0"
	reserveToken                = "res"
	reserveToken2               = "rez"

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
	initAllowSell              = true
	initSigners                = []sdk.AccAddress{initCreator}
	initBatchBlocks            = sdk.NewUint(10)
	initState                  = types.OpenState

	buyPrices = sdk.NewDecCoins(sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 2),
		sdk.NewInt64Coin(reserveToken2, 3),
	))

	sellPrices = sdk.NewDecCoins(sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 4),
		sdk.NewInt64Coin(reserveToken2, 5),
	))

	// Base order
	baseOrderAddress = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	baseOrderAmount  = sdk.NewCoin(token, sdk.OneInt())

	// Buy
	buyerAddress = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	buyAmount    = sdk.NewCoin(token, sdk.OneInt())
	maxPrices    = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 6),
		sdk.NewInt64Coin(reserveToken2, 7),
	)

	// Sell
	sellerAddress = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	sellAmount    = sdk.NewCoin(token, sdk.OneInt())

	// Swapper
	swapperAddress = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	swapFrom       = sdk.NewCoin(reserveToken, sdk.OneInt())
	swapTo         = reserveToken2

	batchBlocks = sdk.NewUint(5)
)

func functionParametersPower() types.FunctionParams {
	return types.FunctionParams{
		types.NewFunctionParam("m", sdk.NewDec(12)),
		types.NewFunctionParam("n", sdk.NewDec(2)),
		types.NewFunctionParam("c", sdk.NewDec(100))}
}

func functionParametersAugmented() types.FunctionParams {
	return types.FunctionParams{
		types.NewFunctionParam("d0", sdk.MustNewDecFromStr("500.0")),
		types.NewFunctionParam("p0", sdk.MustNewDecFromStr("0.01")),
		types.NewFunctionParam("theta", sdk.MustNewDecFromStr("0.4")),
		types.NewFunctionParam("kappa", sdk.MustNewDecFromStr("3.0"))}
}

func powerReserves() []string   { return []string{reserveToken} }
func swapperReserves() []string { return []string{reserveToken, reserveToken2} }

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	return app, ctx
}

func getValidPowerFunctionBond() types.Bond {
	functionType := types.PowerFunction
	functionParams := functionParametersPower()
	reserveTokens := powerReserves()
	return types.NewBond(initToken, initName, initDescription,
		initCreator, functionType, functionParams,
		reserveTokens, initReserveAddress, initTxFeePercentage,
		initExitFeePercentage, initFeeAddress, initMaxSupply,
		initOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks, initState)
}

func getValidAugmentedFunctionBond() types.Bond {
	functionType := types.AugmentedFunction
	functionParams := functionParametersAugmented()
	reserveTokens := powerReserves()
	return types.NewBond(initToken, initName, initDescription,
		initCreator, functionType, functionParams,
		reserveTokens, initReserveAddress, initTxFeePercentage,
		initExitFeePercentage, initFeeAddress, initMaxSupply,
		initOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks, initState)
}

func getValidSwapperBond() types.Bond {
	functionType := types.SwapperFunction
	functionParams := types.FunctionParams(nil)
	reserveTokens := swapperReserves()
	return types.NewBond(initToken, initName, initDescription,
		initCreator, functionType, functionParams,
		reserveTokens, initReserveAddress, initTxFeePercentage,
		initExitFeePercentage, initFeeAddress, initMaxSupply,
		initOrderQuantityLimits, initSanityRate, initSanityMarginPercentage,
		initAllowSell, initSigners, initBatchBlocks, initState)
}

func getValidBond() types.Bond {
	return getValidPowerFunctionBond()
}

func getValidBatch() types.Batch {
	return types.NewBatch(token, batchBlocks)
}

func getValidBaseOrder() types.BaseOrder {
	return types.NewBaseOrder(baseOrderAddress, baseOrderAmount)
}

func getValidBuyOrder() types.BuyOrder {
	return types.NewBuyOrder(buyerAddress, buyAmount, maxPrices)
}

func getValidSellOrder() types.SellOrder {
	return types.NewSellOrder(sellerAddress, sellAmount)
}

func getValidSwapOrder() types.SwapOrder {
	return types.NewSwapOrder(swapperAddress, swapFrom, swapTo)
}
