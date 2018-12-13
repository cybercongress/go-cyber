# Before run new testnet check:

1. loose coins total amount  = initial supply
2. adjust genesis supply in x/mint
3. cyberd repo change chain-id
4. reset old state
5. change genesis file chain-id
6. generate new first validator keys
7. upadte genesis gen tx

```bash
./daemon gentx --amount=100000CBD \
 --pubkey=cbdvalconspub1zcjduepqe2wacj36s63tmytk8v4drpc4wrh7uex692msel2pjegeseeapp0q59t5pz \
 --name=euler-dev1_earth
```