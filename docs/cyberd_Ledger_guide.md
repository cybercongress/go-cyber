# Ledger nano Support

It is possible to use Ledger nano devices with cyberd to store keys and sign transacions.

## Cyberd CLI & Ledger Nano

How to get started. First of all you'll need couple things to be done:

+ Working and synced cyberd node with  (how to: [here](https://github.com/cybercongress/cyberd/blob/0.1.5/docs/run_validator.md) and [here](https://github.com/cybercongress/cyberd/blob/master/docs/ultimate-commands-guide_v2.md))

+ [Setup](https://support.ledger.com/hc/en-us/articles/360000613793-Set-up-as-new-device) your Ledger device and [install Cosmos app onto it](https://github.com/cosmos/ledger-cosmos/blob/master/README.md#installing) (latest firmware for Ledger and Cosmos app v1.5.3 required)

It is necessary to verify that cyberd is build with netgo and ledger tags. To check that we can run `cyberdcli version --long`.

## Add your Ledger key

If you setup Ledger on the different machine that the one with cyberd it is necessary to make sure that Ledger device is generelly working on it. Great way to do so is installing [Ledger Live](https://shop.ledger.com/pages/ledger-live) to that particular machine and trying to connect Ledger device to it, it wil show some possible issues and error codes to work with ([Fix connection issues](https://support.ledger.com/hc/en-us/articles/115005165269-Fix-connection-issues) guide from Ledger)
When you make sure that your Ledger device is successfully interacting with your machine do following:

+ Connect and unlock your Ledger device
+ Open the Cosmos app on your Ledger
+ Create an account in cyberdcli from your ledger key

For actual account creation run:

``` js
cyberdcli keys add <your_key_name> --ledger
```

After submitting this command Ledger device should show generated address and will wait for confirmation. Hit confirm button and in the console you'll see following output:

``` js
- name: <your_key_name>
  type: ledger
  address: cyber1gw5kdey7fs9wdh05w66s0h4s24tjdvtcp5fhky
  pubkey: cyberpub1addwnpepq0lfpdumac47nysl06u95czd4026ahzmjr9stsx4h65kw3dhh60py0m7k6e
  mnemonic: ""
  threshold: 0
  pubkeys: []
  ```

By default, using `...keys add` command account 0 and index 0 of [bip44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) derivation path is used, so in order to add more than one key account and/or index must be specified separately in the following way:

``` js
cyberdcli keys add <your_key2_name> --ledger --account 1 --index 1
```

You don't need to remember which numbers for account and index you used, it will be matched to <your_key_name> automatically.

## Confirm your address

To make sure you have everything added correctly just run:

``` js
cyberdcli keys show <key_name> -d
```

Now it's necessary to confirm that key on Ledger matches one shown in console.

## Signing transactions

You now ready to sign and send transaktions. That could be done using `tx send` command. Ledger device should be connected and unlocked at this step. Run the following to send some cyb's to someone:

``` js
cyberdcli tx send <from_key_name> <destination_address> <ammount>cyb --chain-id <current_chain_id>
```

`<from_key_name>` is your ledger key name, `<destination_address>` is address of recipient in format `cyber1wq7p5qfygxr37vqqufhj5fzwlg55zmm4w0p8sw`.
When prompted with `confirm transaction before signing`, answer Y, and next your Ledger will show and ask to approve transaction. Make sure you'll inspect transaction JSON before signing. When transaction is signed on the Ledger, usuall output will show up in console.
