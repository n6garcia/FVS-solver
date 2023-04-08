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

func LoadDict() dictInterface {
	start := time.Now()

	fmt.Println("loading dictionary...")

	dict := &Dictionary{definitions: make(map[string]*Definition)}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()

	return dict
}

func LoadWNDict() dictInterface {
	start := time.Now()

	fmt.Println("loading dictionary...")

	dict := &WNdict{definitions: make(map[string][]*WNdef), IDMappings: make(map[string]*WNdef)}

	dict.loadData("wn.json")

	dict.PrintSize()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	fmt.Println()

	return dict
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

func getNodes(fn string) []string {
	file, err := os.Open(fn)
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

func Solve(dict dictInterface) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex), pqMap: make(map[string]*Item)}

	dict.AddData(tGraph)
	tGraph.pqInit()

	listFree := tGraph.top()

	write(listFree, folder+"undefWords.json")

	start := time.Now()

	delNodes := tGraph.FVS()

	write(delNodes, folder+"delNodes.json")

	fmt.Println("nodes removed: ", len(delNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func reconstructWord(dict dictInterface, word string, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	defn := dict.getDef(word)

	fmt.Println(defn)

	defn = dict.expandDef(delNodes, word)

	fmt.Println(defn)
}

func exportSol(dict dictInterface, fn string, fn2 string) {
	start := time.Now()

	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	m := dict.export(delNodes)

	b, err := json.MarshalIndent(m, "", "")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile("data/sol/"+fn2, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func cullSolution(dict dictInterface, fn string) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	listFree := tGraph.top()

	delNodes := getNodes(folder + fn)

	start := time.Now()

	cullNodes := tGraph.cullSol(delNodes, listFree)

	write(cullNodes, folder+"cullNodes.json")

	fmt.Println("nodes removed: ", len(cullNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func simulatedAnnealing(dict dictInterface, fn string) {
	folder := dict.getFolder()

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	listFree := tGraph.top()

	delNodes := getNodes(folder + fn)

	start := time.Now()

	simNodes := tGraph.simAnneal(delNodes, listFree)

	write(simNodes, folder+"simNodes.json")

	fmt.Println("nodes removed: ", len(simNodes))

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func graphVerify(dict dictInterface, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	listFree := tGraph.top()

	start := time.Now()

	verified := tGraph.verify(delNodes, listFree)

	fmt.Println("verified: ", verified)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func alternateVerify(dict dictInterface, fn string) {
	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	tGraph := &Graph{vertices: make(map[string]*Vertex)}

	dict.AddData(tGraph)

	start := time.Now()

	tGraph.firstPop()

	for _, k := range delNodes {
		delList := tGraph.DeleteVertex(k)

		pops, delList := tGraph.popList(delList)

		for pops != 0 {
			pops, delList = tGraph.popList(delList)
		}
	}

	fmt.Println(tGraph.Size())

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func dictVerify(dict dictInterface, fn string) {
	start := time.Now()

	folder := dict.getFolder()

	delNodes := getNodes(folder + fn)

	verified := dict.verify(delNodes)

	fmt.Println("verified: ", verified)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)
}

func origHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	val, ok := SOL[word]
	if ok {
		w.Write([]byte(val[0]))
	} else {
		w.Write([]byte(""))
	}

}

func newHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	val, ok := SOL[word]
	if ok {
		w.Write([]byte(val[1]))
	} else {
		w.Write([]byte(""))
	}
}

var SOL map[string][]string

func handleServer(fn string) {
	fmt.Println("starting server...")

	bytes, err := os.ReadFile("data/sol/" + fn)
	if err != nil {
		fmt.Print(err)
	}

	var myData map[string][]interface{}

	json.Unmarshal(bytes, &myData)

	SOL = make(map[string][]string)

	for k, v := range myData {
		for _, u := range v {
			SOL[k] = append(SOL[k], u.(string))
		}
	}

	r := mux.NewRouter()

	r.HandleFunc("/orig", origHandler).Methods("GET")
	r.HandleFunc("/new", newHandler).Methods("GET")

	http.Handle("/", r)

	fmt.Println("server ready!")

	log.Fatal(http.ListenAndServe(":3001", nil))
}

// func exportTree()

func exportNames(dict dictInterface) {
	folder := dict.getFolder()

	export := dict.getNames()

	b, err := json.MarshalIndent(export, "", "")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err = os.WriteFile(folder+"names.json", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exportJson(dict dictInterface) {
	folder := dict.getFolder()

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

	tGraph := &Graph{vertices: make(map[string]*Vertex)}
	dict.AddData(tGraph)

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
		err = os.WriteFile(folder+"expJson.json", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exportCSV(dict dictInterface, fn string) {
	folder := dict.getFolder()

	rows := [][]string{
		{"source", "target"},
	}

	tGraph := &Graph{vertices: make(map[string]*Vertex)}
	dict.AddData(tGraph)

	delNodes := getNodes(folder + fn)
	for _, k := range delNodes {
		tGraph.DeleteVertex(k)
	}

	for _, vert := range tGraph.vertices {

		for _, out := range vert.outList {
			if out.key != "" {
				rows = append(rows, []string{vert.key, out.key})
			}
		}

	}

	csvfile, err := os.Create(folder + "expCSV.csv")

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
