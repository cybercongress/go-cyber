package keeper

import (
	context "context"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

var _ bank.Keeper = (*Proxy)(nil)

type Proxy struct {
	bk	bank.Keeper
	ak  authkeeper.AccountKeeper
	ek  types.EnergyKeeper

	coinsTransferHooks []types.CoinsTransferHook
}

func Wrap(bk *bank.Keeper) *Proxy {
	return &Proxy{
		bk: *bk,
		coinsTransferHooks: make([]types.CoinsTransferHook, 0),
	}
}

func (p *Proxy) AddHook(hook types.CoinsTransferHook) {
	p.coinsTransferHooks = append(p.coinsTransferHooks, hook)
}

func (k *Proxy) SetEnergyKeeper(ek types.EnergyKeeper) {
	k.ek = ek
}

func (k *Proxy) SetAccountKeeper(ak authkeeper.AccountKeeper) {
	k.ak = ak
}

func (p *Proxy) OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	for _, hook := range p.coinsTransferHooks {
		hook(ctx, from, to)
	}
}

func (k Proxy) GetTotalSupplyVolt(ctx sdk.Context) int64 {
	return k.bk.GetSupply(ctx).GetTotal().AmountOf(ctypes.VOLT).Int64()
}

func (k Proxy) GetTotalSupplyAmper(ctx sdk.Context) int64 {
	return k.bk.GetSupply(ctx).GetTotal().AmountOf(ctypes.AMPER).Int64()
}

func (k Proxy) GetAccountUnboundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.bk.GetBalance(ctx, addr, ctypes.CYB).Amount.Int64()
}

func (k Proxy) GetAccountStakePercentageVolt(ctx sdk.Context, addr sdk.AccAddress) float64 {
	a := k.GetAccountTotalStakeVolt(ctx, addr)
	aFloat := float64(a)

	b := k.GetTotalSupplyVolt(ctx)
	bFloat := float64(b)

	c := aFloat / bFloat

	if math.IsNaN(c) { return 0 }
	return c
}

func (k Proxy) GetAccountTotalStakeVolt(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.bk.GetBalance(ctx, addr, ctypes.VOLT).Amount.Int64() + k.GetRoutedTo(ctx, addr).AmountOf(ctypes.VOLT).Int64()
}

func (k Proxy) GetAccountTotalStakeAmper(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.bk.GetBalance(ctx, addr, ctypes.AMPER).Amount.Int64() + k.GetRoutedTo(ctx, addr).AmountOf(ctypes.AMPER).Int64()
}

//func (k Proxy) GetAccountPower(ctx sdk.Context, addr sdk.AccAddress) int64 {
//	power := k.ek.GetRoutedToEnergy(ctx, addr)
//	// TODO return
//	c := sdk.Coin{}
//	if power == c {
//		return 0
//	} else {
//		return power.Amount.Int64()
//	}
//}

func (k Proxy) GetRoutedTo(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.ek.GetRoutedToEnergy(ctx, addr)
}

// -----------------------------------------------------------------

func (p *Proxy) InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	err := p.bk.InputOutputCoins(ctx, inputs, outputs)
	if err == nil {
		for _, i := range inputs {
			p.OnCoinsTransfer(ctx, sdk.AccAddress(i.Address), nil)
		}
		for _, j := range outputs {
			p.OnCoinsTransfer(ctx, nil, sdk.AccAddress(j.Address))
		}
	}
	return err
}

func (p *Proxy) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SubtractCoins(ctx, addr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, addr, nil)
	}
	return err
}

func (p Proxy) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.AddCoins(ctx, addr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, nil, addr)
	}
	return err
}

func (p *Proxy) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SendCoins(ctx, fromAddr, toAddr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, fromAddr, toAddr)
	}
	return err
}

// -----------------------------------------------------------------

func (p *Proxy) DenomMetadata(ctx context.Context, request *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	return p.bk.DenomMetadata(ctx, request)
}

func (p *Proxy) DenomsMetadata(ctx context.Context, request *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	return p.bk.DenomsMetadata(ctx, request)
}

// -----------------------------------------------------------------

func (p *Proxy) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, p.ak.GetModuleAddress(senderModule), recipientAddr)
	}
	return err
}

func (p *Proxy) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, p.ak.GetModuleAddress(senderModule), p.ak.GetModuleAddress(recipientModule))
	}
	return err
}

func (p *Proxy) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, senderAddr, p.ak.GetModuleAddress(recipientModule))
	}
	return err
}

// -----------------------------------------------------------------

func (p *Proxy) ValidateBalance(ctx sdk.Context, addr sdk.AccAddress) error {
	return p.bk.ValidateBalance(ctx, addr)
}

func (p *Proxy) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return p.bk.HasBalance(ctx, addr, amt)
}

func (p *Proxy) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.GetAllBalances(ctx, addr)
}

func (p *Proxy) GetAccountsBalances(ctx sdk.Context) []banktypes.Balance {
	return p.bk.GetAccountsBalances(ctx)
}

func (p *Proxy) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return p.bk.GetBalance(ctx, addr, denom)
}

func (p *Proxy) LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.LockedCoins(ctx, addr)
}

func (p *Proxy) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.SpendableCoins(ctx, addr)
}

func (p *Proxy) IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) (stop bool)) {
	p.bk.IterateAccountBalances(ctx, addr, cb)
}

func (p *Proxy) IterateAllBalances(ctx sdk.Context, cb func(address sdk.AccAddress, coin sdk.Coin) (stop bool)) {
	p.bk.IterateAllBalances(ctx, cb)
}

func (p *Proxy) SetBalance(ctx sdk.Context, addr sdk.AccAddress, balance sdk.Coin) error {
	return p.bk.SetBalance(ctx, addr, balance)
}

func (p *Proxy) SetBalances(ctx sdk.Context, addr sdk.AccAddress, balances sdk.Coins) error {
	return p.bk.SetBalances(ctx, addr, balances)
}

func (p *Proxy) GetParams(ctx sdk.Context) banktypes.Params {
	return p.bk.GetParams(ctx)
}

func (p *Proxy) SetParams(ctx sdk.Context, params banktypes.Params) {
	p.bk.SetParams(ctx, params)
}

func (p *Proxy) SendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool {
	return p.bk.SendEnabledCoin(ctx, coin)
}

func (p *Proxy) SendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	return p.bk.SendEnabledCoins(ctx, coins...)
}

func (p *Proxy) BlockedAddr(addr sdk.AccAddress) bool {
	return p.bk.BlockedAddr(addr)
}

// -----------------------------------------------------------------

func (p *Proxy) InitGenesis(context sdk.Context, state *banktypes.GenesisState) {
	p.bk.InitGenesis(context, state)
}

func (p *Proxy) ExportGenesis(context sdk.Context) *banktypes.GenesisState {
	return p.bk.ExportGenesis(context)
}

func (p *Proxy) GetSupply(ctx sdk.Context) bankexported.SupplyI {
	return p.bk.GetSupply(ctx)
}

func (p *Proxy) SetSupply(ctx sdk.Context, supply bankexported.SupplyI) {
	p.bk.SetSupply(ctx, supply)
}

func (p *Proxy) GetDenomMetaData(ctx sdk.Context, denom string) banktypes.Metadata {
	return p.bk.GetDenomMetaData(ctx, denom)
}

func (p *Proxy) SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata) {
	p.bk.SetDenomMetaData(ctx, denomMetaData)
}

func (p *Proxy) IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool) {
	p.bk.IterateAllDenomMetaData(ctx, cb)
}

// ----------------------------------------------

func (p *Proxy) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return p.bk.DelegateCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (p *Proxy) UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.UndelegateCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}

func (p *Proxy) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return p.bk.MintCoins(ctx, moduleName, amt)
}

func (p *Proxy) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return p.bk.BurnCoins(ctx, moduleName, amt)
}

func (p *Proxy) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
}

func (p *Proxy) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt)
}

// ----------------------------------------------

func (p *Proxy) MarshalSupply(supplyI bankexported.SupplyI) ([]byte, error) {
	return p.bk.MarshalSupply(supplyI)
}

func (p *Proxy) UnmarshalSupply(bz []byte) (bankexported.SupplyI, error) {
	return p.bk.UnmarshalSupply(bz)
}

func (p *Proxy) Balance(ctx context.Context, request *banktypes.QueryBalanceRequest) (*banktypes.QueryBalanceResponse, error) {
	return p.bk.Balance(ctx, request)
}

func (p *Proxy) AllBalances(ctx context.Context, request *banktypes.QueryAllBalancesRequest) (*banktypes.QueryAllBalancesResponse, error) {
	return p.bk.AllBalances(ctx, request)
}

func (p *Proxy) TotalSupply(ctx context.Context, request *banktypes.QueryTotalSupplyRequest) (*banktypes.QueryTotalSupplyResponse, error) {
	return p.bk.TotalSupply(ctx, request)
}

func (p *Proxy) SupplyOf(ctx context.Context, request *banktypes.QuerySupplyOfRequest) (*banktypes.QuerySupplyOfResponse, error) {
	return p.bk.SupplyOf(ctx, request)
}

func (p *Proxy) Params(ctx context.Context, request *banktypes.QueryParamsRequest) (*banktypes.QueryParamsResponse, error) {
	return p.bk.Params(ctx, request)
}