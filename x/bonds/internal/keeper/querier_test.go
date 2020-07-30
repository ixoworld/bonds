package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ixoworld/bonds/x/bonds/internal/keeper"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestNewQuerier(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err := querier(ctx, []string{"foo", "bar"}, req)
	require.Error(t, err)
}

func TestQueryBonds(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.QueryBonds

	// Initially no errors and zero bonds
	res, err := querier(ctx, []string{keeper.QueryBonds}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Len(t, queryResult, 0)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Still not errors but one bond
	res, err = querier(ctx, []string{keeper.QueryBonds}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, types.QueryBonds{token})
}

func TestQueryBond(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.Bond

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryBond, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// No error because of new bond
	res, err = querier(ctx, []string{keeper.QueryBond, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, bond)
}

func TestQueryBatch(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.Batch

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryBatch, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	batch := getValidBatch()
	app.BondsKeeper.SetBatch(ctx, token, batch)

	// No error because of new bond
	res, err = querier(ctx, []string{keeper.QueryBatch, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, batch)
}

func TestQueryLastBatch(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.Batch

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryLastBatch, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	batch := getValidBatch()
	app.BondsKeeper.SetLastBatch(ctx, token, batch)

	// No error because of new bond
	res, err = querier(ctx, []string{keeper.QueryLastBatch, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, batch)
}

func TestQueryCurrentPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult sdk.DecCoins

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Get current price directly
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)
	currentPrices, _ := bond.GetCurrentPricesPT(reserveBalances)

	// Calculate current price manually
	// y = mx^n + c = 12(0^2) + 100 = 0 + 100 = 100
	manualPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}

	// Check that prices are correct
	res, err = querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, currentPrices)
	require.Equal(t, queryResult, manualPrices)

	// Change current supply to 10 for increased price
	bond.CurrentSupply = sdk.NewInt64Coin(bond.Token, 10)
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Get current price directly
	reserveBalances = app.BondsKeeper.GetReserveBalances(ctx, token)
	currentPrices, _ = bond.GetCurrentPricesPT(reserveBalances)

	// Calculate current price manually
	// y = mx^n + c = 12(10^2) + 100 = 1200 + 100 = 1300
	manualPrices = sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 1300)}

	// Check that prices are correct
	res, err = querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, currentPrices)
	require.Equal(t, queryResult, manualPrices)
}

func TestQueryCurrentPriceWithZeroPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult sdk.DecCoins

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	bond := getValidBond()
	bond.FunctionParameters = types.FunctionParams{
		types.NewFunctionParam("m", sdk.NewDec(12)),
		types.NewFunctionParam("n", sdk.NewDec(2)),
		types.NewFunctionParam("c", sdk.NewDec(0))} // set to 0 for P(R=0)=0
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Calculate current price manually
	// y = mx^n + c = 12(0^2) + 0 = 0 + 0 = 0
	manualPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 0)}

	// Check that prices are correct
	res, err = querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, manualPrices, queryResult)
	require.Equal(t, "0.000000000000000000res", queryResult.String())

	// Important note: the fact that queryResult is "0.000000000000000000res"
	// rather than the default "" (for empty coins) is intentional
}

func TestQueryCurrentPriceForSwapper(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult sdk.DecCoins

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add swapper bond
	bond := getValidSwapperBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Get current price directly (error since swapper not initialised)
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)
	_, err = bond.GetCurrentPricesPT(reserveBalances)
	if err == nil {
		panic("expected error")
	}

	// Check that error since swapper not initialised
	res, err = querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Update current supply to 2
	bond.CurrentSupply = sdk.NewInt64Coin(bond.Token, 2)
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Send 200res,300rez to reserve
	newReserve := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 200),
		sdk.NewInt64Coin(reserveToken2, 300),
	)
	_, _ = app.BankKeeper.AddCoins(ctx, bond.ReserveAddress, newReserve)

	// Get current price directly
	reserveBalances = app.BondsKeeper.GetReserveBalances(ctx, token)
	currentPrices, _ := bond.GetCurrentPricesPT(reserveBalances)

	// Calculate current price manually
	// (since 2 tokens (current supply) => 200res,300rez then by the
	//  constant product formula, another 1 token => 100res,150rez)
	manualPrices := sdk.DecCoins([]sdk.DecCoin{
		sdk.NewInt64DecCoin(reserveToken, 100),
		sdk.NewInt64DecCoin(reserveToken2, 150),
	})

	// Check that prices are correct
	res, err = querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, currentPrices)
	require.Equal(t, queryResult, manualPrices)
}

func TestQueryCurrentReserve(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult sdk.Coins

	// Initially error since no bond
	res, err := querier(ctx, []string{keeper.QueryCurrentPrice, token}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Check that reserve balances are correct (initially empty)
	res, err = querier(ctx, []string{keeper.QueryCurrentReserve, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, "0res", queryResult.String())

	// Send 200res,300rez to reserve
	newReserve := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 200),
		sdk.NewInt64Coin(reserveToken2, 300),
	)
	_, _ = app.BankKeeper.AddCoins(ctx, bond.ReserveAddress, newReserve)

	// Get current reserve (now 200token2,300token3)
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)

	// Check that reserve balances are correct
	res, err = querier(ctx, []string{keeper.QueryCurrentReserve, token}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, reserveBalances)
	require.Equal(t, queryResult, newReserve)
}

func TestQueryCustomPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult sdk.DecCoins

	// Initially error since no bond
	dummySupply := sdk.ZeroInt()
	res, err := querier(ctx,
		[]string{keeper.QueryCustomPrice, token, dummySupply.String()}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Get custom price directly with supply=0
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	customSupply := sdk.ZeroInt()
	customPrices, _ := bond.GetPricesAtSupply(customSupply)

	// Calculate current price manually
	// y = mx^n + c = 12(0^2) + 100 = 0 + 100 = 100
	manualPrices := sdk.DecCoins([]sdk.DecCoin{sdk.NewInt64DecCoin(reserveToken, 100)})

	// Check that prices are correct
	res, err = querier(ctx,
		[]string{keeper.QueryCustomPrice, token, customSupply.String()}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, customPrices)
	require.Equal(t, queryResult, manualPrices)

	// Get custom price directly with supply=10
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	customSupply = sdk.NewInt(10)
	customPrices, _ = bond.GetPricesAtSupply(customSupply)

	// Calculate current price manually
	// y = mx^n + c = 12(10^2) + 100 = 1200 + 100 = 1300
	manualPrices = sdk.DecCoins([]sdk.DecCoin{sdk.NewInt64DecCoin(reserveToken, 1300)})

	// Check that prices are correct
	res, err = querier(ctx,
		[]string{keeper.QueryCustomPrice, token, customSupply.String()}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult, customPrices)
	require.Equal(t, queryResult, manualPrices)
}

func TestQueryBuyPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.QueryBuyPrice

	// Initially error since no bond
	dummyAmount := sdk.OneInt().String()
	res, err := querier(ctx,
		[]string{keeper.QueryBuyPrice, token, dummyAmount}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond and batch (batch necessary since buy price considers buy orders)
	bond := getValidBond()
	batch := getValidBatch()
	app.BondsKeeper.SetBond(ctx, token, bond)
	app.BondsKeeper.SetBatch(ctx, token, batch)

	// Get buy price for 10 tokens directly
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	buyAmount := sdk.NewInt(10)
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)
	buyPrices, _ := bond.GetPricesToMint(buyAmount, reserveBalances)
	txFees := bond.GetTxFees(buyPrices)
	roundedPrices := types.RoundReservePrices(buyPrices)
	roundedTotalPrices := roundedPrices.Add(txFees)

	// Calculate buy price manually
	// reserveAt(10) = (m/n+1)x^(n+1) + xc = (12/3)(10^(2+1)) + 10(100) = 5000
	// reserveAt(0) = (m/n+1)x^(n+1) + xc = (12/3)(0^(2+1)) + 0(100) = 0
	// price = reserveAt(10) - reserveAt(0) = 5000
	manualPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 5000)}

	// Adjusted supply will be current (0) + buy orders (0)
	manualSupply := sdk.NewInt64Coin(bond.Token, 0)

	// Check that prices and fees are correct
	res, err = querier(ctx,
		[]string{keeper.QueryBuyPrice, token, buyAmount.String()}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult.AdjustedSupply, manualSupply)
	require.Equal(t, queryResult.Prices, roundedPrices)
	require.Equal(t, queryResult.Prices, manualPrices)
	require.Equal(t, queryResult.TxFees, txFees)
	require.Equal(t, queryResult.TotalFees, txFees)
	require.Equal(t, queryResult.TotalPrices, roundedTotalPrices)

	// Simulate the above buy taking place
	_, _ = app.BankKeeper.AddCoins(ctx, bond.ReserveAddress, queryResult.Prices)
	_, _ = app.BankKeeper.AddCoins(ctx, bond.FeeAddress, queryResult.TotalFees)
	app.BondsKeeper.SetCurrentSupply(ctx, token, sdk.NewCoin(token, buyAmount))

	// Get buy price for 5 MORE tokens directly
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	buyAmount = sdk.NewInt(5)
	reserveBalances = app.BondsKeeper.GetReserveBalances(ctx, token)
	buyPrices, _ = bond.GetPricesToMint(buyAmount, reserveBalances)
	txFees = bond.GetTxFees(buyPrices)
	roundedPrices = types.RoundReservePrices(buyPrices)
	roundedTotalPrices = roundedPrices.Add(txFees)

	// Calculate buy price manually
	// reserveAt(15) = (m/n+1)x^(n+1) + xc = (12/3)(15^(2+1)) + 15(100) = 15000
	// reserveAt(10) = (m/n+1)x^(n+1) + xc = (12/3)(10^(2+1)) + 10(100) = 5000
	// price = reserveAt(15) - reserveAt(10) = 10000
	manualPrices = sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)}

	// Adjusted supply will be current (10) + buy orders (0)
	manualSupply = sdk.NewInt64Coin(bond.Token, 10)

	// Check that prices and fees are correct
	res, err = querier(ctx,
		[]string{keeper.QueryBuyPrice, token, buyAmount.String()}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult.AdjustedSupply, manualSupply)
	require.Equal(t, queryResult.Prices, roundedPrices)
	require.Equal(t, queryResult.Prices, manualPrices)
	require.Equal(t, queryResult.TxFees, txFees)
	require.Equal(t, queryResult.TotalFees, txFees)
	require.Equal(t, queryResult.TotalPrices, roundedTotalPrices)
}

func TestQuerySellPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.QuerySellReturn
	var buyQueryResult types.QueryBuyPrice

	// Initially error since no bond
	dummyAmount := sdk.OneInt().String()
	res, err := querier(ctx,
		[]string{keeper.QuerySellReturn, token, dummyAmount}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add bond and batch (batch necessary since sell returns considers sell orders)
	bond := getValidBond()
	batch := getValidBatch()
	app.BondsKeeper.SetBond(ctx, token, bond)
	app.BondsKeeper.SetBatch(ctx, token, batch)

	// Still an error error since current supply is zero
	dummyAmount = sdk.OneInt().String()
	res, err = querier(ctx,
		[]string{keeper.QuerySellReturn, token, dummyAmount}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Simulate a buy of 10 tokens
	buyAmount := sdk.NewInt(10)
	res, err = querier(ctx,
		[]string{keeper.QueryBuyPrice, token, buyAmount.String()}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &buyQueryResult)
	_, _ = app.BankKeeper.AddCoins(ctx, bond.ReserveAddress, buyQueryResult.Prices)
	_, _ = app.BankKeeper.AddCoins(ctx, bond.FeeAddress, buyQueryResult.TotalFees)
	app.BondsKeeper.SetCurrentSupply(ctx, token, sdk.NewCoin(token, buyAmount))

	// Get sell returns for 10 tokens directly
	bond, _ = app.BondsKeeper.GetBond(ctx, token)
	sellAmount := sdk.NewInt(10)
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)
	sellReturns := bond.GetReturnsForBurn(buyAmount, reserveBalances)
	txFees := bond.GetTxFees(sellReturns)
	exitFees := bond.GetExitFees(sellReturns)
	totalFees := txFees.Add(exitFees)
	roundedReturns := types.RoundReserveReturns(sellReturns)
	roundedTotalReturns := roundedReturns.Sub(totalFees)

	// Calculate sell returns manually
	// reserveAt(10) = (m/n+1)x^(n+1) + xc = (12/3)(10^(2+1)) + 10(100) = 5000
	// reserveAt(0) = (m/n+1)x^(n+1) + xc = (12/3)(0^(2+1)) + 0(100) = 0
	// returns = reserveAt(10) - reserveAt(0) = 5000
	manualReturns := sdk.Coins{sdk.NewInt64Coin(reserveToken, 5000)}

	// Adjusted supply will be current (10) - sell orders (0)
	manualSupply := sdk.NewInt64Coin(bond.Token, 10)

	// Check that returns and fees are correct
	res, err = querier(ctx,
		[]string{keeper.QuerySellReturn, token, sellAmount.String()}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult.AdjustedSupply, manualSupply)
	require.Equal(t, queryResult.Returns, roundedReturns)
	require.Equal(t, queryResult.Returns, manualReturns)
	require.Equal(t, queryResult.TxFees, txFees)
	require.Equal(t, queryResult.TotalFees, totalFees)
	require.Equal(t, queryResult.TotalReturns, roundedTotalReturns)
}

func TestQuerySwapReturn(t *testing.T) {
	app, ctx := createTestApp(false)
	querier := keeper.NewQuerier(app.BondsKeeper)
	req := abci.RequestQuery{}
	var queryResult types.QuerySwapReturn

	// Initially error since no bond
	dummy1, dummy2, dummy3 := reserveToken, "100", reserveToken2
	res, err := querier(ctx,
		[]string{keeper.QuerySwapReturn, token, dummy1, dummy2, dummy3}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Add swapper bond
	bond := getValidSwapperBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Check that error since swapper not initialised
	res, err = querier(ctx,
		[]string{keeper.QuerySwapReturn, token, dummy1, dummy2, dummy3}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// Update current supply to 2
	bond.CurrentSupply = sdk.NewInt64Coin(bond.Token, 2)
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Send 200res,300rez to reserve
	newReserve := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 200),
		sdk.NewInt64Coin(reserveToken2, 300),
	)
	_, _ = app.BankKeeper.AddCoins(ctx, bond.ReserveAddress, newReserve)

	// Get swap return directly
	fromCoin := sdk.NewInt64Coin(reserveToken, 100)
	toToken := reserveToken2
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)
	swapReturns, txFee, _ := bond.GetReturnsForSwap(fromCoin, toToken, reserveBalances)

	// Calculate swap return manually
	// (since k = x.y = 200*300 = 60000 then if x becomes 300, the change in y,
	//  i.e. the returns must be 100rez. However, x in reality becomes 300-fee,
	//  so return is a bit less. This fee can be found to be 1 token. After
	//  rounding, the actual change in y turns out to be 99rez)
	manualSwapReturns := sdk.Coins{sdk.NewInt64Coin(reserveToken2, 99)}

	// Check that prices are correct
	res, err = querier(ctx, []string{keeper.QuerySwapReturn, token,
		fromCoin.Denom, fromCoin.Amount.String(), toToken}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	types.ModuleCdc.MustUnmarshalJSON(res, &queryResult)
	require.Equal(t, queryResult.TotalReturns, swapReturns)
	require.Equal(t, queryResult.TotalReturns, manualSwapReturns)
	require.Equal(t, queryResult.TotalFees, sdk.Coins{txFee})
}
