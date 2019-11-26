# Validator launch with Tendermint KMS + Leger Nano

In this guide you'll find complete path how to setup validator at cyberd using Tenderming KMS app and Ledger nano S as a keystore.

## Preparing Ledger

First of all we'll need to prepare Ledger device to be able to work with Tendermint KMS.
We need to install Tendermint app to Ledger using Ledger Live.

*Note: at the moment, you might need to enable `developer mode` in Ledger Live settings*

Also if you setup Ledger on the different machine that the one with cyberd it is necessary to make sure that Ledger device is generally recognized. Great way to do that - installing [Ledger Live](https://shop.ledger.com/pages/ledger-live) to that particular machine and attempting to connect Ledger device to it. That wil show some possible issues and error codes (if they present, of cource) to deal with ([Fix connection issues](https://support.ledger.com/hc/en-us/articles/115005165269-Fix-connection-issues) guide from Ledger).

## Install Tendermint KMS to the node

Normally it has to be done according Tendermint [guide](https://github.com/tendermint/kms), but we need some additives above it (all instructions tested on KMS v0.6.3).

### Installation

You will need the following prerequisites:

- **Rust** (stable; 1.35+) [install](https://rustup.rs/)
- **C compiler**: e.g. gcc, clang
- **pkg-config**
- **libusb** (1.0+) : for Debian/Ubuntu: `apt install libusb-1.0-0-dev`

NOTE (x86_64 only): Configure `RUSTFLAGS` environment variable:
`export RUSTFLAGS=-Ctarget-feature=+aes,+ssse3`  (generally it would be necessary to add those flags to ~/.bash-profile or ~/.profile to make them available after re-login).

Now we are ready to install KMS. Generally there is 2 ways to do it: compite from source or install with Rust's `cargo-install`. We'll use first option.

### Compiling from source code

`tmkms` can be compiled directly from the git repository source code using the following commands:

``` js
git clone https://github.com/tendermint/kms.git && cd kms
cargo build --release --features=ledgertm,softsign
```

If successful, it will produce the `tmkms` executable located at
`./target/release/tmkms`.

## KMS configuration

After compiling, we should create settings file `tmkms.toml` and file `secret_connection.key` and adjust cyberd node settings.
First of all we need to generate connection key to be used during communication with Ledger device. In oreder to do that `cd` to the directory with tmkms executable and run following (path to save file is optional):

```py
tmkms softsign keygen ~/.tmkms/secret_connection.key
```

Once done with connection key, we can proceed to config file. It could be created anywhere, it is possible to specify path to it during tmkms launch.
Here is example of config `tmkms.toml` file wotking with v0.6.3 of KMS:

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

## Retrieve validator key

The last step is to retrieve the validator key that you will use in cyberd. Ledger device must be connected, unlocked and had Tendermint app opened on it.  Start `tmkms` with following command:

```c
tmkms start -c ~/.tmkms/tmkms.toml
```

The output should look something like that:

```bash
13:28:35 [info] tmkms 0.6.3 starting up...
13:28:35 [info] [keyring:ledgertm] added consensus key cybervalconspub1zcjduepq8jv0uxx2fw4ur6gj2r3wgs374n6ys5edh9pc4rseqqcaq2yyzy2q0fhx5q
13:28:35 [info] KMS node ID: C0E0DBA8AA82D597E8F893993E66398DDE8F89B0
```

KMS may complain about impossibillity to connect to cyberd, that's fine, we'll fix that in the next section.

Output indicates the validator key linked to this particular device is: cybervalconspub1zcjduepq8jv0uxx2fw4ur6gj2r3wgs374n6ys5edh9pc4rseqqcaq2yyzy2q0fhx5q

Take *note* of the validator pubkey that appears on your screen. We will use it in the next section.

## Configuration of Cyberd

Before we start validating, it's necessary to enanble communication port at `cyberd` itself.  In config file `<your_cyberd_location>/cyberd/config/config.toml`, modify `priv_validator_laddr` value to create a listening address/port in. Port here should be set according with `[[validator]] addr = ` value of your `tmkms.toml` file.
Example of `config.toml`:

```py
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "tcp://127.0.0.1:26658"
```

*NOTE: if you will try to launch cyberd after enabling  `priv_validator_laddr` and without `tmkms` app running simultaneously, it likely would not start.*

To process further we assume that `cyberd` is synced and there was no validators created on that particular node, as well as we have an account imported to `cyberdcli`, that has enough  bandwith and EUL to start validator.

Also we assume that Ledger device is connected, unlocked and has Tendermint open on it. And the last thing: `tmkms` should be launched under the different console session on that machine.
