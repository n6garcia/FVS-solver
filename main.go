package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Graph struct {
	vertices []*Vertex
	vertsAdd []string
	edgeAdd  []string
}

type Vertex struct {
	key     string
	outList []*Vertex
	inList  []*Vertex
}

func (g *Graph) AddVertex(k string) {
	if !contains(g.vertices, k) {
		g.vertices = append(g.vertices, &Vertex{key: k})
	}
}

func contains(s []*Vertex, k string) bool {
	for _, v := range s {
		if k == v.key {
			return true
		}
	}
	return false
}

func (g *Graph) AddEdge(from string, to string) {
	fromVertex := g.getVertex(from)
	toVertex := g.getVertex(to)

	if fromVertex == nil || toVertex == nil {
		fmt.Println("cannot add edge")
	} else {
		if !containsEdge(fromVertex, to) {
			fromVertex.outList = append(fromVertex.outList, toVertex)
			toVertex.inList = append(toVertex.inList, fromVertex)
		}
	}
}

func containsEdge(from *Vertex, to string) bool {
	for _, v := range from.outList {
		if v.key == to {
			return true
		}
	}
	return false
}

func (g *Graph) getVertex(k string) *Vertex {
	for i, v := range g.vertices {
		if v.key == k {
			return g.vertices[i]
		}
	}
	return nil
}

func (g *Graph) DeleteVertex(k string, delIndex int) {

	// delete vertex
	g.vertices = remove(g.vertices, delIndex)

	// delete edges
	for _, v := range g.vertices {
		for i, out := range v.outList {
			if out.key == k {
				v.outList = remove(v.outList, i)
			}
		}
		for i, in := range v.inList {
			if in.key == k {
				v.inList = remove(v.inList, i)
			}
		}
	}

}

func remove(s []*Vertex, i int) []*Vertex {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (g *Graph) Print() {
	for _, v := range g.vertices {
		fmt.Printf("\nVertex: %s", v.key)
		fmt.Printf(" outEdges: ")
		for _, v := range v.outList {
			fmt.Printf(" %s ", v.key)
		}
		fmt.Printf(" inEdges: ")
		for _, v := range v.inList {
			fmt.Printf(" %s ", v.key)
		}
	}
}

func (g *Graph) PrintSize() {
	fmt.Println("\ngSize: ", len(g.vertices))
}

type Dictionary struct {
	definitions []*Definiton
}

type Definiton struct {
	name  string
	words []string
}

func (d *Dictionary) addDef(n string, w []string) {
	d.definitions = append(d.definitions, &Definiton{name: n, words: w})
}

func (d *Dictionary) Print() {
	for _, v := range d.definitions {
		fmt.Println("name: ", v.name)
		fmt.Println("words: ", v.words)
	}
}

func (d *Dictionary) PrintSize() {
	fmt.Println("\nsize : ", len(d.definitions))
}

func (d *Dictionary) loadData(fn string) {
	file, err := os.Open("wrangle/cleaned/" + fn)
	if err != nil {
		fmt.Println("error loading json")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var txt string

	for scanner.Scan() {
		line := scanner.Text()
		txt = txt + line
	}

	bytes := []byte(txt)

	fmt.Println("\nisValid: ", json.Valid(bytes))
	fmt.Println("\nfile: ", fn)

	var myData map[string][]interface{}

	json.Unmarshal(bytes, &myData)

	for k, v := range myData {
		var words []string
		for _, u := range v {
			words = append(words, u.(string))
		}
		d.addDef(k, words)
	}
}

func (g *Graph) AddData(d *Dictionary) {
	// add vertices
	for _, v := range d.definitions {
		g.AddVertex(v.name)
		for _, word := range v.words {
			g.AddVertex(word)
		}
	}

	// add edges
	for _, v := range d.definitions {
		for _, word := range v.words {
			g.AddEdge(v.name, word)
		}
	}
}

func (g *Graph) FirstPop() []string {

	var freeWords []string

	for _, v := range g.vertices {
		if len(v.inList) == 0 {
			freeWords = append(freeWords, v.key)
			//g.DeleteVertex(v.key, i)
		}
	}

	return freeWords
}

func (g *Graph) Pop() {
	for i, v := range g.vertices {
		if len(v.inList) == 0 {
			g.DeleteVertex(v.key, i)
		}
	}
}

func main() {

	tGraph := &Graph{}

	dict := &Dictionary{}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	//dict.loadData("A.json")

	dict.PrintSize()

	tGraph.AddData(dict)

	tGraph.PrintSize()

	/*
		listFree := tGraph.FirstPop()

		fmt.Println(listFree)
		fmt.Println(len(listFree))
	*/

	//tGraph.Print()

}
