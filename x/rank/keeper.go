package rank

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ Keeper = &BaseRankKeeper{}

type BaseRankKeeper struct {
	paramSpace *params.Subspace
}

func NewBaseRankKeeper(paramSpace *params.Subspace) *BaseRankKeeper {
	return &BaseRankKeeper{
		paramSpace: paramSpace,
	}
}
