package bonds_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ixoworld/bonds/x/bonds"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestInitAndExportGenesis(t *testing.T) {
	app, ctx := createTestApp(false)
	genesisState := bonds.DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Bonds))
	require.Equal(t, 0, len(genesisState.Batches))

	token := "testtoken"
	name := "test token"
	description := "this is a test token"
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	functionType := types.PowerFunction
	functionParameters := types.FunctionParams{
		types.NewFunctionParam("m", sdk.NewDec(12)),
		types.NewFunctionParam("n", sdk.NewDec(2)),
		types.NewFunctionParam("c", sdk.NewDec(100))}
	reserveTokens := []string{"reservetoken"}
	reserveAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	txFeePercentage := sdk.MustNewDecFromStr("0.1")
	exitFeePercentage := sdk.MustNewDecFromStr("0.2")
	feeAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	maxSupply := sdk.NewInt64Coin(token, 10000)
	orderQuantityLimits := sdk.NewCoins(
		sdk.NewInt64Coin("token1", 1),
		sdk.NewInt64Coin("token2", 2),
		sdk.NewInt64Coin("token3", 3),
	)
	sanityRate := sdk.MustNewDecFromStr("0.3")
	sanityMarginPercentage := sdk.MustNewDecFromStr("0.4")
	allowSell := "true"
	signers := []sdk.AccAddress{creator}
	batchBlocks := sdk.NewUint(10)

	bond := types.NewBond(token, name, description, creator,
		functionType, functionParameters, reserveTokens,
		reserveAddress, txFeePercentage, exitFeePercentage,
		feeAddress, maxSupply, orderQuantityLimits, sanityRate,
		sanityMarginPercentage, allowSell, signers, batchBlocks)
	batch := types.NewBatch(bond.Token, bond.BatchBlocks)

	genesisState = bonds.NewGenesisState(
		[]types.Bond{bond}, []types.Batch{batch})

	bonds.InitGenesis(ctx, app.BondsKeeper, genesisState)

	returnedBond := app.BondsKeeper.MustGetBond(ctx, token)
	require.Equal(t, bond, returnedBond)

	returnedBatch := app.BondsKeeper.MustGetBatch(ctx, token)
	require.Equal(t, batch, returnedBatch)

	exportedGenesisState := bonds.ExportGenesis(ctx, app.BondsKeeper)
	require.Equal(t, genesisState.Bonds, exportedGenesisState.Bonds)
	require.Equal(t, genesisState.Batches, exportedGenesisState.Batches)
}
