package core

import (
	"github.com/cybercongress/cyberd/cosmos/poc/http/util"
	"net/http"
)

func GetHandlerFn(ctx ProxyContext, endpoint string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		resp, err := ctx.Get(endpoint)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}

func GetWithParamHandlerFn(ctx ProxyContext, endpoint string, param string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		paramValue, err := util.GetSingleParamValue(param, r)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		resp, err := ctx.Get(endpoint + "?" + param + "=\"" + paramValue + "\"")
		if err != nil {
			util.HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}
