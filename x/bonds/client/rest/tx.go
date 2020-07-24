package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/ixoworld/bonds/x/bonds/client"
	"github.com/ixoworld/bonds/x/bonds/internal/types"
	"net/http"
	"strings"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/bonds/create_bond",
		createBondHandler(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/bonds/edit_bond",
		editBondHandler(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/bonds/buy",
		buyHandler(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/bonds/sell",
		sellHandler(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/bonds/swap",
		swapHandler(cliCtx),
	).Methods("POST")
}

type createBondReq struct {
	BaseReq                rest.BaseReq `json:"base_req" yaml:"base_req"`
	Token                  string       `json:"token" yaml:"token"`
	Name                   string       `json:"name" yaml:"name"`
	Description            string       `json:"description" yaml:"description"`
	FunctionType           string       `json:"function_type" yaml:"function_type"`
	FunctionParameters     string       `json:"function_parameters" yaml:"function_parameters"`
	ReserveTokens          string       `json:"reserve_tokens" yaml:"reserve_tokens"`
	TxFeePercentage        string       `json:"tx_fee_percentage" yaml:"tx_fee_percentage"`
	ExitFeePercentage      string       `json:"exit_fee_percentage" yaml:"exit_fee_percentage"`
	FeeAddress             string       `json:"fee_address" yaml:"fee_address"`
	MaxSupply              string       `json:"max_supply" yaml:"max_supply"`
	OrderQuantityLimits    string       `json:"order_quantity_limits" yaml:"order_quantity_limits"`
	SanityRate             string       `json:"sanity_rate" yaml:"sanity_rate"`
	SanityMarginPercentage string       `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
	AllowSells             string       `json:"allow_sells" yaml:"allow_sells"`
	Signers                string       `json:"signers" yaml:"signers"`
	BatchBlocks            string       `json:"batch_blocks" yaml:"batch_blocks"`
}

func createBondHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createBondReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse function parameters
		functionParams, err := client.ParseFunctionParams(req.FunctionParameters)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse reserve tokens
		reserveTokens := strings.Split(req.ReserveTokens, ",")

		// Parse tx fee percentage
		txFeePercentageDec, err := sdk.NewDecFromStr(req.TxFeePercentage)
		if err != nil {
			err = types.ErrArgumentMissingOrNonFloat(types.DefaultCodespace, "tx fee percentage")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse exit fee percentage
		exitFeePercentageDec, err := sdk.NewDecFromStr(req.ExitFeePercentage)
		if err != nil {
			err = types.ErrArgumentMissingOrNonFloat(types.DefaultCodespace, "exit fee percentage")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse fee address
		feeAddress, err2 := sdk.AccAddressFromBech32(req.FeeAddress)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		// Parse max supply
		maxSupply, err2 := sdk.ParseCoin(req.MaxSupply)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		// Parse order quantity limits
		orderQuantityLimits, err2 := sdk.ParseCoins(req.OrderQuantityLimits)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		// Parse sanity rate
		sanityRate, err := sdk.NewDecFromStr(req.SanityRate)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse sanity margin percentage
		sanityMarginPercentage, err := sdk.NewDecFromStr(req.SanityMarginPercentage)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse signers
		signers, err := client.ParseSigners(req.Signers)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse batch blocks
		batchBlocks, err2 := sdk.ParseUint(req.BatchBlocks)
		if err2 != nil {
			err := types.ErrArgumentMissingOrNonUInteger(types.DefaultCodespace, "max batch blocks")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateBond(req.Token, req.Name, req.Description,
			creator, req.FunctionType, functionParams, reserveTokens,
			txFeePercentageDec, exitFeePercentageDec, feeAddress, maxSupply,
			orderQuantityLimits, sanityRate, sanityMarginPercentage,
			req.AllowSells, signers, batchBlocks)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type editBondReq struct {
	BaseReq                rest.BaseReq `json:"base_req" yaml:"base_req"`
	Token                  string       `json:"token" yaml:"token"`
	Name                   string       `json:"name" yaml:"name"`
	Description            string       `json:"description" yaml:"description"`
	OrderQuantityLimits    string       `json:"order_quantity_limits" yaml:"order_quantity_limits"`
	SanityRate             string       `json:"sanity_rate" yaml:"sanity_rate"`
	SanityMarginPercentage string       `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
	Signers                string       `json:"signers" yaml:"signers"`
}

func editBondHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editBondReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		editor, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse signers
		signers, err := client.ParseSigners(req.Signers)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgEditBond(req.Token, req.Name, req.Description,
			req.OrderQuantityLimits, req.SanityRate, req.SanityMarginPercentage,
			editor, signers)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type buyReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	BondToken  string       `json:"bond_token" yaml:"bond_token"`
	BondAmount string       `json:"bond_amount" yaml:"bond_amount"`
	MaxPrices  string       `json:"max_prices" yaml:"max_prices"`
}

func buyHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		buyer, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bondCoin, err := client.ParseTwoPartCoin(req.BondAmount, req.BondToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		maxPrices, err := sdk.ParseCoins(req.MaxPrices)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgBuy(buyer, bondCoin, maxPrices)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type sellReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	BondToken  string       `json:"bond_token" yaml:"bond_token"`
	BondAmount string       `json:"bond_amount" yaml:"bond_amount"`
}

func sellHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sellReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		seller, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bondCoin, err := client.ParseTwoPartCoin(req.BondAmount, req.BondToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSell(seller, bondCoin)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type swapReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	BondToken  string       `json:"bond_token" yaml:"bond_token"`
	FromAmount string       `json:"from_amount" yaml:"from_amount"`
	FromToken  string       `json:"from_token" yaml:"from_token"`
	ToToken    string       `json:"to_token" yaml:"to_token"`
}

func swapHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req swapReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		swapper, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check that from amount and token can be parsed to a coin
		fromCoin, err := client.ParseTwoPartCoin(req.FromAmount, req.FromToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSwap(swapper, req.BondToken, fromCoin, req.ToToken)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
