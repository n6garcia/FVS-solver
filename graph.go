package main

import (
	"container/heap"
	"fmt"
)

type Graph struct {
	vertices map[string]*Vertex
	pq       PriorityQueue
	pqIdx    map[string]*Item
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

func (g *Graph) top() []string {

	fmt.Println("finding free words...")

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
	for _, v := range delList {
		item, ok := g.pqIdx[v.key]
		if ok {
			g.pq.update(item, v, len(v.outList))
		}
	}

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
	for _, v := range delList {
		item, ok := g.pqIdx[v.key]
		if ok {
			g.pq.update(item, v, len(v.outList))
		}
	}

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

	g.clearLists()
	for _, v := range delList {
		item, ok := g.pqIdx[v.key]
		if ok {
			g.pq.update(item, v, len(v.outList))
		}
	}

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
	delete(g.pqIdx, item.value.key)

	return item.value
}

/* PQ implementation */

// An Item is something we manage in a priority queue.
type Item struct {
	value    *Vertex // The value of the item; arbitrary.
	priority int     // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Note: always CALL after clearLists
// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value *Vertex, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (g *Graph) removePQ(k string) {
	item, ok := g.pqIdx[k]

	if ok {

		heap.Remove(&g.pq, item.index)
		heap.Init(&g.pq)

		delete(g.pqIdx, k)
	}

}

func (g *Graph) initPQ() {
	fmt.Println("initializing PQ...")

	g.pq = make(PriorityQueue, len(g.vertices))
	i := 0
	for _, v := range g.vertices {
		g.pq[i] = &Item{
			value:    v,
			priority: len(v.outList),
			index:    i,
		}
		i++
	}
	heap.Init(&g.pq)

	for _, item := range g.pq {
		g.pqIdx[item.value.key] = item
	}
}
