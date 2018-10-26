import os
import time
from collections import OrderedDict
from typing import Dict, Tuple

import humanize
import networkx as nx
import math
from networkx import DiGraph

from common.tools import human_readable_interval

Edges = Dict[Tuple[str, str], float]
Edge = Tuple[Tuple[str, str], float]


def load_edges(data_file: str) -> Edges:
    """
    Loads graph edges from given file. Supports multi-edges.

    :param data_file: file to load adjacency list
    :return: edges. for multi-graph, resulted edge weight is sum of all same directed weights
    """
    file_size = humanize.naturalsize(os.path.getsize(data_file))
    start_time = time.time()
    links_count = 0

    print("Reading file '{}' {}".format(data_file, file_size))
    adjacency_list = open(data_file, 'r')
    edges = OrderedDict()
    for row in adjacency_list:
        links_count += 1
        cells = row.split()
        if cells[2] == "0":
            continue
        edge = (cells[0], cells[1])
        weight = math.sqrt(float(cells[2]))

        if edge not in edges:
            edges[edge] = weight
        else:
            edges[edge] += weight

    adjacency_list.close()
    end_time = time.time()
    print("Loading file finished in {}. {} links were read".format(
        human_readable_interval(end_time - start_time), humanize.intword(links_count)
    ))
    return edges


def build_graph(edges: Edges) -> DiGraph:
    graph = nx.DiGraph()
    for edge in edges:
        graph.add_edge(edge[0], edge[1], weight=edges[edge])

    print_graph_info(graph)
    return graph


def print_graph_info(graph: DiGraph):
    edges_count = humanize.intword(graph.number_of_edges())
    nodes_count = humanize.intword(graph.number_of_nodes())
    print("Graph contains {} edges for {} nodes".format(edges_count, nodes_count))
