package proxy

import (
	"net/http"
)

func HealthHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		resp, err := ctx.Get("/health")
		if err != nil {
			HandleError(err, w)
			return
		}

		w.Write(resp)
	}
}
