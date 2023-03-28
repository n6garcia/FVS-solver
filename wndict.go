package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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
				if word != v.name {
					g.AddEdge(word, v.name)
				}
			}
		}
	}
}

func (wn *WNdict) expandDef(delNodes []string, k string) string {
	wordMap := make(map[string]bool)

	if k == "" {
		return ""
	}

	for _, val := range wn.IDMappings {
		wordMap[val.name] = false
		for _, word := range val.regexWords {
			wordMap[word] = false
		}
	}

	for _, val := range delNodes {
		wordMap[val] = true
	}

	defnArr := wn.findDefArr(k)

	var out string = ""

	for idx, defn := range defnArr {
		var str string = defn.regexDef
		for i, val := range defn.regexWords {
			if k == val {
				str = strings.Replace(str, "%s", val, 1)
				continue
			}
			expand := wn.recursiveSearch(wordMap, defn.mappings[i], val)
			if len(expand) != 0 {
				str = strings.Replace(str, "%s", expand, 1)
			} else {
				str = strings.Replace(str, "%s", val, 1)
			}
		}
		out = out + strconv.Itoa(idx+1) + ". " + str + "\n"
	}

	return out
}

func (wn *WNdict) recursiveSearch(wordMap map[string]bool, ID string, k string) string {
	val, ok := wordMap[k]
	if !ok || val {
		return ""
	} else {
		defn := wn.findDef(ID)
		var str string = defn.regexDef
		for i, val := range defn.regexWords {
			if k == val {
				str = strings.Replace(str, "%s", val, 1)
				continue
			}
			expand := wn.recursiveSearch(wordMap, defn.mappings[i], val)
			if len(expand) != 0 {
				str = strings.Replace(str, "%s", expand, 1)
			} else {
				str = strings.Replace(str, "%s", val, 1)
			}
		}
		return str
	}
}

func (wn *WNdict) findDef(ID string) *WNdef {
	defn, ok := wn.IDMappings[ID]
	if ok {
		return defn
	} else {
		return &WNdef{}
	}
}

// for use in recursiveSearch
func (wn *WNdict) findDefArr(k string) []*WNdef {
	defn, ok := wn.definitions[k]
	if ok {
		return defn
	} else {
		return []*WNdef{}
	}
}

// for use in main()
func (wn *WNdict) getDef(k string) string {
	var str string = ""

	val, ok := wn.definitions[k]
	if ok {
		for idx, def := range val {
			str = str + strconv.Itoa(idx+1) + ". " + def.origDef + "\n"
		}
	}

	return str
}

// very slow implementation! BUT TRUTHFUL!
func (wn *WNdict) verify(delNodes []string) bool {

	fmt.Println("verifying...")

	for _, defnArr := range wn.definitions {
		// expands all synsets anyway!
		wn.expandDef(delNodes, defnArr[0].name)
	}

	return true

}
