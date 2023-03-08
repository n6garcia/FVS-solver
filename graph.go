package main

import (
	"container/heap"
	"fmt"
)

type Graph struct {
	vertices map[string]*Vertex
	pq       PriorityQueue
	pqMap    map[string]*Item
}

type Vertex struct {
	key     string
	outList []*Vertex
	inList  []*Vertex
}

/* Graph Population Functions */

// Transfers Data in Dictionary to Graph
func (g *Graph) AddData(d *Dictionary) {
	fmt.Println("adding data to graph...")

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

// adds vertex to graph with key k, will not add duplicates
func (g *Graph) AddVertex(k string) {
	if !g.containsVertex(k) {
		vertex := &Vertex{key: k}
		g.vertices[k] = vertex
	}
}

// function which returns whether the vertex with key k is in the graph
func (g *Graph) containsVertex(k string) bool {
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

// Deletes Vertex from graph, Warning: Doesn't delete null ptrs in adjacency lists
// Please always use modLen()
func (g *Graph) DeleteVertex(k string) []*Vertex {
	val, ok := g.vertices[k]

	if ok {
		//outLi := val.outList
		// ^--- shallow copy?

		outLi := make([]*Vertex, len(val.outList))
		copy(outLi, val.outList)

		*val = Vertex{}
		delete(g.vertices, k)

		g.pqRemove(k)

		return outLi
	}

	return []*Vertex{}

}

func modLen(li []*Vertex) int {
	count := 0

	for _, v := range li {
		if v.key != "" {
			count += 1
		}
	}

	return count

}

/* Print Functions */

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

/* verify Functions */

func (g *Graph) verify(delNodes []string, freeWords []string) bool {
	//fmt.Println("verifying...")

	stopWords := make(map[string]bool)
	for _, v := range g.vertices {
		stopWords[v.key] = false
	}
	for _, k := range delNodes {
		stopWords[k] = true
	}
	for _, k := range freeWords {
		stopWords[k] = true
	}

	whiteSet := make(map[string]bool)
	for k, v := range stopWords {
		if !v {
			whiteSet[k] = !v // true
		}
	}

	graySet := make(map[string]bool)
	for k := range whiteSet {
		graySet[k] = false
	}

	blackSet := make(map[string]bool)
	for k := range whiteSet {
		blackSet[k] = false
	}

	for len(whiteSet) != 0 {
		var current string

		for k := range whiteSet {
			current = k
			break
		}

		if g.dfs(current, whiteSet, graySet, blackSet, stopWords) {
			return false
		}

	}

	return true
}

func (g *Graph) dfs(current string, whiteSet map[string]bool, graySet map[string]bool, blackSet map[string]bool, stopWords map[string]bool) bool {
	// move vertex from whiteSet to graySet
	graySet[current] = true
	delete(whiteSet, current)

	vert, ok := g.vertices[current]
	if ok {
		stopBool, ok := stopWords[current]
		if ok {
			if !stopBool {
				for _, v := range vert.inList {
					neighbor := v.key

					bsBool, ok := blackSet[neighbor]
					if ok {
						if bsBool {
							continue
						}
					}

					gsBool, ok := graySet[neighbor]
					if ok {
						if gsBool {
							return true
						}
					}

					if g.dfs(neighbor, whiteSet, graySet, blackSet, stopWords) {
						return true
					}
				}
			}
		}
	}

	// move vertex from graySet to blackSet
	graySet[current] = false
	blackSet[current] = true

	return false
}

/* Brute Force Functions*/

func (g *Graph) bruteForce(listFree []string) []string {
	fmt.Println("brute forcing optimal solution...")

	undefined = listFree

	g.firstPop()

	var fullSet []string

	for k := range tGraph.vertices {
		fullSet = append(fullSet, k)
	}

	return g.allSubsets(fullSet)
}

var undefined []string
var opt []string

func (g *Graph) allSubsets(fullSet []string) []string {
	var subset []string
	g.subsetHelper(fullSet, subset, 0)

	return opt
}

func (g *Graph) subsetHelper(fullSet []string, subset []string, i int) {
	if i == len(g.vertices) || len(subset) == 1 {
		fmt.Println(subset)
		if len(opt) == 0 {
			if g.verify(subset, undefined) {
				opt = subset
			}
		} else if len(subset) < len(opt) {
			if g.verify(subset, undefined) {
				opt = subset
			}
		}
	} else {
		g.subsetHelper(fullSet, subset, i+1)

		subset = append(subset, fullSet[i])
		g.subsetHelper(fullSet, subset, i+1)
	}
}

/* Mod Cover Functions */

func (g *Graph) modCover() []string {
	fmt.Println("performing mod cover...")

	g.firstPop()

	var delNodes []string

	for g.Size() != 0 {
		delNodes = append(delNodes, g.delHighest())
	}

	return delNodes
}

func (g *Graph) top() []string {
	fmt.Println("finding free words...")

	var freeWords []string

	for _, v := range g.vertices {
		if modLen(v.inList) == 0 {
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

	g.pqShuffle()
	g.pqUpdateList(delList)

	return pops, delList
}

func (g *Graph) popList(outLi []*Vertex) (int, []*Vertex) {
	pops := 0
	var delList []*Vertex
	var li []*Vertex

	for _, v := range outLi {
		if modLen(v.inList) == 0 {
			li = g.DeleteVertex(v.key)
			delList = append(delList, li...)
			pops++
		}

	}

	g.pqShuffle()
	g.pqUpdateList(delList)

	return pops, delList
}

func (g *Graph) firstPop() {
	pops, delList := g.pop()

	for pops != 0 {
		pops, delList = g.popList(delList)
	}
}

func (g *Graph) delHighest() string {
	vert := g.findHighest()
	key := vert.key

	delList := g.DeleteVertex(vert.key)

	g.pqShuffle()
	g.pqUpdateList(delList)

	pops, delList := g.popList(delList)

	for pops != 0 {
		pops, delList = g.popList(delList)
	}

	return key
}

func (g *Graph) findHighest() *Vertex {
	item := heap.Pop(&g.pq).(*Item)
	delete(g.pqMap, item.value.key)

	return item.value
}

/* Priority Queue Functions */

func (g *Graph) pqInit() {
	fmt.Println("initializing PQ...")

	g.pq = make(PriorityQueue, len(g.vertices))
	i := 0
	for _, v := range g.vertices {
		g.pq[i] = &Item{
			value:    v,
			priority: modLen(v.outList),
			index:    i,
		}
		i++
	}
	heap.Init(&g.pq)

	for _, item := range g.pq {
		g.pqMap[item.value.key] = item
	}
}

func (g *Graph) pqRemove(k string) {
	item, ok := g.pqMap[k]

	if ok {

		heap.Remove(&g.pq, item.index)
		delete(g.pqMap, k)

	}

}

// call after all removePQ calls are done to reinitialize PQ ordering!
func (g *Graph) pqShuffle() {
	heap.Init(&g.pq)
}

// Must Call after ClearLists if using PQ!
func (g *Graph) pqUpdateList(delList []*Vertex) {
	for _, v := range delList {
		item, ok := g.pqMap[v.key]
		if ok {
			g.pq.update(item, v, modLen(v.outList))
		}
	}
}
