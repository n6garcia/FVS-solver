package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
		err = os.WriteFile(fn, json, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")

	defn := dict.expandDef(delNodes, word)

	var str string

	for _, val := range defn {
		str = str + val
	}

	w.Write([]byte(str))
}

func handleServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", getHandler).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3002", nil))
}

func main() {

	tGraph = &Graph{vertices: make(map[string]*Vertex)}

	dict = &Dictionary{}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	dict.PrintSize()

	tGraph.AddData(dict)

	tGraph.PrintSize()

	listFree := tGraph.top()

	write(listFree, "freeWords.json")

	fmt.Println("\nlistFree: ", len(listFree))

	delNodes = tGraph.vertCover()

	write(delNodes, "delNodes.json")

	fmt.Println("nodes removed: ", len(delNodes))

	handleServer()
}
