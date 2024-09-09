package utils

import (
	db "github.com/cometbft/cometbft-db"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	clocktypes "github.com/cybercongress/go-cyber/v4/x/clock/types"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v4/x/tokenfactory/types"
)

func NewContextWithMSVersion(db db.DB, version int64, keys map[string]*storetypes.KVStoreKey) (sdk.Context, error) {
	ms := store.NewCommitMultiStore(db)

	delete(keys, ibcfeetypes.ModuleName)
	delete(keys, consensustypes.ModuleName)
	delete(keys, resourcestypes.ModuleName)
	delete(keys, crisistypes.ModuleName)
	delete(keys, tokenfactorytypes.ModuleName)
	delete(keys, clocktypes.ModuleName)
	delete(keys, nft.ModuleName)
	delete(keys, icahosttypes.StoreKey)
	delete(keys, icacontrollertypes.StoreKey)

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
