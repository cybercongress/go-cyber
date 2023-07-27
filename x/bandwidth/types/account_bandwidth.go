package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGenesisNeuronBandwidth(address sdk.AccAddress, bandwidth uint64) NeuronBandwidth {
	return NeuronBandwidth{
		Neuron:           address.String(),
		RemainedValue:    bandwidth,
		MaxValue:         bandwidth,
		LastUpdatedBlock: 0,
	}
}

func (ab *NeuronBandwidth) UpdateMax(newValue, currentBlock, recoveryPeriod uint64) {
	ab.Recover(currentBlock, recoveryPeriod)
	ab.MaxValue = newValue
	ab.LastUpdatedBlock = currentBlock

	if ab.RemainedValue > ab.MaxValue {
		ab.RemainedValue = ab.MaxValue
	}
}

func (ab *NeuronBandwidth) Recover(currentBlock, recoveryPeriod uint64) {
	recoverPerBlock := float64(ab.MaxValue) / float64(recoveryPeriod)
	fullRecoveryAmount := float64(ab.MaxValue - ab.RemainedValue)

	recoverAmount := float64(currentBlock-ab.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	ab.RemainedValue += uint64(recoverAmount)
	ab.LastUpdatedBlock = currentBlock
}

func (ab *NeuronBandwidth) Consume(bandwidthToConsume uint64) error {
	ab.RemainedValue -= bandwidthToConsume
	if ab.RemainedValue < 0 {
		return ErrNotEnoughBandwidth
	}
	return nil
}

func (ab *NeuronBandwidth) ApplyCharge(bandwidthToAdd uint64) {
	ab.RemainedValue += bandwidthToAdd
}

func (ab NeuronBandwidth) HasEnoughRemained(bandwidthToConsume uint64) bool {
	return ab.RemainedValue >= bandwidthToConsume
}
