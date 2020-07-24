package simulation

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/ixoworld/bonds/x/bonds/internal/keeper"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateBond = "op_weight_msg_create_bond"
	OpWeightMsgEditBond   = "op_weight_msg_edit_bond"
	OpWeightMsgBuy        = "op_weight_msg_buy"
	OpWeightMsgSell       = "op_weight_msg_sell"
	OpWeightMsgSwap       = "op_weight_msg_swap"

	DefaultWeightMsgCreateBond = 5
	DefaultWeightMsgEditBond   = 5
	DefaultWeightMsgBuy        = 100
	DefaultWeightMsgSell       = 100
	DefaultWeightMsgSwap       = 100
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec,
	ak auth.AccountKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreateBond int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateBond, &weightMsgCreateBond, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBond = DefaultWeightMsgCreateBond
		},
	)

	var weightMsgEditBond int
	appParams.GetOrGenerate(cdc, OpWeightMsgEditBond, &weightMsgEditBond, nil,
		func(_ *rand.Rand) {
			weightMsgEditBond = DefaultWeightMsgEditBond
		},
	)

	var weightMsgBuy int
	appParams.GetOrGenerate(cdc, OpWeightMsgBuy, &weightMsgBuy, nil,
		func(_ *rand.Rand) {
			weightMsgBuy = DefaultWeightMsgBuy
		},
	)

	var weightMsgSell int
	appParams.GetOrGenerate(cdc, OpWeightMsgSell, &weightMsgSell, nil,
		func(_ *rand.Rand) {
			weightMsgSell = DefaultWeightMsgSell
		},
	)

	var weightMsgSwap int
	appParams.GetOrGenerate(cdc, OpWeightMsgSwap, &weightMsgSwap, nil,
		func(_ *rand.Rand) {
			weightMsgSwap = DefaultWeightMsgSwap
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateBond,
			SimulateMsgCreateBond(ak),
		),
		simulation.NewWeightedOperation(
			weightMsgEditBond,
			SimulateMsgEditBond(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBuy,
			SimulateMsgBuy(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSell,
			SimulateMsgSell(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSwap,
			SimulateMsgSwap(ak, k),
		),
	}
}

func SimulateMsgCreateBond(ak auth.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOpt []simulation.FutureOperation, err error) {

		if totalBondCount >= maxBondCount {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, accs)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)

		token := getNextBondName()
		name := getRandomNonEmptyString(r)
		desc := getRandomNonEmptyString(r)

		creator := address
		signers := []sdk.AccAddress{creator}

		var functionType string
		var reserveTokens []string
		randFunctionType := simulation.RandIntBetween(r, 0, 3)
		if randFunctionType == 0 {
			functionType = types.PowerFunction
			reserveTokens = defaultReserveTokens
		} else if randFunctionType == 1 {
			functionType = types.SigmoidFunction
			reserveTokens = defaultReserveTokens
		} else if randFunctionType == 2 {
			functionType = types.SwapperFunction
			reserveToken1, ok1 := getRandomBondName(r)
			reserveToken2, ok2 := getRandomBondNameExcept(r, reserveToken1)
			if !ok1 || !ok2 {
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}
			reserveTokens = []string{reserveToken1, reserveToken2}
		} else {
			panic("unexpected randFunctionType")
		}
		functionParameters := getRandomFunctionParameters(r, functionType)

		// Max fee is 100, so exit fee uses 100-txFee as max
		txFeePercentage := simulation.RandomDecAmount(r, sdk.NewDec(100))
		exitFeePercentage := simulation.RandomDecAmount(r, sdk.NewDec(100).Sub(txFeePercentage))

		// Since 100 is not allowed, a small number is subtracted from one of the fees
		if txFeePercentage.Add(exitFeePercentage).Equal(sdk.NewDec(100)) {
			if txFeePercentage.GT(sdk.ZeroDec()) {
				txFeePercentage = txFeePercentage.Sub(sdk.MustNewDecFromStr("0.000000000000000001"))
			} else {
				exitFeePercentage = exitFeePercentage.Sub(sdk.MustNewDecFromStr("0.000000000000000001"))
			}
		}

		// Addresses
		feeAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

		// Max supply, allow sells, batch blocks
		maxSupply := sdk.NewCoin(token, sdk.NewInt(int64(
			simulation.RandIntBetween(r, 1000000, 1000000000))))
		allowSells := getRandomAllowSellsValue(r)
		batchBlocks := sdk.NewUint(uint64(
			simulation.RandIntBetween(r, 1, 10)))

		msg := types.NewMsgCreateBond(token, name, desc, creator, functionType,
			functionParameters, reserveTokens, txFeePercentage,
			exitFeePercentage, feeAddress, maxSupply, blankOrderQuantityLimits,
			blankSanityRate, blankSanityMarginPercentage, allowSells, signers, batchBlocks)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil,
				fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		incrementBondCount() // since successfully created
		if msg.FunctionType == types.SwapperFunction {
			newSwapperBond(msg.Token)
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgEditBond(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOpt []simulation.FutureOperation, err error) {

		// Get random bond
		token, ok := getRandomBondName(r)
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		bond, found := k.GetBond(ctx, token)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		name := getRandomNonEmptyString(r)
		desc := getRandomNonEmptyString(r)

		simAccount, _ := simulation.FindAccount(accs, bond.Creator)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)

		editor := address
		signers := []sdk.AccAddress{editor}

		msg := types.NewMsgEditBond(token, name, desc,
			types.DoNotModifyField, types.DoNotModifyField,
			types.DoNotModifyField, editor, signers)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil,
				fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func getBuyIntoSwapper(r *rand.Rand, ctx sdk.Context, k keeper.Keeper,
	bond types.Bond, account exported.Account) (msg types.MsgBuy, err error, ok bool) {
	address := account.GetAddress()
	spendable := account.SpendableCoins(ctx.BlockTime())

	// Come up with max prices based on what is spendable
	spendableReserve1 := spendable.AmountOf(bond.ReserveTokens[0])
	spendableReserve2 := spendable.AmountOf(bond.ReserveTokens[1])
	maxPriceInt1, err := simulation.RandPositiveInt(r, spendableReserve1)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	maxPriceInt2, err := simulation.RandPositiveInt(r, spendableReserve2)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	maxPrices := sdk.NewCoins(
		sdk.NewCoin(bond.ReserveTokens[0], maxPriceInt1),
		sdk.NewCoin(bond.ReserveTokens[1], maxPriceInt2),
	)

	// Get lesser of max possible increase in supply and max order quantity
	var maxBuyAmount sdk.Int
	maxIncreaseInSupply := bond.MaxSupply.Sub(bond.CurrentSupply).Amount
	maxOrderQuantity := bond.OrderQuantityLimits.AmountOf(bond.Token)
	if maxOrderQuantity.IsZero() {
		maxBuyAmount = maxIncreaseInSupply
	} else {
		maxBuyAmount = sdk.MinInt(maxIncreaseInSupply, maxOrderQuantity)
	}

	if maxBuyAmount.IsZero() {
		return types.MsgBuy{}, nil, false
	}

	toBuyInt, err := simulation.RandPositiveInt(r, maxBuyAmount)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	amountToBuy := sdk.NewCoin(bond.Token, toBuyInt)

	// If not the first buy, create order and check if can afford
	if bond.CurrentSupply.IsPositive() {
		_, _, err = k.GetUpdatedBatchPricesAfterBuy(ctx, bond.Token,
			types.NewBuyOrder(address, amountToBuy, maxPrices))
		if err != nil {
			return types.MsgBuy{}, err, true
		}
	}

	return types.NewMsgBuy(address, amountToBuy, maxPrices), nil, true
}

func getBuyIntoPowerOrSigmoid(r *rand.Rand, ctx sdk.Context, k keeper.Keeper,
	bond types.Bond, account exported.Account) (msg types.MsgBuy, err error, ok bool) {
	address := account.GetAddress()
	spendable := account.SpendableCoins(ctx.BlockTime())

	// Come up with max price based on what is spendable
	spendableReserve := spendable.AmountOf(bond.ReserveTokens[0])
	maxPriceInt, err := simulation.RandPositiveInt(r, spendableReserve)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	maxPrices := sdk.Coins{sdk.NewCoin(bond.ReserveTokens[0], maxPriceInt)}

	// Get lesser of max possible increase in supply and max order quantity
	var maxBuyAmount sdk.Int
	maxIncreaseInSupply := bond.MaxSupply.Sub(bond.CurrentSupply).Amount
	maxOrderQuantity := bond.OrderQuantityLimits.AmountOf(bond.Token)
	if maxOrderQuantity.IsZero() {
		maxBuyAmount = maxIncreaseInSupply
	} else {
		maxBuyAmount = sdk.MinInt(maxIncreaseInSupply, maxOrderQuantity)
	}

	if maxBuyAmount.IsZero() {
		return types.MsgBuy{}, nil, false
	}

	toBuyInt, err := simulation.RandPositiveInt(r, maxBuyAmount)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	amountToBuy := sdk.NewCoin(bond.Token, toBuyInt)

	// Create order and check if can afford
	_, _, err = k.GetUpdatedBatchPricesAfterBuy(ctx, bond.Token,
		types.NewBuyOrder(address, amountToBuy, maxPrices))
	if err != nil {
		return types.MsgBuy{}, err, true
	}

	return types.NewMsgBuy(address, amountToBuy, maxPrices), nil, true
}

func SimulateMsgBuy(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// Get random bond
		token, ok := getRandomBondName(r)
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		bond, found := k.GetBond(ctx, token)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// Get accounts that have ALL the reserve tokens
		var filteredAccs []simulation.Account
		dummyNonZeroReserve := getDummyNonZeroReserve(bond.ReserveTokens)
		for _, a := range accs {
			coins := ak.GetAccount(ctx, a.Address).SpendableCoins(ctx.BlockTime())
			if dummyNonZeroReserve.DenomsSubsetOf(coins) {
				filteredAccs = append(filteredAccs, a)
			}
		}

		if len(filteredAccs) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, filteredAccs)
		account := ak.GetAccount(ctx, simAccount.Address)

		var msg types.MsgBuy
		if bond.FunctionType == types.SwapperFunction {
			msg, err, ok = getBuyIntoSwapper(r, ctx, k, bond, account)
		} else {
			msg, err, ok = getBuyIntoPowerOrSigmoid(r, ctx, k, bond, account)
		}

		// If ok, err is not something that should stop the simulation
		if (err != nil && ok) || (err == nil && !ok) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		} else if err != nil && !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		} else if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgSell(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// Get random bond
		token, ok := getRandomBondName(r)
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		bond, found := k.GetBond(ctx, token)
		if !found || bond.AllowSells == types.FALSE || bond.CurrentSupply.IsZero() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// Get accounts that have the token to be sold
		var filteredAccs []simulation.Account
		for _, a := range accs {
			coins := ak.GetAccount(ctx, a.Address).SpendableCoins(ctx.BlockTime())
			if coins.AmountOf(bond.Token).IsPositive() {
				filteredAccs = append(filteredAccs, a)
			}
		}

		if len(filteredAccs) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, filteredAccs)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)
		amount := account.SpendableCoins(ctx.BlockTime()).AmountOf(bond.Token)

		toSellInt, err := simulation.RandPositiveInt(r, amount)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		amountToSell := sdk.NewCoin(bond.Token, toSellInt)

		msg := types.NewMsgSell(address, amountToSell)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgSwap(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// Get swapper function bonds with some reserve
		var filteredBonds []string
		for _, sbToken := range swapperBonds {
			if !k.GetReserveBalances(ctx, sbToken).IsZero() {
				filteredBonds = append(filteredBonds, sbToken)
			}
		}

		if len(filteredBonds) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// Get random bond
		token := filteredBonds[simulation.RandIntBetween(r, 0, len(filteredBonds))]
		bond := k.MustGetBond(ctx, token)

		fromIndex := simulation.RandIntBetween(r, 0, 1)
		toIndex := 1 - fromIndex

		fromToken := bond.ReserveTokens[fromIndex]
		toToken := bond.ReserveTokens[toIndex]

		// Get accounts that have the token to be swapped
		var filteredAccs []simulation.Account
		for _, a := range accs {
			coins := ak.GetAccount(ctx, a.Address).SpendableCoins(ctx.BlockTime())
			if coins.AmountOf(fromToken).IsPositive() {
				filteredAccs = append(filteredAccs, a)
			}
		}

		if len(filteredAccs) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, filteredAccs)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)
		fromBalance := account.SpendableCoins(ctx.BlockTime()).AmountOf(fromToken)

		toSwapInt, err := simulation.RandPositiveInt(r, fromBalance)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		amountToSwap := sdk.NewCoin(fromToken, toSwapInt)

		msg := types.NewMsgSwap(address, token, amountToSwap, toToken)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
