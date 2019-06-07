# 2 of 3 multisig account creation and sending transaction guide

`Cyberd` uses docker container technology for usability. If you don't use docker container and use `gaiacli` or you've installed `cyberd` from binaries this guide is useful for you too. Just skip some docker features, because this guide focused on docker users. Remind: this guide covers all types of transactions, not only send. Also, this guide actual for Cosmos Hub Gaiacli users excepted bandwidth in Cosmos we pay a fee with tokens.

Do not forget about `--chain-id` flag in `cyberd` and even `Cosmos Hub` networks. Current `<chain-id>` you can always get in master branch of product repo.

Multisig account creating and sending transaction is simple and clear but a little bit long.

1. Go inside docker container:

  1.1 Detect `<container_id>`
      ```bash
      docker ps
      ```
  1.2 ... and go inside it
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
Now we have 3 accounts for multisig account generating: `test1` and `test2` on the local machine and we have access to them. `test3` from remote thresholder and we haven't access to it. All created and imported accounts you can check with:
```bash
cyberdcli keys list
```

4. Now we can create test 2-of-3 multisig account named, for example, `multitest1` with keys `test1`,`test2` on local machine and `test3` by remote thresholder:
```bash
cyberdcli keys add multitest1 --multisig=test1,test2,test3 --multisig-threshold 2
```

5. You should top up your balance of your multisig account. Make sure if you have enough bandwidth to make transaction later.

6. Create unsigned transaction from multisig account and store it in `unsigned.json` file:
```bash
cyberdcli tx send <recipient_address> <amount>cyb --from=<multisig_address> --chain-id=<chain_id> --generate-only > unsigned.json
```

7. Sign this transaction with the following command and store signed file in `sign1.json`:
```bash
cyberdcli tx sign unsigned.json --multisig=<multisig_address> --from=<your_account_name> --output-document sign1.json --chain-id=<chain_id>
```

8. Now you need to send the resulting file to remote thresholders for signing. You can see the content of the transaction file with
 ```bash
cat unsigned.json
```
command and copy content to convenient for you `.json` file and send it. Also, you can copy this file from docker container to local machine by following command
```bash
docker cp <container_id>:/unsigned.json .
```
File will been copied to current repo.

9. Remote thresholder should to sign it too like it was two steps below and send you signed file back. For example `sign2.json`


10. Copy signed file from remote thresholder in a docker container by the following command:

```bash
docker cp sign2.json <container_id>:/sign2.json
```

Your docker container should content 3 `.json` files: `unsigned.json`, `sign1.json`, and `sign2.json` at least. This is necessary and sufficient condition because we've set up 2 of 3 multisig account

11. Go bask inside a docker container and generate multisig transaction with all signs.

```bash
cyberdcli tx multisign unsigned.json multitest1 sign1.json sign2.json --chain-id=<chain_id> > signed.json
```

12. Finally we need to broadcast this transaction to network

```bash
cyberdcli tx broadcast signed.json --chain-id=<chain_id>
```

If multisig account has enough bandwidth transaction should be broadcasted.
