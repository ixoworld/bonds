package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	BankKeeper    bank.Keeper
	SupplyKeeper  supply.Keeper
	accountKeeper auth.AccountKeeper
	StakingKeeper staking.Keeper

	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	cdc *codec.Codec
}

func NewKeeper(bankKeeper bank.Keeper, supplyKeeper supply.Keeper,
	accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper,
	storeKey sdk.StoreKey, paramSpace params.Subspace,
	cdc *codec.Codec) Keeper {

	// ensure batches module account is set
	if addr := supplyKeeper.GetModuleAddress(types.BatchesIntermediaryAccount); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BatchesIntermediaryAccount))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		BankKeeper:    bankKeeper,
		SupplyKeeper:  supplyKeeper,
		accountKeeper: accountKeeper,
		StakingKeeper: stakingKeeper,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		cdc:           cdc,
	}
}

// GetParams returns the total set of bonds parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the bonds parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
