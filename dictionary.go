package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

/* COULD BE OPTIMIZED BUT HAS OK RUNTIME FOR NOW */
// could easily be a map of string to []string, ... literally a dictionary ...
// spaghetti code can be optimized to make recursive definition search much faster
// by using a map DS, though doesn't affect time takes to find vertexCover at all

type Dictionary struct {
	definitions []*Definition
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

// unsafe add, we assume no duplicate entries in dictionary
func (d *Dictionary) addDef(n string, w []string) {
	d.definitions = append(d.definitions, &Definition{name: n, words: w})
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

	for _, val := range d.definitions {
		if val.name == k {
			defn = val.words
			break
		}
	}

	var newDefn []string = []string{}
	fmt.Println(defn)
	for _, val := range defn {
		fmt.Println(val)
		expand := d.recurDef(wordMap, val)
		fmt.Println(expand)
		if len(expand) != 0 {
			newDefn = append(newDefn, expand...)
		} else {
			newDefn = append(newDefn, val)
		}
	}

	return newDefn
}

func (d *Dictionary) recurDef(wordMap map[string]bool, k string) []string {
	val, ok := wordMap[k]
	if !ok || val {
		return []string{}
	} else {
		defn := d.findDef(k)
		var newDefn []string = []string{}
		for _, val := range defn {
			expand := d.recurDef(wordMap, val)
			if len(expand) != 0 {
				newDefn = append(newDefn, expand...)
			} else {
				newDefn = append(newDefn, val)
			}
		}
		return newDefn
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

// for use outside recurDef
func (d *Dictionary) getDef(k string) []string {
	if k == "" {
		return nil
	}
	for _, val := range d.definitions {
		if val.name == k {
			return val.words
		}
	}
	return nil
}
