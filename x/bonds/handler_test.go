package bonds_test

import (
	"github.com/ixoworld/bonds/x/bonds"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

// Tests

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
	require.Equal(t, res.Code, bonds.CodeBondAlreadyExists)
}

func TestCreatingABondUsingStakingTokenFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with token set to staking token
	msg := newValidMsgCreateBond()
	msg.Token = app.StakingKeeper.GetParams(ctx).BondDenom
	res := h(ctx, msg)

	require.False(t, res.IsOK())
	require.Equal(t, res.Code, bonds.CodeBondTokenInvalid)
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
	require.Equal(t, res.Code, bonds.CodeBondDoesNotExist)
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
	require.Equal(t, res.Code, bonds.CodeArgumentInvalid)
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
	require.Equal(t, res.Code, bonds.CodeArgumentMissingOrIncorrectType)
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
	require.Equal(t, res.Code, bonds.CodeArgumentInvalid)
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
	require.Equal(t, res.Code, bonds.CodeArgumentMissingOrIncorrectType)
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
	require.Equal(t, res.Code, bonds.CodeBondDoesNotExist)
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
	require.Equal(t, res.Code, bonds.CodeReserveDenomsMismatch)
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
	require.Equal(t, res.Code, bonds.CodeOrderQuantityLimitExceeded)
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
	require.Equal(t, res.Code, bonds.CodeMaxPriceExceeded)
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
	require.Equal(t, res.Code, bonds.CodeMaxPriceExceeded)
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
	require.Equal(t, res.Code, bonds.CodeMaxPriceExceeded)
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

	userBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, initFeeAddress)
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
	require.Equal(t, res.Code, bonds.CodeBondDoesNotExist)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))
}

func TestSellingABondWhichCannotBeSoldFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond
	createMsg := newValidMsgCreateBond()
	createMsg.AllowSells = types.FALSE
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
	require.Equal(t, res.Code, bonds.CodeBondDoesNotAllowSelling)
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
	require.Equal(t, res.Code, bonds.CodeOrderQuantityLimitExceeded)
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

	userBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, userAddress)
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

	userBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, userAddress)
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

	userBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, userAddress)
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

	userBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, initFeeAddress)
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
	require.Equal(t, res.Code, bonds.CodeBondDoesNotExist)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))
}

func TestSwapOrderQuantityLimitExceededFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := bonds.NewHandler(app.BondsKeeper)

	// Create bond with order quantity limit
	createMsg := newValidMsgCreateBond()
	createMsg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(reserveToken, 4))
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 4 tokens
	h(ctx, newValidMsgBuy(4, 10000))
	bonds.EndBlocker(ctx, app.BondsKeeper)

	// Swap tokens
	msg := types.NewMsgSwap(userAddress, token, sdk.NewInt64Coin(reserveToken, 5), reserveToken2)
	res := h(ctx, msg)

	userBalance := app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.False(t, res.IsOK())
	require.Equal(t, res.Code, bonds.CodeOrderQuantityLimitExceeded)
	require.Equal(t, sdk.NewInt(4), userBalance.AmountOf(token))
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

	userBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.BondsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.BondsKeeper.CoinKeeper.GetCoins(ctx, initFeeAddress)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(89990), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(90008), userBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(10009), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(9992), reserveBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken))
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
