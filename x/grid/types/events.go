package types

const (
	EventTypeCreateRoute   	= "create_route"
	EventTypeEditRoute 		= "edit_route"
	EventTypeDeleteRoute 	= "delete_route"
	EventTypeEditRouteAlias = "edit_route_alias"


	AttributeKeySource      = "source" // TODO change to sender?
	AttributeKeyDestination = "destination"
	AttributeKeyAlias       = "alias"
	AttributeKeyValue       = "value"

	AttributeValueCategory  = ModuleName
)
