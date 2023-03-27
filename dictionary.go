package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type dictInterface interface {
	Print()
	PrintSize()
	loadData(string)
	AddData(*Graph)
	getDef(string) string
	expandDef([]string, string) string
	verify([]string) bool
}

type Dictionary struct {
	definitions map[string]*Definition
	//definitions []*Definition
	// ^--- old DS
}

type Definition struct {
	name  string
	words []string
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

func (d *Dictionary) addDef(n string, w []string) {
	defn := &Definition{name: n, words: w}
	d.definitions[n] = defn
}

func (d *Dictionary) loadData(fn string) {
	bytes, err := os.ReadFile("wrangle/cleaned/" + fn) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("isValid: ", json.Valid(bytes))

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

// Transfers Data in Dictionary to Graph
func (d *Dictionary) AddData(g *Graph) {
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

func (d *Dictionary) expandDef(delNodes []string, k string) string {
	wordMap := make(map[string]bool)
	var defn []string

	if k == "" {
		return ""
	}

	for _, val := range d.definitions {
		wordMap[val.name] = false
		for _, word := range val.words {
			wordMap[word] = false
		}
	}

	for _, val := range delNodes {
		wordMap[val] = true
	}

	defn = d.findDef(k)

	var newDefn []string = []string{}
	for _, val := range defn {
		expand := d.recursiveSearch(wordMap, val)
		if len(expand) != 0 {
			newDefn = append(newDefn, expand...)
		} else {
			newDefn = append(newDefn, val)
		}
	}

	var str string

	for i, val := range newDefn {
		if i == 0 {
			str = str + val
		} else {
			str = str + " " + val
		}
	}

	return str
}

func (d *Dictionary) recursiveSearch(wordMap map[string]bool, k string) []string {
	val, ok := wordMap[k]
	if !ok || val {
		return []string{}
	} else {
		defn := d.findDef(k)
		var newDefn []string = []string{}
		for _, val := range defn {
			expand := d.recursiveSearch(wordMap, val)
			if len(expand) != 0 {
				newDefn = append(newDefn, expand...)
			} else {
				newDefn = append(newDefn, val)
			}
		}
		return newDefn
	}
}

// for use in recursiveSearch
func (d *Dictionary) findDef(k string) []string {
	defn, ok := d.definitions[k]
	if ok {
		return defn.words
	} else {
		return []string{}
	}
}

// for use in main()
func (d *Dictionary) getDef(k string) string {
	if k == "" {
		return ""
	}
	defn, ok := d.definitions[k]
	if ok {

		var str string

		for i, val := range defn.words {
			if i == 0 {
				str = str + val
			} else {
				str = str + " " + val
			}
		}

		return str
	} else {
		return ""
	}
}

// very slow implementation!
// implementation takes about 1hr40m on my computer to run (average hardware)
func (d *Dictionary) verify(delNodes []string) bool {

	fmt.Println("verifying...")

	for _, val := range d.definitions {
		d.expandDef(delNodes, val.name)
	}

	return true

}
