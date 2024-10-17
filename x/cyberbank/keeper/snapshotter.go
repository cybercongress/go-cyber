package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	snapshot "github.com/cosmos/cosmos-sdk/snapshots/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v5/x/cyberbank/types"
)

var _ snapshot.ExtensionSnapshotter = &CyberbankSnapshotter{}

const SnapshotFormat = 1

type CyberbankSnapshotter struct {
	indexedKeeper *IndexedKeeper
	cms           sdk.MultiStore
}

func NewCyberbankSnapshotter(cms sdk.MultiStore, indexedKeeper *IndexedKeeper) *CyberbankSnapshotter {
	return &CyberbankSnapshotter{
		indexedKeeper: indexedKeeper,
		cms:           cms,
	}
}

func (cs *CyberbankSnapshotter) SnapshotName() string {
	return types.ModuleName
}

func (cs *CyberbankSnapshotter) SnapshotFormat() uint32 {
	return SnapshotFormat
}

func (cs *CyberbankSnapshotter) SupportedFormats() []uint32 {
	// If we support older formats, add them here and handle them in Restore
	return []uint32{SnapshotFormat}
}

func (cs *CyberbankSnapshotter) SnapshotExtension(_ uint64, _ snapshot.ExtensionPayloadWriter) error {
	return nil
}

func (cs *CyberbankSnapshotter) RestoreExtension(height uint64, format uint32, _ snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {
		freshCtx := sdk.NewContext(cs.cms, tmproto.Header{Height: int64(height)}, false, log.NewNopLogger())

		// TODO revisit with get params and case of increased rank computation blocks
		calculationPeriod := int64(5)
		rankRoundBlockNumber := (freshCtx.BlockHeight() / calculationPeriod) * calculationPeriod

		store, err := cs.cms.CacheMultiStoreWithVersion(rankRoundBlockNumber)
		if err != nil {
			println("Error: ", err)
		}
		rankCtx := sdk.NewContext(store, tmproto.Header{Height: rankRoundBlockNumber}, false, log.NewNopLogger())

		cs.indexedKeeper.LoadState(rankCtx, freshCtx)

		return nil
	}
	return snapshot.ErrUnknownFormat
}
