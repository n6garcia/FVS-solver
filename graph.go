package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
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

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value *Vertex, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
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

// only should be used if encountering errors with PQ after pqUpdateList
func (g *Graph) pqReshuffle() {
	heap.Init(&g.pq)
}

func (g *Graph) pqUpdateList(delList []*Vertex) {
	for _, v := range delList {
		item, ok := g.pqMap[v.key]
		if ok {
			g.pq.update(item, v, modLen(v.outList))
		}
	}
}

/* Graph Population Functions */

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
// Please always use modLen() to not count nil values in list
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

/* FVS Functions */

func (g *Graph) FVS() []string {
	fmt.Println("searching for FVS...")

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

func (g *Graph) findHighest() *Vertex {
	item := heap.Pop(&g.pq).(*Item)
	delete(g.pqMap, item.value.key)

	return item.value
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
			whiteSet[k] = true
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

					stopBool, ok := stopWords[neighbor]
					if ok {
						if !stopBool {
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
		}
	}

	// move vertex from graySet to blackSet
	delete(graySet, current)
	blackSet[current] = true

	return false
}

/* Cull Functions */

func (g *Graph) cullSol(delNodes []string, listFree []string) []string {
	fmt.Println("culling solution...")
	count := 0
	i := 0

	length := len(delNodes)

	for count != length {

		b, s := g.cullHelper(delNodes, listFree, i)
		if b {
			delNodes = s
		} else {
			i += 1
		}
		count += 1

	}

	return delNodes

}

func (g *Graph) cullHelper(delNodes []string, listFree []string, i int) (bool, []string) {
	var subset []string = make([]string, len(delNodes))
	copy(subset, delNodes)
	subset = RemoveIndex(subset, i)

	if g.verify(subset, listFree) {
		return true, subset
	} else {
		return false, delNodes
	}
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

/* Simulated Annealing Functions */

func (g *Graph) simAnneal(initial []string, listFree []string) []string {
	fmt.Println("simulating annealing...")

	var T0 float64 = 5
	var T float64 = T0
	var t float64 = 0
	var remCutoff = 5

	// current <-- problem.INIITAL
	var current []string = make([]string, len(initial))
	copy(current, initial)

	var currMap map[string]bool = make(map[string]bool)
	for _, v := range g.vertices {
		currMap[v.key] = false
	}
	for _, k := range current {
		currMap[k] = true
	}

	var compliment []string
	for k, v := range currMap {
		if !v {
			compliment = append(compliment, k)
		}
	}

	// for t = 1 to inf do
	for {
		t += 1

		// T <-- schedule(t)
		T = T0 - (t * 0.0001) // should run in about 100m!

		// if T = 0 then return current
		if T == 0 {
			return current
		}

		// next <-- a randomly selected successor of current

		var next []string
		var nextKey string
		var rem int

		for {

			next = make([]string, len(current))
			copy(next, current)

			rem = rand.Intn(remCutoff + 1)

			if rem < remCutoff {
				nextIdx := rand.Intn(len(next))
				nextKey = next[nextIdx]
				next = RemoveIndex(next, nextIdx)
				if g.verify(next, listFree) {
					break
				}

			} else {
				compIdx := rand.Intn(len(compliment))
				nextKey = compliment[compIdx]
				next = append(next, nextKey)
				break
			}
		}

		// △E <-- VALUE(current) - VALUE(next)

		var E float64 = float64(len(current) - len(next))
		//fmt.Println("E:", E, " curr:", len(current), " next:", len(next))

		// if △E > 0 then current <-- next
		if E > 0 {
			current = next
			if rem < remCutoff {
				compliment = append(compliment, nextKey)
			} else {
				for i, k := range compliment {
					if k == nextKey {
						compliment = RemoveIndex(compliment, i)
						break
					}
				}
			}

			// else current <-- next only with prob. e^(-△E/T)
		} else {
			prob := math.Exp(E / T)
			//fmt.Println(E, " ", T, " ", prob)
			sample := rand.Float64()

			if sample <= prob {
				current = next
				if rem < remCutoff {
					compliment = append(compliment, nextKey)
				} else {
					for i, k := range compliment {
						if k == nextKey {
							compliment = RemoveIndex(compliment, i)
							break
						}
					}
				}
			}

		}

	}

}
