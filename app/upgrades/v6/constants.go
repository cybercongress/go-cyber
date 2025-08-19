package v6

import (
	"github.com/cybercongress/go-cyber/v6/app/upgrades"
)

const (
	UpgradeName = "v6"

	UpgradeHeight = 42_000_000
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
