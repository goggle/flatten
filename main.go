package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/goggle/flatten/filesystem"
	"github.com/goggle/flatten/flatten"
	"github.com/goggle/flatten/osabstraction"
)

const version = "0.8.0"

func ask(question string, defaultYes bool) bool {
	var defaultString string
	if defaultYes {
		defaultString = "[y]/n:"
	} else {
		defaultString = "y/[n]:"
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(question + " " + defaultString + " ")
	test, _ := reader.ReadString('\n')
	test = strings.Trim(test, "\n ")
	if defaultYes {
		if test == "n" || test == "N" || test == "no" || test == "No" || test == "NO" {
			return false
		}
		return true
	}
	if test == "y" || test == "Y" || test == "yes" || test == "Yes" || test == "YES" {
		return true
	}
	return false
}

func simulate(sourceFI osabstraction.FileInfo, destinationFI osabstraction.FileInfo, copyOnly bool, includeSourceFiles bool) (string, error) {
	fs := filesystem.Filesystem{}
	fs.Init()
	err := fs.AddFromRealFilesystem(sourceFI.FullPath())
	if err != nil {
		return "", err
	}
	err = fs.AddFromRealFilesystem(destinationFI.FullPath())
	if err != nil {
		return "", err
	}
	err = flatten.Flatten(sourceFI, destinationFI, fs, copyOnly, includeSourceFiles)
	if err != nil {
		return "", err
	}
	tree := filesystem.Tree{}
	tree.Create(destinationFI, fs)
	return fmt.Sprintf("%v", tree), nil
}

func main() {
	usage := `flatten.

Usage:
  flatten [SOURCE] [DESTINATION] [-c | --copy-only] [-f | --force] [--include-source-files] [-s | --simulate-only] [--verbose]
  flatten -h | --help
  flatten -v

Recursively flatten the directory structure from SOURCE to DESTINATION.

Arguments:
  SOURCE                    Optional source directory (default is current directory).
  DESTINATION               Optional destination directory (default is current directory).

Options:
  -c --copy-only            Do not remove anything from the source directory.
  -f --force                Do not propose a simulation first, immediately execute the command.
  --include-source-files    Include the files which are directly located in the SOURCE directory.
  -s --simulate-only        Do not move or copy any files on the system,
                            just output the expected result.
  --verbose                 Explain what is being done.
  -v --version              Show version.
  -h --help                 Show this screen.`

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
	simulateOnly := arguments["--simulate-only"].(bool)
	force := arguments["--force"].(bool)
	performSimulation := false
	askSecondQuestion := true

	sourceFI := osabstraction.File(path.Clean(source))
	destinationFI := osabstraction.File(path.Clean(destination))

	if simulateOnly {
		performSimulation = true
	}

	// Propose to do a simulation first as long we are not in the
	// "force" or "simulation-only" mode:
	if !force && !simulateOnly {
		anw := ask("Flatten performs changes on the file system. Do you want to simulate this process first?", true)
		if anw {
			performSimulation = true
		} else {
			performSimulation = false
			askSecondQuestion = false
		}
	}

	if performSimulation {
		treeString, err := simulate(sourceFI, destinationFI, copyOnly, includeSourceFiles)
		if err != nil {
			fmt.Println("Could not simulate the process. The following error occured:")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf(treeString)
	}

	// If we are in "simulation-only" mode, we can exit the program
	// at this point.
	if simulateOnly {
		os.Exit(0)
	}

	// Ask the user to continue the process as long as we are not in
	// the "force" mode.
	if !force && askSecondQuestion {
		anw := ask("The above changes will be performed. Do you want to continue?", false)
		if !anw {
			os.Exit(0)
		}
	}

	// Perform the flattening process on the real filesystem:
	osWrapper := osabstraction.RealOS{}
	err := flatten.Flatten(sourceFI, destinationFI, osWrapper, copyOnly, includeSourceFiles)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
