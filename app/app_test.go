package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBondsdExport(t *testing.T) {
	memDB := db.NewMemDB()
	app := NewBondsApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), memDB, nil, true, map[int64]bool{}, 0)
	err := setGenesis(app)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newApp := NewBondsApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), memDB, nil, true, map[int64]bool{}, 0)
	_, _, err = newApp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	memDB := db.NewMemDB()
	app := NewBondsApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), memDB, nil, true, map[int64]bool{}, 0)

	for acc := range maccPerms {
		blacklisted := !allowedReceivingModAcc[acc]
		require.Equal(t, blacklisted, app.BankKeeper.BlacklistedAddr(app.SupplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(app *BondsApp) error {
	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	app.Commit()
	return nil
}
