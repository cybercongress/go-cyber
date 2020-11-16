package rest
// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	// "bytes"
	// "net/http"

	"bytes"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/energy/create_route",
		postCreateEnergyRouteHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/energy/edit_route",
		postEditEnergyRouteHandlerFn(cliCtx),
	).Methods("PUT")
	r.HandleFunc(
		"/energy/delete_route",
		postDeleteEnergyRouteHandlerFn(cliCtx),
	).Methods("DELETE")
	r.HandleFunc(
		"/energy/edit_alias",
		postEditEnergyRouteAliasHandlerFn(cliCtx),
	).Methods("PUT")
}

type CreateEnergyRouteReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	Source 		 string 		`json:"source" yaml:"source"`
	Destination  string 		`json:"destination" yaml:"destination"`
	Alias 		 string 		`json:"alias" yaml:"alias"`
}

type EditEnergyRouteReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	Source 		 string 		`json:"source" yaml:"source"`
	Destination  string 		`json:"destination" yaml:"destination"`
	Value 		 string 		`json:"value" yaml:"value"`
}

type DeleteEnergyRouteReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	Source 		 string 		`json:"source" yaml:"source"`
	Destination  string 		`json:"destination" yaml:"destination"`
}

type EditEnergyRouteAliasReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	Source 		 string		    `json:"source" yaml:"source"`
	Destination  string		    `json:"destination" yaml:"destination"`
	Alias 		 string 		`json:"alias" yaml:"alias"`
}

func postCreateEnergyRouteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateEnergyRouteReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		src, err := sdk.AccAddressFromBech32(req.Source)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		dst, err := sdk.AccAddressFromBech32(req.Destination)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateEnergyRoute(src, dst, req.Alias)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, src) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postEditEnergyRouteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EditEnergyRouteReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		src, err := sdk.AccAddressFromBech32(req.Source)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		dst, err := sdk.AccAddressFromBech32(req.Destination)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		value, err := sdk.ParseCoin(req.Value)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgEditEnergyRoute(src, dst, value)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, src) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postDeleteEnergyRouteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeleteEnergyRouteReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		src, err := sdk.AccAddressFromBech32(req.Source)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		dst, err := sdk.AccAddressFromBech32(req.Destination)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteEnergyRoute(src, dst)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, src) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postEditEnergyRouteAliasHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EditEnergyRouteAliasReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		src, err := sdk.AccAddressFromBech32(req.Source)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		dst, err := sdk.AccAddressFromBech32(req.Destination)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgEditEnergyRouteAlias(src, dst, req.Alias)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, src) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

