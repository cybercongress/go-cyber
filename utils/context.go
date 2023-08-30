package utils

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewContextWithMSVersion(db db.DB, version int64, keys map[string]*sdk.KVStoreKey) (sdk.Context, error) {
	ms := store.NewCommitMultiStore(db)

	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}

	err := ms.LoadVersion(version)
	if err != nil {
		return sdk.Context{}, err
	}

	return sdk.NewContext(ms, tmproto.Header{Height: version}, false, nil), nil
}
