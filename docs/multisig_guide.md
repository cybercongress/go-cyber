# A guide for 2 of 3 multisig account creation and for sending transactions

`Cyberd` uses docker container technology for usability.
If you are not using a docker container, and are using a `gaiacli`, or you have installed `cyberd` from binaries, then this guide can be also useful for you. Just skip some of the docker features, because this guide is focused more on the docker users. 
A reminder: this guide covers all types of transactions, not only send transactions. This guide is also actual for Cosmos Hub Gaiacli users excepted bandwidth in Cosmos we pay a fee with tokens.

Do not forget about the `--chain-id` flag in `cyberd`, and in the `Cosmos Hub` networks. 
You can always het the current `<chain-id>` in the master branch of the product repository.

## Creating a multisig

Multisig account creation and sending transactions is simple and clear, but can be a little long.

1. Go inside the docker container:

* 1.1 Detect `<container_id>`

```bash
docker ps
```

* 1.2 ... and go inside it

```bash
docker exec -ti <container_id> bash
```

2. Import or create thresholders accounts for multisig:

```bash
cyberdcli keys add test1
cyberdcli keys add test2
```

3. Add pubkeys of remote thresholders accounts:

```bash
cyberdcli keys add test3 --pubkey=<thresholder_pub_key>
```

We now have 3 accounts for multisig account generating:
`test1` and `test2` on the local machine that we have access to.
`test3` from a remote thresholder that we do not have access to.
All created and imported accounts can be checked with:

```bash
cyberdcli keys list
```

4. Now we can create and test 2-of-3 multisig account named for example: `multitest1` with keys `test1`,`test2` on local machine and `test3` on a remote thresholder:

```bash
cyberdcli keys add multitest1 --multisig=test1,test2,test3 --multisig-threshold 2
```

5. You should top up the balance of your multisig account. Make sure that you have enough bandwidth in order to execute transactions later.

## Spending from a multisig account

6. Create an unsigned transaction from the multisig account and store it in the `unsigned.json` file:

```bash
cyberdcli tx send <recipient_address> <amount>cyb --from=<multisig_address> --chain-id=<chain_id> --generate-only > unsigned.json
```

7. Sign this transaction with the following command, and then store the signed file in `sign1.json`:

```bash
cyberdcli tx sign unsigned.json --multisig=<multisig_address> --from=<your_account_name> --output-document sign1.json --chain-id=<chain_id>
```

8. You now need to send the obtained file to a remote thresholders for signing. You can see the content of the file containing the transaction with:

 ```bash
cat unsigned.json
```

command, and you may now copy the content that is convenient to your `.json` file and send it.
Also, you can copy this file from the docker container to a local machine via the following command:

```bash
docker cp <container_id>:/unsigned.json .
```

The file has been copied to current repo.

9. You should aslo sign the remote thresholder like you did two steps above, and send your signed file back.
For example `sign2.json`

10. Copy the signed file from the remote thresholder in a docker container via the following command:

```bash
docker cp sign2.json <container_id>:/sign2.json
```

Your docker container should content 3 `.json` files: 
`unsigned.json`, `sign1.json`, and `sign2.json` (at least). This is the necessary and sufficient conditions because we've set up a     2-out-of 3 multisig account.

11. Go bask inside the docker container and generate a multisig transaction with all the signatures:

```bash
cyberdcli tx multisign unsigned.json multitest1 sign1.json sign2.json --chain-id=<chain_id> > signed.json
```

12. Finally we need to broadcast this transaction to the network:

```bash
cyberdcli tx broadcast signed.json --chain-id=<chain_id>
```

If the multisig account has enough bandwidth, the transaction should be broadcasted to the network.
