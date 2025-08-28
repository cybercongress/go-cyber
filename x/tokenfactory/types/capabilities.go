package types

import "slices"

const (
	EnableSetMetadata   = "enable_metadata"
	EnableForceTransfer = "enable_force_transfer"
	EnableBurnFrom      = "enable_burn_from"
)

func IsCapabilityEnabled(enabledCapabilities []string, capability string) bool {
	if len(enabledCapabilities) == 0 {
		return true
	}

	return slices.Contains(enabledCapabilities, capability)
}
