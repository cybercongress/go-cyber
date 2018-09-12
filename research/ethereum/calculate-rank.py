import copy
import glob

from SpringRank import calculate_SpringRank
from adjacency_to_graph import load_edges, build_graph
from graph_to_matrix import build_matrix
from recalculate_rank import update_rank
from tools import ranks_diff

print("")
print("-----------------------------------------------")
print("Loading files with links")

data_files = [file for file in glob.glob("data/*_blocks_data")]

edges = {}
for data_file in data_files:
    load_edges(data_file, edges)
print("-----------------------------------------------")

graph = build_graph(edges)
nodes = list(graph)
print("-----------------------------------------------")

A = build_matrix(graph, nodes)
print("-----------------------------------------------")

iterations, raw_rank = calculate_SpringRank(A)  # raw rank is array with values, where indices is nodes list indices
rank = dict(zip(nodes, raw_rank))
print(f"Spring Rank calculated in {iterations} iterations")
print("-----------------------------------------------")

print("Storing results")
initial_rank_file = open("initial-rank", "w")
for node, node_rank in rank.items():
    initial_rank_file.write(f"{node} {node_rank}\r\n")
initial_rank_file.close()
print("-----------------------------------------------")

initial_count = 100
final_count = 1_000
step = 100
print("Going to recalculate rank")
recalculate_rank_file = open(f"recalculating-rank-info-{step}", "w")
recalculate_rank_file.write(f"0 {iterations} 0 0")
for new_edges_count in range(initial_count, final_count, step):
    iterations, new_rank = update_rank(copy.deepcopy(graph), copy.deepcopy(rank), new_edges_count)
    mean_rank_diff, average_rank_diff = ranks_diff(rank, new_rank)
    recalculate_rank_file.write(f"{new_edges_count} {iterations} {mean_rank_diff} {average_rank_diff}\r\n")
    print(f"For {new_edges_count} new links SpringRank recalculated in {iterations} iterations")
    print(f"Mean ranks diff {mean_rank_diff}, average ranks diff {average_rank_diff}")
recalculate_rank_file.close()
print("-----------------------------------------------")
