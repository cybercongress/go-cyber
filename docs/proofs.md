Since rank is stored within a one-dimensional array, where indices are the CID numbers (we can say that they are ‘managed’ by the CID numbers) the leaves of the Merkle tree, from left to right, are \code{SHA-256} hashes of the rank value. The index of the leaf is the CID number. This helps us to easily find proofs for specified CID (\code{log n} iterations where \code{n} is the number of leaves).

To store a Merkle tree, it is necessary to split the tree further into subtrees, with the number of leaves multiplied by the power of 2. The smallest tree is a subtree with just one leaf (and therefore \code{height == 0}). Leaf addition looks as follows: each new leaf is added as a subtree, with \code{height == 0}. Then, subtrees are sequentially merged with the same \code{height} from right to left.

To obtain a hash of a Merkle root one needs to link subtree roots from right to left.

The Merkle tree rank can be stored variously:

\begin{itemize}
\item[] Full tree - all subtrees with all leaves and intermediary nodes
\item[] Short tree - contains only the subtree roots
\end{itemize}

The trick is that a \verb|full tree| is only necessary for providing Merkle proofs. For consensus purposes and for updating the tree, it is enough to have a \verb|short tree|. To store a Merkle tree in the database one has to use a \verb|short tree|. Marshaling of a short tree with \code{n} subtrees (each subtree takes 40 bytes):

\begin{lstlisting}
<subtree_1_root_hash_bytes><subtree_1_height_bytes>
....
<subtree_n_root_hash_bytes><subtree_n_height_bytes>
\end{lstlisting}

For every \code{1,099,511,627,775} leaves, a \verb|short tree| will only contain 40 subtree roots and will only take 1600 bytes.

Let us express rank state calculation:

\begin{itemize}
    \item[] \code{p} - rank calculation period
    \item[] \code{lbn} - last confirmed block number
    \item[] \code{cbn} - current block number
    \item[] \code{lr} -  length of rank-values array
\end{itemize}

For storing rank and for its calculation, we have two separate in-memory circumstances:

\begin{enumerate}
\item Current rank settings: this includes the last calculated rank state (values and Merkle tree), plus,
all the links and all the agents' stakes that are submitted for the moment of the current rank submission.
\item New rank settings: it is the currently calculated (or the already calculated and waiting for a submission) rank state. It consists of a new calculated rank state (values and Merkle tree), plus, of new and incoming links, and the updated agents' stakes.
\end{enumerate}

Calculation of a new rank state occurs once per \code{p} blocks and proceed in parallel.

The iteration starts from the block number that \code{$\equiv$ 0 (mod p)} and proceeds until the next block number that \code{$\equiv$ 0 (mod p)}.

For block number \code{cbn $\equiv$ 0 (mod p)} (this includes block number 1, as blocks in Cosmos start from 1):

\begin{enumerate}
  \item Check if rank calculation is done. If yes then go to (2.) if not - wait until the calculation is completed
  (actually, this situation should not happen because this would mean that the rank calculation period is too short)
  \item Submit rank, links and agents' stakes move from 'new rank settings' to 'current rank settings'
  \item Store the latest calculated rank Merkle tree root hash
  \item Start new rank calculations in parallel (for links and the stakes from the current rank settings)
\end{enumerate}

For each block:

\begin{enumerate}
  \item All links advance to 'new rank settings'
  \item New, incoming CIDs have a rank that equals to zero. We could check this with the number of the last CIDs and \code{lr} (it equals to the number of CIDs that already hold a rank). Then, add the CIDs with number \code{>lr} to the end of this array with a value equals to zero
  \item Update the current context Merkle tree with CIDs from the previous step
  \item Store the latest Merkle tree from the current settings (let us name this 'last block Merkle tree')
  \item Check if the new rank calculations are completed. If yes, go to (4.) if not, continue to the next block
  \item Push the calculated rank state to 'new rank settings'. Store the Merkle tree of the newly calculated rank
\end{enumerate}
To sum up: in \textit{current rank settings}, we have a rank state from the last calculated iteration (plus, it updates with new CIDs every block). Additionally, we have links and agents' stakes that are participating in the current rank calculation (whether they are complete or not). The \textit{new rank settings} contains links and stakes that will roll to the next rank calculation and the newly calculated rank state (if a calculation is completed) that awaits submission.

If for some reason, we first need to restart the node, then we need to restore both of the settings (the current and the new).
Load the links and the agents' stakes from the database by using one of the following:

\begin{enumerate}
  \item Links and stakes from the last calculated rank version \code{v = lbn - (lbn mod n)} go to current rank settings
  \item Links and stakes that are between versions \code{v} and \code{lbn} go to new rank settings
\end{enumerate}

To restart a node correctly we have to store the following entities within the database:

\begin{enumerate}
  \item The last calculated rank hash (Merkle tree root)
  \item A newly calculated rank - short Merkle tree
  \item Last blocks' short Merkle tree
\end{enumerate}

With the \textit{last calculated rank hash}, and the \textit{newly calculated rank Merkle tree} we can now check if the rank
calculations were done before the restart of the node. If they are equal, then the rank wasn't calculated and we should run the rank calculation again.
If they are not equal we can skip rank calculations and use the \textit{newly calculated rank Merkle tree} in order to participate in the consensus when it reaches block number \code{cbn $\equiv$ 0 (mod p)} (note that rank values will not be available until the rank calculations happen in the next iteration. The validator is still able to participate in the consensus so nothing bad actually happened).

The \textit{last block Merkle tree} is necessary to participate in the consensus until the start of the next rank calculation. After a restart, we could end up with two states:

\begin{enumerate}
\item Restored current rank settings and new rank settings without rank values (links, agent stakes, and Merkle tree)
\item Restored current rank settings without rank values. Restored new rank settings with only the links and the agents' stakes
\end{enumerate}

A node can participate in the consensus but cannot provide rank values (and Merkle proofs) until two iterations of rank calculation are achieved (the current and the next). The search index should be run in parallel and will not influence the work of the consensus machine. The validator should be able to turn off the index support.
