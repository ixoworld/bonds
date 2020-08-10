package bonds_test

import (
	"github.com/ixoworld/bonds/x/bonds"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestInvalidMsgFails(t *testing.T) {
	_, ctx := createTestApp(false)
	h := bonds.NewHandler(bonds.Keeper{})

	msg := sdk.NewTestMsg()
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "Unrecognized bonds Msg type: "+sdk.NewTestMsg().Type())
}

func TestCreateValidBond(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	res := h(ctx, newValidMsgCreateBond())

	require.True(t, res.IsOK())
	require.True(t, app.BondsKeeper.BondExists(ctx, token))

	// Check assigned initial state
	bond := app.BondsKeeper.MustGetBond(ctx, token)
	require.Equal(t, types.OpenState, bond.State)
}

func TestCreateValidAugmentedBondHatchState(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create augmented function bond
	res := h(ctx, newValidMsgCreateAugmentedBond())

	require.True(t, res.IsOK())
	require.True(t, app.BondsKeeper.BondExists(ctx, token))

	// Check initial state (hatch since augmented)
	bond := app.BondsKeeper.MustGetBond(ctx, token)
	require.Equal(t, types.HatchState, bond.State)

	// Check function params (R0, S0, V0 added)
	paramsMap := bond.FunctionParameters.AsMap()
	d0, _ := paramsMap["d0"]
	p0, _ := paramsMap["p0"]
	theta, _ := paramsMap["theta"]
	kappa, _ := paramsMap["kappa"]

	initialParams := functionParametersAugmented().AsMap()
	require.Equal(t, d0, initialParams["d0"])
	require.Equal(t, p0, initialParams["p0"])
	require.Equal(t, theta, initialParams["theta"])
	require.Equal(t, kappa, initialParams["kappa"])

	R0 := d0.Mul(sdk.OneDec().Sub(theta))
	S0 := d0.Quo(p0)
	V0 := types.Invariant(R0, S0, kappa.TruncateInt64())

	require.Equal(t, R0, paramsMap["R0"])
	require.Equal(t, S0, paramsMap["S0"])
	require.Equal(t, V0, paramsMap["V0"])
	require.Len(t, bond.FunctionParameters, 7)
}

func TestCreateBondThatAlreadyExistsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	bond := types.Bond{Token: token}
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Create bond with same token
	res := h(ctx, newValidMsgCreateBond())

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondAlreadyExists)
}

func TestCreatingABondUsingStakingTokenFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with token set to staking token
	msg := newValidMsgCreateBond()
	msg.Token = app.StakingKeeper.GetParams(ctx).BondDenom
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondTokenInvalid)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))
}

func TestEditingANonExistingBondFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "",
		"0", "0", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondDoesNotExist)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))
}

func TestEditingABondWithDifferentCreatorAndSignersFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "",
		"0", "0", initCreator, []sdk.AccAddress{anotherAddress})
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "List of signers does not match the one in the bond")
}

func TestEditingABondWithNegativeOrderQuantityLimitsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "-10testtoken",
		"0", "0", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "invalid coin expression")
}

func TestEditingABondWithFloatOrderQuantityLimitsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "10.5testtoken",
		"0", "0", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "invalid coin expression")
}

func TestEditingABondWithSanityRateEmptyStringMakesSanityFieldsZero(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	bond := newSimpleBond()
	bond.SanityRate = sdk.OneDec()
	bond.SanityMarginPercentage = sdk.OneDec()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Check sanity values before
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	require.NotEqual(t, sdk.ZeroDec(), bond.SanityRate)
	require.NotEqual(t, sdk.ZeroDec(), bond.SanityMarginPercentage)

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "10testtoken",
		"", "", initCreator, initSigners)
	res := h(ctx, msg)

	// Check sanity values after
	require.True(t, res.IsOK())
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	require.Equal(t, sdk.ZeroDec(), bond.SanityRate)
	require.Equal(t, sdk.ZeroDec(), bond.SanityMarginPercentage)
}

func TestEditingABondWithNegativeSanityRateFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "10testtoken",
		"-10", "", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeArgumentInvalid)
}

func TestEditingABondWithNonFloatSanityRateFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "10testtoken",
		"20t", "", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeArgumentMissingOrIncorrectType)
}

func TestEditingABondWithNegativeSanityMarginPercentageFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "10testtoken",
		"10", "-5", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeArgumentInvalid)
}

func TestEditingABondWithNonFloatSanityMarginPercentageFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	msg := types.NewMsgEditBond(token, initName, initDescription, "10testtoken",
		"20", "20t", initCreator, initSigners)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeArgumentMissingOrIncorrectType)
}

func TestEditingABondCorrectlyPasses(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Set bond to simulate creation
	app.BondsKeeper.SetBond(ctx, token, newSimpleBond())

	// Edit bond
	newName := "a new name"
	newDescription := "a new description"
	msg := types.NewMsgEditBond(token, newName, newDescription, "",
		"0", "0", initCreator, initSigners)
	res := h(ctx, msg)

	require.True(t, res.IsOK())
	bond, _ := app.BondsKeeper.GetBond(ctx, token)
	require.Equal(t, newName, bond.Name)
	require.Equal(t, newDescription, bond.Description)
	require.Equal(t, sdk.Coins(nil), bond.OrderQuantityLimits)
	require.Equal(t, sdk.ZeroDec(), bond.SanityRate)
	require.Equal(t, sdk.ZeroDec(), bond.SanityMarginPercentage)
}

func TestBuyingANonExistingBondFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Buy 1 token
	res := h(ctx, newValidMsgBuy(1, 10))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondDoesNotExist)
}

func TestBuyingABondWithNonExistentToken(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Buy tokens of another bond
	msg := newValidMsgBuy(amountLTMaxSupply, 0) // 0 max prices replaced below
	msg.MaxPrices = sdk.Coins{sdk.NewInt64Coin(token2, 10)}
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeReserveDenomsMismatch)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingABondWithMaxPriceBiggerThanBalanceFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 10 tokens
	res := h(ctx, newValidMsgBuy(10, 5000))

	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "insufficient account funds; 4000res < 5000res")
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingBondWithOrderQuantityLimitExceededFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with order quantity limit
	msg := newValidMsgCreateBond()
	msg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(token, 4))
	h(ctx, msg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 10 tokens
	res := h(ctx, newValidMsgBuy(10, 4000))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeOrderQuantityLimitExceeded)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingABondExceedingMaxSupplyFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 6000)})
	require.Nil(t, err)

	// Buy an amount greater than max supply
	res := h(ctx, newValidMsgBuy(amountGTMaxSupply, 10))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeInvalidResultantSupply)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingABondExceedingMaxPriceFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 6000)})
	require.Nil(t, err)

	// Buy an amount less than max supply but with low max prices
	res := h(ctx, newValidMsgBuy(amountLTMaxSupply, 1))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeMaxPriceExceeded)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingABondWithoutSufficientFundsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 10 tokens
	res := h(ctx, newValidMsgBuy(10, 4000))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeMaxPriceExceeded)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingABondWithoutSufficientFundsDueToTxFeeFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 5000)})
	require.Nil(t, err)

	// Buy 10 tokens
	res := h(ctx, newValidMsgBuy(10, 5000))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeMaxPriceExceeded)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingABondCorrectlyPasses(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 2 tokens
	res := h(ctx, newValidMsgBuy(2, 4000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(3767), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(232), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(2), currentSupply.Amount)
}

func TestSellingANonExistingBondFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Sell 10 tokens
	res := h(ctx, newValidMsgSell(10))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondDoesNotExist)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))
}

func TestSellingABondWhichCannotBeSoldFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	createMsg := newValidMsgCreateBond()
	createMsg.AllowSells = false
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Sell 10 tokens
	bondPreSell := app.BondsKeeper.MustGetBond(ctx, token)
	res := h(ctx, newValidMsgSell(10))
	bondPostSell := app.BondsKeeper.MustGetBond(ctx, token)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondDoesNotAllowSelling)
	require.Equal(t, bondPostSell.CurrentSupply.Amount, bondPreSell.CurrentSupply.Amount)
}

func TestSellBondExceedingOrderQuantityLimitFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with order quantity limit
	msg := newValidMsgCreateBond()
	msg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(token, 4))
	h(ctx, msg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Sell 10 tokens
	bondPreSell := app.BondsKeeper.MustGetBond(ctx, token)
	res := h(ctx, newValidMsgSell(10))
	bondPostSell := app.BondsKeeper.MustGetBond(ctx, token)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeOrderQuantityLimitExceeded)
	require.Equal(t, bondPostSell.CurrentSupply.Amount, bondPreSell.CurrentSupply.Amount)
}

func TestSellingABondWithAmountGreaterThanBalanceFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Sell 11 tokens
	bondPreSell := app.BondsKeeper.MustGetBond(ctx, token)
	res := h(ctx, newValidMsgSell(11))
	bondPostSell := app.BondsKeeper.MustGetBond(ctx, token)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "insufficient account funds; 4995"+reserveToken+",10"+token+" < 11"+token)
	require.Equal(t, bondPostSell.CurrentSupply.Amount, bondPreSell.CurrentSupply.Amount)
	require.Equal(t, sdk.NewInt(10), userBalance.AmountOf(token))
}

func TestSellingABondWhichSellerDoesNotOwnFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create first bond
	h(ctx, newValidMsgCreateBond())

	// Create second bond (different token)
	bond2Msg := newValidMsgCreateBond()
	bond2Msg.Token = token2
	h(ctx, bond2Msg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Sell 11 of a different bond
	msg := newValidMsgSell(0) // 0 amount replaced below
	msg.Amount = sdk.NewInt64Coin(token2, 11)
	res := h(ctx, msg)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	require.False(t, res.IsOK())
	require.Contains(t, res.Log,
		"insufficient account funds; 4995"+reserveToken+",10"+token+" < 11"+token2)
	require.Equal(t, sdk.NewInt(10), userBalance.AmountOf(token))
	require.Equal(t, sdk.ZeroInt(), userBalance.AmountOf(token2))
}

func TestSellingMoreTokensThanThereIsSupplyFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Sell an amount greater than the max supply
	bondPreSell := app.BondsKeeper.MustGetBond(ctx, token)
	res := h(ctx, newValidMsgSell(amountGTMaxSupply))
	bondPostSell := app.BondsKeeper.MustGetBond(ctx, token)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	require.False(t, res.IsOK())
	require.Contains(t, res.Log, "insufficient account funds; 4995"+reserveToken+",10"+token+" < 10001"+token)
	require.Equal(t, sdk.NewInt(10), userBalance.AmountOf(token))
	require.Equal(t, bondPreSell.CurrentSupply.Amount, bondPostSell.CurrentSupply.Amount)
}

func TestSellingABondCorrectlyPasses(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 2 tokens
	h(ctx, newValidMsgBuy(2, 4000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Sell 2 tokens
	msg := newValidMsgSell(2)
	res := h(ctx, msg)
	bonds.EndBlocker(ctx, app.BondsKeeper)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	currentSupply := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(3997), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.ZeroInt(), userBalance.AmountOf(token))
	require.Equal(t, sdk.ZeroInt(), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(3), feeBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.ZeroInt(), currentSupply.Amount)
}

func TestSwapBondDoesNotExistFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Swap tokens
	res := h(ctx, newValidMsgSwap(reserveToken, reserveToken2, 1))

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeBondDoesNotExist)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))
}

func TestSwapOrderInvalidReserveDenomsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateSwapperBond())

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Perform swap (invalid instead of reserveToken)
	res := h(ctx, newValidMsgSwap("invalid", reserveToken2, 10))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	userBalance := app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeReserveDenomsMismatch)
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))

	// Perform swap (invalid instead of reserveToken2)
	res = h(ctx, newValidMsgSwap(reserveToken, "invalid", 10))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	userBalance = app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeReserveDenomsMismatch)
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
}

func TestSwapOrderQuantityLimitExceededFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with order quantity limit
	createMsg := newValidMsgCreateSwapperBond()
	createMsg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(reserveToken, 4))
	h(ctx, createMsg)

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Perform swap
	msg := types.NewMsgSwap(userAddress, token, sdk.NewInt64Coin(reserveToken, 5), reserveToken2)
	res := h(ctx, msg)

	userBalance := app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, types.CodeOrderQuantityLimitExceeded)
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
}

func TestSwapInvalidAmount(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateSwapperBond())

	// Add reserve tokens to user (but not enough)
	nineReserveTokens := sdk.NewInt64Coin(reserveToken, 9)
	tenReserveTokens := sdk.NewInt64Coin(reserveToken, 10)
	err := addCoinsToUser(app, ctx, sdk.Coins{nineReserveTokens})
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Perform swap
	msg := types.NewMsgSwap(userAddress, token, tenReserveTokens, reserveToken2)
	res := h(ctx, msg)

	require.False(t, res.IsOK())
}

func TestSwapValidAmount(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateSwapperBond())

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Perform swap
	res := h(ctx, newValidMsgSwap(reserveToken, reserveToken2, 10))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(89990), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(90008), userBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(10009), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(9992), reserveBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken))
}

func TestSwapValidAmountReversed(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateSwapperBond())

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Perform swap
	res := h(ctx, newValidMsgSwap(reserveToken2, reserveToken, 10))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(90008), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(89990), userBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(9992), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(10009), reserveBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken2))
}

func TestMakeOutcomePayment(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with 100k outcome payment
	bondMsg := newValidMsgCreateBond()
	bondMsg.OutcomePayment = sdk.NewCoins(sdk.NewInt64Coin(reserveToken, 100000))
	h(ctx, bondMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 100000)})
	require.Nil(t, err)

	// Make outcome payment
	res := h(ctx, newValidMsgMakeOutcomePayment())
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Check that outcome payment is now in the bond reserve
	userBalance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.ZeroInt(), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(100000), reserveBalance.AmountOf(reserveToken))

	// Check that the bond is now in SETTLE state
	require.Equal(t, types.SettleState, app.BondsKeeper.MustGetBond(ctx, token).State)
}

func TestWithdrawShare(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	h(ctx, newValidMsgCreateBond())

	// Set bond current supply to 3 and state to SETTLE
	bond := app.BondsKeeper.MustGetBond(ctx, token)
	bond.CurrentSupply = sdk.NewCoin(bond.Token, sdk.NewInt(3))
	bond.State = types.SettleState
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Mint 3 bond tokens and send [2 to user 1] and [1 to user 2]
	err := app.SupplyKeeper.MintCoins(ctx, types.BondsMintBurnAccount,
		sdk.NewCoins(sdk.NewInt64Coin(token, 3)))
	require.Nil(t, err)
	err = app.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.BondsMintBurnAccount,
		userAddress, sdk.NewCoins(sdk.NewInt64Coin(token, 2)))
	require.Nil(t, err)
	err = app.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.BondsMintBurnAccount,
		anotherAddress, sdk.NewCoins(sdk.NewInt64Coin(token, 1)))
	require.Nil(t, err)

	// Simulate outcome payment by depositing (freshly minted) 100k into reserve
	hundredK := sdk.NewCoins(sdk.NewCoin(reserveToken, sdk.NewInt(100000)))
	err = app.SupplyKeeper.MintCoins(ctx, types.BondsMintBurnAccount, hundredK)
	require.Nil(t, err)
	err = app.BondsKeeper.DepositReserveFromModule(
		ctx, bond.Token, types.BondsMintBurnAccount, hundredK)
	require.Nil(t, err)

	// User 1 withdraws share
	res := h(ctx, newValidMsgWithdrawShareFrom(userAddress))
	require.True(t, res.IsOK())
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// User 1 had 2 tokens out of the supply of 3 tokens, so user 1 gets 2/3
	user1Balance := app.BondsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(66666), user1Balance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(33334), reserveBalance.AmountOf(reserveToken))

	// Note: rounding is rounded to floor, so despite user 1 being owed 66666.67
	// tokens, user 1 gets 66666 and not 66667 tokens. Then, since user 2 now owns
	// the entire share of the bond tokens, they will get 100% of the remaining
	// 33334 tokens, which is more than what was initially owed (33333.33).

	// User 2 withdraws share
	res = h(ctx, newValidMsgWithdrawShareFrom(anotherAddress))
	require.True(t, res.IsOK())
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// User 2 had 1 token out of the remaining supply of 1 token, so user 2 gets all remaining
	user2Balance := app.BondsKeeper.BankKeeper.GetCoins(ctx, anotherAddress)
	reserveBalance = app.BondsKeeper.GetReserveBalances(ctx, initToken)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(33334), user2Balance.AmountOf(reserveToken))
	require.Equal(t, sdk.ZeroInt(), reserveBalance.AmountOf(reserveToken))
}

func TestDecrementRemainingBlocksCountAfterEndBlock(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	two := sdk.NewUint(2)
	one := sdk.NewUint(1)

	// Create bond
	createMsg := newValidMsgCreateBond()
	createMsg.BatchBlocks = two
	h(ctx, createMsg)

	require.Equal(t, two, app.BondsKeeper.MustGetBatch(ctx, token).BlocksRemaining)
	bonds.EndBlocker(ctx, app.BondsKeeper)
	require.Equal(t, one, app.BondsKeeper.MustGetBatch(ctx, token).BlocksRemaining)
}

func TestEndBlockerDoesNotPerformOrdersBeforeASpecifiedNumberOfBlocks(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with batch blocks set to 2
	createMsg := newValidMsgCreateBond()
	createMsg.BatchBlocks = sdk.NewUint(2)
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Buy 4 tokens
	h(ctx, newValidMsgBuy(2, 10000))
	h(ctx, newValidMsgBuy(2, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	require.Equal(t, len(app.BondsKeeper.MustGetBatch(ctx, token).Buys), 2)
}

func TestEndBlockerPerformsOrdersAfterASpecifiedNumberOfBlocks(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	createMsg := newValidMsgCreateBond()
	createMsg.BatchBlocks = sdk.NewUint(3)
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Buy 4 tokens
	h(ctx, newValidMsgBuy(2, 10000))
	h(ctx, newValidMsgBuy(2, 10000))

	// Run EndBlocker for N times, where N = BatchBlocks
	batchBlocksInt := int(createMsg.BatchBlocks.Uint64())
	for i := 0; i <= batchBlocksInt; i++ {
		bonds.EndBlocker(ctx, app.BondsKeeper)
	}

	// Buys have been performed
	require.Equal(t, 0, len(app.BondsKeeper.MustGetBatch(ctx, token).Buys))
}

func TestEndBlockerAugmentedFunction(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with augmented function type
	createMsg := newValidMsgCreateAugmentedBond()
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Get bond to confirm allowSells==false, S0==50000, state==hatch
	bond := app.BondsKeeper.MustGetBond(ctx, token)
	require.False(t, bond.AllowSells)
	require.Equal(t, sdk.NewDec(50000), bond.FunctionParameters.AsMap()["S0"])
	require.Equal(t, types.HatchState, bond.State)

	// - Buy 49999 tokens; just below S0
	// - Cannot buy 2 tokens in the meantime, since this exceeds S0
	// - Cannot sell tokens (not even 1) in hatch state
	res := h(ctx, newValidMsgBuy(49999, 100000))
	require.True(t, res.IsOK())
	res = h(ctx, newValidMsgBuy(2, 100000))
	require.False(t, res.IsOK())
	res = h(ctx, newValidMsgSell(1))
	require.False(t, res.IsOK())
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Confirm allowSells and state still the same
	bond = app.BondsKeeper.MustGetBond(ctx, token)
	require.False(t, bond.AllowSells)
	require.Equal(t, types.HatchState, bond.State)

	// Buy 1 more token, to reach S0 => state is now open
	h(ctx, newValidMsgBuy(1, 100000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Confirm allowSells==true, state==open
	bond = app.BondsKeeper.MustGetBond(ctx, token)
	require.True(t, bond.AllowSells)
	require.Equal(t, types.OpenState, bond.State)

	// Check user balance of tokens
	balance := app.BankKeeper.GetCoins(ctx, userAddress).AmountOf(token).Int64()
	require.Equal(t, int64(50000), balance)

	// Can now sell tokens (all 50,000 of them)
	res = h(ctx, newValidMsgSell(50000))
	require.True(t, res.IsOK())
	bonds.EndBlocker(ctx, app.BondsKeeper)
	balance = app.BankKeeper.GetCoins(ctx, userAddress).AmountOf(token).Int64()
	require.Equal(t, int64(0), balance)
}

func TestEndBlockerAugmentedFunctionDecimalS0(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with augmented function type
	createMsg := newValidMsgCreateAugmentedBond()
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Change bond's S0 parameter to 49999.5
	decimalS0 := sdk.MustNewDecFromStr("49999.5")
	bond := app.BondsKeeper.MustGetBond(ctx, token)
	for i, p := range bond.FunctionParameters {
		if p.Param == "S0" {
			bond.FunctionParameters[i].Value = decimalS0
			break
		}
	}
	app.BondsKeeper.SetBond(ctx, bond.Token, bond)

	// Get bond to confirm S0==49999.5, allowSells==false, state==hatch
	bond = app.BondsKeeper.MustGetBond(ctx, token)
	require.Equal(t, decimalS0, bond.FunctionParameters.AsMap()["S0"])
	require.False(t, bond.AllowSells)
	require.Equal(t, types.HatchState, bond.State)

	// Buy 49999 tokens; just below S0
	h(ctx, newValidMsgBuy(49999, 100000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Confirm allowSells and state still the same
	bond = app.BondsKeeper.MustGetBond(ctx, token)
	require.False(t, bond.AllowSells)
	require.Equal(t, types.HatchState, bond.State)

	// Buy 1 more token, to exceed S0
	h(ctx, newValidMsgBuy(1, 100000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Confirm allowSells==true, state==open
	bond = app.BondsKeeper.MustGetBond(ctx, token)
	require.True(t, bond.AllowSells)
	require.Equal(t, types.OpenState, bond.State)
}
