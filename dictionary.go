package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type Dictionary struct {
	definitions []*Definition
}

type Definition struct {
	name  string
	words []string
}

func (d *Dictionary) addDef(n string, w []string) {
	d.definitions = append(d.definitions, &Definition{name: n, words: w})
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
	// create map
	wordMap := make(map[string]bool)

	// init map with false vals
	for _, val := range d.definitions {
		wordMap[val.name] = false
	}

	// init map with vertex Cover (true vals)
	// WARNING: possibly adds values not in dict.definitions
	for _, val := range delNodes {
		wordMap[val] = true
	}

	var defn []string

	// find definition of testName
	for _, val := range d.definitions {
		if val.name == k {
			defn = val.words
			break
		}
	}

	fmt.Println(defn)

	for i, val := range defn {
		defn = slices.Insert(defn, i, d.recurDef(wordMap, val)...)
	}

	return defn
}

func (d *Dictionary) recurDef(wordMap map[string]bool, k string) []string {
	val, ok := wordMap[k]
	if !ok || val {
		return []string{}
	} else {
		defn := d.findDef(k)
		for i, val := range defn {
			defn = slices.Insert(defn, i, d.recurDef(wordMap, val)...)
		}
		return defn
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
