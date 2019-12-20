# Validator launch with Tendermint KMS + Leger Nano

In this guide, you'll find the full information on how to set up your validator at cyberd by using Tenderming KMS and Ledger nano S as a keystore.

## Preparing your Ledger

We assume you have one. If not, only buy them from trusted sources!

First of all, we'll need to prepare the Ledger to be able to work with Tendermint KMS.
We need to install the Tendermint app onto the Ledger using Ledger Live.

*Note: at the moment, you might need to enable the `developer mode` in Ledger Live settings*

If you have (initially) set up your Ledger on a different machine than the one with cyberd, you should make sure that the Ledger device is recognized by Legder Live. The best way to do this is by installing [Ledger Live](https://shop.ledger.com/pages/ledger-live) onto that particular machine and attempting to connect the Ledger device to it. This will show any possible issues or/and error codes (if they exist, of course). To deal with them, please use the [Fix connection issues](https://support.ledger.com/hc/en-us/articles/115005165269-Fix-connection-issues) guide from Ledger.

## Installing Tendermint KMS onto the node

Normally, this has to be done according to the Tendermint [guide](https://github.com/tendermint/kms), but we will need a few extras (all instructions have been tested on KMS v0.6.3).

### Installation

You will need the following prerequisites:

- **Rust** (stable; 1.35+) [install](https://rustup.rs/)
- **C compiler**: e.g. gcc, clang
- **pkg-config**
- **libusb** (1.0+) : for Debian/Ubuntu: `apt install libusb-1.0-0-dev`

NOTE (x86_64 only): Configure `RUSTFLAGS` environment variable:
`export RUSTFLAGS=-Ctarget-feature=+aes,+ssse3`  (generally it would be necessary to add these flags to ~/.bash-profile or ~/.profile to make them available after re-login).

We are ready to install KMS. There are 2 ways to do this: compile from source or install with Rusts `cargo-install`. We'll use the first option.

### Compiling from source code

`tmkms` can be compiled directly from the git repository source code, using the following commands:

``` js
git clone https://github.com/tendermint/kms.git && cd kms
cargo build --release --features=ledgertm,softsign
```

If successful, it will produce the `tmkms` executable located at:
`./target/release/tmkms`.

## KMS configuration

After compiling, we should create a settings file - `tmkms.toml`, the `secret_connection.key` file and adjust cyberd node settings.

First of all, we need to generate a connection key that will be used for communication with the Ledger device. In order to do this `cd` to the directory with tmkms executable and run the following (path to save the file is optional):

```py
cd <your_KMS_directory>/target/release/
./tmkms softsign keygen ~/.tmkms/secret_connection.key
```

Once you are done with the connection key, proceed to the config file. It could be created anywhere. It is possible to specify the path to it during tmkms launch.

Here is an example of the `tmkms.toml` config file, working with v0.6.3 of KMS:

```py
[[chain]]
id = "<current_chain_id>"
key_format = { type = "bech32", account_key_prefix = "cyberpub", consensus_key_prefix = "cybervalconspub" }
## Validator configuration
[[validator]]
addr = "tcp://localhost:26658"
chain_id = "<current_chain_id>"
reconnect = true # true is the default
secret_key = "<path_to_secret_connection.key>"
# enable the `ledger` feature
[[providers.ledgertm]]
chain_ids = ["<current_chain_id>"]
```

## Retrieve the validator key

The last step is to retrieve the validator key that you will be using in cyberd. The Ledger device must be connected, unlocked and has the Tendermint app opened on it.  

Start `tmkms` with the following command:

```c
tmkms start -c ~/.tmkms/tmkms.toml
```

The output should look something like this:

```bash
13:28:35 [info] tmkms 0.6.3 starting up...
13:28:35 [info] [keyring:ledgertm] added consensus key cybervalconspub1zcjduepq8jv0uxx2fw4ur6gj2r3wgs374n6ys5edh9pc4rseqqcaq2yyzy2q0fhx5q
13:28:35 [info] KMS node ID: C0E0DBA8AA82D597E8F893993E66398DDE8F89B0
```

KMS may complain about the impossibility to connect to cyberd. That's fine, we'll fix this in the next section.

The output indicates the validator key that is linked to this particular device:

```bash
cybervalconspub1zcjduepq8jv0uxx2fw4ur6gj2r3wgs374n6ys5edh9pc4rseqqcaq2yyzy2q0fhx5q
```

Take *note* of the validator pubkey that appears on your screen. We will use it in the next section.

## Configuration of Cyberd

Before we start validating, it's necessary to enanble a communication port at `cyberd` itself.  In the config file: `<your_cyberd_location>/cyberd/config/config.toml`, modify the `priv_validator_laddr` value to create a listening address/port. The port should be set according with the `[[validator]] addr =` value of your `tmkms.toml` file.

Example of `config.toml`:

```py
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "tcp://127.0.0.1:26658"
```

*NOTE: if you will try to launch cyberd after enabling  `priv_validator_laddr`, and without the `tmkms` app running simultaneously, it will likely NOT start.*

To process further, we assume that `cyberd` node is synced on your machine and there were no validators created previously on that particular node. Also, you should have an imported account to `cyberdcli` that has enough bandwidth and EUL to start validator.

We also assume that the Ledger device is connected, unlocked and has the Tendermint app open.

One last thing: `tmkms` should be launched under a different console session on the same machine.

## Run Validator

TO DO

## Maintenance of validator

TO DO
