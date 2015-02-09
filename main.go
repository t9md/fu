package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Color int

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func colorize(s string, color Color) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
}

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

/*
Example
browse/sort-by-votes - All commands sorted by votes
tagged/163/grep - Commands tagged with 'grep', sorted by date (the default sort order)
matching/ssh/c3No - Search results for the query 'ssh' (note that the final segment is a base64-encoding of the search query)

api_comand_set = [ :browse, :tagged, :matching ]
api_format = [ :plaintext, :json, :rss ]
api_url = "http://www.commandlinefu.com/commands/<command-set>/<format>/"
*/

type Fu struct {
	page    int
	format  string
	command string
	search  string
}

func (fu *Fu) result() string {
	resp, err := http.Get(fu.url())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n\n")
	for i, line := range lines {
		lines[i] = setColor(line, fu.search)
	}
	return string(strings.Join(lines, "\n\n"))
}

func colorizeMatch(s, word string) string {
	r, _ := regexp.Compile(word)
	return r.ReplaceAllStringFunc(s, _colorizeMatch)
}

func _colorizeMatch(s string) string {
	return colorize(s, Yellow)
}

func setColor(s, search string) string {
	r := strings.Split(s, "\n")
	for i, line := range r {
		if strings.HasPrefix(line, "# ") {
			r[i] = colorize(line, Green)
		} else {
			r[i] = colorizeMatch(line, search)
		}
	}
	return strings.Join(r, "\n")
}

func (fu *Fu) commandPart() string {
	switch fu.command {
	case "browse":
		return "browse"
	case "using":
		return "using/" + fu.search
	case "by":
		return "by/" + fu.search
	case "matching":
		enc := base64.StdEncoding.EncodeToString([]byte(fu.search))
		return "matching/" + fu.search + "/" + enc
	default:
		return ""
	}
}

func (fu *Fu) url() string {
	// api_url = "http://www.commandlinefu.com/commands/<command-set>/<format>/"
	api_url := "http://www.commandlinefu.com/commands/%s/%s"
	command := fu.commandPart()
	return fmt.Sprintf(api_url, command, fu.format)
}

func NewFu(config interface{}) *Fu {
	return &defaultConfig
}

var defaultConfig = Fu{
	page:    0,
	format:  "plaintext",
	command: "browse",
	search:  "",
}

var p = fmt.Println

func help(name string) {
	fmt.Printf("%s help\n", name)
}

var commands = []string{"browse", "using", "by", "matching"}
var abbrevTable map[string]string = NewAbbrev(commands)

func main() {
	ProgramName := os.Args[0]
	// Args := os.Args[1:]

	// var commands = []string{"browse", "using", "by", "matching"}
	command := os.Args[1]
	search := os.Args[2]
	// command := "using"
	// search := "grep"

	if val, ok := abbrevTable[command]; ok {
		command = val
	} else {
		help(ProgramName)
		os.Exit(0)
	}

	fu := Fu{
		page:    0,
		format:  "plaintext",
		command: command,
		search:  search,
	}
	p(fu.result())
	os.Exit(0)
}
