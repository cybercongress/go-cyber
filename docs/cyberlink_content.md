# Cyberlink content with Cyber-js

## Script preparation

A small, ready-made repository exists so you can experiment with Cyber-js. Clone it from here. You need [NodeJs](https://nodejs.org/en/download/). If you open the folder in [Visual Studio Code](https://code.visualstudio.com/Download), the IDE should give you all the coding help you require. In the cloned folder you need to install the required modules:

```sh
$ npm install
```

Create a new file named `experiment.ts`. In it, put these lines to confirm it works:

```typescript
const runAll = async(): Promise<void> => {
    console.log("TODO")
}

runAll()
```

To execute, this TypeScript file needs to be compiled into JavaScript before being interpreted by NodeJs. Add this as a run target in `package.json`:

```json
...
    "scripts": {
        ...
        "experiment": "ts-node experiment.ts"
    }
...
```

Confirm that it does what you want:

```sh
$ npm run experiment
```

This returns:

```
> ts-node experiment.ts

TODO
```

You will soon make this script more meaningful. With the basic script ready, you need to prepare some elements.

## Testnet preparation

The Bostrom has a number of testnets running. The Bostrom is currently running a [public testnet](https://github.com/cybercongress/cybernode) for the Space-pussy-1 upgrade that you are connecting to and running your script on. You need to connect to a public node so that you can query information and broadcast transactions. One of the available nodes is:

```
RPC: https://rpc.space-pussy-1.cybernode.ai
```

You need a wallet address on the testnet and you must create a 24-word mnemonic in order to do so. CosmJS can generate one for you. Create a new file `generate_mnemonic.ts` with the following script:

```typescript
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing"

const generateKey = async (): Promise<void> => {
    const wallet: DirectSecp256k1HdWallet = await DirectSecp256k1HdWallet.generate(24)
    process.stdout.write(wallet.mnemonic)
    const accounts = await wallet.getAccounts()
    console.error("Mnemonic with 1st account:", accounts[0].address)
}

generateKey()
```

Now create a **key** for our imaginary user **Alice**:

*Note: You likely need to update Node.js to a later version if this fails. Find a guide [here](https://phoenixnap.com/kb/update-node-js-version).

```sh
$ npx ts-node generate_mnemonic.ts > testnet.alice.mnemonic.key
```

When done, it should also tell you the address of the first account:

```
Mnemonic with 1st account: bostrom1sw8xv3mv2n4xfv6rlpzsevusyzzg78r3e78xnp
```

Temporarily keep this address for convenience, although CosmJS can always recalculate it from the mnemonic. Privately examine the file to confirm it contains your 24 words.

<HighlightBox type="warn">

Important considerations:

1. `process.stdout.write` was used to avoid any line return. Be careful not to add any empty lines or any other character in your `.key` file (this occurs with VSCode under certain conditions). If you add any characters, ComsJs may not be able to parse it.
2. Adjust the `.gitignore` file to not commit your `.key` file by mistake:

    ```
    node_modules
    *.key
    ```

</HighlightBox>

<!-- <HighlightBox type="tip">

For your convenience, we have a branch available [here](https://github.com/b9lab/cosmjs-sandbox/tree/file-preparation) that contains all the code and files you've added so far.

</HighlightBox> -->

## Add your imports

You need a small, simple interface to a blockchain, one which could eventually have users. Good practice is to refrain from requesting a user address until necessary (e.g. when a user clicks a relevant button). Therefore, in `experiment.ts` you first use the read-only client. Import it at the top of the file:

```typescript
import { CyberClient } from "@cybercongress/cyber-js"
```

<HighlightBox type="info">

Note that VSCode assists you to auto-complete [`CyberClient`](https://github.com/cybercongress/soft3.js/blob/main/src/cyberclient.ts#L144) if you type <kbd>CTRL-Space</kbd> inside the `{}` of the `import` line.

</HighlightBox>

## Define your connection

Next, you need to tell the client how to connect to the RPC port of your blockchain:

```typescript
const rpc = "https://rpc.space-pussy-1.cybernode.ai"
```

Inside the `runAll` function you [initialize the connection](https://github.com/cybercongress/soft3.js/blob/main/src/cyberclient.ts#L165) and immediately [check](https://github.com/cybercongress/soft3.js/blob/main/src/cyberclient.ts#L244) you connected to the right place:

```typescript
const runAll = async(): Promise<void> => {
    const client = await CyberClient.connect(rpc)
    console.log("With client, chain id:", await client.getChainId(), ", height:", await client.getHeight())
}
```

Run again to check with `npm run experiment`, and you get:

```
With client, chain id: space-pussy-1 , height: 9507032
```

## Prepare a signing client

If you go through the methods inside [`CyberClient`](https://github.com/cybercongress/soft3.js/blob/main/src/cyberclient.ts#L144), you see that it only contains query-type methods and none for sending transactions.

Now, for Alice to send transactions, she needs to be able to sign them. And to be able to sign transactions, she needs access to her _private keys_ or _mnemonics_. Or rather she needs a client that has access to those. That is where [`SigningCyberClient`](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L239) comes in. Conveniently, `SigningCyberClient` inherits from `CyberClient`.

Update your import line:

```typescript
import { SigningCyberClient, CyberClient } from "@cybercongress/cyber-js"
```

Look at its declaration by right-clicking on the `SigningCyberClient` in your imports and choosing <kbd>Go to Definition</kbd>.

When you instantiate `SigningCyberClient` by using the [`connectWithSigner`](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L247) method, you need to pass it a [**signer**](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L249). In this case, use the [`OfflineDirectSigner`](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L29) interface.

<HighlightBox type="info">

The recommended way to encode messages is by using `OfflineDirectSigner`, which uses Protobuf. However, hardware wallets such as Ledger do not support this and still require the legacy Amino encoder. If your app requires Amino support, you have to use the `OfflineAminoSigner`.
<br></br>
Read more about encoding [here](https://docs.cosmos.network/master/core/encoding.html).

</HighlightBox>

The signer needs access to Alice's **private key**, and there are several ways to accomplish this. In this example, use Alice's saved **mnemonic**. To load the mnemonic as text in your code you need this import:

```typescript
import { readFile } from "fs/promises"
```

There are several implementations of `OfflineDirectSigner` available. The [`DirectSecp256k1HdWallet`](https://github.com/cosmos/cosmjs/blob/0f0c9d8/packages/proto-signing/src/directsecp256k1hdwallet.ts#L133) implementation is most relevant to us due to its [`fromMnemonic`](https://github.com/cosmos/cosmjs/blob/0f0c9d8/packages/proto-signing/src/directsecp256k1hdwallet.ts#L140-L141) method. Add the import:

```typescript
import { DirectSecp256k1HdWallet, OfflineDirectSigner } from "@cosmjs/proto-signing"
```

The `fromMnemonic` factory function needs a string with the mnemonic. You read this string from the mnemonic file. Create a new top-level function that returns an `OfflineDirectSigner`:

```typescript [https://github.com/b9lab/cosmjs-sandbox/blob/4168b97/experiment.ts#L9-L13]
const getAliceSignerFromMnemonic = async (): Promise<OfflineDirectSigner> => {
    return DirectSecp256k1HdWallet.fromMnemonic((await readFile("./testnet.alice.mnemonic.key")).toString(), {
        prefix: "bostrom",
    })
}

```

The Bostrom Testnet uses the `bostrom` address prefix. This is the default used by `DirectSecp256k1HdWallet`, but you are encouraged to explicitly define it as you might be working with different prefixes on different blockchains. In your `runAll` function, add:

```typescript
const aliceSigner: OfflineDirectSigner = await getAliceSignerFromMnemonic()
```

As a first step, confirm that it recovers Alice's address as expected:

```typescript
const alice = (await aliceSigner.getAccounts())[0].address
console.log("Alice's address from signer", alice)
```

Now add the line that finally creates the signing client:

```typescript
const signingClient = await SigningCyberClient.connectWithSigner(rpc, aliceSigner)
```

Check that it works like the read-only client that you used earlier, and from which [it inherits](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L239), by adding:

```typescript
console.log(
    "With signing client, chain id:",
    await signingClient.getChainId(),
    ", height:",
    await signingClient.getHeight()
)
```

## Cyberlinks

A cyberlink (noun) is a link between two particles registered in Bostrom blockchain by a particular neuron.

To cyberlink (verb) - to create a cyberlink between two particles.

## Bandwidth

The bandwidth module process and stores neuron's bandwidth in the network, dynamically adjust bandwidth price to network load. Neurons use bandwidth to add cyberlinks to the network and not need to pay gas fees.

## Bandwidth model

Bandwidth is used for billing for cyberlinks creation instead of fee billing based on gas. Using the bandwidth model removes the cognitive gap for cyberlinks creation because not force neurons to pay a fee. Holding 1 volt represents the possibility to create 1 cyberlinks per `RecoveryPeriod` blocks. (example 1 cyberlink per day holding 1 Volt) Total supply of VOLTs (volts' holdings of all neurons) it's desirable network bandwidth. By investminting of HYDROGEN to VOLT then neuron increase personal bandwidth and desirable bandwidth of network. If the network has low bandwidth consumption, then the network provides a discount using `BasePrice` multiplier for neurons.

## Personal Bandwidth

The volt's stake of the given neuron are easy to understand as the size of his battery. The creation of cyberlinks will consume battery charge, and the battery will be fully recharged during `RecoveryPeriod` blocks period. If a neuron consumes half of its bandwidth, its battery will be fully charged in the `RecoveryPeriod/2` blocks period. If a neuron act when network bandwidth consumption is low, then he will consume less personal bandwidth.

## Get Bandwidth

Alise need to [`investmint`](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L370) `HYDROGEN(stake token)` to `VOLT` and `AMPERE` of  get `Bandwidth`

Alise would like to investmint 1000000000 BOOT (1 GBOOT) to VOLT resource with lock for 30 DAYS (no spendable) and Alise will get newly minted 1 VOLT to my account locked for 30 DAYS (no spendable).

```
(1 GBOOT | 30 DAYS | VOLT) ---investmint---> locked (1 GBOOT | 30 DAYS) + minted and locked (1 VOLT | 30 DAYS)
```

Alise would like to investmint 4000000000 BOOT (4 GBOOT) to AMPERE resource with lock for 7 DAYS (no spendable) and Alise will get newly minted 1 AMPERE to my account locked for 7 DAYS (no spendable).

```
(4.2 GBOOT | 7 DAYS | AMPERE) ---investmint---> locked (4.2 GBOOT | 7 DAYS) + minted and locked (1 AMPERE | 7 DAYS)

```

## Cyberlink content

Alice can send cyberlink, but to do so she also needs to pay the network's gas fee. How much gas should she use, and at what price?

She can copy this:

```
Gas fee: [ { denom: 'boot', amount: '0' } ]
Gas limit: 200000
```

With the gas information now decided, how does Alice structure her command so that she cyberlinks content to network ? `SigningCyberClient`'s [`cyberlink`](https://github.com/cybercongress/soft3.js/blob/main/src/signingcyberclient.ts#L316) function takes a `CidFrom` and `CidTo` as input.:

```typescript
cyberlink(
    neuron: string,
    from: string,
    to: string,
    fee: StdFee,
    memo = "",
  ): Promise<DeliverTxResponse | string[]> 
```

### Upload content to IPFS

Before sending the cyberlink, Alice needs to upload the content to `ipfs`. She can use [`js-ipfs`](https://github.com/ipfs/js-ipfs/)

### Running js-IPFS in your application

If you do not need to run a command line daemon, use the ipfs-core package - it has all the features of ipfs but in a lighter package:

```sh
$ npm install ipfs-core
```

Then start a node in your app:

```js
import * as IPFS from 'ipfs-core'

const ipfs = await IPFS.create()
const { cid } = await ipfs.add('Hello world')
console.info(cid)
// QmXXY5ZxbtuYj6DnfApLiGstzPN7fvSyigrRee3hDWPCaf
```

### Running js-IPFS on the CLI

Installing `ipfs` globally will give you the `jsipfs` command which you can use to start a daemon running:

```sh
$ npm install -g ipfs
$ jsipfs daemon
Initializing IPFS daemon...
js-ipfs version: x.x.x
System version: x64/darwin
Node.js version: x.x.x
Swarm listening on /ip4/127.0
.... more output
```

You can then add a file:

```sh
$ jsipfs add ./hello-world.txt
added QmXXY5ZxbtuYj6DnfApLiGstzPN7fvSyigrRee3hDWPCaf hello-world.txt
```

Before the cyberlink, Alisa needs to check `PersonalBandwidth` for the possibility of making a cyberlink:

```typescript
const checkPersonalBandwidth = async (client: CyberClient, alice: string): Promise<boolean> => {
  try {
    const response = await client.price()
    const priceLink = response.price.dec * 10 ** -18

    const responseAccountBandwidth = await client.accountBandwidth(alice)
    const { maxValue, remainedValue } = responseAccountBandwidth.neuronBandwidth

    if (maxValue === 0 || remainedValue === 0) {
      return false
    } else if (Math.floor(remainedValue / (priceLink * 1000)) === 0) {
      return false
    }

    return true
  } catch (error) {
    console.log('error', error)
    return false
  }
}
```

With this gas and cids, add the command:

```typescript
// Execute the cyberlink Tx and store the result
const result = await signingClient.cyberlink(
    alice,
    CidFrom,
    CidTo,
    {
        amount: [{ denom: "boot", amount: "0" }],
        gas: "200000",
    },
)
// Output the result of the Tx
console.log("Transfer result:", result)
```

Run this with `npm run experiment` and you should get:

```
...
Transfer result: {
  code: 0,
  height: 0,
  rawLog: '[]',
  transactionHash: '2A2F6D0610FF60458EFC05A820B610CEFBCB5C9EF1CA322808DD8B88D369B5E0',
  gasUsed: 0,
  gasWanted: 0
}
```
Check that Alice upload information with [`getTx`](https://github.com/cybercongress/soft3.js/blob/main/src/cyberclient.ts#L314)

```typescript
const result = await client.getTx("2A2F6D0610FF60458EFC05A820B610CEFBCB5C9EF1CA322808DD8B88D369B5E0")
```
