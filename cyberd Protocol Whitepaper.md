# Motivated Search and Evaluation Engine for Permanent Web

[Dima Starodubcev](https://steemit.com/@hipster)

Draft v 0.1

January-February 2017, Bali

[Simplified](https://docs.google.com/presentation/d/1nsPwHLBu5UArQ4zbN1yaGUj8LBgEpCEFqVe6iKP5XLg/pub?start=false&loop=false&delayms=3000&slide=id.g20ef74535d_0_33)

## Abstract

Existing general purpose search engines are restrictive centralized databases everybody forced to trust. These search engines were designed primarily for client-server architecture based on DNS, HTTP, and IP protocols. The emergence of content-addressable storage and distributed ledger technology creates an opportunity for the Internet to acquire new features such as more efficient computing, storing and broadband consumption, more resilient, secure and private access, no middleman for digital property. This can shift existing web's ubiquitously used client-server architecture to truly peer-to-peer interactions based on stateless IPFS and variety of stateful consensus computers such as Ethereum. This creates a challenge and opportunity for a search engine based on emerging technologies and specifically designed for them. Surprisingly the blockchain architecture itself allows organizing general purpose search engine in a way inaccessible for previous architectures. In this paper, we discuss opportunities behind blockchain based search engine, challenges of crawling, indexing and evaluation for the next generation web and propose a blockchain based experimental set of smart contracts to address discussed issues.

## Introduction

Let us start a discussion from disadvantages of conventional general purpose search engines:

- No Transparency. Nobody outside of Google understands how the ranking really works. This creates a market for black and white SEO. The truth is that if e.g. Google discloses complete details of the ranking algorithm it would be easy for adversaries to game organic search results that kill the quality of results and ad revenue streams. Pagerank [PR] has no inherent trust mechanism resistant to Sybil attacks. This problem can be addressed adding transparent and accountable blockchain based ledger with properly designed economic incentives built into the system.
- No Access. Currently, all search engines are centralized. Nobody is able to add to index as well as participate in improving the quality of search results. However, Google itself internally uses a workforce of _search evaluators_. It is our belief that user-generated search engine could have a higher quality of results as in a story with almost every website in existence.
- Broken Incentives. The vast majority of contribution to a search quality is made by users. Then any user searches something she extends semantic core. Then any user clicks on search results she trains a model. This creates an opportunity to continuously improve ranking model at the expense of users. Then search engines sell users to advertisers at the expense of harming user experience and acquire revenue streams which are not returned back to users at all. This simple loop created Alphabet's $550 billion capitalizations (~$80 per Earth capita) in 18 years. We want to change that.
- Central Control. Google become too powerful. It is scary to imagine a future where _everything_ about _everybody_ is known and controlled by **closed** AI corporation. Imagine the world where (1) the only country exist, (2) nobody can control its government and (3) everybody should obey the decision of government without any explanation. There should be open, transparent and accessible _alternative_ with _decentralized control_ built on principles of modern distributed interplanetary content-addressable cyberspace [IPFS] and DAO like governance [RALF].
- Annoying Ads. Separation for organic and ad search results is unnecessary. In fact, all organic ranking decisions are being made by search authority. But for paid search Google use free market solution to determine a fair ad price for every word in its gigantic semantic core. Historically free market solutions are proven to be more efficient in virtually any area of decision making. Why do not use the same principle for the ranking itself disintermediating annoying ads? Let us imagine that every link can be (1) curated or audited by everybody, (2) based on this trusted metric cyber•rank (page rank based on economically incentivized curation and auditing) is calculated and then (3) everybody can promote this link further by burning some money automatically bringing value for everybody in existence. For every action, everybody earns a share proportionally to contributions. This non-zero-sum game is significantly more Sybil-resistant and that is there we are heading.
- One-Way Trust. Everybody use to trust Google, Baidu, and Yandex. But Google, Baidu, and Yandex don't trust users. E.g. you cannot report some kind of proof that given link is a lie and should not be indexed so high. It can count your attention during ranking but can reject to count it. You cannot know what happens inside because Google, Baidu, and Yandex don't trust us. We want to establish a system there trust is bidirectional between the search engine and users because search engine ownership is distributed across all its users based on which all ranking decisions are made.
- Zero Privacy. All search engines will answer you only if they explicitly know how to map your device with your real identity or pseudo-identity which is tracked by RTB [RTB]. Otherwise, you should prove that you are not a robot every time you search. That harm our privacy. Moreover, robot abuse is another hot topic that is about to happen. Nonetheless, nothing should harm our privacy.
- Censorship. Though it's well known that Google working hard to prevent censorship we all know about China [GCHINA](https://en.wikipedia.org/wiki/Google_China) case and [Transparency Report](https://www.google.com/transparencyreport/). A good search should be resistant to censorship without exceptions and build for interplanetary scale in mind.
- Online only. Worth to note that you cannot search offline even if necessary information is stored next door. If we are cut from the wire or backbone we powerless. Global offline search is not a feature which can be easily deployed even by a multibillion corporation. This goal is nearly impossible to achieve based on centralized architecture. Only accessible distributed systems can solve this fundamental problem for the next generation Internet. This future is not about gateway keepers in form of ISPs but about mesh networking and peer-to-peer communications.
- Weak Security. What happens if tomorrow my Google account will be blocked? Do we have something to prevent this? Do we have the necessary level of assurance that _guarantee_ us our security based on math and not on the complicated legal tender? All technical solutions are here but to solve this important issue we need to do a lot of work because security is a foundation for life, liberty, and property.

A pretty huge amount of problems to fix. It would be naive to bootstrap a search engine from a scope where Google, Baidu, and Yandex exist. We need to find a special area there general purpose search engines sucks. A variety of distributed ledgers such as blockchains and tangles can be primary content-addressable data suppliers and this is a scope where current search engines are not the best at work. Moreover, blockchain technology evolves very rapidly and has a lot of promises so it is a sure bet.

The idea is to initially deploy a blockchain based search engine for the purpose of searching against other blockchains so that can be useful from the first day. At the same time, we design the cyber•Fund _application_ [CFUND] based on cyber•Chain to solve a problem of trustless realtime blockchain asset valuation which adds some useful capabilities to the search. But we need to design cyber•Chain in a way to be scalable for a more broad definition of a general purpose _search and valuation engine_, so more applications can emerge. Currently, about 15 trusted (and thousands of not so trusted) and globally available distributed ledgers exist, totaling about 1 billion transactions accumulated. Last year (2016) amount of accumulated blockchain transactions increased tenfold. Not all of them were financial transactions in some sense. E.g. Steem [STM] blockchain and it's Russian sister Golos [GLS] primarily store user generated text such as posts and votes (? votes are text? confusing...). As transactions are the only way of changing states in databases currently used by any web application we foresee distributed ledgers to become the primary source of public information in the 21st century due to tremendous benefits of the technology [ENIGMA, ....].

Thus we are to declare the _principles_ of a general purpose decentralized and distributed search engine for the upcoming age:

- Privacy and Security. Just it.
- Ubiquitous Ownership and Access. Everybody should have a right to possess a piece of it.
- Mesh networks future proof. It should work in every connected surrounding.
- Interplanetary scale. It should work on Earth and Mars.
- Tolerant. In the era of machine learning, it should work for any kind of thinking beasts.
- Open and Accessible. Everybody should be able to bring a bit to a quality of search results.
- Blockchain Agnostic. Foundations behind its design should not rely on any protocol or stack rather be explicitly derived from the nature of the information itself.
- Beautiful. The business model should not harm the user experience.
- Transparency and Trustfulness. Every piece of its reasoning and behavior must be auditable by everybody.
- No Single Point of Failure. Nobody should have a single key to modify or change it.
- Sybil Attacks Resistant. This resistance should be derived from the properties of a free market but not from some single authority.
- Intelligent. It should answer _natural_ questions with _easy to read_ and _provable_ answers no matter text, media or natural numbers should be involved in the answer.

Alongside these principles, there are _requirements_ such as performance, usability, and scalability. Keep reading. We discuss everything step by step.

## Design Rationale

The idea of a permanent web behind IPFS is beautiful in its simplicity. Every unique _piece of data_ has the unique address:

>

Using this address this piece of data can be found in a global data structure called MerkleDAG [DAG] (logical representation of data storage) across peer-to-peer IPFS network using bulletproof DHT. Nodes are weakly incentivized for fair data exchange using BitSwap protocol. This link can point to _any_ piece of data such as GIT object, BitTorent link, Bitcoin block, JSON document, picture, video, plain binary data or even a small piece of text:

>

Let me explain the power of this solution for a search engine:

_Data is unique and self-authenticated_. If you know this hash is a piece of data you trust, you should not care where it comes from. This property free search engine from the storage of documents.

_Direct Content Distribution_. Weak incentivization of the BitSwap protocol has a very interesting side effect: the more popular the file => the more people store pieces of it => the faster you can get it. This is in contrast with a conventional client-server architecture where the more people want the file => the more resources a server needs (and more expensive distribution become) => the slower you can get it. This property significantly reduces resource usage to _deliver_ results right to a user through distributed network.

_Flexible Data Discovery_. Big pieces of data are permanently linkable, thanks to MerkleDAG, so you can trivially reach a small chunk of data. That enable a lot of powerful applications impossible for a peer-to-peer network before. E.g. SQL-like data queries to a database distributed across the network or REST-like queries.

Thus documents which are located at `/ipfs` are _immutable_. But what if we want _mutability_? IPFS offer the self-signed naming system which doesn't rely on conventional centralized and slow DNS (URLs we use to use). Everybody can publish a link in a namespace `/ipns` with a mutable pointer to any `/ipfs` piece of data and sign it with its node's private key:

>

After an owner is able to mutate this link pointer to another /ipfs piece of data. This gives us another important property for advanced search engine:

_Fast Discovery_. This property free search engine from rescanning work significantly reducing resources necessary for keeping index fresh. Workflow of conventional search index is straightforward: web crawler fetches a page, then follow up links, fetch the following pages and so on. After some loops, web crawler is able to get virtually all links in existence. The information contained in web pages is mutable. Thus search engine should decide then and how to fetch indexed links for new versions of documents. As web pages are loosely structured than it should somehow filter all the noise. In our case publishing is visible in DHT and can be nearly instantly visible for all participants across distributed network without the necessity to continuously send, receive and process HTTP requests. An average request is about 800 bytes long, so do every response [HTTP]. This overhead goes through the wire every time search engine want to know the freshness of the web page and then if the page is new it should crawl all new page even if 1 bit has been changed. Interestingly, though HTTP/2 bring new awesome features it still relies on the plain old location-based addressing and ancient DNS, thus this cannot be faster by design. In our case, we need sniff changes in DHT to know all news from a network. Currently, this is about 100kb/minute in a quite small IPFS network. There is estimation [QGTC] that about 10% of Internet traffic consumption has been made by Google crawlers. It easy to imagine how much overhead can be eliminated across the entire planet. This overhead currently is being paid by businesses around the globe.

_Flexibility_. Search results can be not ugly static snippets from 20 century but small dynamic programs _owned_ by creators. It is hard to overestimate this idea. We will cover this topic further.

_No Middleman_. A search engine should not rely on a DNS middleman and can communicate directly with resource creators.

Ok. Now we understand that we have a usable way to store, reach and mutate the data across the globe in a more lean way. But what about an ability to _comprehend_ what data is behind this meaningless hashes? Interplanetary File System uses a novel multihash, multicodec, multiformat, multibase [<https://github.com/ipld/cid#cidv1>] formats for Interplanetary Linked Data (it is not a joke either). IPLD is a project aims to define CID that being used across IFPS network. Content IDentifier or CID is a self-describing content-addressed identifier. That creates enormous opportunity for optimization. Imagine a 32-byte link can contain everything to _independently_ understand how to reach a piece of data and how to programmatically interpret it. Finally, CID is the thing which makes simplistic design of our proposed search engine possible inside consensus computer.

A nice property is that CID are based on well-known cryptography so we don't need to rely on IPFS if something better came. The same addresses can be used to exchange a piece of data in a different peer-to-peer network. Though this doesn't solve a problem of extensibility. CID itself contain version bits, thus if something better came we can switch to it.

But now IPFS and IPLD is a perfect design choice with huge momentum across academia, open source engineering community and (the most important) blockchain engineering teams. It is our belief that it should lay as foundation for the next generation search engine. Everything is perfectly likable (MerkleDAG), fast (DHT), accessible (BitSwap) and comprehendible (IPLD) without a single point of failure. IPFS is a protocol and it can be more useful with a good search engine specifically designed to it. It's possible to compute a PageRank for the whole MerkleDAG. But there are two problems with it:

- An amount of links in MerkleDAG grows `O(n^2)`. That is not either conventional web pages with 20-100 links per page. For 1 Mb file can be thousands of links. Once the network starts to take off, complexity inevitably increases.
- Even if we address some algorithm to extract relevances from these links we should address even more algorithms to extract a meaning.

What we need is to find a way to incentivize extraction from this data fog a meaning that is _relevant to users queries_.

## Search Workflow

Our proposed design of a search engine is based on advanced blockchain technology and enables everybody to participate and be rewarded. Everybody with an account can search. To execute queries user should sign a `search` transaction with CIDv1 as payload and broadcast it.

> cyberd search {CIDv1} hipster true
[where s the privacy declared before? with every search query signed with account name on a public blockchain? (tomarcafe)]

A document should be a valid CIDv1. Post headers are purposefully unique. The rationale is the following:

- strongly encourage to bring valuable documents first and get a reward
- simplify search execution and ranking calculation
- make possible a high-performant flat data structure

Worth to note that search request "buy tesla x" is also can be represented as CIDv1. Thus search index store queries itself as documents. Hence search queries and target documents both will acquire cyber•rank.

Then cyber•chain verifies validity of CIDv1, the correctness of signature, broadband allowance and if all conditions are met either of two things happens:

- if CIDv1 is not in search index => Write to search index
- else broadcast a vote for existing CIDv1 and return sorted links of relevant CIDv1

Based on this data client-side application can deliver documents by 3 ways:

- using full javascript implementation of IPFS [JSIPFS] (fast, but requires initialization)
- using REST API provided by IPFS HTTP gateway (depends on, can be fast or not)
- using local IPFS node (the fastest)

This approach is simple and powerful. A developer has enough choices to balance between usability and performance.

CIDv1 can mean any piece of data. Different data types that can be returned to a user depending on the query:

1. Plain text for autocomplete, e.g. ['apple', 'asos', 'amazon']
2. A piece of media content which a user can play without necessity to go somewhere
3. Static formatted snippet with a link to the conventional web.
4. Static formatted snippet with a text answer.
5. IPNS link pointing to javascript that can return dynamic snippet.

It depends on developers (mainly submitting answers) and users (mainly ranking answers) what kind of things they want to answer questions. Possibilities are limited with imagination. Thus we propose a free market for answers on search queries everybody is encouraged to participate. How does it work?

Information about indexed CIDv1 as well as about its rankings is available for everybody. Thus those interested in rewards can monitor the blockchain for semantic core updates and submit links nearly instantly. Everybody can sign `answer` transaction with a link from

<search cidv1=""> to <answer cidv1=""> as payload and broadcast it:</answer></search>

> answer {search CIDv1} {answer CIDv1} hipster true

A document should be a valid CIDv1 and unique outbound link from an answer. So for any given question, the only unique answer is possible. Then cyber•Chain verify the correctness of signature and either of two things happens:

- if CIDv1 is not in search index => Write to search index
- if answer CIDv1 has no link to question CIDv1 => Write to answer index
- else broadcast a vote for existing answer CIDv1

> We follow black box rule. In order to answer a question right, you need a full comprehension neither the question nor the answer. You just need to match a query with the most relevant links.

In order to increase the rank, everybody can promote either link or query.

> promote_search {CIDv1} hipster 20 true
> promote_answer {search CIDv1} {answer CIDv1} hipster 20 true

Every sent token for promoting is destroyed thus creating value for every token holder.

That is a core API for the entire blockchain. Other methods accomplish support role for the thing. Such compact design opens huge opportunity for performance optimizations. Also, a clean and comprehensible experience is very important for those who want to be involved. That is. The entire graph of the semantic core with weights is open for everybody and available for data mining or any kind of weird AI stuff. But to make it work we need to find a way to calculate relevance.

We can represent our data structure as directed acyclic graph where vertices are indexed documents and edges are directed links between them.

![](https://docs.google.com/drawings/d/1--Uj85OiU-uwj0gxUFWZDbjPuby3IEMBBVFSmHTNkDc/pub?w=785&h=436)

We equate terms `document` / `query` and `link` / `answer` as for our use case these are practically the same. We will stick to `query` and `answer` terms in order to avoid confusion.

Hence if a user searches a document `CID 3` (query) search engine will return links (answers) to `CID 1`, `CID 4`, `CID 5` documents (queries) sorted by cyber•rank. Let us discuss it in details.

## cyber•rank

The idea is to combine two simple yet powerful algorithms: Google's PageRank and Steem's reward mechanism: ![cyber•rank](https://docs.google.com/drawings/d/1yvkyeOpVZoiyUOv0Gqu-kywsaWF16A9NTw8GgJq-rX8/pub?w=1521&h=777)

Where `t_rank` = `n_rank` + `s_rank`

`n_rank`, or natural rank is a rank based on Steem reward system. `s_rank`, or synthetic is a plain old PageRank used by Google and others.

Natural rank is acquired in a process of auditing and curation. Based on this rank payout to those who involved (queries and answers submitting, auditing and curation) a made. Each piece of submitted data gets paid in 7 days. Each piece of data can be voted by cyber•power token holders. We use `auditing` term primary for verification of submitted data by automated scripts and `curation` term to denote manual curation. From a technical standpoint, those are primary the same as both utilize the same method for interactions. Implementation of natural rank is almost identical to Steem. Thus details of implementation can be found in our Github [] or Steem whitepaper.

Synthetic rank is taken as initial natural rank expressed in `rshares`. This makes possible to start calculating PageRank before payouts have started. Conventional search engine work with nearly zero trust data. More than 200 factors used to calculate initial PageRank for a new document in the graph and find relevance to the search terms. Our novel approach allows assigning initial value based on Sybil-resistant voting. It is our belief that proposed approach can significantly simplify ranking and be more precise for information that nature is subjective. Though PageRank calculation is a trivial task we should estimate the feasibility of doing so in a consensus computer such as Steem where scalability is limited with a less performant node and parallel processing [Steem Roadmap] or sharding is yet to be discovered [Ethereum Mauve].

A recent study [http://www.vldb.org/pvldb/vol8/p1804-ching.pdf] shows that Facebook scale PageRank computation is doable (using Java based Giraph) for 1 trillion edges and 1.4B vertices in 600 minutes per iteration. Hence consensus computer made of commodity hardware will be able to process 1 iteration per 2 days for a 10B unique document (SWAG for all blockchains + Git + BitTorrent (https://arxiv.org/pdf/1009.3681.pdf) + IPFS). Our implementation is based on C++ thus can be more performant though is not guaranteed. Also, we have an opportunity to use the most performant GPUs available operated by witnesses and not that is being used by cloud providers [http://www.hipc.org/hipc2011/studsym-papers/1569513051.pdf]. Anyway, our estimation proves that it is practically enough to use CPU for proof-of-concept stage. GPU implementation can multiply computation capacity up to 1000x. But further research in the field of parallel consensus computing is necessary to achieve Google scale (10000x more documents) realtime decentralized computation of cyber•rank.

Our model is recursive and requires the enormous amount of calculations which are limited within blockchain design. Model recalculation does not happen on a periodic basis rather it continuous. We consider introducing consensus variable, in addition to a block size, in order to target processing capacity of the network. Let's call it a _computing target of documents per block_ or CTD. Any witness will be able to set a number of documents the network should recompute every block. The blockchain takes as input computing target of legitimate witnesses and computes CTD as daily moving average. Based on CTD blockchain can schedule the range of CIDs that should be recomputed by every witness per round.

The semantic core is open. Click-through information is stored on-chain. Every time a user follow a link positive voting transaction is broadcasted e.g. with grade 1\. Voting on a protocol level is a number in a range from -100 to 100\. Thus application developers have a tool to implement different grades for different kind of interactions. Such design is crucial to train the model and acquire a data about search popularity of semantic core and its volume. Currently, search engines are very careful in revealing this information because this information is the most important part of the ranking. We want to change that. Every time a user click on a snippet developer earn a fair portion of emission and on chain model is trained. Application acquires the more rank the more rank acquired by its links. The more cyber•rank acquired - the more revenue streams for an application developer.

Both algorithms have strong proof in form of Google's $550 B capitalizations in 18 years and Steem $40 M capitalization in 9 months. Combining both it is possible to empower the world with a new kind of search quality that has been (1) designed to index relevant document fast and (2) has inherent Sybil protection.

## Self Indexing Dilemma

Proposed approach has very unexpected limitation. What if we want to index cyber•chain using cyber•chain itself? Let us say that we have an awesome transaction that happens inside the chain and everybody are talking about it. It is popular thus we should display it in our search results. Adding it to an index spawn another transaction which (surprise) also should be indexed. This entanglement creates an infinite loop that bloat cyber•chain. This can not be a problem either. Consensus computer capacity and power are limited by the market forces. So we have two possible decisions:

1. Let it be. The market is a king. A bit of bloat can be a good piece of a knowledge about itself.
2. Strictly forbid indexing of the blockchain itself. Fortunately, it is not so hard to implement on a consensus level. All we need is to check that CID has not been included in a cyber•chain before. That mean that cyber•chain transaction itself remain unindexed because in order to achieve this we (1) either should mutate data for hashed and timestamped transactions, (2) nor create a possibility for a self-bloat. None of the options is valid. Again, fortunately, direct search without ranking among internal cyber•chain transactions can be available inside search results with the help of CID magic.

It is our belief it is better to stick to restrictive policy without further research.

## Challenges and Advantages of Indexing Distributed Ledgers

Conventional general purpose search engines were built on the _last mile assumption_ of stateless HTTP(S) protocol. Indeed the last mile approach worked well on the Internet where information emerges inside stateful private databases and only after become publicly available using HTTP. The emergence of content-addressable systems such as Git and BitTorrent didn't change much as those can not be compared with stateful private databases in any sense even though these protocols are represented origin of a content. But with an emergence of distributed ledger technology become possible to get know about public content in the moment than it actually has been born. In this sense, blockchains and tangles can be viewed as a real alternative to conventional private databases. That breakthrough create enormous opportunity for better and faster indexing but at the same time has inherent problems solving of which create advantages:

_Chain Validation_. Protocol diversity is the hardest part of solving an issue of chain validation using another chain because solving this requires full implementation of one consensus computer inside another consensus computer. There are two efforts exist that try to solve issues of validating one chain using another. Among them are [Polkadot](https://github.com/polkadot-io/polkadotpaper/blob/master/PolkaDotPaper.pdf) and [Cosmos](https://cosmos.network/whitepaper). Both projects aim to solve a problem of trestles inter blockchain communications. Polkadot aims to design a system which doesn't require a trust to parachains (or interconnected chains). The design of Polkadot is very complex and has inherent scalability limitation. The design of Cosmos is much more simple, but require a trust to zones (interconnected chains). Our design doesn't try to solve a problem of inter blockchain communications rather try to implement probabilistic search across blockchains. Thus it is significantly more simple. Instead of recording information that has been verified from point of view of exogenous protocol we let this information come to index letting market forces and relevance algorithm determine which chains are correct and which are not.

_Meaning Extraction_. Practically there are no common fields that are used by all blockchains and tangles. The only common pattern is `tx_id`. There is no way to extract any meaningful information from `tx_id`. That is practically mean that there is no easy solution exist. For every protocol can be a large number of different approaches to extract useful information. Using advanced machine learning is infeasible for consensus computers at the moment. We offer set of smart contracts that create a free market for extraction of the meaning.

The huge advantage of blockchains that they are implicitly told what happens. Blockchains are highly optimized databases. Every transaction cost money so they are neutrally protected from spam and contain only data that really matter. That reduce noise filtering to nearly zero level. Structured raw data about blocks and transaction is available for everybody. This fact also reduces consumption of resources and make extraction of better meaning possible.

![Google vs Blockchain](https://docs.google.com/drawings/d/1AkdiCTalgqKqkOKgeJCtD4CCwt68buP1rDN2J95FGSY/pub?w=1440&h=743)

Blockchains provide real-time high-quality structured data which don't require *crawling* in a traditional sense. One node of any given blockchain is enough to provide verified data about transactions within one ledger without the necessity to continuously revisit resources significantly reducing costs.

All these factors create a free market opportunity for emergence the diverse set of highly specialized (on a very limited set of a semantic core) but highly efficient *broadcasters*.

_Probabilistic Settlement_. Blockchain designs, especially Proof-of-Work based, implies that finality of a transaction is probabilistic. A moment than a given fact can be considered as truth is blurred. [Previous researches](https://medium.com/@lopp/the-challenges-of-block-chain-indexing-30527cf4bfbd) show that it is expensive to achieve real-time blockchain indexing due to reorganization issues. The intention of discussed in the article design is to answer the question deterministically. Our architecture based on the principle that indexing a system with probabilistic settlement require probabilistic answering. Instead of deciding whether or not this particular block has been included in the canonical chain we can index all of them calculating the probability of finalization using cyber•Rank.

So solve of discussed issued of probabilistic indexing we propose a flexible approach we call lazy oracles.

## Lazy Oracles

One specific ability is crucial for the next generation search engine. Application developers should have a motivation to provide structured arbitrary data feeds. Thus search engine can answer natural questions aggregating data from _highly structured and defined_ feeds. This makes possible a user get high-quality _calculated_ answers in real-time about the state of reality expressed not only in links (which intelligent agent don't know how to parse) but in actionable numbers based on which it's possible to make independent economic decisions. It is hard to find a tool to agree on publicly available and continuously evolving facts. We propose an approach to solving this.

Today different blockchain have the functionality necessary to implement this. For eg. Ethereum enable construction of smart contracts that can validate and incentivize data feeds. But Ethereum has the strong limitation: a price. Due to a network design, every operation should be validated across the network of 5k nodes. For every put operation developer of such contract should pay in hope that somebody in the future will use this feed in the future returning costs. The current cost of a permanent storage inside Ethereum contracts is around $200/megabyte. Worth to note that Ethereum has consensus variable gas limit. Currently, a network load is around 10% of established limit. Once the demand for computation reaches a limit we will have a situation very similar to Bitcoin block size debate and price for storage can reach $2000k/megabyte easily without validation costs. Pretty expensive for unlimited possibilities. There is the alternative - Factom [FACT]. Its consensus design relies on a small amount of paid servers. Thus the cost is around $1/megabyte. Such low price comes with high limitations. You can only put data to Factom and read it. There is no validation and incentivization built-in. There are permissioned blockchain designs such as BigChainDB [BCDB] and Hyperledger [HYPL] which solves validation problem perfectly but require strong efforts for developers to program and establish a network and then somehow monetize it. Lazy Oracles are going to fill this gap providing robust, cheap and reliable way _for monetizing_ structured public data. Any cyber•chain account has a share in a broadband depending on its cyber•Power thus granting lifetime assurance of network usage. A network programmed with decaying inflation a part of which goes to all who participate in indexing depending on the valuation of subjective contributions.

The process consist of 5 steps:

Step 1: Everybody can declare a soft protocol for data feed by posting a CID with the inbound link to _oracle-defenition_ CID pointing to the document with the following structure:

```
// Basic Validation
doc_type: market_update // should be fixed
  exchange: string // domain name of the exchange
  base: string // link to a definition of a referenced namespace e.g. ISO or chaingear
  quote: string // link to a definition of a referenced namespace e.g. ISO or chaingear
  price: number
  volume: number
// Audit Rules
Document format and data types should match a protocol
For every unique exchange and pair `market_update` report can be submitted no more frequently than once per minute.
If either a price or volume for unique exchange and pair has not been changed more than 0.1% a report should not be submitted.
// Inbound Links
<oracle> <market-update> <exchange> <base> <quote>
// Reference Crawler Implementation
CID
// Reference Auditing Implementation
CID
```

Protocol declaration should be agnostic from language implementation and unambiguous.

Step 2: Now every reporter can submit data according to a soft protocol. If data begin to be submitted without protocol definition a value of this data can be significantly lower. Thus auditors are encouraged to flag such documents. If false or malicious data come auditors have strong incentive to flag such data. All lazy oracles should be feed with self-descriptive links so other participants will able to faster process auditing.

Step 3: Auditors validate any given document by scripts.

As result reporters, auditors (where objective script-based validation is possible) and curators (where subjective human-made evaluation is more appropriate) have strong incentive to (1) bring important data, (2) curate and audit data conscientiously. Proposed approach doesn't have strong guaranty of the truthfulness of data such as Augur does. Rather it is correct to say that we can have strong assurance that data is correct. But it is cheaper, faster and significantly more flexible. The process is robust in its unstructured simplicity.

Step 4: Payouts for auditors are made. Payouts information is input information for cyber•rank. Using cyber•rank it is trivial to filter feeds using plugins to give significantly more precise data with a higher level of assurance.

Step 5: High-quality data feeds (or lazy oracles) are available for a consensus engine of a search engine (thus cyber•rank can be continuously improved) and everybody on the planet.

We call this type of oracles lazy because they don't require strict rules of validation at the expense of reducing the level of accuracy (but still enough to reason with a high level of assurance). Also, they are lazy because the don't require to think about monetization for participants rather consensus engine print rewards based on a valuation of subjective contributions. This approach is superior to Ledgys [] than reporters should sell encrypted data pieces using costly Ethereum storage.

The most obvious use cases for Lazy Oracles

- Block indexing
- Transaction indexing
- Balance calculation
- Identity crawling
- Token valuation
- Token rating
- Token description
- Crowdsale tracking
- Weather measuring
- Traffic measuring
- Business performance reporting

Worth to note that use cases are not limited to mentioned above. Data structure and validation rules are arbitrary. So we can think of it as general purpose tool for _popular structured public data_ auditing and curation. Now we can bootstrap and index with amounts of useful information. But what about meaningful search results?

## Dynamic Snippets

This can be thought as serverless micro javascript(is javascript futureproof? ;)) applications that can dynamically take input data from the following sources:

- a search query itself
- asynchronously from background search query
- from a browser APIs
- from device sensory information
- from information about a user stored on cyber•Chain
- from other blockchains
- from IPFS and IPNS
- from conventional HTTPS or WebSocket APIs.

Every application is CID written to search index and answer index. Before developing an application developer should target it for a specific semantic core before defining it. Semantic core and its statistic are publicity available in a blockchain. Developer shall need to develop an application and submit links to an application using either immutable IPFS documents or mutable IPNS pointers for a targeted semantic core. Thus it is up to a developer to define what search queries are a better fit for a particular application. Keep in mind that dynamic snippets are naturally competing for a higher position in search results. Dynamic snippets can be sorted by cyber•rank, and as result become trustful. Worth to note that developers can significantly reduce spendings on app infrastructure as dynamic snippets can be delivered through content-addressable distributed network. Mutable IPNS pointers allow developing snippets for a targeted semantic core and not for every unique query. An implication of this approach is hard to overestimate. E.g. dynamic snippets combined with blockchain wallet make possible to shop right from search results.

The only potential problem with proposed approach is the safety of third party javascript code. Sandboxed third party code is able to mitigate this risks. Web technologies such as [*web workers*](https://www.w3.org/TR/workers/) and [*web components*](https://developer.mozilla.org/en-US/docs/Web/Web_Components) are being actively developed. E.g. javascript library *jailed* [https://github.com/asvd/jailed] is able to do exactly what we need. Further adoption of web components is also one of a possible solution to do that safely.

## Dealing with Long Tail

It is well-known fact that every day Google receive up to 20% of new search queries never seen before. Thus we need to find a way to deal with it. Here are 3 popular use cases with simple solutions (surely non-exhaustive):
- _Misspellings_. These queries contribute fair half to ever-growing unique query set. But it is not a rocket science to stem such queries client side without any indexed knowledge. After a user clicks on client side suggestion correct link can be submitted to cyber•chain thus improving the global model.
- _Phrases_. E.g. `forest gump imdb`. These queries can be a combination of well-known terms and contribute another half to ever-growing unique query set. Even if we get relevant links from indexed `forest`, `gump` and `imdb` separately we would not able to combine answer to return meaningful movie rating to a user. But we can find _the closest documents_ between `imdb`, `forest` and `gump` and sort them by relevance. Those likely be the most relevant answers. This simple method can significantly increase efficiency for long tail queries without rocket science. As API is open for everybody there is no limits on using advanced technics.
- _Unique queries_. These are about 10% of never-seen-before queries. We expect that the market of linking will create a segment for on-demand answers. We remember that any query can be seen in memory pool by every network participant. Thus opportunity to earn can create a healthy and competitive market for an on-demand answer in an environment without technical limitations.

## Spam Protection

In the center of spam protection system is an assumption that write operations can be executed only by those who have vested interest in the success of the search engine. Every 1% of stake in search engine gives the ability to use 1% of possible network broadband. As nobody uses all possessed broadband we use fractional reserves while limiting broadband like ISPs do. Details of an approach can be found in a Steem white paper.

Auditing and curation are based on Steem reward mechanism. It is Sybil-resistant approach as votes are quadratic based on principle 1 token in system = 1 vote. In order to vote one should vest in shares for at least for 20 weeks. That solve a problem entirely because those who have a right to vote are strongly incentivized in a growth of his wealth. In order to prevent abuse of auditing and curation voting power decay implemented exactly as in Steem.

## Applications

It is hard to imagine what kind of applications can be built on top of proposed foundation. I'd like to mention some outstanding opportunities which can be build using cyber•Chain and IPFS:

- Relevance everywhere
- Blockchain browser
_ Multi-protocol wallets
- Offline search
- Smart command tools
- Autonomous robots
- Language convergence

_Relevance Everywhere_. Proposed approach enable social, geo, money or anything aware search inside any application. It is trivial to implement a search relevant to a particular identity using proposed algorithm. The more a user train a model the more behavioral data can be associated with her. This personalized information can be stored locally for (1) faster retrieval and (2) offline access.

_Blockchain browser_. It easy to imagine the emergence of a full-blown blockchain browser. Currently, there are several efforts for developing browsers around blockchains and distributed tech. Among them are Beaker, Mist and Brave .. All of them suffer from very limited functionality. Our developments can be useful for teams who are developing such tools.

_Multi-protocol wallets_. Currently, there are several efforts for developing easy to use the universal wallet for blockchain assets. Jaxx and Exodus are among them. Developers of such applications suffer from a diversity of protocols around blockchain tech. There is no fully functional multi-asset wallet yet. Our developments can help teams who are developing such tools.

_Actions in search_. Proposed design enable native support for blockchain asset related activity. It is possible to design applications which are (1) owned by creators, (2) appear right in search results and (3) allow a transact-able call to actions with (4) provable attribution of a conversion to search query. e-Commerce has never been so easy for everybody.

_Offline search_. IPFS make possible easy retrieval of documents from surroundings without the internet connection. cyber•Chain itself can be distributed using IPFS. That create a possibility for ubiquitous offline search.

_Smart Command Tools_. Command line tools can rely on relevant and structured answers from a search engine. That practically means that the following CLI tool is possible to implement

```
>  mylovelybot earn using hdd -100GB

searching for opportunities:
cyberd search <earn using hdd>

The following answers received:
- apt install siad /// 0.0001 btc per month per GB
- apt install storjd /// 0.00008 btc per month per GB
- apt install filecoind /// 0.00006 btc per month
...

Made a decision try `apt install siad`
Git clone ...
Building siad
Starting siad
Creating wallet using your standard seed
You address is ....
Placing bids ...
Waiting for incoming storage requests ...

```
Search from CLI tools will inevitably create a highly competitive market of a dedicated semantic core for bots.

_Autonomous robots_.
Blockchain technology enables creation of devices which are able to earn, store, spend and invest digital assets by themselves.

> If a robot can earn, store, spend and invest she can do everything you can do

What is needed is simple yet powerful API about the state of reality evaluated in transact-able assets. Our solution offers minimalistic but continuously self-improving API that provides necessary tools for programming economically rational robots.

_Language convergence_. A programmer should not care about what language do the user use. We don't need to have knowledge of what language user is searching in. Entire UTF-8 spectrum is at work. A semantic core is open so competition for answering can become distributed across different domain-specific areas, including semantic cores of different languages. The unified approach creates an opportunity for cyber•Bahasa. Since the Internet, we observe a process of rapid language convergence. We use more truly global words across the entire planet independently of our nationality, language and race, Name the Internet. The dream of truly global language is hard to deploy because it is hard to agree on what mean what. But we have tools to make that dream come true. It is not hard to predict that the shorter a word the more it's cyber•rank will be. Global publicly available list of symbols, words and phrases sorted by cyber•rank with corresponding links provided by cyber•chain can be the foundation for the emergence of truly global language everybody can accept. Recent scientific advances in machine translation [GNMT] are breathtaking but meaningless for those who wish to apply them without Google scale trained model. Proposed cyber•rank semantic core offers exactly this.

This is sure not the exhaustive list of possible applications but very exciting, though.

## Incentive Structure and Distribution Mechanism

To make cyber•rank economically resistant to Sybil attack and to incentivize all participant for rational behavior a system uses 3 types of tokens: CYBER (or cybers), CP (or cyber•power) and CD (cyber•dollar)

CYBER is a transferable equity token which is analog of STEEM. The intrinsic value of CYBER came from the ability to convert it to CP.

CP is a non-transferrable equity token which is analog of SP in Steem. CP can be converted to CYBER in 20 weeks using proportional weekly payments. The intrinsic value of CP came from the right to (1) write to an index according to a bandwidth limit, (2) rank objects, (3) promote objects (4) make consensus decisions. CP can be converted to CYBER in one year.

CD is a debt token with a relatively stable value which came from an ability to convert it into CYBER within 3 days by the price submitted by witnesses and calculated according to cyber•rating methodology [] (don't confuse). 1 CD tracks 1/10^12 of _provable_ blockchain economy.

Reward Pool is defined as 100% of emission and split among the following groups:

Infrastructure Reward Pool

- Witnesses - 5%
- Investors - 10%

Indexing Reward Pool - 30%

- Reporters - 10%
- Auditors - 20%

Linking Reward Pool - 65%

- Responders ~ 20%
- Trainers ~ 40%

Our implementation also offers an incentive for CD holders. They receive APR on holding according to rate defined by witnesses.

As our primary focus is developers of decentralized and distributed applications we offer a convenient way of revenue generation in addition to possible revenue on dynamic snippets development. Every search and answer transaction can define a smart contract with up to 6 beneficiaries with a distribution of potential reward from this particular action. This creates a solid foundation for bootstrapping application ecosystem across a conventional web, mobile app stores and upcoming VR platforms.

A conventional advertiser is also has a tool to participate. Any link can be promoted thus increasing cyber•rank. Revenue from the promotion is burning thus decreasing the supply of CYBER and bringing value to all shareholders of a system. We do believe that this approach is fair and don't harm user experience but open opportunity to compete with conventional search ad a bit.

Virtual loop of the business models for our decentralized autonomous organization is pretty simple: `More indexing => More people search => More developers build => More people earn, rank and promote => Better infrastructure => More indexing`.

A network starts from 100 000 000 tokens

- Crowdsale - 70%
- cyber•Fund reward - 10%
- Founders reward - 10%
- Sharedrop & License - 5%
- Development reserve - 5%

Since inception, a network prints 3 CYBER every block. Every 1 million blocks it reduces print rate on 1%. Thus starting from ~35% print rate per year inflation begin to reduce gradually until it reaches 1%.

There is one problem with proposed incentive structure. We call it _language incentivization bias_. In the core of cyber•chain is quadratic voting. The system needs it to effectively incentivize participants for quality ranking. But that natively leads to weak incentives across different language groups. E.g blockchain Golos was deployed as Russian alternative to Steem because Russian posts acquired only 0.2% of rewards though providing 10% of the content. The idea of deploying cyber•Chain is great if it can be truly global from the start. The only way to overcome this bias is a global deployment from day 0, because otherwise we need significantly increase the complexity of the reward system. We offer good incentives for translation of this white paper to 50 languages worldwide as well as call for action to all blockchain communities across the globe.

Crowdsale will start after 30 days of network launch and will go 30 days. Crowdsale will be capped with 10000 bitcoins without limit in the first day - rush day. During crowd sale, the following discounts will be applied.

- 30% in the first rush day
- 20% from 2 until 15 day
- 15% from 16 until 18 day
- 10% from 19 until 21 day
- 5% from 22 until 24 day
- 0% from 25 until the end

All participant will receive CP proportionally to invested bitcoins after discounts.

## Extensibility and Governance

Currently, our implementation has the following functionality available for application developers.

- Custom Operations. Better alternative for OP-RETURN
- Plugins. Allow implementing API based on custom operations.
- Escrow. The core smart-contract enables 3d party arbitrage for arbitrary transactions.
- Private Messaging. Enable private communications between accounts.
- Dynamic Account Permissions. You can think about it as better multi-sig.

The following possibilities can be available in a distant future:

- Sidechains
- State channels
- Permissionless smart contracts

Consensus can be changed in case of 17 of 21 elected delegates accept a hard fork.

## Search and Evaluation Appliance

For bootstrapping a network we are going to offer software configuration (Cybernode - Github) on top of well tested open source specs for hardware configuration of commodity computer (Enterprise - Github) which cost around $10k-$30k depending on RAM and SSD capacity and is able to participate and earn by itself executing different network tasks:
- operate as witness node
- operate as indexer and auditor
- operate as answering node
- operate as fast and cheap backend for decentralized and distributed application

We need a network of high performant computers in order to achieve our goals. Necessity comes from  the following assumptions:
- all blockchain nodes and IPFS should live inside one machine to remove slow network communications from all middleware. Vast information for processing can be in memory.
Communications inside one bus enable to execute required tasks significantly faster [[https://gist.github.com/jboner/2841832]]. Thus we can achieve nearly live indexing of reality from the very beginning.
- extension with GPU. Currently, data centers cannot compete with commodity GPU. E.g. Amazon offers very expensing professional Nvidia Tesla cards. For our purposes, commodity cards such as GTX 1080 are much more cost effective.

*Enterprise*. Currently, it is not hard to assemble a 2 CPU computer with 1 TB of RAM and 40 TB of SSD using commodity hardware. Such appliance can cost about $30k so we can think of it as affordable for those who are seriously want to be involved in a project. Also, we have an option to extend the capability of proposed search appliance based on 2 CPU motherboards built on Intel C612 chipset. Usually, it has 7 PCI-E slots for GPU which can be dedicated for cyber•rank calculation. Thus a price for an ultimate (2 CPU * Xeon E5 * 22 cores + 7 GPU * GTX 1080 * 2560 Cuda cores) search and evaluation appliance can be around $50k. Currently, such computer will be able to process, index, audit and linking all blockchains.

*Cybernode*. We implementing the following software configuration that is based on docker containers.

![](https://habrastorage.org/files/31c/e9c/05c/31ce9c05c71d44ddbda2e7abc3ac10d7.png)

Cybernode allows everybody fast deployment of decentralized and distributed application powered with cyber•chain search capabilities.

## Performance and Scalability

Proposed blockchain design is based on DPOS consensus algorithm and has fast and predictable 3 seconds block confirmation time and 1 minute finality time. Average confirmation timeframe is 1,5 seconds thus conformations can be asynchronous and nearly invisible for users. A good thing is that users don't need confirmations at all before getting search response as there is no risk associated with that.

Current node implementation theoretically [https://bitshares.org/blog/2015/06/08/measuring-performance/] can process about 100k transactions per second. This theoretical bound is primarily limited with the possibility to replay a blockchain [https://steemit.com/blockchain/@dantheman/how-to-process-100m-transfers-second-on-a-single-blockchain]. As of now, all blockchains are about 1B immutable documents which size is about 200 GB with average tx 200 kb. We need to store all hashes which are on average 64 bytes long. We estimated that storing in the index all blockchain documents as CIDs and votes are roughly the same as storing all raw blockchain data. Linking 1B documents create significant overhead as blockchain index size can be up to 100 times more. Given this, we can assume that indexing all existing blockchains require about 4TB of SSD space. This is affordable for commodity hardware with 10x scaling capability without a necessity for sharding across several machines. We assume this is enough scalability margin for proof-of-concept.

Initial indexing of 1B documents and 100B links will require a continuous load of the network at the upper bound of its capacity in the first year of its existence. If we assume that network will be able to process 10k transactions per second with 2MB block size we will be able to index all blockchains in 4 months. Further operations will require significantly less capacity as currently, not more than 1000 transactions per second happen among all blockchains.

Based on the proposed search appliance we estimate that participants will require investing around $1M for dedicated hardware (21 witnesses) and the same amount for backup nodes. Thus overall costs of hardware network infrastructure can be around $2M. after full deployment.

Worth to note that the network doesn't require ultimate configuration at the start and is able to optimize initial investments by the costs of time to index all blockchains. Thus costs at launch can be around $200k. Given that mining industry has been rapidly developed last years this can not be a showstopper for a project. We expect huge interest from miners as slots are limited with 21 fully paid nodes and ~20 of partially paid nodes (depend on the market).

Possible scalability improvements include:
- Hardware. This year Intel Optane [http://www.intel.com/content/www/us/en/architecture-and-technology/intel-optane-technology.html] That creates an opportunity to converge RAM and SSD. Our design has 3-year hardware margin for moving to cheaper and more dense next generation memory.
- Software. The future of consensus computer optimization is in parallel processing. We are going to seriously invest in research of this field. Solving this issue will enable the network to scale nearly infinitely.

## Deployment

We split a process of network deployment into the following milestones which are not bounded to any timeframe at the moment:

1. Exploration Phase
- a white paper published
- CID verification on consensus level
- cyber•rank implementation
- cyber•node release

The purpose of seed stage is to implement the blockchain and prepare it for technical launch.

2. Validation Phase
- Blockchain launched
- cyber•Fund basic application release
- Crowdsale starts.
- Crowdsale ends. Token distribution

The purpose of validation phase is to verify the feasibility of an idea and prospects of technological design across blockchain community. Successful deployment of MVP in the form of basic technical infrastructure in a decentralized fashion and quality of support from the blockchain community and investors around idea will be enough to understand what kind of future the blockchain has. Objective metric is an amount of bitcoins raised during crowd sale for already working blockchain design.

3. Build Phase
- Indexing 10 blockchains
- Indexing 1000 market pairs
- Evaluating 1000 tokens
- GPU cyber•rank calculation implemented
- Historical records of balance valuations is available for indexed blockchain
- Developers run 10 experiments
- First payouts to indexers and auditors

The purpose of build phase is to reach a very *basic* product/market fit around *one specific use case* or this use case will emerge from experiments. A number of payouts which will be calculated based on current capitalization is objective metric. 6 month is expected the duration of phase. If build phase will be successful there are infinite opportunities ahead.

4. Scaling Phase
- Indexing all blockchains
- Indexing Git, BitTorent, IPFS and DAT
- Indexing 10 blockchains
- Autocomplete is fully functional
- Top 1 mln. search queries return useful answers

This is an infinite phase in which the network start continuously grow indexing more and more relevant and meaningful data and the most important answering questions better and for the better. A key scope of work during this stage is to continuously improve developers experience:

> `More indexing => More people search => *More developers build* => More people earn, rank and promote => Better infrastructure => More indexing`.

## The power of cyber•Chain

The key purpose of our proposed design is not just replicate abilities of existing search engines which return only links but enable answering new class of question:

- How much value of X do I possess now?
- What probability of event Y?
- What packages do I need to install in order to improve ROI on available resources?

Our proposed design has all necessary components to bootstrap a market for a new generation of answer applications.

Proposed economics model disintermediate conventional ad model there users are sold to an advertiser and enable any business or people or robot benefit from pure peer-to-peer interactions which bring value for every involved participant.

_Free Market of Indexing and Auditing_. Everybody can connect any blockchain or content-addressable protocol. A decentralized approach to indexing and auditing create an opportunity for those who want to earn on the contributions to cyber•Chain. Proposed solution is not more than a way to _outsource_ these complicated and unstructured efforts for the entire community.

_Free Market of Answering. After all, we have recent advances in machine learning enable to reason about a piece of data quite well. All these algorithms require enormous highly distributed computation which as nearly impossible to achieve in a trestles consensus computer. With the current state of blockchain technology implementing these algorithms using decentralized computational network seems unfeasible. We find a way to _outsource_ this computation for the entire community.

_Self Hosted Search API_. Everybody can deploy self-hosted API. In comparison with what Google offer ($5 per 1000 answers). Our solution can be much more cost effective for high performant applications. One node can process at least 10k queries per second in a read-only mode. That is about 1B requests per month. That is about 100 times cheaper ($0,05 per 1000 answers) even if payback period of search and evaluation appliance ($50k) will be one month. In reality, the affordable payback period is about 10 months. Thus self-hosted search, in theory, can be 1000x more cost effective than Google offering.

## Conclusion

We describe and implement a motivated blockchain based search engine for the permanent web. A search engine is based on the content-addressable peer-to-peer paradigm and uses IPFS as a foundation. IPFS provide significant benefits in terms of resources consumption. CIDs as a primary object is robust in its simplicity. For every CID cyber•rank is computed by a consensus computer with no single point of failure. cyber•rank is a combination of Google's PageRank and Steem's rewards system. cyber•Rank is resistant to Sybil attacks and is computed based on interactions with a graph of CIDs and it's internal relations. Embedded smart contracts offer fair compensations for those who participate in indexing, linking, auditing and curation process. The primary goal is indexing of peer-to-peer systems with self-authenticated data either stateless, such as IFSS, DAT, GIT, BitTorent, or stateful such as Bitcoin, Ethereum and other blockchains and tangles. Proposed market of linking offers necessary incentives for outsourcing computing part responsible for finding meaningful relations between objects. The proposed market of curation and auditing creates essential incentives for ranking high-quality links and objects. Dynamic snippets in search results make possible functionality necessary for the next generation search. Lazy oracles enable indexing of structured publicly verifiable data feeds in a highly competitive environment. A source code of a search engine is open source. Every bit of data accumulated by a blockchain is available for everybody for free. The performance of proposed software-hardware implementation is sufficient for seamless user interactions. Scalability of proposed implementation is enough to index all self-authenticated data that exist today. The blockchain is managed by a decentralized autonomous organization which functions under DPOS consensus algorithm. Thought a system provide necessary utility to offer an alternative for conventional search engines it is not limited to this use case either. The system is extendable for numerous applications and makes possible to design economically rational self-owned robots to spawn a market for AI outsourcing.

## References - not ready yet

[QGTC] <https://www.quora.com/How-many-pages-is-Google-crawling-every-day> PR [] RALF <http://merkle.com/papers/DAOdemocracyDraft.pdf>

ENIGMA <https://www2.deloitte.com/content/dam/Deloitte/uk/Documents/Innovation/deloitte-uk-blockchain-full-report.pdf>

RTB <https://en.wikipedia.org/wiki/Real-time_bidding>

## TODO

_Auditing and Curation_. Probably need more details

_Anonymity_. Explain an economic difference between read search queries and write search queries.
