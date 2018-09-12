# First SpringRank Research on Ethereum addresses

## Task

1. Calculate SpringRank for simplest possible graph. As a graph ethereum transactions were choosen.
2. Test perfomance of first time SpringRank calculation.
3. Test perfomance of adjusting rank by new nodes/links including.

## Model

Nodes - Ethereum Addresses. Links from addresses txes. Link direction: for tx A -> B link is B -> A and vice versa. Links weight is total ETH amout of same-directed txes.

## Current values

As for 29.08.18 Ethereum network  consists of 300.851 millions transaction and 63.225 millions addresses. //todo add link for cybernode.live


## Generate random matrix and calculate SpringRank
```
A = scipy.sparse.rand(100000, 100000, density=0.0015, format="csr", random_state=42)
print(A.todense())
calculate_SpringRank(A)
```