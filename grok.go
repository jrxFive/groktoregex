package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

const grokNameIndex int = 1
const grokValueIndex int = 2
const grokSyntaxIndex int = 0
const grokInnerSyntaxIndex int = 1

// Holds source map, grok keynames and their associated regex values.
//Substitution pattern is a compiled regex:
//	regexp.Compile(`\%\{+(\w+)(?:(?:[\,\:\w+]+)\}|\})`)
//
//Matches %{....}
//
//Fileloading pattern is used to create a map[string]string from source files
//	regexp.Compile(`\"([\w\d]+)\"(?:(?:\s+\=\s+)|(?:\=))\"(.*)\"`)
//
//Arguments pointer to CLI
type Groker struct {
	LookupDict          map[string]string
	SubstitutionPattern *regexp.Regexp
	FileloadingPattern  *regexp.Regexp
	Arguments           *CLI
}

// Loads and scrubs source files to create a map[string]string assigned to
// Groker LookupDict
func (g *Groker) Map() error {

	files, err := g.Arguments.grokfiles()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, filename := range files {
		fh, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		splitFh := bytes.Split(fh, []byte{'\n'})

		for _, line := range splitFh {
			rgroups := g.FileloadingPattern.FindStringSubmatch(string(line))

			if len(rgroups) > 0 {
				g.LookupDict[rgroups[grokNameIndex]] = rgroups[grokValueIndex]
			}
		}
	}

	return nil
}

// Checks to see if any grok syntax keys are in the current value of the Grok Pattern
func (g *Groker) checkForKeys() bool {
	return g.SubstitutionPattern.MatchString(g.Arguments.GrokPattern)
}

// Returns [][]string of all grok syntax items in current value of Grok Pattern
func (g *Groker) grabKeys() [][]string {
	return g.SubstitutionPattern.FindAllStringSubmatch(g.Arguments.GrokPattern, -1)
}

// Replaces grok syntax with regex value, recursively. Checks to see if any
// subsitituion patterns are found, and replaces the grok syntax with the mapped
// value of that key
func (g *Groker) Convert() {

	if g.checkForKeys() {
		foundKeys := g.grabKeys()

		for _, value := range foundKeys {

			mapValue, ok := g.LookupDict[value[grokInnerSyntaxIndex]]
			if !ok {
				fmt.Sprintln("Could not find key: %s", value[grokInnerSyntaxIndex])
				continue
			}
			g.Arguments.GrokPattern = strings.Replace(g.Arguments.GrokPattern,
				value[grokSyntaxIndex],
				mapValue,
				-1)
		}

		if g.checkForKeys() {
			g.Convert()
		} else {
			fmt.Println(g.Arguments.GrokPattern)
		}
	} else {
		fmt.Println(g.Arguments.GrokPattern)
	}

}
