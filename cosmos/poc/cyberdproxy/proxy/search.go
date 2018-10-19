package proxy

import (
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

		resp, err := ctx.Get("/search?cid=\"" + cid + "\"")
		if err != nil {
			HandleError(err, w)
			return
		}

		w.Write(resp)

	}
}
