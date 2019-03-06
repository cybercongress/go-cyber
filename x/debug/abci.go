package debug

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"math/rand"
	"os"
)

func BeginBlocker(state *State, req abci.RequestBeginBlock, log log.Logger) {

	if state.Opts.FailRandomlyInNextNBlocks != 0 {
		randNum := rand.Int63n(state.Opts.FailRandomlyInNextNBlocks) + 1
		state.Opts.FailBeforeBlock = state.StartupBlock + randNum
		state.Opts.FailRandomlyInNextNBlocks = 0
		log.Info("Scheduled forced panic at begin blocker ", "height", state.Opts.FailBeforeBlock)
	}

	failBeforeBlock := state.Opts.FailBeforeBlock
	if failBeforeBlock != 0 && req.Header.Height == failBeforeBlock {
		log.Info("Forced panic at begin blocker", "height", failBeforeBlock)
		os.Exit(1)
	}
}
