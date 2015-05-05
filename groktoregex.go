// groktoregex is a utility that converts logstash grok alias patterns to its regex value.
// Pattern files are in the format of:
// "<GROK_NAME>" = "<VALUE>"
//
// Logstash patterns can be at:
// https://github.com/logstash-plugins/logstash-patterns-core

package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	os.Exit(realmain())
}

func realmain() int {

	c, err := Parser()
	if err != nil {
		return 2
	}

	grokSubstitutePattern, err := regexp.Compile(`\%\{+(\w+)(?:(?:[\,\:\w+]+)\}|\})`)
	if err != nil {
		fmt.Println(err)
		return 2
	}

	grokFileLoadingPattern, err := regexp.Compile(`\"([\w\d]+)\"(?:(?:\s+\=\s+)|(?:\=))\"(.*)\"`)
	if err != nil {
		fmt.Println(err)
		return 2
	}

	g := Groker{
		LookupDict:          make(map[string]string),
		SubstitutionPattern: grokSubstitutePattern,
		FileloadingPattern:  grokFileLoadingPattern,
		Arguments:           c,
	}

	if err = g.Map(); err != nil {
		return 2
	}

	g.Convert()

	return 0
}
