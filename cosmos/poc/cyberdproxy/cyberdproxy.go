package main

import (
	"github.com/cybercongress/cyberd/cosmos/poc/cyberdproxy/core"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	ctx := core.NewProxyContext("http://localhost:26657")

	mux := http.NewServeMux()
	mux.HandleFunc("/link", core.LinkHandlerFn(ctx))
	mux.HandleFunc("/search", core.SearchHandlerFn(ctx))
	mux.HandleFunc("/account", core.AccountHandlerFn(ctx))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":8081", handler)
}
