package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	snapshot "github.com/cosmos/cosmos-sdk/snapshots/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
)

var _ snapshot.ExtensionSnapshotter = &BandwidthSnapshotter{}

const SnapshotFormat = 1

type BandwidthSnapshotter struct {
	keeper *BandwidthMeter
	cms    sdk.MultiStore
}

func NewBandwidthSnapshotter(cms sdk.MultiStore, keeper *BandwidthMeter) *BandwidthSnapshotter {
	return &BandwidthSnapshotter{
		keeper: keeper,
		cms:    cms,
	}
}

func (bs *BandwidthSnapshotter) SnapshotName() string {
	return types.ModuleName
}

func (bs *BandwidthSnapshotter) SnapshotFormat() uint32 {
	return SnapshotFormat
}

func (bs *BandwidthSnapshotter) SupportedFormats() []uint32 {
	// If we support older formats, add them here and handle them in Restore
	return []uint32{SnapshotFormat}
}

func (bs *BandwidthSnapshotter) SnapshotExtension(_ uint64, _ snapshot.ExtensionPayloadWriter) error {
	return nil
}

func (bs *BandwidthSnapshotter) RestoreExtension(height uint64, format uint32, _ snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {
		ctx := sdk.NewContext(bs.cms, tmproto.Header{Height: int64(height)}, false, log.NewNopLogger())
		bs.keeper.LoadState(ctx)

		return nil
	}
	return snapshot.ErrUnknownFormat
}
