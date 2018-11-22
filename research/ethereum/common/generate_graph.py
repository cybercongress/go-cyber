import numpy as np
from scipy.sparse import csr_matrix
from scipy.sparse import lil_matrix


def generate_graph(N, beta=1, c=1) -> (csr_matrix, [float]):
    """
    Generate a graph.
    Edges are drawn from a Poisson distribution, assuming all spring constants are 1
    Input:
        beta:   inversed temperature (noise)
        c:      sparsity constant
    Output:
        A: an adjacency matrix where A[i,j] is the weight of an edge from node i to j
        raw_rank:
    """

    rank = np.random.normal(scale=10, size=N)
    A = lil_matrix((N, N), dtype=float)

    for i in range(N):
        for j in range(N):
            mu = np.exp(-0.5 * beta * (rank[i] - rank[j] - 1) ** 2)
            A[i, j] = np.random.poisson(c * mu)

    return (A.tocsr(), rank)
