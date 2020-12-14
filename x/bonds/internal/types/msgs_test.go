package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

// MsgCreateBond: Missing arguments

func TestValidateBasicMsgCreateTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Token = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNameArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Name = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateDescriptionArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Description = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateCreatorMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Creator = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateReserveTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ReserveTokens = nil

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgFeeAddressArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FeeAddress = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgFunctionTypeArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionType = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Bond token denomination

func TestValidateBasicMsgCreateInvalidTokenArgumentGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Token = "123abc" // starts with number
	err := message.ValidateBasic()
	require.NotNil(t, err)

	message.Token = "a" // too short
	err = message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Function parameters and function type

func TestValidateBasicMsgCreateMissingFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionParameters = []FunctionParam{
		message.FunctionParameters[0],
		message.FunctionParameters[1],
	}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateTypoFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionParameters = []FunctionParam{
		NewFunctionParam("invalidParam", message.FunctionParameters[0].Value),
		message.FunctionParameters[1],
		message.FunctionParameters[2],
	}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNegativeFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionParameters[0].Value = sdk.NewDec(-1)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgFunctionTypeArgumentInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionType = "invalid_function_type"

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Reserve tokens

func TestValidateBasicMsgCreateReserveTokenArgumentInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ReserveTokens[0] = "123abc" // invalid denomination

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNoReserveTokensInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ReserveTokens = nil

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateReserveTokensWrongAmountInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateSwapperBond()
	message.ReserveTokens = append(message.ReserveTokens, "extra")

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Max supply validity

func TestValidateBasicMsgCreateInvalidMaxSupplyGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.MaxSupply.Amount = message.MaxSupply.Amount.Neg() // negate

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Order quantity limits validity

func TestValidateBasicMsgCreateInvalidOrderQuantityLimitGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.OrderQuantityLimits = sdk.NewCoins(sdk.NewCoin("abc", sdk.OneInt()))
	message.OrderQuantityLimits[0].Amount = message.OrderQuantityLimits[0].Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Max supply denom matches bond token denom

func TestValidateBasicMsgCreateMaxSupplyDenomTokenDenomMismatchGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Token = message.MaxSupply.Denom + "a" // to ensure different

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Sanity values must be positive

func TestValidateBasicMsgCreateNegativeSanityRateGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.SanityRate = sdk.OneDec().Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNegativeSanityPercentageGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.SanityMarginPercentage = sdk.OneDec().Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Fee percentages must be positive and not add up to 100

func TestValidateBasicMsgCreateTxFeeIsNegativeGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.TxFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateTxFeeIsZeroGivesNoError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.TxFeePercentage = sdk.ZeroDec()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

func TestValidateBasicMsgCreateExitFeeIsNegativeGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ExitFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateExitFeeIsZeroGivesNoError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ExitFeePercentage = sdk.ZeroDec()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

func TestValidateBasicMsgCreate100PercentFeeGivesError(t *testing.T) {
	message := newValidMsgCreateBond()

	message.TxFeePercentage = sdk.NewDec(100)
	message.ExitFeePercentage = sdk.ZeroDec()
	err := message.ValidateBasic()
	require.NotNil(t, err)

	message.TxFeePercentage = sdk.NewDec(50)
	message.ExitFeePercentage = sdk.NewDec(50)
	err = message.ValidateBasic()
	require.NotNil(t, err)

	message.TxFeePercentage = sdk.ZeroDec()
	message.ExitFeePercentage = sdk.NewDec(100)
	err = message.ValidateBasic()
	require.NotNil(t, err)

	message.TxFeePercentage = sdk.MustNewDecFromStr("49.999999")
	message.ExitFeePercentage = sdk.NewDec(50)
	require.Nil(t, message.ValidateBasic())

	message.TxFeePercentage = sdk.NewDec(50)
	message.ExitFeePercentage = sdk.MustNewDecFromStr("49.999999")
	require.Nil(t, message.ValidateBasic())
}

// MsgCreateBond: Batch blocks and max supply cannot be zero

func TestValidateBasicMsgCreateZeroBatchBlocksGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.BatchBlocks = sdk.ZeroUint()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateZeroMaxSupplyGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.MaxSupply = sdk.NewCoin(token, sdk.ZeroInt())

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateBond: Valid bond creation

func TestValidateBasicMsgCreateBondCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgCreateBond()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgEditBond: missing arguments

func TestValidateBasicMsgEditBondTokenArgumentMissingGivesError(t *testing.T) {
	message := newEmptyStringsMsgEditBond()
	message.Token = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditBondNameArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.Name = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditBondDescriptionArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.Description = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditBondOrderQuantityLimitsArgumentMissingGivesNoError(t *testing.T) {
	message := newValidMsgEditBond()
	message.OrderQuantityLimits = ""

	err := message.ValidateBasic()
	require.Nil(t, err)
}

func TestValidateBasicMsgEditBondSanityRateArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.SanityRate = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditBondSanityMarginPercentageArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.SanityMarginPercentage = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditBondEditorArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.Editor = nil

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgEditBond: no edits

func TestValidateBasicMsgEditBondNoEditsGivesError(t *testing.T) {
	message := NewMsgEditBond(DoNotModifyField, DoNotModifyField,
		DoNotModifyField, DoNotModifyField, DoNotModifyField,
		DoNotModifyField, initCreator, initSigners)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgEditBond: correct edit

func TestValidateBasicMsgEditBondCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgEditBond()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgBuy: missing arguments

func TestValidateBasicMsgBuyBuyerArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Buyer = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgBuy: invalid arguments

func TestValidateBasicMsgBuyInvalidAmountGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Amount.Amount = message.Amount.Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgBuyZeroAmountGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Amount.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgBuyMaxPricesInvalidGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.MaxPrices[0].Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgBuy: correct buy

func TestValidateBasicMsgBuyCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgBuy()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgSell: missing arguments

func TestValidateBasicMsgSellSellerArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Seller = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSell: invalid arguments

func TestValidateBasicMsgSellInvalidAmountGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Amount.Amount = message.Amount.Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSellZeroAmountGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Amount.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSell: correct sell

func TestValidateBasicMsgSellCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgSell()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgSwap: missing arguments

func TestValidateBasicMsgSwapSwapperArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.Swapper = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapBondTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.BondToken = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapToTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.ToToken = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSwap: invalid arguments

func TestValidateBasicMsgSwapInvalidFromAmountGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From.Amount = message.From.Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapInvalidToTokenGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.ToToken = "123abc"

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapZeroFromAmountGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSwap: fromToken==toToken

func TestValidateBasicMsgSwapFromAndToSameTokenGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From = sdk.NewInt64Coin(reserveToken, 10)
	message.ToToken = message.From.Denom

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSwap: correct swap

func TestValidateBasicMsgSwapCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgSwap()

	err := message.ValidateBasic()
	require.Nil(t, err)
}
