package bindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cybercongress/go-cyber/x/cyberbank/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	tokenfactorykeeper "github.com/cybercongress/go-cyber/x/tokenfactory/keeper"
)

func RegisterCustomPlugins(
	bank *keeper.BankProxyKeeper,
	tokenFactory *tokenfactorykeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(bank, tokenFactory)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, tokenFactory),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
