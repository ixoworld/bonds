package bonds

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/ixoworld/bonds/x/bonds/internal/keeper"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
	"strings"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCreateBond:
			return handleMsgCreateBond(ctx, keeper, msg)
		case types.MsgEditBond:
			return handleMsgEditBond(ctx, keeper, msg)
		case types.MsgBuy:
			return handleMsgBuy(ctx, keeper, msg)
		case types.MsgSell:
			return handleMsgSell(ctx, keeper, msg)
		case types.MsgSwap:
			return handleMsgSwap(ctx, keeper, msg)
		case types.MsgMakeOutcomePayment:
			return handleMsgMakeOutcomePayment(ctx, keeper, msg)
		case types.MsgWithdrawShare:
			return handleMsgWithdrawShare(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized bonds Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) []abci.ValidatorUpdate {

	iterator := keeper.GetBondIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		bond := keeper.MustGetBondByKey(ctx, iterator.Key())
		batch := keeper.MustGetBatch(ctx, bond.Token)

		// Subtract one block
		batch.BlocksRemaining = batch.BlocksRemaining.SubUint64(1)
		keeper.SetBatch(ctx, bond.Token, batch)

		// If blocks remaining > 0 do not perform orders
		if !batch.BlocksRemaining.IsZero() {
			continue
		}

		// Perform orders
		keeper.PerformOrders(ctx, bond.Token)

		// Get bond again just in case current supply was updated
		// Get batch again just in case orders were cancelled
		bond = keeper.MustGetBond(ctx, bond.Token)
		batch = keeper.MustGetBatch(ctx, bond.Token)

		// For augmented, if hatch phase and newSupply >= S0, go to open phase
		if bond.FunctionType == types.AugmentedFunction &&
			bond.State == types.HatchState {
			args := bond.FunctionParameters.AsMap()
			if sdk.NewDecFromInt(bond.CurrentSupply.Amount).GTE(args["S0"]) {
				keeper.SetBondState(ctx, bond.Token, types.OpenState)
				bond = keeper.MustGetBond(ctx, bond.Token) // get bond again
				bond.AllowSells = true                     // enable sells
				keeper.SetBond(ctx, bond.Token, bond)      // update bond
			}
		}

		// Save current batch as last batch and reset current batch
		keeper.SetLastBatch(ctx, bond.Token, batch)
		keeper.SetBatch(ctx, bond.Token, types.NewBatch(bond.Token, bond.BatchBlocks))
	}
	return []abci.ValidatorUpdate{}
}

func handleMsgCreateBond(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateBond) sdk.Result {
	if keeper.BankKeeper.BlacklistedAddr(msg.FeeAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.FeeAddress)).Result()
	}

	if keeper.BondExists(ctx, msg.Token) {
		return types.ErrBondAlreadyExists(DefaultCodespace, msg.Token).Result()
	} else if msg.Token == keeper.StakingKeeper.GetParams(ctx).BondDenom {
		return types.ErrBondTokenCannotBeStakingToken(DefaultCodespace).Result()
	}

	// Get reserve address
	reserveAddress := supply.NewModuleAddress(
		fmt.Sprintf("bonds/%s/reserveAddress", msg.Token))

	// Set state to open by default (overridden below if augmented function)
	state := types.OpenState

	// TODO: investigate ways to prevent reserve address from receiving transactions

	// Not critical since as is no tokens can be taken out of the reserve, unless
	// programmatically. However, increases in balance still affect calculations.
	// Two possible solutions are (i) add new reserve addresses to the bank module
	// blacklisted addresses (but no guarantee that this will be sufficient), or
	// (ii) use a global res. address and store (in the bond) the share of the pool.

	// If augmented, add R0, S0, V0 as parameters for quick access
	// Also, override AllowSells and set to False if S0 > 0
	if msg.FunctionType == types.AugmentedFunction {
		paramsMap := msg.FunctionParameters.AsMap()
		d0, _ := paramsMap["d0"]
		p0, _ := paramsMap["p0"]
		theta, _ := paramsMap["theta"]
		kappa, _ := paramsMap["kappa"]

		R0 := d0.Mul(sdk.OneDec().Sub(theta))
		S0 := d0.Quo(p0)
		V0 := types.Invariant(R0, S0, kappa.TruncateInt64())
		// TODO: consider calculating these on-the-fly, especially R0 and S0

		msg.FunctionParameters = append(msg.FunctionParameters,
			types.FunctionParams{
				types.NewFunctionParam("R0", R0),
				types.NewFunctionParam("S0", S0),
				types.NewFunctionParam("V0", V0),
			}...)

		// Set state to Hatch and disable sells. Note that it is never the case
		// that we start with OpenState because S0>0, since S0=d0/p0 and d0>0
		state = types.HatchState
		msg.AllowSells = false
	}

	bond := types.NewBond(msg.Token, msg.Name, msg.Description, msg.Creator,
		msg.FunctionType, msg.FunctionParameters, msg.ReserveTokens,
		reserveAddress, msg.TxFeePercentage, msg.ExitFeePercentage,
		msg.FeeAddress, msg.MaxSupply, msg.OrderQuantityLimits, msg.SanityRate,
		msg.SanityMarginPercentage, msg.AllowSells, msg.Signers,
		msg.BatchBlocks, state)

	keeper.SetBond(ctx, msg.Token, bond)
	keeper.SetBatch(ctx, msg.Token, types.NewBatch(bond.Token, msg.BatchBlocks))

	logger := keeper.Logger(ctx)
	logger.Info(fmt.Sprintf("bond %s with reserve(s) [%s] created by %s",
		msg.Token, strings.Join(bond.ReserveTokens, ","), msg.Creator.String()))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateBond,
			sdk.NewAttribute(types.AttributeKeyBond, msg.Token),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
			sdk.NewAttribute(types.AttributeKeyFunctionType, msg.FunctionType),
			sdk.NewAttribute(types.AttributeKeyFunctionParameters, msg.FunctionParameters.String()),
			sdk.NewAttribute(types.AttributeKeyReserveTokens, types.StringsToString(msg.ReserveTokens)),
			sdk.NewAttribute(types.AttributeKeyReserveAddress, reserveAddress.String()),
			sdk.NewAttribute(types.AttributeKeyTxFeePercentage, msg.TxFeePercentage.String()),
			sdk.NewAttribute(types.AttributeKeyExitFeePercentage, msg.ExitFeePercentage.String()),
			sdk.NewAttribute(types.AttributeKeyFeeAddress, msg.FeeAddress.String()),
			sdk.NewAttribute(types.AttributeKeyMaxSupply, msg.MaxSupply.String()),
			sdk.NewAttribute(types.AttributeKeyOrderQuantityLimits, msg.OrderQuantityLimits.String()),
			sdk.NewAttribute(types.AttributeKeySanityRate, msg.SanityRate.String()),
			sdk.NewAttribute(types.AttributeKeySanityMarginPercentage, msg.SanityMarginPercentage.String()),
			sdk.NewAttribute(types.AttributeKeyAllowSells, strconv.FormatBool(msg.AllowSells)),
			sdk.NewAttribute(types.AttributeKeySigners, types.AccAddressesToString(msg.Signers)),
			sdk.NewAttribute(types.AttributeKeyBatchBlocks, msg.BatchBlocks.String()),
			sdk.NewAttribute(types.AttributeKeyState, state),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgEditBond(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgEditBond) sdk.Result {

	bond, found := keeper.GetBond(ctx, msg.Token)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, msg.Token).Result()
	}

	if !bond.SignersEqualTo(msg.Signers) {
		errMsg := fmt.Sprintf("List of signers does not match the one in the bond")
		return sdk.ErrInternal(errMsg).Result()
	}

	if msg.Name != types.DoNotModifyField {
		bond.Name = msg.Name
	}
	if msg.Description != types.DoNotModifyField {
		bond.Description = msg.Description
	}

	if msg.OrderQuantityLimits != types.DoNotModifyField {
		orderQuantityLimits, err := sdk.ParseCoins(msg.OrderQuantityLimits)
		if err != nil {
			return sdk.ErrInternal(err.Error()).Result()
		}
		bond.OrderQuantityLimits = orderQuantityLimits
	}

	if msg.SanityRate != types.DoNotModifyField {
		var sanityRate, sanityMarginPercentage sdk.Dec
		if msg.SanityRate == "" {
			sanityRate = sdk.ZeroDec()
			sanityMarginPercentage = sdk.ZeroDec()
		} else {
			parsedSanityRate, err := sdk.NewDecFromStr(msg.SanityRate)
			if err != nil {
				return types.ErrArgumentMissingOrNonFloat(types.DefaultCodespace, "sanity rate").Result()
			} else if parsedSanityRate.IsNegative() {
				return types.ErrArgumentCannotBeNegative(types.DefaultCodespace, "sanity rate").Result()
			}
			parsedSanityMarginPercentage, err := sdk.NewDecFromStr(msg.SanityMarginPercentage)
			if err != nil {
				return types.ErrArgumentMissingOrNonFloat(types.DefaultCodespace, "sanity margin percentage").Result()
			} else if parsedSanityMarginPercentage.IsNegative() {
				return types.ErrArgumentCannotBeNegative(types.DefaultCodespace, "sanity margin percentage").Result()
			}
			sanityRate = parsedSanityRate
			sanityMarginPercentage = parsedSanityMarginPercentage
		}
		bond.SanityRate = sanityRate
		bond.SanityMarginPercentage = sanityMarginPercentage
	}

	logger := keeper.Logger(ctx)
	logger.Info(fmt.Sprintf("bond %s edited by %s",
		msg.Token, msg.Editor.String()))

	keeper.SetBond(ctx, msg.Token, bond)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditBond,
			sdk.NewAttribute(types.AttributeKeyBond, msg.Token),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
			sdk.NewAttribute(types.AttributeKeyOrderQuantityLimits, msg.OrderQuantityLimits),
			sdk.NewAttribute(types.AttributeKeySanityRate, msg.SanityRate),
			sdk.NewAttribute(types.AttributeKeySanityMarginPercentage, msg.SanityMarginPercentage),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Editor.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBuy(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBuy) sdk.Result {

	token := msg.Amount.Denom
	bond, found := keeper.GetBond(ctx, token)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, token).Result()
	}

	// Check max prices
	if !bond.ReserveDenomsEqualTo(msg.MaxPrices) {
		return types.ErrReserveDenomsMismatch(types.DefaultCodespace, msg.MaxPrices.String(), bond.ReserveTokens).Result()
	}

	// Check if order quantity limit exceeded
	if bond.AnyOrderQuantityLimitsExceeded(sdk.Coins{msg.Amount}) {
		return types.ErrOrderQuantityLimitExceeded(types.DefaultCodespace).Result()
	}

	// For the swapper, the first buy is the initialisation of the reserves
	// The max prices are used as the actual prices and one token is minted
	// The amount of token serves to define the price of adding more liquidity
	if bond.CurrentSupply.IsZero() && bond.FunctionType == types.SwapperFunction {
		return performFirstSwapperFunctionBuy(ctx, keeper, msg)
	}

	// Take max that buyer is willing to pay (enforces maxPrice <= balance)
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Buyer,
		types.BatchesIntermediaryAccount, msg.MaxPrices)
	if err != nil {
		return err.Result()
	}

	// Create order
	order := types.NewBuyOrder(msg.Buyer, msg.Amount, msg.MaxPrices)

	// Get buy price and check if can add buy order to batch
	buyPrices, sellPrices, err := keeper.GetUpdatedBatchPricesAfterBuy(ctx, token, order)
	if err != nil {
		return err.Result()
	}

	// Add buy order to batch
	keeper.AddBuyOrder(ctx, token, order, buyPrices, sellPrices)

	// Cancel unfulfillable orders
	keeper.CancelUnfulfillableOrders(ctx, token)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBuy,
			sdk.NewAttribute(types.AttributeKeyBond, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyMaxPrices, msg.MaxPrices.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func performFirstSwapperFunctionBuy(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBuy) sdk.Result {

	// TODO: investigate effect that a high amount has on future buyers' ability to buy.

	token := msg.Amount.Denom
	bond, found := keeper.GetBond(ctx, token)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, token).Result()
	}

	// Check if initial liquidity violates sanity rate
	if bond.ReservesViolateSanityRate(msg.MaxPrices) {
		return types.ErrValuesViolateSanityRate(types.DefaultCodespace).Result()
	}

	// Use max prices as the amount to send to the liquidity pool (i.e. price)
	err := keeper.BankKeeper.SendCoins(ctx, msg.Buyer, bond.ReserveAddress, msg.MaxPrices)
	if err != nil {
		return err.Result()
	}

	// Mint bond tokens
	err = keeper.SupplyKeeper.MintCoins(ctx, types.BondsMintBurnAccount,
		sdk.Coins{msg.Amount})
	if err != nil {
		return err.Result()
	}

	// Send bond tokens to buyer
	err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
		types.BondsMintBurnAccount, msg.Buyer, sdk.Coins{msg.Amount})
	if err != nil {
		return err.Result()
	}

	// Update supply
	keeper.SetCurrentSupply(ctx, bond.Token, bond.CurrentSupply.Add(msg.Amount))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeInitSwapper,
			sdk.NewAttribute(types.AttributeKeyBond, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyChargedPrices, msg.MaxPrices.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSell(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSell) sdk.Result {

	token := msg.Amount.Denom
	bond, found := keeper.GetBond(ctx, token)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, token).Result()
	}

	if !bond.AllowSells {
		return types.ErrBondDoesNotAllowSelling(types.DefaultCodespace).Result()
	}

	// Check if order quantity limit exceeded
	if bond.AnyOrderQuantityLimitsExceeded(sdk.Coins{msg.Amount}) {
		return types.ErrOrderQuantityLimitExceeded(types.DefaultCodespace).Result()
	}

	// Send coins to be burned from seller (enforces sellAmount <= balance)
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Seller,
		types.BondsMintBurnAccount, sdk.Coins{msg.Amount})
	if err != nil {
		return err.Result()
	}

	// Burn bond tokens to be sold
	err = keeper.SupplyKeeper.BurnCoins(ctx, types.BondsMintBurnAccount,
		sdk.Coins{msg.Amount})
	if err != nil {
		return err.Result()
	}

	// Create order
	order := types.NewSellOrder(msg.Seller, msg.Amount)

	// Get sell price and check if can add sell order to batch
	buyPrices, sellPrices, err := keeper.GetUpdatedBatchPricesAfterSell(ctx, token, order)
	if err != nil {
		return err.Result()
	}

	// Add sell order to batch
	keeper.AddSellOrder(ctx, token, order, buyPrices, sellPrices)

	//// Cancel unfulfillable orders (Note: no need)
	//keeper.CancelUnfulfillableOrders(ctx, token)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSell,
			sdk.NewAttribute(types.AttributeKeyBond, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Seller.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSwap(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSwap) sdk.Result {

	bond, found := keeper.GetBond(ctx, msg.BondToken)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, msg.BondToken).Result()
	}

	// Confirm that function type is swapper_function
	if bond.FunctionType != types.SwapperFunction {
		return types.ErrFunctionNotAvailableForFunctionType(types.DefaultCodespace).Result()
	}

	// Check that from and to use reserve token names
	fromAndTo := sdk.NewCoins(msg.From, sdk.NewCoin(msg.ToToken, sdk.OneInt()))
	fromAndToDenoms := msg.From.Denom + "," + msg.ToToken
	if !bond.ReserveDenomsEqualTo(fromAndTo) {
		return types.ErrReserveDenomsMismatch(types.DefaultCodespace, fromAndToDenoms, bond.ReserveTokens).Result()
	}

	// Check if order quantity limit exceeded
	if bond.AnyOrderQuantityLimitsExceeded(sdk.Coins{msg.From}) {
		return types.ErrOrderQuantityLimitExceeded(types.DefaultCodespace).Result()
	}

	// Take coins to be swapped from swapper (enforces swapAmount <= balance)
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Swapper,
		types.BatchesIntermediaryAccount, sdk.Coins{msg.From})
	if err != nil {
		return err.Result()
	}

	// Create order
	order := types.NewSwapOrder(msg.Swapper, msg.From, msg.ToToken)

	// Add swap order to batch
	keeper.AddSwapOrder(ctx, msg.BondToken, order)

	//// Cancel unfulfillable orders (Note: no need)
	//keeper.CancelUnfulfillableOrders(ctx, token)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeyBond, msg.BondToken),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.From.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySwapFromToken, msg.From.Denom),
			sdk.NewAttribute(types.AttributeKeySwapToToken, msg.ToToken),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Swapper.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgMakeOutcomePayment(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgMakeOutcomePayment) sdk.Result {

	bond, found := keeper.GetBond(ctx, msg.BondToken)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, msg.BondToken).Result()
	}

	// Confirm that function type is augmented_function and that state is OPEN
	if bond.FunctionType != types.AugmentedFunction {
		return types.ErrFunctionNotAvailableForFunctionType(types.DefaultCodespace).Result()
	} else if bond.State != types.OpenState {
		return types.ErrInvalidNextState(types.DefaultCodespace).Result()
	}

	// Send outcome payment to reserve address
	// TODO: amount should not be hard-coded
	outcomePayment := sdk.NewCoins(sdk.NewCoin("res", sdk.NewInt(100000)))
	err := keeper.BankKeeper.SendCoins(
		ctx, msg.Sender, bond.ReserveAddress, outcomePayment)
	if err != nil {
		return err.Result()
	}

	// Set bond state to SETTLE and save bond
	bond.State = types.SettleState
	keeper.SetBond(ctx, bond.Token, bond)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOutcomePayment,
			sdk.NewAttribute(types.AttributeKeyBond, msg.BondToken),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Sender.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, outcomePayment.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawShare(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgWithdrawShare) sdk.Result {

	bond, found := keeper.GetBond(ctx, msg.BondToken)
	if !found {
		return types.ErrBondDoesNotExist(types.DefaultCodespace, msg.BondToken).Result()
	}

	// Confirm that function type is augmented_function and state is SETTLE
	if bond.FunctionType != types.AugmentedFunction {
		return types.ErrFunctionNotAvailableForFunctionType(types.DefaultCodespace).Result()
	} else if bond.State != types.SettleState {
		return types.ErrInvalidStateForAction(types.DefaultCodespace).Result()
	}

	// Get number of bond tokens owned by the recipient
	bondTokensOwnedAmount := keeper.BankKeeper.GetCoins(ctx, msg.Recipient).AmountOf(msg.BondToken)
	if bondTokensOwnedAmount.IsZero() {
		return types.ErrNoBondTokensOwned(types.DefaultCodespace).Result()
	}
	bondTokensOwned := sdk.NewCoin(msg.BondToken, bondTokensOwnedAmount)

	// Send coins to be burned from recipient
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx, msg.Recipient, types.BondsMintBurnAccount, sdk.NewCoins(bondTokensOwned))
	if err != nil {
		return err.Result()
	}

	// Burn bond tokens
	err = keeper.SupplyKeeper.BurnCoins(ctx, types.BondsMintBurnAccount,
		sdk.NewCoins(sdk.NewCoin(msg.BondToken, bondTokensOwnedAmount)))
	if err != nil {
		return err.Result()
	}

	// Calculate amount owned
	remainingReserve := keeper.GetReserveBalances(ctx, bond.Token)
	bondTokensShare := sdk.NewDecFromInt(bondTokensOwnedAmount).QuoInt(bond.CurrentSupply.Amount)
	reserveOwedDec := sdk.NewDecCoins(remainingReserve).MulDec(bondTokensShare)
	reserveOwed, _ := reserveOwedDec.TruncateDecimal()

	// Send coins owed to recipient
	err = keeper.BankKeeper.SendCoins(
		ctx, bond.ReserveAddress, msg.Recipient, reserveOwed)
	if err != nil {
		return err.Result()
	}

	// Update supply
	keeper.SetCurrentSupply(ctx, bond.Token, bond.CurrentSupply.Sub(bondTokensOwned))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawShare,
			sdk.NewAttribute(types.AttributeKeyBond, msg.BondToken),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Recipient.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, reserveOwed.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Recipient.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
