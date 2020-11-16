package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/energy/exported"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

var _ = exported.EnergyKeeper(nil)

// Keeper of the power store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	supplyKeeper  types.SupplyKeeper
	accountKeeper types.AccountKeeper
	proxyKeeper   types.BankKeeper
	paramspace    types.ParamSubspace
}

func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey,
	sk types.SupplyKeeper, bk types.BankKeeper,
	ak types.AccountKeeper, paramspace types.ParamSubspace,
) Keeper {
	if addr := sk.GetModuleAddress(types.EnergyPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.EnergyPoolName))
	}

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		supplyKeeper: sk,
		proxyKeeper:   bk,
		accountKeeper: ak,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) CreateEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress, alias string) error {
	if src.Equals(dst) {
		return types.ErrSelfRoute
	}

	_, found := k.GetRoute(ctx, src, dst)
	if found {
		return types.ErrRouteExist
	}

	_, found = k.GetRoute(ctx, dst, src)
	if found {
		return types.ErrReverseRoute
	}

	routes := k.GetSourceRoutes(ctx, src, k.MaxSourceRoutes(ctx))
	if uint16(len(routes)) > k.MaxSourceRoutes(ctx) {
		return  types.ErrMaxRoutes
	}

	acc := k.accountKeeper.GetAccount(ctx, dst)
	if acc == nil {
		acc = k.accountKeeper.NewAccountWithAddress(ctx, dst)
		k.accountKeeper.SetAccount(ctx, acc)
	}

	k.SetRoute(ctx, types.NewRoute(src, dst, alias, sdk.ZeroInt()))

	return nil
}

func (k Keeper) EditEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress, value sdk.Coin) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	if value.Denom != k.EnergyDenom(ctx) {
		return types.ErrWrongDenom
	}

	energy := k.GetRoutedToEnergy(ctx, dst)

	if value.Amount.GT(route.Amount) {
		diff := value.Amount.Sub(route.Amount)
		coins := sdk.NewCoins(sdk.NewCoin(k.EnergyDenom(ctx), diff))
		if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, src, types.EnergyPoolName, coins); err != nil {
			return err
		}
		k.SetRoutedEnergy(ctx, dst, energy.Add(diff))
	} else {
		diff := route.Amount.Sub(value.Amount)
		coins := sdk.NewCoins(sdk.NewCoin(k.EnergyDenom(ctx), diff))
		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.EnergyPoolName, src, coins); err != nil {
			return err
		}
		k.SetRoutedEnergy(ctx, dst, energy.Sub(diff))
	}
	k.SetRoute(ctx, types.NewRoute(route.Source, route.Destination, route.Alias, value.Amount))

	k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)

	return nil
}

func (k Keeper) DeleteEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.EnergyDenom(ctx), route.Amount))
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.EnergyPoolName, route.Source, coins); err != nil {
		return err
	}

	energy := k.GetRoutedToEnergy(ctx, dst)
	k.SetRoutedEnergy(ctx, dst, energy.Sub(route.Amount))

	k.RemoveRoute(ctx, route)

	k.proxyKeeper.OnCoinsTransfer(ctx, dst, nil)

	return nil
}

func (k Keeper) EditEnergyRouteAlias(ctx sdk.Context, src, dst sdk.AccAddress, alias string) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	k.SetRoute(ctx, types.NewRoute(route.Source, route.Destination, alias, route.Amount))

	return nil
}

func (k Keeper) SetRoutes(ctx sdk.Context, routes types.Routes) error {
	for _, route := range routes {
		energy := k.GetRoutedToEnergy(ctx, route.Destination)
		coins := sdk.NewCoins(sdk.NewCoin(k.EnergyDenom(ctx), route.Amount))
		if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, route.Source, types.EnergyPoolName, coins); err != nil {
			return err
		}
		k.SetRoutedEnergy(ctx, route.Destination, energy.Add(route.Amount))
		k.SetRoute(ctx, types.NewRoute(route.Source, route.Destination, route.Alias, route.Amount))
		k.proxyKeeper.OnCoinsTransfer(ctx, route.Destination, nil)
	}
	// TODO imported routes amount for src should not exceed max source routes param (validate genesis)
	return nil
}

//______________________________________________________________________

func (k Keeper) EnergyDenom(ctx sdk.Context) (res string) {
	k.paramspace.Get(ctx, types.KeyEnergyDenom, &res)
	return
}

func (k Keeper) MaxSourceRoutes(ctx sdk.Context) (res uint16) {
	k.paramspace.Get(ctx, types.KeyMaxRoutes, &res)
	return
}

//______________________________________________________________________

func (k Keeper) SetRoute(ctx sdk.Context, route types.Route) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalRoute(k.cdc, route)
	store.Set(types.GetRouteKey(route.Source, route.Destination), b)
}

func (k Keeper) RemoveRoute(ctx sdk.Context, route types.Route) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetRouteKey(route.Source, route.Destination))
}

func (k Keeper) SetRoutedEnergy(ctx sdk.Context, dst sdk.AccAddress, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRoutedEnergyByDestinationKey(dst), k.cdc.MustMarshalBinaryBare(amount))
}

//func (k Keeper) RemoveRoutedPower(ctx sdk.Context, delegate sdk.AccAddress) {
//	store := ctx.KVStore(k.storeKey)
//	store.Delete(types.GetRoutedEnergyByDestinationKey(delegate))
//}

//______________________________________________________________________

func (k Keeper) GetRoute(ctx sdk.Context, src, dst sdk.AccAddress) (route types.Route, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRouteKey(src, dst)

	value := store.Get(key)
	if value == nil {
		return route, false
	}

	route = types.MustUnmarshalRoute(k.cdc, value)

	return route, true
}

func (k Keeper) GetAllRoutes(ctx sdk.Context) (routes []types.Route) {
	k.IterateAllRoutes(ctx, func(route types.Route) bool {
		routes = append(routes, route)
		return false
	})

	return routes
}

func (k Keeper) IterateAllRoutes(ctx sdk.Context, cb func(route types.Route) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RouteKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalRoute(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

func (k Keeper) GetDestinationRoutes(ctx sdk.Context, dst sdk.AccAddress) (routes []types.Route) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RouteKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		route := types.MustUnmarshalRoute(k.cdc, iterator.Value())
		if route.Destination.Equals(dst) {
			routes = append(routes, route)
		}
	}

	return routes
}

func (k Keeper) GetSourceRoutes(ctx sdk.Context, src sdk.AccAddress,
	maxRetrieve uint16) (routes []types.Route) {
	routes = make([]types.Route, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	sourcePrefixKey := types.GetRoutesKey(src)

	iterator := sdk.KVStorePrefixIterator(store, sourcePrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalRoute(k.cdc, iterator.Value())
		routes[i] = delegation
		i++
	}

	return routes[:i] // trim if the array length < maxRetrieve
}

func (k Keeper) GetRoutedToEnergy(ctx sdk.Context, dst sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRoutedEnergyByDestinationKey(dst))
	if bz == nil {
		return sdk.NewInt(0)
	}
	amount := sdk.Int{}
	k.cdc.MustUnmarshalBinaryBare(bz, &amount)

	return amount
}

func (k Keeper) GetRoutedFromEnergy(ctx sdk.Context, src sdk.AccAddress) (amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	sourcePrefixKey := types.GetRoutesKey(src)

	iterator := sdk.KVStorePrefixIterator(store, sourcePrefixKey)
	amount = sdk.ZeroInt()

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		route := types.MustUnmarshalRoute(k.cdc, iterator.Value())
		amount = amount.Add(route.GetAmount())
	}

	return amount
}