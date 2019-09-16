package acc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/x/acc/types"
)

var (
	_ types.AccountIndexKeeper = AccountIndexKeeper{}
)

type AccountIndexKeeper struct {
	auth.AccountKeeper

	accounts map[types.AccNumber]sdk.AccAddress
}

func NewAccountIndexKeeper(acc auth.AccountKeeper) AccountIndexKeeper {
	return AccountIndexKeeper{AccountKeeper: acc, accounts: make(map[types.AccNumber]sdk.AccAddress)}
}

func (ak AccountIndexKeeper) AddToIndex(acc auth.Account) {
	ak.accounts[types.AccNumber(acc.GetAccountNumber())] = acc.GetAddress()
}

func (ak AccountIndexKeeper) GetAccountAddress(number types.AccNumber) (addr sdk.AccAddress, ok bool) {
	addr, ok = ak.accounts[number]
	return addr, ok
}

func (ak AccountIndexKeeper) GetAccountAddresses(numbers []types.AccNumber) []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0, len(numbers))
	for _, n := range numbers {
		addr, ok := ak.GetAccountAddress(n)
		if ok {
			addresses = append(addresses, addr)
		}
	}

	return addresses
}

func (ak AccountIndexKeeper) RefreshIndex(ctx sdk.Context) {
	ak.AccountKeeper.IterateAccounts(ctx, func(account auth.Account) (stop bool) {
		ak.AddToIndex(account)
		return false
	})
}

func (ak AccountIndexKeeper) GetAccountKeeper() auth.AccountKeeper {
	return ak.AccountKeeper
}

func (ak AccountIndexKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account {
	acc := ak.AccountKeeper.NewAccountWithAddress(ctx, addr)
	ak.AddToIndex(acc)
	return acc
}

func (ak AccountIndexKeeper) NewAccount(ctx sdk.Context, acc auth.Account) auth.Account {
	newAcc := ak.NewAccount(ctx, acc)
	ak.AddToIndex(newAcc)
	return newAcc
}

func (ak AccountIndexKeeper) SetAccount(ctx sdk.Context, acc auth.Account) {
	ak.AddToIndex(acc)
	ak.AccountKeeper.SetAccount(ctx, acc)
}
