package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/resources/types"
)

func WrapStakingHandler(
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *stakingtypes.MsgDelegate:
			result, err :=  ProcessDelegate(ctx, bk, msg)
			if err != nil {
				return nil, err
			}
			if result == false {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized message: %T", msg)
			}
		case *stakingtypes.MsgUndelegate:
			result, err := ProcessUndelegate(ctx, bk, msg)
			if err != nil {
				return nil, err
			}
			if result == false {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized message: %T", msg)
			}
		case *stakingtypes.MsgCreateValidator:
			result, err := ValidateCreateValidator(ctx, bk , msg)
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

func ProcessDelegate(ctx sdk.Context,
	bk bankkeeper.Keeper, msg *stakingtypes.MsgDelegate,
) (bool, error) {
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); if err != nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "bad account")
	}

	toMint := sdk.NewCoin(ctypes.SCYB, msg.Amount.Amount)
	err = bk.MintCoins(ctx, types.ResourcesName, sdk.NewCoins(toMint))
	if err != nil {
		return false, sdkerrors.Wrapf(types.ErrMintCoins, delegator.String())
	}
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ResourcesName, delegator, sdk.NewCoins(toMint))
	if err != nil {
		return false, sdkerrors.Wrapf(types.ErrSendMintedCoins, delegator.String())
	}

	return true, nil
}

func ProcessUndelegate(ctx sdk.Context,
	bk bankkeeper.Keeper, msg *stakingtypes.MsgUndelegate,
) (bool, error) {
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); if err != nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "bad account")
	}

	toBurn := sdk.NewCoin(ctypes.SCYB, msg.Amount.Amount)
	err = bk.SendCoinsFromAccountToModule(ctx, delegator, types.ResourcesName, sdk.NewCoins(toBurn))
	if err != nil {
		return false, sdkerrors.Wrapf(types.ErrSendMintedCoins, delegator.String())
	}
	err = bk.BurnCoins(ctx, types.ResourcesName, sdk.NewCoins(toBurn))
	if err != nil {
		return false, sdkerrors.Wrapf(types.ErrBurnCoins, delegator.String())
	}

	return true, nil
}

func ValidateCreateValidator(ctx sdk.Context,
	bk bankkeeper.Keeper, msg *stakingtypes.MsgCreateValidator,
) (bool, error) {
	//if msg.MinSelfDelegation.GTE(sdk.NewInt(10*ctypes.Giga)) {
	//	return true, nil
	//}
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); if err != nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "bad account")
	}

	toMint := sdk.NewCoin(ctypes.SCYB, msg.Value.Amount)
	err = bk.MintCoins(ctx, types.ResourcesName, sdk.NewCoins(toMint))
	if err != nil {
		return false, sdkerrors.Wrapf(types.ErrMintCoins, delegator.String())
	}
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ResourcesName, delegator, sdk.NewCoins(toMint))
	if err != nil {
		return false, sdkerrors.Wrapf(types.ErrSendMintedCoins, delegator.String())
	}


	return false, nil
}
