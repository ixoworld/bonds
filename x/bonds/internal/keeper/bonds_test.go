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

func TestGetNumberOfBonds(t *testing.T) {
	app, ctx := createTestApp(false)

	// No bond exists yet
	require.Equal(t, sdk.ZeroInt(), app.BondsKeeper.GetNumberOfBonds(ctx))

	// Add bond 1/3
	app.BondsKeeper.SetBond(ctx, token1, getValidBondWithToken(token1))
	require.Equal(t, sdk.NewInt(1), app.BondsKeeper.GetNumberOfBonds(ctx))

	// Add bond 2/3
	app.BondsKeeper.SetBond(ctx, token2, getValidBondWithToken(token2))
	require.Equal(t, sdk.NewInt(2), app.BondsKeeper.GetNumberOfBonds(ctx))

	// Add bond 3/3
	app.BondsKeeper.SetBond(ctx, token3, getValidBondWithToken(token3))
	require.Equal(t, sdk.NewInt(3), app.BondsKeeper.GetNumberOfBonds(ctx))
}

func TestGetReserveAddressByBondCount(t *testing.T) {
	app, _ := createTestApp(false)
	const maxInt64 = int64(^uint64(0) >> 1) // 9223372036854775807

	testCases := []struct {
		input        sdk.Int
		truncatedHex string
	}{
		{sdk.NewInt(0), "7B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A0"},
		{sdk.NewInt(5), "7B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A5"},
		{sdk.NewInt(10), "B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A10"},
		{sdk.NewInt(50), "B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A50"},
		{sdk.NewInt(1000000000), "AF4A1D3EC729DC422C6341BAEEDC9A1000000000"},
		{sdk.NewInt(5000000000), "AF4A1D3EC729DC422C6341BAEEDC9A5000000000"},
		{sdk.NewInt(maxInt64), "729DC422C6341BAEEDC9A9223372036854775807"},
	}

	for _, tc := range testCases {
		expectedAddr, err := sdk.AccAddressFromHex(tc.truncatedHex)
		require.Nil(t, err)

		addr := app.BondsKeeper.GetReserveAddressByBondCount(tc.input)
		require.Equal(t, expectedAddr, addr)
	}
}

func TestGetNextUnusedReserveAddress(t *testing.T) {
	app, ctx := createTestApp(false)

	testCases := []struct {
		truncatedHex  string
		nextBondToken string
	}{
		{"7B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A0", token1},
		{"7B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A1", token2},
		{"7B2E13A94AF4A1D3EC729DC422C6341BAEEDC9A2", token3},
	}

	for _, tc := range testCases {
		expectedAddr, err := sdk.AccAddressFromHex(tc.truncatedHex)
		require.Nil(t, err)

		addr := app.BondsKeeper.GetNextUnusedReserveAddress(ctx)
		require.Equal(t, expectedAddr, addr)

		app.BondsKeeper.SetBond(ctx, tc.nextBondToken, getValidBondWithToken(tc.nextBondToken))
	}
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
