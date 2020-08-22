package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace = ModuleName

	// General
	CodeArgumentInvalid                = 301
	CodeArgumentMissingOrIncorrectType = 302
	CodeIncorrectNumberOfValues        = 303
	CodeActionInvalid                  = 304

	// Bonds
	CodeBondDoesNotExist        = 305
	CodeBondAlreadyExists       = 306
	CodeBondDoesNotAllowSelling = 307
	CodeDidNotEditAnything      = 308
	CodeInvalidSwapper          = 309
	CodeInvalidBond             = 310
	CodeInvalidState            = 311

	// Function types and function parameters
	CodeUnrecognizedFunctionType             = 312
	CodeInvalidFunctionParameter             = 313
	CodeFunctionNotAvailableForFunctionType  = 314
	CodeFunctionRequiresNonZeroCurrentSupply = 315

	// Token/coin names
	CodeReserveTokenInvalid     = 316
	CodeMaxSupplyDenomInvalid   = 317
	CodeBondTokenInvalid        = 318
	CodeReserveDenomsMismatch   = 319
	CodeInvalidCoinDenomination = 320

	// Amounts and fees
	CodeInvalidResultantSupply     = 321
	CodeMaxPriceExceeded           = 322
	CodeSwapAmountInvalid          = 323
	CodeOrderQuantityLimitExceeded = 324
	CodeSanityRateViolated         = 325
	CodeFeeTooLarge                = 326
	CodeNoBondTokensOwned          = 327
	CodeInsufficientReserveToBuy   = 328
)

var (
	ErrArgumentMustBePositive               = sdkerrors.Register(ModuleName, CodeArgumentInvalid, "argument must be a positive value")
	ErrArgumentMustBeInteger                = sdkerrors.Register(ModuleName, CodeArgumentInvalid, "argument must be an integer value")
	ErrArgumentMustBeBetween                = sdkerrors.Register(ModuleName, CodeArgumentInvalid, "argument must be between")
	ErrArgumentCannotBeEmpty                = sdkerrors.Register(ModuleName, CodeArgumentInvalid, "argument cannot be empty")
	ErrArgumentCannotBeNegative             = sdkerrors.Register(ModuleName, CodeArgumentInvalid, "argument cannot be negative")
	ErrArgumentMissingOrNonFloat            = sdkerrors.Register(ModuleName, CodeArgumentMissingOrIncorrectType, "Argument is missing or is not a float")
	ErrBondDoesNotExist                     = sdkerrors.Register(ModuleName, CodeBondDoesNotExist, "Bond does not exist")
	ErrBondAlreadyExists                    = sdkerrors.Register(ModuleName, CodeBondAlreadyExists, "Bond already exists")
	ErrBondTokenCannotBeStakingToken        = sdkerrors.Register(ModuleName, CodeBondTokenInvalid, "Bond token cannot be staking token")
	ErrInvalidStateForAction                = sdkerrors.Register(ModuleName, CodeInvalidState, "Cannot perform that action at the current state")
	ErrReserveDenomsMismatch                = sdkerrors.Register(ModuleName, CodeReserveDenomsMismatch, "Denom do not match reserve")
	ErrOrderQuantityLimitExceeded           = sdkerrors.Register(ModuleName, CodeOrderQuantityLimitExceeded, "Order quantity limits exceeded")
	ErrValuesViolateSanityRate              = sdkerrors.Register(ModuleName, CodeSanityRateViolated, "Values violate sanity rate")
	ErrBondDoesNotAllowSelling              = sdkerrors.Register(ModuleName, CodeBondDoesNotAllowSelling, "Bond does not allow selling at the moment")
	ErrFunctionNotAvailableForFunctionType  = sdkerrors.Register(ModuleName, CodeFunctionNotAvailableForFunctionType, "Function is not available for the function type")
	ErrCannotMakeZeroOutcomePayment         = sdkerrors.Register(ModuleName, CodeActionInvalid, "Cannot make outcome payment because outcome payment is set to nil")
	ErrNoBondTokensOwned                    = sdkerrors.Register(ModuleName, CodeNoBondTokensOwned, "No bond tokens of this bond are owned")
	ErrCannotBurnMoreThanSupply             = sdkerrors.Register(ModuleName, CodeInvalidResultantSupply, "Cannot burn more tokens than the current supply")
	ErrFeesCannotBeOrExceed100Percent       = sdkerrors.Register(ModuleName, CodeFeeTooLarge, "Sum of fees is or exceeds 100 percent")
	ErrFromAndToCannotBeTheSameToken        = sdkerrors.Register(ModuleName, CodeInvalidSwapper, "From and To tokens cannot be the same token")
	ErrCannotMintMoreThanMaxSupply          = sdkerrors.Register(ModuleName, CodeInvalidResultantSupply, "Cannot mint more tokens than the max supply")
	ErrMaxPriceExceeded                     = sdkerrors.Register(ModuleName, CodeMaxPriceExceeded, "Max price exceeded")
	ErrInsufficientReserveToBuy             = sdkerrors.Register(ModuleName, CodeInsufficientReserveToBuy, "Insufficient reserve was supplied to perform buy order")
	ErrIncorrectNumberOfFunctionParameters  = sdkerrors.Register(ModuleName, CodeIncorrectNumberOfValues, "Incorrect number of function parameters")
	ErrFunctionParameterMissingOrNonFloat   = sdkerrors.Register(ModuleName, CodeArgumentMissingOrIncorrectType, "Parameter is missing or is not a float")
	ErrFunctionRequiresNonZeroCurrentSupply = sdkerrors.Register(ModuleName, CodeFunctionRequiresNonZeroCurrentSupply, "Function requires the current supply to be non zero")
	ErrTokenIsNotAValidReserveToken         = sdkerrors.Register(ModuleName, CodeReserveTokenInvalid, "Token is not a valid reserve token")
	ErrSwapAmountTooSmallToGiveAnyReturn    = sdkerrors.Register(ModuleName, CodeSwapAmountInvalid, "Swap amount too small to give any return")
	ErrSwapAmountCausesReserveDepletion     = sdkerrors.Register(ModuleName, CodeSwapAmountInvalid, "Swap amount too large and causes reserve to be depleted")
)

/*
func ErrArgumentMissingOrNonUInteger(codespace sdk.CodespaceType, arg string) sdk.Error {
	errMsg := fmt.Sprintf("%s argument is missing or is not an unsigned integer", arg)
	return sdk.NewError(codespace, CodeArgumentMissingOrIncorrectType, errMsg)
}

func ErrArgumentMissingOrNonBoolean(codespace sdk.CodespaceType, arg string) sdk.Error {
	errMsg := fmt.Sprintf("%s argument is missing or is not true or false", arg)
	return sdk.NewError(codespace, CodeArgumentMissingOrIncorrectType, errMsg)
}

func ErrIncorrectNumberOfReserveTokens(codespace sdk.CodespaceType, expected int) sdk.Error {
	errMsg := fmt.Sprintf("Incorrect number of reserve tokens; expected: %d", expected)
	return sdk.NewError(codespace, CodeIncorrectNumberOfValues, errMsg)
}

func ErrDidNotEditAnything(codespace sdk.CodespaceType) sdk.Error {
	errMsg := "Did not edit anything from the bond"
	return sdk.NewError(codespace, CodeDidNotEditAnything, errMsg)
}

func ErrDuplicateReserveToken(codespace sdk.CodespaceType) sdk.Error {
	errMsg := "Cannot have duplicate tokens in reserve tokens"
	return sdk.NewError(codespace, CodeInvalidBond, errMsg)
}

func ErrUnrecognizedFunctionType(codespace sdk.CodespaceType) sdk.Error {
	errMsg := "Unrecognized function type"
	return sdk.NewError(codespace, CodeUnrecognizedFunctionType, errMsg)
}

func ErrInvalidFunctionParameter(codespace sdk.CodespaceType, parameter string) sdk.Error {
	errMsg := fmt.Sprintf("Invalid function parameter '%s'", parameter)
	return sdk.NewError(codespace, CodeInvalidFunctionParameter, errMsg)
}

func ErrMaxSupplyDenomDoesNotMatchTokenDenom(codespace sdk.CodespaceType) sdk.Error {
	errMsg := "Max supply denom does not match token denom"
	return sdk.NewError(codespace, CodeMaxSupplyDenomInvalid, errMsg)
}

func ErrBondTokenCannotAlsoBeReserveToken(codespace sdk.CodespaceType) sdk.Error {
	errMsg := "Token cannot also be a reserve token"
	return sdk.NewError(codespace, CodeBondTokenInvalid, errMsg)
}

func ErrInvalidCoinDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	errMsg := fmt.Sprintf("Invalid coin denomination '%s'", denom)
	return sdk.NewError(codespace, CodeInvalidCoinDenomination, errMsg)
}
*/
