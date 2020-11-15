package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AcсountBandwidth struct {
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	RemainedValue    uint64         `json:"remained" yaml:"remained"`
	LastUpdatedBlock uint64         `json:"last_updated_block" yaml:"last_updated_block"`
	MaxValue         uint64         `json:"max_value" yaml:"max_value"`
}

func NewGenesisAccountBandwidth(address sdk.AccAddress, bandwidth uint64) AcсountBandwidth {
	return AcсountBandwidth{
		Address:            address,
		RemainedValue:      bandwidth,
		MaxValue:           bandwidth,
		LastUpdatedBlock:   0,
	}
}

func (bs *AcсountBandwidth) UpdateMax(newValue uint64, currentBlock uint64, recoveryPeriod uint64) {
	bs.Recover(currentBlock, recoveryPeriod)
	bs.MaxValue = newValue
	bs.LastUpdatedBlock = currentBlock

	if bs.RemainedValue > bs.MaxValue {
		bs.RemainedValue = bs.MaxValue
	}
}

func (bs *AcсountBandwidth) Recover(currentBlock uint64, recoveryPeriod uint64) {
	recoverPerBlock := float64(bs.MaxValue) / float64(recoveryPeriod)
	fullRecoveryAmount := float64(bs.MaxValue - bs.RemainedValue)

	recoverAmount := float64(currentBlock-bs.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	bs.RemainedValue = bs.RemainedValue + uint64(recoverAmount)
	bs.LastUpdatedBlock = currentBlock
}

func (bs *AcсountBandwidth) Consume(bandwidthToConsume uint64) {
	bs.RemainedValue = bs.RemainedValue - bandwidthToConsume
	if bs.RemainedValue < 0 {
		panic("Negative bandwidth!")
	}
}

func (bs AcсountBandwidth) HasEnoughRemained(bandwidthToConsume uint64) bool {
	return bs.RemainedValue >= bandwidthToConsume
}

