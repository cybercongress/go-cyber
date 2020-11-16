package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

var _ bank.Keeper = (*Proxy)(nil)

type Proxy struct {
	bk	bank.Keeper
	sk	types.StakingKeeper
	sp	types.SupplyKeeper
	ek  types.EnergyKeeper

	coinsTransferHooks []types.CoinsTransferHook
}

func Wrap(bk *bank.Keeper, sk types.StakingKeeper, sp supply.Keeper) *Proxy {
	return &Proxy{
		bk: *bk,
		sk: sk,
		sp: sp,
		coinsTransferHooks: make([]types.CoinsTransferHook, 0),
	}
}

func (p *Proxy) AddHook(hook types.CoinsTransferHook) {
	p.coinsTransferHooks = append(p.coinsTransferHooks, hook)
}

func (k *Proxy) SetEnergyKeeper(ek types.EnergyKeeper) {
	k.ek = ek
}

func (p *Proxy) OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	for _, hook := range p.coinsTransferHooks {
		hook(ctx, from, to)
	}
}

func (k Proxy) GetTotalSupply(ctx sdk.Context) int64 {
	keeperSupply := k.sp.GetSupply(ctx)
	return keeperSupply.GetTotal().AmountOf(ctypes.CYB).Int64()
}

func (k Proxy) GetAccountUnboundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.bk.GetCoins(ctx, addr).AmountOf(ctypes.CYB).Int64()
}

func (k Proxy) GetAccountBoundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, addr)
	boundedStake := int64(0)
	for _, del := range delegations {
		boundedStake += del.GetShares().TruncateInt64()
	}
	return boundedStake
}

func (k Proxy) GetAccountStakePercentage(ctx sdk.Context, addr sdk.AccAddress) float64 {
	a := k.GetAccountTotalStake(ctx, addr)
	aFloat := float64(a)

	b := k.GetTotalSupply(ctx)
	bFloat := float64(b)

	c := aFloat / bFloat
	return c
}

func (k Proxy) GetAccountTotalStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.GetAccountUnboundedStake(ctx, addr) + k.GetAccountBoundedStake(ctx, addr) + k.GetAccountPower(ctx, addr)
}

func (k Proxy) GetAccountPower(ctx sdk.Context, addr sdk.AccAddress) int64 {
	power := k.ek.GetRoutedToEnergy(ctx, addr)
	return power.Int64()
}

// -----------------------------------------------------------------


func (p Proxy) InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) error {
	return p.bk.InputOutputCoins(ctx, inputs, outputs)
}

func (p Proxy) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SendCoins(ctx, fromAddr, toAddr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, fromAddr, toAddr)
	}
	return err
}

func (p Proxy) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	coins, err := p.bk.SubtractCoins(ctx, addr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, nil, addr)
	}
	return coins, err
}

func (p Proxy) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	coins, err := p.bk.AddCoins(ctx, addr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, nil, addr)
	}
	return coins, err
}

func (p Proxy) GetSendEnabled(ctx sdk.Context) bool {
	return p.bk.GetSendEnabled(ctx)
}

func (p Proxy) SetSendEnabled(ctx sdk.Context, enabled bool) {
	p.bk.SetSendEnabled(ctx, enabled)
}

func (p Proxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return p.bk.BlacklistedAddr(addr)
}

func (p Proxy) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, delegatorAddr, nil)
	}
	return err
}

func (p Proxy) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, nil, delegatorAddr)
	}
	return err
}

func (p Proxy) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.GetCoins(ctx, addr)
}

func (p Proxy) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return p.bk.HasCoins(ctx, addr, amt)
}

func (p Proxy) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.SetCoins(ctx, addr, amt)
}


