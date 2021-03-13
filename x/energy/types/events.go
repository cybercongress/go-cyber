package types

// power module event types
const (
	EventTypeCreateEnergyRoute   	= "create_energy_route"
	EventTypeEditEnergyRoute 		= "edit_energy_route"
	EventTypeDeleteEnergyRoute 		= "delete_energy_route"
	EventTypeEditEnergyRouteAlias 	= "edit_energy_route_alias"


	AttributeKeySource    			= "source"
	AttributeKeyDestination 		= "destination"
	AttributeKeyAlias 				= "alias"
	AttributeKeyAmount	  			= "value"

	AttributeValueCategory   		= ModuleName
)
