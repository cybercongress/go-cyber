### Rank State

Rank State = rank values (stored in one dimensional array) and merkle tree of those values.

Each new CID gets unique number. Number starts from zero and incrementing by one for each new CID.
So that we can store rank in one dimensional array where indices are CID numbers.

### Rank Merkle Tree

Calculated based on RFC-6962 standard (https://tools.ietf.org/html/rfc6962#section-2.1).

Since rank stored in one dimensional array where indices are CID numbers (we could say that it ordered by CID numbers)
leafs in merkle tree from left to right are `SHA-256` hashes of rank value. Index of leaf is CID number. It helps to
easily find proofs for specified CID (`log n` iterations where `n` is number of leafs).

To store merkle tree is necessary to split tree on subtrees with number of leafs multiply of power of 2.
The smallest one is obviously subtree with only one leaf (and therefore `height == 0`).

Leaf addition looks as follows. Each new leaf is added as subtree with `height == 0`.
Then sequentially merge subtrees with the same `height` from right to left.

Example:

                                            
      ┌──┴──┐   │      ┌──┴──┐     │         ┌──┴──┐     │   │
    ┌─┴─┐ ┌─┴─┐ │    ┌─┴─┐ ┌─┴─┐ ┌─┴─┐     ┌─┴─┐ ┌─┴─┐ ┌─┴─┐ │
    (5-leaf)         (6-leaf)              (7-leaf)

To get merkle root hash - join subtree roots from right to left.

#### How to store merkle tree?

<b>Full tree</b> - all subtrees with all leafs and intermediary nodes  
<b>Short tree</b> - contains only subtrees roots

The trick is that <b>full tree</b> is only necessary for providing merkle proofs. 
For consensus purposes and updating tree it's enough to have <b>short tree</b>.

To store merkle tree in database use only <b>short tree</b>.
Marshaling of short tree with `n` subtrees (each subtree takes 40 bytes):  
`<subtree_1_root_hash_bytes><subtree_1_height_bytes>....<subtree_n_root_hash_bytes><subtree_n_height_bytes>`

For `1,099,511,627,775` leafs <b>short tree</b> would contain only 40 subtrees roots and take only 1600 bytes.

### Rank State Calculation

Lets denote:
  
`p` - rank calculation period  
`lbn` - last confirmed block number  
`cbn` - current block number  
`lr` -  length of rank values array  

For rank storing and calculation we have two separate in-memory contexts:
1. Current rank context. It includes last calculated rank state (values and merkle tree) plus
all links and user stakes submitted to the moment of this rank submission.
2. New rank context. It's currently calculating (or already calculated and waiting for submission) rank state.
Consists of new calculated rank state (values and merkle tree) plus new incoming links and updated user stakes.

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

To sum up.  
In <b>current rank context</b> we have rank state from last calculated iteration (plus, every block, it updates with new CIDs).
And we have links and user stakes that are participating in current rank calculation iteration (whether it finished or not).
<b>New rank context</b> contains links and stakes that will go to next rank calculation and newly calculated rank state
(if calculation is finished) that waiting for submitting.

#### Node restart

Firstly, need to restore both contexts (current and new).
Load links and user stakes from database using different versions:
1. Links and stakes from last calculated rank version `v = lbn - (lbn mod n)` goes to current rank context.
2. Links and stakes between versions `v` and `lbn` goes to new rank context.

Also to restart node correctly we have to store next entities in database:

1. Last calculated rank hash (merkle tree root)
2. Newly calculated rank short merkle tree
3. Last block short merkle tree

With <b>last calculated rank hash</b> and <b>newly calculated rank merkle tree</b> we could check if rank calculation
was finished before node restart. If they are equal then rank wasn't calculated and we should run rank calculation.
If not we could skip rank calculation and use <b>newly calculated rank merkle tree</b> to participate in consensus
when it comes to block number `cbn ≡ 0 (mod p)` (rank values will not be available until rank calculation happens
in next iteration. Still validator can participate in consensus so nothing bad).

<b>Last block merkle tree</b> necessary to participate in consensus till start of next rank calculation iteration.

So, after restart we could end up with two states:
1. Restored current rank context and new rank context without rank values (links, user stakes and merkle tree).
2. Restored current rank context without rank values. Restored new rank context only with links and user stakes.

Node is able to participate in consensus but cannot provide rank values (and merkle proofs)
till two rank calculation iterations finished (current and next).

### Search Index

Index should be ran in parallel and do not influence work of consensus machine.
Validator should be able to turn off index support. May be even make it a separate daemon.

<b>Base idea.</b> Always submit new links to index and take rank values from current context (insert in sorted array operation).
When new rank state is submitted trigger index to update rank values and do sortings
(in most cases new arrays will be almost sorted).

Need to solve problem of adjusting arrays capacity (to not copy arrays each time new linked cid added).
Possible solution is to adjust capacity with reserve before resorting array.

Therefore for building index we need to find sorting algorithm that will be fast on almost sorted arrays. 
Also we should implement it for GPU (so it should better be parallelizable). Mergesort(Timsort), Heapsort, Smoothsort???