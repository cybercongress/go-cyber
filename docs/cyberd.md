# cyberd: Computing the knowledge from web3
Notes on [`cyber`](https://github.com/cybercongress/cyberd/releases/tag/v0.1.0) release of `cyber://` protocol [reference implementation](https://github.com/cybercongress/cyberd) using Go.

[cyber•Congress](https://cybercongress.ai/): @xhipster, @litvintech, @hleb-albau, @arturalbov, @npopeka

## Abstract
A consensus computer that allows for computing of provably relevant answers. Without the opinionated blackbox intermediaries, such as Google, Youtube, Amazon or Facebook. Stateless, content-addressable, peer-to-peer communication networks, such as IPFS, and stateful consensus computers, such as Ethereum, provide part of the solution. But, there are at least 3 problems associated with that implementation. Of course, the first problem is the subjective nature of relevance. The second problem is that it is hard to scale consensus computers for an over-sized knowledge graph. The third problem is that the quality of such knowledge graph will suffer from different surface attacks such as: sybil attacks and selfish behaviour of the interacting agents. In this paper, we (1) define a protocol for provable consensus computing of relevance between IPFS objects based on Tendermint consensus of cyber•rank, which is computed on GPU, (2) discuss its implementation details, and (3) design a distribution and an incentive scheme, based on our experience. We believe that the minimalistic architecture of the protocol is critical for the formation of a network of domain-specific knowledge in consensus computers. As a result of our work, some applications that never existed before will emerge. We expand the work with our "after Genesis vision" on features and apps.

## Introduction to web3
Original protocols of the Internet such as - TCP/IP, DNS, URL, and HTTPS brought the web into the point where it is located as of now. Along with all the benefits those protocols have created, along with them, they brought problems to the table. Globality - being a vital property of the web since its inception, is under a real threat. The speed of the connection keeps on degrading, while the network grows. Ubiquitous government interventions into our privacy and security of web users, are also an issue. One property, not evident in the beginning, becomes important with everyday usage of the Internet: it's ability to exchange permanent hyperlinks - thus they [would not break after some time had passed](https://ipfs.io/ipfs/QmNhaUrhM7KcWzFYdBeyskoNyihrpHvUEBQnaddwPZigcN). Reliance on "one at a time ISP" architecture, allows governments to effectively censor packets. It is the last straw in a conventional web stack for every engineer who is concerned about the future of our children.

Other properties, while being might not be so critical, are very desirable: offline and real-time connections. The average internet user that is in an offline mode, must have the ability to work with the state he has. After acquiring connection, he must be able to sync with the global state and continue to verify the validity of the state in realtime, while establishing a connection. Currently, these properties are only offered on the application level, while such properties must be integrated into lower level protocols.

The emergence of [web3 stack](https://ipfs.io/ipfs/Qmf3eHU9idMUZgx6MKhCsFPWL24X9pDUi2ECqyH8UtBAMQ) creates an opportunity for a new kind of Internet. The community calls it web3. We call it - "The Great Web", as it is expected that some low-level conventions will become immutable and havent been changed for decades, i.e. immutable content links. They seem promising at removing problems of the conventional protocol stack, and they add to the web - better speed and a more accessible connection. However, as usual, with any story that offers a new stack, new problems begin to emerge. One of such issues, is general-purpose search. The existing general-purpose search engines are restrictive centralized databases that everybody is forced to trust. These search engines were designed primarily for client-server architecture based on TCP/IP, DNS, URL, and HTTPS protocols. Web3 creates a challenge and an opportunity for a search engine based on developing technologies and is specifically designed for them. Surprisingly, the permissionless blockchain architecture itself allows organizing general purpose search engine in a way inaccessible to previous architectures.

## On adversarial examples problem
[The conventional architecture of search engines](https://ipfs.io/ipfs/QmeS4LjoL1iMNRGuyYSx78RAtubTT2bioSGnsvoaupcHR6) is a one entity process. It ranks all the shit that suffers from one difficult, but a particular problem, that is still has not been solved, even by the brilliant scientists from Google: [examples of the adversarial problem](https://ipfs.io/ipfs/QmNrAFz34SLqkzhSg4wAYYJeokfJU5hBEpkT4hPRi226y9). The problem that Google acknowledges, is that it is rather hard to algorithmically reason, wether or not this particular sample is adversarial. That is - independently of how cool the learning technology itself is. Obviously, a cryptoeconomic approach can change beneficiaries in this game effectively, by removing possible sybil attack vectors, and by removing the necessity to make a decision on example crawling and meaning extraction from one entity to the whole world. A learning sybil-resistant model will probably lead to orders of magnitude with more predictive results.

## Cyber protocol at `cyber`
- compute `cyber` inception of cyber protocol based on the Genesis distribution rules
- define knowledge graph state
- take cyberlinks
- check the validity of signatures
- check the bandwidth limit
- check the validity of CIDv0
- if - signatures, bandwidth limit, and CIDv0 are valid - than cyberlink is valid
- every round, calculate cyber•rank delta for the knowledge graph

## Knowledge graph
We represent a knowledge graph as a weighted graph of directed links between - content addresses or content identifications or CIDs or simply IPFS hashes. In this paper we will use them as synonyms.

![knowledge_graph.png](https://ipfs.io/ipfs/QmejVRS9irYb6eXGDZNM9YEuFyb3a5jn4EWh3MRC3LVRij)

Content addresses are essentially web3 links. Instead of using non-obvious and mutable things:
```
https://github.com/cosmos/cosmos/blob/master/WHITEPAPER.md
```
we can, pretty much, use the exact thing:
```
Qme4z71Zea9xaXScUi6pbsuTKCCNFp5TAv8W5tjdfH7yuH
```

Using content addresses for building a knowledge graph we get [the so much needed](https://steemit.com/web3/@hipster/an-idea-of-decentralized-search-for-web3-ce860d61defe5est) superpowers of [ipfs](https://ipfs.io/ipfs/QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps)-[like](https://ipfs.io/ipfs/QmXHGmfo4sjdHVW2MAxczAfs44RCpSeva2an4QvkzqYgfR) p2p protocols for a search engine:

- future proofs for mesh-networks 
- interplanetary communication
- tolerant communication
- accessible communication
- technologicly agnostic

Web3 agents generate our knowledge graph. Web3 agents include themselves to the knowledge graph with just the one transactiion. Thereby, they prove the existence of their private keys for content addresses of the revealed public keys. Using this basic proof mechanism, the consensus computer could have provable differentiation between subjects and objects in a knowledge graph.

Our `cyber` implementation is based on [`cosmos-sdk`](https://github.com/cosmos/cosmos-sdk) identities and [`cidv0`](https://github.com/multiformats/cid#cidv0) content addresses.

Web 3 agents generate knowledge graph by applying `cyberlinks`.

## Cyberlinks
To understand cyberlinks, we need to understand the difference between a `URL link`, aka hyperlink, and an `IPFS link`. A `URL link` points to the location of content, wether an `IPFS link` points to the content itself. The difference of web architecture that is based on location links and content links is drastical, hence it requires a new approach.

`Cyberlink` is an approach to link two content addresses or `IPFS links` semantically:

```
QmdvsvrVqdkzx8HnowpXGLi88tXZDsoNrGhGvPvHBQB6sH
.
Qme4z71Zea9xaXScUi6pbsuTKCCNFp5TAv8W5tjdfH7yuH
```

This `cyberlink` means that cyberd presentation on cyberc0n is referencing to the Cosmos whitepaper. A concept of `cyberlink` is a convention around simple semantics of communication format in any peer to peer network:

`<content-address x>.<content-address z>`

You can see that `cyberlink` represents a link between two links. Easy peasy!

Cyberlink is a simple, yet a powerful semantic construction for building a predictive model of the universe. That means, that by using `cyberlinks` instead of `hyperlinks` gives us the superpowers that were inaccessible to previous architectures of general-purpose search engines.

Cyberlinks can be extended, e.g. they can form chain links if two series of cyberlinks exist from one agent, in which the second link in the first cyberlink is equal to the first link in the second cyberlink:

```
<content-address x>.<content-address z>
<content-address z>.<content-address z>
```

Using this simple principle, all interacting agents can reach consensus around interpreting clauses. That makes chain-links helpful for interpreting any communications around the topic of relevance.

![link_chains.png](https://ipfs.io/ipfs/QmNd15Pa1pzFVv98jmf1uuvPGLxRpptQMVoMidYkj1YvDC)

Also, by using the following link: `QmNedUe2wktW65xXxWqcR8EWWssHVMXm3Ly4GKiRRSEBkn`, one can signal the start and the end of an execution in a knowledge graph. A lot of cool stuff can be done using cyberlinks.

If web3 agents expand native `IPFS links` with something semantically richer, such as [`DURA` links](https://github.com/cybercongress/cyb/blob/dev/docs/dura.md), than those agents can reach consensus in an easier manner, on the rules of program execution. Indeed, the `DURA` protocol is a real implementation of a cyberlink concept.

[`cyber`](https://github.com/cybercongress/cyberd) implementation of `cyberlinks` based on [`DURA`](https://github.com/cybercongress/cyb/blob/dev/docs/dura.md) specification is available in [`.cyber`](https://github.com/cybercongress/.cyber) app of browser [`cyb`](https://github.com/cybercongress/cyb).

Based on `cyberlinks` we can compute the relevance of subjects and objects in a knowledge graph. That is why we need a consensus computer.

## Notion of consensus computer
Consensus computer is an abstract computing machine that emerges from agent interaction.

A consensus computer has capacity in terms of fundamental computing resources, such as memory and computing. To interact with agents, a computer needs bandwidth.

Ideal consensus computer is a computer where:

```
the sum of all computations and memory available for individuals
is equal to
the sum of all verified computations and memory of a *consensus computer*
```

We know that:

```
verifications of computations < computations + verifications of computations
```

Hence we will never be able to achieve an ideal consensus computer. The CAP theorem and the scalability trilemma also add proof to this statement.

However, this theory can work as a performance indicator of a consensus computer.

After 6 years of investing into consensus computers, we realize that [Tendermint](https://ipfs.io/ipfs/QmaMtD7xDgghqgjN62zWZ5TBGFiEjGQtuZBjJ9sMh816KJ) consensus has a good balance between the coolness required for our task, and the readiness for its production. So, we decide to implement the `cyber` protocol using Tendermint, which is very close to the Cosmos Hub settings.

The `cyber` implementation is a 64-bit tendermint consensus computer, relevant to 64-byte string space that is as far from ideal. At least as "1/146". This is because we have 146 validators who verify the same computation using the knowledge graph of the same size.

We must bind the computation, storage and the bandwidth supply of the consensus computer with a maximized demand for queries. Computation and storage in case of a basic relevance machine, can be easily predicted based on bandwidth, but bandwidth requires a limiting mechanism.

## Bandwidth
`Bonded stake` - is a stake, that deducted from your acc coins and be put as a deposit to take part in consensus. Due to the passive inflation model and slashing, the deposit will not match 1-to-1 to the final reward. So, for example, stakeholders may wish to set up a script that will periodically withdraw and rebound rewards to increase their bonded stake.

`Active stake` - currently available for direct transfer; non-bonded stake.

`Bandwidth stake` = `active stake + bonded stake`.

Cyberd uses a very simple bandwidth model. The main goal of that model is to reduce daily network growth to a given constant, say 3gb per day.

Thus, here we introduce `resource credits`, or RS. Each message type has assigned RS cost. There is a constant `DesirableNetworkBandwidthForRecoveryPeriod` determened desirable for `RecoveryPeriod` of spent RS value. `RecoveryPeriod` is defining how fast agent can recover their bandwidth from 0 to max bandwidth. An agent has maximum RS proportional to his stake, by the following formula:

`agent_max_rc = bandwidth_stake * DesirableNetworkBandwidthForRecoveryPeriod`

There is an `AdjustPricePeriod`, summing how much RS was spent for that period `AdjustPricePeriodTotalSpent`. Also, there is a constant `AdjustPricePeriodDesiredSpent`, used to calculate the network load.

The `AdjustPricePeriodTotalSpent / AdjustPricePeriodDesiredSpent` ratio is called "fractional reserve ratio". If network usage is low, the fractional reserve ratio adjusts message cost (via the use of simple multiplication) to allow agents with a lower stake to execute more transactions.
If resource demand increases than the fractional reserve ratio goes `>1`, thus increasing messaging cost and limiting final tx count for a long-term period (RC recovery will be `<` than RC spending).

There are only two ways to change the account bandwidth stake:

1. Direct coin transfer
2. When distribution payouts occur 
For example, when a validator changes his commission rates, all delegations will be automatically unbounded. Another example is when the delegator himself unbonds a part or his full share.

The agents must have CYB tokens in accordance with their will to learn the knowledge graph. However, the proposed mechanics of CYB tokens works not only as spam protection, but also as the economic regulation mechanism, that aligns the ability of validators to process knowledge graphs and markets the demand for processing.

## Relevance machine
We define the relevance machine as a machine that transitions a knowledge graph state, based on the will of agents to learn the knowledge graph. The more agents that learn the knowledge graph, the more valuable the graph itself becomes. This machine enables simple construction for search questions, querying and answers delivering.

The will to learn is projected on every agent's cyberlink. A simple rule prevents abuse by agents: one content address can be voted on by a coin, only once. So, for ranking purposes, it does not matter from how many accounts one has voted. Only the sum of their balances matters.

A useful property of a relevance machine is that it must have inductive reasoning properties or follow the blackbox principle.

```
It must be able to interfere predictions
without any knowledge about objects
except those cyberlinked, when cyberlinked, and that were cyberlinked.
```

If we assume that a consensus computer must have some information about linked objects, the complexity of such a model will grow unpredictably, hence, there must be a requirement for memory and computation for that computer. That is, deduction of meaning inside the consensus computer is expensive, thus, our design depends on the "blindness of assumption". Instead of deducting meaning inside the consensus computer, we design a system, in which meaning extraction is incentivized - because the agents need CYB in order to compute relevance.

Also, thanks to content addressing the relevance machine following the blackbox principle does not need to store data, but, can effectively operate on it.

Human intelligence is organized in such a way, that the pruning of none-relevant, and none-important memories - passes with time. The same can be applied to the relevance machine. Another useful property of the relevance machine is that it needs to store neither the past state, nor the full-current state in order to remain useful. Or more precisely: _relevant_. So, the relevance machine can implement [aggressive pruning strategies](QmP81EcuNDZHQutvdcDjbQEqiTYUzU315aYaTyrVj6gtJb) such as pruning all the history of a knowledge graph or "forgetting" links that become non-relevant.

`cyber` implementation of the relevance machine is based on the most straightforward mechanism which is called - cyber•Rank.

## cyber•Rank
Ranking using consensus computers is difficult because consensus computers bring serious resources bounds, e.g. [Nebulas](https://ipfs.io/ipfs/QmWTZjDZNbBqcJ5b6VhWGXBQ5EQavKKDteHsdoYqB5CBjh) still fails to deliver something useful on-chain. At first, we must ask ourselves: why do we need to compute and store on-chain ranking, rather then do it the [Colony](https://ipfs.io/ipfs/QmZo7eY5UdJYotf3Z9GNVBGLjkCnE1j2fMdW2PgGCmvGPj) way or the [Truebit](https://ipfs.io/ipfs/QmTrxXp2xhB2zWGxhNoLgsztevqKLwpy5HwKjLjzFa7rnD) way?

If a rank is computed inside a consensus computer, one has an easy content distribution of the rank, as well as, an easy way to build provable applications on top of the rank. Hence, we decided to follow more cosmic architecture. In the next section, we describe the proof of relevance mechanism which allows the network to scale with the help of domain-specific relevance machines that work in parallel, thanks to the IBC protocol.

Eventually, the relevance machine needs to find (1) a deterministic algorithm, that will allow for a computing a rank of a continuously appended network, that can scale the consensus computer to orders of magnitude of the likes of Google. A perfect algorithm (2) must have linear memory and computation complexity. Most importantly, it must have the (3) highest provable prediction capabilities for the existence of relevant cyberlinks.   

After [some research](https://arxiv.org/pdf/1709.09002.pdf), we found that we can not find the silver bullet here. So, we decided to find a more basic bulletproof way that will bootstrap the network: [the rank](http://ilpubs.stanford.edu:8090/422/1/1999-66.pdf) from which Larry and Sergey have bootstrapped a previous network. The key problem with the original PageRank is that it is not resistant to sybil attacks.

Token weighted [PageRank](http://ilpubs.stanford.edu:8090/422/1/1999-66.pdf) is limited by token-weighted bandwidth and do not have inhereted problems, such as the naive PageRank, and, are resistant to sybil attacks. For the time being, we will call it cyber•Rank until something better will emerge.

In the center of the spam protection system, is an assumption that 'write' operations can be only executed by those who have a vested interest in the evolutionary success of the relevance machine. Every 1% of the stake in consensus computer, gives the ability to use 1% of possible network bandwidth and computing capabilities.

As nobody uses all the possessed bandwidth, we can safely use up to 100x of the fractional reserves with a 2-minute recalculation target. This mechanics offers a discount for cyberlinking, thus effectively maximizing demand for it.

We would love to discuss the problem of vote-buying. Vote-buying by itself is not such a bad thing. The problem with vote buying appears in systems where voting affects the allocation of the inflation in that system. For example, [Steem](QmepU77tqMAHHuiSASUvUnu8f8ENuPF2Kfs97WjLn8vAS3) or any other state-based systems. Vote buying can become easily profitable for adversary that employs a zero-sum game, without a necessity to add value. Our original idea of a decentralized search was based on this approach. But, we have rejected that idea, removing the incentive on consensus level for the knowledge graph formation. In our setting, in which every participant must add some value to the system in order to affect the predictive model, vote-buying becomes NP-hard problem, hence it is useful for the system.

We understand that the ranking mechanism will always remain a "red herring". That is why we expect to rely on on-chain governance mechanisms in order to define a winning one. We would love to switch from one algorithm to another, based on economic a/b testing through hard spoons of domain-specific relevance machines thought.  

The current implementation of the consensus computer, is based on the relevance machine for cyber•Rank, that can answer and deliver relevant results for any given search request in a 64 byte CID space. However, it is not enough, to build a network of domain-specific relevance machines. Consensus computers must have the ability to prove relevance to each other.

## Proof of relevance
We design a system under the assumption that such a thing as bad behaviour, does not exist  in regards to search. That is, because anything bad cannot be in the intention of finding answers, but ony in the asnwers themselves. Also, this approach significantly reduces attack surfaces.

> Ranks are computed on the only fact that something has been searched for. Thus linked, and as a result - affected the predictive model.

A good analogy can be observed in quantum mechanics. That is why we do not require such things as negative voting. By doing this, we remove subjectivity out of the protocol and can define proof of relevance.

```
Rank state =
rank values stored in a one-dimensional array
and merkle tree of those values
```

Each new CID gets a unique number. The number starts with zero and increments by one, for each new CID. This is done, so that we can store rank in a one-dimensional array, where indices are the CID numbers.

A Merkle Tree is calculated based on the [RFC-6962 standard](https://tools.ietf.org/html/rfc6962#section-2.1). Since rank is stored in a one-dimensional array where indices are CID numbers (we can say that it is "brought to order" by CID numbers), this leaves in the merkle tree, from left to right, `SHA-256` hashes of rank value. The index of the leaf is the CID number. It helps to easily find proofs for specified CID (`log n` iterations, where `n` is a number of leaves).

Storing a merkle tree is necessary, im order to split the tree into subtrees with a number of leaves multiplied to the power of 2. The smallest one is obviously a subtree with only the one leaf (and therefore its `height == 0`). Any new leaf addition will look as follows: each new leaf is added as a subtree with a `height == 0`. Then, it is sequentially merged with subtrees, with the same `height`, from right to left.

Example:


      ┌──┴──┐   │      ┌──┴──┐     │         ┌──┴──┐     │   │
    ┌─┴─┐ ┌─┴─┐ │    ┌─┴─┐ ┌─┴─┐ ┌─┴─┐     ┌─┴─┐ ┌─┴─┐ ┌─┴─┐ │
    (5-leaf)         (6-leaf)              (7-leaf)

To obtain the merkle root hash - join the subtree roots from right to left.

"Rank merkle tree", can be stored differently:

_Full tree_ - all subtrees with all leaves and intermediary nodes  
_Short tree_ - contains only subtrees roots

The trick is that a _full tree_ is only necessary for providing merkle proofs. For consensus purposes and updating a tree, it's enough to have a _short tree_. In order to store a merkle tree in a database, one can use a _short tree_. Marshaling of a short tree with `n` subtrees (each subtree takes 40 bytes):  

```
<subtree_1_root_hash_bytes><subtree_1_height_bytes>
....
<subtree_n_root_hash_bytes><subtree_n_height_bytes>
```

For `1,099,511,627,775` leaves a _short tree_ would contain only 40 subtrees roots and take only 1600 bytes.

Let us denote rank state calculation:

`p` - rank calculation period  
`lbn` - last confirmed block number  
`cbn` - current block number  
`lr` -  length of rank values array  

For rank storing and calculations we have two separate in-memory contexts:

1. Current rank context. 
It includes the last calculated rank state (values and merkle tree) plus all links and agent stakes submitted to the moment of this particular rank submission.
2. New rank context. 
Currently calculated (or already calculated and waiting for a submission) rank state. Consists of new calculated rank states (values and merkle tree) plus new incoming links and updated agent stakes.

Calculation of a new rank state happens once per `p` blocks and happens in parallel.

The iteration starts from a block number that `≡ 0 (mod p)` and goes utill the next block that `≡ 0 (mod p)`.

For blocks numbered `cbn ≡ 0 (mod p)` (including block numbered 1, 'cause in cosmos blocks starts from 1):

1. Check if the rank calculation is finished. If yes then go to (2.) if not - wait till calculation is finished
(actually this situation should not happen because it means that the rank calculation period is too short).
2. Submit rank, links and agent stakes from new rank context to current rank context.
3. Store the latter calculated rank merkle tree root hash.
4. Start a new rank calculation in parallel (on links and stakes from current rank context).

For each block:

1. All links go to a new rank context.
2. New coming CIDs gets their rank equals to zero. We could do this by checking the last CIDs number and `lr` (it equals to the number of CIDs that already have rank). Then, add the CIDs with the number `>lr` to the end of this array with a value that equals to zero.
3. Update the current context merkle tree with CIDs from the previous step.
4. Store the latest merkle tree from current context (let us name it "last block merkle tree").
5. Check if the calculation of the new rank had finished. If yes go to (4.) if not go to next block.
6. Push calculated rank state to the new rank context. Store merkle tree of the newly calculated rank.

To sum up: in _current rank context_, we have a rank state from the last calculated iteration (plus, with every block, it updates with new CIDs). Moreover, we have links and agent stakes that are participating in the current rank calculation iteration (whether it is finished or not). The _new rank context_ contains links and stakes that will roll over to the next rank calculations, and newly calculated rank state (if a calculation is finished) that is waiting for submitting.

If we need to restart a node firstly, we need to restore both contexts (the current and the new).
Load links and agent stakes from a database using different versions:
1. Links and stakes from last calculated rank version `v = lbn - (lbn mod n)` go to the current rank context.
2. Links and stakes between versions `v` and `lbn` go to new rank context.

Also to restart the node correctly, we have to store the following entities in the database:

1. The last calculated hash of a rank (merkle tree root)
2. A newly calculated rank - short merkle tree
3. Last blocks short merkle tree

With _last calculated rank hash_ and _newly calculated rank merkle tree_ we could check if the rank calculation was finished before the restart of a node. If they are equal, then the rank wasn't calculated, and we should run the rank calculation. If not, we could skip rank calculation and use _newly calculated rank merkle tree_ to participate in consensus when it comes to block number `cbn ≡ 0 (mod p)` (rank values will not be available until rank calculation happens in the next iteration. The validator can still participate in consensus, so nothing bad happens).

_Last block merkle tree_ is necessary to participate in consensus till the start of the next rank calculation. After the restart we could end up with two states:
1. Restored current rank context and new rank context without rank values (links, agent stakes, and merkle tree).
2. Restored current rank context without rank values. Restored new rank context only with links and agent stakes.

A node can participate in consensus but cannot provide rank values (and merkle proofs) 'till two iterations of rank calculation have finished (the current and the next one). Search index should be run in parallel, and should not influence the work of the consensus machine. The validator should be able to turn off index support.

Now we have proof of rank of any given content address. While the relevance is still subjective by nature, we have a collective proof that something was relevant for a certain community at some point in time.

```
For any given CID it is possible to prove the relevance
```

Using this type of proof any two [IBC compatible](https://ipfs.io/ipfs/QmdCeixQUHBjGnKfwbB1dxf4X8xnadL8xWmmEnQah5n7x2) consensus computers can proof the relevance to each other so that domain-specific relevance machines can flourish. Thanks to inter-blockchain communication protocol, one can basically, either launch an own domain-specific search engine by forking cyberd, which is focused on the _common public knowledge_, or, plug cyberd as a module into an existing chain, e.g. Cosmos Hub.  Hence - our search architecture, based on domain-specific relevance machine, can learn from common knowledge.

![rm-network.png](https://ipfs.io/ipfs/QmdfgdkaU8CKXD7ow983vZ2LjJjz8Um9JA5buwQ1aaXT6Q)

In our relevance for common `cyber` implementations of proof of relevance root, a hash is computed on Cuda GPUs every round.

## Speed and scalability
We need quick confirmation times to feel like a usual web app. It is a strong architecture requirement that shapes an economic topology and scalability of the cyber protocol.

Proposed blockchain design is based on the [Tendermint consensus](https://ipfs.io/ipfs/QmaMtD7xDgghqgjN62zWZ5TBGFiEjGQtuZBjJ9sMh816KJ) algorithm with 146 validators and has a very fast - 2 second finality. Average confirmation timeframe, at half a second with asynchronous interaction make complex blockchain search almost invisible to agents.

Let us assume, that our node implementation based on `cosmos-sdk` can process 10k tps. Thus, every day at least 8.64 million agents can submit 100 cyberlinks each, and impact results simultaneously. That is enough to verify all the assumptions in the world. As blockchain technology evolves we want to check that every hypothesis works before scaling it any further. Moreover, the proposed design needs demand for full bandwidth in order for the relevance to becomes valuable. That is why we strongly focus on accessible, but provable distribution from inception.

## Approach toward distribution
While designing the initial distribution structure for Cyber protocol we aimed to achieve the following goals:
- Develop a provable and transparent distribution, in accordance with the best industry practices
- Allow equal participation irrespectively of political, regulatory or any other restrictions, which may be imposed by outside agents
- Prevent attacks on privacy, such as - installment of KYC requirements
- Spread distribution over time in order to grant equal access to all agents to initial distribution; without any limitations such as - hard caps or any other restrictions
- Honour genesis Cosmos investors, for the development of a technology which made it possible for a simplified development of Cyber protocol
- Attract the most professional validators from the Cosmos ecosystem for bootstrapping the network
- Allow easy early access for active agents of the Ethereum ecosystem in order to accelerate the growth of the knowledge graph and solve the chicken and egg problem
- Decentralize the management of auction donations, starting from day 0
- Honour 30 months of cyber•Congress R&D and the community behind it

The goal of creating an alternative to a Google-like structure requires extraordinary effort of different groups. So, we have decide to set up cyber•Foundation as fund managed, via a decentralized engine such as - Aragon DAO filled with ETH and managed by agents who participated in the initial distribution. This approach will allow for safeguarding from an excessive market dumping of the platforms native CYB tokens in the first years of work, thereby ensuring stable development. Additionally, this allows to diversify the underlying platform and extend the protocol to other consensus computing architecture should the need arise.

While choosing tokens for donations we followed three main criteria: the token must be (1) one of the most liquid on the market, (2) the most promising, so a community can secure a solid investment bag to be competitive even when compared to giants like Google and (3) have the technical ability to execute an auction and resulting the organization without relying on any third party. The only system that matches this criteria, is - Ethereum, hence the primary token of donations will be ETH. That is why we have decided to create 2 tokens: THC and CYB:

- THC is a creative cyber proto substance. THC being an Ethereum ERC-20 compatible token, has utility value in the form of control over  the cyber•Foundations (Aragon DAO) ETH funds from the initial auction proceeds. THC was emitted during the creation of cyber•Foundation as an Aragon organization. Creative powers of THC came from ability to receive 1 CYB per each 1 THC for locking it during cyber•Auction.
- CYB is the native token of sovereign Cyber protocol under the Tendermint consensus algorithm. It also has 2 primary uses: (1) it is staked for consensus and (2) it is bandwidth limiting for submitting links and when computing rank.

Both tokens remain functional and will track value independently due to very different utility by nature.

The initial distribution happens in 3 different epochs, by nature and goals, and are spread in a time-frame of almost 2 years:

1. Pre-genesis: 21 days. Is needed to launch cyber protocol in a decentralized fashion with independent genesis validators.

2. Genesis: 21 days. Launch of the main net and Game of Thrones: needed in order to involve the most active Ethereum and Cosmos crypto players into an engaging game with the ability to fully understand, and test the software, using the main network and incentives. Game of Thrones is a sort of a distribution game, with some discount for the attraction of the necessary critical mass, in order to make it possible of learning the early knowledge graph, by the most intelligent community members.

3. Post-Genesis: 600 days. Continuous distribution of CYB based on cyber•Auction proceeds. It is needed in order, to dilute the initial distribution and evolve into an existing crypto community, and beyond...

## Pre-genesis
2 distribution events have happened prior to Genesis:

1. 700 000 000 000 000 THC tokens are minted by cyber•Foundation. Allocations of THC tokens is the following:

- 100 000 000 000 000 THC is allocated to cyber•Congress
- 600 000 000 000 000 THC is allocated to cyber•Auction contract

2. At the start of `euler-5` donation round in ATOMs started. Purpose of this round is to involve real validators at a genesis. 5% of CYB will be allocated to participants of this donation round.

## Genesis and Game of Thrones
Genesis of `cyber` protocol will contains 1 000 000 000 000 000 CYB (One Quadrillion CYB) broken down as follows:

- 600 000 000 000 000 CYB under multisig managed by cyberCongress for manual distributions during cyber•Auction for those who stake THC until the end of cyber•Auction
- 200 000 000 000 000 CYB under multisig managed by cyberCongress: the Game of Thrones for ATOM and ETH holders, 100 TCYB for each.
- 100 000 000 000 000 CYB for top 80% ETH holders by stake excluding contracts
-  50 000 000 000 000 CYB as a drop for all ATOM stakeholders
-  50 000 000 000 000 CYB for pre-genesis contributors in ATOM

Game of Thrones - is a game between ATOM and ETH holders for being the greatest. As a result of a 21-day auction after Genesis, every community will earn 10% of CYB. In order to make the game run smoothly we concisely adding arbitrage opportunity in the form of significant discount to ATOM holders because the system needs provably professional validators and delegators at the beginning, and basically employs them for free.

We can describe the discount in the following terms: 
The current buying power of all ATOM against all ETH based on current (market) caps is about 1/24. Given that 10% of CYB will be distributed based on donation in ATOM and 10% of CYB will be distributed based on donations in ETHб the discount for every ATOM donation during Game of Thrones is about 24x, which is significant enough to encourage participation based on arbitrage opportunity during the first 21 days of Genesis auction and stimulate the price of ATOM as an appreciation for all cosmic community.

## Post-genesis and cyber•Auction
The post Genesis stage is called cyber•Auction. It starts after the end of the Game of Thrones and lasts for 600 rounds, 23 hours each. During this phase, CYB are continuously distributed based on locked THC bought in the continuous auction.

The role of cyber•Auction is twofold:
- It creates a non-exlusive long lasting and provable game of initial distribution without necessity to spend energy on proof of work. It is crucial, that the early knowledge graph was created, in some sense - fairly, by an engaged community, which was formed during a non-exclusive game.
- As a result of the auction, the community will has access to all the raised resources under the Aragon organisation. We believe in a true decentralized nature of the thing we create, so we do not want to grab all the money from the funding, as we have already funded the creation of the system by ourselves, and we kindly ask for a fair 10% CYB cut from the pre-genesis investors, founders and developers. Competing with Google is challenging and will be more viable if the community will sit on a bag of ever-growing ETH. Given the current growth-rate of ETH, this bag can become very compelling in some years to come. Also, this bag can be the source of alternative implementation of the protocol if Cosmos-based system will fail, or, in the case the community just wants to diversify technology the involved, e.g. ETH2, Polkadot or whatever else.

After genesis, CYB tokens can only be created by validators, based on staking and slashing parameters. The basic consensus is that newly created CYB tokens are distributed to validators as they do the essential work to make the relevance machine run, both in regards to energy consumed for the computation, and the cost for storage capacity. This means that stakeholders decide where the tokens can flow further.

After Genesis inflation is adjusted using the `TokensPerBlock` parameter. Given that the network has a 2 second target block and ~7% target inflation, the starting parameter will be 50 MCYB.

There is currently, no such thing, as the maximum amount of CYB due to continuous inflation paid to the validators. Currently, CYB is implemented using 64int so the creation of more CYB makes significantly more expensive to compute state changes and rank. We expect that a lifetime monetary strategy must be established by a governance system, after the complete initial distribution of CYB tokens and the activation of functionality of smart contracts.

The following rules apply to CYBs under cyber•Auction multisig:
- it will not delegate its stake and as a result it will remain as a passive stake until it becomes distributed
- after the end of cyber•Auction, all remaining CYB tokens will be provably burned

## Role of ATOMs
Overall 15% of CYB will be distributed based on donations in ATOM during 2 rounds:
-  50 000 000 000 000 CYB for genesis ATOM contributors
- 100 000 000 000 000 CYB for ATOM contributors at the start of Smith Epoch

All ATOM donations go to cyber•Congress multisig. The role of ATOM donations is the following: thanks to ATOM we want to secure a lifetime commitment of cyber•Congress in the development of Cosmos and Cyber ecosystems. ATOM donations to cyber•Congress will allow us to use the staking rewards for continuously funding of the Cyber protocol without the necessity to dump CYB tokens.

## Learning the Graph
We assume that the proposed algorithms does not guarantee high quality knowledge by default. Like a child it needs to learn in order to prosper. The protocol itself provides only one simple tool: the ability to create a cyberlink with a certain weight between two content addresses.

Analysis of the semantic core, behavioral factors, anonymous data about the interests of agents and other tools that determine the quality of the search can be done in smart contracts and off-chain applications, such as web3 browsers, decentralized social networks and content platforms. So it is the goal of the community and the agents to build an initial knowledge graph and maintain it to provide the most relevant search.

We suggest that the knowledge graph could be created using the following methods:
- agents can use different applications, like [cyb.virus](https://github.com/cybercongress/cyb-virus) chrome extension, in which agents save their content or any content they like, into IPFS, and link it with relevant keywords
- validators and search enthusiasts use automatic scripts to create and index popular resources, such as Wikipedia - or transfer and edit current search result from Google and other search engines
- In the process of informational search by agents in open source application (like the cyb web3 browser or any other), the application itself analyzes its current rank, the semantic core of the content, the behavior of the agent, when searching and viewing content  (agent transitions between content, bounce rate, viewing time and others). Based on this data, and by using the balance of agent CYB tokens, the application automatically, or, with agent input - links content. Agents will choose which open-source application to use for search and how it will affect the rank. In this case, creating a knowledge graph ceases to be a black box and becomes absolutely transparent. Data about the agents search queries and their behavior remains on their devices, ensuring anonymity. This data can be synced between devices by agent's private instance of the `cyber` chain.
- Existing search engines, such as Google, use thousands of professional assessors to improve the quality of search rank. Validators have the maximum economic interest in the entire protocol and can voluntarily allocate part of their commission in CYB tokens to smart contract that pay for professional search assessors. Each validator, can provide an application for assessor to make a search more relevant. Applications index all transactions and stores information about changes of weights in the knowledge graph. Assessors receive assignments to check the relevance of cyberlinks and, if necessary, they increase the weight of cyberlinks. The validator automatically checks for increased or decreased cyberlink weights with which assessor has worked. If the assessor increased the weight of the link and the final weight due to the actions of the other agents also increased, then the cyberlink was more relevant and the assessor receives a reward from the validator.

We are confident that the measures described above will allow us to build an effective search, created by "the people" for "the people", without the need to transfer their personal data to intermediaries, such as Google or/and Facebook.

## In-browser implementation
We wanted to imagine how that could work in a web3 browser. To our disappointment we [were not able](https://github.com/cybercongress/cyb/blob/master/docs/comparison.md) to find a web3 browser that can showcase the coolness of the proposed approach in action. That is why we decide to develop the web3 browser [cyb](https://github.com/cybercongress/cyb/blob/master/docs/cyb.md) that has a sample application ".cyber" for interacting with the `cyber://` protocol.

![search-main](https://user-images.githubusercontent.com/410789/60329155-151be200-9990-11e9-8d78-e32e72285abc.png)

As another good example, we created [a Chrome extension](https://github.com/cybercongress/cyb-virus) that allows anybody to pin any web page to IPFS and index it by any keywords, thus making it searchable.

Current search snippets are ugly, but we expect that they can be easily extended using IPLD for different types of content so that they can be even more perfect than that of Google.

During the implementation of the proposed architecture, we realized at least 3 key benefits that Google, probably, would not be able to deliver with its conventional approach:
- the search result can be easily delivered from a p2p network right into search results: eg. .cyber can play video.
- payment buttons can be embedded right into search snippets, so web3 agent can interact with search results, e.g. an agent can buy an item right in `.cyber`. So e-commerce can flourish because of transparent conversion attributes.
- search snippets must not be static but can be interactive, e.g. `.cyber` will eventually answer location-based answers.

## Roadmap
We foresee the demand for the following protocol features, that the community could work on after launch:
- Parametrization
- Universal oracle
- IBC
- WASM VM for gas
- Onchain upgrades
- CUDA VM for gas
- Privacy by default

## Applications of the knowledge graph
A lot of cool applications can be built on top of the proposed architecture:

_Web3 browsers_: 
Its easy to imagine the emergence of a full-blown blockchain browser. Currently, there are several efforts for developing browsers around blockchains and distributed tech. Among them are: Beaker, ~~Mist~~, Brave, Saifmade and Metamask. All of them suffer from trying to embed web2 into web3. Our approach is a bit different. We consider web2 as an unsafe subset for web3. That is why we have decided to develop a web3 browser that can showcase the cyber approach to answer questions better.

_Programmable semantic cores_: 
Currently, the most popular keywords in a gigantic semantic core of Google, are keywords of apps such as youtube, facebook, github, etc. However, developers have a very limited possibility to explain Google how to better structure results. The cyber approach brings this power back to developers. On any given input string in any application, relevant answer can be computed either globally, in the context of an app, an agent, geo-position, or in all of the above combined.

_Search actions_: 
The proposed design enables native support for blockchain asset related activity. It is trivial to design applications which are (1) owned by creators, (2) appear right in search results, and (3) allow for a transactable call to action with (4) provable attribution of a conversion to search query. e-Commerce has never been so easy for anybody.

_Offline search_: 
IPFS makes it possible for easy retrieval of documents from surroundings without a global internet connection. cyberd can itself be distributed using IPFS. That creates a possibility for ubiquitous offline search.

_Command tools_: 
Command line tools can rely on relevant and structured answers from a search engine. That practically means that the following CLI tool is possible to implement

```
>  cyberd earn using 100 gb hdd

Enjoy the following predictions:
- apt install go-filecoin:     0.001   BTC per month per GB
- apt install siad:            0.0001  BTC per month per GB
- apt install storjd:          0.00008 BTC per month per GB

According to the best prediction, I made a decision try `mine go-filecoin`

Git clone ...
Building go-filecoin
Starting go-filecoin
Creating a wallet using @xhipster seed
You address is ....
Placing bids ...
Waiting for incoming storage requests ...

```
A search from CLI tools will inevitably create a highly competitive market of a dedicated semantic core for bots.

_Autonomous robots_: 
Blockchain technology enables the creation of devices which can earn, store, spend and invest digital assets by themselves.

> If a robot can earn, store, spend and invest, it can do everything you can!

What is needed, is a simple yet a powerful state-reality tool with the ability to find particular things. cyberd offers minimalistic but continuously self-improving data source, that provides the necessary tools for programming economically rational robots. According to [top-10000 english words](https://github.com/first20hours/google-10000-english) the most popular word in English is (the - ironically) defined article `the`, which means a pointer to a particular thing. This fact can be explained as the following: particular things are the most important to us. So, the nature of our current semantic computing is to find unique things. Hence, the understanding of unique things will become essential for robots also.

_Language convergence_: 
A programmer should not care about what language does an agent use. We don't need to know in what language an agent is searching in. The entire UTF-8 spectrum is at work. A semantic core is open, so competition for answering can become distributed across different domain-specific areas, including semantic cores of different languages. This unified approach creates an opportunity for cyber•Bahasa. Since the Internet, we observe a process of rapid language convergence. We use (more) truly global words across the entire planet independently of our nationality, language and race, Name the Internet. The dream of a truly global language is hard to deploy because it is hard to agree on what means what. However, we have the tools to make that dream come true. It is not hard to predict that the shorter a word, the higher its cyber•rank will be. Global publicly available list of symbols, words, and phrases sorted by cyber•rank with corresponding links provided by cyberd can be the foundation for the emergence of genuinely global language everybody can accept. Recent [scientific advances](https://ipfs.io/ipfs/QmQUWBhDMfPKgFt3NfbxM1VU22oU8CRepUzGPBDtopwap1) in machine translation are breathtaking but meaningless for those who wish to apply them without Googles scale-trained model. cyber•rank offers precisely this.

This is sure not the full list of possible applications, though, a very exciting one indeed.

## Apps on top of knowledge graph
Our approach to the economics of consensus computing is that agents buy a certain amount of RAM, CPU, and GPU as they want to execute programs. OpenCypher or a GraphQL like language can be provided to explore the semantics of the knowledge graph. The following list is a list of simple programs [we can envision](https://medium.com/@karpathy/software-2-0-a64152b37c35) that can be built on top of a simple relevance machine with the support of on-chain WASM-like VM.

_Self prediction_: 
A consensus computer that can continuously build a knowledge graph by itself predicting the existence of cyberlinks and applying these predictions to a state of itself. Hence, a consensus computer can participate in the economic consensus of the cyber protocol.

_Universal oracle._:
A consensus computer that can store the most relevant data in a key-value store, where the key is CID, and the value are bytes of the actual content. It is done by making a decision on every round, on which CID value it wants to prune and which value it wants to apply, based on the utility measure of content addresses in the knowledge graph. To compute utility measure validators check availability and the size of the content for the top-ranked content addresses in the knowledge graph, then weights on the size of CIDs and its rank. The emergent key-value store will be available to write for the consensus computer only and not for agents, but values can be used in programs.

_Proof of location_: 
It is possible to construct cyberlinks with proof-of-location based on some existing protocol such as [Foam](https://ipfs.io/ipfs/QmZYKGuLHf2h1mZrhiP2FzYsjj3tWt2LYduMCRbpgi5pKG). Location-based search can also become provable if web3 agents will mine triangulations and attach proof-of-location for every chain link.

_Proof of web3 agent_: 
Agents are a subset of content addresses with one fundamental property: consensus computer can prove the existence of private keys for content addresses for the subset of a knowledge graph, even if those addresses had never transacted on their own on the chain. Hence, it is possible to compute much of provable stuff on top of that knowledge, e.g. some inflation can be distributed to addresses that have never transacted in the cyber network, but have a provable link.

_Motivation for read requests_: 
It would be great to create cybernomics not only to write requests to consensus computers, but to read those requests too. So, read requests can be "two order of magnitude" cheaper, but guaranteed. Read requests to a search engine can be provided by the second tier of nodes which earn CYB tokens in state channels. We consider implementing state channels based on HTLC and proof verification which unlocks amount earned for already served requests.

_Prediction markets on link relevance_: 
We can move the idea further, by the ranking of knowledge graph based on prediction market on links relevance. An app that allows betting on link relevance can become a unique source of truth for the direction of terms, as well as, to motivate agents to submit more links.

_Private cyberlinks_: 
Privacy is foundational. While we are committed to privacy, achieving implementation of private cyberlinks is unfeasible for our team up to Genesis. Hence, it is up to the community to work on WASM-like programs that can be executed on top of the protocol. The problem is to compute cyberRank-based on cyberlink submitted by a web3 agent, without revealing neither previous request nor the public keys of a web3 agent. Zero-knowledge proofs, in general, are very expensive. We believe that the privacy of search should be a must by design, but not sure that we know how to implement it. [Coda](https://ipfs.io/ipfs/Qmdje3AmtsfjX9edWAxo3LFhV9CTAXoUvwGR7wHJXnc2Gk) like recursive snarks and [mimblewimble](https://ipfs.io/ipfs/Qmd99xmraYip9cVv8gRMy6Y97Bkij8qUYArGDME7CzFasg) constructions, in theory, can solve part of the privacy issue, but they are new, untested and anyways will be more expensive in regarding to computation, than a transparent alternative.

## Conclusion
We define and implement a protocol for provable communications of consensus computers on relevance. The protocol is based on a simple idea of content, defined by knowledge graphs, which are generated by web3 agents by using cyberlinks. Cyberlinks are processed by a consensus computer using a concept we call a relevance machine. `cyber` consensus computer is based on `CIDv0` and uses `go-ipfs` and `cosmos-sdk` as its foundation. IPFS provides significant benefits in regards to resource consumption. CIDv0 as primary objects are robust in their simplicity. For every CIDv0 cyber•rank is computed by a consensus computer with no single point of failure. Cyber•rank is CYB weighted PageRank, with economic protection from sybil attacks and selfish voting. Every round, a merkle root of the rank tree is published so every computer can prove to any other computer - a relevance value for a given CID. Sybil resistance is based on bandwidth limiting. Embedded ability to execute programs offer inspiring apps. The starting primary goal, is the indexing of peer-to-peer systems with self-authenticated data, either stateless, such as - IPFS, Swarm, DAT, Git, BitTorrent, or, stateful, such as - Bitcoin, Ethereum and other blockchains (and tangles). The proposed semantics of linking architecture, offers a robust mechanism for predicting meaningful relations between objects by a consensus computer itself. The source-code of a relevance machine is open-source. Every bit of data accumulated by a consensus computer is available for everybody, if one has resources to process it. The performance of the proposed software implementation is sufficient for seamless agent interactions. Scalability of the proposed implementation is enough to index all the self-authenticated data that exist today, and serve it to millions of web3 agents. The blockchain is managed by a decentralized autonomous organization which functions under the Tendermint consensus algorithm with a standard governance module. Though a system provides the necessary utility to offer an alternative for conventional search engines, it is not limited to this use-case alone. The system is extendable for numerous applications, e.g. making it possible to design economically rational self-owned robots that can autonomously understand objects around them.

## References
- [cyberd](https://github.com/cybercongress/cyberd)
- [Scholarly context adrift](https://ipfs.io/ipfs/QmNhaUrhM7KcWzFYdBeyskoNyihrpHvUEBQnaddwPZigcN)
- [Web3 stack](https://ipfs.io/ipfs/Qmf3eHU9idMUZgx6MKhCsFPWL24X9pDUi2ECqyH8UtBAMQ)
- [Search engines information retrieval in practice](https://ipfs.io/ipfs/QmeS4LjoL1iMNRGuyYSx78RAtubTT2bioSGnsvoaupcHR6)
- [Motivating game for adversarial example research](https://ipfs.io/ipfs/QmNrAFz34SLqkzhSg4wAYYJeokfJU5hBEpkT4hPRi226y9.ifps)
- [An idea of decentralized search](https://steemit.com/web3/@hipster/an-idea-of-decentralized-search-for-web3-ce860d61defe5est)
- [IPFS](https://ipfs.io/ipfs/QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps)
- [DAT](https://ipfs.io/ipfs/QmXHGmfo4sjdHVW2MAxczAfs44RCpSeva2an4QvkzqYgfR)
- [cosmos-sdk](https://github.com/cosmos/cosmos-sdk)
- [CIDv0](https://github.com/multiformats/cid#cidv0)
- [Thermodynamics of predictions](https://ipfs.io/ipfs/QmP81EcuNDZHQutvdcDjbQEqiTYUzU315aYaTyrVj6gtJb)
- [DURA](https://github.com/cybercongress/cyb/blob/dev/docs/dura.md)
- [Nebulas](https://ipfs.io/ipfs/QmWTZjDZNbBqcJ5b6VhWGXBQ5EQavKKDteHsdoYqB5CBjh)
- [Colony](https://ipfs.io/ipfs/QmZo7eY5UdJYotf3Z9GNVBGLjkCnE1j2fMdW2PgGCmvGPj)
- [Truebit](https://ipfs.io/ipfs/QmTrxXp2xhB2zWGxhNoLgsztevqKLwpy5HwKjLjzFa7rnD)
- [SpringRank presentation](https://ipfs.io/ipfs/QmNvxWTXQaAqjEouZQXTV4wDB5ryW4PGcaxe2Lukv1BxuM)
- [PageRank](http://ilpubs.stanford.edu:8090/422/1/1999-66.pdf)
- [RFC-6962](https://tools.ietf.org/html/rfc6962#section-2.1)
- [IBC protocol](https://ipfs.io/ipfs/QmdCeixQUHBjGnKfwbB1dxf4X8xnadL8xWmmEnQah5n7x2)
- [Tendermint](https://ipfs.io/ipfs/QmaMtD7xDgghqgjN62zWZ5TBGFiEjGQtuZBjJ9sMh816KJ)
- [Comparison of web3 browsers](https://github.com/cybercongress/cyb/blob/master/docs/comparison.md)
- [Cyb](https://github.com/cybercongress/cyb/blob/master/docs/cyb.md)
- [Cyb virus](https://github.com/cybercongress/cyb-virus)
- [SpringRank](https://arxiv.org/pdf/1709.09002.pdf)
- [How to become validator in cyber protocol](/docs/how_to_become_validator.md)
- [Top 10000 english words](https://github.com/first20hours/google-10000-english)
- [Multilingual neural machine translation](https://ipfs.io/ipfs/QmQUWBhDMfPKgFt3NfbxM1VU22oU8CRepUzGPBDtopwap1)
- [Foam](https://ipfs.io/ipfs/QmZYKGuLHf2h1mZrhiP2FzYsjj3tWt2LYduMCRbpgi5pKG)
- [Coda](https://ipfs.io/ipfs/Qmdje3AmtsfjX9edWAxo3LFhV9CTAXoUvwGR7wHJXnc2Gk)
- [Mimblewimble](https://ipfs.io/ipfs/Qmd99xmraYip9cVv8gRMy6Y97Bkij8qUYArGDME7CzFasg)
- [Tezos](https://ipfs.io/ipfs/QmdSQ1AGTizWjSRaVLJ8Bw9j1xi6CGLptNUcUodBwCkKNS)
- [Software 2.0](https://medium.com/@karpathy/software-2-0-a64152b37c35)
