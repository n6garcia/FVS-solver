package main

import (
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {
	fmt.Println("nothing!")
}

func BenchmarkVerify(b *testing.B) {

	var delNodes []string
	var tGraph *Graph
	var dict *Dictionary

	dict = &Dictionary{definitions: make(map[string]*Definition)}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	dict.PrintSize()

	tGraph = &Graph{vertices: make(map[string]*Vertex), pqMap: make(map[string]*Item)}

	tGraph.AddData(dict)
	tGraph.initPQ()

	listFree := tGraph.top()

	write(listFree, "freeWords.json")

	delNodes = tGraph.modCover()

	write(delNodes, "delNodes.json")

	fmt.Println("nodes removed: ", len(delNodes))

}
