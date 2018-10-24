# cyberd: A search consensus computer for web3

@xhipster, @litvintech

v 0.3. Research notes

Kenig and Minsk

```
cyb:
- nick. a friendly software robot who helps you explore universes

cyber:
- noun. a superintelligent network computer for answers
- verb. to do something intelligent, to be very smart

CYB:
- ticker. transferable token expressing will to become smarter

CYBER:
- ticker. non-transferable token expressing intelligence
```

## Content
<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->
<!-- code_chunk_output -->

* [cyberd: A search consensus computer for web3](#cyberd-a-search-consensus-computer-for-web3)
	* [Content](#content)
	* [Abstract](#abstract)
	* [Introduction to web3](#introduction-to-web3)
	* [The protocol](#the-protocol)
	* [Knowledge graph](#knowledge-graph)
	* [Agents of knowledge](#agents-of-knowledge)
	* [Link chains](#link-chains)
	* [Notion of consensus computer](#notion-of-consensus-computer)
	* [Relevance machine](#relevance-machine)
	* [cyber•Rank](#cyberrank)
	* [Proof of relevance](#proof-of-relevance)
	* [State grow history problem](#state-grow-history-problem)
	* [Motivation for read requests](#motivation-for-read-requests)
	* [Self prediction](#self-prediction)
	* [Universal oracle](#universal-oracle)
	* [Smart contracts](#smart-contracts)
	* [Selfish predictions](#selfish-predictions)
	* [Spam protection](#spam-protection)
	* [Distribution](#distribution)
	* [Incentive structure](#incentive-structure)
	* [Applications](#applications)
	* [Evolvability](#evolvability)
	* [Decentralization](#decentralization)
	* [Performance](#performance)
	* [Scalability](#scalability)
	* [Conclusion](#conclusion)
	* [References](#references)

<!-- /code_chunk_output -->

## Abstract

An incentivized consensus computer would allow to compute provably relevant answers without opinionated blackbox intermediaries such as Google. Stateless content-addressable peer-to-peer communication networks such as IPFS and stateful consensus computers such as Ethereum provide part of the solution but there are at least three problems associated with implementation. Of course, the first problem is subjective nature of relevance. The second problem is that it is hard to scale consensus computer of huge knowledge graph. The third problem is that the quality of such knowledge graph will suffer from different attack surfaces such as sybil and selfish behavior of interacting agents. In this paper we (1) define a protocol for provable consensus computing of relevance between IPFS objects based on some offline observation and theory behind prediction markets, (2) solve a problem of implementation inside consensus computer based on SpringRank and propose workaround for pruning historical state and (3) design distribution and incentive scheme based on our experience and known attacks. Also we discuss some considerations on minimalistic architecture of the protocol as we believe that is critical for formation of a network of domain specific search consensus computers. As result of our work some applications never existed before emerge.

## Introduction to web3

Original protocols of the Internet such as TCP/IP, DNS, URL and HTTPS brought a web into the point there it is now. Along with all benefits they has created they brought more problem into the table. Globality being a key property of the the web since inception is under real threat. Speed of connections degrade with network grow and from ubiquitous government interventions into privacy and security of web users. One property, not obvious in the beginning, become really important with everyday usage of the Internet: its ability to exchange permanent hyperlinks thus they would not break after time have pass. Reliance on one at a time internet service provider architecture allow governments censor packets is the last straw in conventional web stack for every engineer who is concerned about the future of our children. Other properties while being not so critical are very desirable: offline and real-time. Average internet user being offline must have ability to work with the state it has and after acquiring connection being able to sync with global state and continue verify state's validity in realtime while having connection. Now this properties offered on app level while such properties must be integrated into lower level protocols.

The emergence of a distributed protocol stack [W3S] creates an opportunity for a new kind of Internet. We call it web3. It has a promise to remove problems of conventional protocol stack and add to the web better speed and more accessible connection. But as usually in a story with a new stack, new problems emerge. One of such problem is general purpose search. Existing general purpose search engines are restrictive centralized databases everybody forced to trust. These search engines were designed primarily for client-server architecture based on TCP/IP, DNS, URL and HTPPS protocols. Web3 create a challenge and opportunity for a search engine based on developing technologies and specifically designed for them. Surprisingly the permission-less blockchain architecture itself allows organizing general purpose search engine in a way inaccessible for previous architectures.

## The protocol

- def knowledge graph state
- take link chains
- check that linkchain signatures are valid
- check that resource consumption of a signer is not exceed 24 moving average
- emit prediction of links for every valid link chain
- every block calculate cyber•rank deltas for the knowledge graph
- every block distribute 42 CYB based on links CYBERs
- for links with proven keys distribute payouts based on CYBER
- for links without proven keys distribute payouts according to CYBER weight of incoming links with proven keys
- every block apply predictions of consensus computer according to prediction bound
- every block write data to key/value store according to storage bound based on size and cyber•rank
- every epoch nodes reach consensus around pruned state history via ipfs hash of state blob

## Knowledge graph

We represent a knowledge graph as weighted graph of directed links between content addresses.

> Illustration


## Agents of knowledge

Our knowledge graph include everything it can including agents themselves. This subset of content addresses with proven knowledge of private keys and cuckoo cycle pow ids is being used by incentive scheme of a consensus algorithm.

## Link chains

A concept of linkchain is a convention around simple semantics of communication format with a consensus computer.

- a tx format is `<peer id>` `<up to 7 ipfs hashes of links>` `<signature>`
- `<signature>` must be valid from `<peer id>`

Explain the concept of link chains and demonstrate a case of semantic linking based on interpretation of an ERC-20 transfer transaction or something similar

> Illustration

Amount of content addresses ipfs-addresses: up to 7
link-chain-proof: ...
cyber-protocol-current: ...

## Notion of consensus computer

Consensus computer is an abstract machine that has capacity in terms of fundamental computing resources such as memory, processing and bandwidth.

///
Ideal consensus computer is a computer in which sum of agents contribution as actual . Its like a perforamnce indicator resources simple part of memory, processing and bandwidth of all computers.
///

(1) maximization of knowledge graph, aligned with computational, storage and broadband bound and (2) reduction for attack surfaces such as sybil and selfish behavior

> Illustration

## Relevance machine

Ranking of knowledge graph based on prediction market on links relevance.

A useful property of a computer is that it must have inductive reasoning property. She must be able to interfere predictions without any knowledge about objects except when, who and where some prediction was asked. If we assume that a consensus computer must have some information about linked objects the complexity of such model growth unpredictably, hence a requirements for a computer for memory and computations. That is, deduction of a meaning inside consensus computer is expensive thus our design hardly depend on the blindness assumption. Instead we design incentives around meaning extractions

Proof of who. Digital signatures and zero knowledge proofs. Privacy is foundational. The problem is to compute rank based on link of an agent based on its ranking without revealing identity. Zero knowledge proofs in general are very expensive. Privacy of search by design or as an option?

Proof of when. Proof-of-history + Tendermint

Proof of where. Mining triangulations and attaching proof of location for every link chain

## cyber•Rank

As input for every edge value get signer account's:

    sums of `<CYBER>`

    plus

    (sum of square root from `<CYB>`) squared

Ranking using consensus computer is hard because consensus computers bring serious resource bounds. e.g. Nebulas fail to deliver something useful onchain.

> No rank computed inside consensus computer => no possibility to incentivize network participants to form predictions on relevance.

A problem here is that computational complexity of conventional ranks grow sublineary with the growth of the network. So we need to find (1) deterministic algorithm that allow to compute a rank for continuously appended network to scale the consensus computer to orders of magnitude that of Google. Perfect algorithm (2) must have linear memory and computation complexity. The most importantly without having (3) good prediction capabilities for existence of relevant links it proves to be useless. One of recent algorithms: SpringRank.

Original idea came from physics. Links represented as system of springs with some energy.

> Illustration

H(s) = 1/2 ...

## Proof of relevance

Payouts ...

Based on linkchains its possible to know which content addresses proved that private keys are exist, and which do not.

We design a system under assumption that in terms of search such thing as bad behavior does not exist as nothing bad can be in the intention of finding answers.

> Ranks is computed on the only fact that something has been searched, thus linked and as result affected predictive model.

Good analogy is observing in quantum mechanics. So no negative voting is implemented. Doing this we remove subjectivity out of the protocol and can define one possible method for proof of relevance. Also this approach significantly reduce attack surface. Implication of this assumption is that we must bind resource supply of relevance machine with demand of queries.

## State grow history problem

Every N blocks cybernodes prune blockchain/history and calculate IPFS hash for them/publish to IPFS, add hash to block and validate them with consensus. This state/blob economicaly finalized and new node start from them. Cybernodes motivated to store/provide this blob cause this cause network grow. Dynamically recalculate N with economy, network size, rank score...

## Motivation for read requests

Explain an economic difference and censorship impact between read search queries and write search queries

Idea:
All nodes run payment channels to serve request for their users and take tokens for request processing. Because this is heavy computation (only cybernodes will serve this) nodes will serve this only with payments for them. Cybernodes run protocol and earn tokens for read requests.

Solution is payment channels based on HLTC and proof verification which unlocks amount earned for already served request (new signatures post via requester/user to cybernode via whisper/ipfs-pub-sub)

## Self prediction

A consensus computer is able to continuously build a knowledge graph by itself predicting existence of links and applying this predictions to a state of itself.

Idea: Everything that has been earned by a consensus computer can form a validators budget.

Forgetting links: Prune min possible rank / 2

## Universal oracle

A consensus computer is able to store the most relevant data in key value store. She is doing it by making a decision every block about what record he want to prune and what he want to apply. This key-value store can be ...

## Smart contracts

Ability to programmatically extend state based on proven knowledge graph is of paramount importance. Thus we consider that WASM programs will be available for execution on top of knowledge graph.

## Selfish predictions

Protection from selfish predictions: from linear in the beginning to sum of square roots being mature.

## Spam protection

In the center of spam protection system is an assumption that write operations can be executed only by those who have vested interest in the evolutionary success of a consensus computer. Every 1% of stake in consensus computer gives the ability to use 1% of possible network broadband and computing capabilities. As nobody uses all possessed broadband we can use fractional reserves while limiting broadband like ISPs do, but we cant because it degrades value of rank. So we just compute 24 hours moving avarage and

In order to automate governance process we will not add a feature of staking and unstaking, but instead implement s-curve type of automatic staking time depending on which all bandwith and computational limits will be accounted

## Distribution

No ERC-20.
Reasons: expensive and non deterministic.

- Compute SpringRank for Ethereum (Link is tx, weight is amount)
- Create initial genesis for our network

Every snapshot amount to be distributed is announced for the next snapshot. Every new PoC do distribution based on snapshot of previous chain merged with new Ethereum snapshot. For each new PoC amount of distribution which goes to Ethereum snapshot decreases.

## Incentive structure

To make cyber•rank economically resistant to Sybil attack and to incentivize all participant for rational behavior a system uses CYBER token.

Since inception, a network prints 42 CYB every block.

Reward pool is defined as 100% of emission and split among the contracts according to (?):

- Validators
- Linkers

## Applications

_Web3 browsers_. It easy to imagine the emergence of a full-blown blockchain browser. Currently, there are several efforts for developing browsers around blockchains and distributed tech. Among them are Beaker, Mist and Brave .. All of them suffer from very limited functionality. Our developments can be useful for teams who are developing such tools.

_Programmable semantic cores_. Relevance everywhere means that on any given user input string in any application relevant answer can be computed either globally, in the context of an app or in the context of a user.

_Actions in search_. Proposed design enable native support for blockchain asset related activity. It is trivial to design an applications which are (1) owned by creators, (2) appear right in search results and (3) allow a transact-able call to actions with (4) provable attribution of a conversion to search query. e-Commerce has never been so easy for everybody.

_Offline search_. IPFS make possible easy retrieval of documents from surroundings without global internet connection. cyberd itself can be distributed using IPFS. That create a possibility for ubiquitous offline search.

_Command tools_. Command line tools can rely on relevant and structured answers from a search engine. That practically means that the following CLI tool is possible to implement

```
>  cyberd earn using 100 gb hdd

Enjoy the following predictions:
- apt install go-filecoin /// 0.001 btc per month
- apt install siad /// 0.0001 btc per month per GB
- apt install storjd /// 0.00008 btc per month per GB

According to the best prediction I made a decision try `go get go-filecoin`

Git clone ...
Building go-filecoin
Starting go-filecoin
Creating wallet using your standard seed
You address is ....
Placing bids ...
Waiting for incoming storage requests ...

```
Search from CLI tools will inevitably create a highly competitive market of a dedicated semantic core for bots.

_Autonomous robots_. Blockchain technology enables creation of devices which are able to earn, store, spend and invest digital assets by themselves.

> If a robot can earn, store, spend and invest she can do everything you can do

What is needed is simple yet powerful API about the state of reality evaluated in transact-able assets. Our solution offers minimalistic but continuously self-improving API that provides necessary tools for programming economically rational robots.

_Language convergence_. A programmer should not care about what language do the user use. We don't need to have knowledge of what language user is searching in. Entire UTF-8 spectrum is at work. A semantic core is open so competition for answering can become distributed across different domain-specific areas, including semantic cores of different languages. The unified approach creates an opportunity for cyber•Bahasa. Since the Internet, we observe a process of rapid language convergence. We use more truly global words across the entire planet independently of our nationality, language and race, Name the Internet. The dream of truly global language is hard to deploy because it is hard to agree on what mean what. But we have tools to make that dream come true. It is not hard to predict that the shorter a word the more it's cyber•rank will be. Global publicly available list of symbols, words and phrases sorted by cyber•rank with corresponding links provided by cyberd can be the foundation for the emergence of truly global language everybody can accept. Recent scientific advances in machine translation [GNMT] are breathtaking but meaningless for those who wish to apply them without Google scale trained model. Proposed cyber•rank offers exactly this.

This is sure not the exhaustive list of possible applications but very exciting, though.

## Evolvability

Following ideas from Tezos we can define the current state of a protocol as immutable content address. Thus she can adopt to a new environment changing content address of current protocol following the rules hidden behind previous protocol. We would love to check the hypothesis that it is possible to have a protocol which follows simple rule:

> The more closer some content address to literally `cyber-protocol-current` ipfs address the more probability than it will become winning. The most close protocol `cyber-protocol-current` is the protocol which is the most relevant to users.

Thus nodes must always signal connection with `cyber-protocol-current` by sending linkchains with semantics like `<{cyber-protocol-current}> <cid>`.

## Decentralization

Decentralization comes with costs and slowness. We want to find a good balance between those as we believe both are sensitive for widespread web3 adoption. That is the area of research for us now.

## Performance

We need very fast conformation times.

Proposed blockchain design is based on Tendermint consensus algorithm and has fast and predictable before seconds block confirmation time and very fast finality time. Average confirmation timeframe is less than second thus conformations can be asynchronous and nearly invisible for users. A good thing is that users don't need confirmations at all before getting search response as there is no risk associated with that.

## Scalability

Let us say that our node implementation based on Cosmos-SDK can process 10k transactions per second. Thus every day at least 8.64 million nodes can submit 100 predictions and get back search results simultaneously. That is enough to verify all assumptions in the wild. As blockchain technology evolve we want to check that every hypothesis work before scale it further.

## Conclusion

We describe a motivated blockchain based search engine for web3. A search engine is based on the content-addressable peer-to-peer paradigm and uses IPFS as a foundation. IPFS provide significant benefits in terms of resources consumption. IPFS addresses as a primary objects are robust in its simplicity. For every IPFS hash cyber•rank is computed by a consensus computer with no single point of failure. Cyber•rank is a spring rank with economic protection from selfish predictions. Sybil resistance is also implemented on two levels: during id generation and during bandwidth limiting. Embedded smart contracts offer fair compensations for those who is able to predict relevance of content addresses. The primary goal is indexing of peer-to-peer systems with self-authenticated data either stateless, such as IPFS, Swarm, DAT, Git, BitTorent, or stateful such as Bitcoin, Ethereum and other blockchains and tangles. Proposed semantics of  linking offers robust mechanism for predicting meaningful relations between objects. A source code of a relevance machine is open source. Every bit of data accumulated by a consensus computer is available for everybody if the one has resources to process it. The performance of proposed software implementation is sufficient for seamless user interactions. Scalability of proposed implementation is enough to index all self-authenticated data that exist today. The blockchain is managed by a decentralized autonomous organization which functions under Tendermint consensus algorithm. Thought a system provide necessary utility to offer an alternative for conventional search engines it is not limited to this use case either. The system is extendable for numerous applications and e.g. makes possible to design economically rational self-owned robots that are able to autonomously understand objects around them.

## References

W3S: [Web3 stack](https://github.com/w3f/Web3-wiki/wiki)
