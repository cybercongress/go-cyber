package proxy

import (
	"errors"
	"net/http"
)

func GetHandlerFn(ctx ProxyContext, endpoint string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		resp, err := ctx.Get(endpoint)
		if err != nil {
			HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}

func GetWithParamHandlerFn(ctx ProxyContext, endpoint string, param string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		addresses, ok := r.URL.Query()[param]

		if !ok || addresses == nil {
			HandleError(errors.New("Cannot find param "+param), w)
			return
		}

		paramValue := ""
		if addresses != nil && len(addresses[0]) == 1 {
			paramValue = addresses[0]
		}

		resp, err := ctx.Get(endpoint + "?" + param + "=\"" + paramValue + "\"")
		if err != nil {
			HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}
