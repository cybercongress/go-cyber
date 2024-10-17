package keeper

import (
	"fmt"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"

	ctypes "github.com/cybercongress/go-cyber/v5/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/grid/exported"
	"github.com/cybercongress/go-cyber/v5/x/grid/types"
)

var _ = exported.EnergyKeeper(nil)

type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	proxyKeeper   types.CyberbankKeeper

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	bk types.CyberbankKeeper,
	ak types.AccountKeeper,
	authority string,
) Keeper {
	if addr := ak.GetModuleAddress(types.GridPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.GridPoolName))
	}

	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		proxyKeeper:   bk,
		accountKeeper: ak,
		authority:     authority,
	}
	return keeper
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string { return k.authority }

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&p)
	store.Set(types.ParamsKey, bz)

	return nil
}

func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

func (k Keeper) CreateEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress, name string) error {
	if src.Equals(dst) {
		return types.ErrSelfRoute
	}

	_, found := k.GetRoute(ctx, src, dst)
	if found {
		return types.ErrRouteExist
	}

	routes := k.GetSourceRoutes(ctx, src, k.MaxSourceRoutes(ctx))
	if uint32(len(routes)) == k.MaxSourceRoutes(ctx) {
		return types.ErrMaxRoutes
	}

	acc := k.accountKeeper.GetAccount(ctx, dst)
	if acc == nil {
		acc = k.accountKeeper.NewAccountWithAddress(ctx, dst)
		k.accountKeeper.SetAccount(ctx, acc)
	}

	k.SetRoute(ctx, src, dst, types.NewRoute(src, dst, name, sdk.Coins{}))

	k.proxyKeeper.OnCoinsTransfer(ctx, nil, dst)

	return nil
}

func (k Keeper) EditEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress, value sdk.Coin) error {
	defer telemetry.IncrCounter(1.0, types.ModuleName, "grid routed")

	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	energy := k.GetRoutedToEnergy(ctx, dst)

	if !route.Value.IsValid() {
		coins := sdk.NewCoins(value)
		if err := k.proxyKeeper.SendCoinsFromAccountToModule(ctx, src, types.GridPoolName, coins); err != nil {
			return err
		}
		k.SetRoutedEnergy(ctx, dst, sdk.NewCoins(value))
	} else {
		if value.Amount.GT(route.Value.AmountOf(value.Denom)) {
			diff := sdk.NewCoin(value.Denom, value.Amount.Sub(route.Value.AmountOf(value.Denom)))
			coins := sdk.NewCoins(diff)

			if err := k.proxyKeeper.SendCoinsFromAccountToModule(ctx, src, types.GridPoolName, coins); err != nil {
				return err // should never happen
			}
			k.SetRoutedEnergy(ctx, dst, energy.Sort().Add(diff))
		} else {
			diff := sdk.NewCoin(value.Denom, route.Value.AmountOf(value.Denom).Sub(value.Amount))
			coins := sdk.NewCoins(diff)

			if err := k.proxyKeeper.SendCoinsFromModuleToAccount(ctx, types.GridPoolName, src, coins); err != nil {
				return err // should never happen
			}

			k.SetRoutedEnergy(ctx, dst, energy.Sort().Sub(coins...))
		}
	}

	ampers := route.Value.AmountOf(ctypes.AMPERE)
	volts := route.Value.AmountOf(ctypes.VOLT)

	var newValues sdk.Coins
	if value.Denom == ctypes.VOLT {
		newValues = sdk.NewCoins(value, sdk.NewCoin(ctypes.AMPERE, ampers))
	} else {
		newValues = sdk.NewCoins(sdk.NewCoin(ctypes.VOLT, volts), value)
	}

	k.SetRoute(ctx, src, dst, types.NewRoute(src, dst, route.Name, newValues.Sort()))

	k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)

	return nil
}

func (k Keeper) DeleteEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	if err := k.proxyKeeper.SendCoinsFromModuleToAccount(ctx, types.GridPoolName, src, route.Value); err != nil {
		return err // should never happen
	}

	energy := k.GetRoutedToEnergy(ctx, dst)
	k.SetRoutedEnergy(ctx, dst, energy.Sub(route.Value...))

	k.RemoveRoute(ctx, src, dst)

	k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)

	return nil
}

func (k Keeper) EditEnergyRouteName(ctx sdk.Context, src, dst sdk.AccAddress, name string) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	k.SetRoute(ctx, src, dst, types.NewRoute(src, dst, name, route.Value))

	return nil
}

func (k Keeper) SetRoutes(ctx sdk.Context, routes types.Routes) error {
	// TODO add check that source and destinations addresses already exist
	for _, route := range routes {
		src, err := sdk.AccAddressFromBech32(route.Source)
		if err != nil {
			return err
		}
		dst, err := sdk.AccAddressFromBech32(route.Destination)
		if err != nil {
			return err
		}

		energy := k.GetRoutedToEnergy(ctx, dst)
		if !energy.IsValid() {
			k.SetRoutedEnergy(ctx, dst, route.Value)
		} else {
			k.SetRoutedEnergy(ctx, dst, energy.Add(route.Value...))
		}

		k.SetRoute(ctx, src, dst, types.NewRoute(src, dst, route.Name, route.Value))
		k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)
	}
	return nil
}

func (k Keeper) MaxSourceRoutes(ctx sdk.Context) (res uint32) {
	return k.GetParams(ctx).MaxRoutes
}

func (k Keeper) SetRoute(ctx sdk.Context, src, dst sdk.AccAddress, route types.Route) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetRouteKey(src, dst), types.MustMarshalRoute(k.cdc, route))
}

func (k Keeper) RemoveRoute(ctx sdk.Context, src, dst sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetRouteKey(src, dst))
}

func (k Keeper) SetRoutedEnergy(ctx sdk.Context, dst sdk.AccAddress, amount sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	value := types.NewValue(amount)

	store.Set(types.GetRoutedEnergyByDestinationKey(dst), k.cdc.MustMarshal(&value))
}

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

// Deprecated: enable after debug with empty routes
func (k Keeper) IterateAllRoutes(ctx sdk.Context, cb func(route types.Route) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RouteKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		route, err := types.UnmarshalRoute(k.cdc, iterator.Value())
		if err != nil {
			// TODO fix empty case
			//"grid" : {
			//	"routes" : [{
			//		"destination" : "",
			//		"value" : [],
			//		"source" : "",
			//		"name" : ""
			//	}],
			//	"params" : {
			//		"max_routes" : 16
			//	}
			//}
			break
		}
		if cb(route) {
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
		rdst, _ := sdk.AccAddressFromBech32(route.Destination)
		if rdst.Equals(dst) {
			routes = append(routes, route)
		}
	}

	return routes
}

func (k Keeper) GetSourceRoutes(ctx sdk.Context, src sdk.AccAddress, maxRetrieve uint32) (routes []types.Route) {
	routes = make([]types.Route, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	sourcePrefixKey := types.GetRoutesKey(src)

	iterator := sdk.KVStorePrefixIterator(store, sourcePrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		route := types.MustUnmarshalRoute(k.cdc, iterator.Value())
		routes[i] = route
		i++
	}

	return routes[:i] // trim if the array length < maxRetrieve
}

func (k Keeper) GetRoutedToEnergy(ctx sdk.Context, dst sdk.AccAddress) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRoutedEnergyByDestinationKey(dst))
	if bz == nil {
		return sdk.Coins{}
	}
	amount := types.Value{}
	k.cdc.MustUnmarshal(bz, &amount)

	return amount.Value
}

func (k Keeper) GetRoutedFromEnergy(ctx sdk.Context, src sdk.AccAddress) (amount sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	sourcePrefixKey := types.GetRoutesKey(src)

	iterator := sdk.KVStorePrefixIterator(store, sourcePrefixKey)
	amount = sdk.Coins{}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var route types.Route
		k.cdc.MustUnmarshal(iterator.Value(), &route)
		if amount.IsValid() {
			amount = amount.Add(route.Value...)
		} else {
			amount = route.Value
		}
	}

	return amount
}
