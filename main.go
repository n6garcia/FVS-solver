package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// 2nd Draft

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

	dict = &Dictionary{definitions: make(map[string]*Definition)}

	for ch := 'A'; ch <= 'Z'; ch++ {
		dict.loadData(string(ch) + ".json")
	}

	dict.PrintSize()

	/*
		tGraph = &Graph{vertices: make(map[string]*Vertex)}

		tGraph.AddData(dict)

		tGraph.PrintSize()

		listFree := tGraph.top()

		write(listFree, "freeWords.json")

		fmt.Println("\nlistFree: ", len(listFree))

		delNodes = tGraph.vertCover()

		write(delNodes, "delNodes.json")

		fmt.Println("nodes removed: ", len(delNodes))
	*/

	//delNodes = getNodes()

	//verified := dict.verify(delNodes)

	//log.Println(verified)

	//handleServer()
}
