package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestMoreEqualBuysSells(t *testing.T) {
	zero := sdk.NewInt64Coin("token", 0)
	one := sdk.NewInt64Coin("token", 1)
	batch := NewBatch("token", sdk.OneUint())
	testCases := []struct {
		buys      sdk.Coin
		sells     sdk.Coin
		moreBuys  bool
		moreSells bool
		equal     bool
	}{
		{zero, zero, false, false, true},
		{zero, one, false, true, false},
		{one, zero, true, false, false},
		{one, one, false, false, true},
	}
	for _, tc := range testCases {
		batch.TotalBuyAmount = tc.buys
		batch.TotalSellAmount = tc.sells
		require.Equal(t, tc.moreBuys, batch.MoreBuysThanSells())
		require.Equal(t, tc.moreSells, batch.MoreSellsThanBuys())
		require.Equal(t, tc.equal, batch.EqualBuysAndSells())
	}
}

func TestBaseOrderIsCancelled(t *testing.T) {
	testCases := []struct {
		order       BaseOrder
		valueToSet  bool
		isCancelled bool
	}{
		{NewBaseOrder(sdk.AccAddress{}, sdk.Coin{}), true, true},
		{NewBaseOrder(sdk.AccAddress{}, sdk.Coin{}), false, false},
	}
	for i, tc := range testCases {
		tc.order.Cancelled = tc.valueToSet
		require.Equal(t, tc.order.IsCancelled(), tc.isCancelled,
			"unexpected result for test case #%d, input: %v", i, tc.valueToSet)
	}
}

func TestNewBatchDefaultValues(t *testing.T) {
	token := "token"
	blocksRemaining := sdk.NewUint(100)
	batch := NewBatch(token, blocksRemaining)

	require.Equal(t, token, batch.Token)
	require.Equal(t, blocksRemaining, batch.BlocksRemaining)
	require.Equal(t, sdk.NewInt64Coin(token, 0), batch.TotalBuyAmount)
	require.Equal(t, sdk.NewInt64Coin(token, 0), batch.TotalSellAmount)
	require.Nil(t, batch.BuyPrices)
	require.Nil(t, batch.SellPrices)
	require.Nil(t, batch.Buys)
	require.Nil(t, batch.Sells)
	require.Nil(t, batch.Swaps)
}

func TestNewBaseOrderDefaultValues(t *testing.T) {
	address := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount := sdk.NewInt64Coin("token", 1000)
	order := NewBaseOrder(address, amount)

	require.Equal(t, address, order.Address)
	require.Equal(t, amount, order.Amount)
	require.False(t, order.Cancelled)
	require.Empty(t, order.CancelReason)
}

func TestNewBuyOrderDefaultValues(t *testing.T) {
	address := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount1 := sdk.NewInt64Coin("token1", 1000)
	amount2 := sdk.NewInt64Coin("token2", 2000)
	amount3 := sdk.NewInt64Coin("token3", 3000)
	maxPrices := sdk.NewCoins(amount2, amount3)
	order := NewBuyOrder(address, amount1, maxPrices)

	require.Equal(t, address, order.Address)
	require.Equal(t, amount1, order.Amount)
	require.False(t, order.Cancelled)
	require.Empty(t, order.CancelReason)
	require.Equal(t, maxPrices, order.MaxPrices)
}

func TestNewSellOrderDefaultValues(t *testing.T) {
	address := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount := sdk.NewInt64Coin("token", 1000)
	order := NewSellOrder(address, amount)

	require.Equal(t, address, order.Address)
	require.Equal(t, amount, order.Amount)
	require.False(t, order.Cancelled)
	require.Empty(t, order.CancelReason)
}

func TestNewSwapOrderDefaultValues(t *testing.T) {
	address := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	fromAmount := sdk.NewInt64Coin("token1", 1000)
	toToken := "token2"
	order := NewSwapOrder(address, fromAmount, toToken)

	require.Equal(t, address, order.Address)
	require.Equal(t, fromAmount, order.Amount)
	require.Equal(t, toToken, order.ToToken)
	require.False(t, order.Cancelled)
	require.Empty(t, order.CancelReason)
}
