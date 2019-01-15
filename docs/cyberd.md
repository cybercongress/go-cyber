# cyberd: Computing knowledge from web3

Notes for [euler](https://github.com/cybercongress/cyberd/releases/tag/v0.1.0) release of [cyberd](https://github.com/cybercongress/cyberd) reference implementation of `cyber://` protocol in Go

[cyber•Congress](https://cybercongress.ai/): @xhipster, @litvintech, @hleb-albau, @arturalbov

```
cyb:
- nick. a friendly software robot who helps you explore universes

cyber:
- noun. a superintelligent network computer for answers
- verb. to do something intelligent, to be very smart

cyber://
- web3 protocol for computing answers and knowledge exchange

CYB:
- ticker. transferable token expressing will to become smarter

CYBER:
- ticker. non-transferable token measuring intelligence

CBD:
- ticker. erc-20 proto token representing substance from which CYB emerge

cyberlink:
- link type. express connection from one link to another as link-x.link-y

```

## Content
<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->
<!-- code_chunk_output -->

* [cyberd: Computing knowledge from web3](#cyberd-computing-knowledge-from-web3)
	* [Content](#content)
	* [Abstract](#abstract)
	* [Introduction to web3](#introduction-to-web3)
	* [Cyber protocol at `euler`](#cyber-protocol-at-euler)
	* [Knowledge graph](#knowledge-graph)
	* [Cyberlinks](#cyberlinks)
	* [Notion of consensus computer](#notion-of-consensus-computer)
	* [Relevance machine](#relevance-machine)
	* [cyber•Rank](#cyberrank)
	* [Proof of relevance](#proof-of-relevance)
	* [Speed](#speed)
	* [Implementation in browser](#implementation-in-browser)
	* [From Inception to Genesis](#from-inception-to-genesis)
	* [Satoshi Lottery](#satoshi-lottery)
	* [Possible applications](#possible-applications)
	* [Economic protection is `smith`](#economic-protection-is-smith)
	* [Ability to evolve is `darwin`](#ability-to-evolve-is-darwin)
	* [`turing` is about computing more](#turing-is-about-computing-more)
	* [In search for equilibria is `nash`](#in-search-for-equilibria-is-nash)
	* [On faster evolution at `weiner`](#on-faster-evolution-at-weiner)
	* [Genesis is secure as `merkle`](#genesis-is-secure-as-merkle)
	* [Conclusion](#conclusion)
	* [References](#references)
	* [Notes](#notes)

<!-- /code_chunk_output -->

## Abstract

A consensus computer allow to compute provably relevant answers without opinionated blackbox intermediaries such as Google, Youtube, Amazon or Facebook. Stateless content-addressable peer-to-peer communication networks such as IPFS and stateful consensus computers such as Ethereum provide part of the solution but there are at least three problems associated with implementation. Of course, the first problem is subjective nature of relevance. The second problem is that it is hard to scale consensus computer of huge knowledge graph. The third problem is that the quality of such knowledge graph will suffer from different attack surfaces such as sybil, selfish behaviour of interacting agents and other attack vectors. In this paper we (1) define a protocol for provable consensus computing of relevance between IPFS objects based on Tendermint consensus of cyberRank computed on GPU (2) discus implementation details and (3) design distribution and incentive scheme based on our experience and known attacks. We believe the minimalistic architecture of the protocol is critical for formation of a network of domain specific search consensus computers. As result of our work some applications never existed before emerge. We expand the work including our vision on features we expect to work up to Genesis.

## Introduction to web3

Original protocols of the Internet such as TCP/IP, DNS, URL and HTTPS brought a web into the point there it is now. Along with all benefits they has created they brought more problem into the table. Globality being a key property of the the web since inception is under real threat. Speed of connections degrade with network grow and from ubiquitous government interventions into privacy and security of web users. One property, not obvious in the beginning, become really important with everyday usage of the Internet: its ability to exchange permanent hyperlinks thus they would not break after time have pass. Reliance on one at a time internet service provider architecture allow governments censor packets is the last straw in conventional web stack for every engineer who is concerned about the future of our children. Other properties while being not so critical are very desirable: offline and real-time. Average internet user being offline must have ability to work with the state it has and after acquiring connection being able to sync with global state and continue verify state's validity in realtime while having connection. Now this properties offered on app level while such properties must be integrated into lower level protocols.

The emergence of a [web3 stack](https://github.com/w3f/Web3-wiki/wiki) creates an opportunity for a new kind of Internet. We call it web3. It has a promise to remove problems of conventional protocol stack and add to the web better speed and more accessible connection. But as usually in a story with a new stack, new problems emerge. One of such problem is general purpose search. Existing general purpose search engines are restrictive centralized databases everybody forced to trust. These search engines were designed primarily for client-server architecture based on TCP/IP, DNS, URL and HTPPS protocols. Web3 create a challenge and opportunity for a search engine based on developing technologies and specifically designed for them. Surprisingly the permission-less blockchain architecture itself allows organizing general purpose search engine in a way inaccessible for previous architectures.

## Cyber protocol at `euler`

- compute `euler` inception of cyber protocol based on satoshi lottery
- def knowledge graph state
- take cyberlinks
- check that signatures are valid
- check that resource consumption of a signer is not exceed 24 moving average (bandwidth)
- if signatures and bandwidth ok than cyberlink is valid
- for every valid cyberlink emit prediction as array of cyberlinks
- every round calculate cyber•rank deltas for the knowledge graph
- every round distribute CYB based on defined rules
- apply more secure consensus state based on CBD balances
- continue up to `merkle`

## Knowledge graph

We represent a knowledge graph as weighted graph of directed links between content addresses or content identifications or CIDs. In this paper we will use them as synonyms.

> Illustration

Content addresses are essentially a web3 links. Instead of using non obvious and mutable thing:
```
https://github.com/cosmos/cosmos/blob/master/WHITEPAPER.md
```
we can use pretty much exact thing:
```
Qme4z71Zea9xaXScUi6pbsuTKCCNFp5TAv8W5tjdfH7yuHhttps
```

Using content addresses for building a knowledge graph we get [so much needed](https://steemit.com/web3/@hipster/an-idea-of-decentralized-search-for-web3-ce860d61defe5est) superpowers of [ipfs](dura://QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps.ipfs)-[like](dura://QmXHGmfo4sjdHVW2MAxczAfs44RCpSeva2an4QvkzqYgfR.ipfs) p2p protocols for a search engine:

- mesh-network future proof
- interplanetary
- tolerant
- accessible
- technology agnostic

> Illustration

Our knowledge graph is generated by web3 agents. Web3 agents includes itself to the knowledge graph transacting only once. Thereby they proof the existence of private keys for content addresses of public keys.

Our `euler` implementation is based on [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) identities and [cidv0](https://github.com/multiformats/cid#cidv0) content addresses.

Web 3 agents generate knowledge graph by applying cyberlinks.

## Cyberlinks

In order to understanding cyberlinks we need to understand the difference between url link and ipfs link. Url link points to location of content, but ipfs link point to the content itself. The difference in web architecture based on location links and content links is drastical, hence require new approaches.

Cyberlink is an approach to semantically link two content address links.

```
QmdvsvrVqdkzx8HnowpXGLi88tXZDsoNrGhGvPvHBQB6sH.QmdSQ1AGTizWjSRaVLJ8Bw9j1xi6CGLptNUcUodBwCkKNS
```

This cyberlink means that cyberd presentation on cyberc0n is referencing tezos whitepaper.

A concept of cyberlink is a convention around simple semantics of communication format in any peer to peer network:

`<content-address x>.<content-address y>`

You can see that cyberlink represents a link between two links. Easy peasy!

Cyberlink is simple yet powerful semantic construction. Cyberlinks can form link chains if exist a series of two cyberlinks from one agent in which the second link in first cyberlink is equal to first link in second cyberlink:

```
<content-address x>.<content-address y>
<content-address y>.<content-address z>
```

Using this simple principle agents can reach consensus around interpreting clauses. So link chains are helpful for interpreting rich communications around relevance.

> Illustration

Also using the following link: `QmNedUe2wktW65xXxWqcR8EWWssHVMXm3Ly4GKiRRSEBkn` the one can signal start and stop of execution in the knowledge graph.

If web3 agents instead of native ipfs links use something semantically more rich as [dura]() links than web3 agents can easier to reach consensus on the rules for program execution.

Certainly DURA protocol is a good implementation of a cyberlinks concept.

`euler` implementation of cyberlinks based on DURA specification is available in `.cyber` app of browser `cyb`.

Based on cyberlinks we can compute the relevance of subjects and objects in a knowledge graph. That is why we need a consensus computer.

## Notion of consensus computer

Consensus computer is an abstract computing machine that emerge from agents interactions.

A consensus computer has a capacity in terms of fundamental computing resources such as memory and computing. In order to interact with agents a computer need a bandwidth.

Ideal consensus computer is a computer in which:

```
sum of all *individual agents* computations and memory
is equal to
sum of all verified by agents computations and memory of a *consensus computer*
```

We know that:

```
verifications of computations < computations + verifications of computations

```
Hence we will not be able to achieve ideal consensus computer ever.

But this theory can work as a performance indicator of a consensus computer.

Our current implementation is a 64 bit consensus computer for relevance of 64 byte string space that is as far from ideal at least as 1/146.

> Illustration

We must bind computational, storage and broadband supply of relevance machine with maximised demand of queries.

Computation and storage in case of basic relevance machine can be easily predicted based on broadband, but broadband require limiting mechanism. Broadband limit can be imported from [/docs](https://github.com/cybercongress/cyberd/blob/master/docs/bandwidth.md). CYB tokens must appears here.

## Relevance machine

Relevance machine is a machine that transition knowledge graph state based on some reputation score of agents.

This machine is enable simple construction for search question querying and answers delivering.

The reputation score is projected on every agent's cyberlink. Agents abuse is prevented by a simple rule: one content address can be voted by a token only once. So it does not matter for ranking from how much accounts you voted. Only sum of their balances matters.

> Illustration

A useful property of a relevance machine is that it must have inductive reasoning property or simply follows black box principle.

```
She must be able to interfere predictions
without any knowledge about objects
except who linked, when linked and what was linked.
```

If we assume that a consensus computer must have some information about linked objects the complexity of such model growth unpredictably, hence a requirements for a computer for memory and computations. That is, deduction of a meaning inside consensus computer is expensive thus our design hardly depend on the blindness assumption. Instead we design a system in which meaning extraction is incentivezed because agents need CYB to compute relevance.

Also thanks to content addresses the relevance machine following black box principle do not need to store the data but can effectively operate on it.

Human intelligence organized in a way to prune non relevant and non important memories with time has pass. The same way can do relevance machine.

Also one useful property of relevance machine is that it doesn't need to store neither past state nor full current state to remain useful, or more precisely: "relevant".

So relevance machine can implement aggressive pruning strategies such as pruning all history of knowledge graph formation or forgetting links that become non relevant.

Forgetting links: Prune min possible rank / 2

The pruning group of features features can be implemented in `nash`.

`euler` implementation of relevance machine is based on the simplest mechanism which is called cyber•Rank.

## cyber•Rank

Ranking using consensus computer is hard because consensus computers bring serious resource bounds. e.g. [Nebulas](dura://QmefxTSFG1W95yg3PLfKV2mshh6TtyRxwv5yPiyZCGyPmG.ipfs) fail to deliver something useful onchain. First we must ask ourselves why do we need to compute and store the rank on chain, and not go Colony or Truebit way?

If rank computed inside consensus computer you have easy content distribution of the rank as well as easy way to build provable applications on top of the rank. Hence we decide to follow more cosmic architecture. In the next section we describe proof of relevance mechanism which allow network to scale with help of domain specific relevance machines that works in parallel.

Eventually relevance machine need to find (1) deterministic algorithm that allow to compute a rank for continuously appended network to scale the consensus computer to orders of magnitude that of Google. Perfect algorithm (2) must have linear memory and computation complexity. The most importantly it must have (3) highest provable prediction capabilities for existence of relevant links.

After some research we found that we can not find silver bullet here. We find an algorithm that is probably satisfy our criteria: SpringRank. Original idea of an algorithm came to Caterina from physics. Links represented as system of springs with some energy and the task of computing the ranks is the task of finding relaxed state of springs.

But we got at least 3 problems with SpringRank:
1. We were not able to implement it onchain fast using Go in `euler`.
2. We was not able to prove it for knowledge graph because we just did not have provable knowledge graph yet.
3. Also we was not able to prove it applying it for the Ethereum blockchain during computing the genesis file for `euler`. It could work, but for the time being it is better to call this kind of distribution a lottery.

So we decide to find some more basic bulletproof way to bootstrap the network: a rank from which a previous network has been bootstrapped by Lary and Sergey. The problem with original PageRank is that it is not resistant to Sybil Attack.

Token weighted [PageRank](http://ilpubs.stanford.edu:8090/422/1/1999-66.pdf) limited by token weighted bandwidth do not have inherent problems of naive PageRank and is resistant to sybil attacks. For the time being we will call it cyber•Rank until something better emerge.

In the centre of spam protection system is an assumption that write operations can be executed only by those who have vested interest in the evolutionary success of a relevance machine. Every 1% of stake in consensus computer gives the ability to use 1% of possible network broadband and computing capabilities.

As nobody uses all possessed broadband we can safely use 10x fractional reserves with 2 minute recalculation target.

In order to switch from one algorithm to another we are going to make simulations and experiment with economic a/b testing based on winning chains through hard spoons.

Consensus computer based on relevance machine for cyber•Rank is able to answer and deliver relevant results for any given search request in the 64 byte CID space. But in order to build a network of domain specific relevance machines it is not enough. Consensus computers must have ability to prove relevance for each other.

## Proof of relevance

We design a system under assumption that in terms of search such thing as bad behaviour does not exist as nothing bad can be in the intention of finding answers. Also this approach significantly reduce attack surfaces.

> Ranks is computed on the only fact that something has been searched, thus linked and as result affected predictive model.

Good analogy is observing in quantum mechanics. That is why we do not need such things as negative voting. Doing this we remove subjectivity out of the protocol and can define proof of relevance.

```
Rank state = rank values stored in one dimensional array and merkle tree of those values
```

Each new CID gets unique number. Number starts from zero and incrementing by one for each new CID. So that we can store rank in one dimensional array where indices are CID numbers.

Merkle Tree calculated based on RFC-6962 standard (https://tools.ietf.org/html/rfc6962#section-2.1). Since rank stored in one dimensional array where indices are CID numbers (we could say that it ordered by CID numbers) leafs in merkle tree from left to right are `SHA-256` hashes of rank value. Index of leaf is CID number. It helps to easily find proofs for specified CID (`log n` iterations where `n` is number of leafs).

To store merkle tree is necessary to split tree on subtrees with number of leafs multiply of power of 2. The smallest one is obviously subtree with only one leaf (and therefore `height == 0`). Leaf addition looks as follows. Each new leaf is added as subtree with `height == 0`.
Then sequentially merge subtrees with the same `height` from right to left.

Example:


      ┌──┴──┐   │      ┌──┴──┐     │         ┌──┴──┐     │   │
    ┌─┴─┐ ┌─┴─┐ │    ┌─┴─┐ ┌─┴─┐ ┌─┴─┐     ┌─┴─┐ ┌─┴─┐ ┌─┴─┐ │
    (5-leaf)         (6-leaf)              (7-leaf)

To get merkle root hash - join subtree roots from right to left.

Rank merkle tree can be stored differently:

<b>Full tree</b> - all subtrees with all leafs and intermediary nodes  
<b>Short tree</b> - contains only subtrees roots

The trick is that <b>full tree</b> is only necessary for providing merkle proofs. For consensus purposes and updating tree it's enough to have <b>short tree</b>. To store merkle tree in database use only <b>short tree</b>. Marshaling of short tree with `n` subtrees (each subtree takes 40 bytes):  

```
<subtree_1_root_hash_bytes><subtree_1_height_bytes>
....
<subtree_n_root_hash_bytes><subtree_n_height_bytes>
```

For `1,099,511,627,775` leafs <b>short tree</b> would contain only 40 subtrees roots and take only 1600 bytes.

Lets denote rank state calculation:

`p` - rank calculation period  
`lbn` - last confirmed block number  
`cbn` - current block number  
`lr` -  length of rank values array  

For rank storing and calculation we have two separate in-memory contexts:

1. Current rank context. It includes last calculated rank state (values and merkle tree) plus
all links and user stakes submitted to the moment of this rank submission.
2. New rank context. It's currently calculating (or already calculated and waiting for submission) rank state. Consists of new calculated rank state (values and merkle tree) plus new incoming links and updated user stakes.

Calculation of new rank state happens once per `p` blocks and going in parallel.
Iteration starts from block number that `≡ 0 (mod p)` and goes till next block number that `≡ 0 (mod p)`.

For block number `cbn ≡ 0 (mod p)` (including block number 1 cause in cosmos blocks starts from 1):

1. Check if rank calculation is finished. If yes then go to (2.) if not - wait till calculation finished
(actually this situation should not happen because it means that rank calculation period is too short).
2. Submit rank, links and user stakes from new rank context to current rank context.
3. Store last calculated rank merkle tree root hash.
4. Start new rank calculation in parallel (on links and stakes from current rank context).

For each block:

1. All links goes to new rank context.
2. New coming CIDs gets rank equals to zero. We could do it by checking last CIDs number and `lr`
(it's obviously equals to number of CIDs that already have rank). Then add CIDs with number `>lr`
to the end of this array with value equal to zero.
3. Update current context merkle tree with CIDs from previous step
4. Store latest merkle tree from current context (lets call it last block merkle tree).
4. Check if new rank calculation finished. If yes go to (4.) if not go to next block.
5. Push calculated rank state to new rank context. Store merkle tree of newly calculated rank.

To sum up. In <b>current rank context</b> we have rank state from last calculated iteration (plus, every block, it updates with new CIDs). And we have links and user stakes that are participating in current rank calculation iteration (whether it finished or not). <b>New rank context</b> contains links and stakes that will go to next rank calculation and newly calculated rank state (if calculation is finished) that waiting for submitting.

If we need to restart node firstly, we need to restore both contexts (current and new).
Load links and user stakes from database using different versions:
1. Links and stakes from last calculated rank version `v = lbn - (lbn mod n)` goes to current rank context.
2. Links and stakes between versions `v` and `lbn` goes to new rank context.

Also to restart node correctly we have to store next entities in database:

1. Last calculated rank hash (merkle tree root)
2. Newly calculated rank short merkle tree
3. Last block short merkle tree

With <b>last calculated rank hash</b> and <b>newly calculated rank merkle tree</b> we could check if rank calculation was finished before node restart. If they are equal then rank wasn't calculated and we should run rank calculation. If not we could skip rank calculation and use <b>newly calculated rank merkle tree</b> to participate in consensus when it comes to block number `cbn ≡ 0 (mod p)` (rank values will not be available until rank calculation happens
in next iteration. Still validator can participate in consensus so nothing bad).

<b>Last block merkle tree</b> necessary to participate in consensus till start of next rank calculation iteration. So, after restart we could end up with two states:
1. Restored current rank context and new rank context without rank values (links, user stakes and merkle tree).
2. Restored current rank context without rank values. Restored new rank context only with links and user stakes.

Node is able to participate in consensus but cannot provide rank values (and merkle proofs)
till two rank calculation iterations finished (current and next).

Search index should be ran in parallel and do not influence work of consensus machine.
Validator should be able to turn off index support. May be even make it a separate daemon.

<b>Base idea.</b> Always submit new links to index and take rank values from current context (insert in sorted array operation). When new rank state is submitted trigger index to update rank values and do sortings (in most cases new arrays will be almost sorted).

Need to solve problem of adjusting arrays capacity (to not copy arrays each time new linked cid added). Possible solution is to adjust capacity with reserve before resorting array.

Therefore for building index we need to find sorting algorithm that will be fast on almost sorted arrays. Also we should implement it for GPU (so it should better be parallelizable:
Mergesort(Timsort), Heapsort, Smoothsort ... ?

Now we have proof of rank of any given content address. While the relevance is still subjective by nature we have a collective proof that something was relevant for some community at some point in time.

```
For any given CID it is possible to prove the relevance
```

Using this type of proof any two [IBC compatible](dura://QmdCeixQUHBjGnKfwbB1dxf4X8xnadL8xWmmEnQah5n7x2.ipfs) consensus computers can proof the relevance to each other, so domain specific relevance machines can flourish. Thanks to inter-blockchain communication protocol you basically can launch you own domain specific search engine either private or public by forking cyberd which is focuses on the public common knowledge. So in our search architecture domain specific relevance machine can learn from common knowledge. We are going to work on IBC during `smith` implementation.

In our relevance for commons `euler` implementation proof of relevance root hash is computed on cuda gpus every round.

## Speed and scalability

We need very fast conformation times in order to feels like usual web app. It is strong architecture requirement that shape an economic topology and scalability of cyber protocol.

Proposed blockchain design is based on Tendermint consensus algorithm with 146 validators and has very fast 1 second finality time. Average confirmation timeframe at half the second with asynchronous interaction make complex blockchain search almost invisible for agents.

Let us say that our node implementation based on cosmos-sdk can process 10k transactions per second. Thus every day at least 8.64 million agents can submit 100 cyberlinks each and impact results simultaneously. That is enough to verify all assumptions in the wild. As blockchain technology evolve we want to check that every hypothesis work before scale it further. Moreover, proposed design needs demand for full bandwith in order the relevance become valuable. That is why we strongly focus on accesable, but provable distribution to millions from inception.

## Implementation in browser

We wanted to imagine how that could work in web3 browser. To our disappointment we [was not able](https://github.com/cybercongress/cyb/blob/master/docs/comparison.md) to find web3 browser that is able to showcase the coolness of proposed approach in action. That is why we decide to develop web3 browser [cyb](https://github.com/cybercongress/cyb/blob/master/docs/cyb.md) that has sample application .cyber for interacting with `cyber://` protocol.

## From Inception to Genesis

It is trivial to develop `euler` like proof-of-concept implementation but it is hard to achieve stable protocol `merkle` a lot of CYB value on which can exist. `euler` is Inception that already happened, `merkle` is Genesis that is far away. That is why we decide to innovate a bit on the going mainnet process. We do not have CYB balances and rank guaranties before `merkle` but we can have exponentially growing semantic core which can be improved based on measurements and observations during development and gradual transfer of value since `euler`. So think that Genesis or `merkle` is very stable and can store semantic core and value, but all releases before can store the whole semantic core and only part of the value you would love to store due to weak security guaranties. The percents of CYB value to be distributed in every testnet:

```toml
euler = 0.7
smith = 3.3
darwin = 8
turing = 15
nash = 21
weiner = 25
merkle = 27
```

In order to secure value of CYB before genesis 100 CBD ERC-20 tokens [are issued](https://etherscan.io/token/0x136c1121f21c29415D8cd71F8Bb140C7fF187033) by cyberFoundation. Based on CBD cyber proto substance snapshot balances are recomputed 7 times according to defined proportions.

``` Need triple check !!!
100 CBD in euler got 70 000 000 000 000 CYB
```

Essentially CBD substance is distributed by cyberFoundation in the following proportion:

- `Proof-of-use : 70%` is allocated to web3 agents according to some probabilistic algorithm. E.g. first `euler` proof-of-use distribution we call Satoshi Lottery is allocated to top 1M key owned Ethereum addresses by [SpringRank](dura://QmNvxWTXQaAqjEouZQXTV4wDB5ryW4PGcaxe2Lukv1BxuM.ipfs).
- `Proof-of-code: 15%` is allocated for direct contribution into the code base. E.g. as assigned by cyberFoundation to cyberCongress contribution including team is 11.2% and the other 3.8% allocated to developers community projects such as Gitcoin community and cyberColony based experimental organization.
- `Proof-of-value: 15%` is allocated for direct contribution of funds. 8% of this value either has been already contributed nor has some reservation for ongoing contributions by close friends and 7% is going to be distributed during eos like auction not defined exactly yet. All contribution will go to Aragon based cyberFoundation and will be managed by CBD token holders.

Details of code and value distribution will be produced by cyberFoundation.

Except 7 proof-of-use lotteries CYB tokens can be created only by validators based on default staking and slashing parameters in accordance with the following approximate percents of inflation per year:

```toml
euler = 200
smith = 134
darwin = 90
turing = 60
nash = 40
weiner = 27
merkle = 18
```

Basic consensus is that newly created CYB tokens are distributed to validators as they do the most essential work to make relevance machine run both in terms of energy consumed for computation and cost for storage capacity. So validators decide where the tokens can flow further.

After Genesis starting inflation rate will become fixed at `4 200 000 CYB` per block.

## Satoshi Lottery

Satoshi Lottery is the inception version of proof-of-use distribution that happens in the tenth birthday of Bitcoin Genesis at 3 Jan 2019. It is highly experimental way of provable distribution.

Basic algorithm is of 5 steps:

```
- Compute SpringRank for Ethereum addresses
- Sort by SpringRank
- Filter top 1M addresses by SpringRank
- Compute CYB balances based on CBD
- Create genesis for cyber protocol
```

Tolik's article goes here

## Possible applications

A lot of cool applications can be build on top of proposed architecture:

_Web3 browsers_. It easy to imagine the emergence of a full-blown blockchain browser. Currently, there are several efforts for developing browsers around blockchains and distributed tech. Among them are Beaker, Mist, Brave and Metamask .. All of them suffer from trying to embed web2 in we3. Our approach is a bit different. We develop truly web3 experience first.

_Programmable semantic cores_. Currently the most popular keywords in gigantic semantic core of Google are keywords of apps such as youtube, facebook, github etc. But developers has very limited possibility to explain Google how to better structure results. Cyber approach bring this power back to developers. On any given user input string in any application relevant answer can be computed either globally, in the context of an app, a user, a geo or in all of them combined.

_Search actions_. Proposed design enable native support for blockchain asset related activity. It is trivial to design an applications which are (1) owned by creators, (2) appear right in search results and (3) allow a transact-able call to actions with (4) provable attribution of a conversion to search query. e-Commerce has never been so easy for everybody.

_Offline search_. IPFS make possible easy retrieval of documents from surroundings without global internet connection. cyberd can itself can be distributed using IPFS. That create a possibility for ubiquitous offline search.

_Command tools_. Command line tools can rely on relevant and structured answers from a search engine. That practically means that the following CLI tool is possible to implement

```
>  cyberd earn using 100 gb hdd

Enjoy the following predictions:
- apt install go-filecoin: 0.001 btc per month
- apt install siad: 0.0001 btc per month per GB
- apt install storjd: 0.00008 btc per month per GB

According to the best prediction I made a decision try `mine go-filecoin`

Git clone ...
Building go-filecoin
Starting go-filecoin
Creating wallet using @xhipster seed
You address is ....
Placing bids ...
Waiting for incoming storage requests ...

```
Search from CLI tools will inevitably create a highly competitive market of a dedicated semantic core for bots.

_Autonomous robots_. Blockchain technology enables creation of devices which are able to earn, store, spend and invest digital assets by themselves.

> If a robot can earn, store, spend and invest she can do everything you can do

What is needed is simple yet powerful state reality tool with ability to find particular things. cyberd offers minimalistic but continuously self-improving datasource that provides necessary tools for programming economically rational robots. According to [top-10000 english words](https://github.com/first20hours/google-10000-english) the most popular word in english is defined article `the` that mean a pointer to a particular thing. That fact can be explained as the following: particular things are the most important for us. So current nature of our current semantic computing is to find very unique things. Hence understanding of unique things become essential for robots too.

_Language convergence_. A programmer should not care about what language do the user use. We don't need to have knowledge of what language user is searching in. Entire UTF-8 spectrum is at work. A semantic core is open so competition for answering can become distributed across different domain-specific areas, including semantic cores of different languages. The unified approach creates an opportunity for cyber•Bahasa. Since the Internet, we observe a process of rapid language convergence. We use more truly global words across the entire planet independently of our nationality, language and race, Name the Internet. The dream of truly global language is hard to deploy because it is hard to agree on what mean what. But we have tools to make that dream come true. It is not hard to predict that the shorter a word the more it's cyber•rank will be. Global publicly available list of symbols, words and phrases sorted by cyber•rank with corresponding links provided by cyberd can be the foundation for the emergence of truly global language everybody can accept. Recent scientific advances in machine translation [GNMT] are breathtaking but meaningless for those who wish to apply them without Google scale trained model. Proposed cyber•rank offers exactly this.

This is sure not the exhaustive list of possible applications but very exciting, though.

## Economic protection is `smith`

About private knowledge on relevance.

## Ability to evolve is `darwin`

About the importance of alternative implementation.

## `turing` is about computing more

Ability to programmatically extend state based on proven knowledge graph is of paramount importance. Thus we consider that WASM programs will be available for execution in cyber consensus computer on top of knowledge graph.

Our approach to economics of consensus computer is that users buys amount of ram, cpu and gpu as they want to execute programs. The following list is simple programs we can envision that can be build on top of simple relevance machine.

_Self prediction_. A consensus computer is able to continuously build a knowledge graph by itself predicting existence of cyberlinks and applying this predictions to a state of itself. Hence a consensus computer can participate in economic consensus of cyber protocol.

_Universal oracle._ A consensus computer is able to store the most relevant data in key value store, where key is cid and value is bytes of actual content. She is doing it by making a decision every round about which value she want to prune and which she want to apply based on utility measure of content addresses in the knowledge graph. In order to compute utility measure validators check availability and size of a content for top ranked content address in the knowledge graph. Then weighted on size of cids and rank of . N is variable defined by consensus. The most  The logic is simple:   The good question os Every available content Emergent key-value store will be available to write for consensus computer only and not agents. but values can be used in programs.

_Proof of location_. It is possible to construct cyberlinks with proof-of-location based on some existing protocol such as [Foam](dura://QmZYKGuLHf2h1mZrhiP2FzYsjj3tWt2LYduMCRbpgi5pKG.ipfs). So location based search also can become provable if web3 agents will mine triangulations and attaching proof of location for every link chain.

_Proof of web3 agent_. Agents are a subset of content addresses with one very important property: consensus computer can prove an existence of private keys for content addresses for the subset of knowledge graph even if those addresses has never transacted in its own chain. Hence it is possible to compute a lot of provable stuff on top of that knowledge. Eg. some inflation can be distributed to addresses that have never transacted in the cyber network but have provable link.

_Motivation for read requests_. It would be great to create cybernomics not only for write requests to consensus computer, but from read requests also. So read requests can be two order of magnitude cheaper. Read requests to a search engine can be provided by the second tier of nodes which earn CYB tokens in state channels. We consider to implement state channels based on HTLC and proof verification which unlocks amount earned for already served request (new signatures post via requester/user to cybernode via whisper/ipfs-pub-sub)

_Prediction markets on link relevance_. We can move idea further by ranking of knowledge graph based on prediction market on links relevance. App that allow to bet on link relevance can become unique source of truth for direction of terms as well motivate to submit more links.

_Private cyberlinks_. Privacy is foundational. While we are committed to privacy achieving implementation of private is unfeasible for our team hence we leave it on after genesis times. zero knowledge proofs. . The problem is to compute rank based on link of an agent based on its ranking without revealing identity. Zero knowledge proofs in general are very expensive. Privacy of search by design or as an option? Interactive proofs.

## In search for equilibria is `nash`

On scalability trilema ...

Decentralization comes with costs and slowness. We want to find a good balance between speed, relience and ability to scale, as we believe all three are sensitive for widespread web3 adoption.

That is the area of research for us now.

Let us say that our node implementation based on cosmos-sdk can process 10k transactions per second. Thus every day at least 8.64 million agents can submit 100 cyberlinks each and impact results simultaneously. That is enough to verify all assumptions in the wild. As blockchain technology evolve we want to check that every hypothesis work before scale it further.

## On faster evolution at `weiner`

Basic purpose of wiener stage is to be able to update the consensus of a network from a consensus computer state using some onchain upgrade mechanism.

Evolvability and governance are connected tightly.

Ability to reflect input from the world and output changes of itself is essential evolutionary feature. Hence, thanks to cosmos-sdk `euler` implementation has basic but very power features such as onchain voting with vetos and abstains that drastically simplified open discussions for change. So we are going to use this feature from the inception of the network.

But we can go to different direction than cosmos-sdk offers. Following ideas from Tezos in Weiner we can define the current state of a protocol as immutable content address that is included in round merkle root.

Also instead of formal governance procedure we would love to check the hypothesis that changing state of a protocol is possible indeed using relevance machine itself.

Starting protocol can be as simple as follows:

> The more closer some content address to literally `cyber-protocol-current` content address the more probability than it will become winning. The most close protocol `cyber-protocol-current` is the protocol which is the most relevant to users.

Hence it is up to nodes to signal `cyber-protocol-current` by sending cyberlinks with semantics like `<cyber-protocol-current> <cid-of-protocol>`.

## Genesis is secure as `merkle`

Before unleashing our creature we need to have strong assurance that implementations are secure. Merkle is our final genesis release after security audits and more formalism.

After this release the network of relevance machines become fully functional and evolvable.

## Conclusion

We describe a motivated blockchain based search engine for web3. A search engine is based on the content-addressable peer-to-peer paradigm and uses IPFS as a foundation. IPFS provide significant benefits in terms of resources consumption. IPFS addresses as a primary objects are robust in its simplicity. For every IPFS hash cyber•rank is computed by a consensus computer with no single point of failure. Cyber•rank is a spring rank with economic protection from selfish predictions. Sybil resistance is also implemented on two levels: during id generation and during bandwidth limiting. Embedded smart contracts offer fair compensations for those who is able to predict relevance of content addresses. The primary goal is indexing of peer-to-peer systems with self-authenticated data either stateless, such as IPFS, Swarm, DAT, Git, BitTorent, or stateful such as Bitcoin, Ethereum and other blockchains and tangles. Proposed semantics of  linking offers robust mechanism for predicting meaningful relations between objects. A source code of a relevance machine is open source. Every bit of data accumulated by a consensus computer is available for everybody if the one has resources to process it. The performance of proposed software implementation is sufficient for seamless user interactions. Scalability of proposed implementation is enough to index all self-authenticated data that exist today. The blockchain is managed by a decentralized autonomous organization which functions under Tendermint consensus algorithm. Thought a system provide necessary utility to offer an alternative for conventional search engines it is not limited to this use case either. The system is extendable for numerous applications and e.g. makes possible to design economically rational self-owned robots that are able to autonomously understand objects around them.

## References

...

## Notes

Abstract and Conclusion must be rewritten

Other parts require proof read, error correction and thorough discussion.

Two fundamental problems associated has not been discussed: vote buying and selfish predictions.
