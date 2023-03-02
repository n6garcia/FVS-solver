package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

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
	file, err := os.Open("wrangle/cleaned/" + fn)
	if err != nil {
		fmt.Println("error loading json")
		return
	}
	defer file.Close()

	fmt.Println("file: ", fn)

	scanner := bufio.NewScanner(file)

	var txt string

	for scanner.Scan() {
		line := scanner.Text()
		txt = txt + line
	}

	bytes := []byte(txt)

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

func (d *Dictionary) expandDef(delNodes []string, k string) []string {
	wordMap := make(map[string]bool)
	var defn []string

	if k == "" {
		return defn
	}

	for _, val := range d.definitions {
		wordMap[val.name] = false
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

	return newDefn
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
func (d *Dictionary) getDef(k string) []string {
	if k == "" {
		return nil
	}
	defn, ok := d.definitions[k]
	if ok {
		return defn.words
	} else {
		return nil
	}
}

// very slow implementation, for next verifier use graph and bfs!
// implementation takes about 40m on my computer to run (average hardware)
func (d *Dictionary) verify(delNodes []string) bool {

	fmt.Println("verifying...")

	for _, val := range d.definitions {
		d.expandDef(delNodes, val.name)
	}

	return true

}
