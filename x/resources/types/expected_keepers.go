package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SupplyKeeper defines the expected supply keeper
//type SupplyKeeper interface {
//	GetModuleAddress(name string) sdk.AccAddress
//	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
//	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
//	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
//	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
//	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
//}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

// StakingKeeper defines the expected keeper interface for the staking keeper
type StakingKeeper interface {
	IterateLastValidators(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateAllDelegations(ctx sdk.Context, cb func(delegation stakingtypes.Delegation) (stop bool))
	GetBondedPool(ctx sdk.Context) (bondedPool authtypes.ModuleAccountI)
	BondDenom(ctx sdk.Context) (res string)
}

type BankKeeper interface {
	GetSupply(ctx sdk.Context) bankexported.SupplyI
	//GetCoins(ctx sdk.Context) bankexported.SupplyI
}

