# Dictionary-Problem

## The Problem
What is the smallest possible set of words that can be used to define every other word in the dictionary?

## The Solution
My current solution to the problem is to create a directed graph from an entire dictionary dataset. One word points to another if that word is in the other words definition. Note: if a word has no in-degree it is considered defined and is thus deleted from the graph. I plan on coding a greedy solution to this problem which involves deleting the words with the highest out-degree first iteratively until the entire graph is gone.

## Plans
I plan on hosting an application of the minimal set of words on my website where definitions of any word can be searched up and defined using only words in the minimal set.

