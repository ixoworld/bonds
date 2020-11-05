package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
)

func (k Keeper) MustGetBatch(ctx sdk.Context, token string) types.Batch {
	store := ctx.KVStore(k.storeKey)
	if !k.BatchExists(ctx, token) {
		panic(fmt.Sprintf("batch not found for %s\n", token))
	}

	bz := store.Get(types.GetBatchKey(token))
	var batch types.Batch
	k.cdc.MustUnmarshalBinaryBare(bz, &batch)

	return batch
}

func (k Keeper) MustGetLastBatch(ctx sdk.Context, token string) types.Batch {
	store := ctx.KVStore(k.storeKey)
	if !k.LastBatchExists(ctx, token) {
		panic(fmt.Sprintf("last batch not found for %s\n", token))
	}

	bz := store.Get(types.GetLastBatchKey(token))
	var batch types.Batch
	k.cdc.MustUnmarshalBinaryBare(bz, &batch)

	return batch
}

func (k Keeper) BatchExists(ctx sdk.Context, token string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBatchKey(token))
}

func (k Keeper) LastBatchExists(ctx sdk.Context, token string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetLastBatchKey(token))
}

func (k Keeper) SetBatch(ctx sdk.Context, token string, batch types.Batch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBatchKey(token), k.cdc.MustMarshalBinaryBare(batch))
}

func (k Keeper) SetLastBatch(ctx sdk.Context, token string, batch types.Batch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastBatchKey(token), k.cdc.MustMarshalBinaryBare(batch))
}

func (k Keeper) AddBuyOrder(ctx sdk.Context, token string, bo types.BuyOrder, buyPrices, sellPrices sdk.DecCoins) {
	batch := k.MustGetBatch(ctx, token)
	batch.TotalBuyAmount = batch.TotalBuyAmount.Add(bo.Amount)
	batch.BuyPrices = buyPrices
	batch.SellPrices = sellPrices
	batch.Buys = append(batch.Buys, bo)
	k.SetBatch(ctx, token, batch)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("added buy order for %s from %s", bo.Amount.String(), bo.Address.String()))
}

func (k Keeper) AddSellOrder(ctx sdk.Context, token string, so types.SellOrder, buyPrices, sellPrices sdk.DecCoins) {
	batch := k.MustGetBatch(ctx, token)
	batch.TotalSellAmount = batch.TotalSellAmount.Add(so.Amount)
	batch.BuyPrices = buyPrices
	batch.SellPrices = sellPrices
	batch.Sells = append(batch.Sells, so)
	k.SetBatch(ctx, token, batch)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("added sell order for %s from %s", so.Amount.String(), so.Address.String()))
}

func (k Keeper) AddSwapOrder(ctx sdk.Context, token string, so types.SwapOrder) {
	batch := k.MustGetBatch(ctx, token)
	batch.Swaps = append(batch.Swaps, so)
	k.SetBatch(ctx, token, batch)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("added swap order for %s to %s from %s", so.Amount.String(), so.ToToken, so.Address.String()))
}

func (k Keeper) GetBatchBuySellPrices(ctx sdk.Context, token string, batch types.Batch) (buyPricesPT, sellPricesPT sdk.DecCoins, err error) {
	bond := k.MustGetBond(ctx, token)

	buyAmountDec := batch.TotalBuyAmount.Amount.ToDec()
	sellAmountDec := batch.TotalSellAmount.Amount.ToDec()

	reserveBalances := k.GetReserveBalances(ctx, token)
	currentPricesPT, err := bond.GetCurrentPricesPT(reserveBalances)
	if err != nil {
		return nil, nil, err
	}

	// Get (amount of) matched and (actual) curve-calculated value for the remaining amount
	// - The matched amount is the least of the buys and sells (i.e. greatest common amount)
	// - The curved values are the prices/returns for the extra unmatched buys/sells
	var matchedAmount sdk.Dec
	var curvedValues sdk.DecCoins
	if batch.EqualBuysAndSells() {
		// Since equal, both prices are current prices
		return currentPricesPT, currentPricesPT, nil
	} else if batch.MoreBuysThanSells() {
		matchedAmount = sellAmountDec // since sells < buys, greatest common amount is sells
		extraBuys := batch.TotalBuyAmount.Sub(batch.TotalSellAmount)
		curvedValues, err = bond.GetPricesToMint(extraBuys.Amount, reserveBalances) // buy prices
		if err != nil {
			return nil, nil, err
		}
	} else {
		matchedAmount = buyAmountDec // since buys < sells, greatest common amount is buys
		extraSells := batch.TotalSellAmount.Sub(batch.TotalBuyAmount)
		curvedValues = bond.GetReturnsForBurn(extraSells.Amount, reserveBalances) // sell returns
	}

	// Get (actual) matched values
	matchedValues := types.MultiplyDecCoinsByDec(currentPricesPT, matchedAmount)

	// If buys > sells, totalValues is the total buy prices
	// If sells > buys, totalValues is the total sell returns
	totalValues := matchedValues.Add(curvedValues...)

	// Calculate buy and sell prices per token
	if batch.MoreBuysThanSells() {
		buyPricesPT = types.DivideDecCoinsByDec(totalValues, buyAmountDec)
		sellPricesPT = currentPricesPT
	} else {
		buyPricesPT = currentPricesPT
		sellPricesPT = types.DivideDecCoinsByDec(totalValues, sellAmountDec)
	}
	return buyPricesPT, sellPricesPT, nil
}

func (k Keeper) GetUpdatedBatchPricesAfterBuy(ctx sdk.Context, token string, bo types.BuyOrder) (buyPrices, sellPrices sdk.DecCoins, err error) {
	bond := k.MustGetBond(ctx, token)
	batch := k.MustGetBatch(ctx, token)

	// Max supply cannot be less than supply (max supply >= supply)
	adjustedSupply := k.GetSupplyAdjustedForBuy(ctx, token)
	adjustedSupplyWithBuy := adjustedSupply.Add(bo.Amount)
	if bond.MaxSupply.IsLT(adjustedSupplyWithBuy) {
		return nil, nil, sdkerrors.Wrap(types.ErrCannotMintMoreThanMaxSupply, bond.MaxSupply.String())
	}

	// If augmented in hatch phase and adjusted supply exceeds S0, disallow buy
	// since it is not allowed for a batch to cross over to the open phase.
	//
	// S0 is rounded to ceil for the case that it has a decimal, otherwise it
	// cannot be reached without being exceeded, when using integer buy amounts
	// (e.g. if supply is 100 and S0=100.5, we cannot reach S0 by performing
	// the minimum buy of 1 token [101>100.5], so S0 is rounded to ceil; S0=101)
	if bond.FunctionType == types.AugmentedFunction &&
		bond.State == types.HatchState {
		args := bond.FunctionParameters.AsMap()
		if adjustedSupplyWithBuy.Amount.ToDec().GT(args["S0"].Ceil()) {
			return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Buy exceeds initial supply S0. Consider buying less tokens.")
		}
	}

	// Simulate buy by bumping up total buy amount
	batch.TotalBuyAmount = batch.TotalBuyAmount.Add(bo.Amount)
	buyPrices, sellPrices, err = k.GetBatchBuySellPrices(ctx, token, batch)
	if err != nil {
		return nil, nil, err
	}

	err = k.CheckIfBuyOrderFulfillableAtPrice(ctx, token, bo, buyPrices)
	if err != nil {
		return nil, nil, err
	}

	return buyPrices, sellPrices, nil
}

func (k Keeper) GetUpdatedBatchPricesAfterSell(ctx sdk.Context, token string, so types.SellOrder) (buyPrices, sellPrices sdk.DecCoins, err error) {
	batch := k.MustGetBatch(ctx, token)

	// Cannot burn more tokens than what exists
	adjustedSupply := k.GetSupplyAdjustedForSell(ctx, token)
	if adjustedSupply.IsLT(so.Amount) {
		return nil, nil, sdkerrors.Wrap(types.ErrCannotBurnMoreThanSupply, so.Amount.String())
	}

	// Simulate sell by bumping up total sell amount
	batch.TotalSellAmount = batch.TotalSellAmount.Add(so.Amount)
	buyPrices, sellPrices, err = k.GetBatchBuySellPrices(ctx, token, batch)
	if err != nil {
		return nil, nil, err
	}

	return buyPrices, sellPrices, nil
}

func (k Keeper) PerformBuyAtPrice(ctx sdk.Context, token string, bo types.BuyOrder, prices sdk.DecCoins) (err error) {
	bond := k.MustGetBond(ctx, token)
	var extraEventAttributes []sdk.Attribute

	// Mint bond tokens
	err = k.SupplyKeeper.MintCoins(ctx, types.BondsMintBurnAccount,
		sdk.Coins{bo.Amount})
	if err != nil {
		return err
	}

	// Send bond tokens bought to buyer
	err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
		types.BondsMintBurnAccount, bo.Address, sdk.Coins{bo.Amount})
	if err != nil {
		return err
	}

	reservePrices := types.MultiplyDecCoinsByInt(prices, bo.Amount.Amount)
	reservePricesRounded := types.RoundReservePrices(reservePrices)
	txFees := bond.GetTxFees(reservePrices)
	totalPrices := reservePricesRounded.Add(txFees...)

	if totalPrices.IsAnyGT(bo.MaxPrices) {
		sdkerrors.Wrapf(types.ErrMaxPriceExceeded, "Actual prices %s exceed max prices %s", totalPrices, bo.MaxPrices)
	}

	// Add new reserve to reserve (reservePricesRounded should never be zero)
	// TODO: investigate possibility of zero reservePricesRounded
	if bond.FunctionType == types.AugmentedFunction &&
		bond.State == types.HatchState {
		args := bond.FunctionParameters.AsMap()
		theta := args["theta"]

		// Get current reserve
		var currentReserve sdk.Int
		if bond.CurrentReserve.Empty() {
			currentReserve = sdk.ZeroInt()
		} else {
			// Reserve balances should all be equal given that we are always
			// applying the same additions/subtractions to all reserve balances.
			// Thus we can pick the first reserve balance as the global balance.
			currentReserve = k.GetReserveBalances(ctx, token)[0].Amount
		}

		// Calculate expected new reserve (as fraction 1-theta of new total raise)
		newSupply := bond.CurrentSupply.Add(bo.Amount).Amount
		newTotalRaise := args["p0"].Mul(newSupply.ToDec())
		newReserve := newTotalRaise.Mul(
			sdk.OneDec().Sub(theta)).Ceil().TruncateInt()

		// Calculate amount that should go into initial reserve
		toInitialReserve := newReserve.Sub(currentReserve)
		if reservePricesRounded[0].Amount.LT(toInitialReserve) {
			// Reserve supplied by buyer is insufficient
			return sdkerrors.Wrap(types.ErrInsufficientReserveToBuy, toInitialReserve.String())
		}
		coinsToInitialReserve, _ := bond.GetNewReserveDecCoins(
			toInitialReserve.ToDec()).TruncateDecimal()

		// Calculate amount that should go into funding pool
		coinsToFundingPool := reservePricesRounded.Sub(coinsToInitialReserve)

		// Send reserve tokens to initial reserve
		err = k.DepositReserveFromModule(ctx, bond.Token,
			types.BatchesIntermediaryAccount, coinsToInitialReserve)
		if err != nil {
			return err
		}

		// Send reserve tokens to funding pool
		err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
			types.BatchesIntermediaryAccount, bond.FeeAddress, coinsToFundingPool)
		if err != nil {
			return err
		}

		extraEventAttributes = append(extraEventAttributes,
			sdk.NewAttribute(types.AttributeKeyChargedPricesReserve, toInitialReserve.String()),
			sdk.NewAttribute(types.AttributeKeyChargedPricesFunding, coinsToFundingPool.String()),
		)
	} else {
		err = k.DepositReserveFromModule(
			ctx, bond.Token, types.BatchesIntermediaryAccount, reservePricesRounded)
		if err != nil {
			return err
		}
	}

	// Add charged fee to fee address
	if !txFees.IsZero() {
		err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
			types.BatchesIntermediaryAccount, bond.FeeAddress, txFees)
		if err != nil {
			return err
		}
	}

	// Add remainder to buyer address
	returnToBuyer := bo.MaxPrices.Sub(totalPrices)
	if !returnToBuyer.IsZero() {
		err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
			types.BatchesIntermediaryAccount, bo.Address, returnToBuyer)
		if err != nil {
			return err
		}
	}

	// Update supply (max supply exceeded check done during MsgBuy)
	k.SetCurrentSupply(ctx, token, bond.CurrentSupply.Add(bo.Amount))

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("performed buy order for %s from %s", bo.Amount.String(), bo.Address.String()))

	// Get new bond token balance
	bondTokenBalance := k.BankKeeper.GetCoins(ctx, bo.Address).AmountOf(bond.Token)

	event := sdk.NewEvent(
		types.EventTypeOrderFulfill,
		sdk.NewAttribute(types.AttributeKeyBond, bond.Token),
		sdk.NewAttribute(types.AttributeKeyOrderType, types.AttributeValueBuyOrder),
		sdk.NewAttribute(types.AttributeKeyAddress, bo.Address.String()),
		sdk.NewAttribute(types.AttributeKeyTokensMinted, bo.Amount.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyChargedPrices, reservePricesRounded.String()),
		sdk.NewAttribute(types.AttributeKeyChargedFees, txFees.String()),
		sdk.NewAttribute(types.AttributeKeyReturnedToAddress, returnToBuyer.String()),
		sdk.NewAttribute(types.AttributeKeyNewBondTokenBalance, bondTokenBalance.String()),
	)
	if len(extraEventAttributes) > 0 {
		event = event.AppendAttributes(extraEventAttributes...)
	}
	ctx.EventManager().EmitEvent(event)

	return nil
}

func (k Keeper) PerformSellAtPrice(ctx sdk.Context, token string, so types.SellOrder, prices sdk.DecCoins) (err error) {
	bond := k.MustGetBond(ctx, token)

	reserveReturns := types.MultiplyDecCoinsByInt(prices, so.Amount.Amount)
	reserveReturnsRounded := types.RoundReserveReturns(reserveReturns)
	txFees := bond.GetTxFees(reserveReturns)
	exitFees := bond.GetExitFees(reserveReturns)

	totalFees := types.AdjustFees(txFees.Add(exitFees...), reserveReturnsRounded) // calculate actual total fees
	totalReturns := reserveReturnsRounded.Sub(totalFees)                          // calculate actual reserveReturns

	// Send total returns to seller (totalReturns should never be zero)
	// TODO: investigate possibility of zero totalReturns
	err = k.WithdrawReserve(ctx, bond.Token, so.Address, totalReturns)
	if err != nil {
		return err
	}

	// Send total fee to fee address
	if !totalFees.IsZero() {
		err = k.WithdrawReserve(ctx, bond.Token, bond.FeeAddress, totalFees)
		if err != nil {
			return err
		}
	}

	// Update supply (burn more than supply check done during MsgSell)
	k.SetCurrentSupply(ctx, token, bond.CurrentSupply.Sub(so.Amount))

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("performed sell order for %s from %s", so.Amount.String(), so.Address.String()))

	// Get new bond token balance
	bondTokenBalance := k.BankKeeper.GetCoins(ctx, so.Address).AmountOf(bond.Token)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOrderFulfill,
		sdk.NewAttribute(types.AttributeKeyBond, bond.Token),
		sdk.NewAttribute(types.AttributeKeyOrderType, types.AttributeValueSellOrder),
		sdk.NewAttribute(types.AttributeKeyAddress, so.Address.String()),
		sdk.NewAttribute(types.AttributeKeyTokensBurned, so.Amount.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyChargedFees, txFees.String()),
		sdk.NewAttribute(types.AttributeKeyReturnedToAddress, totalReturns.String()),
		sdk.NewAttribute(types.AttributeKeyNewBondTokenBalance, bondTokenBalance.String()),
	))

	return nil
}

func (k Keeper) PerformSwap(ctx sdk.Context, token string, so types.SwapOrder) (err error, ok bool) {
	bond := k.MustGetBond(ctx, token)

	// WARNING: do not return ok=true if money has already been transferred when error occurs

	// Get return for swap
	reserveBalances := k.GetReserveBalances(ctx, token)
	reserveReturns, txFee, err := bond.GetReturnsForSwap(so.Amount, so.ToToken, reserveBalances)
	if err != nil {
		return err, true
	}
	adjustedInput := so.Amount.Sub(txFee) // same as during GetReturnsForSwap

	// Check if new rates violate sanity rate
	newReserveBalances := reserveBalances.Add(adjustedInput).Sub(reserveReturns)
	if bond.ReservesViolateSanityRate(newReserveBalances) {
		return sdkerrors.Wrap(types.ErrValuesViolateSanityRate, newReserveBalances.String()), true
	}

	// Give resultant tokens to swapper (reserveReturns should never be zero)
	err = k.WithdrawReserve(ctx, bond.Token, so.Address, reserveReturns)
	if err != nil {
		return err, false
	}

	// Add fee-reduced coins to be swapped to reserve (adjustedInput should never be zero)
	err = k.DepositReserveFromModule(
		ctx, bond.Token, types.BatchesIntermediaryAccount, sdk.Coins{adjustedInput})
	if err != nil {
		return err, false
	}

	// Add fee (taken from swapper) to fee address
	if !txFee.IsZero() {
		err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
			types.BatchesIntermediaryAccount, bond.FeeAddress, sdk.Coins{txFee})
		if err != nil {
			return err, false
		}
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("performed swap order for %s to %s from %s",
		so.Amount.String(), reserveReturns, so.Address.String()))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOrderFulfill,
		sdk.NewAttribute(types.AttributeKeyBond, bond.Token),
		sdk.NewAttribute(types.AttributeKeyOrderType, types.AttributeValueSwapOrder),
		sdk.NewAttribute(types.AttributeKeyAddress, so.Address.String()),
		sdk.NewAttribute(types.AttributeKeyTokensSwapped, adjustedInput.String()),
		sdk.NewAttribute(types.AttributeKeyChargedFees, txFee.String()),
		sdk.NewAttribute(types.AttributeKeyReturnedToAddress, reserveReturns.String()),
	))

	return nil, true
}

func (k Keeper) PerformBuyOrders(ctx sdk.Context, token string) {
	batch := k.MustGetBatch(ctx, token)

	// Perform buys or return to buyer
	for _, bo := range batch.Buys {
		if !bo.IsCancelled() {
			err := k.PerformBuyAtPrice(ctx, token, bo, batch.BuyPrices)
			if err != nil {
				// Panic here since all calculations should have been done
				// correctly to prevent any errors during the buy
				panic(err)
			}
		}
	}

	// Update batch with any new changes (shouldn't be any)
	k.SetBatch(ctx, token, batch)
}

func (k Keeper) PerformSellOrders(ctx sdk.Context, token string) {
	batch := k.MustGetBatch(ctx, token)

	// Perform sells or return to seller
	for _, so := range batch.Sells {
		if !so.IsCancelled() {
			err := k.PerformSellAtPrice(ctx, token, so, batch.SellPrices)
			if err != nil {
				// Panic here since all calculations should have been done
				// correctly to prevent any errors during the sell
				panic(err)
			}
		}
	}

	// Update batch with any new changes (shouldn't be any)
	k.SetBatch(ctx, token, batch)
}

func (k Keeper) PerformSwapOrders(ctx sdk.Context, token string) {
	logger := ctx.Logger()
	batch := k.MustGetBatch(ctx, token)

	// Perform swaps
	// TODO: implement swaps front-running prevention
	for i, so := range batch.Swaps {
		if !so.IsCancelled() {
			err, ok := k.PerformSwap(ctx, token, so)
			if err != nil {
				if ok {
					batch.Swaps[i].Cancelled = true
					batch.Swaps[i].CancelReason = err.Error()

					logger.Info(fmt.Sprintf("cancelled swap order for %s to %s from %s", so.Amount.String(), so.ToToken, so.Address.String()))
					logger.Info(fmt.Sprintf("cancellation reason: %s", err.Error())) // set to Info for visibility

					ctx.EventManager().EmitEvent(sdk.NewEvent(
						types.EventTypeOrderCancel,
						sdk.NewAttribute(types.AttributeKeyBond, token),
						sdk.NewAttribute(types.AttributeKeyOrderType, types.AttributeValueSwapOrder),
						sdk.NewAttribute(types.AttributeKeyAddress, so.Address.String()),
						sdk.NewAttribute(types.AttributeKeyCancelReason, batch.Swaps[i].CancelReason),
					))

					// Return from amount to swapper
					err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
						types.BatchesIntermediaryAccount, so.Address, sdk.Coins{so.Amount})
					if err != nil {
						panic(err)
					}
				} else {
					// Panic here since all calculations should have been done
					// correctly to prevent any errors during the swap
					panic(err)
				}
			}
		}
	}

	// Update batch with any new cancellations
	k.SetBatch(ctx, token, batch)
}

func (k Keeper) PerformOrders(ctx sdk.Context, token string) {
	k.PerformBuyOrders(ctx, token)
	k.PerformSellOrders(ctx, token)
	k.PerformSwapOrders(ctx, token)
}

func (k Keeper) CheckIfBuyOrderFulfillableAtPrice(ctx sdk.Context, token string, bo types.BuyOrder, prices sdk.DecCoins) error {
	bond := k.MustGetBond(ctx, token)

	reservePrices := types.MultiplyDecCoinsByInt(prices, bo.Amount.Amount)
	reserveRounded := types.RoundReservePrices(reservePrices)
	txFees := bond.GetTxFees(reservePrices)
	totalPrices := reserveRounded.Add(txFees...)

	// Check that max prices not exceeded
	if totalPrices.IsAnyGT(bo.MaxPrices) {
		return sdkerrors.Wrapf(types.ErrMaxPriceExceeded, "Actual prices %s exceed max prices %s", totalPrices, bo.MaxPrices)
	}

	return nil
}

func (k Keeper) CancelUnfulfillableBuys(ctx sdk.Context, token string) (cancelledOrders int) {
	logger := k.Logger(ctx)
	batch := k.MustGetBatch(ctx, token)

	// Cancel unfulfillable buys
	for i, bo := range batch.Buys {
		if !bo.IsCancelled() {
			err := k.CheckIfBuyOrderFulfillableAtPrice(ctx, token, bo, batch.BuyPrices)
			if err != nil {
				// Cancel (important to use batch.Buys[i] and not bo!)
				batch.Buys[i].Cancelled = true
				batch.Buys[i].CancelReason = err.Error()
				batch.TotalBuyAmount = batch.TotalBuyAmount.Sub(bo.Amount)
				cancelledOrders += 1

				logger.Info(fmt.Sprintf("cancelled buy order for %s from %s", bo.Amount.String(), bo.Address.String()))
				logger.Info(fmt.Sprintf("cancellation reason: %s", err.Error())) // set to Info for visibility

				ctx.EventManager().EmitEvent(sdk.NewEvent(
					types.EventTypeOrderCancel,
					sdk.NewAttribute(types.AttributeKeyBond, token),
					sdk.NewAttribute(types.AttributeKeyOrderType, types.AttributeValueBuyOrder),
					sdk.NewAttribute(types.AttributeKeyAddress, bo.Address.String()),
					sdk.NewAttribute(types.AttributeKeyCancelReason, batch.Buys[i].CancelReason),
				))

				// Return reserve to buyer
				err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
					types.BatchesIntermediaryAccount, bo.Address, bo.MaxPrices)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// Save batch and return number of cancelled orders
	k.SetBatch(ctx, token, batch)
	return cancelledOrders
}

func (k Keeper) CancelUnfulfillableOrders(ctx sdk.Context, token string) (cancelledOrders int) {
	batch := k.MustGetBatch(ctx, token)
	cancelledOrders = 0

	cancelledOrders += k.CancelUnfulfillableBuys(ctx, token)
	//cancelledOrders += k.CancelUnfulfillableSells(ctx, token) // Sells always fulfillable
	//cancelledOrders += k.CancelUnfulfillableSwaps(ctx, token) // Swaps only cancelled while they are being performed

	// Update buy and sell prices if any cancellation took place
	if cancelledOrders > 0 {
		batch = k.MustGetBatch(ctx, token) // get batch again
		buyPrices, sellPrices, err := k.GetBatchBuySellPrices(ctx, token, batch)
		if err != nil {
			panic(err)
		}
		batch.BuyPrices = buyPrices
		batch.SellPrices = sellPrices
	}

	// Save batch and return number of cancelled orders
	k.SetBatch(ctx, token, batch)
	return cancelledOrders
}
