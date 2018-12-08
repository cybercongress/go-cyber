package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	RecoveryPeriod = 100
)

type AccountBandwidthState struct {
	Address          types.AccAddress `json:"address"`
	Remained         int64            `json:"remained"`
	LastUpdatedBlock int64            `json:"last_updated_block"`
}

func (bs *AccountBandwidthState) Recover(maxBandwidth int64, currentBlock int64) {

	recoverPerBlock := maxBandwidth / RecoveryPeriod
	fullRecoveryAmount := maxBandwidth - bs.Remained

	recoverAmount := (currentBlock - bs.LastUpdatedBlock) * recoverPerBlock
	if recoverAmount > fullRecoveryAmount {
		recoverAmount = fullRecoveryAmount
	}

	bs.Remained = bs.Remained + recoverAmount
}

func (bs AccountBandwidthState) HasEnoughRemained(bandwidthToConsume int64) bool {
	return bs.Remained >= bandwidthToConsume
}


func (bs *AccountBandwidthState) Consume(bandwidthToConsume int64) {
	bs.Remained = bs.Remained - bandwidthToConsume
}
