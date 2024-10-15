package keeper

import (
	context "context"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	ctypes "github.com/cybercongress/go-cyber/v4/types"
	"github.com/cybercongress/go-cyber/v4/x/cyberbank/types"
)

var _ bank.Keeper = (*Proxy)(nil)

type Proxy struct {
	bk bank.Keeper
	ak authkeeper.AccountKeeper
	ek types.EnergyKeeper

	coinsTransferHooks []types.CoinsTransferHook
}

func (p *Proxy) SpendableCoin(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return p.bk.SpendableCoin(ctx, addr, denom)
}

func (p *Proxy) IsSendEnabledDenom(ctx sdk.Context, denom string) bool {
	return p.bk.IsSendEnabledDenom(ctx, denom)
}

func (p *Proxy) GetSendEnabledEntry(ctx sdk.Context, denom string) (banktypes.SendEnabled, bool) {
	return p.bk.GetSendEnabledEntry(ctx, denom)
}

func (p *Proxy) SetSendEnabled(ctx sdk.Context, denom string, value bool) {
	p.bk.SetSendEnabled(ctx, denom, value)
}

func (p *Proxy) SetAllSendEnabled(ctx sdk.Context, sendEnableds []*banktypes.SendEnabled) {
	p.bk.SetAllSendEnabled(ctx, sendEnableds)
}

func (p *Proxy) DeleteSendEnabled(ctx sdk.Context, denoms ...string) {
	p.bk.DeleteSendEnabled(ctx, denoms...)
}

func (p *Proxy) IterateSendEnabledEntries(ctx sdk.Context, cb func(denom string, sendEnabled bool) (stop bool)) {
	p.bk.IterateSendEnabledEntries(ctx, cb)
}

func (p *Proxy) GetAllSendEnabledEntries(ctx sdk.Context) []banktypes.SendEnabled {
	return p.bk.GetAllSendEnabledEntries(ctx)
}

func (p *Proxy) GetBlockedAddresses() map[string]bool {
	return p.bk.GetBlockedAddresses()
}

func (p *Proxy) GetAuthority() string {
	return p.bk.GetAuthority()
}

func (p *Proxy) WithMintCoinsRestriction(fn bank.MintingRestrictionFn) bank.BaseKeeper {
	return p.bk.WithMintCoinsRestriction(fn)
}

func (p *Proxy) HasDenomMetaData(ctx sdk.Context, denom string) bool {
	return p.bk.HasDenomMetaData(ctx, denom)
}

func (p *Proxy) GetAllDenomMetaData(ctx sdk.Context) []banktypes.Metadata {
	return p.bk.GetAllDenomMetaData(ctx)
}

func (p *Proxy) SpendableBalanceByDenom(ctx context.Context, request *banktypes.QuerySpendableBalanceByDenomRequest) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	return p.bk.SpendableBalanceByDenom(ctx, request)
}

func (p *Proxy) DenomOwners(ctx context.Context, request *banktypes.QueryDenomOwnersRequest) (*banktypes.QueryDenomOwnersResponse, error) {
	return p.bk.DenomOwners(ctx, request)
}

func (p *Proxy) SendEnabled(ctx context.Context, request *banktypes.QuerySendEnabledRequest) (*banktypes.QuerySendEnabledResponse, error) {
	return p.bk.SendEnabled(ctx, request)
}

func Wrap(bk bank.Keeper) *Proxy {
	return &Proxy{
		bk:                 bk,
		coinsTransferHooks: make([]types.CoinsTransferHook, 0),
	}
}

func (p *Proxy) AddHook(hook types.CoinsTransferHook) {
	p.coinsTransferHooks = append(p.coinsTransferHooks, hook)
}

func (p *Proxy) SetGridKeeper(ek types.EnergyKeeper) {
	p.ek = ek
}

func (p *Proxy) SetAccountKeeper(ak authkeeper.AccountKeeper) {
	p.ak = ak
}

func (p *Proxy) OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	for _, hook := range p.coinsTransferHooks {
		hook(ctx, from, to)
	}
}

func (p Proxy) GetTotalSupplyVolt(ctx sdk.Context) int64 {
	return p.bk.GetSupply(ctx, ctypes.VOLT).Amount.Int64()
}

func (p Proxy) GetTotalSupplyAmper(ctx sdk.Context) int64 {
	return p.bk.GetSupply(ctx, ctypes.AMPERE).Amount.Int64()
}

func (p Proxy) GetAccountStakePercentageVolt(ctx sdk.Context, addr sdk.AccAddress) float64 {
	a := p.GetAccountTotalStakeVolt(ctx, addr)
	aFloat := float64(a)

	b := p.GetTotalSupplyVolt(ctx)
	bFloat := float64(b)

	c := aFloat / bFloat

	if math.IsNaN(c) {
		return 0
	}
	return c
}

func (p Proxy) GetAccountTotalStakeVolt(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return p.bk.GetBalance(ctx, addr, ctypes.VOLT).Amount.Int64() + p.GetRoutedTo(ctx, addr).AmountOf(ctypes.VOLT).Int64()
}

func (p Proxy) GetAccountTotalStakeAmper(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return p.bk.GetBalance(ctx, addr, ctypes.AMPERE).Amount.Int64() + p.GetRoutedTo(ctx, addr).AmountOf(ctypes.AMPERE).Int64()
}

func (p Proxy) GetRoutedTo(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.ek.GetRoutedToEnergy(ctx, addr)
}

// -----------------------------------------------------------------

func (p *Proxy) InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	err := p.bk.InputOutputCoins(ctx, inputs, outputs)
	if err == nil {
		for _, i := range inputs {
			inAddress, _ := sdk.AccAddressFromBech32(i.Address)
			p.OnCoinsTransfer(ctx, inAddress, nil)
		}
		for _, j := range outputs {
			outAddress, _ := sdk.AccAddressFromBech32(j.Address)
			p.OnCoinsTransfer(ctx, nil, outAddress)
		}
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

func (p *Proxy) DenomMetadata(ctx context.Context, request *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	return p.bk.DenomMetadata(ctx, request)
}

func (p *Proxy) DenomsMetadata(ctx context.Context, request *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	return p.bk.DenomsMetadata(ctx, request)
}

func (p *Proxy) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, p.ak.GetModuleAddress(senderModule), recipientAddr)
	}
	return err
}

func (p *Proxy) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
	return err
}

func (p *Proxy) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, senderAddr, p.ak.GetModuleAddress(recipientModule))
	}
	return err
}

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

func (p *Proxy) GetParams(ctx sdk.Context) banktypes.Params {
	return p.bk.GetParams(ctx)
}

func (p *Proxy) SetParams(ctx sdk.Context, params banktypes.Params) error {
	return p.bk.SetParams(ctx, params)
}

func (p *Proxy) BlockedAddr(addr sdk.AccAddress) bool {
	return p.bk.BlockedAddr(addr)
}

func (p *Proxy) InitGenesis(context sdk.Context, state *banktypes.GenesisState) {
	p.bk.InitGenesis(context, state)
}

func (p *Proxy) ExportGenesis(context sdk.Context) *banktypes.GenesisState {
	return p.bk.ExportGenesis(context)
}

func (p *Proxy) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	return p.bk.GetSupply(ctx, denom)
}

func (p *Proxy) IsSendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool {
	return p.bk.IsSendEnabledCoin(ctx, coin)
}

func (p *Proxy) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	return p.bk.IsSendEnabledCoins(ctx, coins...)
}

func (p *Proxy) GetPaginatedTotalSupply(ctx sdk.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error) {
	return p.bk.GetPaginatedTotalSupply(ctx, pagination)
}

func (p *Proxy) IterateTotalSupply(ctx sdk.Context, cb func(sdk.Coin) bool) {
	p.bk.IterateTotalSupply(ctx, cb)
}

func (p *Proxy) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	return p.bk.GetDenomMetaData(ctx, denom)
}

func (p *Proxy) SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata) {
	p.bk.SetDenomMetaData(ctx, denomMetaData)
}

func (p *Proxy) IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool) {
	p.bk.IterateAllDenomMetaData(ctx, cb)
}

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

func (p *Proxy) SpendableBalances(ctx context.Context, request *banktypes.QuerySpendableBalancesRequest) (*banktypes.QuerySpendableBalancesResponse, error) {
	return p.bk.SpendableBalances(ctx, request)
}

func (p *Proxy) HasSupply(ctx sdk.Context, denom string) bool {
	return p.bk.HasSupply(ctx, denom)
}
