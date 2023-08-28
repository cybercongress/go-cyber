package types

const (
	EventTypeCreateThought         = "create_thought"
	EventTypeForgetThought         = "forget_thought"
	EventTypeChangeThoughtParticle = "change_thought_particle"
	EventTypeChangeThoughtName     = "change_thought_name"
	EventTypeChangeThoughtInput    = "change_thought_input"
	EventTypeChangeThoughtGasPrice = "change_thought_gas_price"
	EventTypeChangeThoughtPeriod   = "change_thought_period"
	EventTypeChangeThoughtBlock    = "change_thought_block"

	AttributeKeyThoughtProgram  = "program"
	AttributeKeyThoughtTrigger  = "trigger"
	AttributeKeyThoughtLoad     = "load"
	AttributeKeyThoughtName     = "name"
	AttributeKeyThoughtParticle = "particle"
	AttributeKeyThoughtInput    = "input"
	AttributeKeyThoughtGasPrice = "gas_price"
	AttributeKeyThoughtPeriod   = "period"
	AttributeKeyThoughtBlock    = "block"

	AttributeValueCategory = ModuleName
)
