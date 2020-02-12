package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cybercongress/cyberd/x/rank/exported"
	"github.com/cybercongress/cyberd/x/rank/internal/types"
)

var _ exported.Keeper = &BaseRankKeeper{}

type BaseRankKeeper struct {
	paramSpace params.Subspace
}

func NewBaseRankKeeper(paramSpace params.Subspace) *BaseRankKeeper {
	return &BaseRankKeeper{
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}
