import os

import humanize
import networkx as nx
import time

from tools import human_readable_interval
from tools import record_execution_time

"""
Takes an adjacency_list like: "23 41 18" or 18 times  "23 41 1"   (edge from 23 --> 41)
possibly having multiple edges and build a graph with no multiple edges but weigths representing how many of them there are
"""


def load_edges(data_file, edges):
    file_size = humanize.naturalsize(os.path.getsize(data_file))
    start_time = time.time()
    links_count = 0

    print("Reading file '{}' {}".format(data_file, file_size))
    adjacency_list = open(data_file, 'r')

    for row in adjacency_list:
        links_count += 1
        cells = row.split()
        edge = (cells[0], cells[1])
        weight = cells[2]

        if edge not in edges:
            edges[edge] = weight
        else:
            edges[edge] += weight

    adjacency_list.close()
    end_time = time.time()
    print("Loading file finished in {}. {} links added".format(
        human_readable_interval(end_time - start_time), humanize.intword(links_count)
    ))


@record_execution_time(
    before_message="Building graph",
    after_message="Graph build in {}"
)
def build_graph(edges):
    graph = nx.DiGraph()
    for edge in edges:
        graph.add_edge(edge[0], edge[1], weight=edges[edge])

    edges_count = humanize.intword(graph.number_of_edges())
    nodes_count = humanize.intword(graph.number_of_nodes())
    print("Graph contains {} edges for {} nodes".format(edges_count, nodes_count))
    return graph
