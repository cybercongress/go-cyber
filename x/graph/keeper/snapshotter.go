package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	snapshot "github.com/cosmos/cosmos-sdk/snapshots/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v4/x/graph/types"
)

var _ snapshot.ExtensionSnapshotter = &GraphSnapshotter{}

const SnapshotFormat = 1

type GraphSnapshotter struct {
	graphKeeper *GraphKeeper
	indexKeeper *IndexKeeper
	cms         sdk.MultiStore
}

func NewGraphSnapshotter(cms sdk.MultiStore, graphKeeper *GraphKeeper, indexKeeper *IndexKeeper) *GraphSnapshotter {
	return &GraphSnapshotter{
		graphKeeper: graphKeeper,
		indexKeeper: indexKeeper,
		cms:         cms,
	}
}

func (gs *GraphSnapshotter) SnapshotName() string {
	return types.ModuleName
}

func (gs *GraphSnapshotter) SnapshotFormat() uint32 {
	return SnapshotFormat
}

func (gs *GraphSnapshotter) SupportedFormats() []uint32 {
	// If we support older formats, add them here and handle them in Restore
	return []uint32{SnapshotFormat}
}

func (gs *GraphSnapshotter) SnapshotExtension(_ uint64, _ snapshot.ExtensionPayloadWriter) error {
	return nil
}

func (gs *GraphSnapshotter) RestoreExtension(height uint64, format uint32, _ snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {
		freshCtx := sdk.NewContext(gs.cms, tmproto.Header{Height: int64(height)}, false, log.NewNopLogger())
		//calculationPeriod := gs.rankKeeper.GetParams(freshCtx).CalculationPeriod
		//// TODO remove this after upgrade to v4 because on network upgrade block cannot access rank params
		//if calculationPeriod == 0 {
		//	calculationPeriod = int64(5)
		//}
		calculationPeriod := int64(5)
		rankRoundBlockNumber := (freshCtx.BlockHeight() / calculationPeriod) * calculationPeriod
		if rankRoundBlockNumber == 0 && freshCtx.BlockHeight() >= 1 {
			rankRoundBlockNumber = 1
		}

		store, err := gs.cms.CacheMultiStoreWithVersion(rankRoundBlockNumber)
		if err != nil {
			println("Error: ", err)
		}
		rankCtx := sdk.NewContext(store, tmproto.Header{Height: rankRoundBlockNumber}, false, log.NewNopLogger())

		gs.indexKeeper.LoadState(rankCtx, freshCtx)
		gs.graphKeeper.LoadNeudeg(rankCtx, freshCtx)

		return nil
	}
	return snapshot.ErrUnknownFormat
}
