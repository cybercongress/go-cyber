package types

const (
	ModuleName 		= "resources"
	StoreKey 		= ModuleName
	RouterKey 		= ModuleName
	QuerierRoute 	= ModuleName

	ResourcesName = "resources"

	ActionConvert    	  = "convert"
	ActionCreateResource  = "create_resource"
	ActionRedeemResource  = "redeem_resource"
)

var (
	ResourcePrefix = []byte{0x01}
)

func ResourceStoreKey(denom string) []byte {
	return append(ResourcePrefix, []byte(denom)...)
}