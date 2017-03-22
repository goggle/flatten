package flatten

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/goggle/flatten/osabstraction"
)

func countFileNames(files []osabstraction.FileInfo) map[string]int {
	countMap := map[string]int{}
	for _, file := range files {
		countMap[file.Name()]++
	}
	return countMap
}

func evaluateAppendixLength(destination string, filename string, occurences int, osw osabstraction.OSWrapper) int {
	if occurences == 1 {
		if !osw.Exists(filepath.Join(destination, filename)) {
			return 0
		}
	}
	// TODO: Add maximal length of occurences and return error, if this number is exceeded
	minNumberDigits := len(fmt.Sprintf("%v", occurences))
	numberDigits := minNumberDigits
	for {
		works := true
		for i := 1; i <= occurences; i++ {
			fname := generateFilename(filename, i, numberDigits)
			fullpath := filepath.Join(destination, fname)
			if !osw.Exists(fullpath) {
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
	if length == 0 {
		return oldFilename
	}
	lengthString := fmt.Sprintf("%v", length)
	numberExt := fmt.Sprintf("%0"+lengthString+"v", index)
	base := baseName(oldFilename)
	fname := base + "_" + numberExt + path.Ext(oldFilename)
	return fname
}

func baseName(filename string) string {
	j := strings.LastIndex(filename, path.Ext(filename))
	return filename[:j]
}

func Flatten(source, destination osabstraction.FileInfo, osw osabstraction.OSWrapper, copyOnly bool, includeBaseFiles bool) error {
	if !osw.IsDirectory(source.FullPath()) {
		return errors.New(source.FullPath() + " is not a directory")
	}
	if !osw.IsDirectory(destination.FullPath()) {
		return errors.New(destination.FullPath() + " is not a directory")
	}

	files, err := osw.GetFiles(source.FullPath(), includeBaseFiles)
	if err != nil {
		return errors.New("could not retrieve files in " + source.FullPath())
	}
	countMap := countFileNames(files)
	lenAppendixMap := map[string]int{}
	currentIndexMap := map[string]int{}
	for k, v := range countMap {
		l := evaluateAppendixLength(destination.FullPath(), k, v, osw)
		lenAppendixMap[k] = l
		currentIndexMap[k] = 1
	}

	for _, srcFile := range files {
		name := srcFile.Name()
		lenAppendix, _ := lenAppendixMap[name]
		currIndex, _ := currentIndexMap[name]
		currentIndexMap[name]++
		newName := generateFilename(name, currIndex, lenAppendix)
		newNameFullpath := filepath.Join(destination.FullPath(), newName)
		if copyOnly {
			// fmt.Println("Copying " + srcFile.FullPath() + " to " + newNameFullpath)
			err := osw.Copy(srcFile.FullPath(), newNameFullpath)
			if err != nil {
				return err
			}
		} else {
			// fmt.Println("Moving " + srcFile.FullPath() + " to " + newNameFullpath)
			err := osw.Move(srcFile.FullPath(), newNameFullpath)
			if err != nil {
				return err
			}
		}
	}

	if !copyOnly {
		err := osw.RemoveSubDirectories(source.FullPath())
		if err != nil {
			return err
		}
	}
	return nil
}
