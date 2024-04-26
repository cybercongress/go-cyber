package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v4/x/bandwidth/exported"
	v2 "github.com/cybercongress/go-cyber/v4/x/bandwidth/migrations/v2"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	meter          BandwidthMeter
	legacySubspace exported.Subspace
}

func NewMigrator(bm BandwidthMeter, ss exported.Subspace) Migrator {
	return Migrator{
		meter:          bm,
		legacySubspace: ss,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.Migrate(ctx, ctx.KVStore(m.meter.storeKey), m.legacySubspace, m.meter.cdc)
}
