package types

const (
	EventTypeAddJob            = "add_job"
	EventTypeRemoveJob         = "remove_job"
	EventTypeChangeJobParticle = "change_job_particle"
	EventTypeChangeJobLabel    = "change_job_label"
	EventTypeChangeJobCallData = "change_job_call_data"
	EventTypeChangeJobGasPrice = "change_job_gas_price"
	EventTypeChangeJobPeriod   = "change_job_period"
	EventTypeChangeJobBlock    = "change_job_block"

	AttributeKeyJobProgram     = "program"
	AttributeKeyJobTrigger     = "trigger"
	AttributeKeyJobLoad        = "load"
	AttributeKeyJobLabel       = "label"
	AttributeKeyJobParticle    = "particle"
	AttributeKeyJobCallData    = "call_data"
	AttributeKeyJobGasPrice    = "gas_price"
	AttributeKeyJobPeriod      = "period"
	AttributeKeyJobBlock       = "block"

	AttributeValueCategory     = ModuleName
)
