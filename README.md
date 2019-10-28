
# Cyberd - Knowledge consensus computer for The Great Web

<div align="center">
 <img src="docs/img/header.png" width="320" />
</div>

### Code
[![version](https://img.shields.io/github/release/cybercongress/cyberd.svg?style=flat-square)](https://github.com/cybercongress/cyberd/releases/latest)
[![CircleCI](https://img.shields.io/circleci/project/github/cybercongress/cyberd.svg?style=flat-square)](https://circleci.com/gh/cybercongress/cyberd/tree/master)
[![license](https://img.shields.io/badge/License-Cyber-brightgreen.svg?style=flat-square)](https://github.com/cybercongress/cyberd/blob/master/LICENSE)
![Cosmos-SDK](https://img.shields.io/static/v1.svg?label=cosmos-sdk&message=0.35.0&color=blue&style=flat-square)
![Tendermint](https://img.shields.io/static/v1.svg?label=tendermint&message=0.31.5&color=blue&style=flat-square)
[![LoC](https://tokei.rs/b1/github/cybercongress/cyberd?style=flat)](https://github.com/cybercongress/cyberd)
[![contributors](https://img.shields.io/github/contributors/cybercongress/cyberd.svg?style=flat-square)](https://github.com/cybercongress/cyberd/graphs/contributors)
[![Coverage Status](https://img.shields.io/coveralls/github/cybercongress/cyberd/master?style=flat-square)](https://coveralls.io/github/cybercongress/cyberd?branch=master)

### Blockchain
[![chain](https://img.shields.io/badge/Chain-Euler--4-success.svg?style=flat-square)](https://github.com/cybercongress/cyberd/blob/master/docs/run_validator.md)
[![block](https://img.shields.io/badge/dynamic/json?color=blue&label=Block%20Height&query=%24.result.height&url=http%3A%2F%2F93.125.26.210%3A34657%2Findex_stats&style=flat-square)]()
[![cyberlinks](https://img.shields.io/badge/dynamic/json?color=blue&label=Cyberlinks&query=%24.result.linksCount&url=http%3A%2F%2F93.125.26.210%3A34657%2Findex_stats&style=flat-square)]()
[![cids](https://img.shields.io/badge/dynamic/json?color=blue&label=CIDs&query=%24.result.cidsCount&url=http%3A%2F%2F93.125.26.210%3A34657%2Findex_stats&style=flat-square)]()
[![agents](https://img.shields.io/badge/dynamic/json?color=blue&label=Web3%20Agents&query=%24.result.accsCount&url=http%3A%2F%2F93.125.26.210%3A34657%2Findex_stats&style=flat-square)]()
[![validators](https://img.shields.io/badge/dynamic/json?label=Validators&query=%24.result.validators.length&url=http%3A%2F%2F93.125.26.210%3A34657%2Fvalidators%3F&style=flat-square)]()

### Community
[![telegram](https://img.shields.io/badge/Join%20Us%20On-Telegram-2599D2.svg?style=flat-square)](https://t.me/fuckgoogle)
[![gitcoin](https://img.shields.io/badge/Join%20Us%20On-Gitcoin-2599D2.svg?style=flat-square)](https://t.me/fuckgoogle)
[![twitter](https://img.shields.io/twitter/follow/cyber_devs?label=Follow)](https://twitter.com/@cyber_devs)
[![reddit](https://img.shields.io/reddit/subreddit-subscribers/cybercongress?style=social)](https://www.reddit.com/r/cybercongress)

## Why

The Great Web is coming. New search systems will drive its growth. Google is the most powerful religion ever, and now is the time to abandon it.

<div align="center">
 <img src="./docs/img/cyber.png"/>
</div>

## What is Cyberd

Cyberd is a knowledge consensus computer or a search engine, which computes the [cyber•Rank]() like token, that is a weighted [Page Rank]() of a knowledge graph of the [Content IDentificators (CIDs)](), that are linked to each other with the help of [cyberlinks](): 

```
CID1 -----> CID2
```

Cyberlinks are committed by Web3 agents. They are links between two CIDs. In its current implementation, a CID is an IPFS hash of the CIDv0 or of the CIDv1 versions. A web3 agent can link any keyword with a CID and cyberd will wrap these keywords in an IPFS hash; and then, create a link between the two hashes, with a weight corresponding to the users' stake.

All the cyberlinks with a given weight are stored within the knowledge graph. The graph is re-computed by the validators every given number of blocks. For the calculations, we've implemented "the proof of relevance" root-hash, which is computed on the Cuda GPUs every round. As this is the best way to calculate a Merkle tree - faster. The volume is important because we must load the whole graph into the memory to do the calculations.

Cyberd is the first implementation of the cyber protocol. It is based on the [cosmos-SDK](https://github.com/cosmos/cosmos-sdk) and [tendermint BFT Consensus](https://github.com/tendermint/tendermint). 

This implementation uses a very simple bandwidth model. The main goal of the model is to reduce the daily networks growth to a given constant.

Thus, here we introduce [resource credits(RS)](). Each message, of a transaction type - for example, a "link" or a "send" have been assigned an RS cost. We call this "Bandwidth". A users bandwidth depends on its balance and is equal to the sum of their liquid and staked tokens. The users' bandwidth is a recoverable value. Full recovery of the bandwidths quantity, from 0 to maximum value - takes 24 hours.

There is a period `AdjustPricePeriod`, summing how much RS or bandwidth was spent during that period (`AdjustPricePeriodTotalSpent`). Also, there is a constant `AdjustPricePeriodDesiredSpent`, which is used to calculate the network load. The `AdjustPricePeriodTotalSpent/AdjustPricePeriodDesiredSpent` ratio defines the so-called current `price multiplier`. If the network usage is low, the `price multiplier` adjusts the message cost (simply by multiplying) to allow a user with a lower stake to make more transactions. If the demand for resources increases, the `price multiplier` goes `>1` thus, increasing message cost and limiting final tx count for a long-term period (RC recovery will be `<` then RC spending).


## For validators

Each validator participates in the tendermint consensus and computes/validates cyber•Rank within the knowledge graph.


## For rank providers

Rank providers crawls/indexes the web, and then cyberlinks CIDs of any given data to cyberd by consuming its bandwidth.


## For search users

A valuable, censorship-resistant and a provable search of the web for any kind of species. A search is transactionally-based, and only possible if the agent has enough bandwidth.


## For developers

The chance to create a new and a decentralized Google with affiliated services like: SEO, crawlers, web indexers, decentralized platforms and so on. You can be the first to do so. 

## For data/content producers

The opportunity to move their content to web3 and save it from any type of censorship. Your content is yours. Make sure others will be able to see it.

## For miners / GPUs holders

With the growth of the network, we will need cards, a lot of cards. Join us.


## cyber•Rank
```
0. cyber•Rank - a token weighted Page Rank (initial implementation).
1. A knowledge graph consists of CIDs, which are connected with cyberlinks. 
2. A cyberlink may be cast once only; between any given CIDs. 
3. The weight of cyberlink is the token balance of the agent, which cyber linked given CIDs. 
4. A rank computes within the current ranks calculation of the networks window on CUDA kernel.
5. After cyber•Ranks computation, each CID take a given rank.
6. [Very important] A rank computation is based on the current (computation window) agents balances.
```

## Bandwidth
```
0. The network has a desirable network bandwidth (max bandwidth).
1. An agent's bandwidth is proportional to the stake that he owns against the total network supply.
2. Linking and other chain operations consume bandwidth.
3. It takes 24 hours for a full recovery of an agent's bandwidth.
4. The network gives a discount of up to 100X of the operational costs for incentivized loads.
```

## Search index and proofs
```
0. A node can be launched in an "ALLOW SEARCH" mode, which allows searching with this node, within the knowledge graph.
1. A node that is in search mode, also constructs a full Merkle tree for cyberlinks and calculates link ranks.
2. Proof of rank of a given CID provides a Merkle path, which allows the client to validate the returned rank of a given CID and the existence of a given cyberlink.
3. The Merkle root of rank and the Merke roots of the link Merkle trees are used for calculating the app hash for each block, and are part of the protocol/consensus.
```

## Technologies
```
1. Tendermint
2. Cosmos-SDK
3. NVIDIA CUDA
4. IPFS
```

## Scaling
#### CUDA Kernel
```
1. Single host, single GPU <--- We are here
2. Single host, multiple GPUs <--- community-based R&D
3. Multiple hosts / Data Centers <--- 2021
```

#### Tendrmint
```
```

## Paper
```
Current state: Community Preview (Euler-5/Mainnet)
```

## Status

**WARNING**: Project status: **testnet**. We are at research state at the moment.
Read our [whitepaper](./docs/cyberd.md)

## cyberd Public Testnets

To run a full-node or a validator node on the latest public testnet of cyberd, please follow [this guide](./docs/help/run_validator.md).

## Explorers

The [cyberd](https://cyberd.ai/) explorer is based on [bigDipper](https://cosmos.bigdipper.live) by Forbole(https://www.forbole.com/)

## Game of Links

The "game of links" is a game between cyber•Congress and between Cosmos stakeholders for a place in Genesis. It should bootstrap and load the network at Euler-5 testnet. The greatest project will come on top with the significant number of followers. The game is finished if both of the following criteria are met:

**\>** 146 validators are in consensus for a period of 10k blocks
**\>** 500000 ATOM donated or 90 days have passed


#### Goals:
```
0. Intentiavize professional and long-term committed validators.
1. Bootstrap a full validator set, which will drive and grow the network, the knowledge graph, and will build a strong community and have the delegators backing it.
2. Initialize the knowledge graph with valuable knowledge domains; and develop tools for crawling and indexing specific knowledge domains
3. Start to form responsibility and activity within the governance of the network and within decision-making mechanisms.
4. Distribute tokens in the most valuable form and play the game with cyber•Congress.
```

#### Rewards for:
```
0. Summary uptime of every validator
1. A load of an agent (consumed networks bandwidth)
2. Amount of delegated tokens to the validator (Game of Stakes)
3. The relevance of links submitted (TOP-1000 CIDs)
```

## Distribution

## Docs

Explore the docs in our [knowledge base](https://cybercongress.ai/docs/cyberd/cyberd/).

## IBC

### Subgraph transffering
### CIDs Interchange
### Ranks Interchange

## Community

**\>** [devChat](https://t.me/fuckgoogle) for web3 agnets
**\>** [TG channel](https://t.me/cybercongress) with hot updates
**\>** [Twitter](https://twitter.com/cyber_devs) for updates and memes
**\>** [Steemit blog](https://steemit.com/@cybercongress)
**\>** [Own blog](https://cybercongress.ai/post/) with rss and useful articles

##  The 10 min. Development Setup

Take a part with [set up dev environment](https://cybercongress.ai/docs/cyberd/setup_dev_env/) in the "10 minutes challenge".

## Research and Development
```
- cyber•Rank scaling
— Online parametrization
— Onchain upgrades
— IBC
— Universal oracle
— WASM VM for gas
— CUDA VM for gas
— Privacy
- PoRep/PoST
- Autonomous onchain agents
```



## Let's #fuckgoogle together
```
Bring your PRs and text us to take on board.
```

## Issues

If you have any problems with, or questions, about search - please contact us via
a [GitHub issue](https://github.com/cybercongress/cyberd/issues).

## Contribute

You are invited to contribute new features, fixes, or updates - large or small; we are always thrilled to receive pull requests and do our best to process them as fast as we can. You can find detailed information in our
[contribution guide](./docs/contributing/contributing.md).

## Gitcoin program

We want to reward you for your contributions! We constantly fund our issues on [gitcoin](https://gitcoin.co/profile/cybercongress) and attach good descriptions to them, along with the current project state and along with user stories. We try to answer comments regularly in issues and in our [devChat](https://t.me/fuckgoogle).

<a href="https://gitcoin.co/explorer?q=cyberd">
 <img src="https://gitcoin.co/funding/embed?repo=https://github.com/cybercongress/cyberd">
</a>

## Team
<table>
  <tr>
    <td align="center"><a href=https://github.com/xhipster><img src="https://avatars0.githubusercontent.com/u/410789?s=400&v=4" width="100px;" alt="xhipster"/><br /><sub><b>Dima Starodubcev</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=xhipster" title="Documentation">📖</a> <a href="#maintenance-xhipster" title="Maintenance">🚧</a></td><td align="center"><a href="https://github.com/litvintech"><img src="https://avatars2.githubusercontent.com/u/1690657?v=4" width="100px;" alt="Valery Litvin"/><br /><sub><b>Valery Litvin</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=litvintech" title="Code">💻</a> <a href="#projectManagement-litvintech" title="Project Management">📆</a> <a href="https://github.com/cybercongress/cyberd/commits?author=litvintech" title="Documentation">📖</a></td><td align="center"><a href="https://github.com/SaveTheAles"><img src="https://avatars0.githubusercontent.com/u/36516972?v=4" width="100px;" alt="Ales Puchilo"/><br /><sub><b>Ales Puchilo</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=SaveTheAles" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/mrlp4"><img src="https://avatars2.githubusercontent.com/u/38989909?s=400&v=4" width="100px;" alt="Kiryl Laptanovich"/><br /><sub><b>Kiryl Laptanovich</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=mrlp4" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/npopeka"><img src="https://avatars1.githubusercontent.com/u/29427584?s=400&v=4" width="100px;" alt="Nick_Popeka"/><br /><sub><b>Nick Popeka</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=npopeka" title="Documentation">📖</a></td>
  </tr>
</table>


## Linked Projects

**\>** Cosmos-SDK
**\>** Tendermint
**\>** IPFS

## GIFs

## License

Cyber License - Don’t believe, don’t fear, don’t ask.

We will be happy if you fork and launch your own network and set up a knowledge graph. Eventually, we will meet each other with the help of IBC.


## The End of Google

<div align="center">
 <img src="docs/img/End-of-Google.jpg" width="600" />
</div>

## Changelog

Stay tuned with our [Changelog](./CHANGELOG.md).

<div align="center">
 <sub>Built by
 <a href="https://twitter.com/cyber_devs">cyber•Congress</a> and
 <a href="https://github.com/cybercongress/cyberd/graphs/contributors">contributors</a>
</div>
