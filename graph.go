package main

import (
	"fmt"
)

type Graph struct {
	vertices map[string]*Vertex
}

type Vertex struct {
	key     string
	outList []*Vertex
	inList  []*Vertex
}

// adds vertex to graph with key k, will not add duplicates
func (g *Graph) AddVertex(k string) {
	if !g.contains(k) {
		vertex := &Vertex{key: k}
		g.vertices[k] = vertex
	}
}

// function which returns whether the vertex with key k is in the graph
func (g *Graph) contains(k string) bool {
	_, ok := g.vertices[k]
	return ok
}

// Adds Edge to graph going (from) --> (to) if it doesn't already exist
func (g *Graph) AddEdge(from string, to string) {
	fromVertex := g.getVertex(from)
	toVertex := g.getVertex(to)

	if !(fromVertex == nil || toVertex == nil) {
		if !containsEdge(fromVertex, to) {
			fromVertex.outList = append(fromVertex.outList, toVertex)
			toVertex.inList = append(toVertex.inList, fromVertex)
		}
	}
}

// returns whether edge exists in from's outlist
func containsEdge(from *Vertex, to string) bool {
	for _, v := range from.outList {
		if v.key == to {
			return true
		}
	}
	return false
}

// retrieves vertex from graph
func (g *Graph) getVertex(k string) *Vertex {
	val, ok := g.vertices[k]
	if ok {
		return val
	} else {
		return nil
	}
}

// Deletes Vertex from graph, Warning: Doesn't delete null ptrs in adjacency lists, call clearLists()
func (g *Graph) DeleteVertex(k string) []*Vertex {

	val, ok := g.vertices[k]
	if ok {
		outLi := val.outList
		*val = Vertex{}
		delete(g.vertices, k)

		return outLi
	}

	return []*Vertex{}

}

// Must call after DeleteVertex calls
func (g *Graph) clearLists() {
	for _, v := range g.vertices {
		v.inList = rmList(v.inList)
		v.outList = rmList(v.outList)
	}
}

// Removes null ptrs from adjacency list
func rmList(li []*Vertex) []*Vertex {
	for i, v := range li {
		if v.key == "" {
			li = remove(li, i)
			li = rmList(li)
			break
		}
	}

	return li

}

// list removal method
func remove(s []*Vertex, i int) []*Vertex {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Prints Graph
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

// Prints Vertex from Graph
func (g *Graph) PrintVert(k string) {
	for _, v := range g.vertices {
		if v.key == k {
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
}

// Prints Graph Size
func (g *Graph) PrintSize() {
	fmt.Println("\ngSize: ", len(g.vertices))
}

// Returns Graph Size
func (g *Graph) Size() int {
	return len(g.vertices)
}

// Transfers Data in Dictionary to Graph
func (g *Graph) AddData(d *Dictionary) {
	for _, v := range d.definitions {
		g.AddVertex(v.name)
		for _, word := range v.words {
			g.AddVertex(word)
		}
	}

	for _, v := range d.definitions {
		for _, word := range v.words {
			// word defines name
			g.AddEdge(word, v.name)
		}
	}
}

/* OPTIMIZE RUNTIME
IDEAS :
1. Add PQ for findHighest!
*/

func (g *Graph) top() []string {

	var freeWords []string

	for _, v := range g.vertices {
		if len(v.inList) == 0 {
			freeWords = append(freeWords, v.key)
		}
	}

	return freeWords
}

// only used for first pop
func (g *Graph) pop() (int, []*Vertex) {
	pops := 0
	var delList []*Vertex
	var li []*Vertex

	for _, v := range g.vertices {
		if len(v.inList) == 0 {
			li = g.DeleteVertex(v.key)
			delList = append(delList, li...)
			pops++
		}
	}

	g.clearLists()

	return pops, delList
}

func (g *Graph) popList(outLi []*Vertex) (int, []*Vertex) {
	pops := 0
	var delList []*Vertex
	var li []*Vertex

	for _, v := range outLi {
		if len(v.inList) == 0 {
			li = g.DeleteVertex(v.key)
			delList = append(delList, li...)
			pops++
		}

	}

	g.clearLists()

	return pops, delList
}

func (g *Graph) firstPop() {
	pops, delList := g.pop()

	for pops != 0 {
		pops, delList = g.popList(delList)
	}
}

// OK
func (g *Graph) vertCover() []string {
	g.firstPop()

	var delNodes []string

	for g.Size() != 0 {
		delNodes = append(delNodes, g.delHighest())
	}

	return delNodes
}

func (g *Graph) delHighest() string {
	vert := g.findHighest()
	key := vert.key

	delList := g.DeleteVertex(vert.key)
	g.clearLists()

	pops, delList := g.popList(delList)

	for pops != 0 {
		pops, delList = g.popList(delList)
	}

	return key
}

// BAD O(N) runtime, fix with PQ?
func (g *Graph) findHighest() *Vertex {
	var vert *Vertex
	top := 0

	for _, val := range g.vertices {
		if len(val.outList) > top {
			top = len(val.outList)
			vert = val
		}
	}

	return vert
}
