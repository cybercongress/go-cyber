package core

import (
	"io/ioutil"
	"net/http"
)

func AccountHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		addresses, ok := r.URL.Query()["address"]

		if !ok || len(addresses[0]) < 1 {
			w.WriteHeader(404)
		}

		address := addresses[0]

		resp, _ := http.Get(ctx.NodeUrl + "/account?address=\"" + address + "\"")

		respBytes, _ := ioutil.ReadAll(resp.Body)

		w.Write(respBytes)

	}
}
