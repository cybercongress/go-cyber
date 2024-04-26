package v3

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"

	"github.com/cybercongress/go-cyber/v4/app/upgrades"
)

const UpgradeName = "v4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV4UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			crisistypes.ModuleName,
			consensustypes.ModuleName,
			resourcestypes.ModuleName,
		},
	},
}
