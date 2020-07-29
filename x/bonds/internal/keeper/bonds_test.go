package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBondExistsSetGet(t *testing.T) {
	app, ctx := createTestApp(false)

	// Try to get bond
	_, found := app.BondsKeeper.GetBond(ctx, token)

	// Bond doesn't exist yet
	require.False(t, found)
	require.False(t, app.BondsKeeper.BondExists(ctx, token))

	// Add bond
	bondAdded := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bondAdded)

	// Bond now exists
	require.True(t, app.BondsKeeper.BondExists(ctx, token))

	// Option 1: get bond
	bondFetched1, found := app.BondsKeeper.GetBond(ctx, token)
	// Option 2: must get bond
	bondFetched2 := app.BondsKeeper.MustGetBond(ctx, token)
	// Option 2: must get bond
	bondFetched3 := app.BondsKeeper.MustGetBondByKey(ctx, types.GetBondKey(token))

	// Batch fetched is equal to added batch
	require.Equal(t, bondAdded, bondFetched1)
	require.Equal(t, bondAdded, bondFetched2)
	require.Equal(t, bondAdded, bondFetched3)
	require.True(t, found)
}

func TestGetReserveBalances(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Reserve is initially empty
	require.True(t, app.BondsKeeper.GetReserveBalances(ctx, token).IsZero())

	// Send coins to reserve address
	reserveCoins, _ := sdk.ParseCoins("12.34res1,56.78res2")
	_, _ = app.BankKeeper.AddCoins(ctx, bond.ReserveAddress, reserveCoins)

	// Reserve now equal to amount sent
	reserveBalances := app.BondsKeeper.GetReserveBalances(ctx, token)
	require.Equal(t, reserveCoins, reserveBalances)
}

func TestGetSupplyAdjustedForBuy(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add bond and batch
	bond := getValidBond()
	batch := getValidBatch()
	app.BondsKeeper.SetBond(ctx, token, bond)
	app.BondsKeeper.SetBatch(ctx, token, batch)

	// Supply is initially zero
	require.True(t, app.BondsKeeper.GetSupplyAdjustedForBuy(ctx, token).IsZero())

	// Increase current supply
	increaseInSupply := sdk.NewInt64Coin(token, 100)
	bond.CurrentSupply = increaseInSupply
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Supply has increased
	supply := app.BondsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, increaseInSupply, supply)

	// Increase supply by adding a buy order
	increaseDueToOrder := sdk.NewInt64Coin(token, 11)
	buyOrder := getValidBuyOrder()
	buyOrder.Amount = increaseDueToOrder
	app.BondsKeeper.AddBuyOrder(ctx, token, buyOrder, nil, nil)

	// Supply has increased
	expectedSupply := increaseInSupply.Add(increaseDueToOrder)
	supply = app.BondsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding sell order does not affect supply
	sellOrder := getValidSellOrder()
	sellOrder.Amount = sdk.NewInt64Coin(token, 100)
	app.BondsKeeper.AddSellOrder(ctx, token, sellOrder, nil, nil)

	// Supply has not increased
	supply = app.BondsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding swap order does not affect supply
	app.BondsKeeper.AddSwapOrder(ctx, token, getValidSwapOrder())

	// Supply has not increased
	supply = app.BondsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, expectedSupply, supply)
}

func TestGetSupplyAdjustedForSell(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add bond and batch
	bond := getValidBond()
	batch := getValidBatch()
	app.BondsKeeper.SetBond(ctx, token, bond)
	app.BondsKeeper.SetBatch(ctx, token, batch)

	// Supply is initially zero
	require.True(t, app.BondsKeeper.GetSupplyAdjustedForSell(ctx, token).IsZero())

	// Increase current supply
	increaseInSupply := sdk.NewInt64Coin(token, 100)
	bond.CurrentSupply = increaseInSupply
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Supply has increased
	supply := app.BondsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, increaseInSupply, supply)

	// Decrease supply by adding a sell order
	decreaseDueToOrder := sdk.NewInt64Coin(token, 11)
	sellOrder := getValidSellOrder()
	sellOrder.Amount = decreaseDueToOrder
	app.BondsKeeper.AddSellOrder(ctx, token, sellOrder, nil, nil)

	// Supply has decreased
	expectedSupply := increaseInSupply.Sub(decreaseDueToOrder)
	supply = app.BondsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding buy order does not affect supply
	buyOrder := getValidBuyOrder()
	buyOrder.Amount = sdk.NewInt64Coin(token, 100)
	app.BondsKeeper.AddBuyOrder(ctx, token, buyOrder, nil, nil)

	// Supply has not increased
	supply = app.BondsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding swap order does not affect supply
	app.BondsKeeper.AddSwapOrder(ctx, token, getValidSwapOrder())

	// Supply has not increased
	supply = app.BondsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, expectedSupply, supply)
}

func TestSetCurrentSupply(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// Supply is initially zero
	require.True(t, app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply.IsZero())

	// Change supply
	newSupply := sdk.NewInt64Coin(token, 100)
	app.BondsKeeper.SetCurrentSupply(ctx, token, newSupply)

	// Check that supply changed
	supplyFetched := app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.Equal(t, newSupply, supplyFetched)

	// Change supply again
	newSupply = sdk.NewInt64Coin(token, 50)
	app.BondsKeeper.SetCurrentSupply(ctx, token, newSupply)

	// Check that supply changed
	supplyFetched = app.BondsKeeper.MustGetBond(ctx, token).CurrentSupply
	require.Equal(t, newSupply, supplyFetched)
}

func TestSetBondState(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add bond
	bond := getValidBond()
	app.BondsKeeper.SetBond(ctx, token, bond)

	// State is initially "initState"
	require.Equal(t, initState, app.BondsKeeper.MustGetBond(ctx, token).State)

	// Change state
	newState := "some_other_state"
	app.BondsKeeper.SetBondState(ctx, token, newState)

	// Check that state changed
	stateFetched := app.BondsKeeper.MustGetBond(ctx, token).State
	require.Equal(t, newState, stateFetched)

	// Change supply again
	newState = "yet another state"
	app.BondsKeeper.SetBondState(ctx, token, newState)

	// Check that supply changed
	stateFetched = app.BondsKeeper.MustGetBond(ctx, token).State
	require.Equal(t, newState, stateFetched)
}
