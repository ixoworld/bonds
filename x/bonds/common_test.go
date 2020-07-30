package bonds_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/ixoworld/bonds/x/bonds/app"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	token  = "testtoken"
	token2 = "testtoken2"

	blankSanityRate             = "0"
	blankSanityMarginPercentage = "0"
	reserveToken                = "res"
	reserveToken2               = "rez"

	anotherAddress = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	userAddress    = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	functionParametersPower = types.FunctionParams{
		types.NewFunctionParam("m", sdk.NewDec(12)),
		types.NewFunctionParam("n", sdk.NewDec(2)),
		types.NewFunctionParam("c", sdk.NewDec(100)),
	}

	functionParametersAugmented = types.FunctionParams{
		types.NewFunctionParam("d0", sdk.MustNewDecFromStr("500.0")),
		types.NewFunctionParam("p0", sdk.MustNewDecFromStr("0.01")),
		types.NewFunctionParam("theta", sdk.MustNewDecFromStr("0.4")),
		types.NewFunctionParam("kappa", sdk.MustNewDecFromStr("3.0")),
	}

	powerReserves   = []string{reserveToken}
	swapperReserves = []string{reserveToken, reserveToken2}

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
	initBatchBlocks            = sdk.OneUint()

	amountLTMaxSupply = initMaxSupply.Amount.Sub(sdk.OneInt()).Int64()
	amountGTMaxSupply = initMaxSupply.Amount.Add(sdk.OneInt()).Int64()

	extraMaxSupply = sdk.NewInt64Coin(initToken, 1000000)
)

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	return app, ctx
}

// Helpers

func newSimpleBond() types.Bond {
	return types.Bond{
		Token:   token,
		Creator: initCreator,
		Signers: initSigners,
	}
}

func newValidMsgCreateSwapperBond() types.MsgCreateBond {
	validMsg := newValidMsgCreateBond()
	validMsg.FunctionType = types.SwapperFunction
	validMsg.FunctionParameters = nil
	validMsg.ReserveTokens = swapperReserves
	return validMsg
}

func newValidMsgCreateAugmentedBond() types.MsgCreateBond {
	validMsg := newValidMsgCreateBond()
	validMsg.FunctionType = types.AugmentedFunction
	validMsg.FunctionParameters = functionParametersAugmented
	validMsg.MaxSupply = extraMaxSupply
	return validMsg
}

func newValidMsgCreateBond() types.MsgCreateBond {
	functionType := types.PowerFunction
	functionParams := functionParametersPower
	reserveTokens := powerReserves
	return types.NewMsgCreateBond(token, initName, initDescription,
		initCreator, functionType, functionParams, reserveTokens,
		initTxFeePercentage, initExitFeePercentage, initFeeAddress,
		initMaxSupply, initOrderQuantityLimits, initSanityRate,
		initSanityMarginPercentage, initAllowSell, initSigners, initBatchBlocks)
}

func newValidMsgBuy(amount int64, maxPrice int64) types.MsgBuy {
	amountCoin := sdk.NewInt64Coin(token, amount)
	maxPrices := sdk.NewCoins(sdk.NewInt64Coin(reserveToken, maxPrice))
	return types.NewMsgBuy(userAddress, amountCoin, maxPrices)
}

func newValidMsgSell(amount int64) types.MsgSell {
	amountCoin := sdk.NewInt64Coin(token, amount)
	return types.NewMsgSell(userAddress, amountCoin)
}

func newValidMsgSwap(fromToken, toToken string, amount int64) types.MsgSwap {
	fromAmount := sdk.NewInt64Coin(fromToken, amount)
	return types.NewMsgSwap(userAddress, token, fromAmount, toToken)
}

func addCoinsToUser(app *simapp.SimApp, ctx sdk.Context, coins sdk.Coins) sdk.Error {
	_, err := app.BondsKeeper.BankKeeper.AddCoins(ctx, userAddress, coins)
	return err
}
