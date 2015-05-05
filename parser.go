package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// CLI arguments, source folder and grok pattern.
type CLI struct {
	LookupFolder string
	GrokPattern  string
}

// Validates CLI arguments given, both grok and source are required.
// Returns a pointer to the CLI struct and error.
func Parser() (*CLI, error) {

	var grok string
	var source string

	flags := flag.NewFlagSet("grok", flag.ExitOnError)
	flags.Usage = func() { helper() }
	flags.StringVar(&grok, "grok", "", "Grok pattern to convert to regex")
	flags.StringVar(&source, "source", "", "Folder location of grok names and patterns")

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
	}

	if len(grok) == 0 {
		fmt.Fprintf(os.Stderr, "pattern folder location required\n")
		return nil, errors.New("Zero length pattern")
	}

	if len(source) == 0 {
		fmt.Fprintf(os.Stderr, "grok pattern required to convert to regex")
		return nil, errors.New("Zero length lookup\n")
	}

	return &CLI{
		LookupFolder: source,
		GrokPattern:  grok,
	}, nil

}

// Creates a file path glob using a wildcard, returns all files found in
// the given source path, returns an array of filepath strings and error
func (c *CLI) grokfiles() ([]string, error) {

	globPath := filepath.Join(c.LookupFolder, "*")

	matches, err := filepath.Glob(globPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return matches, nil
}

// Displays help and usage
func helper() {
	fmt.Fprintf(os.Stderr, message)
}

const message string = `Usage: groktoregex [options]

Options:
  -grok=""            Grok pattern to convert to regex
  -source=  		  Folder location for grok patterns

Examples:
	groktoregex --grok "%{HOSTPORT}" --source patterns/`
