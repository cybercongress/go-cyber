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

func (ab *AccountBandwidth) UpdateMax(newValue uint64, currentBlock uint64, recoveryPeriod uint64) {
	ab.Recover(currentBlock, recoveryPeriod)
	ab.MaxValue = newValue
	ab.LastUpdatedBlock = currentBlock

	if ab.RemainedValue > ab.MaxValue {
		ab.RemainedValue = ab.MaxValue
	}
}

func (ab *AccountBandwidth) Recover(currentBlock uint64, recoveryPeriod uint64) {
	recoverPerBlock := float64(ab.MaxValue) / float64(recoveryPeriod)
	fullRecoveryAmount := float64(ab.MaxValue - ab.RemainedValue)

	recoverAmount := float64(currentBlock-ab.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	ab.RemainedValue = ab.RemainedValue + uint64(recoverAmount)
	ab.LastUpdatedBlock = currentBlock
}

func (ab *AccountBandwidth) Consume(bandwidthToConsume uint64) error {
	ab.RemainedValue = ab.RemainedValue - bandwidthToConsume
	if ab.RemainedValue < 0 {
		return ErrNotEnoughBandwidth
	}
	return nil
}

func (ab *AccountBandwidth) ApplyCharge(bandwidthToAdd uint64) {
	ab.RemainedValue += bandwidthToAdd
}

func (ab AccountBandwidth) HasEnoughRemained(bandwidthToConsume uint64) bool {
	return ab.RemainedValue >= bandwidthToConsume
}

