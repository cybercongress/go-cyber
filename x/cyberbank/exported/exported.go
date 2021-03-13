package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

type Keeper interface {
	sdkbank.Keeper

	AddHook(types.CoinsTransferHook)
	//GetAccountUnboundedStake(sdk.Context, sdk.AccAddress) int64
	//GetAccountBoundedStake(sdk.Context, sdk.AccAddress) int64
	//GetAccountTotalStake(sdk.Context, sdk.AccAddress) int64
	//GetAccStakePercentage(sdk.Context, sdk.AccAddress) float64
	//GetTotalSupply(sdk.Context) int64
	OnCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress)
}