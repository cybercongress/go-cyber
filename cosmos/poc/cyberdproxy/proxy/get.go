package proxy

import "net/http"

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

		paramValue, err := getSingleParamValue(param, r)
		if err != nil {
			HandleError(err, w)
			return
		}

		resp, err := ctx.Get(endpoint + "?" + param + "=\"" + paramValue + "\"")
		if err != nil {
			HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}
