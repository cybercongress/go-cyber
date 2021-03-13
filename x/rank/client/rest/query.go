package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	//"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cybercongress/go-cyber/types/query"
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
		"/rank/karmas",
		queryKarmas(cliCtx)).Methods("GET")
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
		var cid string

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			cid = v
		}
		// TODO put check to existence

		params := types.QueryRankRequest{cid}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
		var req types.QuerySearchRequest
		var pr query.PageRequest

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			req.Cid = v
		}
		// TODO put check to existence

		if v := r.URL.Query().Get("page"); len(v) != 0 {
			page, _ := strconv.ParseInt(v, 10, 32)
			pr.Page = uint32(page)
		} else {
			pr.Page = uint32(0)
		}
		if v := r.URL.Query().Get("limit"); len(v) != 0 {
			limit, _ := strconv.ParseInt(v, 10, 32)
			pr.PerPage = uint32(limit)
		} else {
			pr.PerPage = uint32(10)
		}

		params := types.QuerySearchRequest{req.Cid, &pr}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
		var req types.QuerySearchRequest
		var pr query.PageRequest

		if v := r.URL.Query().Get("cid"); len(v) != 0 {
			req.Cid = v
		}
		// TODO put check to existence

		if v := r.URL.Query().Get("page"); len(v) != 0 {
			page, _ := strconv.ParseInt(v, 10, 32)
			pr.Page = uint32(page)
		} else {
			pr.Page = uint32(0)
		}
		if v := r.URL.Query().Get("limit"); len(v) != 0 {
			limit, _ := strconv.ParseInt(v, 10, 32)
			pr.PerPage = uint32(limit)
		} else {
			pr.PerPage = uint32(10)
		}

		params := types.QuerySearchRequest{req.Cid, &pr}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
		var pr query.PageRequest

		if v := r.URL.Query().Get("page"); len(v) != 0 {
			page, _ := strconv.ParseInt(v, 10, 32)
			pr.Page = uint32(page)
		} else {
			pr.Page = uint32(0)
		}
		if v := r.URL.Query().Get("limit"); len(v) != 0 {
			limit, _ := strconv.ParseInt(v, 10, 32)
			pr.PerPage = uint32(limit)
		} else {
			pr.PerPage = uint32(10)
		}

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, pr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
		var req types.QueryIsLinkExistRequest

		if v := r.URL.Query().Get("from"); len(v) != 0 {
			req.From = v
		}
		if v := r.URL.Query().Get("to"); len(v) != 0 {
			req.To = v
		}
		if v := r.URL.Query().Get("address"); len(v) != 0 {
			req.Address = v
		}
		// TODO put check to existence

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
		var req types.QueryIsAnyLinkExistRequest

		if v := r.URL.Query().Get("from"); len(v) != 0 {
			req.From = v
		}
		if v := r.URL.Query().Get("to"); len(v) != 0 {
			req.To = v
		}
		// TODO put check to existence

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

func queryKarmas(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryKarmasRequest

		bz, err := codec.MarshalJSONIndent(cliCtx.LegacyAmino, req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryKarmas)
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


