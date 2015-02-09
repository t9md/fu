package main

import (
	"fmt"
	"os"
)

func NewAbbrev(words []string) map[string]string {
	s := ""
	m := make(map[string]int)
	table := make(map[string]string)
	str := []rune{}

	for _, word := range words {
		str = []rune(word)
		for i := 0; i < len(str); i++ {
			s = string(str[0:i])
			m[s] += 1
			table[s] = word
		}
	}
	for k, v := range m {
		if v > 1 {
			delete(table, k)
		}
	}
	for _, word := range words {
		table[word] = word
	}
	return table
}

func main() {
	words := []string{
		"Hydrogen",
		"Helium",
		"アイランド",
		"アイメイク",
		"ユニコード",
		"ユニコ",
		"Lithium",
		"Beryllium",
		"Boron",
		"Carbon",
		"Nitrogen",
		"Oxygen",
		"Fluorine",
		"Neon",
	}
	table := NewAbbrev(words)
	// fmt.Println(NewAbb([]string{"abc", "def"}))
	// os.Exit(0)
	// fmt.Println(abbrev_table)
	var input string
	for {
		fmt.Printf("input? ")
		fmt.Scanf("%s", &input)

		switch input {
		case "quit":
			os.Exit(0)
		default:
			if val, ok := table[input]; ok {
				fmt.Printf("Found!: %s => %s\n", input, val)
			} else {
				fmt.Printf("NotFound: %s\n", input)
			}
		}
	}
}
