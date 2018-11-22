from typing import List

import humanize
import numpy as np
import scipy.sparse
import scipy.sparse.linalg
from networkx import DiGraph
from scipy.sparse import csr_matrix

from common.adjacency_list_to_graph import print_graph_info, Edge
from common.graph_to_matrix import build_matrix
from common.tools import Rank
from common.tools import print_with_time


def calculate_spring_rank(A, initial_x=None):
    """
    Main routine to calculate SpringRank by solving linear system.
    Default parameters are initialized as in the standard SpringRank model.

    Matrix equation to solve: Bx=b --> [Dout+Din−A_o]x=[Dout−Din]*1+(doutN−dinN)*1
    where B = [Dout+Din−A_o]=[Dout+Din−A-A.T-Anj-Ajn]
    where b = [Dout−Din]*1+(doutN−dinN)*1

    :param initial_x: initial solution guess (solution start point)
    :param A: network adjacency matrix (can be weighted)
    :return: N-dim array, indices represent the nodes' indices used in ordering the matrix A
    """

    N = A.shape[0]  # N is matrix size - NxN. shape is [N,N] array
    One = np.ones(N)  # [1,1,1..., 1]

    d_in_matrix = A.sum(axis=0)  # returns 2-array [[sum(col_1), sum(col_2), ..., sum(col_N)]]
    d_out_matrix = A.sum(axis=1)  # returns 2-array [[sum(row_1)], [sum(row_2)], ...., [sum(row_N)]]

    d_in = [d_in_matrix[0, j] for j in range(N)]  # [sum(col_1), sum(col_2), ..., sum(col_N)]
    d_out = [d_out_matrix[i, 0] for i in range(N)]  # [sum(row_1), sum(row_2), ..., sum(row_N)]

    D_int = scipy.sparse.diags(d_in)  # NxN matrix, where elements of d_in array located on diagonals and other are 0
    D_out = scipy.sparse.diags(d_out)  # NxN matrix, where elements of d_out array located on diagonals and other are 0

    # get last row of A and create new matrix NxN matrix, where each row is last row of A
    print_with_time("Calculating Anj ....")
    A_N_j = scipy.sparse.vstack([A.getrow(N - 1)] * N)

    # get last column of A and create new matrix NxN matrix, where each column is last column of A
    print_with_time("Calculating Ajn ....")
    A_j_N = scipy.sparse.hstack([A.getcol(N - 1).tocsc()] * N, format="csc")

    print_with_time("Calculating A_o ....")
    A_o = A + A.transpose() + A_N_j + A_j_N

    print_with_time("Calculating B ....")
    B = D_out + D_int - A_o
    size_of_B = humanize.naturalsize(B.data.nbytes + B.indptr.nbytes + B.indices.nbytes)
    print_with_time("Matrix B takes {} RAM".format(size_of_B))

    print_with_time("Calculating b ....")
    b = (D_out - D_int) * One + (d_out[N - 1] - d_in[N - 1]) * One

    # ----------------------------------------------------------------------------
    iterations = 0

    def bicgstab_callback(x):
        nonlocal iterations
        iterations += 1

    print_with_time("Solving Bx=b equation using 'bicgstab' iterative method")
    result = scipy.sparse.linalg.bicgstab(B, b, x0=initial_x, callback=bicgstab_callback, tol=1e-3)

    if result[1] != 0:
        print_with_time("Can't solve Bx=b")
        raise ArithmeticError("Can't solve Bx=b")

    return iterations, result[0]


def update_rank(graph: DiGraph, rank: Rank, new_edges: List[Edge]) -> (int, Rank):
    """
     Updates graph rank for n new links.

    :param graph: graph to add new links
    :param rank: graph current ranks
    :param n: number of new links to add
    :return: (number of iterations, new rank)
    """
    for edge, weight in new_edges:
        if graph.has_edge(edge[0], edge[1]):
            graph[edge[0]][edge[1]]['weight'] += weight
        else:
            graph.add_edge(edge[0], edge[1], weight=weight)

    nodes = list(graph)
    print_graph_info(graph)

    A = build_matrix(graph, nodes)

    initial_x = []
    for node in nodes:
        initial_x.append(rank.get(node, 0))

    iterations, raw_rank = calculate_spring_rank(A, initial_x)
    rank = dict(zip(nodes, raw_rank))

    return iterations, rank


def calculate_violations(A: csr_matrix, rank: [float]) -> (float, float, float):
    """
    Calculate number of violations in a graph given SpringRank scores
    A violaton is an edge going from a lower ranked node to a higher ranked one
    Input:
        A: graph adjacency matrix where A[i,j] is the weight of an edge from node i to j
        rank: SpringRank result
    Output:
        number of violations
        proportion of all edges against violations
    """

    rank_sort = np.argsort(rank)[::-1]  # Get sorted by rank A indices
    A_sort = A[rank_sort,][:, rank_sort]  # Create new matrix, where elements sorted by rank desc

    # All elements below main triangle is connection from low-ranked node
    # to high-ranked node. So to get all violation just sum all elements
    # below main triangle.
    viol = scipy.sparse.tril(A_sort).sum()

    m = A_sort.sum()  # All matrix weights sum

    return (viol, viol / m, m)


def calculate_min_violations(A: csr_matrix) -> (float, float):
    """
    Calculate the minimum number of violations in a graph for all possible rankings
    A violaton is an edge going from a lower ranked node to a higher ranked one
    Minimum number is calculated by summing bidirectional interactions.
    Input:
        A: graph adjacency matrix where A[i,j] is the weight of an edge from node i to j
    Output:
        minimum number of violations
        proportion of all edges against minimum violations
    """

    ii, ji, v = scipy.sparse.find(A)  # I,J,V contain the row, column indices, and values of the nonzero entries.

    min_viol = 0.0
    for e in range(len(v)):  # for all nodes interactions
        i, j = ii[e], ji[e]
        if A[i, j] > 0 and A[j, i] > 0:
            min_viol = min_viol + min(A[i, j], A[j, i])

    m = A.sum()
    return (min_viol, min_viol / m)


def calculate_system_violated_energy(A: csr_matrix, rank: [float]) -> (float, float, float):
    """
    Calculate number of violations in a graph given SpringRank scores
    A violaton is an edge going from a lower ranked node to a higher ranked one
        weighted by the difference between these two nodes
    Input:
        A: graph adjacency matrix where A[i,j] is the weight of an edge from node i to j
        rank: SpringRank scores
    Output:
        system violated energy
        proportion of system energy against system violated energy
    """
    i, j, v = scipy.sparse.find(A)  # I,J,V contain the row indices, column indices, and values of the nonzero entries.
    normed_rank = (rank - min(rank)) / (max(rank) - min(rank))  # normalize
    wv = 0.0
    for e in range(len(v)):  # for all nodes interactions
        if normed_rank[i[e]] < normed_rank[j[e]]:  # compare ranks of two nodes i and j
            wv = wv + v[e] * (normed_rank[j[e]] - normed_rank[i[e]])

    H = calculate_Hamiltonion(A, rank)
    return (wv, wv / H, H)


def calculate_Hamiltonion(A: csr_matrix, rank: [float]) -> float:
    """
    Calculate the Hamiltonion of the network
    Input:
        A: graph adjacency matrix where matrix[i,j] is the weight of an edge from node i to j
        rank: SpringRank scores
    Output:
        H: Hamiltonion energy of the system
    """

    H = 0.0
    normed_rank = (rank - min(rank)) / (max(rank) - min(rank))  # normalize

    ii, ji, v = scipy.sparse.find(A)  # I,J,V contain the row, column indices, and values of the nonzero entries.
    for e in range(len(v)):  # for all nodes interactions
        i, j = ii[e], ji[e]
        H = H + 0.5 * A[i, j] * (normed_rank[i] - normed_rank[j] - 1) ** 2

    return H
