package types

const (
	ModuleName 		= "investment"
	StoreKey 		= ModuleName
	RouterKey 		= ModuleName
	QuerierRoute 	= ModuleName

	InvestmentsName = "investments"

	ActionInvest    	  = "invest"
	ActionCreateResource  = "create_resource"
	ActionRedeemResource  = "redeem_resource"
)

var (
	ResourcePrefix = []byte{0x01}
)

func ResourceStoreKey(denom string) []byte {
	return append(ResourcePrefix, []byte(denom)...)
}