package v6

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cybercongress/go-cyber/v6/app/upgrades"
)

const UpgradeName = "v6"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV6UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
