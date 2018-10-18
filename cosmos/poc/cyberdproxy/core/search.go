package core

import (
	"io/ioutil"
	"net/http"
)

func SearchHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		cid, err := getSingleParamValue("cid", r)
		if err != nil {
			HandleError(err, w)
			return
		}

		resp, err := ctx.HttpClient.Get(ctx.NodeUrl + "/search?cid=\"" + cid + "\"")
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
