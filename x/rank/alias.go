package rank

import (
	"github.com/cybercongress/cyberd/x/rank/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
)

type (
	Params = types.Params
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	ParamKeyTable = types.ParamKeyTable
)
