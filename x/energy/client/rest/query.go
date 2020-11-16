package rest

import (
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/energy/params",
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/source_routes", Source),
		querySourceRoutesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/destination_routes", Destination),
		queryDestinationRoutesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/source_routed_energy", Source),
		querySourceRoutedEnergyHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/{%s}/destination_routed_energy", Destination),
		queryDestinationRoutedEnergyHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/energy/route/{%s}/{%s}", Source, Destination),
		queryRouteHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/energy/routes",
		queryRoutesHandlerFn(cliCtx),
	).Methods("GET")

}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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

func querySourceRoutesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		src  := vars[Source]

		addr, err := sdk.AccAddressFromBech32(src)
		params := types.NewQuerySourceParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySourceRoutes), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDestinationRoutesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		dst  := vars[Destination]

		addr, err := sdk.AccAddressFromBech32(dst)
		params := types.NewQueryDestinationParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDestinationRoutes)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySourceRoutedEnergyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		src  := vars[Source]

		addr, err := sdk.AccAddressFromBech32(src)
		params := types.NewQuerySourceParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySourceRoutedEnergy)

		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDestinationRoutedEnergyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		dst  := vars[Destination]

		addr, err := sdk.AccAddressFromBech32(dst)
		params := types.NewQueryDestinationParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDestinationRoutedEnergy)

		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRouteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		src  := vars[Source]
		dst  := vars[Destination]

		addrS, err := sdk.AccAddressFromBech32(src)
		addrD, err := sdk.AccAddressFromBech32(dst)
		params := types.NewQueryRouteParams(addrS, addrD)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRoute)

		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRoutesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRoutes)

		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}