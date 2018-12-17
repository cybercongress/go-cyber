package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBandwidthCost func(msg sdk.Msg) int64

type AcсBandwidth struct {
	Address          sdk.AccAddress `json:"address"`
	RemainedValue    int64          `json:"remained"`
	LastUpdatedBlock int64          `json:"last_updated_block"`
	MaxValue         int64          `json:"max_value"`
}

func (bs *AcсBandwidth) UpdateMax(newValue int64, currentBlock int64, recoveryPeriod int64) {
	bs.Recover(currentBlock, recoveryPeriod)
	bs.MaxValue = newValue
	bs.LastUpdatedBlock = currentBlock

	if bs.RemainedValue > bs.MaxValue {
		bs.RemainedValue = bs.MaxValue
	}
}

func (bs *AcсBandwidth) Recover(currentBlock int64, recoveryPeriod int64) {
	recoverPerBlock := float64(bs.MaxValue) / float64(recoveryPeriod)
	fullRecoveryAmount := float64(bs.MaxValue - bs.RemainedValue)

	recoverAmount := float64(currentBlock-bs.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	bs.RemainedValue = bs.RemainedValue + int64(recoverAmount)
	bs.LastUpdatedBlock = currentBlock
}

func (bs AcсBandwidth) HasEnoughRemained(bandwidthToConsume int64) bool {
	return bs.RemainedValue >= bandwidthToConsume
}

func (bs *AcсBandwidth) Consume(bandwidthToConsume int64) {
	bs.RemainedValue = bs.RemainedValue - bandwidthToConsume
	if bs.RemainedValue < 0 {
		panic("Negative bandwidth!!")
	}
}

func NewGenesisAccBandwidth(address sdk.AccAddress, bandwidth int64) AcсBandwidth {
	return AcсBandwidth{Address: address, RemainedValue: bandwidth, MaxValue: bandwidth, LastUpdatedBlock: 0}
}
