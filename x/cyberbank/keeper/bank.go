package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

var _ bankkeeper.Keeper = (*BankProxyKeeper)(nil)

// TODO revisit check that all need accounts indexed

type (
	CoinsTransferHook = func(sdk.Context, []sdk.AccAddress)
	BankProxyKeeper   struct {
		bk        bankkeeper.Keeper
		ak        types.AccountKeeper
		listeners []CoinsTransferHook
	}
)

func WrapBank(bk bankkeeper.Keeper, ak types.AccountKeeper) *BankProxyKeeper {
	return &BankProxyKeeper{
		bk:        bk,
		ak:        ak,
		listeners: make([]CoinsTransferHook, 0),
	}
}

func (pk *BankProxyKeeper) AddBalanceListener(l func(sdk.Context, []sdk.AccAddress)) {
	pk.listeners = append(pk.listeners, l)
}

func (pk BankProxyKeeper) NotifyListeners(ctx sdk.Context, accounts ...sdk.AccAddress) {
	accounts = deduplicate(accounts)
	for _, l := range pk.listeners {
		l(ctx, accounts)
	}
}

func deduplicate(accounts []sdk.AccAddress) []sdk.AccAddress {
	idx := make(map[string]struct{}, len(accounts))
	r := make([]sdk.AccAddress, 0, len(accounts))
	for _, a := range accounts {
		if _, exists := idx[string(a)]; exists {
			continue
		}
		r = append(r, a)
		idx[string(a)] = struct{}{}
	}
	return r
}

func (pk BankProxyKeeper) InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	err := pk.bk.InputOutputCoins(ctx, inputs, outputs)
	if err != nil {
		return err
	}

	accounts := make([]sdk.AccAddress, 0, len(inputs)+len(outputs))
	for _, a := range inputs {
		// invalid addresses were handled before in the wrapped keeper
		addr, _ := sdk.AccAddressFromBech32(a.Address)
		accounts = append(accounts, addr)
	}
	for _, a := range outputs {
		addr, _ := sdk.AccAddressFromBech32(a.Address)
		accounts = append(accounts, addr)
	}

	pk.NotifyListeners(ctx, accounts...)
	return nil
}

func (pk BankProxyKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	err := pk.bk.SendCoins(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}
	pk.NotifyListeners(ctx, fromAddr, toAddr)
	return nil
}

func (pk *BankProxyKeeper) GetPaginatedTotalSupply(ctx sdk.Context, pagination *query.PageRequest,
) (sdk.Coins, *query.PageResponse, error) {
	return pk.bk.GetPaginatedTotalSupply(ctx, pagination)
}

func (pk *BankProxyKeeper) IterateTotalSupply(ctx sdk.Context, cb func(sdk.Coin) bool) {
	pk.bk.IterateTotalSupply(ctx, cb)
}

func (pk BankProxyKeeper) ValidateBalance(ctx sdk.Context, addr sdk.AccAddress) error {
	return pk.bk.ValidateBalance(ctx, addr)
}

func (pk BankProxyKeeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return pk.bk.HasBalance(ctx, addr, amt)
}

func (pk BankProxyKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return pk.bk.GetAllBalances(ctx, addr)
}

func (pk BankProxyKeeper) GetAccountsBalances(ctx sdk.Context) []banktypes.Balance {
	return pk.bk.GetAccountsBalances(ctx)
}

func (pk BankProxyKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return pk.bk.GetBalance(ctx, addr, denom)
}

func (pk BankProxyKeeper) LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return pk.bk.LockedCoins(ctx, addr)
}

func (pk BankProxyKeeper) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return pk.bk.SpendableCoins(ctx, addr)
}

func (pk BankProxyKeeper) IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) (stop bool)) {
	pk.bk.IterateAccountBalances(ctx, addr, cb)
}

func (pk BankProxyKeeper) IterateAllBalances(ctx sdk.Context, cb func(address sdk.AccAddress, coin sdk.Coin) (stop bool)) {
	pk.bk.IterateAllBalances(ctx, cb)
}

func (pk BankProxyKeeper) GetParams(ctx sdk.Context) banktypes.Params {
	return pk.bk.GetParams(ctx)
}

func (pk BankProxyKeeper) IsSendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool {
	return pk.bk.IsSendEnabledCoin(ctx, coin)
}

func (pk BankProxyKeeper) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	return pk.bk.IsSendEnabledCoins(ctx, coins...)
}

func (pk BankProxyKeeper) BlockedAddr(addr sdk.AccAddress) bool {
	return pk.bk.BlockedAddr(addr)
}

func (pk BankProxyKeeper) SetParams(ctx sdk.Context, params banktypes.Params) {
	pk.bk.SetParams(ctx, params)
}

func (pk BankProxyKeeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	return pk.bk.GetSupply(ctx, denom)
}

func (pk *BankProxyKeeper) InitGenesis(ctx sdk.Context, state *banktypes.GenesisState) {
	pk.bk.InitGenesis(ctx, state)
}

func (pk *BankProxyKeeper) ExportGenesis(ctx sdk.Context) *banktypes.GenesisState {
	return pk.bk.ExportGenesis(ctx)
}

func (pk *BankProxyKeeper) MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error {
	return pk.bk.MintCoins(ctx, moduleName, amounts)
}

func (pk *BankProxyKeeper) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	return pk.bk.GetDenomMetaData(ctx, denom)
}

func (pk *BankProxyKeeper) SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata) {
	pk.bk.SetDenomMetaData(ctx, denomMetaData)
}

func (pk *BankProxyKeeper) HasSupply(ctx sdk.Context, denom string) bool {
	return pk.bk.HasSupply(ctx, denom)
}

func (pk *BankProxyKeeper) SpendableBalances(ctx context.Context, request *banktypes.QuerySpendableBalancesRequest) (*banktypes.QuerySpendableBalancesResponse, error) {
	return pk.bk.SpendableBalances(ctx, request)
}

// GetAllDenomMetaData cloned here from the bank keeper as it is not exposed in
// the bank keeper interface
func (pk *BankProxyKeeper) GetAllDenomMetaData(ctx sdk.Context) []banktypes.Metadata {
	denomMetaData := make([]banktypes.Metadata, 0)
	pk.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		denomMetaData = append(denomMetaData, metadata)
		return false
	})

	return denomMetaData
}

func (pk *BankProxyKeeper) IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool) {
	pk.bk.IterateAllDenomMetaData(ctx, cb)
}

func (pk *BankProxyKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	err := pk.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
	if err != nil {
		return err
	}
	// TODO remove notification about module balance change, state breaking because of bandwidth updates
	pk.NotifyListeners(ctx, recipientAddr, pk.ak.GetModuleAddress(senderModule))
	return nil
}

func (pk *BankProxyKeeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	return pk.bk.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
}

func (pk *BankProxyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	err := pk.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
	if err != nil {
		return err
	}
	// TODO remove notification about module balance change, state breaking because of bandwidth updates
	pk.NotifyListeners(ctx, senderAddr, pk.ak.GetModuleAddress(recipientModule))
	return nil
}

func (pk *BankProxyKeeper) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	err := pk.bk.DelegateCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
	if err != nil {
		return err
	}
	//pk.NotifyListeners(ctx, senderAddr)
	return nil
}

func (pk *BankProxyKeeper) UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	err := pk.bk.UndelegateCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
	if err != nil {
		return err
	}
	//pk.NotifyListeners(ctx, recipientAddr)
	return nil
}

func (pk *BankProxyKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return pk.bk.BurnCoins(ctx, moduleName, amt)
}

func (pk *BankProxyKeeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	err := pk.bk.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
	if err != nil {
		return err
	}
	//pk.NotifyListeners(ctx, delegatorAddr)
	return nil
}

func (pk *BankProxyKeeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	err := pk.bk.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt)
	if err != nil {
		return err
	}
	//pk.NotifyListeners(ctx, delegatorAddr)
	return nil
}

func (pk *BankProxyKeeper) Balance(ctx context.Context, request *banktypes.QueryBalanceRequest) (*banktypes.QueryBalanceResponse, error) {
	return pk.bk.Balance(ctx, request)
}

func (pk *BankProxyKeeper) AllBalances(ctx context.Context, request *banktypes.QueryAllBalancesRequest) (*banktypes.QueryAllBalancesResponse, error) {
	return pk.bk.AllBalances(ctx, request)
}

func (pk *BankProxyKeeper) TotalSupply(ctx context.Context, request *banktypes.QueryTotalSupplyRequest) (*banktypes.QueryTotalSupplyResponse, error) {
	return pk.bk.TotalSupply(ctx, request)
}

func (pk *BankProxyKeeper) SupplyOf(ctx context.Context, request *banktypes.QuerySupplyOfRequest) (*banktypes.QuerySupplyOfResponse, error) {
	return pk.bk.SupplyOf(ctx, request)
}

func (pk *BankProxyKeeper) Params(ctx context.Context, request *banktypes.QueryParamsRequest) (*banktypes.QueryParamsResponse, error) {
	return pk.bk.Params(ctx, request)
}

func (pk *BankProxyKeeper) DenomMetadata(ctx context.Context, request *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	return pk.bk.DenomMetadata(ctx, request)
}

func (pk *BankProxyKeeper) DenomsMetadata(ctx context.Context, request *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	return pk.bk.DenomsMetadata(ctx, request)
}

func (pk *BankProxyKeeper) GetBankKeeper() *bankkeeper.Keeper {
	return &pk.bk
}
