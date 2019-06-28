# Setup development environment

## Prestart
* Install Golang 1.11+
* Install GoLand IDE

## Import project to GoLand
Open Project in GoLand by selecting: Open Project -> selecting cloned repository root folder
![Open project](https://ipfs.io/ipfs/QmYQxSKzQkkpCofuHmbJZqYoTo4cuvHQuUNYkJRBKgSduL)

Enable **go mod** package management
![Enable go mod](https://ipfs.io/ipfs/Qmaz3o7LjAG9bNhE8VkHSUokM8EqgfcVXT9R5qB3XxLZSe)
Wait for dependency downloading and indexation

## Add Run Configurations
Add testnet configuration
![Generate testnet](https://ipfs.io/ipfs/QmVbNBnFAwiPELL41iKh9MNDR2uywcFby6qtP5cYHc3Jvv)

Add run configuration
![Run node](https://ipfs.io/ipfs/QmUuPDbwYJpdibQRt8yGbqFiseoyMZf99DT8ULjt8YwDNh)

Add reset configuration
![Reset node data](https://ipfs.io/ipfs/Qmdp5MnbqCh38Su8eiBSr6ZVYPVL9hg8ap6hTFGJPYwd22)

## Running Node

Before node running, execute **generate testnet** run configuration.
Folder **mytestnet** will be added to the project root.
In **node0** subfolder you can find daemon and cli folders.
Daemon folder will contain validator node data.
In cli folder you can find initial validator seed.

After, just run **run node** run configuration.
You can reset chains data to genesis at any time by executing **reset** run configuration
