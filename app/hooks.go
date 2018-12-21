package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

//types check
var _ sdk.StakingHooks = Hooks{}

// Combined Staking Hooks
type Hooks struct {
	sh slashing.Hooks
	dh distr.Hooks
}

func NewStakeHooks(dh distr.Hooks, sh slashing.Hooks) Hooks {
	return Hooks{dh: dh, sh: sh}
}

// nolint
func (h Hooks) OnValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.sh.OnValidatorCreated(ctx, valAddr)
	h.dh.OnValidatorCreated(ctx, valAddr)
}
func (h Hooks) OnValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.sh.OnValidatorModified(ctx, valAddr)
	h.dh.OnValidatorModified(ctx, valAddr)
}
func (h Hooks) OnValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.sh.OnValidatorRemoved(ctx, consAddr, valAddr)
	h.dh.OnValidatorRemoved(ctx, consAddr, valAddr)
}
func (h Hooks) OnValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.sh.OnValidatorBonded(ctx, consAddr, valAddr)
	h.dh.OnValidatorBonded(ctx, consAddr, valAddr)
}
func (h Hooks) OnValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.sh.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
	h.dh.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
}
func (h Hooks) OnValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.sh.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
	h.dh.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
}
func (h Hooks) OnDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.sh.OnDelegationCreated(ctx, delAddr, valAddr)
	h.dh.OnDelegationCreated(ctx, delAddr, valAddr)
}
func (h Hooks) OnDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.sh.OnDelegationSharesModified(ctx, delAddr, valAddr)
	h.dh.OnDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h Hooks) OnDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.sh.OnDelegationRemoved(ctx, delAddr, valAddr)
	h.dh.OnDelegationRemoved(ctx, delAddr, valAddr)
}
