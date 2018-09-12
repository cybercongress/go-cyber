from typing import Dict

from networkx import DiGraph

from SpringRank import calculate_SpringRank
from adjacency_to_graph import load_edges, print_graph_info
from graph_to_matrix import build_matrix

print("")
print("-----------------------------------------------")
new_edges = {}
load_edges("2700000-2799999_blocks_data", new_edges)


def update_rank(graph: DiGraph, rank: Dict[str, float], n: int) -> (int, dict):
    """
     Updates graph rank for n new links.

    :param graph: graph to add new links
    :param rank: graph current ranks
    :param n: number of new links to add
    :return: (number of iterations, new rank)
    """
    print(f"Updating graph with {n} new links")
    added_edges_count = 0
    for edge in new_edges:
        if added_edges_count > n - 1:
            break
        if graph.has_edge(edge[0], edge[1]):
            graph[edge[0]][edge[1]]['weight'] += new_edges[edge]
        else:
            graph.add_edge(edge[0], edge[1], weight=new_edges[edge])
        added_edges_count += 1

    nodes = list(graph)
    print_graph_info(graph)
    print("-----------------------------------------------")

    print("Regenerate matrix A")
    A = build_matrix(graph, nodes)
    print("-----------------------------------------------")

    print("Computing initial guess")
    initial_x = []
    for node in nodes:
        initial_x.append(rank.get(node, 0))

    print("-----------------------------------------------")

    iterations, raw_rank = calculate_SpringRank(A, initial_x)
    rank = dict(zip(nodes, raw_rank))
    print("-----------------------------------------------")

    return iterations, rank
