package v5

import (
	"github.com/cybercongress/go-cyber/v5/app/upgrades"
)

const (
	UpgradeName = "v5"

	UpgradeHeight = 373421
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
