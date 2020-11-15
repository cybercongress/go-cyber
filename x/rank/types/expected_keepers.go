package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/link"
)

type StakeKeeper interface {
	DetectUsersStakeChange(ctx sdk.Context) bool
	GetTotalStakes() map[uint64]uint64
}

type GraphIndexedKeeper interface {
	FixLinks()
	EndBlocker() bool

	GetOutLinks() link.Links
	GetInLinks() link.Links

	GetLinksCount(sdk.Context) uint64
	GetCurrentBlockNewLinks() []link.CompactLink
	GetCidsCount(sdk.Context) uint64
}

type GraphKeeper interface {
	GetCidsCount(sdk.Context) uint64
}
