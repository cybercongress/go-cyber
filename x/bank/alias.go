package bank

import (
	"github.com/cybercongress/cyberd/x/bank/exported"
	"github.com/cybercongress/cyberd/x/bank/internal/keeper"
	"github.com/cybercongress/cyberd/x/bank/internal/types"
)

type (
	Keeper        = exported.Keeper
	IndexedKeeper = exported.IndexedKeeper

	CoinsTransferHook = types.CoinsTransferHook
)

var (
	NewKeeper        = keeper.NewKeeper
	NewIndexedKeeper = keeper.NewIndexedKeeper
)
