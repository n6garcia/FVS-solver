package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type WNdict struct {
	IDMappings  map[string]*WNdef
	definitions map[string][]*WNdef
}

type WNdef struct {
	name       string
	origDef    string
	regexDef   string
	regexWords []string
	mappings   []string
}

func (wn *WNdict) Print() {
	for _, v := range wn.definitions {
		for _, d := range v {
			fmt.Println("name: ", d.name)
			fmt.Println("origDef: ", d.origDef)
		}
	}
}

func (wn *WNdict) PrintSize() {
	fmt.Println("\nsize : ", len(wn.definitions))
}

func (wn *WNdict) addDef(ID string, def *WNdef) {
	wn.IDMappings[ID] = def
	wn.definitions[def.name] = append(wn.definitions[def.name], def)
}

func (wn *WNdict) loadData(fn string) {
	bytes, err := os.ReadFile("wrangle/wordnet/" + fn) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("isValid: ", json.Valid(bytes))

	var myData map[string][]interface{}

	json.Unmarshal(bytes, &myData)

	for k, v := range myData {
		ID := k
		name := v[0].(string)
		origDef := v[1].(string)
		regexDef := v[2].(string)

		regexWordsInterface := v[3].([]interface{})
		var regexWords []string
		for _, word := range regexWordsInterface {
			regexWords = append(regexWords, word.(string))
		}

		mappingsInterface := v[4].([]interface{})
		var mappings []string
		for _, word := range mappingsInterface {
			mappings = append(mappings, word.(string))
		}

		def := &WNdef{name: name, origDef: origDef, regexDef: regexDef, regexWords: regexWords, mappings: mappings}

		wn.addDef(ID, def)
	}
}

func (wn *WNdict) AddData(g *Graph) {
	fmt.Println("adding data to graph...")

	// add words
	for _, li := range wn.definitions {
		for _, v := range li {
			g.AddVertex(v.name)
			for _, word := range v.regexWords {
				g.AddVertex(word)
			}
		}
	}

	// add edges (has to happen once all words are in graph!)
	for _, li := range wn.definitions {
		for _, v := range li {
			for _, word := range v.regexWords {
				// word defines name
				g.AddEdge(word, v.name)
			}
		}
	}
}

func (wn *WNdict) expandDef(delNodes []string, k string) string {
	return ""
}

func (wn *WNdict) recursiveSearch(wordMap map[string]bool, k string) []string {
	return []string{}
}

// for use in recursiveSearch
func (wn *WNdict) findDef(k string) []string {
	return []string{}
}

// for use in main()
func (wn *WNdict) getDef(k string) string {
	return ""
}

// very slow implementation!
// implementation takes about 1hr40m on my computer to run (average hardware)
func (wn *WNdict) verify(delNodes []string) bool {

	fmt.Println("verifying...")

	return false

}
