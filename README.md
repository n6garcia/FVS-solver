# Dictionary-Problem

### `What is the smallest set of words that can be used to define every word in the dictionary?`

## Intro

Using the algorithmn that I created to solve this question I was able to define every word in a 110,301 word dictionary by defining only 7,508 words. How do you define a word using only 7,508 words! We define every word in the dictionary by recursively defining words in their definition and replacing them with those recursions. For example the definition for 'handle' could be "the broom stick", in this case we replace the with it's definition, broom with it's defiintion and stick with it's definition. This is unless they are already defined words which (our set of 7,508 words) then we don't recurse on those definitions. We repeatedly do this with all definitions we expand until the recursion ends. Imagine the dictionary as a directed graph G where for all words in the dictionary, a->b meas a defines b or a is in b's definition. The idea is that finite recursion is only possible if the words not in the defined set area are all within a directed acyclic graph (DAG). Without cycles a DFS which is how I implemented my recursive search will always be finite. We try to maximize the acycylic subgraph (MAS) problem by trying to define as few words as possible which means that we are also minimizing the innverse which is the Feedback Vertex Set (FVS). The answer to "what is the smallest set of words that can be used to define every word in the dictionary?" is the minimum FVS of the graph of every words in the dictionary where a->b means a defines b. This is a brand new application of the FVS problem. Hopefully with more work on this problem that truly good applications in fields like ML can be found.

## Directed FVS Approximation Algorithm (My Algorithmn)

Start with any directed graph.

1. Cut any nodes from the graph with no in-degree do this repeatedly until no nodes are found
2. Pick the node with the highest out-degree and cut it from the graph and add it to the set X
3. If graph has no nodes stop else repeat #1-3

The set X is your FVS.

## Website

An application of the FVS algorithm on my website where definitions of any word can be searched up and defined using only words in the FVS.

[coming soon!]

## Dataset(s)

https://www.bragitoff.com/2016/03/english-dictionary-in-csv-format/

WordNet
