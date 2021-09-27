package keeper

import (
	"fmt"

	"encoding/hex"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	cdc           codec.BinaryCodec
	wasmKeeper 	  wasm.Keeper
	accountKeeper types.AccountKeeper
	proxyKeeper   types.BankKeeper
	paramspace    paramstypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	bk  types.BankKeeper,
	ak  types.AccountKeeper,
	paramSpace paramstypes.Subspace,
) *Keeper {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:   key,
		cdc:        cdc,
		proxyKeeper:   bk,
		accountKeeper: ak,
		paramspace: paramSpace,
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

func (k Keeper) SaveJob(
	ctx sdk.Context, program sdk.AccAddress,
	trigger types.Trigger, load types.Load,
	label string, cid graphtypes.Cid,
) error {

	if trigger.Block != 0 && ctx.BlockHeight() > int64(trigger.Block) {
		return types.ErrBadTrigger
	}

	// if there are full slots but new one with higher fee than delete job with
	// the smallest one fee and add new one with higher fee
	jobs := k.GetAllJobs(ctx)
	jobs.Sort()
	if uint32(len(jobs)) == k.MaxJobs(ctx) {
		if jobs[len(jobs)-1].Load.GasPrice.IsLT(load.GasPrice) {
			pr, _ := sdk.AccAddressFromBech32(jobs[len(jobs)-1].Program)
			k.DeleteJob(ctx, pr, jobs[len(jobs)-1].Label)
			k.DeleteJobStats(ctx, pr, label)
		} else {
			return types.ErrExceededMaxJobs
		}
	}

	k.SetJob(ctx, types.NewJob(
		program.String(),
		trigger, load,
		label, string(cid),
	))
	// set last_block to current height to start count future ttl fee
	k.SetJobStats(ctx, program, label,
		types.NewStats(
			program.String(), label,
			0,0, 0, uint64(ctx.BlockHeight()),
		),
	)

	return nil
}

func (k Keeper) RemoveJobFull(
	ctx sdk.Context, program sdk.AccAddress, label string,
) error {
	_, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}

	k.DeleteJob(ctx, program, label)
	k.DeleteJobStats(ctx, program, label)

	return nil
}

func (k Keeper) UpdateJobCID(
	ctx sdk.Context, program sdk.AccAddress, label string,
	cid graphtypes.Cid,
) error {
	job, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}

	k.SetJob(ctx, types.NewJob(
		job.Program,
		job.Trigger, job.Load,
		job.Label, string(cid),
	))

	return nil
}

func (k Keeper) UpdateJobLabel(
	ctx sdk.Context, program sdk.AccAddress, label string,
	labelNew string,
) error {
	job, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}
	jobStats, _ := k.GetJobStats(ctx, program, label)

	if job.Label == labelNew {
		return types.ErrBadLabel
	}

	k.DeleteJob(ctx, program, label)
	k.DeleteJobStats(ctx, program, label)

	k.SetJob(ctx, types.NewJob(
		job.Program,
		job.Trigger, job.Load,
		labelNew, job.Cid,
	))

	k.SetJobStats(ctx, program, labelNew,
		types.NewStats(
			program.String(), labelNew,
			jobStats.Calls, jobStats.Fees, jobStats.Fees, jobStats.LastBlock,
	))

	return nil
}

func (k Keeper) UpdateJobCallData(
	ctx sdk.Context, program sdk.AccAddress, label string,
	calldata string,
) error {
	job, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}

	k.SetJob(ctx, types.NewJob(
		job.Program,
		job.Trigger, types.NewLoad(calldata, job.Load.GasPrice),
		job.Label, job.Cid,
	))

	return nil
}

func (k Keeper) UpdateJobGasPrice(
	ctx sdk.Context, program sdk.AccAddress, label string,
	gasprice sdk.Coin,
) error {
	job, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}

	k.SetJob(ctx, types.NewJob(
		job.Program,
		job.Trigger, types.NewLoad(job.Load.CallData, gasprice),
		job.Label, job.Cid,
	))

	return nil
}

func (k Keeper) UpdateJobPeriod(
	ctx sdk.Context, program sdk.AccAddress, label string,
	period uint64,
) error {
	job, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}

	if job.Trigger.Block > 0 {
		return types.ErrConvertTrigger
	}

	k.SetJob(ctx, types.NewJob(
		job.Program,
		types.NewTrigger(period, job.Trigger.Block), job.Load,
		job.Label, job.Cid,
	))

	return nil
}

func (k Keeper) UpdateJobBlock(
	ctx sdk.Context, program sdk.AccAddress, label string,
	block uint64,
) error {
	job, found := k.GetJob(ctx, program, label)
	if !found {
		return types.ErrJobNotExist
	}

	if ctx.BlockHeight() >= int64(block) {
		return types.ErrBadTrigger
	}

	if job.Trigger.Period > 0 {
		return types.ErrConvertTrigger
	}

	k.SetJob(ctx, types.NewJob(
		job.Program,
		types.NewTrigger(job.Trigger.Period, block), job.Load,
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
	b := k.cdc.MustMarshal(&job)

	program, _ := sdk.AccAddressFromBech32(job.Program)
	store.Set(types.GetJobKey(program, job.Label), b)
}

func (k Keeper) DeleteJob(ctx sdk.Context, program sdk.AccAddress, label string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetJobKey(program, label))
}

func (k Keeper) SetJobs(ctx sdk.Context, jobs types.Jobs) error {
	for _, job := range jobs {
		k.SetJob(ctx, types.NewJob(
			job.Program,
			job.Trigger, job.Load,
			job.Label, job.Cid,
		))
	}
	return nil
}

func (k Keeper) SetJobStats(ctx sdk.Context, program sdk.AccAddress, label string, stats types.JobStats) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&stats)
	store.Set(types.GetJobStatsKey(program, label), b)
}

func (k Keeper) DeleteJobStats(ctx sdk.Context, program sdk.AccAddress, label string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetJobStatsKey(program, label))
}

//______________________________________________________________________

func (k Keeper) GetJob(ctx sdk.Context, program sdk.AccAddress, label string) (job types.Job, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetJobKey(program, label)

	value := store.Get(key)
	if value == nil {
		return job, false
	}

	k.cdc.MustUnmarshal(value, &job)

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
		k.cdc.MustUnmarshal(iterator.Value(), &jobStats)
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
		k.cdc.MustUnmarshal(iterator.Value(), &job)
		if cb(job) {
			break
		}
	}
}

func (k Keeper) GetJobStats(ctx sdk.Context, program sdk.AccAddress, label string) (stats types.JobStats, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetJobStatsKey(program, label)

	value := store.Get(key)
	if value == nil {
		return stats, false
	}

	k.cdc.MustUnmarshal(value, &stats)

	return stats, true
}

func (k Keeper) GetLowestFee(ctx sdk.Context) sdk.Coin {
	jobs := k.GetAllJobs(ctx)
	if len(jobs) == 0 {
		return ctypes.NewCybCoin(0)
	} else {
		jobs.Sort()
		return jobs[len(jobs)-1].Load.GasPrice
	}
}

func (k Keeper) ExecuteJobsQueue(ctx sdk.Context) {
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

	jobs := k.GetAllJobs(ctx)
	jobs.Sort()

	maxGas := k.MaxGas(ctx)
	gasBeforeCron := ctx.GasMeter().GasConsumed()
	gasUsedTotal := sdk.Gas(0)

	feeTTL := k.FeeTTL(ctx)
	maxJobs := k.MaxJobs(ctx)
	maxGasPerJob := maxGas / maxJobs

	if jobs.Len() > 0 {
		k.Logger(ctx).Info("Jobs in queue", "size", jobs.Len())
	}

	var jobsTriggered uint32
	for i, job := range jobs {
		if (job.Trigger.Period != 0 && ctx.BlockHeight()%int64(job.Trigger.Period) == 0) ||
			(job.Trigger.Period == 0 && ctx.BlockHeight() == int64(job.Trigger.Block)) {
			price := job.Load.GasPrice

			k.Logger(ctx).Info("Started job", "number", i, "gas price", price)
			jobsTriggered = jobsTriggered+1

			cacheContext, writeFn := ctx.CacheContext()
			cacheContext = cacheContext.WithGasMeter(sdk.NewGasMeter(sdk.Gas(maxGasPerJob)))

			k.Logger(ctx).Info("Context gas stats before job execution",
				"consumed", ctx.GasMeter().GasConsumed(),
			)

			remained := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumedToLimit()
			if remained < uint64(maxGasPerJob) {
				k.Logger(ctx).Info("Job break, not enough gas", "job #", i)
				break
			}

			program, _ := sdk.AccAddressFromBech32(job.Program)
			_, errExecute := k.executeJobWithSudo(
				cacheContext, program, job.Load.CallData,
			)

			gasUsedByJob := cacheContext.GasMeter().GasConsumed()
			ctx.GasMeter().ConsumeGas(gasUsedByJob, "job execution")
			if gasUsedTotal + gasUsedByJob > uint64(maxGas) {
				break
			} else {
				gasUsedTotal += gasUsedByJob
			}

			js, _ := k.GetJobStats(ctx, program, job.Label)
			// TODO move to more advanced fee system, 10X fee reducer applied (min 0.1 per gas)
			amtGasFee := price.Amount.Int64() * int64(gasUsedByJob) / 10
			amtTTLFee := (ctx.BlockHeight() - int64(js.LastBlock))*int64(feeTTL)
			amtTotalFee := amtGasFee + amtTTLFee

			k.Logger(ctx).Info("Gas job execution stats",
				"used", gasUsedByJob,
				"gas fee", amtGasFee,
				"ttl fee", amtTTLFee,
				"total fee", amtTotalFee,
			)

			fee := sdk.NewCoin(ctypes.CYB, sdk.NewInt(amtTotalFee))

			errSend := k.proxyKeeper.SendCoins(
				ctx, program, k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName), sdk.NewCoins(fee))
			if errSend != nil {
				k.DeleteJob(ctx, program, job.Label)
				k.DeleteJobStats(ctx, program, job.Label)

				k.Logger(ctx).Info("Not enough contract balance, state not applied, job killed", "Job #", i)
				continue
			}

			if errExecute != nil {
				k.Logger(ctx).Info("Job failed, state not applied", "Job #", i)
				k.Logger(ctx).Info("Failed with error: ", errExecute.Error())
			} else {
				writeFn() // apply cached context
				k.Logger(ctx).Info("Job finished, state applied", "Job #", i)
			}

			k.SetJobStats(ctx, program, job.Label,
				types.NewStats(
					program.String(), job.Label,
					js.Calls+1, js.Fees+uint64(amtTotalFee),
					js.Gas+gasUsedByJob, uint64(ctx.BlockHeight())),
			)

			if ctx.BlockHeight() == int64(job.Trigger.Block) {
				k.DeleteJob(ctx, program, job.Label)
				k.DeleteJobStats(ctx, program, job.Label)

				k.Logger(ctx).Info("Job executed at given block, deleted from queue", "Job #", i)
			}
		}
	}

	gasAfterCron := ctx.GasMeter().GasConsumed()
	if jobsTriggered > 0 {
		k.Logger(ctx).Info("Total cron gas used", "Gas used", gasAfterCron-gasBeforeCron)
	}
}

func (k Keeper) executeJobWithSudo(ctx sdk.Context, program sdk.AccAddress, msg string) ([]byte, error) {
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
		telemetry.IncrCounter(1.0, types.ModuleName, "executed")
	}()

	callData, _ := hex.DecodeString(msg)
	return k.wasmKeeper.Sudo(ctx, program, callData)
}