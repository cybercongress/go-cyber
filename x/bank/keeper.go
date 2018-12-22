package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cybercongress/cyberd/types/coin"
)

type Keeper struct {
	bank.Keeper

	ak auth.AccountKeeper
	sk *stake.Keeper

	coinsTransferHooks []CoinsTransferHook
}

func NewBankKeeper(ak auth.AccountKeeper, sk *stake.Keeper) Keeper {
	return Keeper{
		Keeper:             bank.NewBaseKeeper(ak),
		ak:                 ak,
		sk:                 sk,
		coinsTransferHooks: make([]CoinsTransferHook, 0),
	}
}

func (k *Keeper) AddHook(hook CoinsTransferHook) {
	k.coinsTransferHooks = append(k.coinsTransferHooks, hook)
}

/* Override methods */
// sdk acc keeper is not interface yet
func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	coins, tags, err := k.Keeper.AddCoins(ctx, addr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, addr)
	}
	return coins, tags, err
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	coins, tags, err := k.Keeper.SubtractCoins(ctx, addr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, addr)
	}
	return coins, tags, err
}

func (k Keeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	err := k.Keeper.SetCoins(ctx, addr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, addr)
	}
	return err
}

func (k Keeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	tags, err := k.Keeper.SendCoins(ctx, fromAddr, toAddr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, fromAddr, toAddr)
	}
	return tags, err
}

func (k Keeper) InputOutputCoins(
	ctx sdk.Context, inputs []bank.Input, outputs []bank.Output,
) (sdk.Tags, sdk.Error) {
	tags, err := k.Keeper.InputOutputCoins(ctx, inputs, outputs)
	if err == nil {
		for _, i := range inputs {
			k.onCoinsTransfer(ctx, i.Address, nil)
		}
		for _, j := range outputs {
			k.onCoinsTransfer(ctx, nil, j.Address)
		}
	}
	return tags, err
}

func (k Keeper) onCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	for _, hook := range k.coinsTransferHooks {
		hook(ctx, from, to)
	}
}

/* Query methods */

func (k Keeper) GetAccountUnboundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	acc := k.ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0
	}
	return acc.GetCoins().AmountOf(coin.CBD).Int64()
}

func (k Keeper) GetAccountBoundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, addr)
	boundedStake := int64(0)
	for _, del := range delegations {
		boundedStake += del.Shares.TruncateInt64()
	}
	return boundedStake
}

// Returns bounded plus unbounded account tokens
func (k Keeper) GetAccountTotalStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.GetAccountUnboundedStake(ctx, addr) + k.GetAccountBoundedStake(ctx, addr)
}

func (k Keeper) GetAccStakePercentage(ctx sdk.Context, addr sdk.AccAddress) float64 {
	return float64(k.GetAccountTotalStake(ctx, addr)) / float64(k.GetTotalSupply(ctx))
}

func (k Keeper) GetTotalSupply(ctx sdk.Context) int64 {
	pool := k.sk.GetPool(ctx)
	return pool.BondedTokens.RoundInt64() + pool.LooseTokens.RoundInt64()
}
