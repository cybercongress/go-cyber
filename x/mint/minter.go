package mint

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake/keeper"
)

// mint new cbd tokens for every block
type Minter struct {
	fck         auth.FeeCollectionKeeper
	stakeKeeper keeper.Keeper
	blockReward uint64
}

func NewMinter(fck auth.FeeCollectionKeeper, stakeKeeper keeper.Keeper) Minter {

	// for 1sec block and 1kkk genesis will be 31
	blockReward := (GenesisSupply / BlocksPerYear) * InflationRatePerYear
	return Minter{
		fck:         fck,
		stakeKeeper: stakeKeeper,
		blockReward: uint64(blockReward),
	}
}

func (m *Minter) BlockReward() uint64 {
	return m.blockReward
}
