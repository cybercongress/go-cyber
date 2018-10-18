package core

import (
	"io/ioutil"
	"net/http"
)

func AccountHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		address, err := getSingleParamValue("address", r)
		if err != nil {
			HandleError(err, w)
			return
		}

		resp, err := ctx.HttpClient.Get(ctx.NodeUrl + "/account?address=\"" + address + "\"")
		if err != nil {
			HandleError(err, w)
			return
		}

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			HandleError(err, w)
			return
		}

		w.Write(respBytes)

	}
}
