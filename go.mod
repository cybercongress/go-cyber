module github.com/cybercongress/go-cyber

go 1.16

require (
	github.com/armon/go-metrics v0.3.8 // indirect
	github.com/CosmWasm/wasmd v0.18.0
	github.com/CosmWasm/wasmvm v0.16.0
	github.com/cosmos/cosmos-sdk v0.42.9
	github.com/cosmos/iavl v0.16.0
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/gravity-devs/liquidity v1.2.9
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/ipfs/go-cid v0.0.7
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.4.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	github.com/tendermint/tendermint v0.34.11
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.38.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/gravity-devs/liquidity => github.com/litvintech/liquidity v1.2.10
