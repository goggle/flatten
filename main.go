package main

import (
	"fmt"
	"os"
	"path"

	docopt "github.com/docopt/docopt-go"
	"github.com/goggle/flatten/flatten"
	"github.com/goggle/flatten/osabstraction"
)

const version = "0.5"

func main() {
	usage := `flatten.

Usage:
  flatten [SOURCE] [DESTINATION] [--include-source-files] [-s | --simulate-only] [-c | --copy-only] [--verbose]
  flatten -h | --help
  flatten -v

Recursively flatten the directory structure from SOURCE to DESTINATION.

Arguments:
  SOURCE                    Optional source directory (default is current directory).
  DESTINATION               Optional destination directory (default is current directory).

Options:
  -h --help                 Show this screen.
  -v --version              Show version.
  --include-source-files    Include the files which are directly located in the SOURCE directory.
  -s --simulate-only        Do not move or copy any files on the system,
                            just output the expected result.
  -c --copy-only            Do not remove anything from the source directory.
  --verbose                 Explain what is being done.`

	arguments, _ := docopt.Parse(usage, nil, true, "flatten "+version, false)

	var source string
	var destination string

	src := arguments["SOURCE"]
	if src == nil {
		p, err := os.Getwd()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		source = p
	} else {
		source = src.(string)
	}

	dst := arguments["DESTINATION"]
	if dst == nil {
		p, err := os.Getwd()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		destination = p
	} else {
		destination = dst.(string)
	}

	verbose := arguments["--verbose"].(bool)
	if verbose {
		flatten.SetVerbose()
	}

	includeSourceFiles := arguments["--include-source-files"].(bool)

	copyOnly := arguments["--copy-only"].(bool)

	sourceFI := osabstraction.File(path.Clean(source))
	destinationFI := osabstraction.File(path.Clean(destination))
	osWrapper := osabstraction.RealOS{}
	err := flatten.Flatten(sourceFI, destinationFI, osWrapper, copyOnly, includeSourceFiles)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
