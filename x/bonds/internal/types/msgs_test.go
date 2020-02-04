package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupMsgsTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func TestValidateBasicMsgCreateTokenArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.Token = ""
	err := message.ValidateBasic()
	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateNameArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.Name = ""
	err := message.ValidateBasic()
	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateDescriptionArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.Description = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateCreatorMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.Creator = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateReserveTokenArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.ReserveTokens = nil

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgReserveAddressArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.ReserveAddress = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgFeeAddressArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.FeeAddress = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgFunctionTypeArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.FunctionType = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateAllowSellsArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.AllowSells = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateAllowSellsIsNotTrueOrFalseGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.AllowSells = "neitherTrueNorFalse"

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentMissingOrIncorrectType, err.Code())
}

func TestValidateBasicMsgCreateTxFeeIsNegativeGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.TxFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateTxFeeIsZeroGivesNoError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.TxFeePercentage = sdk.ZeroDec()

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgCreateExitFeeIsNegativeGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.ExitFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgCreateExitFeeIsZeroGivesNoError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.ExitFeePercentage = sdk.ZeroDec()

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgCreateFeeAddressEqualsReserveAddressGivesError(t *testing.T) {
	message := NewValidMsgCreateBond()
	message.FeeAddress = initReserveAddress

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeInvalidBond, err.Code())
}

func TestValidateBasicMsgCreateBondCorrectlyGivesNoError(t *testing.T) {
	message := NewValidMsgCreateBond()

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgEditBondTokenArgumentMissingGivesError(t *testing.T) {
	message := NewEmptyStringsMsgEditBond()
	message.Token = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondNameArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgEditBond()
	message.Name = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondDescriptionArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgEditBond()
	message.Description = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondOrderQuantityLimitsArgumentMissingGivesNoError(t *testing.T) {
	message := NewValidMsgEditBond()
	message.OrderQuantityLimits = ""

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgEditBondSanityRateArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgEditBond()
	message.SanityRate = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondSanityMarginPercentageArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgEditBond()
	message.SanityMarginPercentage = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgEditBondNoEditsGivesError(t *testing.T) {
	message := NewMsgEditBond(DoNotModifyField, DoNotModifyField,
		DoNotModifyField, DoNotModifyField, DoNotModifyField,
		DoNotModifyField, initCreator, initSigners)

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeDidNotEditAnything, err.Code())
}

func TestValidateBasicMsgEditBondCorrectlyGivesNoError(t *testing.T) {
	message := NewValidMsgEditBond()

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgBuyBondBuyerArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgBuy()
	message.Buyer = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgBuyBondCorrectlyGivesNoError(t *testing.T) {
	message := NewValidMsgBuy()

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgSellBondBuyerArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgSell()
	message.Seller = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgSellBondCorrectlyGivesNoError(t *testing.T) {
	message := NewValidMsgSell()

	err := message.ValidateBasic()

	require.Nil(t, err)
}

func TestValidateBasicMsgSwapSwapperArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgSwap()
	message.Swapper = sdk.AccAddress{}

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgSwapBondTokenArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgSwap()
	message.BondToken = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgSwapToTokenArgumentMissingGivesError(t *testing.T) {
	message := NewValidMsgSwap()
	message.ToToken = ""

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeArgumentInvalid, err.Code())
}

func TestValidateBasicMsgSwapFromAndToSameTokenGivesError(t *testing.T) {
	message := NewValidMsgSwap()
	message.From = sdk.NewInt64Coin(reserveToken, 10)
	message.ToToken = message.From.Denom

	err := message.ValidateBasic()

	require.NotNil(t, err)
	require.Equal(t, CodeInvalidSwapper, err.Code())
}

func TestValidateBasicMsgSwapCorrectlyGivesNoError(t *testing.T) {
	message := NewValidMsgSwap()

	err := message.ValidateBasic()

	require.Nil(t, err)
}
