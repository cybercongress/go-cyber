# How to upgrade cyber node with cosmosd

Upgrade to newer version of binaries will be covered in this guide. Go-cyber using [cosmosd](https://github.com/regen-network/cosmosd) upgrade manager and [upgrade module](https://docs.cosmos.network/master/modules/upgrade/) for on-chain upgrades.

If you set up your node using this [Run validator guide](https://github.com/cybercongress/go-cyber/blob/master/docs/run_validator.md) it should be pretty eady to prepare for upgrade.

## Definition

Current upgrade named Darwin, so all further referencies will be made to that name. Overal upgrade process is basically tied to changing binary file of the cyberd to new version, at the same time for all validator nodes. Governance proposal for upgrade will define name for upgrade and block in the chain, when this upgrade should happen. New binary file must be pre-compiled (or downloaded), and placed rigth inside `.cyberd` folder. Whole uprade process is as simple as theese steps:

- newly compiled `cyberd` binary file is being placed into correct directory

- when desired upgrade block comes `cosmosd` catches event and stops cyber node

- `cosmosd` starting node from new binary file from the directory with upgrade name

### Cosmosd explained

[Cosmosd](https://github.com/regen-network/cosmosd) is a tool designed by [Regen-Network](https://github.com/regen-network) devs. It is designed for smooth and configurable management of upgrading binaries as a live chain is upgraded, as well as makes syncing a full node for genesis simplier. The upgrade manager will monitor the stdout of the daemon to look for messages from the upgrade module indicating a pending or required upgrade and act appropriately.

Since the version 0.1.6 (Euler-6) upgrade module has been integrated to cyber, which add possibility to make new kind of proposal called `software-upgrade`. This proposal, if passed, set specific block or time when upgrade should happen as well as specifies unique upgrade name. By default cosmosd is launching deamon file that is located under `DAEMON_HOME/upgrade_manager/genesis/` folder. When upgrade time comes `cyberd` generates to it's stdout sting looks like following:

```bash
Jul 16 17:25:11 mars cosmosd[5791]: E[2020-07-16|17:25:11.325] UPGRADE "darwin" NEEDED at height: 2000:     module=main
```

Cosomosd catching that string at deamon stdout, stops cyberd and changes symbolic link current from

```bash
current -> /root/.cyberd_testnet/upgrade_manager/genesis
```

to new one according with upgrade name

```bash
current -> /root/.cyberd_testnet/upgrade_manager/upgrades/darwin
```

After link changed to new binary `cosmosd` will attempt to restart node using new daemon binary file.

## Preparation

There is a few checkpoints each validator must comlpy with to make upgrade succesfull:

1. Correct version of cosmosd

2. Correctly set `cyberd.service` file

3. New `cyberd` binary file is compiled and placed correctly

Let's dig into each point separately.

### Cosmosd version check

Cosmosd must be compiled from the latest commit of [cosmosd-repo](https://github.com/regen-network/cosmosd). Correct commit is `984175ff7c0a3815dcf2e99329cb3fbf7d7bb0b9`, so if you have cosmosd repo cloned on your node just do:

```bash
cd /path_to_cosmosd/cosmosd
git log
```

and make sure current commit is same or newer than `984175ff7c0a3815dcf2e99329cb3fbf7d7bb0b9`. Released version 0.2.0 has a bug, which will fail upgrade. So if your cosmosd folder loks like this

```bash
cosmosd-0.2.0
```

you are on the wrong version, and will have to correct this.

In case when you removed cosmosd directory, it's advised to clone repo again and rebuild cosmosd to make sure you on 100% correct version.

To get correct version do the following:

```bash
git clone https://github.com/regen-network/cosmosd.git
cd cosmosd/
go build
service cyberd stop
mv cosmosd $DAEMON_HOME/
chmod +x $DAEMON_HOME/cosmosd
service cyberd start
```

Attention - please, build cosmosd from latest master and double-check your version

### Cyberd.service check

Cyberd.service responsibe for cyberd operation must be 100% similar to current [doc](https://github.com/cybercongress/go-cyber/blob/master/docs/run_validator.md).

To check it open

```bash
sudo nano /etc/systemd/system/cyberd.service
```

and verify that it identical to followng (of cource, except usernames and path to folders):

```json
[Unit]
Description=Cyber Node
After=network-online.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/.cyberd/
ExecStart=/home/ubuntu/.cyberd/cosmosd start --compute-rank-on-gpu=true
Environment=DAEMON_HOME=/home/ubuntu/.cyberd
Environment=DAEMON_NAME=cyberd
Environment=GAIA_HOME=/home/ubuntu/.cyberd
Environment=DAEMON_RESTART_AFTER_UPGRADE=on
Restart=always
RestartSec=3m
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```

The most important part is

```bash
Environment=DAEMON_RESTART_AFTER_UPGRADE=on
```

Because it was missing in out earlier versions of run validator guide, and higly important for node start after upgrade.

If service file was modified, systemd must be reloaded to apply changes:

```bash
systemctl daemon-reload
```

### Prepraing new cyberd binary

Upgrade will be released with tag v0.1.6.3, so technically all that required to get it is pull changes inside cloned go-cyber repo:

```bash
cd /<path_to_go-cyber>/go-cyber/
git reset --hard
git checkout master
git pull
```

and then build all in the same way as during setup, cuda code first, `cyberd` binary next:

```bash
cd /<path_to_go-cyber>/go-cyber/x/rank/cuda/
make
cd /<path_to_go-cyber>/go-cyber/
make build
```

After all above done, new `cyberd` and `cyberdcli` binaries must appear inside ~/go-cyber/build/.

Next we need to create folders tree with upgrade name and place binaries to appropriate location (all of the following execute form `go-cyber` directory):

```bash
mkdir $DAEMON_HOME/upgrade_manager/upgrades/
mkdir $DAEMON_HOME/upgrade_manager/upgrades/darwin/
mkdir $DAEMON_HOME/upgrade_manager/upgrades/darwin/bin/
cp build/cyberd $DAEMON_HOME/upgrade_manager/upgrades/darwin/bin/
cp build/cyberdcli /usr/local/bin/
cp build/cyberd /usr/local/bin/
chmod +x $DAEMON_HOME/upgrade_manager/upgrades/darwin/bin/cyberd
```

To check that everything builded correclty run:

```bash
$DAEMON_HOME/upgrade_manager/upgrades/darwin/bin/cyberd version
```

Version must be `v0.1.6.3`.

And verify upgrade_manager folder from .cyberd directory looks like following, run

```bash
tree
```

**it migth require installing tree**

then check that your upgrade_manager directory looks similar:

```bash
.
├── current -> /root/.cyberd/upgrade_manager/genesis
├── genesis
│   └── bin
│       └── cyberd
└── upgrades
    └── darwin
        └── bin
            └── cyberd
```

And here thats it for preparation side. It everythng set correctly upgrade will happen automatically.

After upgrade current link will be automatically changed to ```current -> /root/.cyberd/upgrade_manager/upgrades/darwin```

## Upgrade itself

When upgrade proposal will be submited and upgrade block will come, cosmosd will do following on all nodes across the chain automatically:

- Stop cyber node

- Change `current` link to` -> .cyberd/upgrade_manager/upgrades/darwin`

- Start node

As soon as there will be more than ~66.66% of voting power at succesfully upgraded validators, chain will continue block production.

**Important**

This manual would not cover docker setup, for that case you'll have to manually build new conatiner and start it over the node state, instead of old one, manually at/or after the time of the upgrade.

## Manual upgrade

If something goes wrong and cosmosd didn't start with new binary at upgrade block you need do following:

1. Stop cyber node:

    ```bash
     service cyberd stop
    ```

2. Go to .cyberd directory and remove current link:

    ```bash
    cd /$DAEMON_HOME/upgrade_manager/
    rm current
    ```

3. Create new symbolic link to directory with upgraded binary

    ```bash
    ln -s /$DAEMON_HOME/upgrade_manager/upgrades/darwin/ /$DAEMON_HOME/upgrade_manager/current
    ```

    that will point cosmosd to correct binary location.

4. Start cyberd service

    ```bash
    service cyberd start
    ```
