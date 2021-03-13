package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	Contract        = "contract"
	Creator         = "creator"
	Label           = "label"
)

// RegisterRoutes registers power-related REST handlers to a router
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type AddJobReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Block 		 	 uint64 	`json:"block" yaml:"block"`
	Period 		 	 uint64 	`json:"period" yaml:"period"`
	CallData 		 string 	`json:"call_data" yaml:"call_data"`
	GasPrice 		 sdk.Coin 	`json:"gas_price" yaml:"gas_price"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	CID 		 	 string 	`json:"cid" yaml:"cid"`

}

type RemoveJobReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	  `json:"contract" yaml:"contract"`
	Creator  		 string	      `json:"creator" yaml:"creator"`
	Label 		 	 string 	  `json:"label" yaml:"label"`
}

type ChangeJobLabelReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	NewLabel 		 string 	`json:"new_label" yaml:"new_label"`
}

type ChangeJobCIDReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	CID 		 	 string 	`json:"cid" yaml:"cid"`
}

type ChangeJobCallDataReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	CallData 		 string 	`json:"call_data" yaml:"call_data"`
}

type ChangeJobGasPriceReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	GasPrice 		 sdk.Coin 	`json:"gas_price" yaml:"gas_price"`
}

type ChangeJobPeriodReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	Period 		 	 uint64 	`json:"period" yaml:"period"`
}

type ChangeJobBlockReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	Block 		 	 uint64 	`json:"block" yaml:"block"`
}