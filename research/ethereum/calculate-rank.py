import glob

from SpringRank import calculate_SpringRank
from adjacency_to_graph import load_edges, build_graph
from graph_to_matrix import build_matrix

print("")
print("-----------------------------------------------")
print("Loading files with links")

data_files = [file for file in glob.glob("data/*_blocks_data")]

edges = {}
for data_file in data_files:
    load_edges(data_file, edges)
print("-----------------------------------------------")

graph = build_graph(edges)
nodes = list(graph.nodes())
print("-----------------------------------------------")

A = build_matrix(graph, nodes)
print("-----------------------------------------------")

iterations, raw_rank = calculate_SpringRank(A)  # raw rank is array with values, where indices is nodes list indices
print(f"Spring Rank calculated in {iterations} iterations")
print("-----------------------------------------------")
