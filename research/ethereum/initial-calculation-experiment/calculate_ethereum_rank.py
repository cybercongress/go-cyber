import os
import sys


sys.path.append(os.path.join(os.path.dirname(__file__), '../'))

import glob
from collections import OrderedDict

from common.adjacency_list_to_graph import load_edges, build_graph
from common.calculate_spring_rank import calculate_spring_rank
from common.graph_to_matrix import build_matrix
from common.tools import print_with_time


"""
Calculate ethereum rank from `../data` folders data.
"""

print("")
print("-----------------------------------------------")
print("Loading files with links")

data_files = [file for file in glob.glob("../data/*_blocks_data")]

edges = OrderedDict()
for data_file in data_files:
    new_edges = load_edges(data_file)
    for edge, weight in new_edges.items():
        edges[edge] = edges.get(edge, 0) + weight
print("-----------------------------------------------")

graph = build_graph(edges)
nodes = list(graph)
print("-----------------------------------------------")

A = build_matrix(graph, nodes)
print("-----------------------------------------------")

iterations, raw_rank = calculate_spring_rank(A)  # raw rank is array with values, where indices is nodes list indices
rank = dict(zip(nodes, raw_rank))
print_with_time(f"Spring Rank calculated in {iterations} iterations")
print("-----------------------------------------------")

print("Storing results")
initial_rank_file = open("../result/calculated-rank", "w")
for node, node_rank in rank.items():
    initial_rank_file.write(f"{node} {node_rank}\r\n")
initial_rank_file.close()
print("-----------------------------------------------")
