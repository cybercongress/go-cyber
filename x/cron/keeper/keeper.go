package keeper

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/cybercongress/go-cyber/types"

	"github.com/cybercongress/go-cyber/x/cron/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
)

// Keeper of the power store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	wasmKeeper 	  wasm.Keeper
	accountKeeper authkeeper.AccountKeeper
	proxyKeeper   types.BankKeeper
	paramspace    paramstypes.Subspace
	router 		  types.Router
}

func NewKeeper(
	cdc codec.BinaryMarshaler,
	key sdk.StoreKey,
	wk wasm.Keeper,
	bk types.BankKeeper,
	ak authkeeper.AccountKeeper,
	paramSpace paramstypes.Subspace,
	router types.Router,
) Keeper {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		wasmKeeper: wk,
		proxyKeeper:   bk,
		accountKeeper: ak,
		paramspace: paramSpace,
		router: router,
	}
	return keeper
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

func (k Keeper) AddJob(
	ctx sdk.Context, creator, contract sdk.AccAddress,
	trigger types.Trigger, load types.Load,
	label string, cid graphtypes.Cid,
) error {

	if trigger.Block != 0 && ctx.BlockHeight() > int64(trigger.Block) {
		return types.ErrBadTrigger
	}

	// TODO add push with higher fee an pull with lowest

	if uint32(len(k.GetAllJobs(ctx))) >= k.MaxJobs(ctx) {
		return types.ErrExceededMaxJobs
	}

	k.SetJob(ctx, types.NewJob(
		creator.String(), contract.String(),
		trigger, load,
		label, string(cid),
	))
	k.SetJobStats(ctx, contract, creator, label,
		types.NewStats(0,0, 0, uint64(ctx.BlockHeight())),
	)

	return nil
}

func (k Keeper) RemoveJob(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}

	k.DeleteJob(ctx, contract, creator, label)

	return nil
}

func (k Keeper) ChangeJobCID(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	cid graphtypes.Cid,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		*job.Trigger, *job.Load,
		job.Label, string(cid),
	))

	return nil
}

func (k Keeper) ChangeJobLabel(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	labelNew string,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}
	jobStats, _ := k.GetJobStats(ctx, contract, creator, label)

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		*job.Trigger, *job.Load,
		labelNew, job.Cid,
	))

	ct, _ := sdk.AccAddressFromBech32(job.Contract)
	cr, _ := sdk.AccAddressFromBech32(job.Creator)
	k.SetJobStats(ctx, ct, cr, labelNew,
		types.NewStats(
			jobStats.Calls, jobStats.Fees, jobStats.Fees, jobStats.LastBlock,
	))

	k.DeleteJob(ctx, contract, creator, label)
	k.DeleteJobStats(ctx, contract, creator, label)

	return nil
}

func (k Keeper) ChangeJobCallData(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	calldata string,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		*job.Trigger, types.NewLoad(calldata, job.Load.GasPrice),
		job.Label, job.Cid,
	))

	return nil
}

func (k Keeper) ChangeJobGasPrice(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	gasprice sdk.Coin,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}



	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		*job.Trigger, types.NewLoad(job.Load.CallData, gasprice),
		job.Label, job.Cid,
	))

	return nil
}

func (k Keeper) ChangeJobPeriod(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	period uint64,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}

	if job.Trigger.Block > 0 {
		return types.ErrConvertTrigger
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		types.NewTrigger(period, job.Trigger.Block), *job.Load,
		job.Label, job.Cid,
	))

	return nil
}

func (k Keeper) ChangeJobBlock(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	block uint64,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if ctx.BlockHeight() >= int64(block) {
		return types.ErrBadTrigger
	}

	if job.Creator != creator.String() {
		return types.ErrNotAuthorized
	}

	if job.Trigger.Period > 0 {
		return types.ErrConvertTrigger
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		types.NewTrigger(job.Trigger.Period, block), *job.Load,
		job.Label, job.Cid,
	))

	return nil
}


//______________________________________________________________________


func (k Keeper) MaxJobs(ctx sdk.Context) (res uint32) {
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

func (k Keeper) SetJob(ctx sdk.Context, job types.Job) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&job)

	ct, _ := sdk.AccAddressFromBech32(job.Contract)
	cr, _ := sdk.AccAddressFromBech32(job.Creator)
	store.Set(types.GetJobKey(ct, cr, job.Label), b)
}

func (k Keeper) DeleteJob(ctx sdk.Context, contract, creator sdk.AccAddress, label string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetJobKey(contract, creator, label))
}

func (k Keeper) SetJobs(ctx sdk.Context, jobs types.Jobs) error {
	for _, job := range jobs {
		k.SetJob(ctx, types.NewJob(
			job.Creator, job.Contract,
			*job.Trigger, *job.Load,
			job.Label, job.Cid,
		))
	}
	return nil
}

func (k Keeper) SetJobStats(ctx sdk.Context, contract, creator sdk.AccAddress, label string, stats types.JobStats) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&stats)
	store.Set(types.GetJobStatsKey(contract, creator, label), b)
}

func (k Keeper) DeleteJobStats(ctx sdk.Context, contract, creator sdk.AccAddress, label string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetJobStatsKey(contract, creator, label))
}

//______________________________________________________________________

func (k Keeper) GetJob(ctx sdk.Context, contract, creator sdk.AccAddress, label string) (job types.Job, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetJobKey(contract, creator, label)

	value := store.Get(key)
	if value == nil {
		return job, false
	}

	k.cdc.MustUnmarshalBinaryBare(value, &job)

	return job, true
}

func (k Keeper) GetAllJobs(ctx sdk.Context) (jobs types.Jobs) {
	k.IterateAllJobs(ctx, func(job types.Job) bool {
		jobs = append(jobs, job)
		return false
	})

	return jobs
}

func (k Keeper) GetAllJobsStats(ctx sdk.Context) (jobsStats types.JobsStats) {
	k.IterateAllJobsStats(ctx, func(jobStats types.JobStats) bool {
		jobsStats = append(jobsStats, jobStats)
		return false
	})

	return jobsStats
}

func (k Keeper) IterateAllJobsStats(ctx sdk.Context, cb func(jobStats types.JobStats) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.JobStatsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var jobStats types.JobStats
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &jobStats)
		if cb(jobStats) {
			break
		}
	}
}

func (k Keeper) IterateAllJobs(ctx sdk.Context, cb func(job types.Job) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.JobKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var job types.Job
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &job)
		if cb(job) {
			break
		}
	}
}

func (k Keeper) GetJobStats(ctx sdk.Context, contract, creator sdk.AccAddress, label string) (stats types.JobStats, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetJobStatsKey(contract, creator, label)

	value := store.Get(key)
	if value == nil {
		return stats, false
	}

	k.cdc.MustUnmarshalBinaryBare(value, &stats)

	return stats, true
}

//______________________________________________________________________

func (k Keeper) EndBlocker(ctx sdk.Context) {

	jobs := k.GetAllJobs(ctx)
	jobs.Sort()

	maxGas := k.MaxGas(ctx)
	gasBeforeCron := ctx.GasMeter().GasConsumed()
	gasUsedTotal := sdk.Gas(0)

	// Allow recovery from OutOfGas panics so that we don't crash
	// Deletes job which compromised cron execution
	//var currentJob types.Job
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				//k.DeleteJob(ctx, currentJob.Contract, currentJob.Creator, currentJob.Label)
				//k.DeleteJobStats(ctx, currentJob.Contract, currentJob.Creator, currentJob.Label)
				log := fmt.Sprintf(
					"out of gas in location: %v; gasUsed: %d",
					rType.Descriptor, ctx.GasMeter().GasConsumed(),
				)
				k.Logger(ctx).Error(sdkerrors.Wrap(sdkerrors.ErrOutOfGas, log).Error())
			default:
				// Not ErrorOutOfGas, so panic again.
				panic(r)
			}
		}
	}()

	k.Logger(ctx).Debug("Jobs in queue", "size", jobs.Len())

	for i, job := range jobs {
		if (job.Trigger.Period != 0 && ctx.BlockHeight()%int64(job.Trigger.Period) == 0) ||
			(job.Trigger.Period == 0 && ctx.BlockHeight() == int64(job.Trigger.Block)) {
			//currentJob = job
			price := job.Load.GasPrice

			k.Logger(ctx).Debug("Started job", "number", i, "price", price)

			cacheContext, writeFn := ctx.CacheContext()
			gasBefore := ctx.GasMeter().GasConsumed()

			// if not enought remained so not call at all cause we cannot guarantee execution
			// what if infinite gas (gas limit will be 0)
			k.Logger(ctx).Debug("Gas stats",
				"limit", ctx.GasMeter().Limit(),
				"consumed to limit", ctx.GasMeter().GasConsumedToLimit(),
				"consumed", ctx.GasMeter().GasConsumed(),
			)

			remained := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumedToLimit()
			if remained < uint64(k.MaxGas(ctx)) {
				k.Logger(ctx).Debug("Job break, not enough gas", "number", i)
				break
			}

			// because we need events applied
			msg := wasm.MsgExecuteContract{
				job.Creator, // TODO job Contract self call
				job.Contract,
				[]byte(job.Load.CallData),
				sdk.Coins{},
			}
			result, errExecute := k.runJob(cacheContext, &msg)
			k.Logger(ctx).Debug("Job executed", "result", result)

			//if i > 0 {
			//	panic(sdk.ErrorOutOfGas{"Wasmer function execution"})
			//}

			gasUsed := ctx.GasMeter().GasConsumed() - gasBefore
			if gasUsedTotal + gasUsed > uint64(maxGas) {
				break
			} else {
				gasUsedTotal += gasUsed
			}

			contract, _ := sdk.AccAddressFromBech32(job.Contract)
			creator, _ := sdk.AccAddressFromBech32(job.Creator)
			js, _ := k.GetJobStats(ctx, contract, creator, job.Label)
			amtGasFee := price.Amount.Int64() * int64(gasUsed)
			amtTTLFee := (ctx.BlockHeight() - int64(js.LastBlock))*int64(k.FeeTTL(ctx))
			amtTotalFee := amtGasFee + amtTTLFee

			k.Logger(ctx).Debug("Gas",
				"used", gasUsed,
				"gas fee", amtGasFee,
				"ttl fee", amtTTLFee,
				"total fee", amtTotalFee,
			)

			fee := sdk.NewCoin(ctypes.CYB, sdk.NewInt(amtTotalFee))

			errSend := k.proxyKeeper.SendCoins(
				ctx, contract, k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName), sdk.NewCoins(fee))
			if errSend != nil {
				k.DeleteJob(ctx, contract, creator, job.Label)
				k.DeleteJobStats(ctx, contract, creator, job.Label)

				k.Logger(ctx).Debug("Not enough contract balance, state not applied, job killed", "number", i)
				continue
			}

			if errExecute != nil {
				k.Logger(ctx).Debug("Job failed, state not applied", "number", i)
				k.Logger(ctx).Debug("Failed with error: ", errExecute.Error())
			} else {
				writeFn()
				k.Logger(ctx).Debug("Job finished, state applied", "number", i)
			}

			k.SetJobStats(ctx, contract, creator, job.Label,
				types.NewStats(
					js.Calls+1, js.Fees+uint64(amtTotalFee),
					js.Gas+gasUsed, uint64(ctx.BlockHeight())),
			)

			if ctx.BlockHeight() == int64(job.Trigger.Block) {
				k.DeleteJob(ctx, contract, creator, job.Label)
				k.DeleteJobStats(ctx, contract, creator, job.Label)

				k.Logger(ctx).Debug("Job executed at given block, cleaned", "number", i)
			}
		}
	}

	gasAfterCron := ctx.GasMeter().GasConsumed()
	k.Logger(ctx).Debug("Total used", "gas", gasAfterCron-gasBeforeCron)
}

func (k Keeper) runJob(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	handler := k.router.Route(ctx, msg.Route())
	if handler == nil {
		return nil, types.ErrInvalidRoute
	}

	return handler(ctx, msg)
}