package core

import (
	"github.com/cybercongress/cyberd/http/util"
	"net/http"
	"net/url"
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

func GetWithParamsHandlerFn(ctx ProxyContext, endpoint string, params []string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		paramValues := make(url.Values)
		for _, param := range params {
			value, err := util.GetSingleParamValue(param, r)
			if err == nil {
				paramValues[param] = []string{"\"" + value + "\""}
			}
		}

		resp, err := ctx.Get(endpoint + "?" + paramValues.Encode())
		if err != nil {
			util.HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}
