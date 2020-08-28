package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

const (
	TypeMsgCreateBond         = "create_bond"
	TypeMsgEditBond           = "edit_bond"
	TypeMsgBuy                = "buy"
	TypeMsgSell               = "sell"
	TypeMsgSwap               = "swap"
	TypeMsgMakeOutcomePayment = "make_outcome_payment"
	TypeMsgWithdrawShare      = "withdraw_share"
)

type MsgCreateBond struct {
	Token                  string           `json:"token" yaml:"token"`
	Name                   string           `json:"name" yaml:"name"`
	Description            string           `json:"description" yaml:"description"`
	FunctionType           string           `json:"function_type" yaml:"function_type"`
	FunctionParameters     FunctionParams   `json:"function_parameters" yaml:"function_parameters"`
	Creator                sdk.AccAddress   `json:"creator" yaml:"creator"`
	ReserveTokens          []string         `json:"reserve_tokens" yaml:"reserve_tokens"`
	TxFeePercentage        sdk.Dec          `json:"tx_fee_percentage" yaml:"tx_fee_percentage"`
	ExitFeePercentage      sdk.Dec          `json:"exit_fee_percentage" yaml:"exit_fee_percentage"`
	FeeAddress             sdk.AccAddress   `json:"fee_address" yaml:"fee_address"`
	MaxSupply              sdk.Coin         `json:"max_supply" yaml:"max_supply"`
	OrderQuantityLimits    sdk.Coins        `json:"order_quantity_limits" yaml:"order_quantity_limits"`
	SanityRate             sdk.Dec          `json:"sanity_rate" yaml:"sanity_rate"`
	SanityMarginPercentage sdk.Dec          `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
	AllowSells             bool             `json:"allow_sells" yaml:"allow_sells"`
	Signers                []sdk.AccAddress `json:"signers" yaml:"signers"`
	BatchBlocks            sdk.Uint         `json:"batch_blocks" yaml:"batch_blocks"`
	OutcomePayment         sdk.Coins        `json:"outcome_payment" yaml:"outcome_payment"`
}

func NewMsgCreateBond(token, name, description string, creator sdk.AccAddress,
	functionType string, functionParameters FunctionParams, reserveTokens []string,
	txFeePercentage, exitFeePercentage sdk.Dec, feeAddress sdk.AccAddress, maxSupply sdk.Coin,
	orderQuantityLimits sdk.Coins, sanityRate, sanityMarginPercentage sdk.Dec,
	allowSell bool, signers []sdk.AccAddress, batchBlocks sdk.Uint,
	outcomePayment sdk.Coins) MsgCreateBond {
	return MsgCreateBond{
		Token:                  token,
		Name:                   name,
		Description:            description,
		Creator:                creator,
		FunctionType:           functionType,
		FunctionParameters:     functionParameters,
		ReserveTokens:          reserveTokens,
		TxFeePercentage:        txFeePercentage,
		ExitFeePercentage:      exitFeePercentage,
		FeeAddress:             feeAddress,
		MaxSupply:              maxSupply,
		OrderQuantityLimits:    orderQuantityLimits,
		SanityRate:             sanityRate,
		SanityMarginPercentage: sanityMarginPercentage,
		AllowSells:             allowSell,
		Signers:                signers,
		BatchBlocks:            batchBlocks,
		OutcomePayment:         outcomePayment,
	}
}

func (msg MsgCreateBond) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.Token) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Token")
	} else if strings.TrimSpace(msg.Name) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Name")
	} else if strings.TrimSpace(msg.Description) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Description")
	} else if msg.Creator.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Creator")
	} else if len(msg.ReserveTokens) == 0 {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Reserve Token")
	} else if msg.FeeAddress.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Fee Address")
	} else if strings.TrimSpace(msg.FunctionType) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Function Type")
	}
	// Note: FunctionParameters can be empty

	// Check that bond token is a valid token name
	err := CheckCoinDenom(msg.Token)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidCoinDenomination, msg.Token)
	}

	// Validate function parameters
	if err := msg.FunctionParameters.Validate(msg.FunctionType); err != nil {
		return err
	}

	// Validate reserve tokens
	if err = CheckReserveTokenNames(msg.ReserveTokens, msg.Token); err != nil {
		return err
	} else if err = CheckNoOfReserveTokens(msg.ReserveTokens, msg.FunctionType); err != nil {
		return err
	}

	// Validate coins
	if !msg.MaxSupply.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "max supply is invalid")
	} else if !msg.OrderQuantityLimits.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "order quantity limits are invalid")
	} else if !msg.OutcomePayment.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "outcome payment is invalid")
	}

	// Check that max supply denom matches token denom
	if msg.MaxSupply.Denom != msg.Token {
		return sdkerrors.Wrap(ErrMaxSupplyDenomDoesNotMatchTokenDenom, msg.Token)
	}

	// Check that Sanity values not negative
	if msg.SanityRate.IsNegative() {
		return sdkerrors.Wrap(ErrArgumentCannotBeNegative, "SanityRate")
	} else if msg.SanityMarginPercentage.IsNegative() {
		return sdkerrors.Wrap(ErrArgumentCannotBeNegative, "SanityMarginPercentage")
	}

	// Check FeePercentages not negative and don't add up to 100
	if msg.TxFeePercentage.IsNegative() {
		return sdkerrors.Wrap(ErrArgumentCannotBeNegative, "TxFeePercentage")
	} else if msg.ExitFeePercentage.IsNegative() {
		return sdkerrors.Wrap(ErrArgumentCannotBeNegative, "ExitFeePercentage")
	} else if msg.TxFeePercentage.Add(msg.ExitFeePercentage).GTE(sdk.NewDec(100)) {
		return sdkerrors.Wrap(ErrFeesCannotBeOrExceed100Percent, msg.TxFeePercentage.Add(msg.ExitFeePercentage).String())
	}

	// Check that not zero
	if msg.BatchBlocks.IsZero() {
		return sdkerrors.Wrap(ErrArgumentMustBePositive, "BatchBlocks")
	} else if msg.MaxSupply.Amount.IsZero() {
		return sdkerrors.Wrap(ErrArgumentMustBePositive, "MaxSupply")
	}

	// Note: uniqueness of reserve tokens checked when parsing

	return nil
}

func (msg MsgCreateBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateBond) GetSigners() []sdk.AccAddress {
	return msg.Signers
}

func (msg MsgCreateBond) Route() string { return RouterKey }

func (msg MsgCreateBond) Type() string { return TypeMsgCreateBond }

type MsgEditBond struct {
	Token                  string           `json:"token" yaml:"token"`
	Name                   string           `json:"name" yaml:"name"`
	Description            string           `json:"description" yaml:"description"`
	OrderQuantityLimits    string           `json:"order_quantity_limits" yaml:"order_quantity_limits"`
	SanityRate             string           `json:"sanity_rate" yaml:"sanity_rate"`
	SanityMarginPercentage string           `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
	Editor                 sdk.AccAddress   `json:"editor" yaml:"editor"`
	Signers                []sdk.AccAddress `json:"signers" yaml:"signers"`
}

func NewMsgEditBond(token, name, description, orderQuantityLimits, sanityRate,
	sanityMarginPercentage string, editor sdk.AccAddress,
	signers []sdk.AccAddress) MsgEditBond {
	return MsgEditBond{
		Token:                  token,
		Name:                   name,
		Description:            description,
		OrderQuantityLimits:    orderQuantityLimits,
		SanityRate:             sanityRate,
		SanityMarginPercentage: sanityMarginPercentage,
		Editor:                 editor,
		Signers:                signers,
	}
}

func (msg MsgEditBond) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.Token) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Token")
	} else if strings.TrimSpace(msg.Name) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Name")
	} else if strings.TrimSpace(msg.Description) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Description")
	} else if strings.TrimSpace(msg.SanityRate) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "SanityRate")
	} else if strings.TrimSpace(msg.SanityMarginPercentage) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "SanityMarginPercentage")
	} else if msg.Editor.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Editor")
	}
	// Note: order quantity limits can be blank

	// Check that at least one editable was edited. Fields that will not
	// be edited should be "DoNotModifyField", and not an empty string
	inputList := []string{
		msg.Name, msg.Description, msg.OrderQuantityLimits,
		msg.SanityRate, msg.SanityMarginPercentage,
	}
	atLeaseOneEdit := false
	for _, e := range inputList {
		if e != DoNotModifyField {
			atLeaseOneEdit = true
			break
		}
	}
	if !atLeaseOneEdit {
		return ErrDidNotEditAnything
	}

	return nil
}

func (msg MsgEditBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgEditBond) GetSigners() []sdk.AccAddress {
	return msg.Signers
}

func (msg MsgEditBond) Route() string { return RouterKey }

func (msg MsgEditBond) Type() string { return TypeMsgEditBond }

type MsgBuy struct {
	Buyer     sdk.AccAddress `json:"buyer" yaml:"buyer"`
	Amount    sdk.Coin       `json:"amount" yaml:"amount"`
	MaxPrices sdk.Coins      `json:"max_prices" yaml:"max_prices"`
}

func NewMsgBuy(buyer sdk.AccAddress, amount sdk.Coin, maxPrices sdk.Coins) MsgBuy {
	return MsgBuy{
		Buyer:     buyer,
		Amount:    amount,
		MaxPrices: maxPrices,
	}
}

func (msg MsgBuy) ValidateBasic() error {
	// Check if empty
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Buyer")
	}

	// Check that amount valid and non zero
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is invalid")
	} else if msg.Amount.Amount.IsZero() {
		return sdkerrors.Wrap(ErrArgumentMustBePositive, "Amount")
	}

	// Check that maxPrices valid
	if !msg.MaxPrices.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "maxprices is invalid")
	}

	return nil
}

func (msg MsgBuy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBuy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

func (msg MsgBuy) Route() string { return RouterKey }

func (msg MsgBuy) Type() string { return TypeMsgBuy }

type MsgSell struct {
	Seller sdk.AccAddress `json:"seller" yaml:"seller"`
	Amount sdk.Coin       `json:"amount" yaml:"amount"`
}

func NewMsgSell(seller sdk.AccAddress, amount sdk.Coin) MsgSell {
	return MsgSell{
		Seller: seller,
		Amount: amount,
	}
}

func (msg MsgSell) ValidateBasic() error {
	// Check if empty
	if msg.Seller.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Seller")
	}

	// Check that amount valid and non zero
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is invalid")
	} else if msg.Amount.Amount.IsZero() {
		return sdkerrors.Wrap(ErrArgumentMustBePositive, "Amount")
	}

	return nil
}

func (msg MsgSell) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSell) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Seller}
}

func (msg MsgSell) Route() string { return RouterKey }

func (msg MsgSell) Type() string { return TypeMsgSell }

type MsgSwap struct {
	Swapper   sdk.AccAddress `json:"swapper" yaml:"swapper"`
	BondToken string         `json:"bond_token" yaml:"bond_token"`
	From      sdk.Coin       `json:"from" yaml:"from"`
	ToToken   string         `json:"to_token" yaml:"to_token"`
}

func NewMsgSwap(swapper sdk.AccAddress, bondToken string, from sdk.Coin, toToken string) MsgSwap {
	return MsgSwap{
		Swapper:   swapper,
		BondToken: bondToken,
		From:      from,
		ToToken:   toToken,
	}
}

func (msg MsgSwap) ValidateBasic() error {
	// Check if empty
	if msg.Swapper.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Swapper")
	} else if strings.TrimSpace(msg.BondToken) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "BondToken")
	} else if strings.TrimSpace(msg.ToToken) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "ToToken")
	}

	// Validate from amount
	if !msg.From.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "from amount is invalid")
	}

	// Validate to token
	err := CheckCoinDenom(msg.ToToken)
	if err != nil {
		return err
	}

	// Check if from and to the same token
	if msg.From.Denom == msg.ToToken {
		return sdkerrors.Wrap(ErrFromAndToCannotBeTheSameToken, msg.From.Denom)
	}

	// Check that non zero
	if msg.From.Amount.IsZero() {
		return sdkerrors.Wrap(ErrArgumentMustBePositive, "FromAmount")
	}

	// Note: From denom and amount must be valid since sdk.Coin
	return nil
}

func (msg MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Swapper}
}

func (msg MsgSwap) Route() string { return RouterKey }

func (msg MsgSwap) Type() string { return TypeMsgSwap }

type MsgMakeOutcomePayment struct {
	Sender    sdk.AccAddress `json:"sender" yaml:"sender"`
	BondToken string         `json:"bond_token" yaml:"bond_token"`
}

func NewMsgMakeOutcomePayment(sender sdk.AccAddress, bondToken string) MsgMakeOutcomePayment {
	return MsgMakeOutcomePayment{
		Sender:    sender,
		BondToken: bondToken,
	}
}

func (msg MsgMakeOutcomePayment) ValidateBasic() error {
	// Check if empty
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Sender")
	} else if strings.TrimSpace(msg.BondToken) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "BondToken")
	}

	// Validate bond token
	err := CheckCoinDenom(msg.BondToken)
	if err != nil {
		return err
	}

	return nil
}

func (msg MsgMakeOutcomePayment) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMakeOutcomePayment) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgMakeOutcomePayment) Route() string { return RouterKey }

func (msg MsgMakeOutcomePayment) Type() string { return TypeMsgMakeOutcomePayment }

type MsgWithdrawShare struct {
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	BondToken string         `json:"bond_token" yaml:"bond_token"`
}

func NewMsgWithdrawShare(recipient sdk.AccAddress, bondToken string) MsgWithdrawShare {
	return MsgWithdrawShare{
		Recipient: recipient,
		BondToken: bondToken,
	}
}

func (msg MsgWithdrawShare) ValidateBasic() error {
	// Check if empty
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "Recipient")
	} else if strings.TrimSpace(msg.BondToken) == "" {
		return sdkerrors.Wrap(ErrArgumentCannotBeEmpty, "BondToken")
	}

	// Validate bond token
	err := CheckCoinDenom(msg.BondToken)
	if err != nil {
		return err
	}

	return nil
}

func (msg MsgWithdrawShare) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawShare) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Recipient}
}

func (msg MsgWithdrawShare) Route() string { return RouterKey }

func (msg MsgWithdrawShare) Type() string { return TypeMsgWithdrawShare }
