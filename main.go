package main

import (
	"bufio"
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

func getNodes() []string {
	file, err := os.Open("data/delNodes.json")
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
	r := mux.NewRouter()

	r.HandleFunc("/orig", origHandler).Methods("GET")
	r.HandleFunc("/new", newHandler).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3001", nil))
}

func main() {

	/* set-up dictionary */

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

	/* set-up graph and solve */

	/*
		tGraph = &Graph{vertices: make(map[string]*Vertex), pqMap: make(map[string]*Item)}

		tGraph.AddData(dict)
		tGraph.pqInit()

		listFree := tGraph.top()

		write(listFree, "freeWords.json")

		start = time.Now()

		delNodes = tGraph.modCover()

		write(delNodes, "delNodes.json")

		fmt.Println("nodes removed: ", len(delNodes))

		t = time.Now()
		elapsed = t.Sub(start)
		fmt.Println("\ntime elapsed : ", elapsed)
	*/

	/* verify sol. (graph) */

	delNodes = getNodes()

	tGraph = &Graph{vertices: make(map[string]*Vertex)}

	tGraph.AddData(dict)

	listFree := tGraph.top()

	start = time.Now()

	verified := tGraph.verify2(delNodes, listFree)

	fmt.Println("verified: ", verified)

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println("\ntime elapsed : ", elapsed)

	/* verify sol. (dictionary) */

	/*
		start = time.Now()

		delNodes = getNodes()

		verified := dict.verify(delNodes)

		fmt.Println("verified: ", verified)

		t = time.Now()
		elapsed = t.Sub(start)
		fmt.Println("\ntime elapsed : ", elapsed)
	*/

	/* Export Graph Json*/

	/*
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

		var export expGraph

		for _, vert := range tGraph.vertices {
			n := node{vert.key}

			export.Nodes = append(export.Nodes, n)

			for _, out := range vert.outList {
				l := link{vert.key, out.key}

				export.Links = append(export.Links, l)
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
	*/

	/* Export Graph CSV*/

	/*
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
	*/

	/* handle online service */

	//handleServer()
}
