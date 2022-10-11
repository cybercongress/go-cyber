# A guide for creating a 2 of 3 multisig account and sending transactions

To follow this guide you'll need `pussy` installed and connected to any pussy node (refer to our cli [guide](https://github.com/joinresistance/space-pussy/blob/main/docs/ultimate-commands-guide.md).
A reminder: this guide covers all types of transactions, not only send transactions. This guide is also relevant for Cosmos Hub Gaiacli users, except for the bandwidth params, in Cosmos we pay a fee using tokens.

Do not forget about the `--chain-id` flag in `pussy`, and in the `Cosmos Hub` networks.
You can always get the current `<chain-id>` in the master branch of the [repository](https://github.com/joinresistance/space-pussy).

## Creating a multisig

The multisig account creation and sending transactions are simple and clear but can be a little long.

1. Import or create a thresholder accounts for multisig:

```bash
pussy keys add test1
pussy keys add test2
```

2. Add pubkeys of remote thresholder accounts:

```bash
pussy keys add test3 --pubkey=<thresholder_pub_key>
```

We now have 3 accounts for multisig account generating:
`test1` and `test2` on a local machine that we have access to.
`test3` from a remote thresholder that we do not have access to.
All the created and imported accounts can be checked with:

```bash
pussy keys list
```

3. Now, we can create and test the 2-of-3 multisig account, named for example: `multitest1` with keys `test1`,`test2` on a local machine and `test3` on a remote thresholder:

```bash
pussy keys add multitest1 --multisig=test1,test2,test3 --multisig-threshold 2
```

4. You should top up the balance of your multisig account. Make sure that you have enough bandwidth to execute transactions later.

## Spending out of a multisig account

5. Create an unsigned transaction from the multisig account and store it in the `unsigned.json` file:

```bash
pussy tx send <recipient_address> <amount>pussy --from=<multisig_address> --chain-id=<chain_id> --generate-only > unsigned.json
```

6. Sign this transaction with the following command and then store the signed file in `sign1.json`:

```bash
pussy tx sign unsigned.json --multisig=<multisig_address> --from=<your_account_name> --output-document sign1.json --chain-id=<chain_id>
```

7. You need to send the obtained file to a remote thresholders for signing. You can see the content of the file containing the transaction with:

 ```bash
cat unsigned.json
```

You may now copy the content that is convenient to your `.json` file and send it.

8. You should also sign the remote thresholder, just like you did two steps above, and send your signed file back.
For example `sign2.json`

9. Copy the signed file from the remote thresholder into your cli home directory with the following command:

```bash
cp sign2.json $HOME/.pussy
```

Your cli-home folder should content 3 `.json` files:
`unsigned.json`, `sign1.json`, and `sign2.json` (at least). Those are the necessary and sufficient conditions, because we've set up a 2-out-of 3 multisig account.

10. Generate a multisig transaction with all of the signatures:

```bash
pussy tx multisign unsigned.json multitest1 sign1.json sign2.json --chain-id=<chain_id> > signed.json
```

11. Finally, we need to broadcast this transaction to the network:

```bash
pussy tx broadcast signed.json --chain-id=<chain_id>
```

If the multisig account has enough bandwidth, the transaction should be broadcasted to the network.
