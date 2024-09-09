package v4

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	clocktypes "github.com/cybercongress/go-cyber/v4/x/clock/types"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v4/x/tokenfactory/types"

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
			tokenfactorytypes.ModuleName,
			nft.ModuleName,
			clocktypes.ModuleName,
			icacontrollertypes.StoreKey,
			icahosttypes.StoreKey,
		},
	},
}
