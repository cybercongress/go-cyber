package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/rank/exported"
	v2 "github.com/cybercongress/go-cyber/v5/x/rank/migrations/v2"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper         StateKeeper
	legacySubspace exported.Subspace
}

func NewMigrator(k StateKeeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.Migrate(ctx, ctx.KVStore(m.keeper.storeKey), m.legacySubspace, m.keeper.cdc)
}
