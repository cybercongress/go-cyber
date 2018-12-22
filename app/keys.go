package app

import sdk "github.com/cosmos/cosmos-sdk/types"

type CyberdAppDbKeys struct {
	main          *sdk.KVStoreKey
	acc           *sdk.KVStoreKey
	accIndex      *sdk.KVStoreKey
	cidNum        *sdk.KVStoreKey
	cidNumReverse *sdk.KVStoreKey
	links         *sdk.KVStoreKey
	rank          *sdk.KVStoreKey
	stake         *sdk.KVStoreKey
	tStake        *sdk.TransientStoreKey
	fees          *sdk.KVStoreKey
	distr         *sdk.KVStoreKey
	slashing      *sdk.KVStoreKey
	params        *sdk.KVStoreKey
	tParams       *sdk.TransientStoreKey
	accBandwidth  *sdk.KVStoreKey
}

func NewCyberdAppDbKeys() CyberdAppDbKeys {
	return CyberdAppDbKeys{
		main:          sdk.NewKVStoreKey("main"),
		acc:           sdk.NewKVStoreKey("acc"),
		cidNum:        sdk.NewKVStoreKey("cid_index"),
		cidNumReverse: sdk.NewKVStoreKey("cid_index_reverse"),
		links:         sdk.NewKVStoreKey("links"),
		rank:          sdk.NewKVStoreKey("rank"),
		stake:         sdk.NewKVStoreKey("stake"),
		fees:          sdk.NewKVStoreKey("fee"),
		tStake:        sdk.NewTransientStoreKey("transient_stake"),
		distr:         sdk.NewKVStoreKey("distr"),
		slashing:      sdk.NewKVStoreKey("slashing"),
		params:        sdk.NewKVStoreKey("params"),
		tParams:       sdk.NewTransientStoreKey("transient_params"),
		accBandwidth:  sdk.NewKVStoreKey("acc_bandwidth"),
	}
}

func (k CyberdAppDbKeys) GetStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		k.main, k.acc, k.cidNum, k.cidNumReverse, k.links, k.rank, k.stake,
		k.slashing, k.params, k.distr, k.fees, k.accBandwidth,
	}
}

func (k CyberdAppDbKeys) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		k.tStake, k.tParams,
	}
}
