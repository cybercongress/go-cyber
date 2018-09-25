from typing import List

import humanize
import numpy as np
import scipy.sparse
import scipy.sparse.linalg
from networkx import DiGraph

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
    result = scipy.sparse.linalg.bicgstab(B, b, x0=initial_x, callback=bicgstab_callback)

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
