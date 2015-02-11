package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

func colorizeMatch(s, word string) string {
	r, _ := regexp.Compile(word)
	return r.ReplaceAllStringFunc(s, func(ms string) string {
		return colorize(ms, Yellow)
	})
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

	lines := Map(
		strings.Split(string(body), "\n\n"), func(v string) string {
			return setColor(v, fu.search)
		})
	return string(strings.Join(lines, "\n\n")) + colorize("Page: ", Magenta) + fmt.Sprintf("%2d", fu.page)
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
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
	page_idx := (fu.page - 1) * 25
	return fmt.Sprintf(
		"http://www.commandlinefu.com/commands/%s/%s/%d",
		fu.commandPart(), fu.format, page_idx)
}

func help(name string) {
	s := `# Usage

    %s COMMAND [PAGE]

      COMMAND: browse, using WORD, by USER, matching WORD
      PAGE: 1-999 (defaut: 1)

  # Example

    %s browse
    %s using grep
    %s by USER
    %s matching find

  # Abbreviation
    you can abbreviate COMMAND

    %s br
    %s u grep 2
    %s by USER
    %s m find

`

	r, _ := regexp.Compile("%s")
	fmt.Printf(r.ReplaceAllStringFunc(s, func(s string) string {
		return name
	}))
}

var commands = []string{"browse", "using", "by", "matching"}
var abbrevTable map[string]string = NewAbbrev(commands)

func parsePage(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	} else {
		return 1
	}
}

func main() {
	ProgramName := os.Args[0]
	OtherArgs := os.Args[1:]
	if len(OtherArgs) < 1 {
		help(ProgramName)
		os.Exit(0)
	}

	command := OtherArgs[0]
	OtherArgs = OtherArgs[1:] // shift

	if val, ok := abbrevTable[command]; ok {
		command = val
	} else {
		help(ProgramName)
		os.Exit(0)
	}

	search := ""
	page := 1

	switch command {
	case "using", "by", "matching":
		if len(OtherArgs) == 0 {
			help(ProgramName)
			os.Exit(0)
		}
		search = OtherArgs[0]
		OtherArgs = OtherArgs[1:] // shift
	default:
		search = ""
	}

	if len(OtherArgs) > 0 {
		page = parsePage(OtherArgs[0])
	}

	fu := &Fu{
		page:    page,
		format:  "plaintext",
		command: command,
		search:  search,
	}
	fmt.Println(fu.result())
}
