package rest
// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	"bytes"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cybercongress/go-cyber/x/cron/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/cron/add_job",
		addJobHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(
		"/cron/remove_job",
		removeJobHandlerFn(cliCtx)).Methods("DELETE")
	r.HandleFunc(
		"/cron/change_job_cid",
		changeJobCIDHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_label",
		changeJobLabelHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_calldata",
		changeJobCallDataHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_gasprice",
		changeJoGasPriceHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_period",
		changeJobPeriodHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_block",
		changeJobBlockHandlerFn(cliCtx)).Methods("PUT")
}

func addJobHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddJobReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgCronAddJob(
			cr, co,
			types.NewTrigger(req.Period, req.Block),
			types.NewLoad(req.CallData, req.GasPrice),
			req.Label, req.CID,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func removeJobHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveJobReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronRemoveJob(
			cr, co, req.Label,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func changeJobCIDHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobCIDReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronChangeJobCID(
			cr, co, req.Label, req.CID,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func changeJobLabelHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobLabelReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronChangeJobLabel(
			cr, co, req.Label, req.NewLabel,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func changeJobCallDataHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobCallDataReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronChangeCallData(
			cr, co, req.Label, req.CallData,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func changeJoGasPriceHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobGasPriceReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronChangeGasPrice(
			cr, co, req.Label, req.GasPrice,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func changeJobBlockHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobBlockReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronChangeJobBlock(
			cr, co, req.Label, req.Block,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func changeJobPeriodHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobPeriodReq

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCronChangeJobPeriod(
			cr, co, req.Label, req.Period,
		)

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

