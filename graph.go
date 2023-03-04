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

func modLen(li []*Vertex) int {
	count := 0

	for _, v := range li {
		if v.key != "" {
			count += 1
		}
	}

	return count

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
		//outLi := val.outList
		// ^--- shallow copy?

		outLi := make([]*Vertex, len(val.outList))
		copy(outLi, val.outList)

		*val = Vertex{}
		delete(g.vertices, k)

		g.removePQ(k)

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

func (g *Graph) verify(delNodes []string, freeWords []string) bool {
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

	for k, b := range stopWords {
		if !b {
			g.dfs(k, stopWords)
		}
	}

	return true
}

func (g *Graph) dfs(k string, stopWords map[string]bool) {
	v, ok := g.vertices[k]
	if ok {
		b, ok := stopWords[k]
		if ok {
			if !b {
				for _, u := range v.inList {
					g.dfs(u.key, stopWords)
				}
			}
		}
	}
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

	g.pqUpdateList(delList)

	return pops, delList
}

func (g *Graph) firstPop() {
	pops, delList := g.pop()

	for pops != 0 {
		pops, delList = g.popList(delList)
	}
}

func (g *Graph) modCover() []string {
	fmt.Println("performing mod cover...")

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

	g.pqUpdateList(delList)

	pops, delList := g.popList(delList)

	for pops != 0 {
		pops, delList = g.popList(delList)
	}

	return key
}

// BAD O(N^N) runtime, fix with PQ?
// note: need PQ with delete function
// would runtime be any better?
// (constant reshuffling and updating)
/*
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
}*/

func (g *Graph) findHighest() *Vertex {
	item := heap.Pop(&g.pq).(*Item)
	delete(g.pqMap, item.value.key)

	return item.value
}

func (g *Graph) initPQ() {
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

func (g *Graph) removePQ(k string) {
	item, ok := g.pqMap[k]

	if ok {

		heap.Remove(&g.pq, item.index)
		heap.Init(&g.pq)

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
