# Minimum Directed Feedback Vertex Set Approximation Algorithm AKA Dictionary Problem

### Motivation/Question (Dictionary Problem) - `What is the smallest set of words that can be used to define every word in the dictionary?`

### Author: Noel Garcia - n6garcia@ucsd.edu - https://noeldev.site/

### Algorithmn Time Complexity : O(N\*Log(N)) : Solves 100k nodes in ~400ms

## Introduction

Using the algorithmn that I found to solve this question I was able to define every word in a 110,301 word dictionary by defining only 7,508 words. We re-define every word in the dictionary by recursively defining words in their definition and replacing them with those recursions. For example the definition for 'handle' could be "the broom stick", in this case we replace 'the' with it's definition, 'broom' with it's defiintion and 'stick' with it's definition. This is unless they are already defined words which (our set of 7,508 words) then we don't recurse on those words. We repeatedly do this with all definitions we expand until the recursion ends. Imagine the dictionary as a directed graph G where for all words in the dictionary, a->b means a defines b or a is in b's definition. The idea is that finite recursion is only possible if the words not in the defined set area are all within a directed acyclic graph (DAG). Without cycles a DFS which is how I implemented my recursive search will always be finite. We try to maximize the acycylic subgraph (MAS) problem by trying to define as few words as possible which means that we are also minimizing the inverse which is the Feedback Vertex Set (FVS). The answer to our original question is the minimum FVS of the graph of all words in the dictionary where a->b means a defines b. This is a brand new application of the FVS problem. Hopefully with more work on this problem that truly good applications in fields like ML can be found. I have also theorized that the inverse of the Dictionary Problem could be used for data compression, something I think could be called quantification.

## Preliminaries

Let's explain the general setup for the graph from the dictionary. let W be a set of words and w ∈ W then def(w) = S where S is the set of unique words in the definition of the word w. LET G = (V,E) where V = (v ∈ W ∩ def(W)) and E = (eij | i ∈ def(w), j ∈ w where w ∈ W). This will exactly give you the graph G where one can find "the smallest set of words" from a minimum FVS. Let's define FVS just so we are completely clear on that, a subset S of V(G) is a directed Feedback Vertex Set (FVS), if the induced subgraph V(G) \ S is acyclic.

## Min DFVS Approximation Algorithm

Start with any directed graph G.

1. Cut any nodes from the graph with no in-degree do this repeatedly until no nodes are found
2. Pick the node with the highest out-degree and cut it from the graph and add it to the set X
3. If graph has no nodes stop else repeat #1-3

The set X is your FVS.

Disclaimer: I created and developed this algorithm entirely with litte knowledge of the sparse literature on DFVS approximation.
I created and began work on this algorithm in time with the PACE 2022 DFVS competition and I have found my exact algorithm published by Levy & Low on 1988 in [1].

## Website

An application of the FVS algorithm on my website where definitions of any word can be searched up and defined using only words in the FVS.

https://noeldev.site/search/search.html

## Dataset(s)

https://www.bragitoff.com/2016/03/english-dictionary-in-csv-format/

WordNet®

## Reference(s)

[1] Hanoch Levy, David W Low, A contraction algorithm for finding small cycle cutsets, Journal of Algorithms, Volume 9, Issue 4, 1988, Pages 470-493,
    ISSN 0196-6774, https://doi.org/10.1016/0196-6774(88)90013-2. (https://www.sciencedirect.com/science/article/pii/0196677488900132)
