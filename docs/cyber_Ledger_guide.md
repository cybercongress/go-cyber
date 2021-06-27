# Ledger nano Support

It is possible to use your Ledger device with cyber to store keys and sign transactions.

## Cyberd CLI & Ledger Nano

How to get started. First of all, you'll need a couple of things to be done:

+ A running and synced cyber node (how to: [here](https://github.com/cybercongress/go-cyber/blob/bostrom-dev/docs/run_validator.md) and [here](https://github.com/cybercongress/cyber/blob/master/docs/ultimate-commands-guide_v2.md))

+ [Setup](https://support.ledger.com/hc/en-us/articles/360000613793-Set-up-as-new-device) your Ledger device and [install Cosmos app on it](https://github.com/cosmos/ledger-cosmos/blob/master/README.md#installing) (the latest firmware for Ledger and Cosmos app required)

It is necessary to verify that cyber is built with netgo and Ledger tags. To check that, we can run: `cyber version --long`.

## Add your Ledger key

If you have set up your Ledger device on a different machine then the one running cyber, it is necessary to make sure that the Ledger device is generally working on this machine. A great way to do so is installing [Ledger Live](https://shop.ledger.com/pages/ledger-live) on the machine and trying to connect your Ledger device to it. This will show possible issues and error codes to work with ([Fix connection issues](https://support.ledger.com/hc/en-us/articles/115005165269-Fix-connection-issues) guide from Ledger).
When you made sure that your Ledger device is successfully interacting with your machine do the following:

+ Connect and unlock your Ledger device
+ Open the Cosmos app on your Ledger
+ Create an account in cyber from your Ledger key

For account creation run:

``` js
cyber keys add <your_key_name> --ledger
```

After submitting this command your Ledger device should show a generated address and will wait for confirmation. Hit confirm the button and in the console, you'll see the following output:

``` js
- name: <your_key_name>
  type: ledger
  address: cyber1gw5kdey7fs9wvh05w66s034s24tjdvxcp5fhkz
  pubkey: cyberpub1addwnpepq0lfpdumac47nyel06u95czd4026ahzmjr8stsx4h65kw3dhh60py0m7k6r
  mnemonic: ""
  threshold: 0
  pubkeys: []
  ```

By default, the `...keys add` command with account and index set to 0 of [bip44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) derivation path is used. To add more than one key account and/or index it must be specified separately in the following way:

``` js
cyber keys add <your_key2_name> --ledger --account 1 --index 1
```

You don't need to remember which numbers for account and index you've used, it will be matched to <your_key_name> automatically.

## Confirm your address

To make sure you have added everything correctly just run:

``` js
cyber keys show <key_name> -d
```

It's necessary to confirm that the key on your Ledger matches the one shown in the console.

## Signing transactions

You are now ready to sign and send transactions. This could be done by using the `tx bank send` command. Your Ledger device should be connected and unlocked at this step. Run the following to send some CYB tokens:

``` js
cyber tx bank send <from_key_name> <destination_address> <ammount>cyb --chain-id <current_chain_id>
```

`<from_key_name>` is your ledger key name, `<destination_address>` is the address of the recipient in the following format: `cyber1wq7p5qfygxr37vqqufhj5fzwlg55zmm4w0p8sw`.
When prompted with `confirm transaction before signing`, answer Y. Your Ledger will ask to approve the transaction. Make sure you'll inspect the transaction JSON before signing it. When the transaction is signed on the Ledger, usually, the output will show up in the console.
