//nolint
package simapp

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	tmkv "github.com/tendermint/tendermint/libs/kv"
	"io/ioutil"
)

// SimulationOperations retrieves the simulation params from the provided file path
// and returns all the modules weighted operations
func SimulationOperations(app *SimApp, cdc *codec.Codec, config simulation.Config) []simulation.WeightedOperation {
	simState := module.SimulationState{
		AppParams: make(simulation.AppParams),
		Cdc:       cdc,
	}

	if config.ParamsFile != "" {
		bz, err := ioutil.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}

		app.cdc.MustUnmarshalJSON(bz, &simState.AppParams)
	}

	simState.ParamChanges = app.sm.GenerateParamChanges(config.Seed)
	simState.Contents = app.sm.GetProposalContents(simState)
	return app.sm.WeightedOperations(simState)
}

//---------------------------------------------------------------------
// Simulation Utils

// ExportStateToJSON util function to export the app state to JSON
func ExportStateToJSON(app *SimApp, path string) error {
	fmt.Println("exporting app state...")
	appState, _, err := app.ExportAppStateAndValidators(false, nil)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(appState), 0644)
}

// ExportParamsToJSON util function to export the simulation parameters to JSON
func ExportParamsToJSON(params simulation.Params, path string) error {
	fmt.Println("exporting simulation params...")
	paramsBz, err := json.MarshalIndent(params, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, paramsBz, 0644)
}

// GetSimulationLog unmarshals the KVPair's Value to the corresponding type based on the
// each's module store key and the prefix bytes of the KVPair's key.
func GetSimulationLog(storeName string, sdr sdk.StoreDecoderRegistry, cdc *codec.Codec, kvAs, kvBs []tmkv.Pair) (log string) {
	for i := 0; i < len(kvAs); i++ {

		if len(kvAs[i].Value) == 0 && len(kvBs[i].Value) == 0 {
			// skip if the value doesn't have any bytes
			continue
		}

		decoder, ok := sdr[storeName]
		if ok {
			log += decoder(cdc, kvAs[i], kvBs[i])
		} else {
			log += fmt.Sprintf("store A %X => %X\nstore B %X => %X\n", kvAs[i].Key, kvAs[i].Value, kvBs[i].Key, kvBs[i].Value)
		}
	}

	return
}
