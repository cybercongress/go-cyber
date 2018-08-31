import glob
import time
import networkx as nx
import humanize
import SpringRank as sr

from tools import human_readable_time_interval
from adjacency_to_graph import build_graph_from_adjacency_lists


dataFiles = [file for file in glob.glob("data/*_blocks_data")]

print("")
print("-----------------------------------------------")
print("About to download calculate SpringRank for {} files".format(dataFiles))

graph = build_graph_from_adjacency_lists(dataFiles)
nodes = list(graph.nodes())
edgesNumber = humanize.intword(graph.number_of_edges())
nodesNumber = humanize.intword(graph.number_of_nodes())

print("Graph contains {} edges for {} nodes".format(edgesNumber, nodesNumber))
print("")

alpha=0.
l0=1.
l1=1.

print("Building Matrix......")
A = nx.to_numpy_matrix(graph, nodelist=nodes)
print("Calculating Rank......")

startTime = time.time()
rank=sr.SpringRank(A)
endTime = time.time()
print("Rank calculated in {}".format(human_readable_time_interval(endTime - startTime)))

