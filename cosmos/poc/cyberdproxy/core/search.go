package core

import (
	"io/ioutil"
	"net/http"
)

func SearchHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		cids, ok := r.URL.Query()["cid"]

		if !ok || len(cids[0]) < 1 {
			w.WriteHeader(404)
		}

		cid := cids[0]

		resp, _ := http.Get(ctx.NodeUrl + "/search?cid=\"" + cid + "\"")

		respBytes, _ := ioutil.ReadAll(resp.Body)

		w.Write(respBytes)

	}
}
