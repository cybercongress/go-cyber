package util

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tm-db"
)

func NewContextWithMSVersion(db db.DB, version int64, keys ...*sdk.KVStoreKey) (sdk.Context, error) {

	ms := store.NewCommitMultiStore(db)

	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}

	err := ms.LoadVersion(version)

	if err != nil {
		return sdk.Context{}, err
	}

	return sdk.NewContext(ms, types.Header{Height: version}, false, nil), nil
}
