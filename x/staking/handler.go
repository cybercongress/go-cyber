package staking

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ctypes "github.com/cybercongress/go-cyber/types"
)

func WrapStakingHandler(
	ak authkeeper.AccountKeeper, sk stakingkeeper.Keeper,
) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *stakingtypes.MsgDelegate:
			result, err :=  ValidateDelegate(ctx, ak, msg)
			if err != nil {
				return nil, err
			}
			if result == false {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized message: %T", msg)
			}
		case *stakingtypes.MsgUndelegate:
			result, err := ValidateUndelegate(ctx, ak, msg)
			if err != nil {
				return nil, err
			}
			if result == false {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized message: %T", msg)
			}
		case *stakingtypes.MsgCreateValidator:
			result, err := ValidateCreateValidator(msg)
			if err != nil {
				return nil, err
			}
			if result == false {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized message: %T", msg)
			}
		}
		handler := staking.NewHandler(sk)
		return handler(ctx, msg)
	}
}

func ValidateDelegate(ctx sdk.Context,
	ak authkeeper.AccountKeeper, msg *stakingtypes.MsgDelegate,
) (bool, error) {
	fmt.Println("[*] Process Delegate")

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); if err != nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "bad account")
	}

	account := ak.GetAccount(ctx, delegator)

	pva, ok := account.(*vestingtypes.PeriodicVestingAccount)
	if !ok {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account is not vesting account type")
	}

	vi := pva.GetVestingCoins(time.Now()).AmountOf(ctypes.CYB)
	dv := pva.GetDelegatedVesting().AmountOf(ctypes.CYB)

	fmt.Println("[~] VESTING: ", vi)
	fmt.Println("[~] DELEGATED VESTING: ", dv)
	fmt.Println("[~] DELEGATE AMOUNT: ", msg.Amount.Amount)
	//fmt.Println("[&] ACCOUNT: ", ak.GetAccount(ctx, delegator).String())
	//fmt.Println("**************************\n")

	if vi.Sub(dv).GTE(msg.Amount.Amount) {
		fmt.Println("[!] Process Delegate - TRUE")
		//fmt.Println("[&] ACCOUNT: ", pva.String())
		//fmt.Println("----------------------------\n")
		return true, nil
	}

	fmt.Println("[!] Process Delegate - FALSE")
	return false, nil
}

func ValidateUndelegate(ctx sdk.Context,
	ak authkeeper.AccountKeeper, msg *stakingtypes.MsgUndelegate,
) (bool, error) {
	fmt.Println("[*] Process Undelegate")

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); if err != nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "bad account")
	}
	account := ak.GetAccount(ctx, delegator)

	pva, ok := account.(*vestingtypes.PeriodicVestingAccount)
	if !ok {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account is not vesting account type")
	}

	dv := pva.GetDelegatedVesting().AmountOf(ctypes.CYB)
	vi := pva.GetVestingCoins(time.Now()).AmountOf(ctypes.CYB)

	fmt.Println("[~] VESTING: ", vi)
	fmt.Println("[~] DELEGATED VESTING: ", dv)
	fmt.Println("[~] UNDELEGATE AMOUNT: ", msg.Amount.Amount)
	//fmt.Println("[&] ACCOUNT: ", ak.GetAccount(ctx, delegator).String())
	//fmt.Println("**************************\n")

	if dv.Sub(vi).GTE(msg.Amount.Amount) {
		fmt.Println("[!] Process Undelegate - TRUE")
		//pva.TrackUndelegation(sdk.NewCoins(msg.Amount))
		//fmt.Println("[&] ACCOUNT: ", pva.String())
		//fmt.Println("----------------------------\n")
		return true, nil
	}

	fmt.Println("[!] Process Undelegate - FALSE")
	return false, nil
}

func ValidateCreateValidator(msg *stakingtypes.MsgCreateValidator) (bool, error) {
	fmt.Println("[*] Process Create Validator")

	//if msg.MinSelfDelegation.GTE(sdk.NewInt(10*ctypes.Giga)) {
	if true {
		fmt.Println("[!] Process Create Validator - TRUE")
		return true, nil
	}

	fmt.Println("[!] Process Create Validator - FALSE")
	return false, nil
}
