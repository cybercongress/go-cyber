package keeper

import (
	"fmt"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/cybercongress/go-cyber/x/energy/exported"
	"github.com/cybercongress/go-cyber/x/energy/types"
)

var _ = exported.EnergyKeeper(nil)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	accountKeeper authkeeper.AccountKeeper
	proxyKeeper   types.BankKeeper
	paramSpace    paramstypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryMarshaler,
	key sdk.StoreKey,
	bk types.BankKeeper,
	ak authkeeper.AccountKeeper,
	paramSpace paramstypes.Subspace,
) Keeper {
	if addr := ak.GetModuleAddress(types.EnergyPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.EnergyPoolName))
	}

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		proxyKeeper:   bk,
		accountKeeper: ak,
		paramSpace:    paramSpace,
	}
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
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
	if uint32(len(routes)) == k.MaxSourceRoutes(ctx) {
		return  types.ErrMaxRoutes
	}

	acc := k.accountKeeper.GetAccount(ctx, dst)
	if acc == nil {
		acc = k.accountKeeper.NewAccountWithAddress(ctx, dst)
		k.accountKeeper.SetAccount(ctx, acc)
	}

	k.SetRoute(ctx, types.NewRoute(src, dst, alias, sdk.Coins{}))

	return nil
}

func (k Keeper) EditEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress, value sdk.Coin) error {
	if value.Denom != ctypes.VOLT && value.Denom != ctypes.AMPER {
		return types.ErrWrongDenom
	}

	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	energy := k.GetRoutedToEnergy(ctx, dst)
	fmt.Println("ENERGY: ", energy)

	if !route.Value.IsValid() {
		coins := sdk.NewCoins(value)
		if err := k.proxyKeeper.SendCoinsFromAccountToModule(ctx, src, types.EnergyPoolName, coins); err != nil {
			return err
		}
		k.SetRoutedEnergy(ctx, dst, sdk.NewCoins(value))
	} else {
		if value.Amount.GT(route.Value.AmountOf(value.Denom)) { //TODO test here
			diff := sdk.NewCoin(value.Denom, value.Amount.Sub(route.Value.AmountOf(value.Denom)))
			coins := sdk.NewCoins(diff)

			fmt.Println("Energy +++ diff (sent A->M)", diff)
			fmt.Println("Energy +++ energy+diff (set routed)", energy.Add(diff))

			if err := k.proxyKeeper.SendCoinsFromAccountToModule(ctx, src, types.EnergyPoolName, coins); err != nil {
				return err
			}
			k.SetRoutedEnergy(ctx, dst, energy.Sort().Add(diff))
		} else {
			diff := sdk.NewCoin(value.Denom, route.Value.AmountOf(value.Denom).Sub(value.Amount))
			coins := sdk.NewCoins(diff)

			fmt.Println("Energy --- diff (sent M->A)", diff)
			fmt.Println("Energy --- energy-diff (set routed)", energy.Sub(coins))

			if err := k.proxyKeeper.SendCoinsFromModuleToAccount(ctx, types.EnergyPoolName, src, coins); err != nil {
				return err
			}

			k.SetRoutedEnergy(ctx, dst, energy.Sub(coins))
		}
	}

	ampers := route.Value.AmountOf(ctypes.AMPER)
	volts  := route.Value.AmountOf(ctypes.VOLT)
	newValues := sdk.Coins{}
	if value.Denom == ctypes.VOLT {
		newValues = sdk.NewCoins(value, sdk.NewCoin(ctypes.AMPER, ampers))
	} else {
		newValues = sdk.NewCoins(sdk.NewCoin(ctypes.VOLT, volts), value)
	}

	fmt.Println("Energy *** newValues", newValues)

	k.SetRoute(ctx, types.NewRoute(src, dst, route.Alias, newValues))

	// TODO migrate from on coins callback to on coins in proxy methods callback
	// TODO but SendCoinsFromModuleToAccount and others will pass modules, not both src,dst
	k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)

	return nil
}

func (k Keeper) DeleteEnergyRoute(ctx sdk.Context, src, dst sdk.AccAddress) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	if err := k.proxyKeeper.SendCoinsFromModuleToAccount(ctx, types.EnergyPoolName, src, route.Value); err != nil {
		return err
	}

	energy := k.GetRoutedToEnergy(ctx, dst)
	k.SetRoutedEnergy(ctx, dst, energy.Sub(route.Value))

	k.RemoveRoute(ctx, route)

	k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)

	return nil
}

func (k Keeper) EditEnergyRouteAlias(ctx sdk.Context, src, dst sdk.AccAddress, alias string) error {
	route, found := k.GetRoute(ctx, src, dst)
	if !found {
		return types.ErrRouteNotExist
	}

	k.SetRoute(ctx, types.NewRoute(src, dst, alias, route.Value))

	return nil
}

func (k Keeper) SetRoutes(ctx sdk.Context, routes types.Routes) error {
	for _, route := range routes {
		src, _ := sdk.AccAddressFromBech32(route.Source)
		dst, _ := sdk.AccAddressFromBech32(route.Destination)

		if err := k.proxyKeeper.SendCoinsFromAccountToModule(ctx, src, types.EnergyPoolName, route.Value); err != nil {
			return err
		}

		energy := k.GetRoutedToEnergy(ctx, dst)
		if !energy.IsValid() {
			k.SetRoutedEnergy(ctx, dst, route.Value)
		} else {
			k.SetRoutedEnergy(ctx, dst, energy.Add(route.Value...))
		}

		k.SetRoute(ctx, types.NewRoute(src, dst, route.Alias, route.Value))
		k.proxyKeeper.OnCoinsTransfer(ctx, src, dst)
	}
	return nil
}

//______________________________________________________________________

func (k Keeper) MaxSourceRoutes(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyMaxRoutes, &res)
	return
}

//______________________________________________________________________

func (k Keeper) SetRoute(ctx sdk.Context, route types.Route) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&route)

	src, _ := sdk.AccAddressFromBech32(route.Source)
	dst, _ := sdk.AccAddressFromBech32(route.Destination)

	store.Set(types.GetRouteKey(src, dst), b)
}

func (k Keeper) RemoveRoute(ctx sdk.Context, route types.Route) {
	store := ctx.KVStore(k.storeKey)

	src, _ := sdk.AccAddressFromBech32(route.Source)
	dst, _ := sdk.AccAddressFromBech32(route.Destination)

	store.Delete(types.GetRouteKey(src, dst))
}

func (k Keeper) SetRoutedEnergy(ctx sdk.Context, dst sdk.AccAddress, amount sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	value := types.NewValue(amount)
	store.Set(types.GetRoutedEnergyByDestinationKey(dst), k.cdc.MustMarshalBinaryBare(&value))
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

	k.cdc.MustUnmarshalBinaryBare(value, &route)

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
		var route types.Route
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &route)
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
		var route types.Route
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &route)
		rdst, _ := sdk.AccAddressFromBech32(route.Destination)
		if rdst.Equals(dst) {
			routes = append(routes, route)
		}
	}

	return routes
}

func (k Keeper) GetSourceRoutes(ctx sdk.Context, src sdk.AccAddress,
	maxRetrieve uint32) (routes []types.Route) {
	routes = make([]types.Route, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	sourcePrefixKey := types.GetRoutesKey(src)

	iterator := sdk.KVStorePrefixIterator(store, sourcePrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		var route types.Route
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &route)
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
	k.cdc.MustUnmarshalBinaryBare(bz, &amount)

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
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &route)
		if amount.IsValid() {
			amount = amount.Add(route.GetValue()...)
		} else {
			amount = route.Value
		}
	}

	return amount
}