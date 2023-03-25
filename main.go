package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var delNodes []string
var tGraph *Graph
var dict *Dictionary

func main() {

	LoadDict()

	reconstructWord("john")

	//cullSolution()

	//simulatedAnnealing()

	//graphVerify()

	//dictVerify()

	//exportJson()

	//exportCSV()

	//handleServer()

}

func LoadDict() {
	start := time.Now()

	dict = &Dictionary{definitions: make(map[string]*Definition)}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()
}

func LoadWNDict() {

}

func write(li []string, fn string) {
	json, err := json.MarshalIndent(li, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile("data/"+fn, json, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNodes(fn string) []string {
	file, err := os.Open("data/" + fn)
	if err != nil {
		fmt.Println("error loading json")
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var txt string

	for scanner.Scan() {
		line := scanner.Text()
		txt = txt + line
	}

	bytes := []byte(txt)

	var myData []string

	json.Unmarshal(bytes, &myData)

	return myData
}

func origHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	defn := dict.getDef(word)

	var str string

	for i, val := range defn {
		if i == 0 {
			str = str + val
		} else {
			str = str + " " + val
		}
	}

	w.Write([]byte(str))
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	defn := dict.expandDef(delNodes, word)

	var str string

	for i, val := range defn {
		if i == 0 {
			str = str + val
		} else {
			str = str + " " + val
		}
	}

	w.Write([]byte(str))
}

func handleServer() {
	delNodes = getNodes("bestSol.json")

	r := mux.NewRouter()

	r.HandleFunc("/orig", origHandler).Methods("GET")
	r.HandleFunc("/new", newHandler).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3001", nil))
}

func reconstructWord(word string) {
	delNodes = getNodes("delNodes.json")

	defn := dict.getDef(word)

	var str string

	for i, val := range defn {
		if i == 0 {
			str = str + val
		} else {
			str = str + " " + val
		}
	}

	fmt.Println(str)

	defn = dict.expandDef(delNodes, word)

	str = ""

	for i, val := range defn {
		if i == 0 {
			str = str + val
		} else {
			str = str + " " + val
		}
	}

	fmt.Println(str)
}

func Solve() {
	tGraph = &Graph{vertices: make(map[string]*Vertex), pqMap: make(map[string]*Item)}

	tGraph.AddData(dict)
	tGraph.pqInit()

	listFree := tGraph.top()

	write(listFree, "freeWords.json")

	start := time.Now()

	delNodes = tGraph.FVS()

	write(delNodes, "delNodes.json")

	fmt.Println("nodes removed: ", len(delNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func cullSolution() {
	tGraph = &Graph{vertices: make(map[string]*Vertex)}

	tGraph.AddData(dict)

	listFree := tGraph.top()

	delNodes = getNodes("simNodes3.json")

	start := time.Now()

	cullNodes := tGraph.cullSol(delNodes, listFree)

	write(cullNodes, "cullNodes.json")

	fmt.Println("nodes removed: ", len(cullNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func simulatedAnnealing() {
	tGraph = &Graph{vertices: make(map[string]*Vertex)}

	tGraph.AddData(dict)

	listFree := tGraph.top()

	delNodes = getNodes("simNodes2.json")

	start := time.Now()

	simNodes := tGraph.simAnneal(delNodes, listFree)

	write(simNodes, "simNodes.json")

	fmt.Println("nodes removed: ", len(simNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func graphVerify() {
	delNodes = getNodes("bestSol.json")

	tGraph = &Graph{vertices: make(map[string]*Vertex)}

	tGraph.AddData(dict)

	listFree := tGraph.top()

	start := time.Now()

	verified := tGraph.verify(delNodes, listFree)

	fmt.Println("verified: ", verified)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func dictVerify() {
	start := time.Now()

	delNodes = getNodes("delNodes.json")

	verified := dict.verify(delNodes)

	fmt.Println("verified: ", verified)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func exportJson() {
	type node struct {
		Name string `json:"name"`
	}

	type link struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}

	type expGraph struct {
		Nodes []node `json:"nodes"`
		Links []link `json:"links"`
	}

	tGraph = &Graph{vertices: make(map[string]*Vertex)}
	tGraph.AddData(dict)

	fmt.Println("exporting graph...")

	var export expGraph

	for _, vert := range tGraph.vertices {
		n := node{vert.key}

		export.Nodes = append(export.Nodes, n)

		for _, out := range vert.outList {
			var l link

			if out.key != "" {
				l = link{vert.key, out.key}

				export.Links = append(export.Links, l)
			}

		}

	}

	b, err := json.MarshalIndent(export, "", "")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile("data/expGraph.json", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exportCSV() {
	rows := [][]string{
		{"source", "target"},
	}

	tGraph = &Graph{vertices: make(map[string]*Vertex)}
	tGraph.AddData(dict)

	for _, vert := range tGraph.vertices {

		for _, out := range vert.outList {
			rows = append(rows, []string{vert.key, out.key})
		}

	}

	csvfile, err := os.Create("data/expCSV.csv")

	if err != nil {
		log.Fatalf("Failed to create file, : %s", err)
	}

	cswriter := csv.NewWriter(csvfile)

	for _, row := range rows {
		_ = cswriter.Write(row)
	}

	cswriter.Flush()
	csvfile.Close()
}
