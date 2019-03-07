package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// should be removed after cosmos refactor
func NewAnteHandler(ak auth.AccountKeeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {

		//todo really need it?
		newCtx = auth.SetGasMeter(true, ctx, 0)

		// all transactions must be of type auth.StdTx
		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			return newCtx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}

		params := ak.GetParams(ctx)
		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err.Result(), true
		}

		if res := auth.ValidateMemo(stdTx, params); !res.IsOK() {
			return newCtx, res, true
		}

		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		signerAddrs := stdTx.GetSigners()
		stdSigs := stdTx.GetSignatures()
		isGenesis := ctx.BlockHeight() == 0
		signerAccs := make([]auth.Account, len(signerAddrs))

		for i := 0; i < len(stdSigs); i++ {

			signerAccs[i], res = auth.GetSignerAcc(newCtx, ak, signerAddrs[i])
			if !res.IsOK() {
				return newCtx, res, true
			}

			// check signature, return account with incremented nonce
			signBytes := auth.GetSignBytes(newCtx.ChainID(), stdTx, signerAccs[i], isGenesis)
			signerAccs[i], res = processSig(signerAccs[i], stdSigs[i], signBytes, simulate)
			if !res.IsOK() {
				return newCtx, res, true
			}

			ak.SetAccount(newCtx, signerAccs[i])
		}

		return newCtx, sdk.Result{}, false // continue...
	}
}

// verify the signature and increment the sequence. If the account doesn't have a pubkey, set it.
func processSig(
	acc auth.Account, sig auth.StdSignature, signBytes []byte, simulate bool,
) (updatedAcc auth.Account, res sdk.Result) {

	pubKey, res := auth.ProcessPubKey(acc, sig, simulate)
	if !res.IsOK() {
		return nil, res
	}

	err := acc.SetPubKey(pubKey)
	if err != nil {
		return nil, sdk.ErrInternal("setting PubKey on signer's account").Result()
	}

	if !simulate && !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return nil, sdk.ErrUnauthorized("signature verification failed").Result()
	}

	if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
		panic(err)
	}

	return acc, res
}
