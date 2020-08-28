package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ixoworld/bonds/types"
)

// Init initializes the application, overriding the default genesis states that should be changed
func Init() {
	mint.DefaultGenesisState = mintGenesisState
	staking.DefaultGenesisState = stakingGenesisState
	gov.DefaultGenesisState = govGenesisState
}

// stakingGenesisState returns the default genesis state for the staking module, replacing the
// bond denom from stake to ubtsg
func stakingGenesisState() staking.GenesisState {
	// Get default staking genesis state and set staking bond denom to BondDenom
	state := staking.DefaultGenesisState()
	state.Params.BondDenom = types.BondDenom
	return state
}

func govGenesisState() gov.GenesisState {
	// Get default gov genesis state and set deposit denom to BondDenom
	state := gov.DefaultGenesisState()
	state.DepositParams = gov.NewDepositParams(
		sdk.NewCoins(sdk.NewCoin(types.BondDenom, govTypes.DefaultMinDepositTokens)),
		gov.DefaultPeriod,
	)
	return state
}

func mintGenesisState() mint.GenesisState {
	// Get default mint genesis state and set mint denom to BondDenom
	state := mint.DefaultGenesisState()
	state.Params.MintDenom = types.BondDenom
	return state
}
