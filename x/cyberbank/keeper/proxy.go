package keeper

import (
	context "context"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

var _ bank.Keeper = (*BankProxyKeeper)(nil)

type BankProxyKeeper struct {
	bk bank.Keeper
	ak authkeeper.AccountKeeper
	ek types.EnergyKeeper

	coinsTransferHooks []types.CoinsTransferHook
}

func WrapBank(bk bank.Keeper) *BankProxyKeeper {
	return &BankProxyKeeper{
		bk:                 bk,
		coinsTransferHooks: make([]types.CoinsTransferHook, 0),
	}
}

func (p *BankProxyKeeper) AddBalanceListener(hook types.CoinsTransferHook) {
	p.coinsTransferHooks = append(p.coinsTransferHooks, hook)
}

func (p *BankProxyKeeper) SetGridKeeper(ek types.EnergyKeeper) {
	p.ek = ek
}

func (p *BankProxyKeeper) SetAccountKeeper(ak authkeeper.AccountKeeper) {
	p.ak = ak
}

func (p *BankProxyKeeper) OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	for _, hook := range p.coinsTransferHooks {
		hook(ctx, from, to)
	}
}

func (p BankProxyKeeper) GetTotalSupplyVolt(ctx sdk.Context) int64 {
	return p.bk.GetSupply(ctx, ctypes.VOLT).Amount.Int64()
}

func (p BankProxyKeeper) GetTotalSupplyAmper(ctx sdk.Context) int64 {
	return p.bk.GetSupply(ctx, ctypes.AMPERE).Amount.Int64()
}

func (p BankProxyKeeper) GetAccountStakePercentageVolt(ctx sdk.Context, addr sdk.AccAddress) float64 {
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

func (p BankProxyKeeper) GetAccountTotalStakeVolt(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return p.bk.GetBalance(ctx, addr, ctypes.VOLT).Amount.Int64() + p.GetRoutedTo(ctx, addr).AmountOf(ctypes.VOLT).Int64()
}

func (p BankProxyKeeper) GetAccountTotalStakeAmper(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return p.bk.GetBalance(ctx, addr, ctypes.AMPERE).Amount.Int64() + p.GetRoutedTo(ctx, addr).AmountOf(ctypes.AMPERE).Int64()
}

func (p BankProxyKeeper) GetRoutedTo(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.ek.GetRoutedToEnergy(ctx, addr)
}

// -----------------------------------------------------------------

func (p *BankProxyKeeper) InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
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

func (p *BankProxyKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SendCoins(ctx, fromAddr, toAddr, amt)
	if err == nil {
		p.OnCoinsTransfer(ctx, fromAddr, toAddr)
	}
	return err
}

func (p *BankProxyKeeper) DenomMetadata(ctx context.Context, request *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	return p.bk.DenomMetadata(ctx, request)
}

func (p *BankProxyKeeper) DenomsMetadata(ctx context.Context, request *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	return p.bk.DenomsMetadata(ctx, request)
}

func (p *BankProxyKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, p.ak.GetModuleAddress(senderModule), recipientAddr)
	}
	return err
}

func (p *BankProxyKeeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)

	//if err == nil {
	//	p.OnCoinsTransfer(ctx, p.ak.GetModuleAddress(senderModule), p.ak.GetModuleAddress(recipientModule))
	//}
	return err
}

func (p *BankProxyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	err := p.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)

	if err == nil {
		p.OnCoinsTransfer(ctx, senderAddr, p.ak.GetModuleAddress(recipientModule))
	}
	return err
}

func (p *BankProxyKeeper) ValidateBalance(ctx sdk.Context, addr sdk.AccAddress) error {
	return p.bk.ValidateBalance(ctx, addr)
}

func (p *BankProxyKeeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return p.bk.HasBalance(ctx, addr, amt)
}

func (p *BankProxyKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.GetAllBalances(ctx, addr)
}

func (p *BankProxyKeeper) GetAccountsBalances(ctx sdk.Context) []banktypes.Balance {
	return p.bk.GetAccountsBalances(ctx)
}

func (p *BankProxyKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return p.bk.GetBalance(ctx, addr, denom)
}

func (p *BankProxyKeeper) LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.LockedCoins(ctx, addr)
}

func (p *BankProxyKeeper) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return p.bk.SpendableCoins(ctx, addr)
}

func (p *BankProxyKeeper) IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) (stop bool)) {
	p.bk.IterateAccountBalances(ctx, addr, cb)
}

func (p *BankProxyKeeper) IterateAllBalances(ctx sdk.Context, cb func(address sdk.AccAddress, coin sdk.Coin) (stop bool)) {
	p.bk.IterateAllBalances(ctx, cb)
}

func (p *BankProxyKeeper) GetParams(ctx sdk.Context) banktypes.Params {
	return p.bk.GetParams(ctx)
}

func (p *BankProxyKeeper) SetParams(ctx sdk.Context, params banktypes.Params) {
	p.bk.SetParams(ctx, params)
}

func (p *BankProxyKeeper) BlockedAddr(addr sdk.AccAddress) bool {
	return p.bk.BlockedAddr(addr)
}

func (p *BankProxyKeeper) InitGenesis(context sdk.Context, state *banktypes.GenesisState) {
	p.bk.InitGenesis(context, state)
}

func (p *BankProxyKeeper) ExportGenesis(context sdk.Context) *banktypes.GenesisState {
	return p.bk.ExportGenesis(context)
}

func (p *BankProxyKeeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	return p.bk.GetSupply(ctx, denom)
}

func (p *BankProxyKeeper) IsSendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool {
	return p.bk.IsSendEnabledCoin(ctx, coin)
}

func (p *BankProxyKeeper) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	return p.bk.IsSendEnabledCoins(ctx, coins...)
}

func (p *BankProxyKeeper) GetPaginatedTotalSupply(ctx sdk.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error) {
	return p.bk.GetPaginatedTotalSupply(ctx, pagination)
}

func (p *BankProxyKeeper) IterateTotalSupply(ctx sdk.Context, cb func(sdk.Coin) bool) {
	p.bk.IterateTotalSupply(ctx, cb)
}

func (p *BankProxyKeeper) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	return p.bk.GetDenomMetaData(ctx, denom)
}

func (p *BankProxyKeeper) SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata) {
	p.bk.SetDenomMetaData(ctx, denomMetaData)
}

func (p *BankProxyKeeper) IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool) {
	p.bk.IterateAllDenomMetaData(ctx, cb)
}

func (p *BankProxyKeeper) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return p.bk.DelegateCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (p *BankProxyKeeper) UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.UndelegateCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}

func (p *BankProxyKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return p.bk.MintCoins(ctx, moduleName, amt)
}

func (p *BankProxyKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return p.bk.BurnCoins(ctx, moduleName, amt)
}

func (p *BankProxyKeeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
}

func (p *BankProxyKeeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	return p.bk.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt)
}

func (p *BankProxyKeeper) Balance(ctx context.Context, request *banktypes.QueryBalanceRequest) (*banktypes.QueryBalanceResponse, error) {
	return p.bk.Balance(ctx, request)
}

func (p *BankProxyKeeper) AllBalances(ctx context.Context, request *banktypes.QueryAllBalancesRequest) (*banktypes.QueryAllBalancesResponse, error) {
	return p.bk.AllBalances(ctx, request)
}

func (p *BankProxyKeeper) TotalSupply(ctx context.Context, request *banktypes.QueryTotalSupplyRequest) (*banktypes.QueryTotalSupplyResponse, error) {
	return p.bk.TotalSupply(ctx, request)
}

func (p *BankProxyKeeper) SupplyOf(ctx context.Context, request *banktypes.QuerySupplyOfRequest) (*banktypes.QuerySupplyOfResponse, error) {
	return p.bk.SupplyOf(ctx, request)
}

func (p *BankProxyKeeper) Params(ctx context.Context, request *banktypes.QueryParamsRequest) (*banktypes.QueryParamsResponse, error) {
	return p.bk.Params(ctx, request)
}

func (p *BankProxyKeeper) SpendableBalances(ctx context.Context, request *banktypes.QuerySpendableBalancesRequest) (*banktypes.QuerySpendableBalancesResponse, error) {
	return p.bk.SpendableBalances(ctx, request)
}

func (p *BankProxyKeeper) HasSupply(ctx sdk.Context, denom string) bool {
	return p.bk.HasSupply(ctx, denom)
}

func (p *BankProxyKeeper) GetAllDenomMetaData(ctx sdk.Context) []banktypes.Metadata {
	denomMetaData := make([]banktypes.Metadata, 0)
	p.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		denomMetaData = append(denomMetaData, metadata)
		return false
	})

	return denomMetaData
}
