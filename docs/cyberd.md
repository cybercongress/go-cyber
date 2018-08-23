# cyberd: A search consensus computer

@xhipster, @litvintech
v 0.2
August 2018, Iceland and Tolyatti

## Abstract

An incentivized consensus computer would allow to compute provably relevant answers without opinionated blackbox intermediaries such as Google. Stateless content-addressable peer-to-peer communication networks such as IPFS and stateful consensus computers such as Ethereum provide part of the solution but there are at least three problems associated with implementation. Of course, the first problem is subjective nature of relevance. The second problem is that it is hard to scale consensus computer of knowledge graph due to non linear nature of provably working solutions for web search such as PageRank and more than exponentially growing size of a knowledge graph including its history. The third problem is that the quality of such knowledge graph will suffer from different attack surfaces such as sybil and selfish behavior of interacting agents. In this paper we (1) define a protocol for provable consensus computing of relevance between IPFS objects based on some offline observation and theory behind prediction markets, (2) solve a problem of implementation inside consensus computer based on linear SpringRank and propose workaround for pruning historical state and (3) design distribution and incentive scheme based on known attacks. Also we discuss some considerations on minimalistic architecture of the protocol as we believe that is critical to formation of a network of domain specific search consensus computers. As result of our work some exciting applications emerge.

## Introduction

Existing general purpose search engines are restrictive centralized databases everybody forced to trust. These search engines were designed primarily for client-server architecture based on DNS, HTTP, and IP protocols. The emergence of a distributed protocol stack creates an opportunity for a new kind of Internet. We call it web3. Protocols such as Ethereum create a challenge and opportunity for a search engine based on developing technologies and specifically designed for them. Surprisingly the permission-less blockchain architecture itself allows organizing general purpose search engine in a way inaccessible for previous architectures.

## The protocol

- def state transition
- take txs
- validate signatures and timestamps
- emit predictions of relevant objects for every valid tx
- every N block calculate spring rank
- distribute inflation via objects's ranks score
- each object - Q(inner, question) and A(outer) and U(user) takes tokens via their ranks
- Q/A - number of objects inner and outer edges define proportions of emission to token
- U - vests reputations to their balances
- every N blocks nodes reach consensus around prunned state history via ipfs hash of state blob

## Search as prediction market on the relevance of links

(1) maximization of knowledge graph based , aligned with computational, storage and broadband bound and (2) reduction for attack surfaces such as sybil and selfish behavior

## Inductive reasoning

A useful property of a computer is that it must know nothing about objects except when, who and where some prediction was asked. If we assume that a consensus computer must have some information about linked objects the complexity of such model growth unpredictably, hence a requirements for a computer for memory and computations. That is, deduction of a meaning inside consensus computer is expensive thus our design hardly depend on the blind assumption. Instead of we design incentives around meaning extractions

## About when

Proof-of-history + Tendermint

## About who

Digital signatures and zero knowledge proofs

## About where

Mining triangulations and attaching proof of location for every search query can significantly

## Link chains

Explain the concept of link chains and demonstrate a case of semantic linking based on interpretation of an ERC-20 transfer transaction or something similar

## Bad behavior

We design a system under assumption that in terms of search such a thing as bad behavior does not exist as nothing bad can be in the intention of finding answers. Ranks is computed on the only fact that something has been searched, thus linked and as result affected predictive model. Good analogy is observing in quantum mechanics. So no negative voting is implemented. Doing this we remove subjectivity out of the protocol and can define one possible method for proof of relevance. Also this approach significantly reduce attack surface.

## Problem of ranking inside consensus computers

Consensus computers bring serious resource bounds

It's possible to compute a ranks for the whole MerkleDAG. But there are two problems with it:

- An amount of links in MerkleDAG grows `O(n^2)`. That is not either conventional web pages with 20-100 links per page. For 1 Mb file can be thousands of links. Once the network starts to take off, complexity inevitably increases.
- Even if we address some algorithm to extract relevances from these links we should address even more algorithms to extract a meaning.

What we need is to find a way to incentivize extraction from this data fog a meaning that is _relevant to users queries_.

We can represent our data structure as directed acyclic graph where vertices are indexed documents and edges are directed links between them.

![](https://docs.google.com/drawings/d/1--Uj85OiU-uwj0gxUFWZDbjPuby3IEMBBVFSmHTNkDc/pub?w=785&h=436)

Our model is recursive (check SpringRank) and requires the enormous amount of calculations which are limited within blockchain design. Model recalculation does not happen on a periodic basis rather it continuous. We consider introducing consensus variable, in addition to a block size, in order to target processing capacity of the network. Let's call it a _computing target of documents per block_ or CTD. Any witness will be able to set a number of documents the network should recompute every block. The blockchain takes as input computing target of legitimate witnesses and computes CTD as daily moving average. Based on CTD blockchain can schedule the range of IPFS hashes that should be recomputed by every witness per round.

## cyber•Rank
Nebulas fail.
No rank computed inside consensus computing => no possibility to incentivize network participants to form predictions on relevance. A problem here is that computational complexity of conventional ranks grow sublineary with the growth of the network. So we need to find deterministic algorithm that allow to compute a rank for continuously appended network to scale the consensus computer to orders of magnitude that of Google. Also an algorithm must have good prediction capability for existence of relevant to an object links.

Spring ranks cons: linear & better ranking (original work, case proof with Steem)

## State grow history problem

Every N blocks cybernodes prune blockchain/history and calculate IPFS hash for them/publish to IPFS, add hash to block and validate them with consensus. This state/blob economicaly finalized and new node start from them. Cybernodes motivated to store/provide this blob cause this cause network grow. Dynamically recalculate N with economy, newtork size, rank score...
~ vs direct IPFS DAG

## Motivation on request's processing problem

Explain an economic difference and censorship impact between read search queries and write search queries

Idea:
All nodes run payment channels to serve request for their users and take tokens for request processing. Because this is heavy computation (only cybernodes will serve this) nodes will serve this only with payments for them. Cybernodes run protocol and earn tokens for read requests.

Solution is payment channels based on HLTC and proof verification which unlocks amount earned for already served request (new signatures post via requester/user to cybernode wia whisper/ipfs-pub-sub)

## Prediction by a consensus computer

## Smart contracts for predictions

## Selfish linking

Manipulating the rank: Sybil Attack: Quadratic voting. Rank squared. Decay from linear in the beginning.

## Spam protection

In the center of spam protection system is an assumption that write operations can be executed only by those who have vested interest in the success of the search engine. Every 1% of stake in consensus computer gives the ability to use 1% of possible network broadband and computing capabilities. As nobody uses all possessed broadband we use fractional reserves while limiting broadband like ISPs do. Details of an approach can be found in a Steem white paper.

Auditing and curation are based on Steem reward mechanism. It is Sybil-resistant approach as votes are quadratic based on principle 1 token in system = 1 vote. In order to vote one should vest in shares for at least for 20 weeks. That solve a problem entirely because those who have a right to vote are strongly incentivized in a growth of his wealth. In order to prevent abuse of auditing and curation voting power decay implemented exactly as in Steem.

## Distribution Mechanism

Describe drops and other built in incentives

## Incentive Structure

To make cyber•rank economically resistant to Sybil attack and to incentivize all participant for rational behavior a system uses CYBER token.

Since inception, a network prints 42 CYBER every block.

Reward pool is defined as 100% of emission and split among the contracts according to (?):

- Validators
- Linkers

## Applications

_Web3 browsers_. It easy to imagine the emergence of a full-blown blockchain browser. Currently, there are several efforts for developing browsers around blockchains and distributed tech. Among them are Beaker, Mist and Brave .. All of them suffer from very limited functionality. Our developments can be useful for teams who are developing such tools.

_Programmable semantic cores_. Relevance everywhere means that on any given user input string in any application relevant answer can be computed either globally, in the context of an app or in the context of a user.

_Actions in search_. Proposed design enable native support for blockchain asset related activity. It is trivial to design an applications which are (1) owned by creators, (2) appear right in search results and (3) allow a transact-able call to actions with (4) provable attribution of a conversion to search query. e-Commerce has never been so easy for everybody.

_Shared mempools_.

_Offline search_. IPFS make possible easy retrieval of documents from surroundings without the internet connection. cyberd itself can be distributed using IPFS. That create a possibility for ubiquitous offline search.

_Smart Command Tools_. Command line tools can rely on relevant and structured answers from a search engine. That practically means that the following CLI tool is possible to implement

```
>  cyberd earn using 100 gb hdd

Enjoy the following predictions:
- apt install siad /// 0.0001 btc per month per GB
- apt install storjd /// 0.00008 btc per month per GB
- apt install filecoind /// 0.00006 btc per month

According to the best prediction I made a decision try `apt install siad`

Git clone ...
Building siad
Starting siad
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

## Extensibility, Upgradability and Governance

That is a big question!

## Performance

Proposed blockchain design is based on tendermint consensus algorithm and has fast and predictable near seconds block confirmation time and very fast finality time. Average confirmation timeframe is less than second thus conformations can be asynchronous and nearly invisible for users. A good thing is that users don't need confirmations at all before getting search response as there is no risk associated with that.

## Scalability

Our node implementation theoretically can process about ? predictions per second. This theoretical bound is primarily limited with the possibility to replay a blockchain [https://steemit.com/blockchain/@dantheman/how-to-process-100m-transfers-second-on-a-single-blockchain]. As of now, all blockchains are about 1B immutable documents which size is about 200 GB with average tx 200 kb. We need to store all hashes which are on average 64 bytes long. We estimated that storing in the index all blockchain documents as IPFS hashes and votes are roughly the same as storing all raw blockchain data. Linking 1B documents create significant overhead as blockchain index size can be up to 100 times more. Given this, we can assume that indexing all existing blockchains require about 4TB of SSD space. This is affordable for commodity hardware with 10x scaling capability without a necessity for sharding across several machines. We assume this is enough scalability margin for proof-of-concept.

Initial indexing of 1B documents and 100B links will require a continuous load of the network at the upper bound of its capacity in the first year of its existence. If we assume that network will be able to process 10k transactions per second with 2MB block size we will be able to index all blockchains in 4 months. Further operations will require significantly less capacity as currently, not more than 1000 transactions per second happen among all blockchains.

Based on the proposed search appliance we estimate that participants will require investing around $1M for dedicated hardware (21 witnesses) and the same amount for backup nodes. Thus overall costs of hardware network infrastructure can be around $2M. after full deployment.

Worth to note that the network doesn't require ultimate configuration at the start and is able to optimize initial investments by the costs of time to index all blockchains. Thus costs at launch can be around $200k. Given that mining industry has been rapidly developed last years this can not be a showstopper for a project. We expect huge interest from miners as slots are limited with 21 fully paid nodes and ~20 of partially paid nodes (depend on the market).

Possible scalability improvements include:
- Hardware. This year Intel Optane [http://www.intel.com/content/www/us/en/architecture-and-technology/intel-optane-technology.html] That creates an opportunity to converge RAM and SSD. Our design has 3-year hardware margin for moving to cheaper and more dense next generation memory.
- Software. The future of consensus computer optimization is in parallel processing. We are going to seriously invest in research of this field. Solving this issue will enable the network to scale nearly infinitely.
- Caching/pin   ing of more scored objects in network (provides faster and in-place gateway to ipfs)

## Conclusion

We describe and implement a motivated blockchain based search engine for the permanent web. A search engine is based on the content-addressable peer-to-peer paradigm and uses IPFS as a foundation. IPFS provide significant benefits in terms of resources consumption. IPFS addresses as a primary objects are robust in its simplicity. For every IPFS hash cyber•rank is computed by a consensus computer with no single point of failure. Cyber•rank is a spring rank enchased with economic based sybil protection. Embedded smart contracts offer fair compensations for those who is able to predict popularity of hashes. The primary goal is indexing of peer-to-peer systems with self-authenticated data either stateless, such as IFSS, DAT, GIT, BitTorent, or stateful such as Bitcoin, Ethereum and other blockchains and tangles. Proposed market of linking offers necessary incentives for outsourcing computing part responsible for finding meaningful relations between objects. A source code of a search engine is open source. Every bit of data accumulated by a consensus computer is available for everybody for free. The performance of proposed software implementation is sufficient for seamless user interactions. Scalability of proposed implementation is enough to index all self-authenticated data that exist today. The blockchain is managed by a decentralized autonomous organization which functions under Tendermint consensus algorithm. Thought a system provide necessary utility to offer an alternative for conventional search engines it is not limited to this use case either. The system is extendable for numerous applications and e.g. makes possible to design economically rational self-owned robots that are able to autonomously understand objects around them.

## References
