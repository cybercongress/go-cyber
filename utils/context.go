package utils

import (
	db "github.com/cometbft/cometbft-db"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewContextWithMSVersion(db db.DB, version int64, keys map[string]*storetypes.KVStoreKey) (sdk.Context, error) {
	ms := store.NewCommitMultiStore(db)

	delete(keys, "feeibc")
	delete(keys, "consensus")
	delete(keys, "resources")
	delete(keys, "crisis")
	delete(keys, "tokenfactory")
	delete(keys, "clock")
	delete(keys, "nft")

	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, nil)
	}

	err := ms.LoadVersion(version)
	if err != nil {
		// TODO put comment here because when added new module during migration this module have 0 store's version
		return sdk.Context{}, err
	}

	return sdk.NewContext(ms, tmproto.Header{Height: version}, false, nil), nil
}
