package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/slices"
)

type Graph struct {
	vertices map[string]*Vertex
}

type Vertex struct {
	key     string
	outList []*Vertex
	inList  []*Vertex
}

func (g *Graph) AddVertex(k string) {
	if !g.contains(k) {
		vertex := &Vertex{key: k}
		g.vertices[k] = vertex
	}
}

func (g *Graph) contains(k string) bool {
	_, ok := g.vertices[k]
	return ok
}

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

func containsEdge(from *Vertex, to string) bool {
	for _, v := range from.outList {
		if v.key == to {
			return true
		}
	}
	return false
}

func (g *Graph) getVertex(k string) *Vertex {
	val, ok := g.vertices[k]
	if ok {
		return val
	} else {
		return nil
	}
}

func (g *Graph) DeleteVertex(k string) {

	val, ok := g.vertices[k]
	if ok {
		*val = Vertex{}
		delete(g.vertices, k)
	}

}

func (g *Graph) clearLists() {
	for _, v := range g.vertices {
		v.inList = rmList(v.inList)
		v.outList = rmList(v.outList)
	}
}

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

func (g *Graph) PrintSize() {
	fmt.Println("\ngSize: ", len(g.vertices))
}

func (g *Graph) Size() int {
	return len(g.vertices)
}

type Dictionary struct {
	definitions []*Definition
}

type Definition struct {
	name  string
	words []string
}

func (d *Dictionary) addDef(n string, w []string) {
	d.definitions = append(d.definitions, &Definition{name: n, words: w})
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

	fmt.Println("\nfile: ", fn)

	scanner := bufio.NewScanner(file)

	var txt string

	for scanner.Scan() {
		line := scanner.Text()
		txt = txt + line
	}

	bytes := []byte(txt)

	fmt.Println("\nisValid: ", json.Valid(bytes))

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
	for _, v := range d.definitions {
		g.AddVertex(v.name)
		for _, word := range v.words {
			g.AddVertex(word)
		}
	}

	for _, v := range d.definitions {
		for _, word := range v.words {
			g.AddEdge(v.name, word)
		}
	}
}

func (g *Graph) top() []string {

	var freeWords []string

	for _, v := range g.vertices {
		if len(v.inList) == 0 {
			freeWords = append(freeWords, v.key)
		}
	}

	return freeWords
}

func (g *Graph) pop() int {
	pops := 0
	for _, v := range g.vertices {
		if len(v.inList) == 0 {
			g.DeleteVertex(v.key)
			pops++
		}
	}

	g.clearLists()

	return pops
}

func (g *Graph) firstPop() {
	pops := g.pop()

	for pops != 0 {
		pops = g.pop()
	}
}

func write(li []string, fn string) {
	json, err := json.MarshalIndent(li, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile(fn, json, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

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

	g.DeleteVertex(vert.key)
	g.clearLists()

	pops := g.pop()

	for pops != 0 {
		pops = g.pop()
	}

	return key
}

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

func (d *Dictionary) expandDef(delNodes []string, k string) []string {
	// create map
	wordMap := make(map[string]bool)

	// init map with false vals
	for _, val := range d.definitions {
		wordMap[val.name] = false
	}

	// init map with vertex Cover (true vals)
	// WARNING: possibly adds values not in dict.definitions
	for _, val := range delNodes {
		wordMap[val] = true
	}

	var defn []string

	// find definition of testName
	for _, val := range d.definitions {
		if val.name == k {
			defn = val.words
			break
		}
	}

	fmt.Println(defn)

	for i, val := range defn {
		defn = slices.Insert(defn, i, d.recurDef(wordMap, val)...)
	}

	return defn
}

func (d *Dictionary) recurDef(wordMap map[string]bool, k string) []string {
	val, ok := wordMap[k]
	if !ok || val {
		return []string{}
	} else {
		defn := d.findDef(k)
		for i, val := range defn {
			defn = slices.Insert(defn, i, d.recurDef(wordMap, val)...)
		}
		return defn
	}
}

func (d *Dictionary) findDef(k string) []string {
	for _, val := range d.definitions {
		if val.name == k {
			return val.words
		}
	}
	return []string{}
}

func main() {

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict := &Dictionary{}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	dict.PrintSize()

	tGraph.AddData(dict)

	tGraph.PrintSize()

	listFree := tGraph.top()

	write(listFree, "freeWords.json")

	fmt.Println("\nlistFree: ", len(listFree))

	delNodes := tGraph.vertCover()

	write(delNodes, "delNodes.json")

	fmt.Println("nodes removed: ", len(delNodes))

	fmt.Println(dict.expandDef(delNodes, "upstart"))

}
