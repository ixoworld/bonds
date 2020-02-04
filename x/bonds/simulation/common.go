package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	defaultReserveTokens = []string{sdk.DefaultBondDenom}

	blankOrderQuantityLimits    = sdk.Coins{}
	blankSanityRate             = sdk.MustNewDecFromStr("0")
	blankSanityMarginPercentage = sdk.MustNewDecFromStr("0")

	tokenPrefix    = "token"
	totalBondCount = 0 // Updated for each bond created
	maxBondCount   = 0 // Set during genesis creation
	gas            = uint64(100000000)

	swapperBonds []string
)
