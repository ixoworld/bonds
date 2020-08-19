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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateNameArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Name = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateDescriptionArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Description = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateCreatorMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Creator = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateReserveTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ReserveTokens = nil

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgFeeAddressArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FeeAddress = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgFunctionTypeArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionType = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgCreateBond: Bond token denomination

func TestValidateBasicMsgCreateInvalidTokenArgumentGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Token = "123abc" // starts with number
	err := message.ValidateBasic()
	require.NotNil(t, err)
	require.Equal(t, CodeInvalidCoinDenomination, err.Code())

	message.Token = "a" // too short
	err = message.ValidateBasic()
	require.NotNil(t, err)
	require.Equal(t, CodeInvalidCoinDenomination, err.Code())
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
	require.Equal(t, CodeIncorrectNumberOfValues, err.Code())
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
	require.Equal(t, CodeArgumentMissingOrIncorrectType, err.Code())
}

func TestValidateBasicMsgCreateNegativeFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionParameters[0].Value = sdk.NewDec(-1)

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgFunctionTypeArgumentInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.FunctionType = "invalid_function_type"

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeUnrecognizedFunctionType, err.Code())
}

// MsgCreateBond: Reserve tokens

func TestValidateBasicMsgCreateReserveTokenArgumentInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ReserveTokens[0] = "123abc" // invalid denomination

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeInvalidCoinDenomination, err.Code())
}

func TestValidateBasicMsgCreateNoReserveTokensInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.ReserveTokens = nil

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateReserveTokensWrongAmountInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateSwapperBond()
	message.ReserveTokens = append(message.ReserveTokens, "extra")

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeIncorrectNumberOfValues, err.Code())
}

// MsgCreateBond: Max supply validity

func TestValidateBasicMsgCreateInvalidMaxSupplyGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.MaxSupply.Amount = message.MaxSupply.Amount.Neg() // negate

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

// MsgCreateBond: Order quantity limits validity

func TestValidateBasicMsgCreateInvalidOrderQuantityLimitGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.OrderQuantityLimits = sdk.NewCoins(sdk.NewCoin("abc", sdk.OneInt()))
	message.OrderQuantityLimits[0].Amount = message.OrderQuantityLimits[0].Amount.Neg()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

// MsgCreateBond: Max supply denom matches bond token denom

func TestValidateBasicMsgCreateMaxSupplyDenomTokenDenomMismatchGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.Token = message.MaxSupply.Denom + "a" // to ensure different

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeMaxSupplyDenomInvalid, err.Code())
}

// MsgCreateBond: Sanity values must be positive

func TestValidateBasicMsgCreateNegativeSanityRateGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.SanityRate = sdk.OneDec().Neg()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateNegativeSanityPercentageGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.SanityMarginPercentage = sdk.OneDec().Neg()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgCreateBond: Fee percentages must be positive and not add up to 100

func TestValidateBasicMsgCreateTxFeeIsNegativeGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.TxFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
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
	require.Equal(t, CodeArgumentInvalid, err.Code())
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
	require.Equal(t, CodeFeeTooLarge, err.Code())

	message.TxFeePercentage = sdk.NewDec(50)
	message.ExitFeePercentage = sdk.NewDec(50)
	err = message.ValidateBasic()
	require.NotNil(t, err)
	require.Equal(t, CodeFeeTooLarge, err.Code())

	message.TxFeePercentage = sdk.ZeroDec()
	message.ExitFeePercentage = sdk.NewDec(100)
	err = message.ValidateBasic()
	require.NotNil(t, err)
	require.Equal(t, CodeFeeTooLarge, err.Code())

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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateZeroMaxSupplyGivesError(t *testing.T) {
	message := newValidMsgCreateBond()
	message.MaxSupply = sdk.NewCoin(token, sdk.ZeroInt())

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondNameArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.Name = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondDescriptionArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.Description = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondSanityMarginPercentageArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.SanityMarginPercentage = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondEditorArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditBond()
	message.Editor = nil

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgEditBond: no edits

func TestValidateBasicMsgEditBondNoEditsGivesError(t *testing.T) {
	message := NewMsgEditBond(DoNotModifyField, DoNotModifyField,
		DoNotModifyField, DoNotModifyField, DoNotModifyField,
		DoNotModifyField, initCreator, initSigners)

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeDidNotEditAnything, err.Code())
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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgBuy: invalid arguments

func TestValidateBasicMsgBuyInvalidAmountGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Amount.Amount = message.Amount.Amount.Neg()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

func TestValidateBasicMsgBuyZeroAmountGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Amount.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgBuyMaxPricesInvalidGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.MaxPrices[0].Amount = sdk.ZeroInt()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, sdk.CodeInvalidCoins, err.Code())
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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgSell: invalid arguments

func TestValidateBasicMsgSellInvalidAmountGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Amount.Amount = message.Amount.Amount.Neg()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

func TestValidateBasicMsgSellZeroAmountGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Amount.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
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
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgSwapBondTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.BondToken = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgSwapToTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.ToToken = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgSwap: invalid arguments

func TestValidateBasicMsgSwapInvalidFromAmountGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From.Amount = message.From.Amount.Neg()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

func TestValidateBasicMsgSwapInvalidToTokenGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.ToToken = "123abc"

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeInvalidCoinDenomination, err.Code())
}

func TestValidateBasicMsgSwapZeroFromAmountGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

// MsgSwap: fromToken==toToken

func TestValidateBasicMsgSwapFromAndToSameTokenGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From = sdk.NewInt64Coin(reserveToken, 10)
	message.ToToken = message.From.Denom

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeInvalidSwapper, err.Code())
}

// MsgSwap: correct swap

func TestValidateBasicMsgSwapCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgSwap()

	err := message.ValidateBasic()

	require.Nil(t, err)
}
