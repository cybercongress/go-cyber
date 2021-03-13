package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cybercongress/go-cyber/x/energy/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/energy/parameters",
		queryParamsHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/source_routes", Source),
		querySourceRoutesHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/destination_routes", Destination),
		queryDestinationRoutesHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/source_routed_energy", Source),
		querySourceRoutedEnergyHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/destination_routed_energy", Destination),
		queryDestinationRoutedEnergyHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/route/{%s}/{%s}", Source, Destination),
		queryRouteHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc(
		"/energy/routes",
		queryRoutesHandlerFn(cliCtx)).Methods("GET")

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

func querySourceRoutesHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		src  := vars[Source]

		addr, err := sdk.AccAddressFromBech32(src)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QuerySourceRequest{addr.String()}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySourceRoutes)
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

func queryDestinationRoutesHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		dst  := vars[Destination]

		addr, err := sdk.AccAddressFromBech32(dst)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QueryDestinationRequest{addr.String()}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDestinationRoutes)
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

func querySourceRoutedEnergyHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		src  := vars[Source]

		addr, err := sdk.AccAddressFromBech32(src)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QuerySourceRequest{addr.String()}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySourceRoutedEnergy)
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

func queryDestinationRoutedEnergyHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		dst  := vars[Destination]

		addr, err := sdk.AccAddressFromBech32(dst)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QueryDestinationRequest{addr.String()}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDestinationRoutedEnergy)
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

func queryRouteHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		src  := vars[Source]
		dst  := vars[Destination]

		addrS, err := sdk.AccAddressFromBech32(src)
		addrD, err := sdk.AccAddressFromBech32(dst)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QueryRouteRequest{addrS.String(), addrD.String()}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRoute)
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

func queryRoutesHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRoutes)
		res, height, err := cliCtx.QueryWithData(route, nil)
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