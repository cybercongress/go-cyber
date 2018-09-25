# First SpringRank Research on Ethereum addresses

## Model

Nodes - Ethereum Addresses. Links from addresses txes. Link direction: for tx A -> B link is B -> A and vice versa. 
Links weight is total ETH amout of same-directed txes.

## Experiments

### Initial Calculation

Goal: Calculate ethereum rank for fist **N** blocks.

1. Run local parity node
2. Execute script from **common** `python3 ethereum_chain_to_adjacency_list.py 0 {N}`
3. Execute script from **initial-calculation-experiment** 'python3 calculate_ethereum_rank.py'

### Initial Guess experiment

Goal: Calculate initial guess performance impact.

1. Execute **Initial Calculation** experiment
2. Execute script from **common** `python3 ethereum_chain_to_adjacency_list.py {N+1} {M}`. Where [N+1,M] blocks should 
contains at least 1kk txes.
3. Copy file from previous step into **initial-guess-experiment** folder.
4. Execute script from **initial-calculation-experiment** 
'python3 calculate_initial_guess_performance_impact.py 1000 1000 "1366461-2000000_blocks_data"' 
5. Plot results using command "python3 plot.py"