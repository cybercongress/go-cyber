import copy
import glob
import os
import sys
from collections import OrderedDict

sys.path.append(os.path.join(os.path.dirname(__file__), '../'))

from common.calculate_spring_rank import update_rank
from common.tools import ranks_diff
from common.adjacency_list_to_graph import load_edges, build_graph

"""
Usage:

1) Calculate initial rank using "Initial calculation experiment"
2) Download more edges from ethereum network using common module.
3) Run script `python calculate_initial_guess_performance_impact.py 1000 1000 "./1366461-2000000_blocks_data"`
4) Plot data

"""

delta_edges_number = int(sys.argv[1])
iterations_number = int(sys.argv[2])
file_with_new_links = str(sys.argv[3])

print("")
print("-----------------------------------------------")
print("Restoring results")

data_files = [file for file in glob.glob("../data/*_blocks_data")]

edges = OrderedDict()
for data_file in data_files:
    new_edges = load_edges(data_file)
    for edge, weight in new_edges.items():
        edges[edge] = edges.get(edge, 0) + weight

graph = build_graph(edges)

initial_rank = OrderedDict()
initial_rank_file = open("../data/initial-rank", "r")
for row in initial_rank_file:
    cells = row.split()
    initial_rank[cells[0]] = cells[1]
initial_rank_file.close()
print("-----------------------------------------------")
print("Loading new edges")

new_edges = [(e, w) for e, w in load_edges(file_with_new_links).items()]
print("-----------------------------------------------")
print(f"Calculating initial guess performance impact for {iterations_number} iterations on {delta_edges_number} delta")

iterations = []
mean_rank_diffs = []
average_rank_diffs = []

new_rank = initial_rank
for i in range(0, iterations_number):
    print("-----------------------------------------------")
    print(f"Iteration {i}: ")
    prev_rank = copy.deepcopy(new_rank)
    start_offset = i * delta_edges_number
    new_edges_for_iteration = new_edges[start_offset: start_offset + delta_edges_number - 1]
    try:
        new_iterations, new_rank = update_rank(graph, prev_rank, new_edges_for_iteration)
        mean_rank_diff, average_rank_diff = ranks_diff(prev_rank, new_rank)

    except:
        mean_rank_diff = -2
        average_rank_diff = -2
        new_iterations = -100

    mean_rank_diffs.append(mean_rank_diff)
    average_rank_diffs.append(average_rank_diff)
    iterations.append(new_iterations)

print("-----------------------------------------------")
print(f"Storing results")
recalculate_rank_file = open("initial_guess_impact_results", "w")
for i in range(0, iterations_number):
    recalculate_rank_file.write(f"{iterations[i]} {mean_rank_diffs[i]} {average_rank_diffs[i]}\r\n")
recalculate_rank_file.close()
print("-----------------------------------------------")
