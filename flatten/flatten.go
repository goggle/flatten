package flatten

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getNestedFileInfo(source string, includeBaseFiles bool) []os.FileInfo {
	files := []os.FileInfo{}
	filepath.Walk(source, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if !includeBaseFiles {
				if path.Dir(p) != path.Clean(source) {
					files = append(files, info)
				}
			} else {
				files = append(files, info)
			}
		}
		return nil
	})
	return files
}

// func generateConflictMap(existingFiles []os.FileInfo, nestedFiles []os.FileInfo) map[string]int {
// 	conflictMap := map[string]int{}
// 	for _, file := range existingFiles {
// 		conflictMap[file.Name()]++
// 	}
// 	for _, file := range nestedFiles {
// 		conflictMap[file.Name()]++
// 	}
// 	return conflictMap
// }

func countFileNames(files []os.FileInfo) map[string]int {
	countMap := map[string]int{}
	for _, file := range files {
		countMap[file.Name()]++
	}
	return countMap
}

func evaluateAppendixLength(destination string, filename string, occurences int) int {
	// TODO: Add maximal length of occurences and return error, if this number is exceeded
	minNumberDigits := len(fmt.Sprintf("%v", occurences))
	numberDigits := minNumberDigits
	for {
		// base := baseName(filename)
		works := true
		for i := 1; i <= occurences; i++ {
			// nStr := fmt.Sprintf("%v", numberDigits)
			// numberExt := fmt.Sprintf("%0"+nStr+"v", i)
			// fname := base + "_" + numberExt + path.Ext(filename)
			// fullpath := filepath.Join(destination, fname)
			fname := generateFilename(filename, i, numberDigits)
			fullpath := filepath.Join(destination, fname)
			if _, err := os.Stat(fullpath); os.IsNotExist(err) {
				continue
			}
			works = false
			break
		}
		if !works {
			numberDigits++
			continue
		}
		break
	}
	return numberDigits
}

func generateFilename(oldFilename string, index int, length int) string {
	indexString := fmt.Sprintf("%v", index)
	numberExt := fmt.Sprintf("%0"+indexString+"v", length)
	base := baseName(oldFilename)
	fname := base + "_" + numberExt + path.Ext(oldFilename)
	return fname
}

func baseName(filename string) string {
	j := strings.LastIndex(filename, path.Ext(filename))
	return filename[:j]
}

func Flatten(source string, destination string, copyOnly bool, includeBaseFiles bool) error {
	if fi, err := os.Stat(source); os.IsNotExist(err) {
		return errors.New(source + " does not exist")
	} else if !fi.IsDir() {
		return errors.New(source + " is not a directory")
	}
	if fi, err := os.Stat(destination); os.IsNotExist(err) {
		return errors.New(destination + " does not exist")
	} else if !fi.IsDir() {
		return errors.New(destination + "is not a directory")
	}

	files := getNestedFileInfo(source, includeBaseFiles)
	countMap := countFileNames(files)
	indexMap := map[string]int{}
	for key, _ := range countMap {
		indexMap[key] = 1
	}
	return nil
}

func copy(source string, destination string) error {
	return nil
}

func move(source string, destination string) error {
	return nil
}
