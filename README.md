
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

The Great Web is coming. New search systems will drive their growth. Google is the most powerful religion ever and now time to go out of that.

<div align="center">
 <img src="./docs/img/cyber.png"/>
</div>

## What is Cyberd

Cyberd is knowledge consensus computer or search engine which computes cyber•Rank like token weighted Page Rank of knowledge graph of Content IDentificators (CIDs) linked with each other with cyberlinks: 

```
CID1 -----> CID2
```

The cyberlinks committed by Web3 agents. It is link between two CIDs. In current implementation the CID is IPFS hash CIDv0 or CIDv1 version. Web3 agent can link any keyword with CID and cyberd will wrap keyword in IPFS hash and create link between two hashes with a weight according to user stake.

All cyberlinks with weight store in the knowledge graph. It's re-computed by validators every number of blocks. For these calculations, we've implemented the proof of relevance root hash is computed on Cuda GPUs every round as the best way to calculate merkle tree faster. We need to load the whole graph in memory for calculating that’s why memory volume is important.

Cyberd is the first implementation of the cyber protocol based on cosmos-SDK and tendermint BFT Consensus. 

This implementation use a very simple bandwidth model. Main goal of that model is to reduce daily network growth to given constant.

Thus, here we introduce resource credits(RS). Each message type of transaction like the link or send have assigned RS cost. We call it Bandwidth. The user's bandwidth depends on its balance as the sum of liquid and staked tokens. The user's bandwidth is recoverable value. Full recovery of bandwidth quantity from 0 to max value is 24 h.

There is period `AdjustPricePeriod` summing how much RS or bandwidth was spent for that period (`AdjustPricePeriodTotalSpent`). Also, there is constant `AdjustPricePeriodDesiredSpent`, used to calculate network loading. The `AdjustPricePeriodTotalSpent/AdjustPricePeriodDesiredSpent` ratio defined so called current `price multiplier`. If network usage is low, `price multiplier` adjust message cost (by simply multiplying) to allow user with lower stake to do more transactions. If resource demand increase, `price multiplier` goes `>1` thus increase messages cost and limiting final tx count for some long-term period (RC recovery will be `<` then RC spending).


## For validators

Each validator participates in tendermint consensus and compute/validate cyber•Rank of the knowledge graph.


## For rank providers

Rank providers crawlers/index web and cyberlink CIDs of given data to cyberd consuming their bandwidth.


## For search users

Valuable, censorship-resistant and provable search in web for all species. Search is transaction-based and allowed if the agent has enough bandwidth.


## For developers

The chance to create new decentralized Google with related services as SEO, crawlers, indexers, decentralized platforms and so on. You can be the first. 

## For data/content producers

The opportunity to move their content to web3 and save it from any kind of censorship. Your content is yours. Make shure user will find it.

## For miners / GPUs holders

With network grow we will need cards, a lot of cards. Come to us, guys.


## cyber•Rank
```
0. cyber•Rank - token weighted Page Rank (first implementation)
1. Knowledge graph consists of CIDs which connected with cyberlinks. 
2. Cyberlink may be cast only one time between given CIDs. 
3. Weight of cyberlink is the token balance of agent which cyber linked given CIDs. 
4. Rank computes with current network rank's calculation window on CUDA kernel.
5. After cyber•Rank's computation, each CID take given rank.
6. [Very important] Rank computes based on current (computation window) agent balances.
```

## Bandwidth
```
0. The network has desirable network bandwidth (max bandwidth)
1. Agent's bandwidth is proportional to the stake he has to whole network's supply.
2. Linking and other chain's operations consume bandwidth.
3. It takes 24 hours for full agent's bandwidth recovery.
4. Network gives a discount to incentivize load up to 100X for operation's cost.
```

## Search index and proofs
```
0. A node may be run in ALLOW SEARCH mode which allows search with this node in the knowledge graph.
1. A node in the search mode also constructs full Merkle trees for cyberlinks and calculated ranks of links.
2. Proof of rank of given CID provide Merkle path which allows the client to validate returned rank of given CID and existence of given cyberlink.
3. Merkle roots of rank's and link's Merkle trees used for calculating of app hash for each block and are part of protocol/consensus.
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
Read [whitepaper](./docs/cyberd.md)

## cyberd Public Testnets

To run a full-node or validator in the latest public testnet of the cyberd follow [the guide](./docs/help/run_validator.md).

## Explorers

[Cyberd](https://cyberd.ai/) explorer based on [bigDipper](https://cosmos.bigdipper.live) by Forbole(https://www.forbole.com/)

## Game of Links

Game of links is a game between cyber•Congress and Cosmos stakehodlers for the place in Genesis. It should bootstrap and load network at Euler-5 testnet. The great project will come with significant number of followers. The game is finished if both of the following criteria met:

**\>** 146 validators in consensus during 10k blocks
**\>** 500000 ATOM donated or 90 days passed


#### Goals:
```
0. Intentiavize professional and long-term committed validators
1. Bootstrap full validator set which will drive and grow the network and knowledge graph and have strong community and delegators behind
2. Initialize knowledge graph with valuable knowledge domains and develop tooling for crawling and indexing specific knowledge domains
3. Start to form responsibility and activity with the network’s governance and decisions
4. Distribute tokens the most valuable way and play the game with cyber•Congress on this
```

#### Rewards for:
```
0. Summary uptime of every validator
1. A load of an agent (consumed network's bandwidth)
2. Amount of delegated to the validator (Game of Stakes)
3. The relevance of links submitted (TOP-1000 CIDs)
```

## Distribution

## Docs

Explore the docs in our [knowledge base](https://cybercongress.ai/docs/cyberd/cyberd/).

## Community

**\>** [devChat](https://t.me/fuckgoogle) for web3 agnets
**\>** [TG channel](https://t.me/cybercongress) with hot updates
**\>** [Twitter](https://twitter.com/cyber_devs) for updates and memes
**\>** [Steemit blog](https://steemit.com/@cybercongress)
**\>** [Own blog](https://cybercongress.ai/post/) with rss and useful articles

## 10 min Development Setup

Take a part with [set up dev environment](https://cybercongress.ai/docs/cyberd/setup_dev_env/) in 10 minutes challenge

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
```



## Let's #fuckgoogle together
```
Bring your PRs and text us to take on board.
```

## Issues

If you have any problems with or questions about search, please contact us through a
 [GitHub issue](https://github.com/cybercongress/cyberd/issues).

## Contribute

You are invited to contribute new features, fixes, or updates, large or small; We are always thrilled to receive pull requests and do our best to process them as fast as we can. You can find detailed information in our
 [contribution guide](./docs/contributing/contributing.md).

## Gitcoin program

We want to pay you for your contribution! We constantly fund our issues on [gitcoin](https://gitcoin.co/profile/cybercongress) and attach good description for them with project state and user stories. We try to answer to comments regular in issues and in our [devChat](https://t.me/fuckgoogle).

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

We will happy if you fork and launch your own network. And setup knowledge graph. These guys will meet each other with IBC.


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

