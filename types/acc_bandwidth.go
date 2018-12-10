package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

type AccountBandwidth struct {
	Address          types.AccAddress `json:"address"`
	RemainedValue    int64            `json:"remained"`
	LastUpdatedBlock int64            `json:"last_updated_block"`
	MaxValue         int64            `json:"max_value"`
}

func (bs *AccountBandwidth) UpdateMax(newValue int64, currentBlock int64, recoveryPeriod int64) {
	bs.Recover(currentBlock, recoveryPeriod)
	bs.MaxValue = newValue
	bs.LastUpdatedBlock = currentBlock

	if bs.RemainedValue > bs.MaxValue {
		bs.RemainedValue = bs.MaxValue
	}
}

func (bs *AccountBandwidth) Recover(currentBlock int64, recoveryPeriod int64) {
	recoverPerBlock := bs.MaxValue / recoveryPeriod
	fullRecoveryAmount := bs.MaxValue - bs.RemainedValue

	recoverAmount := (currentBlock - bs.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	bs.RemainedValue = bs.RemainedValue + recoverAmount
	bs.LastUpdatedBlock = currentBlock
}

func (bs AccountBandwidth) HasEnoughRemained(bandwidthToConsume int64) bool {
	return bs.RemainedValue >= bandwidthToConsume
}

//TODO: Add check for remained bandwidth
func (bs *AccountBandwidth) Consume(bandwidthToConsume int64) {
	bs.RemainedValue = bs.RemainedValue - bandwidthToConsume
}
