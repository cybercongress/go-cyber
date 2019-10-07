module github.com/cybercongress/cyberd

require (
	github.com/arturalbov/atomicf v0.1.1
	github.com/cosmos/cosmos-sdk v0.37.0
	github.com/cosmos/gaia v0.0.0-20190822123916-3c70fee43395
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/ipfs/go-cid v0.0.3
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae // indirect
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.2
	github.com/tendermint/tm-db v0.1.1
	github.com/zondax/ledger-go v0.8.0 // indirect

)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
