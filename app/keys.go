	package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmwasm/wasmd/x/wasm"
)

type cyberdAppDbKeys struct {
	main           *sdk.KVStoreKey
	acc            *sdk.KVStoreKey
	accIndex       *sdk.KVStoreKey
	cidNum         *sdk.KVStoreKey
	cidNumReverse  *sdk.KVStoreKey
	links          *sdk.KVStoreKey
	rank           *sdk.KVStoreKey
	stake          *sdk.KVStoreKey
	tStake         *sdk.TransientStoreKey
	supply         *sdk.KVStoreKey
	distr          *sdk.KVStoreKey
	gov            *sdk.KVStoreKey
	slashing       *sdk.KVStoreKey
	fees           *sdk.KVStoreKey
	params         *sdk.KVStoreKey
	tParams        *sdk.TransientStoreKey
	accBandwidth   *sdk.KVStoreKey
	blockBandwidth *sdk.KVStoreKey
	mint           *sdk.KVStoreKey
	wasm           *sdk.KVStoreKey
}

func NewCyberdAppDbKeys() cyberdAppDbKeys {
	return cyberdAppDbKeys{
		main:     sdk.NewKVStoreKey(bam.MainStoreKey),
		acc:      sdk.NewKVStoreKey(auth.StoreKey),
		stake:    sdk.NewKVStoreKey(staking.StoreKey),
		tStake:   sdk.NewTransientStoreKey(staking.TStoreKey),
		supply:   sdk.NewKVStoreKey(supply.StoreKey),
		fees:     sdk.NewKVStoreKey(auth.FeeCollectorName),
		distr:    sdk.NewKVStoreKey(distr.StoreKey),
		gov:      sdk.NewKVStoreKey(gov.StoreKey),
		slashing: sdk.NewKVStoreKey(slashing.StoreKey),
		params:   sdk.NewKVStoreKey(params.StoreKey),
		tParams:  sdk.NewTransientStoreKey(params.TStoreKey),
		mint:     sdk.NewKVStoreKey(mint.StoreKey),
		wasm:     sdk.NewKVStoreKey(wasm.StoreKey),

		cidNum:         sdk.NewKVStoreKey("cid_index"),
		cidNumReverse:  sdk.NewKVStoreKey("cid_index_reverse"),
		links:          sdk.NewKVStoreKey("links"),
		rank:           sdk.NewKVStoreKey("rank"),
		accBandwidth:   sdk.NewKVStoreKey("acc_bandwidth"),
		blockBandwidth: sdk.NewKVStoreKey("block_spent_bandwidth"),
	}
}

func (k cyberdAppDbKeys) GetStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		k.main, k.acc, k.cidNum, k.cidNumReverse, k.links, k.rank, k.stake, k.supply, k.gov,
		k.slashing, k.params, k.distr, k.fees, k.accBandwidth, k.blockBandwidth, k.mint, k.wasm,
	}
}

func (k cyberdAppDbKeys) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		k.tStake, k.tParams,
	}
}
