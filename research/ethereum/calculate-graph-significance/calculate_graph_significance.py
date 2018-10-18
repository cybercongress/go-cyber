import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), '../'))
from common.calculate_significance import test_ranks_significance

import glob
from collections import OrderedDict

from common.adjacency_list_to_graph import load_edges, build_graph
from common.graph_to_matrix import build_matrix

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

print("Calculating Rank Significance")
p_val, H_array = test_ranks_significance(A, plot_file_name='../result/significance.png')
significance_file = open("../result/significance", "w")
significance_file.write(f"p-value: {p_val}\r\n")
for H in H_array:
    significance_file.write(f"{H}\r\n")
significance_file.close()
print("-----------------------------------------------")
