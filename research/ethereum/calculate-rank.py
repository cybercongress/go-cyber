import glob
import time
import networkx as nx
import humanize

from adjacency_to_graph import build_graph_from_adjacency_lists


dataFiles = [file for file in glob.glob("data/*_blocks_data")]

print("")
print("-----------------------------------------------")
print("About to download calculate SpringRank for {} files".format(dataFiles))

graph = build_graph_from_adjacency_lists(dataFiles)
edgesNumber = humanize.intword(graph.number_of_edges())
nodesNumber = humanize.intword(graph.number_of_nodes())

print("Graph contains {} edges for {} nodes".format(edgesNumber, nodesNumber))
print("")
