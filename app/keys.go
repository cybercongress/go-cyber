package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

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
	distr         *sdk.KVStoreKey
	tDistr        *sdk.TransientStoreKey
	slashing      *sdk.KVStoreKey
	mint          *sdk.KVStoreKey
	params        *sdk.KVStoreKey
	tParams       *sdk.TransientStoreKey
	accBandwidth  *sdk.KVStoreKey
}

func NewCyberdAppDbKeys() CyberdAppDbKeys {

	return CyberdAppDbKeys{

		main:     sdk.NewKVStoreKey(bam.MainStoreKey),
		acc:      sdk.NewKVStoreKey(auth.StoreKey),
		stake:    sdk.NewKVStoreKey(staking.StoreKey),
		tStake:   sdk.NewTransientStoreKey(staking.TStoreKey),
		mint:     sdk.NewKVStoreKey(mint.StoreKey),
		distr:    sdk.NewKVStoreKey(distr.StoreKey),
		tDistr:   sdk.NewTransientStoreKey(distr.TStoreKey),
		slashing: sdk.NewKVStoreKey(slashing.StoreKey),
		params:   sdk.NewKVStoreKey(params.StoreKey),
		tParams:  sdk.NewTransientStoreKey(params.TStoreKey),

		cidNum:        sdk.NewKVStoreKey("cid_index"),
		cidNumReverse: sdk.NewKVStoreKey("cid_index_reverse"),
		links:         sdk.NewKVStoreKey("links"),
		rank:          sdk.NewKVStoreKey("rank"),
		accBandwidth:  sdk.NewKVStoreKey("acc_bandwidth"),
	}
}

func (k CyberdAppDbKeys) GetStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		k.main, k.acc, k.cidNum, k.cidNumReverse, k.links, k.rank, k.stake,
		k.slashing, k.params, k.distr, k.mint, k.accBandwidth,
	}
}

func (k CyberdAppDbKeys) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		k.tStake, k.tParams, k.tDistr,
	}
}
