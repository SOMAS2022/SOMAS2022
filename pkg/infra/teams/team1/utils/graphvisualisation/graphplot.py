import networkx as nx
import matplotlib.pyplot as plt
import numpy as np
import re
import pydot
   
  
# Defining a Class
class GraphVisualization:
   
    def __init__(self):
          
        # visual is a list which stores all 
        # the set of edges that constitutes a
        # graph
        self.visual = []
          
    # addEdge function inputs the vertices of an
    # edge and appends it to the visual list
    def addEdge(self, a, b):
        temp = [a, b]
        self.visual.append(temp)
          
    # In visualize function G is an object of
    # class Graph given by networkx G.add_edges_from(visual)
    # creates a graph with a given list
    # nx.draw_networkx(G) - plots the graph
    # plt.show() - displays the graph
    def visualize(self):
        G = nx.Graph()
        G.add_edges_from(self.visual)
        nx.nx_pydot.write_dot(G, "fig.dot")
        nx.draw_networkx(G)
        #plt.savefig("fig.png")
  
# Driver code
G = GraphVisualization()
f = open('input.txt', 'r')
data_txt = f.read()
f.close()
regex = "[0-9]+"
my_nparray = np.array(re.findall(regex, data_txt), dtype=int)
zipped = zip(my_nparray[::2], my_nparray[1::2])

for a, b in zipped:
    G.addEdge(a, b)

G.visualize()

# With graphviz installed, run
# dot -Kneato -Tsvg fig.dot > fig.svg