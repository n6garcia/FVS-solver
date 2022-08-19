# Dictionary-Problem

### What is the smallest set of words that can be used to define every word in the dictionary?

## The Solution
My solution to the problem is a little complex but it involves creating a directed graph from the dictionary (including the words in the definition) with words as keys. By searching for a modified vertex cover on our modified digraph using the algorithmn we have we get an approximation for the smallest word set. From all the research I've done on the problem it seems to that finding the perfect solution is an NP hard problem so an approximation algorithmn seems to be the best I can do. 

## The Modified Directed Graph and Modified Vertex Cover

### a ---> b
### a defines b

We add edges to the graph we construct where, a points to b if that a is in b's definition. I have added some special functionality to the digraph to make it possible to solve the problem of the smallest word set. Imagine the case where after a cover a word has no in-degree, what this means is that either it has no definition or we've "defined" every word in that words definition by deleting them from the graph. In either case, this word is essentially defined. What we hope to find with our modified vertex cover is a cover in our graph where after deleting the vertices from the graph in the cover the rest of the graph is automatically deleted because they end up all end up being "defined". To simulate this action of a word becoming defined, we "pop" all the vertices from the graph which have no in-degree until none can be found. This should happen every time a vertex is deleted from the graph.

## The Algorithmn

I thought of many algorithmns, some of which involved SCC's, but I went with a simple greedy solution. My solution to finding a modified vertex cover is to delete the highest out-degree vertex one at a time until every vertex in the graph is gone, either by being deleted or popped. By saving every vertex me manually delete we can get a modified vertex cover that can be used to "define" every other word in the graph.

## Findings
Using the algorithm described above I found a set of 16,214 words that can be used to define every word in the 110,447 word dictionary that I used for the test. The words can be found in the vertexCover.json file in the data folder.

### Interesting Words
They, An, In, Let, To, It, I, Cat, God, Bible, Angel...

## Website
An application of the minimal set of words on my website where definitions of any word can be searched up and defined using only words in the minimal set.

### Link
https://noeldev.site/search/search.html

## Dataset
https://www.bragitoff.com/2016/03/english-dictionary-in-csv-format/

