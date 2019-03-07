# Setup development environment

## Prestart
* Install Golang 1.11+
* Install GoLand IDE

## Import project to GoLand
Open Project in GoLand by selecting: Open Project -> selecting cloned repository root folder
![Open project](../img/open-project.png)

Enable **go mod** package management
![Enable go mod](../img/enable-go-mod.png)
Wait for dependency downloading and indexation

## Add Run Configurations
Add testnet configuration
![Generate testnet](../img/generate-testnet.png)

Add run configuration
![Run node](../img/run-node.png)

Add reset configuration
![Reset node data](../img/reset-node-data.png)

## Running Node

Before node running, execute **generate testnet** run configuration. 
Folder **mytestnet** will be added to the project root.
In **node0** subfolder you can find daemon and cli folders.
Daemon folder will contain validator node data.
In cli folder you can find initial validator seed.

After, just run **run node** run configuration.
You can reset chains data to genesis at any time by executing **reset** run configuration
