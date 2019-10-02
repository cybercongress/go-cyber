package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	cbd "github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/x/link"
)

type StakeKeeper interface {
	FixUserStake() bool
	GetTotalStakes() map[cbd.AccNumber]uint64
}

type LinkIndexedKeeper interface {
	FixLinks()
	EndBlocker() bool

	GetOutLinks() link.Links
	GetInLinks() link.Links

	GetLinksCount(sdk.Context) uint64
	GetCurrentBlockNewLinks() []link.CompactLink
}

type CidNumberKeeper interface {
	GetCidsCount(sdk.Context) uint64
}
