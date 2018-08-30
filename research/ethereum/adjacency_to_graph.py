
import networkx as nx
import humanize

import os
import time
from tools import human_readable_time_interval

"""
Takes an adjacency_list like: "23 41 18" or 18 times  "23 41 1"   (edge from 23 --> 41)
possibly having multiple edges and build a graph with no multiple edges but weigths representing how many of them there are
"""


def build_graph_from_adjacency_lists(dataFiles):

    edges = {}

    for dataFile in dataFiles:
        load_edges(dataFile, edges)

    return build_graph(edges)
# End build_graph_from_adjacency_lists


def load_edges(dataFile, edges):

    fileSize = humanize.naturalsize(os.path.getsize(dataFile))
    startTime = time.time()
    linksCount = 0

    print("")
    print("Reading file '{}' {}".format(dataFile, fileSize))
    adjacencyList = open(dataFile, 'r')

    for row in adjacencyList:
        linksCount += 1
        cells = row.split()
        edge = (cells[0], cells[1])
        weight = cells[2]

        if edge not in edges:
            edges[edge] = weight
        else:
            edges[edge] += weight

    adjacencyList.close
    endTime = time.time()
    print("Loading file '{}' finished in {}. {} links added".format(
        dataFile, human_readable_time_interval(endTime - startTime), humanize.intword(linksCount)
    ))
# End load_edges


def build_graph(edges):

    print("")
    print("Building graph.........")

    G = nx.DiGraph()
    startTime = time.time()
    for edge in edges:
        G.add_edge(edge[0], edge[1], weight=edges[edge])
    endTime = time.time()

    print("Graph build in {}".format(human_readable_time_interval(endTime - startTime)))
    return G
# End build_graph
