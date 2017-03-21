package filesystem

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/goggle/flatten/osabstraction"
)

type Filesystem map[string]DummyFile

func (fs Filesystem) Init() {
	fs["/"] = DummyFile{Path: "/", IsDirectory: true}
}

func (fs Filesystem) MkDir(dir string) error {
	cleanPath := filepath.Clean(dir)
	if strings.Index(cleanPath, "/") != 0 {
		return errors.New("invalid path: " + cleanPath)
	}
	_, alreadyExists := fs[cleanPath]
	if alreadyExists {
		return errors.New(dir + " already exists in filesystem")
	}
	chain := strings.Split(cleanPath, "/")
	currentPath := "/"
	for _, name := range chain {
		if name == "" {
			continue
		}
		currentPath += name
		file, exists := fs[currentPath]
		if exists && !file.IsDir() {
			return errors.New(currentPath + " is a file, not a directory")
		} else if !exists {
			break
		}
		currentPath += "/"
	}
	currentPath = "/"
	for _, name := range chain {
		if name == "" {
			continue
		}
		currentPath += name
		_, exists := fs[currentPath]
		if !exists {
			fs[currentPath] = DummyFile{Path: currentPath, IsDirectory: true}
		}
		currentPath += "/"
	}
	return nil
}

func (fs Filesystem) CreateFile(fpath string) error {
	cleanPath := filepath.Clean(fpath)
	if strings.Index(cleanPath, "/") != 0 {
		return errors.New("invalid path: " + cleanPath)
	}
	file := DummyFile{Path: cleanPath, IsDirectory: false}
	directory := file.Directory()
	entry, exists := fs[directory]
	if !exists {
		err := fs.MkDir(directory)
		if err != nil {
			return err
		}
	} else if !entry.IsDir() {
		return errors.New(directory + " is not a directory")
	}
	_, exists = fs[cleanPath]
	if exists {
		return errors.New(cleanPath + " already exists in file system")
	}
	fs[cleanPath] = file
	return nil
}

func (fs Filesystem) RemoveDirectory(path string) error {
	cleanPath := filepath.Clean(path)
	if cleanPath == "/" {
		return errors.New("cannot remove root directory")
	}
	entry, exists := fs[cleanPath]
	if !exists {
		return errors.New(cleanPath + " does not exist in file system")
	} else if !entry.IsDir() {
		return errors.New(cleanPath + " is not a directory")
	}
	for _, v := range fs {
		if v.Directory() == cleanPath {
			return errors.New(cleanPath + " is not empty!")
		}
	}
	delete(fs, cleanPath)
	return nil
}

func (fs Filesystem) RemoveFile(path string) error {
	cleanPath := filepath.Clean(path)
	entry, exists := fs[cleanPath]
	if !exists {
		return errors.New(cleanPath + " does not exist in file system")
	} else if entry.IsDir() {
		return errors.New(cleanPath + " is a directory")
	}
	delete(fs, cleanPath)
	return nil
}

// interface function
func (fs Filesystem) Copy(source string, destination string) error {
	sourcePath := filepath.Clean(source)
	file, exists := fs[sourcePath]
	if !exists {
		return errors.New(sourcePath + " does not exist in file system")
	}
	if file.IsDir() {
		return fs.MkDir(destination)
	}
	return fs.CreateFile(destination)
}

// interface function
func (fs Filesystem) Move(source string, destination string) error {
	err := fs.Copy(source, destination)
	if err != nil {
		return err
	}
	file, _ := fs[filepath.Clean(source)]
	if file.IsDir() {
		return fs.RemoveDirectory(source)
	} else {
		return fs.RemoveFile(source)
	}
}

func (fs Filesystem) Dirs() []string {
	dirs := []string{}
	for _, v := range fs {
		if v.IsDir() {
			dirs = append(dirs, v.FullPath())
		}
	}
	return dirs
}

func (fs Filesystem) RealFiles() []string {
	files := []string{}
	for _, v := range fs {
		if !v.IsDir() {
			files = append(files, v.FullPath())
		}
	}
	return files
}

// interface function
func (fs Filesystem) GetFiles(dir string, includeBaseFiles bool) ([]osabstraction.FileInfo, error) {
	files := []osabstraction.FileInfo{}
	dir = path.Clean(dir)
	for _, v := range fs {
		if !v.IsDir() {
			if v.Directory() == dir && includeBaseFiles {
				files = append(files, v)
			} else if v.Directory() != dir && strings.HasPrefix(v.FullPath(), dir) {
				files = append(files, v)
			}
		}
	}
	return files, nil
}

// interface function
func (fs Filesystem) Stat(name string) (osabstraction.FileInfo, error) {
	cleanPath := path.Clean(name)
	df, exists := fs[cleanPath]
	if !exists {
		return nil, errors.New(name + " does not exist in file system")
	}
	return df, nil
}

// interface function
func (fs Filesystem) IsNotExists(err error) bool {
	if err != nil {
		return true
	}
	return false
}

// interface function
func (fs Filesystem) IsRegularFile(p string) bool {
	df, exists := fs[path.Clean(p)]
	if exists && !df.IsDir() {
		return true
	}
	return false
}

// interface function
func (fs Filesystem) IsDirectory(p string) bool {
	df, exists := fs[path.Clean(p)]
	if exists && df.IsDir() {
		return true
	}
	return false
}

// interface function
func (fs Filesystem) Exists(p string) bool {
	if fs.IsRegularFile(p) || fs.IsDirectory(p) {
		return true
	}
	return false
}

func (fs Filesystem) Equal(f Filesystem) bool {
	if len(fs) != len(f) {
		return false
	}
	for k, v := range fs {
		vv, ok := f[k]
		if !ok {
			return false
		}
		// NOTE: in our file system (which is only used for simulation purposes),
		// two files are considered as equal, if they have the same path and
		// are of the same type (i.e. either regular file or directory)
		if v.FullPath() != vv.FullPath() || v.IsDir() != vv.IsDir() {
			return false
		}
	}
	return true
}

type DummyFile struct {
	Path        string
	IsDirectory bool
}

// type FileInfo interface {
// 	IsDir() bool
// 	FullPath() string
// 	Name() string
// 	Directory() string
// 	Ext() string
// 	BaseName() string
// 	Level() int
// }

func (df DummyFile) IsDir() bool {
	return df.IsDirectory
}

func (df DummyFile) FullPath() string {
	return filepath.Clean(df.Path)
}

func (df DummyFile) Name() string {
	return filepath.Base(df.Path)
}

func (df DummyFile) Directory() string {
	clean := filepath.Clean(df.Path)
	j := strings.LastIndex(clean, "/")
	dir := clean[:j]
	if dir == "" {
		return "/"
	}
	return dir
}

func (df DummyFile) Ext() string {
	return filepath.Ext(df.Path)
}

func (df DummyFile) BaseName() string {
	filename := df.Name()
	ext := df.Ext()
	j := strings.LastIndex(filename, ext)
	return filename[:j]
}

func (df DummyFile) Level() int {
	clean := filepath.Clean(df.Path)
	if clean == "" || clean == "/" {
		return 0
	}
	return strings.Count(clean, "/")
}
