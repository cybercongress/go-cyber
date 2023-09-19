package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	ctypes "github.com/cybercongress/go-cyber/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cybercongress/go-cyber/x/resources/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/resources/parameters",
		queryParamsHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		"/resources/investmint",
		queryInvestmintHandlerFn(cliCtx)).Methods("GET")
}

func queryParamsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)

		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryInvestmintHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var amount sdk.Coin
		var resource string
		var length uint64
		var err error

		if v := r.URL.Query().Get("amount"); len(v) != 0 {
			amount, err = sdk.ParseCoinNormalized(v)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			if amount.Denom != ctypes.SCYB {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
		}
		if v := r.URL.Query().Get("resource"); len(v) != 0 {
			if v != ctypes.VOLT && v != ctypes.AMPERE {
				err = fmt.Errorf("resource %s not a valid resource, please input a valid resource", v)
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			resource = v
		}
		if v := r.URL.Query().Get("length"); len(v) != 0 {
			length, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
		}

		params := types.NewQueryInvestmintParams(amount, resource, length)

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryInvestmint)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
