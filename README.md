# Dictionary-Problem

### `What is the smallest set of words that can be used to define every word in the dictionary?`

### Exclaimer: if you steal my work you will go to hell :3. billionare proj loading...

### Note: If you like the mind expanding material leave a star. :3

## Introduction

Using the algorithmn that I created to solve this question I was able to define every word in a 110,301 word dictionary by defining only 7,508 words. We re-define every word in the dictionary by recursively defining words in their definition and replacing them with those recursions. For example the definition for 'handle' could be "the broom stick", in this case we replace 'the' with it's definition, 'broom' with it's defiintion and 'stick' with it's definition. This is unless they are already defined words which (our set of 7,508 words) then we don't recurse on those words. We repeatedly do this with all definitions we expand until the recursion ends. Imagine the dictionary as a directed graph G where for all words in the dictionary, a->b means a defines b or a is in b's definition. The idea is that finite recursion is only possible if the words not in the defined set area are all within a directed acyclic graph (DAG). Without cycles a DFS which is how I implemented my recursive search will always be finite. We try to maximize the acycylic subgraph (MAS) problem by trying to define as few words as possible which means that we are also minimizing the inverse which is the Feedback Vertex Set (FVS). The answer to our original question is the minimum FVS of the graph of all words in the dictionary where a->b means a defines b. This is a brand new application of the FVS problem. Hopefully with more work on this problem that truly good applications in fields like ML can be found.

## Preliminaries

Let's explain the general setup for the graph from the dictionary. let W be a set of words and w ∈ W then def(w) = S where S is the set of unique words in the definition of the word w. LET G = (V,E) where V = (v ∈ W ∩ def(W)) and E = (eij | i ∈ def(w), j ∈ w where w ∈ W). This will exactly give you the graph G where one can find "the smallest set of words" from a minimum FVS. Let's define FVS just so we are completely clear on that, a subset S of V(G) is a directed Feedback Vertex Set (FVS), if the induced subgraph V(G) \ S is acyclic.

## Directed FVS Approximation Algorithm

Start with any directed graph G.

1. Cut any nodes from the graph with no in-degree do this repeatedly until no nodes are found
2. Pick the node with the highest out-degree and cut it from the graph and add it to the set X
3. If graph has no nodes stop else repeat #1-3

The set X is your FVS.

## Analysis

I have so far been unable to determine the size of the minimum FVS for each of the two dictionaries that I implemented. The reason why that is so is because the best solution to the directed min. FVS problem is O(1.9977^N) by Igor Razgon. That is practically O(2^N), at O(2^N) worst time complexity this gives us a worst time of O(2^110,301) or O(Infinity) according to the google calculator! I've tried using already coded FVS algorithmns found online but even after an entire 24 hours or running what seemed to be the best solution that I figured my solution was better in the long run. In contrast to the algorithm that finds the minimum FVS, my approximation algorithm only takes about 15s to run on my modest hardware. I'd say that my algorithm is really good when it comes to coming up with fast solutions. My approximation algorithm might not get perfect results but I think that around less than 7% of vertices in the FVS I think that this solution is more than optimal for its use in this application. There is the problem of the legibility of reconstructed answers, it could be that our solution also maximizes the amount of connections that aren't being made anymore in a DFS. I noticed that the solutions I got on my first dictionary included a lot of words like 'of', 'that', 'is', etc. getting rid of those words even if not neccessary from the solution would cause the reconstructed words to balloon in size because of how frequently those words are used in definitions. Maybe it is a little good that my algorithmn doesn't provide perfect solutions to the original problem. I've had experience with this because I have done something called culling the solutions my algorithm makes. it turns out not every word that gets deleted during my algorithm is necessary for a complete solution. I had a verifier verify that my solutions would give functional solutions and went about deleting unneccessary words from my solution set one-by-one. I made a solution so good that even though it verifies it causes my word redefiner to crash because of too many stack frames being opened and my expanded definition becoming freakishly large. At that point is it even worth it if the redifined definition becomes so large it's unreadable. I'd like to say that because of the way I set up the WordNet dictionary up that it is much more readable than my first dictionary's results from my algorithm.

## Website

An application of the FVS algorithm on my website where definitions of any word can be searched up and defined using only words in the FVS.

https://noeldev.site/search/search.html

## Dataset(s)

https://www.bragitoff.com/2016/03/english-dictionary-in-csv-format/

WordNet
