package v2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cybercongress/go-cyber/app/upgrades"
)

const UpgradeName = "cyberfrey"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
