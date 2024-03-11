package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Flag struct {
	Name        string
	Type        string
	Description string
}

var (
	flagStart = regexp.MustCompile(`^  -(\S+)\s+(.*)$`)
)

// Parse standard Go --help (when using the flag package).
func parseStandard(in string) []Flag {
	scanner := bufio.NewScanner(strings.NewReader(in))
	var flags []Flag
	var curFlag, curType, curDesc string

	for scanner.Scan() {
		line := scanner.Text()
		match := flagStart.FindStringSubmatch(line)
		if len(match) > 0 {
			if curFlag != "" {
				flags = append(flags, Flag{Name: curFlag, Type: curType, Description: curDesc})
			}
			curFlag = match[1]
			curType = match[2]
			curDesc = ""
		} else if strings.HasPrefix(line, "    ") {
			curDesc += strings.TrimSpace(line[4:])
		}
	}

	if curFlag != "" {
		flags = append(flags, Flag{Name: curFlag, Type: curType, Description: curDesc})
	}

	return flags
}

func generateMD(flags []Flag) string {
	table := "| Flags | Type | Description |\n|-------|------|-------------|\n"
	for _, flag := range flags {
		table += fmt.Sprintf("| `-%s` | %s | %s |\n", flag.Name, flag.Type, flag.Description)
	}
	return table
}

func main() {
	in := bufio.NewReader(os.Stdin)
	bytes, err := io.ReadAll(in)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	flags := parseStandard(string(bytes))
	fmt.Println(generateMD(flags))
}
