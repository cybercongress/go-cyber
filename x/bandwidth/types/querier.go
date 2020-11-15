package types

import (
sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryParameters 		= "parameters"
	QueryLoad				= "load"
	QueryPrice				= "price"
	QueryAccount			= "account"
)

type QueryAccountParams struct {
	Account sdk.AccAddress
}

func NewQueryAccountParams(addr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{
		Account: addr,
	}
}

type ResultPrice struct {
	Price float64 `amino:"unsafe" json:"price" yaml:"price"`
}

func NewResultPrice(price float64) ResultPrice {
	return ResultPrice{
		Price: price,
	}
}

type ResultLoad struct {
	Load float64 `amino:"unsafe" json:"load" yaml:"load"`
}

func NewResultLoad(load float64) ResultLoad {
	return ResultLoad{
		Load: load,
	}
}