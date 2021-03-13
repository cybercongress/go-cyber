package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/tendermint/tendermint/libs/log"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/investments/types"
)

type Keeper struct {
	storeKey      	sdk.StoreKey
	cdc 			codec.BinaryMarshaler
	accountKeeper   types.AccountKeeper
	bankKeeper      bankkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryMarshaler,
	key sdk.StoreKey,
	ak 	authkeeper.AccountKeeper,
	bk  bankkeeper.Keeper,
) Keeper {
	if addr := ak.GetModuleAddress(types.InvestmentsName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.InvestmentsName))
	}

	keeper := Keeper{
		storeKey:   	key,
		cdc:        	cdc,
		accountKeeper: 	ak,
		bankKeeper:     bk,
	}
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) PutInvestment(
	ctx sdk.Context,
	investor sdk.AccAddress,
	amount sdk.Coin,
	resource string,
	length int64,
) error {
	lengthCheck := PeriodCheck(ctx, length)
	if lengthCheck == false {
		return types.ErrNotAvailableLength
	}

	err := k.AddTimeLockedCoinsToAccount(ctx, investor, sdk.NewCoins(amount), length)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrTimeLockCoins, err.Error())
	}
	err = k.Mint(ctx, investor, amount, resource, length)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrIssueCoins, err.Error())
	}
	return err
}

func (k Keeper) AddTimeLockedCoinsToAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins, length int64) error {
	acc := k.accountKeeper.GetAccount(ctx, recipientAddr)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrAccountNotFound, recipientAddr.String())
	}

	switch acc.(type) {
	case *vestingtypes.PeriodicVestingAccount:
		return k.AddTimeLockedCoinsToPeriodicVestingAccount(ctx, recipientAddr, amt, length)
	case *authtypes.BaseAccount:
		return k.AddTimeLockedCoinsToBaseAccount(ctx, recipientAddr, amt, length)
	default:
		return sdkerrors.Wrapf(types.ErrInvalidAccountType, "%T", acc)
	}
}

func (k Keeper) AddTimeLockedCoinsToPeriodicVestingAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins, length int64) error {
	k.addCoinsToVestingSchedule(ctx, recipientAddr, amt, length)
	return nil
}

func (k Keeper) AddTimeLockedCoinsToBaseAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coins, length int64) error {
	acc := k.accountKeeper.GetAccount(ctx, recipientAddr)
	// transition the account to a periodic vesting account:
	//coins := k.bankKeeper.GetBalance(ctx, acc.GetAddress(), ctypes.CYB)
	bacc := authtypes.NewBaseAccount(acc.GetAddress(), acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
	newPeriods := vestingtypes.Periods{types.NewPeriod(amt, length)}
	bva := vestingtypes.NewBaseVestingAccount(bacc, amt, ctx.BlockTime().Unix()+length)
	pva := vestingtypes.NewPeriodicVestingAccountRaw(bva, ctx.BlockTime().Unix(), newPeriods)
	k.accountKeeper.SetAccount(ctx, pva)
	return nil
}

func (k Keeper) addCoinsToVestingSchedule(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins, length int64) {
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
		newPeriodLength := (ctx.BlockTime().Unix() - vacc.EndTime) + length // 110 - 100 + 6 = 16
		newPeriod := types.NewPeriod(amt, newPeriodLength)
		vacc.VestingPeriods = append(vacc.VestingPeriods, newPeriod)
		vacc.EndTime = ctx.BlockTime().Unix() + length
		k.accountKeeper.SetAccount(ctx, vacc)
		return
	}
	// StartTime = 110
	// BlockTime = 100
	// length = 6
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
	return
}

func (k Keeper) Mint(ctx sdk.Context, recipientAddr sdk.AccAddress, amt sdk.Coin, resource string, length int64) error {
	acc := k.accountKeeper.GetAccount(ctx, recipientAddr)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrAccountNotFound, recipientAddr.String())
	}

	if amt.Amount.Mod(sdk.NewInt(ctypes.Mega*1)).IsPositive() { // TODO min amount?
		return sdkerrors.Wrapf(types.ErrMintCoins, recipientAddr.String())
	}

	cycles := sdk.NewDec(length).QuoInt64(30).TruncateDec()
	tmul := TimeMultiplier(cycles)
	smul := SupplyMultiplier(ctx, k.bankKeeper, resource)
	base := amt.Amount.Quo(sdk.NewInt(ctypes.Mega)).ToDec()

	var toMint sdk.Coin
	if resource == ctypes.VOLT {
		toMint = ctypes.NewVoltCoin(base.Mul(cycles).Mul(tmul).Mul(smul).TruncateInt64())
	} else {
		toMint = ctypes.NewAmperCoin(base.Mul(cycles).Mul(tmul).Mul(smul).TruncateInt64())
	}

	fmt.Println("[CYB]-->[VOLT || AMPER]")
	fmt.Println("Amount CYB:", amt.Amount)
	fmt.Println("Length:", length)
	fmt.Println("Cycles:", cycles)
	fmt.Println("Base:", base)
	fmt.Println("TimeMul:", tmul)
	fmt.Println("SupplyMul:", smul)
	fmt.Println("[*] MINT VOLT || AMPER:", toMint.Amount)

	err := k.bankKeeper.MintCoins(ctx, types.InvestmentsName, sdk.NewCoins(toMint))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrMintCoins, recipientAddr.String())
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.InvestmentsName, recipientAddr, sdk.NewCoins(toMint))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrSendMintedCoins, recipientAddr.String())
	}

	return nil
}

// ***JUST SANDBOX***
func PeriodCheck(ctx sdk.Context, length int64) bool {

	//var availableLength int64
	//passed := ctx.BlockTime().Sub(ctx.WithBlockHeight(1).BlockTime()).Seconds()
	//
	//switch {
	////case passed > (time.Second * 3000).Seconds():
	////	availableLength = 3000
	////case passed > (time.Second * 1500).Seconds():
	////	availableLength = 1500
	////case passed > (time.Second * 1000).Seconds():
	////	availableLength = 1000
	//case passed > (time.Second * 500).Seconds():
	//	//availableLength = 500
	//	availableLength = 60 * 60 * 24 * 30 * 12
	//default:
	//	availableLength = 300
	//}
	//return length <= availableLength
	return true
}

func TimeMultiplier(cycles sdk.Dec) sdk.Dec {
	if cycles.LT(sdk.NewDec(100)) {
		return cycles.QuoInt64(10).Add(sdk.OneDec())
	} else {
		return sdk.NewDec(11)
	}
}

func SupplyMultiplier(ctx sdk.Context, bk bankkeeper.Keeper, denom string) sdk.Dec {
	supply := bk.GetSupply(ctx).GetTotal().AmountOf(denom)
	fmt.Println("Supply:", supply)
	if supply.Int64() < ctypes.Giga*9 {
		return sdk.OneDec().Sub((sdk.NewDecFromInt(supply).QuoInt64(ctypes.Giga*10)))
	} else {
		return sdk.NewDecWithPrec(1, 1)
	}
}

//------------------------

func(k Keeper) PutCreateResource(ctx sdk.Context, resource sdk.Coin, sender, receiver sdk.AccAddress) error {
	creator, found := k.GetResource(ctx, resource.Denom)
	if !found {
		// sandbox, will be adding resource with throw params
		err := k.SetResource(ctx, sender, resource.Denom); if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidAccountType, "account: %s", sender.String())
		}
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(resource)); 		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(resource)); 		if err != nil {
			return err
		}
		//return sdkerrors.Wrapf(types.ErrResourceNotFound, "resource: %s", resource.Denom)
	}
	if !creator.Equals(sender) {
		return sdkerrors.Wrapf(types.ErrNotAuthorized, "creator: %s, address: %s", creator, sender)
	}

	err := k.SetResource(ctx, sender, resource.Denom); if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAccountType, "account: %s", sender.String())
	}
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(resource)); if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(resource)); if err != nil {
		return err
	}

	return nil
}

func(k Keeper) PutRedeemResource(ctx sdk.Context, resource sdk.Coin, sender sdk.AccAddress) error {
	creator, found := k.GetResource(ctx, resource.Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrResourceNotFound, "resource: %s", resource.Denom)
	}
	if !creator.Equals(sender) {
		return sdkerrors.Wrapf(types.ErrNotAuthorized, "creator: %s, address: %s", creator, sender)
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(resource)); if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(resource)); if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetResource(ctx sdk.Context, denom string) (addr sdk.AccAddress, exist bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ResourceStoreKey(denom))
	if bz == nil {
		return nil, false
	}
	err := addr.Unmarshal(bz); if err != nil {
		return nil, false
	}
	return addr, true
}

func (k Keeper) SetResource(ctx sdk.Context, sender sdk.AccAddress, denom string) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := sender.Marshal(); if err != nil {
		return err
	}
	store.Set(types.ResourceStoreKey(denom), bz)
	return nil
}