package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-cid"

	//"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cybercongress/go-cyber/x/rank/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/rank/parameters",
		queryParamsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/rank",
		queryRankHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/search",
		querySearchHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/backlinks",
		queryBacklinksHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/top",
		queryTopHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/is_link_exist",
		queryIsLinkExistHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/is_any_link_exist",
		queryIsAnyLinkExistHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/entropy",
		queryEntropyHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/luminosity",
		queryLuminosityHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		"/rank/karma",
		queryKarmaHandlerFn(cliCtx)).Methods("GET")
}

func queryParamsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)

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

func queryRankHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particle string

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			particle = v
		}
		if _, err := cid.Decode(particle); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		params := types.NewQueryRankParams(particle)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRank)
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

func querySearchHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particle string
		var page, limit uint32

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			particle = v
		}
		if _, err := cid.Decode(particle); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		if v := r.URL.Query().Get("page"); len(v) != 0 {
			p, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			page = uint32(p)
		} else {
			page = uint32(0)
		}
		if v := r.URL.Query().Get("limit"); len(v) != 0 {
			l, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			limit = uint32(l)
		} else {
			limit = uint32(10)
		}

		params := types.NewQuerySearchParams(particle, page, limit)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySearch)
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

func queryBacklinksHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particle string
		var page, limit uint32

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			particle = v
		}
		if _, err := cid.Decode(particle); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		if v := r.URL.Query().Get("page"); len(v) != 0 {
			p, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			page = uint32(p)
		} else {
			page = uint32(0)
		}
		if v := r.URL.Query().Get("limit"); len(v) != 0 {
			l, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			limit = uint32(l)
		} else {
			limit = uint32(10)
		}

		params := types.NewQuerySearchParams(particle, page, limit)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBacklinks)
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

func queryTopHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page, limit uint32

		if v := r.URL.Query().Get("page"); len(v) != 0 {
			p, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			page = uint32(p)
		} else {
			page = uint32(0)
		}
		if v := r.URL.Query().Get("limit"); len(v) != 0 {
			l, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			limit = uint32(l)
		} else {
			limit = uint32(10)
		}

		params := types.NewQueryTopParams(page, limit)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTop)
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

func queryIsLinkExistHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particleFrom, particleTo, addr string

		if f := r.URL.Query().Get("from"); len(f) != 0 {
			particleFrom = f
		}
		if _, err := cid.Decode(particleFrom); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if t := r.URL.Query().Get("to"); len(t) != 0 {
			particleTo = t
		}
		if _, err := cid.Decode(particleTo); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if v := r.URL.Query().Get("address"); len(v) != 0 {
			addr = v
		}
		address, err := sdk.AccAddressFromBech32(addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.NewQueryIsLinkExistParams(particleFrom, particleTo, address)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryIsLinkExist)
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

func queryIsAnyLinkExistHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particleFrom, particleTo string

		if f := r.URL.Query().Get("from"); len(f) != 0 {
			particleFrom = f
		}
		if _, err := cid.Decode(particleFrom); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if t := r.URL.Query().Get("to"); len(t) != 0 {
			particleTo = t
		}
		if _, err := cid.Decode(particleTo); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		params := types.NewQueryIsAnyLinkExistParams(particleFrom, particleTo)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryIsAnyLinkExist)
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

func queryEntropyHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particle string

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			particle = v
		}
		if _, err := cid.Decode(particle); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		params := types.NewQueryEntropyParams(particle)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEntropy)
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

func queryLuminosityHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var particle string

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			particle = v
		}
		if _, err := cid.Decode(particle); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		params := types.NewQueryLuminosityParams(particle)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLuminosity)
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

func queryKarmaHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addr string

		if v := r.URL.Query().Get("address"); len(v) != 0 {
			addr = v
		}
		address, err := sdk.AccAddressFromBech32(addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.NewQueryKarmaParams(address)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryKarma)
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


