package simapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *SimApp {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}
