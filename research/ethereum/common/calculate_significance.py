import matplotlib
import matplotlib.pyplot as plt
import numpy as np
import scipy.optimize as op
import scipy.sparse
import scipy.sparse.linalg
from scipy.sparse import csr_matrix

from common.calculate_spring_rank import calculate_Hamiltonion, calculate_spring_rank


# noinspection PyTypeChecker
def estimate_beta(A: csr_matrix, rank: [float]) -> float:
    """
    Use Maximum Likelihood Estimate to estimate inversed temperature beta from a network and its inferred rankings/scores
    Input:
        A: graph adjacency matrix where matrix[i,j] is the weight of an edge from node i to j
        rank: SpringRank inferred scores/rankings of nodes
    Output:
        beta: MLE estimate of inversed tempurature
    """

    def f(beta, A, rnk):
        n = rnk.size
        y = 0.

        for i in range(n):
            for j in range(n):
                d = rnk[i] - rnk[j]
                p = (1 + np.exp(-2 * beta * d)) ** (-1)
                y = y + d * (A[i, j] - (A[i, j] + A[j, i]) * p)

        return y

    normed_rank = rank + abs(min(rank))
    beta_0 = 0.1
    return op.fsolve(f, beta_0, args=(A, normed_rank))[0]


def test_ranks_significance(A: csr_matrix, n_repetitions=100, plot_file_name=None):
    """
    Given an adjacency matrix, test if the hierarchical structure is statitically significant compared to null model
    The null model contains randomized directions of edges while preserving total degree between each pair of nodes
    Input:
        A: graph adjacency matrix where matrix[i,j] is the weight of an edge from node i to j
        n_repetitions: number of null models to generate
        plot: if True shows histogram of null models' energy distribution
    Output:
        p-value: probability of observing the Hamiltonion energy lower than that of the real network if null model is true
        plot: histogram of energy of null models with dashed line as the energy of real network
    """

    # place holder for outputs
    H = np.zeros(n_repetitions)
    H0 = calculate_Hamiltonion(A, calculate_spring_rank(A)[1])

    # generate null models
    i = 0
    while i < n_repetitions:
        try:
            B = randomize_edge_direction(A)
            H[i] = calculate_Hamiltonion(B, calculate_spring_rank(B)[1])
            i += 1
        except:
            print("Filed to calculate Spring Rank")
    p_val = np.sum(H < H0) / n_repetitions

    # Plot
    if plot_file_name:
        matplotlib.use("TkAgg")
        plt.hist(H)
        plt.axvline(x=H0, color='r', linestyle='dashed')
        plt.savefig(plot_file_name)

    return (p_val, H)


def randomize_edge_direction(A: csr_matrix) -> csr_matrix:
    """
    Randomize directions of edges while preserving the total degree.
    Used when values in A are not integers
    Input:
        A: graph adjacency matrix where A[i,j] is the weight of an edge from node i to j
    Output:
        An adjacency matrix representing a null model
    """

    n = A.shape[0]
    (r, c, v) = scipy.sparse.find(A)
    for i in range(v.size):
        if np.random.random_sample() > 0.5:
            temp = r[i]
            r[i] = c[i]
            c[i] = temp
    return scipy.sparse.csr_matrix((v, (r, c)), shape=(n, n))
