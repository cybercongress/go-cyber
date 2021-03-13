package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGenesisAccountBandwidth(address sdk.AccAddress, bandwidth uint64) AccountBandwidth {
	return AccountBandwidth{
		Address:            address.String(),
		RemainedValue:      bandwidth,
		MaxValue:           bandwidth,
		LastUpdatedBlock:   0,
	}
}

func (bs *AccountBandwidth) UpdateMax(newValue uint64, currentBlock uint64, recoveryPeriod uint64) {
	bs.Recover(currentBlock, recoveryPeriod)
	bs.MaxValue = newValue
	bs.LastUpdatedBlock = currentBlock

	if bs.RemainedValue > bs.MaxValue {
		bs.RemainedValue = bs.MaxValue
	}
}

func (bs *AccountBandwidth) Recover(currentBlock uint64, recoveryPeriod uint64) {
	recoverPerBlock := float64(bs.MaxValue) / float64(recoveryPeriod)
	fullRecoveryAmount := float64(bs.MaxValue - bs.RemainedValue)

	recoverAmount := float64(currentBlock-bs.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	bs.RemainedValue = bs.RemainedValue + uint64(recoverAmount)
	bs.LastUpdatedBlock = currentBlock
}

func (bs *AccountBandwidth) Consume(bandwidthToConsume uint64) {
	bs.RemainedValue = bs.RemainedValue - bandwidthToConsume
	if bs.RemainedValue < 0 {
		panic("Negative bandwidth!")
	}
}

func (bs AccountBandwidth) HasEnoughRemained(bandwidthToConsume uint64) bool {
	return bs.RemainedValue >= bandwidthToConsume
}

