package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/cybercongress/go-cyber/types"

	"github.com/cybercongress/go-cyber/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

// Keeper of the power store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	wasmKeeper    wasm.Keeper
	accountKeeper types.AccountKeeper
	proxyKeeper   types.BankKeeper
	paramspace    paramstypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	paramSpace paramstypes.Subspace,
) *Keeper {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:      key,
		cdc:           cdc,
		proxyKeeper:   bk,
		accountKeeper: ak,
		paramspace:    paramSpace,
	}
}

func (k *Keeper) SetWasmKeeper(ws wasm.Keeper) {
	k.wasmKeeper = ws
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}

func (k Keeper) SaveThought(
	ctx sdk.Context, program sdk.AccAddress,
	trigger types.Trigger, load types.Load,
	name string, particle graphtypes.Cid,
) error {

	if trigger.Block != 0 && ctx.BlockHeight() > int64(trigger.Block) {
		return types.ErrBadTrigger
	}

	// if there are full slots but new one with higher fee than delete thought with
	// the smallest one fee and add new one with higher fee
	thoughts := k.GetAllThoughts(ctx)
	thoughts.Sort()
	if uint32(len(thoughts)) == k.MaxThougths(ctx) {
		if thoughts[len(thoughts)-1].Load.GasPrice.IsLT(load.GasPrice) {
			pr, _ := sdk.AccAddressFromBech32(thoughts[len(thoughts)-1].Program)
			k.DeleteThought(ctx, pr, thoughts[len(thoughts)-1].Name)
			k.DeleteThoughtStats(ctx, pr, name)
		} else {
			return types.ErrExceededMaxThoughts
		}
	}

	k.SetThought(ctx, types.NewThought(
		program.String(),
		trigger, load,
		name, string(particle),
	))
	// set last_block to current height to start count future ttl fee
	k.SetThoughtStats(ctx, program, name,
		types.NewStats(
			program.String(), name,
			0, 0, 0, uint64(ctx.BlockHeight()),
		),
	)

	return nil
}

func (k Keeper) RemoveThoughtFull(
	ctx sdk.Context, program sdk.AccAddress, name string,
) error {
	_, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}

	k.DeleteThought(ctx, program, name)
	k.DeleteThoughtStats(ctx, program, name)

	return nil
}

func (k Keeper) UpdateThoughtParticle(
	ctx sdk.Context, program sdk.AccAddress, name string,
	particle graphtypes.Cid,
) error {
	thought, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}

	k.SetThought(ctx, types.NewThought(
		thought.Program,
		thought.Trigger, thought.Load,
		thought.Name, string(particle),
	))

	return nil
}

func (k Keeper) UpdateThoughtName(
	ctx sdk.Context, program sdk.AccAddress, name string,
	nameNew string,
) error {
	thought, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}
	thoughtStats, _ := k.GetThoughtStats(ctx, program, name)

	if thought.Name == nameNew {
		return types.ErrBadName
	}

	k.DeleteThought(ctx, program, name)
	k.DeleteThoughtStats(ctx, program, name)

	k.SetThought(ctx, types.NewThought(
		thought.Program,
		thought.Trigger, thought.Load,
		nameNew, thought.Particle,
	))

	k.SetThoughtStats(ctx, program, nameNew,
		types.NewStats(
			program.String(), nameNew,
			thoughtStats.Calls, thoughtStats.Fees, thoughtStats.Fees, thoughtStats.LastBlock,
		))

	return nil
}

func (k Keeper) UpdateThoughtCallData(
	ctx sdk.Context, program sdk.AccAddress, name string,
	calldata string,
) error {
	thought, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}

	k.SetThought(ctx, types.NewThought(
		thought.Program,
		thought.Trigger, types.NewLoad(calldata, thought.Load.GasPrice),
		thought.Name, thought.Particle,
	))

	return nil
}

func (k Keeper) UpdateThoughtGasPrice(
	ctx sdk.Context, program sdk.AccAddress, name string,
	gasprice sdk.Coin,
) error {
	thought, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}

	k.SetThought(ctx, types.NewThought(
		thought.Program,
		thought.Trigger, types.NewLoad(thought.Load.Input, gasprice),
		thought.Name, thought.Particle,
	))

	return nil
}

func (k Keeper) UpdateThoughtPeriod(
	ctx sdk.Context, program sdk.AccAddress, name string,
	period uint64,
) error {
	thought, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}

	if thought.Trigger.Block > 0 {
		return types.ErrConvertTrigger
	}

	k.SetThought(ctx, types.NewThought(
		thought.Program,
		types.NewTrigger(period, thought.Trigger.Block), thought.Load,
		thought.Name, thought.Particle,
	))

	return nil
}

func (k Keeper) UpdateThoughtBlock(
	ctx sdk.Context, program sdk.AccAddress, name string,
	block uint64,
) error {
	thought, found := k.GetThought(ctx, program, name)
	if !found {
		return types.ErrThoughtNotExist
	}

	if ctx.BlockHeight() >= int64(block) {
		return types.ErrBadTrigger
	}

	if thought.Trigger.Period > 0 {
		return types.ErrConvertTrigger
	}

	k.SetThought(ctx, types.NewThought(
		thought.Program,
		types.NewTrigger(thought.Trigger.Period, block), thought.Load,
		thought.Name, thought.Particle,
	))

	return nil
}

//______________________________________________________________________

func (k Keeper) MaxThougths(ctx sdk.Context) (res uint32) {
	k.paramspace.Get(ctx, types.KeyMaxSlots, &res)
	return
}

func (k Keeper) MaxGas(ctx sdk.Context) (res uint32) {
	k.paramspace.Get(ctx, types.KeyMaxGas, &res)
	return
}

func (k Keeper) FeeTTL(ctx sdk.Context) (res uint32) {
	k.paramspace.Get(ctx, types.KeyFeeTTL, &res)
	return
}

//______________________________________________________________________

func (k Keeper) SetThought(ctx sdk.Context, thought types.Thought) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&thought)

	program, _ := sdk.AccAddressFromBech32(thought.Program)
	store.Set(types.GetThoughtKey(program, thought.Name), b)
}

func (k Keeper) DeleteThought(ctx sdk.Context, program sdk.AccAddress, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetThoughtKey(program, name))
}

func (k Keeper) SetThoughts(ctx sdk.Context, thoughts types.Thoughts) error {
	for _, thought := range thoughts {
		k.SetThought(ctx, types.NewThought(
			thought.Program,
			thought.Trigger, thought.Load,
			thought.Name, thought.Particle,
		))
	}
	return nil
}

func (k Keeper) SetThoughtStats(ctx sdk.Context, program sdk.AccAddress, name string, stats types.ThoughtStats) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&stats)
	store.Set(types.GetThoughtStatsKey(program, name), b)
}

func (k Keeper) DeleteThoughtStats(ctx sdk.Context, program sdk.AccAddress, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetThoughtStatsKey(program, name))
}

//______________________________________________________________________

func (k Keeper) GetThought(ctx sdk.Context, program sdk.AccAddress, name string) (thought types.Thought, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetThoughtKey(program, name)

	value := store.Get(key)
	if value == nil {
		return thought, false
	}

	k.cdc.MustUnmarshal(value, &thought)

	return thought, true
}

func (k Keeper) GetAllThoughts(ctx sdk.Context) (thoughts types.Thoughts) {
	k.IterateAllThoughts(ctx, func(thought types.Thought) bool {
		thoughts = append(thoughts, thought)
		return false
	})

	return thoughts
}

func (k Keeper) GetAllThoughtsStats(ctx sdk.Context) (thoughtsStats types.ThoughtsStats) {
	k.IterateAllThoughtsStats(ctx, func(thoughtStats types.ThoughtStats) bool {
		thoughtsStats = append(thoughtsStats, thoughtStats)
		return false
	})

	return thoughtsStats
}

func (k Keeper) IterateAllThoughtsStats(ctx sdk.Context, cb func(thoughtStats types.ThoughtStats) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ThoughtStatsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var thoughtStats types.ThoughtStats
		k.cdc.MustUnmarshal(iterator.Value(), &thoughtStats)
		if cb(thoughtStats) {
			break
		}
	}
}

func (k Keeper) IterateAllThoughts(ctx sdk.Context, cb func(thought types.Thought) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ThoughtKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var thought types.Thought
		k.cdc.MustUnmarshal(iterator.Value(), &thought)
		if cb(thought) {
			break
		}
	}
}

func (k Keeper) GetThoughtStats(ctx sdk.Context, program sdk.AccAddress, name string) (stats types.ThoughtStats, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetThoughtStatsKey(program, name)

	value := store.Get(key)
	if value == nil {
		return stats, false
	}

	k.cdc.MustUnmarshal(value, &stats)

	return stats, true
}

func (k Keeper) GetLowestFee(ctx sdk.Context) sdk.Coin {
	thoughts := k.GetAllThoughts(ctx)
	if len(thoughts) == 0 {
		return ctypes.NewCybCoin(0)
	} else {
		thoughts.Sort()
		return thoughts[len(thoughts)-1].Load.GasPrice
	}
}

func (k Keeper) ExecuteThoughtsQueue(ctx sdk.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				errorMsg := fmt.Sprintf(
					"out of gas in location: %v; gasUsed: %d",
					rType.Descriptor, ctx.GasMeter().GasConsumed(),
				)
				k.Logger(ctx).Error(sdkerrors.Wrap(sdkerrors.ErrOutOfGas, errorMsg).Error())
			default:
				// Not ErrorOutOfGas, so panic again.
				panic(r)
			}
		}
	}()

	thoughts := k.GetAllThoughts(ctx)
	thoughts.Sort()

	maxGas := k.MaxGas(ctx)
	gasBeforeDmn := ctx.GasMeter().GasConsumed()
	gasUsedTotal := sdk.Gas(0)

	feeTTL := k.FeeTTL(ctx)
	maxThoughts := k.MaxThougths(ctx)
	maxGasPerThought := maxGas / maxThoughts

	if thoughts.Len() > 0 {
		k.Logger(ctx).Info("Thoughts in queue", "size", thoughts.Len())
	}

	var thoughtsTriggered uint32
	for i, thought := range thoughts {
		if (thought.Trigger.Period != 0 && ctx.BlockHeight()%int64(thought.Trigger.Period) == 0) ||
			(thought.Trigger.Period == 0 && ctx.BlockHeight() == int64(thought.Trigger.Block)) {
			price := thought.Load.GasPrice

			k.Logger(ctx).Info("Started thought", "number", i, "gas price", price)
			thoughtsTriggered = thoughtsTriggered + 1

			cacheContext, writeFn := ctx.CacheContext()
			cacheContext = cacheContext.WithGasMeter(sdk.NewGasMeter(sdk.Gas(maxGasPerThought)))

			k.Logger(ctx).Info("Context gas stats before thought execution",
				"consumed", ctx.GasMeter().GasConsumed(),
			)

			remained := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumedToLimit()
			if remained < uint64(maxGasPerThought) {
				k.Logger(ctx).Info("Thought break, not enough gas", "thought #", i)
				break
			}

			program, _ := sdk.AccAddressFromBech32(thought.Program)
			_, errExecute := k.executeThoughtWithSudo(
				cacheContext, program, thought.Load.Input,
			)

			gasUsedByThought := cacheContext.GasMeter().GasConsumed()
			ctx.GasMeter().ConsumeGas(gasUsedByThought, "thought execution")
			if gasUsedTotal+gasUsedByThought > uint64(maxGas) {
				break
			} else {
				gasUsedTotal += gasUsedByThought
			}

			js, _ := k.GetThoughtStats(ctx, program, thought.Name)
			// TODO move to more advanced fee system, 10X fee reducer applied (min 0.1 per gas)
			amtGasFee := price.Amount.Int64() * int64(gasUsedByThought) / 10
			amtTTLFee := (ctx.BlockHeight() - int64(js.LastBlock)) * int64(feeTTL)
			amtTotalFee := amtGasFee + amtTTLFee

			k.Logger(ctx).Info("Gas thought execution stats",
				"used", gasUsedByThought,
				"gas fee", amtGasFee,
				"ttl fee", amtTTLFee,
				"total fee", amtTotalFee,
			)

			fee := sdk.NewCoin(ctypes.CYB, sdk.NewInt(amtTotalFee))

			errSend := k.proxyKeeper.SendCoins(
				ctx, program, k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName), sdk.NewCoins(fee))
			if errSend != nil {
				k.DeleteThought(ctx, program, thought.Name)
				k.DeleteThoughtStats(ctx, program, thought.Name)

				k.Logger(ctx).Info("Not enough program balance, state not applied, thought forgotten", "Thought #", i)
				continue
			}

			if errExecute != nil {
				k.Logger(ctx).Info("Thought failed, state not applied", "Thought #", i)
				k.Logger(ctx).Info("Failed with error: ", "Error", errExecute.Error())
			} else {
				writeFn() // apply cached context
				k.Logger(ctx).Info("Thought finished, state applied", "Thought #", i)
			}

			k.SetThoughtStats(ctx, program, thought.Name,
				types.NewStats(
					program.String(), thought.Name,
					js.Calls+1, js.Fees+uint64(amtTotalFee),
					js.Gas+gasUsedByThought, uint64(ctx.BlockHeight())),
			)

			if ctx.BlockHeight() == int64(thought.Trigger.Block) {
				k.DeleteThought(ctx, program, thought.Name)
				k.DeleteThoughtStats(ctx, program, thought.Name)

				k.Logger(ctx).Info("Thought executed at given block, deleted from queue", "Thought #", i)
			}
		}
	}

	gasAfterDmn := ctx.GasMeter().GasConsumed()
	if thoughtsTriggered > 0 {
		k.Logger(ctx).Info("Total dmn gas used", "Gas used", gasAfterDmn-gasBeforeDmn)
	}
}

func (k Keeper) executeThoughtWithSudo(ctx sdk.Context, program sdk.AccAddress, msg string) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				errorMsg := fmt.Sprintf(
					"out of gas in location: %v; gasUsed: %d",
					rType.Descriptor, ctx.GasMeter().GasConsumed(),
				)
				k.Logger(ctx).Error(sdkerrors.Wrap(sdkerrors.ErrOutOfGas, errorMsg).Error())
			default:
				// Not ErrorOutOfGas, so panic again.
				panic(r)
			}
		}
		telemetry.IncrCounter(1.0, types.ModuleName, "thought")
	}()

	callData, berr := base64.StdEncoding.DecodeString(msg)
	if berr != nil {
		// TODO remove hex later as deprecated
		_, herr := hex.DecodeString(msg)
		if herr != nil {
			return nil, types.ErrBadCallData
		}
	}
	return k.wasmKeeper.Sudo(ctx, program, callData)
}
