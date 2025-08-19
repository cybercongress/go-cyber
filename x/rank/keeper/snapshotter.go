package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	snapshot "github.com/cosmos/cosmos-sdk/snapshots/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v6/x/rank/types"
)

var _ snapshot.ExtensionSnapshotter = &RankSnapshotter{}

const SnapshotFormat = 1

type RankSnapshotter struct {
	keeper *StateKeeper
	cms    sdk.MultiStore
}

func NewRankSnapshotter(cms sdk.MultiStore, keeper *StateKeeper) *RankSnapshotter {
	return &RankSnapshotter{
		keeper: keeper,
		cms:    cms,
	}
}

func (rs *RankSnapshotter) SnapshotName() string {
	return types.ModuleName
}

func (rs *RankSnapshotter) SnapshotFormat() uint32 {
	return SnapshotFormat
}

func (rs *RankSnapshotter) SupportedFormats() []uint32 {
	// If we support older formats, add them here and handle them in Restore
	return []uint32{SnapshotFormat}
}

func (rs *RankSnapshotter) SnapshotExtension(_ uint64, _ snapshot.ExtensionPayloadWriter) error {
	return nil
}

func (rs *RankSnapshotter) RestoreExtension(height uint64, format uint32, _ snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {

		freshCtx := sdk.NewContext(rs.cms, tmproto.Header{Height: int64(height)}, false, log.NewNopLogger())
		rs.keeper.LoadState(freshCtx)
		rs.keeper.StartRankCalculation(freshCtx)

		return nil
	}
	return snapshot.ErrUnknownFormat
}
