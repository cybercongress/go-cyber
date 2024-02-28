package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	graphtypes "github.com/cybercongress/go-cyber/v3/x/graph/types"
)

type StakeKeeper interface {
	DetectUsersStakeAmpereChange(ctx sdk.Context) bool
	GetTotalStakesAmpere() map[uint64]uint64
	GetNextAccountNumber(ctx sdk.Context) uint64
}

type GraphIndexedKeeper interface {
	UpdateRankLinks()
	MergeContextLinks(sdk.Context)

	GetOutLinks() graphtypes.Links
	GetInLinks() graphtypes.Links

	GetLinksCount(sdk.Context) uint64
	GetCurrentBlockNewLinks(ctx sdk.Context) []graphtypes.CompactLink
	GetCidsCount(sdk.Context) uint64
}

type GraphKeeper interface {
	GetCidsCount(sdk.Context) uint64
	GetCidNumber(sdk.Context, graphtypes.Cid) (graphtypes.CidNumber, bool)
	GetCid(ctx sdk.Context, num graphtypes.CidNumber) graphtypes.Cid
	GetNeudegs() map[uint64]uint64
	UpdateRankNeudegs()
}
