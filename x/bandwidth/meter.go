package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/util"
	"github.com/tendermint/tendermint/libs/db"
)

type BandwidthMeter interface {
	Handle(ctx sdk.Context, tx auth.StdTx) error
}

type StdBandwidthMeter struct {
	db              db.DB
	stakeKeeper     stake.Keeper
	accountKeeper   auth.AccountKeeper
	accountKey      *sdk.KVStoreKey
	bandwidthKeeper BandwidthStateKeeper
}

func (bm *StdBandwidthMeter) Handle(ctx sdk.Context, tx auth.StdTx) sdk.Error {

	account, sdkErr := bm.getAccount(ctx, tx)
	if sdkErr != nil {
		return sdkErr
	}

	accountBandwidth, err := bm.bandwidthKeeper.GetAccountBandwidth(account, ctx)
	if err != nil {
		return sdk.ErrInternal("Cannot process tx: " + err.Error())
	}

	contextForStake, err := util.NewContextWithMSVersion(bm.db, ctx.BlockHeight() - RecoveryPeriod, bm.accountKey)
	if err != nil {
		return sdk.ErrInternal("Cannot process tx: " + err.Error())
	}

	addressStake := bm.accountKeeper.GetAccount(contextForStake, account).GetCoins().AmountOf(coin.CBD)

	pool := bm.stakeKeeper.GetPool(ctx)
	totalStake := pool.BondedTokens.RoundInt64() + pool.LooseTokens.RoundInt64()

	addressMaxBandwidth := ((addressStake.Int64() / totalStake) * MaxNetworkBandwidth) / 2

	accountBandwidth.Recover(addressMaxBandwidth, ctx.BlockHeight())
	bandwidthForTx := int64(1) // TODO: add calculation

	if !accountBandwidth.HasEnoughRemained(bandwidthForTx) {
		return sdk.ErrInternal("Bandwidth limit exceeded")
	}

	accountBandwidth.Consume(bandwidthForTx)

	bm.bandwidthKeeper.SetAccountBandwidth(ctx, accountBandwidth)

	return nil
}

func (bm *StdBandwidthMeter) getAccount(ctx sdk.Context, tx auth.StdTx) (sdk.AccAddress, sdk.Error) {

	if tx.Msgs == nil || len(tx.Msgs) == 0 {
		return nil, sdk.ErrInternal("Tx.GetMsgs() must return at least one message in list")
	}

	if err := tx.ValidateBasic(); err != nil {
		return nil, err
	}

	// signers acc [0] bandwidth will be consumed
	account := tx.GetSigners()[0]

	acc := bm.accountKeeper.GetAccount(ctx, account)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(account.String())
	}

	return account, nil
}
