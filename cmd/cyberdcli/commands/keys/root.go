package keys

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {

	keyCommand := keys.Commands()
	subcommands := keyCommand.Commands()
	for _, subcommand := range subcommands {
		if subcommand.Name() == "add" {
			subcommand.AddCommand(importPrivateKeyCmd())
		}
	}

	return keyCommand
}
