package cyberbank

import (
	"github.com/cybercongress/go-cyber/x/cyberbank/keeper"
	//"github.com/litvintech/cyber/x/cyberbank/internal/keeper/proxy"
	"github.com/cybercongress/go-cyber/x/cyberbank/types"
)

const (
	ModuleName        = types.ModuleName
)

var (
	NewIndexedKeeper  = keeper.NewIndexedKeeper
	NewWrap			  = keeper.Wrap
)

type (
	Keeper       	  = keeper.IndexedKeeper
	Proxy		      = keeper.Proxy
)
