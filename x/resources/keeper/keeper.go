package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"math"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/libs/log"

	ctypes "github.com/cybercongress/go-cyber/types"
	bandwithkeeper "github.com/cybercongress/go-cyber/x/bandwidth/keeper"
	"github.com/cybercongress/go-cyber/x/resources/types"
)

type Keeper struct {
	cdc 			codec.BinaryMarshaler
	accountKeeper   types.AccountKeeper
	bankKeeper      types.BankKeeper
	bandwidthMeter  *bandwithkeeper.BandwidthMeter
	paramSpace      paramstypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryMarshaler,
	ak 	types.AccountKeeper,
	bk  types.BankKeeper,
	bm  *bandwithkeeper.BandwidthMeter,
	paramSpace paramstypes.Subspace,
) Keeper {
	if addr := ak.GetModuleAddress(types.ResourcesName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ResourcesName))
	}

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	keeper := Keeper{
		cdc:        	cdc,
		accountKeeper: 	ak,
		bankKeeper:     bk,
		bandwidthMeter: bm,
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

func (k Keeper) ConvertResource(
	ctx sdk.Context,
	agent sdk.AccAddress,
	amount sdk.Coin,
	resource string,
	length uint64,
) (error, sdk.Coin) {
	periodAvailableFlag := k.CheckAvailablePeriod(ctx, length)
	if periodAvailableFlag == false {
		return types.ErrNotAvailablePeriod, sdk.Coin{}
	}

	if k.bankKeeper.SpendableCoins(ctx, agent).AmountOf(ctypes.SCYB).LT(amount.Amount) {
		return sdkerrors.ErrInsufficientFunds, sdk.Coin{}
	}

	if uint32(length) < k.MinInvestmintPeriodSec(ctx) {
		return types.ErrNotAvailablePeriod, sdk.Coin{}
	}

	err, newAccFlag := k.AddTimeLockedCoinsToAccount(ctx, agent, sdk.NewCoins(amount), int64(length))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrTimeLockCoins, err.Error()), sdk.Coin{}
	}
	err, minted := k.Mint(ctx, agent, amount, resource, length, newAccFlag)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrIssueCoins, err.Error()), sdk.Coin{}
	}
	return err, minted
}

func (k Keeper) AddTimeLockedCoinsToAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins, length int64) (error, bool) {
	acc := k.accountKeeper.GetAccount(ctx, recipientAddr)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrAccountNotFound, recipientAddr.String()), false
	}

	switch acc.(type) {
	case *vestingtypes.PeriodicVestingAccount:
		return k.AddTimeLockedCoinsToPeriodicVestingAccount(ctx, recipientAddr, amt, length, false), false
	case *authtypes.BaseAccount:
		return k.AddTimeLockedCoinsToBaseAccount(ctx, recipientAddr, amt, length), true
	default:
		return sdkerrors.Wrapf(types.ErrInvalidAccountType, "%T", acc), false
	}
}

func (k Keeper) AddTimeLockedCoinsToPeriodicVestingAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins, length int64, mergeSlot bool) error {
	err := k.addCoinsToVestingSchedule(ctx, recipientAddr, amt, length, mergeSlot); if err != nil {
		return err
	}
	return nil
}

func (k Keeper) AddTimeLockedCoinsToBaseAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins, length int64) error {
	acc := k.accountKeeper.GetAccount(ctx, recipientAddr)
	bacc := authtypes.NewBaseAccount(acc.GetAddress(), acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
	newPeriods := vestingtypes.Periods{types.NewPeriod(amt, length)}
	bva := vestingtypes.NewBaseVestingAccount(bacc, amt, ctx.BlockTime().Unix()+length)
	pva := vestingtypes.NewPeriodicVestingAccountRaw(bva, ctx.BlockTime().Unix(), newPeriods)
	k.accountKeeper.SetAccount(ctx, pva)
	return nil
}

func (k Keeper) addCoinsToVestingSchedule(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins, length int64, mergeSlot bool) error {
	acc := k.accountKeeper.GetAccount(ctx, addr)
	vacc := acc.(*vestingtypes.PeriodicVestingAccount)

	// Add the new vesting coins to OriginalVesting
	vacc.OriginalVesting = vacc.OriginalVesting.Add(amt...)

	// update vesting periods
	// EndTime = 100
	// BlockTime  = 110
	// length == 6
	if vacc.EndTime < ctx.BlockTime().Unix() {
		// edge case one - the vesting account's end time is in the past (ie, all previous vesting periods have completed)
		// append a new period to the vesting account, update the end time, update the account in the store and return
		//newPeriodLength := (ctx.BlockTime().Unix() - vacc.EndTime) + length // 110 - 100 + 6 = 16
		//newPeriod := types.NewPeriod(amt, newPeriodLength)
		//vacc.VestingPeriods = append(vacc.VestingPeriods, newPeriod)
		//vacc.EndTime = ctx.BlockTime().Unix() + length
		//k.accountKeeper.SetAccount(ctx, vacc)
		//return nil

		// edge case one - the vesting account's end time is in the past (ie, all previous vesting periods have completed)
		// delete all passed periods, add a new period to the vesting account, update the end time, update the account in the store and return
		//newPeriodLength := (ctx.BlockTime().Unix() - vacc.EndTime) + length // 110 - 100 + 6 = 16
		newPeriod := types.NewPeriod(amt, length)
		vacc.VestingPeriods = append(vestingtypes.Periods{}, newPeriod)
		vacc.StartTime = ctx.BlockTime().Unix()
		vacc.EndTime = ctx.BlockTime().Unix() + length
		vacc.OriginalVesting = newPeriod.Amount
		k.accountKeeper.SetAccount(ctx, vacc)
		return nil
	}
	// StartTime = 110
	// BlockTime = 100
	// length = 6
	// this will not happen in case of resource module
	if vacc.StartTime > ctx.BlockTime().Unix() {
		// edge case two - the vesting account's start time is in the future (all periods have not started)
		// update the start time to now and adjust the period lengths in place - a new period will be inserted in the next code block
		updatedPeriods := vestingtypes.Periods{}
		for i, period := range vacc.VestingPeriods {
			updatedPeriod := period
			if i == 0 {
				updatedPeriod = types.NewPeriod(period.Amount, (vacc.StartTime-ctx.BlockTime().Unix())+period.Length) // 110 - 100 + 6 = 16
			}
			updatedPeriods = append(updatedPeriods, updatedPeriod)
		}
		vacc.VestingPeriods = updatedPeriods
		vacc.StartTime = ctx.BlockTime().Unix()
	}

	if len(vacc.VestingPeriods) == int(k.MaxSlots(ctx)) && !mergeSlot {
		// case when there are already filled slots and no one already passed
		if vacc.StartTime + vacc.VestingPeriods[0].Length > ctx.BlockTime().Unix() {
			return sdkerrors.Wrapf(types.ErrFullSlots, "not enough slots")
		} else {
			// TODO refactor next code blocks
			// case when there are passed slots and we are clean them to free space to new ones
			activePeriods := vestingtypes.Periods{}
			accumulatedLength := int64(0)
			shiftStartTime := int64(0)
			for _, period := range vacc.VestingPeriods {
				if vacc.StartTime + period.Length + accumulatedLength > ctx.BlockTime().Unix() {
					activePeriods = append(activePeriods, period)
				} else {
					shiftStartTime += period.Length
				}
				accumulatedLength += period.Length
			}

			updatedPeriods := vestingtypes.Periods{}
			updatedOriginalVesting := sdk.Coins{}
			for _, period := range activePeriods {
				updatedOriginalVesting = updatedOriginalVesting.Add(period.Amount...)
				updatedPeriods = append(updatedPeriods, period)
			}
			vacc.OriginalVesting = updatedOriginalVesting.Add(amt...)
			vacc.VestingPeriods = updatedPeriods
			vacc.StartTime = vacc.StartTime + shiftStartTime
		}
	}

	// logic for inserting a new vesting period into the existing vesting schedule
	remainingLength := vacc.EndTime - ctx.BlockTime().Unix()
	elapsedTime := ctx.BlockTime().Unix() - vacc.StartTime
	proposedEndTime := ctx.BlockTime().Unix() + length
	if remainingLength < length {
		// in the case that the proposed length is longer than the remaining length of all vesting periods, create a new period with length equal to the difference between the proposed length and the previous total length
		newPeriodLength := length - remainingLength
		newPeriod := types.NewPeriod(amt, newPeriodLength)
		vacc.VestingPeriods = append(vacc.VestingPeriods, newPeriod)
		// update the end time so that the sum of all period lengths equals endTime - startTime
		vacc.EndTime = proposedEndTime
	} else {
		// In the case that the proposed length is less than or equal to the sum of all previous period lengths, insert the period and update other periods as necessary.
		// EXAMPLE (l is length, a is amount)
		// Original Periods: {[l: 1 a: 1], [l: 2, a: 1], [l:8, a:3], [l: 5, a: 3]}
		// Period we want to insert [l: 5, a: x]
		// Expected result:
		// {[l: 1, a: 1], [l:2, a: 1], [l:2, a:x], [l:6, a:3], [l:5, a:3]}

		// StartTime = 100
		// Periods = [5,5,5,5]
		// EndTime = 120
		// BlockTime = 101
		// length = 2

		// for period in Periods:
		// iteration  1:
		// lengthCounter = 5
		// if 5 < 101 - 100 + 2 - no
		// if 5 = 3 - no
		// else
		// newperiod = 2 - 0
		newPeriods := vestingtypes.Periods{}
		lengthCounter := int64(0)
		appendRemaining := false
		for _, period := range vacc.VestingPeriods {
			if appendRemaining {
				newPeriods = append(newPeriods, period)
				continue
			}
			lengthCounter += period.Length
			if lengthCounter < elapsedTime+length { // 1
				newPeriods = append(newPeriods, period)
			} else if lengthCounter == elapsedTime+length {
				newPeriod := types.NewPeriod(period.Amount.Add(amt...), period.Length)
				newPeriods = append(newPeriods, newPeriod)
				appendRemaining = true
			} else {
				newPeriod := types.NewPeriod(amt, elapsedTime+length-types.GetTotalVestingPeriodLength(newPeriods))
				previousPeriod := types.NewPeriod(period.Amount, period.Length-newPeriod.Length)
				newPeriods = append(newPeriods, newPeriod, previousPeriod)
				appendRemaining = true
			}
		}
		vacc.VestingPeriods = newPeriods
	}
	k.accountKeeper.SetAccount(ctx, vacc)
	return nil
}

func (k Keeper) Mint(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coin, resource string, length uint64, newAccountFlag bool) (error, sdk.Coin) {
	acc := k.accountKeeper.GetAccount(ctx, recipientAddr)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrAccountNotFound, recipientAddr.String()), sdk.Coin{}
	}

	toMint := k.CalculateInvestmint(ctx, amt, resource, length)

	if toMint.Amount.LT(sdk.NewInt(1000)) {
		return sdkerrors.Wrapf(types.ErrSmallReturn, recipientAddr.String()), sdk.Coin{}
	}

	err := k.bankKeeper.MintCoins(ctx, types.ResourcesName, sdk.NewCoins(toMint))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrMintCoins, recipientAddr.String()), sdk.Coin{}
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ResourcesName, recipientAddr, sdk.NewCoins(toMint))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrSendMintedCoins, recipientAddr.String()), sdk.Coin{}
	}
	// adding converted resources to vesting schedule
	err = k.AddTimeLockedCoinsToPeriodicVestingAccount(ctx, recipientAddr, sdk.NewCoins(toMint), int64(length), true)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrTimeLockCoins, err.Error()), sdk.Coin{}
	}

	if resource == ctypes.VOLT {
		k.bandwidthMeter.AddToDesirableBandwidth(ctx, toMint.Amount.Uint64())
	}
	// Set personal bandwidth of newcomers to 1000
	if newAccountFlag {
		band := k.bandwidthMeter.GetAccountBandwidth(ctx, recipientAddr)
		if band.MaxValue == 0 {
			k.bandwidthMeter.InitChargeAccountBandwidth(ctx, band, 1000)
		}
	}

	return nil, toMint
}

func (k Keeper) CalculateInvestmint(ctx sdk.Context, amt sdk.Coin, resource string, length uint64) sdk.Coin {
	var toMint sdk.Coin
	switch resource {
	case ctypes.VOLT:
		//cycles := sdk.NewDec(int64(length)).QuoInt64(int64(10)) // for local dev
		cycles := sdk.NewDec(int64(length)).QuoInt64(int64(k.BaseInvestmintPeriodVolt(ctx)))
		base := sdk.NewDec(amt.Amount.Int64()).QuoInt64(k.BaseInvestmintAmountVolt(ctx).Amount.Int64())
		halving := sdk.NewDecWithPrec(int64(math.Pow(0.5, float64(ctx.BlockHeight() / int64(k.BaseHalvingPeriodVolt(ctx))))*10000),4)

		toMint = ctypes.NewVoltCoin(base.Mul(cycles).Mul(halving).Mul(sdk.NewDec(1000)).TruncateInt64())

		k.Logger(ctx).Info("Investmint", "cycles", cycles.String(), "base", base.String(), "halving", halving.String(), "mint", toMint.String())
	case ctypes.AMPER:
		//cycles := sdk.NewDec(int64(length)).QuoInt64(int64(10)) // for local dev
		cycles := sdk.NewDec(int64(length)).QuoInt64(int64(k.BaseInvestmintPeriodAmpere(ctx)))
		base := sdk.NewDec(amt.Amount.Int64()).QuoInt64(k.BaseInvestmintAmountAmpere(ctx).Amount.Int64())
		halving := sdk.NewDecWithPrec(int64(math.Pow(0.5, float64(ctx.BlockHeight() / int64(k.BaseHalvingPeriodAmpere(ctx))))*10000),4)

		toMint = ctypes.NewAmperCoin(base.Mul(cycles).Mul(halving).Mul(sdk.NewDec(1000)).TruncateInt64())

		k.Logger(ctx).Info("Investmint", "cycles", cycles.String(), "base", base.String(), "halving", halving.String(), "mint", toMint.String())
	}
	return toMint
}

func (k Keeper) CheckAvailablePeriod(ctx sdk.Context, length uint64) bool {
	var availableLength int64
	passed := ctx.BlockHeight()

	// TODO more advanced logic will be applied after launch
	switch {
	case passed > 600:
		// 12 months max freeze after 3072000 blocks passed after launch
		//availableLength = 31104000
		availableLength = 3600
	case passed > 300:
		// 6 months max freeze after 1536000 blocks passed after launch
		//availableLength = 15552000
		availableLength = 2400
	default:
		// 3 months max freeze before 1536000 blocks passed after launch
		//availableLength = 7776000
		availableLength = 1200
	}
	return length < uint64(availableLength)
}

func (k Keeper) MaxSlots(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyMaxSlots, &res)
	return
}

func (k Keeper) BaseHalvingPeriodVolt(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyBaseHalvingPeriodVolt, &res)
	return
}

func (k Keeper) BaseHalvingPeriodAmpere(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyBaseHalvingPeriodAmpere, &res)
	return
}

func (k Keeper) BaseInvestmintPeriodVolt(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyBaseInvestmintPeriodVolt, &res)
	return
}

func (k Keeper) BaseInvestmintPeriodAmpere(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyBaseInvestmintPeriodAmpere, &res)
	return
}

func (k Keeper) BaseInvestmintAmountVolt(ctx sdk.Context) (res sdk.Coin) {
	k.paramSpace.Get(ctx, types.KeyBaseInvestmintAmountVolt, &res)
	return
}

func (k Keeper) BaseInvestmintAmountAmpere(ctx sdk.Context) (res sdk.Coin) {
	k.paramSpace.Get(ctx, types.KeyBaseInvestmintAmountAmpere, &res)
	return
}

func (k Keeper) MinInvestmintPeriodSec(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyMinInvestmintPeriodSec, &res)
	return
}